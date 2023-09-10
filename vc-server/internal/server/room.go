package server

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
	"varchess/internal/game"

	"github.com/olahol/melody"
)

type room struct {
	id string
	game *game.Game
	members map[*melody.Session]client
}

type client struct {
	username string
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
			members: make(map[*melody.Session]client),
		}
		WriteJSON(w,http.StatusOK,response)
	}
}

func (s *server) handleJoinRoom() http.HandlerFunc{
	type request struct {
		Username string `json:"username"`
		AccessToken string `json:"accessToken"`
		RoomId string `json:"roomId"`
	}

	type response struct{
		GameState string `json:"gameState"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var request request
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			WriteJSON(w,http.StatusBadRequest, err)
			return
		}
		var response response = response{GameState:"Test"}
		WriteJSON(w,http.StatusOK,response)
	}
}
