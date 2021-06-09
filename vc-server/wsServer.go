package  main


type WsServer struct{
	clients map[*Client]bool
	register chan *Client
	unregister chan *Client
	broadcast  chan []byte
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
		
		case message := <-ws.broadcast:
			ws.broadcastToClients(message)
		}

	}
}

func (server *WsServer) registerClient(client *Client) {
	server.clients[client] = true
}

func (server *WsServer) unregisterClient(client *Client) {
	if _, ok := server.clients[client]; ok {
		delete(server.clients, client)
	}
}

func (server *WsServer) broadcastToClients(message []byte) {
	for client := range server.clients {
		client.send <- message
	}
}


/*
package
type WsServer struct{
	clients map[*Client]bool
	register chan *Client
	unregister chan *Client
	broadcast  chan []byte
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
		
		case message := <-ws.broadcast:
			ws.broadcastToClients(message)
		}

	}
}

func (server *WsServer) registerClient(client *Client) {
	server.clients[client] = true
}

func (server *WsServer) unregisterClient(client *Client) {
	if _, ok := server.clients[client]; ok {
		delete(server.clients, client)
	}
}

func (server *WsServer) broadcastToClients(message []byte) {
	for client := range server.clients {
		client.send <- message
	}
}

*/