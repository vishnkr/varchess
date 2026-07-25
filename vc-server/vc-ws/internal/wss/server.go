package wss

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"vc-server/vc-ws/internal/config"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

const SystemAction = "system_action"
const UserAction = "user_action"

type AuthResponse struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
}

func authenticateToken(cfg *config.Config, token string) (*AuthResponse, error) {
	url := fmt.Sprintf("http://%s:%s/api/auth/validate", cfg.ServerHost, cfg.ServerPort)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid token")
	}
	defer resp.Body.Close()

	var userData AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
		return nil, errors.New("invalid response from core")
	}
	return &userData, nil
}

// HandleWSConnection upgrades an HTTP connection to WebSocket and starts auth.
func HandleWSConnection(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	gameID := strings.TrimPrefix(r.URL.Path, "/play/")
	if gameID == "" {
		http.Error(w, "Invalid game ID", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade failed:", err)
		return
	}

	player := &Client{
		conn:          conn,
		gameID:        gameID,
		authenticated: false,
		send:          make(chan []byte, 256),
	}

	go player.handleAuthentication(cfg)
}

func (c *Client) handleAuthentication(cfg *config.Config) {
	_, msg, err := c.conn.ReadMessage()
	if err != nil {
		fmt.Println("Error reading initial message:", err)
		c.conn.Close()
		return
	}

	var connectPayload struct {
		Token     string `json:"token"`
		ColorPref string `json:"colorPref,omitempty"`
	}
	if err := json.Unmarshal(msg, &connectPayload); err != nil {
		fmt.Println("Invalid connection payload:", err)
		c.conn.Close()
		return
	}

	user, err := authenticateToken(cfg, connectPayload.Token)
	if err != nil {
		fmt.Println("Unauthorized access:", err)
		c.conn.Close()
		return
	}

	c.userId = user.UserID
	c.username = user.Username
	c.authenticated = true
	if connectPayload.ColorPref != "" {
		c.colorPref = connectPayload.ColorPref
	}

	c.joinGame()
}

func (c *Client) joinGame() {
	if hub == nil {
		log.Println("joinGame: hub not initialized (Redis unavailable)")
		c.conn.Close()
		return
	}
	game, err := hub.GetOrCreateGame(c.gameID)
	if err != nil {
		log.Printf("joinGame: %v", err)
		c.conn.Close()
		return
	}

	// Reconnect: same user already holds a slot (including during grace window).
	game.mu.Lock()
	var old *Client
	if existing, ok := game.Clients[c.userId]; ok {
		old = existing
		if existing.colorPref != "" && c.colorPref == "" {
			c.colorPref = existing.colorPref
		}
		game.Clients[c.userId] = c
		c.game = game
		game.mu.Unlock()

		// Close the superseded socket so its disconnect handler cannot
		// broadcast offline after we've already marked this user online.
		if old != nil && old != c {
			_ = old.conn.Close()
		}

		go c.listenForWrites()
		go c.listenForMessages()
		c.sendActiveGameSnapshot()
		hub.BroadcastPresence(c.gameID, c.userId, true)
		c.sendPresenceSnapshot()
		log.Printf("Player %s reconnected to game %s", c.userId, c.gameID)
		return
	}
	game.mu.Unlock()

	assignedColor := getPlayerColor(game, c.colorPref)
	if assignedColor == "" {
		log.Printf("Game %s is full, rejecting connection", c.gameID)
		c.conn.Close()
		return
	}
	c.colorPref = assignedColor

	game.mu.Lock()
	game.Clients[c.userId] = c
	game.mu.Unlock()
	c.game = game

	// Publish Join event to vc-core for game state update.
	eventData, _ := json.Marshal(JoinPayload{Color: assignedColor})
	joinEvent := Event{
		GameID: c.gameID,
		UserID: c.userId,
		Type:   Join,
		Data:   eventData,
	}
	joinJSON, err := json.Marshal(joinEvent)
	if err != nil {
		fmt.Println("Error marshalling join event:", err)
		c.conn.Close()
		return
	}
	if err := redisClient.Publish(context.Background(), UserAction, joinJSON).Err(); err != nil {
		fmt.Println("Error publishing join event:", err)
		return
	}

	go c.listenForMessages()
	go c.listenForWrites()
	hub.BroadcastPresence(c.gameID, c.userId, true)
	c.sendPresenceSnapshot()
}

// sendActiveGameSnapshot pushes the current Redis ActiveGame as a start event
// so a reconnecting client can rehydrate board/players without waiting for moves.
func (c *Client) sendActiveGameSnapshot() {
	gameKey := fmt.Sprintf("game:%s", c.gameID)
	gameStateJSON, err := redisClient.Get(context.Background(), gameKey).Result()
	if err != nil {
		log.Printf("sendActiveGameSnapshot: %v", err)
		return
	}
	msg, err := json.Marshal(WSMessage{Type: string(Start), Payload: json.RawMessage(gameStateJSON)})
	if err != nil {
		return
	}
	select {
	case c.send <- msg:
	default:
		log.Printf("send buffer full for reconnecting client %s", c.userId)
	}
}

// sendPresenceSnapshot tells this client who is currently connected in the game.
func (c *Client) sendPresenceSnapshot() {
	if hub == nil || c.game == nil {
		return
	}
	c.game.mu.RLock()
	defer c.game.mu.RUnlock()
	for uid := range c.game.Clients {
		payload, _ := json.Marshal(PresencePayload{UserID: uid, Online: true})
		msg, err := json.Marshal(WSMessage{Type: string(Presence), Payload: json.RawMessage(payload)})
		if err != nil {
			continue
		}
		select {
		case c.send <- msg:
		default:
		}
	}
}

// getPlayerColor assigns a color to a joining player.
// Colors are standardised to "w" and "b".
func getPlayerColor(game *Game, colorPref string) string {
	game.mu.RLock()
	defer game.mu.RUnlock()

	hasW := false
	hasB := false
	for _, client := range game.Clients {
		_ = client
	}
	// Count by peeking at how many clients already have an assigned color.
	// Since we assign by color pref or opposite of existing, count differently.
	clientCount := len(game.Clients)

	if clientCount == 0 {
		if colorPref == "b" {
			return "b"
		}
		return "w"
	}
	if clientCount == 1 {
		// Check which color is taken by looking at what was published via Join events.
		// Simple heuristic: if pref matches available, use it; else use the other.
		for _, client := range game.Clients {
			if client.colorPref == "w" || (!hasB && client.colorPref == "") {
				hasW = true
			} else {
				hasB = true
			}
		}
		if !hasW {
			return "w"
		}
		if !hasB {
			return "b"
		}
	}
	return ""
}
