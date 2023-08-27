package main

import (
	"fmt"
	"os"
	"varchess/internal/logger"
	"varchess/internal/server"
	"varchess/internal/store"

	"github.com/joho/godotenv"
)

var l = logger.Get()

func main() {
	fmt.Println("started")
	err := godotenv.Load(".env")
	if err != nil {
		l.Error().Msg("Error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	fmt.Println("running on",port)
	//var addr = flag.String("addr", ":"+port, "http server address")
	store, err := store.NewStore()
	if err != nil {
		l.Err(err).Msg("error connecting to database")
	}
	var allowedOrigins = os.Getenv("ALLOWED_ORIGINS")
	s := server.NewServer("0.0.0.0:5000", store, allowedOrigins)
	l.Info().Str("port", port).Msgf("Starting Varchess Server")
	s.Start()
}
