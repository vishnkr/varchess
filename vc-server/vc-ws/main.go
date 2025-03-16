package main

import (
	"log"
	"net/http"
	"vc-server/vc-ws/internal/wss"

	"github.com/gorilla/mux"
)

func main() {
    r := mux.NewRouter()
    wss.InitRedis()
    r.HandleFunc("/play/{gameID}", wss.HandleWSConnection)

    log.Println("WebSocket server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
