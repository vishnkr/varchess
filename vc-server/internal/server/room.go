package server

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
	"varchess/internal/game"

	"github.com/olahol/melody"
)

type member struct{
	Session *melody.Session `json:"-"`
	IsHost bool `json:"isHost"`
	PlayerColor game.Color `json:"playerColor"`
}

type memberResponse struct{
	Username string `json:"username"`
	PlayerColor game.Color `json:"color"`
	IsHost bool `json:"isHost"`
}

type room struct {
	id string
	game *game.Game
	members []member
	inGame bool
}

type client struct {
	session *melody.Session
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateRandomString(length int) string {
   b := make([]byte, length)
   for i := range b {
      b[i] = charset[seededRand.Intn(len(charset))]
   }
   return string(b)
}

func (s *server) handleCreateRoom() http.HandlerFunc {
	
	type response struct {
		RoomId string `json:"roomId"`
		AccessToken string `json:"accessToken"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		roomId := generateRandomString(8)
		accessToken := generateRandomString(32)
		var response response = response{ RoomId: roomId, AccessToken: accessToken}
		s.rooms[roomId] = &room{
			id: roomId,
			game: &game.Game{},
			members: []member{},
			inGame: false,
		}
		WriteJSON(w,http.StatusOK,response)
	}
}

//handleGetRoomState : Get current room state, to be used when joining a room or reconnecting to a room after page refresh/websocket disconnect
func (s *server) handleGetRoomState() http.HandlerFunc{
	type request struct {
		RoomId string `json:"roomId"`
	}
	
	type response struct {
		Members []memberResponse `json:"members"`
		InGame bool `json:"inGame"`
		FEN string `json:"fen"`

	}
	return func(w http.ResponseWriter, r *http.Request){
		var request request
		err:= json.NewDecoder(r.Body).Decode(&request)
		if err!=nil{
			WriteJSON(w,http.StatusBadRequest, err)
			return
		}
		room:= s.rooms[request.RoomId]
		
		
		response  := response{
			Members: room.getRoomMembers(),
			InGame: room.inGame,
		}

		WriteJSON(w,http.StatusOK,response)
	}
}

func (r *room) getRoomMembers()[]memberResponse{
	members := []memberResponse{}
	for _,member:= range r.members{
		username,exists := member.Session.Get("username") 
		if exists{
			members = append(members,memberResponse{Username: username.(string),IsHost: member.IsHost})
		}
	}
	return members
}