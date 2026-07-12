package commands

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"

	"vc-server/chess"
	_ "vc-server/chess/variants"
)

// PlayCmd returns the top-level "play" command group.
func PlayCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "play",
		Short: "Multiplayer commands — login, create game, join via WebSocket",
	}
	cmd.AddCommand(loginCmd())
	cmd.AddCommand(createCmd())
	cmd.AddCommand(joinCmd())
	return cmd
}

// ---------------------------------------------------------------------------
// login
// ---------------------------------------------------------------------------

func loginCmd() *cobra.Command {
	var coreURL, username, password, tokenOut string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate and print (or save) an access token",
		Example: `  vc-cli play login --user alice --pass secret
  TOKEN=$(vc-cli play login --user alice --pass secret --token-out -)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogin(coreURL, username, password, tokenOut)
		},
	}
	cmd.Flags().StringVar(&coreURL, "core", "http://localhost:5000", "vc-core base URL")
	cmd.Flags().StringVar(&username, "user", "", "Username")
	cmd.Flags().StringVar(&password, "pass", "", "Password")
	cmd.Flags().StringVar(&tokenOut, "token-out", "", `Write token to this path, or "-" for stdout only`)
	_ = cmd.MarkFlagRequired("user")
	_ = cmd.MarkFlagRequired("pass")
	return cmd
}

func runLogin(coreURL, username, password, tokenOut string) error {
	body, _ := json.Marshal(map[string]string{"username": username, "password": password})
	resp, err := http.Post(coreURL+"/login", "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("login failed (%d): %s", resp.StatusCode, string(b))
	}

	var result struct {
		AccessToken string `json:"accessToken"`
		UID         string `json:"uid"`
		Username    string `json:"username"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode login response: %w", err)
	}

	if tokenOut == "-" {
		// stdout-only mode: print just the token so callers can capture it.
		fmt.Println(result.AccessToken)
		return nil
	}

	fmt.Printf("Logged in as %s (uid: %s)\n", result.Username, result.UID)
	fmt.Printf("Token: %s\n", result.AccessToken)

	if tokenOut != "" {
		if err := os.WriteFile(tokenOut, []byte(result.AccessToken), 0600); err != nil {
			return fmt.Errorf("write token: %w", err)
		}
		fmt.Printf("Token saved to %s\n", tokenOut)
	}
	return nil
}

// ---------------------------------------------------------------------------
// create
// ---------------------------------------------------------------------------

func createCmd() *cobra.Command {
	var coreURL, token, variantType, fen, templateID string
	var targetChecks int

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new game and print its game ID",
		Example: `  GAME=$(vc-cli play create --variant antichess --token $TOKEN)
  vc-cli play create --template <mongoID> --token $TOKEN`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreate(coreURL, token, templateID, variantType, fen, targetChecks)
		},
	}
	cmd.Flags().StringVar(&coreURL, "core", "http://localhost:5000", "vc-core base URL")
	cmd.Flags().StringVar(&token, "token", "", "JWT access token (required)")
	cmd.Flags().StringVar(&variantType, "variant", "checkmate", "Variant type")
	cmd.Flags().StringVar(&fen, "fen", startFEN, "Starting FEN")
	cmd.Flags().StringVar(&templateID, "template", "", "MongoDB template ID (overrides --variant/--fen)")
	cmd.Flags().IntVar(&targetChecks, "target-checks", 3, "Target checks for ncheck variant")
	_ = cmd.MarkFlagRequired("token")
	return cmd
}

func runCreate(coreURL, token, templateID, variantType, fen string, targetChecks int) error {
	var reqBody map[string]interface{}
	if templateID != "" {
		reqBody = map[string]interface{}{"templateId": templateID}
	} else {
		customData := map[string]interface{}{}
		if variantType == "ncheck" {
			customData["targetChecks"] = targetChecks
		}
		reqBody = map[string]interface{}{
			"gc": map[string]interface{}{
				"variantType": variantType,
				"fen":         fen,
				"customData":  customData,
			},
		}
	}

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, coreURL+"/api/games", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("create game request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("create failed (%d): %s", resp.StatusCode, string(b))
	}

	var result struct {
		GameID string `json:"gameId"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode create response: %w", err)
	}

	// Print just the game ID (easy to capture with $(...)).
	fmt.Println(result.GameID)
	return nil
}

// ---------------------------------------------------------------------------
// join
// ---------------------------------------------------------------------------

// wsIncomingMsg is the JSON envelope sent from vc-ws to the client.
type wsIncomingMsg struct {
	Type    string          `json:"t"`
	Payload json.RawMessage `json:"p"`
}

// wsOutgoingMsg is the JSON envelope sent from the client to vc-ws.
type wsOutgoingMsg struct {
	Type    string      `json:"t"`
	Payload interface{} `json:"p"`
}

// activeGameSnapshot mirrors the subset of models.ActiveGame we need.
type activeGameSnapshot struct {
	Config chess.GameConfig `json:"gameConfig"`
	Moves  []string         `json:"moves"`
	Turn   string           `json:"turn"`
	State  string           `json:"state"`
}

func joinCmd() *cobra.Command {
	var wsURL, token, color string

	cmd := &cobra.Command{
		Use:   "join <gameId>",
		Short: "Join a game via WebSocket and play interactively",
		Args:  cobra.ExactArgs(1),
		Example: `  vc-cli play join ABC12345 --color w --token $TOKEN
  vc-cli play join ABC12345 --color b --token $TOKEN --ws ws://localhost:8080`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runJoin(wsURL, token, args[0], color)
		},
	}
	cmd.Flags().StringVar(&wsURL, "ws", "ws://localhost:8080", "vc-ws WebSocket base URL")
	cmd.Flags().StringVar(&token, "token", "", "JWT access token (required)")
	cmd.Flags().StringVar(&color, "color", "w", `Preferred color: "w" or "b"`)
	_ = cmd.MarkFlagRequired("token")
	return cmd
}

func runJoin(wsURL, token, gameID, colorPref string) error {
	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = 10 * time.Second

	endpoint := fmt.Sprintf("%s/play/%s", strings.TrimRight(wsURL, "/"), gameID)
	fmt.Printf("Connecting to %s …\n", endpoint)

	conn, _, err := dialer.Dial(endpoint, nil)
	if err != nil {
		return fmt.Errorf("WebSocket connect failed: %w", err)
	}
	defer conn.Close()

	// 1. Send auth + color preference.
	authMsg, _ := json.Marshal(map[string]string{
		"token":     token,
		"colorPref": colorPref,
	})
	if err := conn.WriteMessage(websocket.TextMessage, authMsg); err != nil {
		return fmt.Errorf("send auth: %w", err)
	}

	fmt.Println("Waiting for game to start (both players must join) …")

	// 2. Wait for the "start" event.
	var eng chess.Engine
	var myColor string

	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("read from server: %w", err)
		}
		var msg wsIncomingMsg
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			continue
		}
		if msg.Type == "start" {
			var snap activeGameSnapshot
			if err := json.Unmarshal(msg.Payload, &snap); err != nil {
				return fmt.Errorf("decode start payload: %w", err)
			}
			eng, err = chess.NewEngine(snap.Config)
			if err != nil {
				return fmt.Errorf("create engine: %w", err)
			}
			// Replay any moves that happened before we joined.
			for _, moveStr := range snap.Moves {
				var mj chess.MoveJSON
				if err := json.Unmarshal([]byte(moveStr), &mj); err != nil {
					continue
				}
				m, err := mj.ToMove()
				if err != nil {
					continue
				}
				if _, err := eng.PerformMove(m); err != nil {
					// Best-effort replay; engine might reject partial states.
					_ = err
				}
			}
			myColor = colorPref
			fmt.Printf("Game started! You are playing as %q\n\n", myColor)
			printBoard(eng)
			break
		}
		fmt.Printf("[server] %s: %s\n", msg.Type, string(msg.Payload))
	}

	// 3. Main game loop.
	incoming := make(chan wsIncomingMsg, 8)
	go func() {
		for {
			_, msgBytes, err := conn.ReadMessage()
			if err != nil {
				close(incoming)
				return
			}
			var msg wsIncomingMsg
			if err := json.Unmarshal(msgBytes, &msg); err == nil {
				incoming <- msg
			}
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	var lastMoves []chess.Move

	printPrompt := func() {
		over, _ := eng.IsGameOver()
		if over {
			return
		}
		pos := eng.GetPosition()
		turn := "white"
		if !pos.WhiteToMove {
			turn = "black"
		}
		if turn == myColor || (myColor == "w" && turn == "white") || (myColor == "b" && turn == "black") {
			fmt.Printf("\nYour turn (%s) > ", turn)
		} else {
			fmt.Printf("\nWaiting for opponent (%s) …\n", turn)
		}
	}

	isMyTurn := func() bool {
		pos := eng.GetPosition()
		if myColor == "w" {
			return pos.WhiteToMove
		}
		return !pos.WhiteToMove
	}

	printPrompt()

	inputCh := make(chan string, 4)
	go func() {
		for scanner.Scan() {
			inputCh <- strings.TrimSpace(scanner.Text())
		}
		close(inputCh)
	}()

	for {
		over, result := eng.IsGameOver()
		if over {
			fmt.Printf("\nGame over — %s wins (%s)\n", result.Winner, result.Reason)
			return nil
		}

		select {
		case msg, ok := <-incoming:
			if !ok {
				fmt.Println("Connection closed by server.")
				return nil
			}
			handleServerMsg(msg, eng, &lastMoves)
			printPrompt()

		case line, ok := <-inputCh:
			if !ok {
				return nil
			}
			if line == "" {
				continue
			}
			if !isMyTurn() {
				fmt.Println("Not your turn yet.")
				continue
			}
			if err := handleUserInput(line, conn, eng, &lastMoves); err != nil {
				fmt.Println("Error:", err)
			}
			printPrompt()
		}
	}
}

func handleServerMsg(msg wsIncomingMsg, eng chess.Engine, lastMoves *[]chess.Move) {
	switch msg.Type {
	case "move":
		// Server confirmed a move — update local engine.
		var payload struct {
			M chess.MoveJSON `json:"m"`
		}
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			fmt.Println("[server] invalid move payload:", err)
			return
		}
		m, err := payload.M.ToMove()
		if err != nil {
			fmt.Println("[server] invalid move:", err)
			return
		}
		if _, err := eng.PerformMove(m); err != nil {
			// Server may send confirmed moves that the local engine already applied.
			_ = err
		}
		*lastMoves = nil
		fmt.Printf("\n[opponent] played %s\n", chess.FormatMove(m, eng.GetPosition()))
		printBoard(eng)

	case "over":
		var result struct {
			Winner string `json:"winner"`
			Reason string `json:"reason"`
		}
		if err := json.Unmarshal(msg.Payload, &result); err == nil {
			fmt.Printf("\nGame over — %s wins (%s)\n", result.Winner, result.Reason)
		}

	case "error":
		fmt.Printf("[server error] %s\n", string(msg.Payload))

	default:
		fmt.Printf("[server] %s: %s\n", msg.Type, string(msg.Payload))
	}
}

func handleUserInput(line string, conn *websocket.Conn, eng chess.Engine, lastMoves *[]chess.Move) error {
	parts := strings.Fields(line)
	switch parts[0] {
	case "board":
		printBoard(eng)

	case "moves":
		*lastMoves = eng.GetLegalMoves()
		fmt.Println(chess.FormatLegalMoves(*lastMoves, eng.GetPosition()))

	case "resign":
		out, _ := json.Marshal(wsOutgoingMsg{Type: "resign"})
		return conn.WriteMessage(websocket.TextMessage, out)

	case "quit", "exit":
		conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		os.Exit(0)

	default:
		return sendMoveFromInput(line, conn, eng, lastMoves)
	}
	return nil
}

func sendMoveFromInput(input string, conn *websocket.Conn, eng chess.Engine, lastMoves *[]chess.Move) error {
	var matched chess.Move
	var ok bool

	if *lastMoves == nil {
		*lastMoves = eng.GetLegalMoves()
	}

	// Try numeric index first.
	pos := eng.GetPosition()
	from, to, promo, err := chess.ParseAlgebraicMove(input, pos)
	if err != nil {
		fmt.Printf("Invalid move %q. Use algebraic notation (e.g. e2e4) or 'moves' to list options.\n", input)
		return nil
	}
	matched, ok = chess.MatchMove(from, to, promo, *lastMoves)
	if !ok {
		fmt.Printf("Move %q not legal. Try 'moves' to see options.\n", input)
		return nil
	}

	// Send to server.
	moveJSON := matched.ToJSON()
	payload := map[string]interface{}{"m": moveJSON}
	out, _ := json.Marshal(wsOutgoingMsg{Type: "move", Payload: payload})
	if err := conn.WriteMessage(websocket.TextMessage, out); err != nil {
		return fmt.Errorf("send move: %w", err)
	}

	// Optimistically apply to local engine.
	if _, err := eng.PerformMove(matched); err != nil {
		fmt.Println("Local engine rejected move (server will be authoritative):", err)
	} else {
		*lastMoves = nil
		printMoveResult(matched, chess.MoveResult{IsCapture: matched.Capture}, eng)
	}
	return nil
}
