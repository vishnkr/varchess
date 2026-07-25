package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	_ "vc-server/chess/variants" // register variant factory
	"vc-server/vc-core/internal/api"
	"vc-server/vc-core/internal/api/handlers"
	"vc-server/vc-core/internal/config"
	"vc-server/vc-core/internal/db"
	"vc-server/vc-core/internal/logger"
	"vc-server/vc-core/internal/middleware"
	mw "vc-server/vc-core/internal/middleware"
	"vc-server/vc-core/internal/worker"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main(){
	cfg, err := config.Load(".env")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	l := logger.New()
	dbConn, err := db.Connect(cfg)
	
	if err != nil {
		log.Fatalf("Error connecting to db: %v", err)
	}
	defer dbConn.Close()

	wp := worker.NewWorkerPool(10,dbConn)
	go dbConn.SetupRedisSubscriber(context.Background(),wp)
	router := chi.NewRouter()
	corsMiddleware := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:5173"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
        ExposedHeaders:   []string{"Link"},
        AllowCredentials: true,
        MaxAge:           300,
    })
	router.Use(corsMiddleware.Handler)
	router.Use(middleware.RequestLogger())
	router.Post("/signup", handlers.HandleSignup(dbConn))
	router.Post("/login", handlers.HandleLogin(dbConn))
	// Refresh must be unauthenticated — access token may already be expired.
	router.Post("/api/refresh", handlers.HandleRefresh(dbConn))

	// OAuth (optional — enabled when provider env vars are set)
	router.Get("/auth/oauth/providers", handlers.HandleOAuthProviders(cfg))
	router.Get("/auth/oauth/github", handlers.HandleOAuthStart(cfg, "github"))
	router.Get("/auth/oauth/github/callback", handlers.HandleOAuthCallback(cfg, dbConn, "github"))
	router.Get("/auth/oauth/google", handlers.HandleOAuthStart(cfg, "google"))
	router.Get("/auth/oauth/google/callback", handlers.HandleOAuthCallback(cfg, dbConn, "google"))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("reached")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	apiServer := api.NewAPI(&l, dbConn)
	router.Route("/api", func(r chi.Router) {
		r.Use(mw.AuthMiddleware)
		apiServer.RegisterHandlers(r, dbConn)
	})

	server := http.Server{
		Addr:    cfg.ServerHost + ":" + cfg.ServerPort,
		Handler: router,
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