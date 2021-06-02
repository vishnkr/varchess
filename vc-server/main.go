package main

import (
	"log"
	"flag"
	"net/http"
	//"github.com/gorilla/websocket"
)
var addr = flag.String("addr", ":5000", "http server address")


func main(){
	flag.Parse()
	wsServer := NewWebsocketServer()
	go wsServer.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(wsServer,w, r)
	})
	log.Println("Serving",*addr)
	log.Fatal(http.ListenAndServe(*addr, nil))

    
    
}