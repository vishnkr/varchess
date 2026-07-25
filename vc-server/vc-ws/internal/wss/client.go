package wss

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// Client represents a connected WebSocket player.
type Client struct {
	conn          *websocket.Conn
	userId        string
	username      string
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
	case "draw-offer":
		drawEvent := Event{
			GameID: c.gameID,
			UserID: c.userId,
			Type:   DrawOffer,
			Data:   nil,
		}
		drawJSON, _ := json.Marshal(drawEvent)
		redisClient.Publish(context.Background(), UserAction, drawJSON)
	case "draw-accept":
		acceptEvent := Event{
			GameID: c.gameID,
			UserID: c.userId,
			Type:   DrawAccept,
			Data:   nil,
		}
		acceptJSON, _ := json.Marshal(acceptEvent)
		redisClient.Publish(context.Background(), UserAction, acceptJSON)
	case "draw-reject":
		rejectEvent := Event{
			GameID: c.gameID,
			UserID: c.userId,
			Type:   DrawReject,
			Data:   nil,
		}
		rejectJSON, _ := json.Marshal(rejectEvent)
		redisClient.Publish(context.Background(), UserAction, rejectJSON)
	case "hints":
		payloadBytes, err := json.Marshal(wsMsg.Payload)
		if err != nil {
			log.Println("Failed to marshal hints payload:", err)
			return
		}
		hintsEvent := Event{
			GameID: c.gameID,
			UserID: c.userId,
			Type:   Hints,
			Data:   payloadBytes,
		}
		hintsJSON, err := json.Marshal(hintsEvent)
		if err != nil {
			log.Println("Failed to marshal hints event:", err)
			return
		}
		if err := redisClient.Publish(context.Background(), UserAction, hintsJSON).Err(); err != nil {
			fmt.Println("Error publishing hints event to Redis:", err)
		}
	case "chat":
		payloadBytes, err := json.Marshal(wsMsg.Payload)
		if err != nil {
			log.Println("Failed to marshal chat payload:", err)
			return
		}
		var incoming struct {
			Message string `json:"msg"`
		}
		if err := json.Unmarshal(payloadBytes, &incoming); err != nil || strings.TrimSpace(incoming.Message) == "" {
			return
		}
		sender := c.username
		if sender == "" {
			sender = c.userId
		}
		out := WSMessage{
			Type: "chat",
			Payload: ChatMessage{
				GameId:  c.gameID,
				Sender:  sender,
				Message: strings.TrimSpace(incoming.Message),
			},
		}
		msgJSON, err := json.Marshal(out)
		if err != nil {
			return
		}
		hub.BroadcastToGame(c.gameID, msgJSON)
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

	// If this socket was superseded by a reconnect, do not announce offline
	// or schedule removal — the new Client owns the slot.
	c.game.mu.Lock()
	current, ok := c.game.Clients[c.userId]
	isCurrent := ok && current == c
	c.game.mu.Unlock()
	if !isCurrent {
		log.Printf("Superseded connection closed for player %s (game %s)", c.userId, c.gameID)
		return
	}

	fmt.Println("Player disconnected:", c.userId)
	hub.BroadcastPresence(c.gameID, c.userId, false)

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
