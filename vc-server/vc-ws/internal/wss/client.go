package wss

import (
	"fmt"
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

func (c *Client) ReadPump() {
    defer func() {
        c.game.UnregisterClient(c)
        c.conn.Close()
    }()
    for {
        _, msg, err := c.conn.ReadMessage()
        if err != nil {
            c.game.UnregisterClient(c)
            break
        }
        c.game.HandleMessage(c, msg)
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
    g.mu.Lock()
    defer g.mu.Unlock()

    fmt.Println("Player disconnected:", c.userId)
    go func() {
        time.Sleep(30 * time.Second)
        g.mu.Lock()
        defer g.mu.Unlock()

        if g.players[c.userId] == c {
            delete(g.players, c.userId)
            fmt.Println("Player removed after timeout:", c.userId)
        }
    }()
}