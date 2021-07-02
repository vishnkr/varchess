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


const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
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
	uniqueRoomId:= genRandSeq(6)
	for ok := true; ok; _, ok = RoomsMap[uniqueRoomId] {
		uniqueRoomId = genRandSeq(6)
	}
	response:= MessageStruct{
		Type: "getRoomId",
		Data: uniqueRoomId,
	}
	json.NewEncoder(w).Encode(response)
}

type GameInfo struct{
	Type string `json:"type"`
	P1 string `json:"p1"`
	P2 string `json:"p2"`
    Turn string `json:"turn"`
	RoomId string `json:"roomId"`
}

func (c *Client) CreateRoom(roomId string,startFen string){
	c.roomId = roomId
	RoomsMap[roomId] = &Room{
		Game: &Game{
			Board: ConvertFENtoBoard(startFen),
			P1: c,
			Turn: "w",
			}, 
		Clients: make(map[*Client]bool), 
		Id: c.roomId,
	}
	DisplayBoardState(RoomsMap[roomId].Game.Board)
	RoomsMap[roomId].Clients[c] = true
	fmt.Println("you are p1")
	gameInfo:= GameInfo{Type:"gameInfo",P1:c.username,Turn:"w",RoomId: roomId}
	marshalledInfo,_ := json.Marshal(gameInfo)
	RoomsMap[roomId].BroadcasToMembers(marshalledInfo)
	fmt.Println("rooms",RoomsMap,*RoomsMap[roomId])
}

func (room *Room) BroadcasToMembers(message []byte){
	for client,_ := range room.Clients{
		client.send <- message
	}
}

func (c *Client) AddtoRoom(roomId string){
	c.mu.Lock()
	defer c.mu.Unlock()
	curRoom:= RoomsMap[roomId]
	var gameInfo GameInfo
	if (len(curRoom.Clients) == 1){
		RoomsMap[roomId].Game.P2 = c
		fmt.Println("you are p2")
		gameInfo = GameInfo{Type:"gameInfo",P1:curRoom.Game.P1.username,P2: c.username,Turn:curRoom.Game.Turn,RoomId: roomId}
	} else { 
		fmt.Println("you are a viewer")
		gameInfo = GameInfo{Type:"gameInfo",P1:curRoom.Game.P1.username,P2: curRoom.Game.P2.username,Turn:curRoom.Game.Turn,RoomId: roomId}
	}
	RoomsMap[roomId].Clients[c] = true
	c.roomId = roomId
	marshalledInfo,_ := json.Marshal(gameInfo)
	RoomsMap[roomId].BroadcasToMembers(marshalledInfo)
}

