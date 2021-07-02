package main

import (
	"log"
	"flag"
	"fmt"
	"encoding/json"
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
var RoomsMap = make(map[string]*Room)

func main(){
	router := mux.NewRouter()
    router.HandleFunc("/getRoomId", roomHandler).Methods("POST")
	router.HandleFunc("/getBoardFen/{roomId}",boardStateHandler).Methods("GET","OPTIONS")
	router.HandleFunc("/", rootHandler)
	wsServer := NewWebsocketServer()
	go wsServer.Run()
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request){
		ServeWsHandler(wsServer,w,r)
	})
	fmt.Print("listening on ", *addr,"\n")
	log.Fatal(http.ListenAndServe(*addr, router))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "home")
}
type BoardState struct{
	Fen string
	RoomId string
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
    (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    (*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
func boardStateHandler(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (r.Method == "OPTIONS") {
        return
    } else {
        params:= mux.Vars(r)
	id:= params["roomId"]
	response:= BoardState{
		Fen: ConvertBoardtoFEN(RoomsMap[id].Game.Board),
		RoomId: id,
	}
	json.NewEncoder(w).Encode(response)
    }
	
}