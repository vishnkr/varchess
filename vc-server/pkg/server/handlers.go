package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"varchess/pkg/game"

	"github.com/gorilla/mux"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "home")
}

func (s *Server) RoomHandler(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	uniqueRoomId := genRandSeq(6)
	for ok := true; ok; _, ok = RoomsMap[uniqueRoomId] {
		uniqueRoomId = genRandSeq(6)
	}
	response := MessageStruct{
		Type: "getRoomId",
		Data: uniqueRoomId,
	}
	WriteJSON(w,http.StatusOK,response)
	return nil
}

func (s *Server) BoardStateHandler(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "OPTIONS" {
		return nil
	}
	params := mux.Vars(r)
	id := params["roomId"]
	room, ok := RoomsMap[id]
	if ok {
		response := BoardState{
			Fen:    game.ConvertBoardtoFEN(room.Game.Board),
			RoomId: id,
		}
		if room.Game.Board.CustomMovePatterns != nil {
			response.MovePatterns = room.Game.Board.CustomMovePatterns
		}
		return WriteJSON(w, http.StatusOK, response)
	} else {
		errResponse := MessageStruct{Type: "error", Data: "Room does not exist/has been closed"}
		return WriteJSON(w, http.StatusOK, errResponse)
	}
}

type PossibleSquaresRequestBody struct{
	SrcRow int `json:"srcRow"`
	SrcCol int `json:"srcCol"`
	Color string `json:"color"`
	RoomId string `json:"roomId"`
	Piece string `json:"piece"`
}

func (s *Server) GetPossibleSquares(w http.ResponseWriter, r *http.Request) error {
	//optimize this request to be done once before every move instead of once after every click, store result in client side
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET")

		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token, Authorization")
		return nil
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	req:= &PossibleSquaresRequestBody{}
	err:= json.NewDecoder(r.Body).Decode(&req)
	if err!=nil{
		return errors.New("Invalid request body format")
	}

	var pColor game.Color
	piece := game.StrToTypeMap[req.Piece]
	room := RoomsMap[req.RoomId]
	if room == nil {
		return errors.New("Room does not exist")
	}
	if req.Color == "white" {
		pColor = game.White
	} else {
		pColor = game.Black
	}
	board := room.Game.Board
	moves := make([][]int, 0)
	valid := board.GetAllValidMoves(pColor)
	for move, p := range valid {
		if p.Type == piece && move.SrcRow == req.SrcRow && move.SrcCol == req.SrcCol {
			moves = append(moves, []int{move.DestRow, move.DestCol})
		}
	}
	response := &PossibleMoves{Moves: moves, Piece: req.Piece}
	return WriteJSON(w, http.StatusOK, response)
}
