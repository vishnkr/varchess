package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"vc-server/chess"
	_ "vc-server/chess/variants"
)

const startFEN = "rnbqkbnr/8/8/8/8/8/8/RNBQKBNR w KQkq - 0 1"

// EngineCmd returns the top-level "engine" command group.
func EngineCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "engine",
		Short: "Local engine commands (no infrastructure required)",
	}
	cmd.AddCommand(engineReplCmd())
	return cmd
}

func engineReplCmd() *cobra.Command {
	var variantType string
	var fen string
	var targetChecks int

	cmd := &cobra.Command{
		Use:   "repl",
		Short: "Interactive chess REPL — test variants and move validation locally",
		Long: `Start an interactive REPL to play through a chess variant without needing
Redis, MongoDB, or the WebSocket server.

Commands inside the REPL:
  board          Print the current board
  moves          List all legal moves (numbered)
  move <n>       Play move number <n> from the last "moves" listing
  move <e2e4>    Play a move by algebraic squares (e.g. e2e4, e7e8Q)
  fen            Print the current FEN string
  variant        Print variant type and serialised variant state
  quit / exit    Exit the REPL`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRepl(variantType, fen, targetChecks)
		},
	}

	cmd.Flags().StringVarP(&variantType, "variant", "v", "checkmate",
		`Variant type: checkmate, antichess, ncheck, archerchess, wormhole, poisonedpawn, hiddenqueen, footballchess, spychess`)
	cmd.Flags().StringVarP(&fen, "fen", "f", startFEN, "Starting FEN position")
	cmd.Flags().IntVar(&targetChecks, "target-checks", 3, "Target number of checks (ncheck variant)")
	return cmd
}

func runRepl(variantType, fen string, targetChecks int) error {
	customData := map[string]interface{}{}
	if variantType == "ncheck" {
		customData["targetChecks"] = targetChecks
	}

	eng, err := chess.NewEngine(chess.GameConfig{
		VariantType: variantType,
		FEN:         fen,
		CustomData:  customData,
	})
	if err != nil {
		return fmt.Errorf("failed to create engine: %w", err)
	}

	fmt.Printf("VarChess Engine REPL — variant: %q\n", variantType)
	fmt.Println(`Commands: board, moves, move <n|e2e4>, fen, variant, quit`)
	fmt.Println()
	printBoard(eng)

	scanner := bufio.NewScanner(os.Stdin)
	var lastMoves []chess.Move

	for {
		over, result := eng.IsGameOver()
		if over {
			fmt.Printf("\nGame over — %s wins (%s)\n", result.Winner, result.Reason)
			return nil
		}

		turn := "White"
		if !eng.GetPosition().WhiteToMove {
			turn = "Black"
		}
		fmt.Printf("\n%s to move > ", turn)

		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		switch parts[0] {
		case "quit", "exit", "q":
			fmt.Println("Bye!")
			return nil

		case "board":
			printBoard(eng)

		case "moves":
			lastMoves = eng.GetLegalMoves()
			if len(lastMoves) == 0 {
				fmt.Println("No legal moves.")
			} else {
				fmt.Println(chess.FormatLegalMoves(lastMoves, eng.GetPosition()))
			}

		case "move":
			if len(parts) < 2 {
				fmt.Println("Usage: move <n> or move <e2e4>")
				continue
			}
			token := parts[1]
			if n, err := strconv.Atoi(token); err == nil {
				// Move by index.
				if lastMoves == nil {
					lastMoves = eng.GetLegalMoves()
				}
				if n < 1 || n > len(lastMoves) {
					fmt.Printf("Move index out of range (1..%d)\n", len(lastMoves))
					continue
				}
				move := lastMoves[n-1]
				result, err := eng.PerformMove(move)
				if err != nil {
					fmt.Println("Error:", err)
					continue
				}
				lastMoves = nil
				printMoveResult(move, result, eng)
			} else {
				// Move by algebraic string.
				pos := eng.GetPosition()
				from, to, promo, err := chess.ParseAlgebraicMove(token, pos)
				if err != nil {
					fmt.Println("Parse error:", err)
					continue
				}
				legalMoves := eng.GetLegalMoves()
				lastMoves = legalMoves
				matched, ok := chess.MatchMove(from, to, promo, legalMoves)
				if !ok {
					fmt.Printf("Move %q not found in legal moves. Try 'moves' to see options.\n", token)
					continue
				}
				result, err := eng.PerformMove(matched)
				if err != nil {
					fmt.Println("Error:", err)
					continue
				}
				lastMoves = nil
				printMoveResult(matched, result, eng)
			}

		case "fen":
			fmt.Println(eng.GetPosition().FEN)

		case "variant":
			// Print variant info.
			fmt.Printf("Variant: %q\n", variantType)
			// Try to access VariantState if the engine exposes it.
			// StandardEngine embeds VariantRules; we access it indirectly by
			// checking if the engine can provide variant state via a type assertion.
			type variantStater interface {
				VariantState() interface{}
			}
			if se, ok := eng.(variantStater); ok {
				state := se.VariantState()
				if state != nil {
					b, _ := json.MarshalIndent(state, "", "  ")
					fmt.Println("State:", string(b))
				} else {
					fmt.Println("State: (none)")
				}
			}

		default:
			// Treat unknown input as an algebraic move attempt.
			pos := eng.GetPosition()
			from, to, promo, err := chess.ParseAlgebraicMove(parts[0], pos)
			if err != nil {
				fmt.Printf("Unknown command %q. Type 'moves' for legal moves.\n", parts[0])
				continue
			}
			legalMoves := eng.GetLegalMoves()
			lastMoves = legalMoves
			matched, ok := chess.MatchMove(from, to, promo, legalMoves)
			if !ok {
				fmt.Printf("Move %q not legal. Type 'moves' to see options.\n", parts[0])
				continue
			}
			result, err := eng.PerformMove(matched)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			lastMoves = nil
			printMoveResult(matched, result, eng)
		}
	}
	return nil
}

func printBoard(eng chess.Engine) {
	pos := eng.GetPosition()
	fmt.Print(chess.PositionToASCII(pos))

	turn := "White"
	if !pos.WhiteToMove {
		turn = "Black"
	}
	fmt.Printf("Turn: %s", turn)

	// Check indicator — type-assert to StandardEngine.
	type checker interface{ IsCheck() bool }
	if se, ok := eng.(checker); ok && se.IsCheck() {
		fmt.Print("  [CHECK]")
	}
	fmt.Println()
}

func printMoveResult(move chess.Move, result chess.MoveResult, eng chess.Engine) {
	pos := eng.GetPosition()
	label := chess.FormatMove(move, pos)
	status := ""
	if result.IsCapture {
		status += " x"
	}
	if result.IsCheck {
		status += " +"
	}
	if result.IsGameOver {
		status += " #"
	}
	fmt.Printf("  Played: %s%s\n", label, status)
	printBoard(eng)
	if result.IsGameOver && result.Result != nil {
		fmt.Printf("\nGame over — %s wins (%s)\n", result.Result.Winner, result.Result.Reason)
	}
}
