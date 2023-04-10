package server

import (
	"encoding/json"
	"net/http"
	"varchess/internal/store"

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

func (s *Server) Start(allowedOrigins string) error {
	router := mux.NewRouter()
	router.Use(setHeadersMiddleware(allowedOrigins))
	router.HandleFunc("/room-state", makeHTTPHandleFunc(s.roomStateHandler)).Methods("GET")
	router.HandleFunc("/create-room", makeHTTPHandleFunc(s.createRoomHandler)).Methods("POST")
	router.HandleFunc("/login", makeHTTPHandleFunc(s.authenticateUserHandler)).Methods("GET")
	router.HandleFunc("/signup", makeHTTPHandleFunc(s.createAccountHandler)).Methods("POST")
	router.HandleFunc("/possible-squares", makeHTTPHandleFunc(s.getPossibleSquares)).Methods("GET")
	router.HandleFunc("/server-status", makeHTTPHandleFunc(func(w http.ResponseWriter, r *http.Request) error{
		w.WriteHeader(http.StatusOK)
		return nil
	})).Methods("GET")

	wsServer := NewWebsocketServer()
	go wsServer.Run()
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWsHandler(wsServer, w, r)
	})
	return http.ListenAndServe(s.listenAddr, router)
}

func setHeadersMiddleware(allowedOrigins string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Content-Type", "application/json")
            w.Header().Set("Access-Control-Allow-Origin", allowedOrigins)
            next.ServeHTTP(w, r)
        })
    }
}