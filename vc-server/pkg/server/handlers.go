package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"varchess/pkg/game"

	"github.com/gorilla/mux"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "home")
}

func (s *Server) ServerStatusHandler(w http.ResponseWriter, r *http.Request) error{
    w.WriteHeader(http.StatusOK)
    return nil
}

func (s *Server) RoomHandler(w http.ResponseWriter, r *http.Request) error {

	uniqueRoomId := genRandSeq(6)
	for ok := true; ok; _, ok = RoomsMap[uniqueRoomId] {
		uniqueRoomId = genRandSeq(6)
	}
	dataBytes,err:=json.Marshal(uniqueRoomId)
	if err!=nil{
		fmt.Println(err)
	}
	response := MessageStruct{
		Type: "getRoomId",
		Data: json.RawMessage(dataBytes),
	}
	WriteJSON(w,http.StatusOK,response)
	return nil
}

func (s *Server) BoardStateHandler(w http.ResponseWriter, r *http.Request) error {
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
		errorMsg := "Room does not exist, connection expired"
		errorMsgBytes, _ := json.Marshal(errorMsg)
		errResponse := MessageStruct{Type: "error", Data: json.RawMessage(errorMsgBytes)}
		return WriteJSON(w, http.StatusOK, errResponse)
	}
}

func (s *Server) GetPossibleSquares(w http.ResponseWriter, r *http.Request) error {
	//optimize this request to be done once before every move instead of once after every click, store result in client side
	query := r.URL.Query()
    roomID := query.Get("roomid")
    color := query.Get("color")
    pieceStr := query.Get("piece")
    startRow := query.Get("src_row")
    startCol := query.Get("src_col")

	var pColor game.Color
	piece := game.StrToTypeMap[pieceStr]
	room := RoomsMap[roomID]
	if room == nil {
		return errors.New("Room does not exist")
	}

	srcRow, err := strconv.Atoi(startRow)
    if err != nil {
        return errors.New("Invalid src row: " + startRow)
    }
	
    srcCol, err := strconv.Atoi(startCol)
    if err != nil {
        return errors.New("Invalid src column: " + startCol)
    }

	if color == "white" {
		pColor = game.White
	} else {
		pColor = game.Black
	}
	board := room.Game.Board
	moves := make([][]int, 0)
	valid := board.GetAllValidMoves(pColor)
	for move, p := range valid {
		if p.Type == piece && move.SrcRow == srcRow && move.SrcCol == srcCol {
			moves = append(moves, []int{move.DestRow, move.DestCol})
		}
	}
	response := &PossibleMoves{Moves: moves, Piece: pieceStr}
	return WriteJSON(w, http.StatusOK, response)
}
