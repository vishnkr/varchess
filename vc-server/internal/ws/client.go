package ws

import (
	"encoding/json"
	"log"
	"time"
	"varchess/internal/user"

	"github.com/gorilla/websocket"
)

type Client struct{
	user *user.User
	hubs map[string]*gameHub
	conn *websocket.Conn
	ws *WebSocket
	send chan []byte
}

const (
	writeWait = 10 * time.Second
	pongWait = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
)


func (c *Client) readPump(){
	defer func(){
		c.unregisterAll()
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string)error{c.conn.SetReadDeadline(time.Now().Add(pongWait));return nil})
	for {
		_,message,err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			 log.Printf("error: %v", err)
			}
			break
		   }
		handleMessage(c,message)
	}
}

func (c *Client) writePump(){
	ticker := time.NewTicker(pingPeriod)
	defer func(){
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select{
		case message, ok:= <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok{
				c.conn.WriteMessage(websocket.CloseMessage,[]byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {
			 return
			}
		}
	}
}

func (c *Client) unregisterAll() {
	for _, hub := range c.hubs {
	 hub.unregister <- c
	}
}

func handleMessage(c *Client, msg []byte) {
	var request Request
	err := json.Unmarshal(msg, &request)
	if err != nil {
		closeConnection(c, websocket.CloseInvalidFramePayloadData, "Bad Event Type")
		return
	}

	switch eventType(request.Event){
	case EventUserConnect:
		handleUserConnect(c,request)
	case EventChatMessage:
		handleChatMessage(c,request)
	case EventGameMakeMove:
		handleMakeMove(c,request)
	case EventGameDrawOffer:
		handleOfferDraw(c,request)
	case EventGameDrawResult:
		handleDrawResult(c,request)
	case EventGameResign:
		handleResign(c,request)
	
	default:
		closeConnection(c, websocket.CloseInvalidFramePayloadData, "Unsupported Event Type")
		return
	}
}