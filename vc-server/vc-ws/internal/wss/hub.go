package wss

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"vc-server/chess"
	_ "vc-server/chess/variants"

	"github.com/redis/go-redis/v9"
)

// activeGameSnapshot is a minimal mirror of models.ActiveGame for the WS hub.
// It avoids importing vc-core/internal, which is forbidden across module boundaries.
type activeGameSnapshot struct {
	Config chess.GameConfig `json:"gameConfig"`
	Moves  []string         `json:"moves"`
	State  string           `json:"state"`
	Turn   string           `json:"turn"`
}

// Hub manages all active WebSocket games and routes Redis events to clients.
type Hub struct {
	mu    sync.RWMutex
	Games map[string]*Game
	redis *redis.Client
}

// Game holds the in-memory state for one live game.
type Game struct {
	mu      sync.RWMutex
	ID      string
	Clients map[string]*Client // userID → Client
	state   chess.Engine
}

func NewHub(r *redis.Client) *Hub {
	return &Hub{
		Games: make(map[string]*Game),
		redis: r,
	}
}

// GetOrCreateGame returns an existing in-memory Game or creates one by loading
// the GameConfig from Redis (key: "game:{gameID}").
func (h *Hub) GetOrCreateGame(gameID string) (*Game, error) {
	h.mu.RLock()
	if g, ok := h.Games[gameID]; ok {
		h.mu.RUnlock()
		return g, nil
	}
	h.mu.RUnlock()

	h.mu.Lock()
	defer h.mu.Unlock()

	// Double-check after acquiring write lock.
	if g, ok := h.Games[gameID]; ok {
		return g, nil
	}

	// Load ActiveGame from Redis — key standardised to "game:{id}".
	gameKey := fmt.Sprintf("game:%s", gameID)
	gameStateJSON, err := h.redis.Get(context.Background(), gameKey).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("game not found: %s", gameID)
	} else if err != nil {
		return nil, fmt.Errorf("redis error: %w", err)
	}

	var activeGame activeGameSnapshot
	if err := json.Unmarshal([]byte(gameStateJSON), &activeGame); err != nil {
		return nil, fmt.Errorf("invalid game state: %w", err)
	}

	engine, err := chess.NewEngine(activeGame.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to create engine: %w", err)
	}

	// Replay recorded moves to restore position.
	for _, moveStr := range activeGame.Moves {
		var mj chess.MoveJSON
		if err := json.Unmarshal([]byte(moveStr), &mj); err != nil {
			log.Printf("skipping invalid recorded move in game %s: %v", gameID, err)
			continue
		}
		m, err := mj.ToMove()
		if err != nil {
			log.Printf("skipping bad recorded move in game %s: %v", gameID, err)
			continue
		}
		if _, err := engine.PerformMove(m); err != nil {
			log.Printf("replay error in game %s: %v", gameID, err)
		}
	}

	game := &Game{
		ID:      gameID,
		Clients: make(map[string]*Client),
		state:   engine,
	}
	h.Games[gameID] = game
	return game, nil
}

// RegisterClient adds a client to the game's client map.
func (h *Hub) RegisterClient(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if g, ok := h.Games[c.gameID]; ok {
		g.mu.Lock()
		g.Clients[c.userId] = c
		g.mu.Unlock()
	}
}

// UnregisterClient removes a client and cleans up empty games.
func (h *Hub) UnregisterClient(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	g, ok := h.Games[c.gameID]
	if !ok {
		return
	}
	g.mu.Lock()
	delete(g.Clients, c.userId)
	remaining := len(g.Clients)
	g.mu.Unlock()

	if remaining == 0 {
		delete(h.Games, c.gameID)
	}
}

// BroadcastToGame sends a message to all clients in a game.
func (h *Hub) BroadcastToGame(gameID string, message []byte) {
	h.mu.RLock()
	g, ok := h.Games[gameID]
	h.mu.RUnlock()
	if !ok {
		return
	}
	g.mu.RLock()
	defer g.mu.RUnlock()
	for _, c := range g.Clients {
		select {
		case c.send <- message:
		default:
			log.Printf("send buffer full for client %s", c.userId)
		}
	}
}

// HandleSystemEvent processes an event arriving from the Redis system_action channel.
func (h *Hub) HandleSystemEvent(event Event) {
	switch event.Type {
	case Start:
		h.handleStartEvent(event)
	case Move:
		h.handleMoveEvent(event)
	case GameOver:
		h.handleGameOverEvent(event)
	default:
		h.BroadcastToGame(event.GameID, marshalEvent(event))
	}
}

func (h *Hub) handleStartEvent(event Event) {
	_, err := h.GetOrCreateGame(event.GameID)
	if err != nil {
		log.Printf("handleStartEvent: %v", err)
		return
	}
	h.BroadcastToGame(event.GameID, marshalEvent(event))
}

func (h *Hub) handleMoveEvent(event Event) {
	h.BroadcastToGame(event.GameID, marshalEvent(event))
}

func (h *Hub) handleGameOverEvent(event Event) {
	h.BroadcastToGame(event.GameID, marshalEvent(event))
	h.mu.Lock()
	delete(h.Games, event.GameID)
	h.mu.Unlock()
}

func marshalEvent(event Event) []byte {
	b, err := json.Marshal(WSMessage{Type: string(event.Type), Payload: event.Data})
	if err != nil {
		return nil
	}
	return b
}
