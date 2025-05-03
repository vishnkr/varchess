package wss

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"vc-server/vc-ws/internal/config"

	"github.com/gorilla/websocket"
)

var (
    upgrader  = websocket.Upgrader{
        CheckOrigin: func(r *http.Request) bool { return true },
    }
)

const SystemAction = "system_action"
const UserAction = "user_action"

type AuthResponse struct {
    UserID  string `json:"userId"`
    Username string `json:"username"`
}

func authenticateToken(config *config.Config, token string) (*AuthResponse, error) {
    url := fmt.Sprintf("http://%s:%s/api/auth/validate", config.ServerHost, config.ServerPort)

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



func HandleWSConnection(w http.ResponseWriter, r *http.Request, config *config.Config) {
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
        conn: conn,
        gameID: gameID,
        authenticated: false,
        send: make(chan []byte,256),
    }

    go player.handleAuthentication(config)
}

func (c *Client) handleAuthentication(config *config.Config) {
    _, msg, err := c.conn.ReadMessage()
    if err != nil {
        fmt.Println("Error reading initial message:", err)
        c.conn.Close()
        return
    }

    var connectPayload struct {
        Token    string `json:"token"`
        ColorPref string `json:"colorPref,omitempty"`
    }

 
    err = json.Unmarshal(msg, &connectPayload)
    if err != nil {
        fmt.Println("Invalid connection payload:", err)
        c.conn.Close()
        return
    }

    user, err := authenticateToken(config,connectPayload.Token)
    if err != nil {
        fmt.Println("Unauthorized access:", err)
        c.conn.Close()
        return
    }

    c.userId = user.UserID
    c.authenticated = true
    if(connectPayload.ColorPref!=""){
        c.colorPref = connectPayload.ColorPref
    }
    
    c.joinGame()
}

func (c *Client) joinGame() {
    game := hub.GetOrCreateGame(c.gameID)
    game.mu.Lock()
    defer game.mu.Unlock()

    assignedColor := getPlayerColor(game, c.colorPref)
    /*if existingPlayer, exists := game.players[assignedColor]; exists && existingPlayer.userId == c.userId {
        existingPlayer.conn.Close()
        delete(game.players, assignedColor)
    }*/

    if assignedColor == "" {
        c.conn.Close()
        return
    }

    /*player := &Client{
        conn:    c.conn,
        game:    game,
        userId:  c.userId,
        send:    make(chan []byte, 256),
        gameID:  c.gameID,
        colorPref: c.colorPref,
    }*/

    game.players[assignedColor] = c

    var eventData JoinPayload = JoinPayload{Color: assignedColor}
    joinEvent := Event{
        GameID: c.gameID,
        UserID: c.userId,
        Type:   Join,
        Data:   eventData,
    }
    c.game=game
    data, err := json.Marshal(joinEvent)
    if err != nil {
        fmt.Println("Error marshalling join event:", err)
        c.conn.Close()
        return
    }

    err = redisClient.Publish(context.Background(),UserAction, data).Err()
    if err != nil {
        fmt.Println("Error publishing join event to Redis:", err)
        return
    }

    go c.listenForMessages()
    go c.listenForWrites()
}


func getPlayerColor(game *Game, colorPref string) string {
    if len(game.players) == 0 {
        return colorPref
    } else if len(game.players) == 1 {
        if _, exists := game.players["w"]; exists {
            return "b"
        }
        return "w"
    }
    return ""
}

