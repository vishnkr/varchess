package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"varchess/pkg/auth"
	"varchess/pkg/server"

	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")
	if port==""{
		port = "5000"
	}
	var addr = flag.String("addr", ":"+port, "http server address")
	router := mux.NewRouter()
	router.HandleFunc("/getRoomId", server.RoomHandler).Methods("POST")
	router.HandleFunc("/getBoardFen/{roomId}", server.BoardStateHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/", rootHandler)
	router.HandleFunc("/login", auth.AuthUserHandler).Methods("GET")
	router.HandleFunc("/signup", auth.CreateAccountHandler).Methods("POST")
	router.HandleFunc("/getPossibleToSquares",server.GetPossibleSquares).Methods("POST")
	wsServer := server.NewWebsocketServer()
	go wsServer.Run()
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		server.ServeWsHandler(wsServer, w, r)
	})
	log.Print("listening on ", *addr, "\n")
	log.Fatal(http.ListenAndServe(*addr, router))
	//replStart()
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "home")
}
