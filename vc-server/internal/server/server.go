package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"varchess/internal/logger"
	"varchess/internal/server/session"
	"varchess/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type ApiError struct {
	Error string
}

type Session struct {
	userId string
}

var (
	AUTH_KEY string = "authenticated"
	USER_ID  string = "user_id"
)

type apiFunction func(http.ResponseWriter, *http.Request) error

type server struct {
	listenAddr   string
	router       *chi.Mux
	store        store.Storage
	sessionStore session.SessionStore[Session]
}

var l = logger.Get()

func makeHTTPHandleFunc(f apiFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func NewServer(listenAddr string, store store.Storage, allowedOrigins string) *server {
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

	var sessionStore session.SessionStore[Session]
	sessionStore.InitStore("auth_store", time.Hour*5)

	s := &server{
		listenAddr:   listenAddr,
		router:       router,
		store:        store,
		sessionStore: sessionStore,
	}
	s.routes()
	return s
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusBadGateway}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (s *server) routes() {
	s.router.Get("/health", makeHTTPHandleFunc(s.handleHealthCheck))

	// Room Handlers
	s.router.Post("/create-room", s.handleCreateRoom())
	//router.HandleFunc("/room-state", makeHTTPHandleFunc(s.roomStateHandler)).Methods("GET")

	//router.HandleFunc("/login", makeHTTPHandleFunc(s.authenticateUserHandler)).Methods("GET")
	//router.HandleFunc("/signup", makeHTTPHandleFunc(s.createAccountHandler)).Methods("POST")
	//router.HandleFunc("/possible-squares", makeHTTPHandleFunc(s.getPossibleSquares)).Methods("GET")

	//wsServer := NewWebsocketServer()
	//go wsServer.Run()
	//router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	//ServeWsHandler(wsServer, w, r)
	//})

}

func (s *server) Start() {
	l.Fatal().Err(http.ListenAndServe(s.listenAddr, requestLogger(s.router))).Msg("Varchess Server closed")
}

func setHeadersMiddleware(allowedOrigins string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigins)
			next.ServeHTTP(w, r)
		})
	}
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, data interface{}, status int) {
	w.WriteHeader(status)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		fmt.Println(err)
	}
}
