package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"varchess/internal/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/olahol/melody"
)

type ApiError struct {
	Error string
}


type apiFunction func(http.ResponseWriter, *http.Request) error

type server struct {
	listenAddr   string
	router       *chi.Mux
	wsManager 	*melody.Melody
	rooms 		map[string]*room
}

var l = logger.Get()

func makeHTTPHandleFunc(f apiFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func NewServer(listenAddr string, allowedOrigins string) *server {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(AllowOptions)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{allowedOrigins},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{},
		AllowCredentials: false,
		MaxAge:           3600,
	}))
	
	s := &server{
		listenAddr : listenAddr,
		router: router,
		rooms:  make(map[string]*room),
	}
	melody := createMelodyForRooms(s)
	s.wsManager = melody
	s.routes()
	return s
}

func (s *server) routes() {	
	s.router.Get("/health", makeHTTPHandleFunc(s.handleHealthCheck))
	// Room Handlers
	s.router.Post("/rooms", s.handleCreateRoom())
	s.router.Post("/room-state", s.handleGetRoomState())
	s.router.HandleFunc("/ws/{roomId}/{username}", func(w http.ResponseWriter, r *http.Request) {
		err:= s.wsManager.HandleRequest(w, r)
		if err!=nil{
			fmt.Println(err.Error())
		}
	})
}

func (s *server) Start() {
	l.Fatal().Err(http.ListenAndServe(s.listenAddr, requestLogger(s.router))).Msg("Varchess Server closed")
}

/*func setHeadersMiddleware(allowedOrigins string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigins)
			next.ServeHTTP(w, r)
		})
	}
}*/

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func (s *server) handleHealthCheck(w http.ResponseWriter, r *http.Request) error {
	return WriteJSON(w, http.StatusOK, "health check OK")
}
