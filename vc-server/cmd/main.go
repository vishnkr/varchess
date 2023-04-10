package main

import (
	"flag"
	"log"
	"os"
	"varchess/internal/server"
	"varchess/internal/store"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	var addr = flag.String("addr", ":"+port, "http server address")
	store, err := store.NewStore()
	if err != nil {
		log.Fatal(err)
	}
	server := server.NewServer(*addr, store)
	log.Print("listening on ", *addr, "\n")
	var allowedOrigins = os.Getenv("ALLOWED_ORIGINS")
	log.Fatal(server.Start(allowedOrigins))
}
