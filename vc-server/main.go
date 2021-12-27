package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
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
	
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	
	router := mux.NewRouter()
    router.HandleFunc("/getRoomId", roomHandler).Methods("POST")
	router.HandleFunc("/getBoardFen/{roomId}",boardStateHandler).Methods("GET","OPTIONS")
	router.HandleFunc("/", rootHandler)
	router.HandleFunc("/login", AuthUserHandler).Methods("GET")
	router.HandleFunc("/signup", CreateAccountHandler).Methods("POST")
	wsServer := NewWebsocketServer()
	go wsServer.Run()
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request){
		ServeWsHandler(wsServer,w,r)
	})
	log.Print("listening on ", *addr,"\n")
	log.Fatal(http.ListenAndServe(*addr, router))
	//replStart()
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
return b / 1024 / 1024
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "home")
}
type BoardState struct{
	Fen string
	RoomId string
	MovePatterns []MovePatterns `json:"movePatterns"`
}

func boardStateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") 
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if (r.Method == "OPTIONS") {
        return
    } else {
        params:= mux.Vars(r)
	id:= params["roomId"]
	room, ok := RoomsMap[id] 
	if ok {
		response:= BoardState{
			Fen: ConvertBoardtoFEN(room.Game.Board),
			RoomId: id,
		}
		if (room.Game.Board.CustomMovePatterns!=nil){
			response.MovePatterns = room.Game.Board.CustomMovePatterns
		}
		json.NewEncoder(w).Encode(response)
	} else { 
		errResponse:= MessageStruct{Type:"error",Data:"Room does not exist/has been closed"}
		json.NewEncoder(w).Encode(errResponse)
	}
    }
}