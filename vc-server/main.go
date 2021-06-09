package main

import (
	"log"
	"flag"
	"fmt"
	
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	//"github.com/gorilla/websocket"
)
var addr = flag.String("addr", ":5000", "http server address")

/*
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
*/

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan string)
var upgrader = websocket.Upgrader{
	ReadBufferSize: 4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		log.Println(origin)
		return origin == "http://localhost:8080"
	},
}


func main(){
	router := mux.NewRouter()
    router.HandleFunc("/getRoomId", roomHandler).Methods("POST")
	router.HandleFunc("/", rootHandler)
	router.HandleFunc("/ws", socketHandler)
	fmt.Print("listening on ", *addr,"\n")
	log.Fatal(http.ListenAndServe(*addr, router))
}

func roomHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json") 
	fmt.Print("reached here")
	response:= responseStruct{
		Type: "getRoomId",
		Data: genRandSeq(6),
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
	//http.HandleFunc("/ws", socketHandler)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "home")
}

/*
func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
			log.Fatal(err)
	}
	// register client
	fmt.Printf("reached")
	if (! clients[ws]){
		clients[ws] = true	
	}
	var message string
	if err := json.NewDecoder(r.Body).Decode(message); err!=nil{
		log.Printf("error %s",err)
		http.Error(w, "Bad request", http.StatusTeapot)
		return
	}
	defer r.Body.Close()
	go writer(message)

}*/

type responseStruct struct {
	Type string `json:"type"`
	Data string `json:"data,omitempty"`
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
    // Upgrade our raw HTTP connection to a websocket based one
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Print("Error during connection upgradation:", err)
        return
    }
    defer conn.Close()
	if (! clients[conn]){
		clients[conn] = true	
	}
    // The event loop
    for {
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            log.Println("Error during message reading:", err)
            break
        }
        log.Printf("Received: %s", message)
		responseData:=responseStruct{}
		json.Unmarshal([]byte(message),&responseData )
		fmt.Println("response: ",responseData)
		switch responseData.Type{
		case "createRoom":
			var roomId string
			roomId = responseData.Data
			fmt.Println("new ro ID:",roomId)
			c := &Client{conn: conn}
			c.CreateRoom(roomId)
			fmt.Println("created room",*c)
			err = conn.WriteMessage(messageType, []byte("successfully created"))
			if err != nil {
				log.Println("Error during message writing:", err)
				break
			}
		case "joinRoom":
			var roomId string
			roomId = responseData.Data
			c := &Client{conn: conn}
			c.AddtoRoom(roomId)
			err = conn.WriteMessage(messageType, []byte("successfully joined"))
			if err != nil {
				log.Println("Error during message writing:", err)
				break
			}
			fmt.Println(RoomList)
		}
			/*message,err = json.Marshal(responseData)
			if err != nil {
				panic(err)
			}*/
        
    }
}