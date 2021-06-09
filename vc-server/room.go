package main


import (
	"fmt"
	"math/rand"
	"time"
)

type Room struct{
	Game *Game
	Clients map[*Client]bool
	Id string
}

var RoomList = make(map[string]*Room)

const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func genRandSeq(length int) string {
	b := make([]byte, length)
	for i := range b {
	  b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
  }

func (c *Client) CreateRoom(roomId string ){
	c.roomId = roomId
	RoomList[roomId] = &Room{Game:&Game{}, Clients: make(map[*Client]bool), Id: c.roomId}
	RoomList[roomId].Clients[c] = true
	fmt.Println("rooms",RoomList,*RoomList[roomId])
}

func (c *Client) AddtoRoom(roomId string){
	RoomList[roomId].Clients[c] = true
	c.roomId = roomId
}
