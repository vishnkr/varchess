package wss

import (
	"encoding/json"
	"errors"
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

type AuthResponse struct {
    UserID  string `json:"user_id"`
    Username string `json:"username"`
}

func authenticateToken(token string) (*AuthResponse, error) {
    
    url := "http://localhost:5000/api/auth/validate"
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


func HandleWSConnection(w http.ResponseWriter, r *http.Request) {
    gameID := strings.TrimPrefix(r.URL.Path, "/play/")
    if gameID == "" {
        http.Error(w, "Invalid game ID", http.StatusBadRequest)
        return
    }

    token := r.Header.Get("X-Auth-Token")
    if token == "" {
        http.Error(w, "Missing token", http.StatusUnauthorized)
        return
    }

    user,err := authenticateToken(token)
    if err!=nil{
        http.Error(w, "Missing token", http.StatusUnauthorized)
        return
    }
    userID := user.UserID
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
    if existingPlayer, exists := game.players[assignedColor]; exists && existingPlayer.userId == userID {
        // reconnect
        existingPlayer.conn.Close()
        delete(game.players, assignedColor)
    }
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
    var eventData JoinPayload = JoinPayload{Color: assignedColor}
    
    joinEvent := Event{
        GameID: gameID,
        UserID: userID,
        Type: Join,
        Data: eventData,
    }
    data,err := json.Marshal(joinEvent)
    if err != nil {
        fmt.Println("Error marshalling join event:", err)
        conn.Close()
        return
    }
    err = redisClient.Publish(r.Context(), "game_events", data).Err()
    if err != nil {
        fmt.Println("Error publishing join event to Redis:", err)
        return
    }
    go player.listenForMessages()
    go player.listenForWrites()
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
    defer func(){
        c.game.handleDisconnect(c)
        c.conn.Close()
    }()
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