package wss

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)


type Client struct {
    conn   *websocket.Conn
	userId string
    game *Game
	send chan []byte
    gameID       string
    colorPref    string
    authenticated bool
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
    var wsMsg WSMessage
	json.Unmarshal(msg, &wsMsg)

	switch wsMsg.Type {
	case "move":
		moveEvent := Event{GameID: c.gameID, UserID: c.userId, Type: Move, Data: wsMsg.Payload}
		//var move MovePayload
		//json.Unmarshal(wsMsg.Payload, &move)
		//r. .Moves = append(r.gameState.Moves, move)

		//hub.SaveGameStateToRedis(r.gameID, r.gameState)
        moveJson,err := json.Marshal(moveEvent)
		//r.Broadcast(msg)
		err = redisClient.Publish(context.Background(), UserAction, moveJson).Err()
		if err != nil {
			fmt.Println("Error publishing join event to Redis:", err)
			return
		}
	case "chat":
		//r.Broadcast(msg)
	case "resign":
		//r.Broadcast(msg)
	case "draw_offer":
		//r.Broadcast(msg)
	}
}

func (c *Client) listenForWrites() {
    defer c.conn.Close()
    for msg := range c.send {
        err := c.conn.WriteMessage(websocket.TextMessage, msg)
        if err != nil {
            fmt.Println("WebSocket write error:", err)
            return
        }
    }
}


func (g *Game) handleDisconnect(c *Client) {
    if g == nil {
        log.Println("Game is nil in handleDisconnect")
        return
    }
    g.mu.Lock()
    defer g.mu.Unlock()

    fmt.Println("Player disconnected:", c.userId)
    go func() {
        time.Sleep(30 * time.Second)
        g.mu.Lock()
        defer g.mu.Unlock()

        if g.players[c.userId] == c {
            delete(g.players, c.userId)
            log.Println("Player removed after timeout:", c.userId)
        }
    }()
}