package server

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
	"varchess/pkg/game"

	"github.com/gorilla/mux"
)

type Room struct{
	Game *game.Game
	Clients map[*Client]bool
	Id string
	P1 *Client
	P2 *Client
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

func RoomHandler(w http.ResponseWriter, r *http.Request){
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
	Members []string `json:"members"`
	Result string `json:"result"`
}

func (c *Client) CreateRoom(roomId string,startFen string) *Room{
	c.mu.Lock()
	defer c.mu.Unlock()
	c.roomId = roomId
	RoomsMap[roomId] = &Room{
		Game: &game.Game{
			Board: game.ConvertFENtoBoard(startFen),
			Turn: "w",
			}, 
		Clients: make(map[*Client]bool), 
		Id: c.roomId,
		P1: c,
	}
	game.DisplayBoardState(RoomsMap[roomId].Game.Board)
	RoomsMap[roomId].Clients[c] = true
	gameInfo:= GameInfo{Type:"gameInfo",P1:c.username,Turn:"w",RoomId: roomId, Members:[]string{}}
	gameInfo.Members = append(gameInfo.Members,c.username)
	marshalledInfo,_ := json.Marshal(gameInfo)
	RoomsMap[roomId].BroadcastToMembers(marshalledInfo)
	return RoomsMap[roomId]
}

func (room *Room) BroadcastToMembers(message []byte){
	for client := range room.Clients{
		client.send <- message
	}
}

func (room *Room) BroadcastToMembersExceptSender(message []byte,c *Client){
	for member := range room.Clients {
		if (member.conn!=c.conn){
			member.send <- message
		}
	}
}

func (c *Client) AddtoRoom(roomId string){
	c.mu.Lock()
	defer c.mu.Unlock()
	curRoom, ok:= RoomsMap[roomId]
	if ok{
		var gameInfo GameInfo
		if (len(curRoom.Clients) == 1){
			RoomsMap[roomId].P2 = c
			gameInfo = GameInfo{Type:"gameInfo",P1:curRoom.P1.username,P2: c.username,Turn:curRoom.Game.Turn,RoomId: roomId,Members:RoomsMap[roomId].getClientUsernames()}
		} else { 
			gameInfo = GameInfo{Type:"gameInfo",P1:curRoom.P1.username,P2: curRoom.P2.username,Turn:curRoom.Game.Turn,RoomId: roomId,Members:RoomsMap[roomId].getClientUsernames()}
		}
		gameInfo.Members = append(gameInfo.Members,c.username)
		RoomsMap[roomId].Clients[c] = true
		c.roomId = roomId
		marshalledInfo,_ := json.Marshal(gameInfo)
		RoomsMap[roomId].BroadcastToMembers(marshalledInfo)
	} else {
		log.Println("Room close")
		message := MessageStruct{Type:"error",Data:"Room does not exist, connection expired"}
		if errMessage,err:= json.Marshal(message);err==nil{
			c.send <- errMessage
		}
	}	
}

func (room *Room) getClientUsernames() []string{
	var clientList = []string{}
	for client := range room.Clients{
		clientList = append(clientList,client.username)
	}
	return clientList
}

type BoardState struct{
	Fen string
	RoomId string
	MovePatterns []game.MovePatterns `json:"movePatterns"`
}

func BoardStateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") 
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if (r.Method == "OPTIONS") {
        return
    } else {
        params:= mux.Vars(r)
	id:= params["roomId"]
	room, ok := RoomsMap[id] 
	if ok {
		response:= BoardState{
			Fen: game.ConvertBoardtoFEN(room.Game.Board),
			RoomId: id,
		}
		if (room.Game.Board.CustomMovePatterns!=nil){
			response.MovePatterns = room.Game.Board.CustomMovePatterns
		}
		json.NewEncoder(w).Encode(response)
	} else { 
		errResponse:= MessageStruct{Type:"error",Data:"Room does not exist/has been closed"}
		json.NewEncoder(w).Encode(errResponse)
	}
    }
}