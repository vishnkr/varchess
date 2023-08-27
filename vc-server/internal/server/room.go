package server

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
)

func generateRandomString(length int) (string, error) {
    bytes := make([]byte, length)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

func (s *server) handleCreateRoom() http.HandlerFunc {
	
	type response struct {
		RoomId string `json:"roomId"`
		AccessToken string `json:"accessToken"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		roomId,err := generateRandomString(8)
		if err!=nil{
			WriteJSON(w,http.StatusBadRequest,err)
			return
		}
		accessToken,err := generateRandomString(32)
		if err!=nil{
			WriteJSON(w,http.StatusBadRequest,err)
			return
		}
		var response response = response{ RoomId: roomId, AccessToken: accessToken}

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
