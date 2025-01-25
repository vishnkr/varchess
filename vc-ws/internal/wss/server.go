package wss

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var (
    upgrader  = websocket.Upgrader{
        CheckOrigin: func(r *http.Request) bool { return true },
    }
)

func HandleWSConnection(w http.ResponseWriter, r *http.Request) {
    gameID := strings.TrimPrefix(r.URL.Path, "/ws/play/")
    if gameID == "" {
        http.Error(w, "Invalid game ID", http.StatusBadRequest)
        return
    }

    userID := r.Header.Get("X-User-ID")
    if userID == "" {
        http.Error(w, "Missing user ID", http.StatusUnauthorized)
        return
    }

    colorPref := r.URL.Query().Get("c")

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Println("WebSocket upgrade failed:", err)
        return
    }

    game := hub.GetOrCreateGame(gameID)
    game.mu.Lock()
    defer game.mu.Unlock()

    assignedColor := getPlayerColor(game, colorPref)
    if assignedColor == "" {
        conn.Close()
        return
    }
    
    player := &Client{
        conn:  conn,
        game:  game,
        userId: userID,
        send:  make(chan []byte, 256),
    }
    game.players[assignedColor] = player

    go player.listenForMessages()
}

func getPlayerColor(game *Game, colorPref string) string {
    if len(game.players) == 0 {
        if colorPref == "b" {
            return "b"
        }
        return "w"
    } else if len(game.players) == 1 {
        if _, exists := game.players["w"]; exists {
            return "b"
        }
        return "w"
    }
    return ""
}

func (c *Client) listenForMessages() {
    defer c.conn.Close()
    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            fmt.Println("WebSocket read error:", err)
            return
        }
        processMessage(c, message)
    }
}

func processMessage(c *Client, msg []byte) {
    // TODO: Validate move with core server via Redis
}