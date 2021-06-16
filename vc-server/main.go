package main

import (
	"log"
	"flag"
	"fmt"
	//"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)
var addr = flag.String("addr", ":5000", "http server address")

var upgrader = websocket.Upgrader{
	ReadBufferSize: 4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		//return origin == "http://localhost:8080"
		return true
	},
}

func main(){
	router := mux.NewRouter()
    router.HandleFunc("/getRoomId", roomHandler).Methods("POST")
	router.HandleFunc("/", rootHandler)
	wsServer := NewWebsocketServer()
	go wsServer.Run()
	//router.HandleFunc("/ws", SocketHandler)
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request){
		ServeWsHandler(wsServer,w,r)
	})
	fmt.Print("listening on ", *addr,"\n")
	log.Fatal(http.ListenAndServe(*addr, router))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "home")
}

