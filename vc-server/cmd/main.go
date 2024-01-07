package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"varchess/internal/config"
	"varchess/internal/db"
	"varchess/internal/game"
	"varchess/internal/logger"
	mw "varchess/internal/middleware"
	"varchess/internal/template"
	"varchess/internal/utils"
	"varchess/internal/ws"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

func main(){
	l := logger.Get()
	cfg, err := config.Load(".env")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	conn, err := db.Connect(cfg.DB)
	if err != nil {
		log.Fatalf("Error connecting to db: %v", err)
	}
	
	defer conn.Close()
	router := chi.NewRouter()
	server := http.Server{
		Addr: cfg.ServerHost + ":" + cfg.ServerPort,
		Handler: serverHandler(router,l,conn),
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %s", err)
		}
	}()

	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Forced server shutdown: %s", err)
	}
}


func serverHandler(r chi.Router, l zerolog.Logger, db *db.Database) chi.Router{
	r.Use(mw.Cors())
	r.Use(mw.RequestLogger(l))
	gameRepository := game.NewRepository(db)
	gameService := game.NewService(gameRepository)
	templateRepository := template.NewRepository(db)
	templateService := template.NewService(templateRepository)
	websocket:= ws.NewWebSocket(gameService,templateService)
	websocket.RegisterHandlers(r)
	r.Get("/health",func (w http.ResponseWriter, _ *http.Request) {
		utils.WriteStatus(w, http.StatusOK, struct {
			Message string `json:"message"`
		}{Message: "health check OK"})
	})
	return r
}

