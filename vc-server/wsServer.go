package main

import (
	"fmt"
	"log"
	"net/http"
)

type WsServer struct{
	clients map[*Client]bool //map clientId/username to Client which contains connection info
	register chan *Client
	unregister chan *Client
	broadcast chan []byte
}


func NewWebsocketServer() *WsServer{
	return &WsServer{
		clients: make(map[*Client]bool),
		register: make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

func (ws *WsServer) Run(){
	for{
		select{
		case client:= <-ws.register:
			ws.registerClient(client)
		
		case client:= <-ws.unregister:
			ws.unregisterClient(client)
		
		}

	}
}

func (server *WsServer) registerClient(client *Client) {
	server.clients[client] = true
	fmt.Println(RoomsMap)
}

func (server *WsServer) unregisterClient(client *Client) {
	var roomId =client.roomId
	if _, ok := RoomsMap[roomId]; ok {
		if _, ok := RoomsMap[roomId].Clients[client]; ok {
			fmt.Println(client.username,"was removed from room")
			delete(RoomsMap[roomId].Clients,client)
		}
		server.deleteEmptyRooms(client)
	}
	if _, ok := server.clients[client]; ok {
		fmt.Println(client.username,"was deleted")
		delete(server.clients, client)
	}
}

func (server *WsServer) deleteEmptyRooms(client *Client){
	client.mu.Lock()
	defer client.mu.Unlock()
	for id,_:= range RoomsMap{
		if(len(RoomsMap[id].Clients)==0){
			fmt.Println(id,"room was deleted since its empty")
			delete(RoomsMap,id)	
		}
	}
	
}

func ServeWsHandler(wsServer *WsServer,w http.ResponseWriter, r *http.Request){
	conn, err:= upgrader.Upgrade(w,r,nil)
	if err!=nil{
		log.Println(err)
		return
	}
	client:= newClient(conn,wsServer)
	go client.Write()
	go client.Read()
	wsServer.register <- client
	fmt.Println("New Client!")
}
