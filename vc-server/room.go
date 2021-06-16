package main


import (
	"fmt"
	"math/rand"
	"time"
	"net/http"
	"encoding/json"
)

type Room struct{
	Game *Game
	Clients map[*Client]bool
	Id string
}

var RoomsMap = make(map[string]*Room)
const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func genRandSeq(length int) string {
	b := make([]byte, length)
	for i := range b {
	  b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
  }

func roomHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json") 
	w.Header().Set("Access-Control-Allow-Origin", "*")
	response:= MessageStruct{
		Type: "getRoomId",
		Data: genRandSeq(6),
	}
	json.NewEncoder(w).Encode(response)
}

func (c *Client) CreateRoom(roomId string){
	c.roomId = roomId
	RoomsMap[roomId] = &Room{Game:&Game{}, Clients: make(map[*Client]bool), Id: c.roomId}
	RoomsMap[roomId].Clients[c] = true
	fmt.Println("rooms",RoomsMap,*RoomsMap[roomId])
}

func (c *Client) AddtoRoom(roomId string){
	RoomsMap[roomId].Clients[c] = true
	c.roomId = roomId
}

