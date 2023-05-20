package server

import (
	"encoding/json"
	"net/http"
	"varchess/internal/logger"
	"varchess/internal/server/handler"
	"varchess/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

type ApiError struct {
	Error string
}

var RoomsMap = make(map[string]bool)
type apiFunction func(http.ResponseWriter, *http.Request) error

type Server struct {
	ListenAddr string
	Store      store.Storage
}

var l = logger.Get()

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
		ListenAddr: listenAddr,
		Store:      store,
	}
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

func (s *Server) Start(allowedOrigins string) {
	
	router := chi.NewRouter()
	//router.Use(setHeadersMiddleware(allowedOrigins))
	router.Use(middleware.Recoverer)
	router.Use(middleware.AllowContentType("application/json"))
	router.Use(cors.Default().Handler)
	router.Get("/health",makeHTTPHandleFunc(handler.HealthCheckHandler))
	//router.HandleFunc("/create-room", makeHTTPHandleFunc(s.createRoomHandler)).Methods("POST")
	//router.HandleFunc("/room-state", makeHTTPHandleFunc(s.roomStateHandler)).Methods("GET")
	
	//router.HandleFunc("/login", makeHTTPHandleFunc(s.authenticateUserHandler)).Methods("GET")
	//router.HandleFunc("/signup", makeHTTPHandleFunc(s.createAccountHandler)).Methods("POST")
	//router.HandleFunc("/possible-squares", makeHTTPHandleFunc(s.getPossibleSquares)).Methods("GET")

	//wsServer := NewWebsocketServer()
	//go wsServer.Run()
	//router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		//ServeWsHandler(wsServer, w, r)
	//})
	l.Fatal().Err(http.ListenAndServe(s.ListenAddr, requestLogger(router))).Msg("Varchess Server closed")
}

/*func setHeadersMiddleware(allowedOrigins string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Access-Control-Allow-Origin", allowedOrigins)
            next.ServeHTTP(w, r)
        })
    }
}*/