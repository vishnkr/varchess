package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WsServer struct {
	clients    map[*Client]bool //map clientId/username to Client which contains connection info
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

func NewWebsocketServer() *WsServer {
	return &WsServer{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (ws *WsServer) Run() {
	for {
		select {
		case client := <-ws.register:
			ws.registerClient(client)

		case client := <-ws.unregister:
			ws.unregisterClient(client)
		}
	}
}

func (server *WsServer) registerClient(client *Client) {
	server.clients[client] = true
}

func (server *WsServer) unregisterClient(client *Client) {
	var roomId = client.roomId
	if _, ok := RoomsMap[roomId]; ok {
		if _, ok := RoomsMap[roomId].Clients[client]; ok {
			delete(RoomsMap[roomId].Clients, client)
		}
		server.deleteEmptyRooms(client)
	}
	if _, ok := server.clients[client]; ok {
		delete(server.clients, client)
	}
}

func (server *WsServer) deleteEmptyRooms(client *Client) {
	client.mu.Lock()
	defer client.mu.Unlock()
	for id := range RoomsMap {
		if len(RoomsMap[id].Clients) == 0 {
			log.Println(id, "room was deleted since its empty")
			delete(RoomsMap, id)
		}
	}

}

func ServeWsHandler(wsServer *WsServer, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	roomId := r.URL.Query().Get("roomId")
    username := r.URL.Query().Get("username")
	if err != nil {
		log.Println(err)
		return
	}
	client := newClient(conn, wsServer,roomId,username)
	go client.Write()
	go client.Read()
	wsServer.register <- client
	log.Println("New Client!")
}