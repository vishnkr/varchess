package main

import (
	"flag"
	"os"
	"varchess/internal/logger"
	"varchess/internal/server"
	"varchess/internal/store"

	"github.com/joho/godotenv"
)

var l = logger.Get()

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		l.Error().Msg("Error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	var addr = flag.String("addr", ":"+port, "http server address")
	store, err := store.NewStore()
	if err != nil {
		l.Err(err).Msg("error connecting to database")
	}
	var allowedOrigins = os.Getenv("ALLOWED_ORIGINS")
	s := server.NewServer(*addr, store, allowedOrigins)
	l.Info().Str("port", port).Msgf("Starting Varchess Server")
	s.Start()
}
