package server

import (
	"math/rand"
	"time"
	"varchess/pkg/game"
)

type DrawOffer struct {
	IsOffered bool
	Color string
}
type Room struct {
	Game    *game.Game
	Clients map[*Client]bool
	Id      string
	P1      *Client
	P2      *Client
	DrawOffer DrawOffer
}

type PossibleMoves struct {
	Piece string  `json:"piece"`
	Moves [][]int `json:"moves"`
}

type RoomState struct {
	Fen          string `json:"fen"`
	RoomId       string `json:"roomId"`
	Members      []string `json:"members"`
	P1 string `json:"p1,omitempty"`
	P2 string `json:"p2,omitempty"`
	MovePatterns []game.MovePatterns `json:"movePatterns"`
	Turn string `json:"turn"`
}

var RoomsMap = make(map[string]*Room)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func genRandSeq(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

type GameInfo struct {
	Type    string   `json:"type"`
	P1      string   `json:"p1"`
	P2      string   `json:"p2"`
	Turn    string   `json:"turn,omitempty"`
	RoomId  string   `json:"roomId"`
	Members []string `json:"members,omitempty"`
	Result  string   `json:"result,omitempty"`
}


func (room *Room) BroadcastToMembers(message []byte) {
	for client := range room.Clients {
		client.send <- message
	}
}

func (room *Room) BroadcastToMembersExceptSender(message []byte, c *Client) {
	for member := range room.Clients {
		if member.conn != c.conn {
			member.send <- message
		}
	}
}

func (room *Room) getClientUsernames() []string {
	var clientList = []string{}
	for client := range room.Clients {
		clientList = append(clientList, client.username)
	}
	return clientList
}

func (room *Room) getViewerClients() []string{
	var clientList = []string{}
	for client := range room.Clients {
		if ((room.P1!=nil && client.username != room.P1.username) && (room.P2!=nil && client.username!= room.P2.username)){
			clientList = append(clientList, client.username)
		}
		
	}
	return clientList
}