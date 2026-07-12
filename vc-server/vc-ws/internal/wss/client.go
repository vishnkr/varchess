package wss

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Client represents a connected WebSocket player.
type Client struct {
	conn          *websocket.Conn
	userId        string
	game          *Game
	send          chan []byte
	gameID        string
	colorPref     string // "w" or "b"
	authenticated bool
}

func (c *Client) listenForMessages() {
	defer func() {
		c.handleDisconnect()
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
	var wsMsg WSMessage
	if err := json.Unmarshal(msg, &wsMsg); err != nil {
		log.Println("Invalid WS message:", err)
		return
	}

	switch wsMsg.Type {
	case "move":
		payloadBytes, err := json.Marshal(wsMsg.Payload)
		if err != nil {
			log.Println("Failed to marshal move payload:", err)
			return
		}
		moveEvent := Event{
			GameID: c.gameID,
			UserID: c.userId,
			Type:   Move,
			Data:   payloadBytes,
		}
		moveJSON, err := json.Marshal(moveEvent)
		if err != nil {
			log.Println("Failed to marshal move event:", err)
			return
		}
		if err := redisClient.Publish(context.Background(), UserAction, moveJSON).Err(); err != nil {
			fmt.Println("Error publishing move event to Redis:", err)
		}
	case "resign":
		resignEvent := Event{
			GameID: c.gameID,
			UserID: c.userId,
			Type:   Resign,
			Data:   nil,
		}
		resignJSON, _ := json.Marshal(resignEvent)
		redisClient.Publish(context.Background(), UserAction, resignJSON)
	case "draw_offer":
		drawEvent := Event{
			GameID: c.gameID,
			UserID: c.userId,
			Type:   DrawOffer,
			Data:   nil,
		}
		drawJSON, _ := json.Marshal(drawEvent)
		redisClient.Publish(context.Background(), UserAction, drawJSON)
	case "chat":
		// TODO: broadcast chat to other clients in game.
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

func (c *Client) handleDisconnect() {
	if c.game == nil {
		log.Println("Game is nil in handleDisconnect")
		return
	}

	fmt.Println("Player disconnected:", c.userId)

	// Give a 30-second reconnect window before removing the player.
	go func() {
		time.Sleep(30 * time.Second)
		c.game.mu.Lock()
		defer c.game.mu.Unlock()
		if existing, ok := c.game.Clients[c.userId]; ok && existing == c {
			delete(c.game.Clients, c.userId)
			log.Println("Player removed after reconnect timeout:", c.userId)
		}
	}()
}
