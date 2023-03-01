package server

import (
	"encoding/json"
	"net/http"
	"varchess/pkg/store"

	"github.com/gorilla/mux"
)

type ApiError struct {
	Error string
}

type apiFunction func(http.ResponseWriter, *http.Request) error

type Server struct {
	listenAddr string
	store      store.Storage
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandleFunc(f apiFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
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
	router.HandleFunc("/getPossibleToSquares", makeHTTPHandleFunc(s.GetPossibleSquares)).Methods("GET")
	wsServer := NewWebsocketServer()
	go wsServer.Run()
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWsHandler(wsServer, w, r)
	})
	return http.ListenAndServe(s.listenAddr, router)
}