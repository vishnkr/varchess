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

func (c *Client) CreateRoom(roomId string,startFen string){
	c.roomId = roomId
	RoomsMap[roomId] = &Room{
		Game: &Game{
			Board: ConvertFENtoBoard(startFen),
			P1: c,
			}, 
		Clients: make(map[*Client]bool), 
		Id: c.roomId,
	}
	DisplayBoardState(RoomsMap[roomId].Game.Board)
	RoomsMap[roomId].Clients[c] = true
	fmt.Println("you are p1")
	fmt.Println("rooms",RoomsMap,*RoomsMap[roomId])
}

func (c *Client) AddtoRoom(roomId string){
	c.mu.Lock()
	defer c.mu.Unlock()
	if (len(RoomsMap[roomId].Clients) == 1){
		RoomsMap[roomId].Game.P2 = c
		fmt.Println("you are p2")
	} else { fmt.Println("you are a viewer")}
	RoomsMap[roomId].Clients[c] = true
	c.roomId = roomId
}

