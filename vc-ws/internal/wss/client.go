package wss

import "github.com/gorilla/websocket"


type Client struct {
    conn   *websocket.Conn
	userId string
    game *Game
	send chan []byte
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

func (c *Client) WritePump() {
    for msg := range c.send {
        c.conn.WriteMessage(websocket.TextMessage, msg)
    }
}
