package main

import (
	"log"
	"net/http"
	"vc-server/vc-ws/internal/config"
	"vc-server/vc-ws/internal/wss"

	"github.com/gorilla/mux"
)

func main() {
    r := mux.NewRouter()
    cfg, err := config.Load(".env")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
    wss.InitRedis(cfg)
    r.HandleFunc("/play/{gameID}", func(w http.ResponseWriter, r *http.Request) {
       wss.HandleWSConnection(w, r, cfg)
    })

    log.Println("WebSocket server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
