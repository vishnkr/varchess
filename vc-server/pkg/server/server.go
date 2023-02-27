package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"varchess/pkg/game"
	"varchess/pkg/store"

	"github.com/gorilla/mux"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

type ApiError struct {
	Error string
}

type apiFunction func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFunc(f apiFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

type Server struct {
	listenAddr string
	store      store.Storage
}

func NewServer(listenAddr string, store store.Storage) *Server {
	return &Server{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *Server) Start() error {
	router := mux.NewRouter()
	router.HandleFunc("/getRoomId", makeHTTPHandleFunc(s.RoomHandler)).Methods("POST")
	router.HandleFunc("/getBoardFen/{roomId}", makeHTTPHandleFunc(s.BoardStateHandler)).Methods("GET", "OPTIONS")
	router.HandleFunc("/", rootHandler)
	router.HandleFunc("/login", makeHTTPHandleFunc(s.AuthenticateUserHandler)).Methods("GET")
	router.HandleFunc("/signup", makeHTTPHandleFunc(s.CreateAccountHandler)).Methods("POST")
	router.HandleFunc("/getPossibleToSquares", makeHTTPHandleFunc(s.GetPossibleSquares)).Methods("POST")
	wsServer := NewWebsocketServer()
	go wsServer.Run()
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWsHandler(wsServer, w, r)
	})
	return http.ListenAndServe(s.listenAddr, router)
}

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
	json.NewEncoder(w).Encode(response)
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

func (s *Server) GetPossibleSquares(w http.ResponseWriter, r *http.Request) error {
	//optimize this request to be done once before every move instead of once after every click, store result in client side
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST")

		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token, Authorization")
		return nil
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var objmap map[string]interface{}
	json.NewDecoder(r.Body).Decode(&objmap)
	var pColor game.Color
	srcRow, srcCol, color := int(objmap["srcRow"].(float64)), int(objmap["srcCol"].(float64)), objmap["color"].(string)
	piece := game.StrToTypeMap[objmap["piece"].(string)]
	room, _ := RoomsMap[objmap["roomId"].(string)]
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
	response := &PossibleMoves{Moves: moves, Piece: objmap["piece"].(string)}
	return WriteJSON(w, http.StatusOK, response)
}
