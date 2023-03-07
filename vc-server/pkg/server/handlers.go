package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"varchess/pkg/game"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "home")
}

func (s *Server) ServerStatusHandler(w http.ResponseWriter, r *http.Request) error{
    w.WriteHeader(http.StatusOK)
    return nil
}

type CreateRoomResponse struct {
	RoomId string `json:"roomId"`
}

func (s *Server) CreateRoomHandler(w http.ResponseWriter,r *http.Request) error{
	uniqueRoomId := genRandSeq(6)
	for ok := true; ok; _, ok = RoomsMap[uniqueRoomId] {
		uniqueRoomId = genRandSeq(6)
	}
	roomInfo := &CreateRoomInfo{}
	err := json.NewDecoder(r.Body).Decode(roomInfo)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
	}
	curRoom := &Room{
		Game: &game.Game{
			Board: game.ConvertFENtoBoard(roomInfo.StartFEN),
			Turn:  game.White,
		},
		Clients: make(map[*Client]bool),
		Id:      uniqueRoomId,
		DrawOffer: DrawOffer{IsOffered: false},
	}
	RoomsMap[uniqueRoomId] = curRoom
	game.DisplayBoardState(curRoom.Game.Board)

	if len(roomInfo.CustomMovePatterns) != 0 {
		RoomsMap[uniqueRoomId].Game.Board.CustomMovePatterns = roomInfo.CustomMovePatterns
	}
	response:= &CreateRoomResponse{RoomId: uniqueRoomId}
	// delete room if no one joins for 20s
	go func(roomId string) {
        time.Sleep(20 * time.Second)
        room, ok := RoomsMap[roomId]
        if ok && len(room.Clients) == 0 {
			fmt.Println("closing room due to inactivity")
            delete(RoomsMap, roomId)
        }
    }(uniqueRoomId)
	
	return WriteJSON(w,http.StatusOK,response)
}

func (s *Server) RoomStateHandler(w http.ResponseWriter, r *http.Request) error{
	query := r.URL.Query()
    roomId := query.Get("roomid")
	curRoom,ok:= RoomsMap[roomId]
	var response RoomState
	if ok{
		response = RoomState{
			Fen:    game.ConvertBoardtoFEN(curRoom.Game.Board),
			RoomId: roomId,
			P1: curRoom.P1.username,
			P2: curRoom.P2.username,
			Members: curRoom.getClientUsernames(),
		}
		if curRoom.Game.Board.CustomMovePatterns != nil {
			response.MovePatterns = curRoom.Game.Board.CustomMovePatterns
		}
	} else {
		return WriteJSON(w, http.StatusBadRequest,ApiError{Error: "Invalid Room ID"})
	}
	
	return WriteJSON(w,http.StatusOK,response)
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
