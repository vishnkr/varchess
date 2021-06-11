package  main


type WsServer struct{
	clients map[*Client]bool
	register chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

type ResponseStruct struct {
	Type string `json:"type"`
	Data string `json:"data,omitempty"`
}

type ChatMessage struct{
	RoomId string `json:"roomId"`
	Message string `json:"message,omitempty"`
	Username string `json:"username"`
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


func SocketHandler(w http.ResponseWriter, r *http.Request) {
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
        fmt.Print("received a message",string(message))
		responseData:=ResponseStruct{}
		json.Unmarshal([]byte(message),&responseData )
		fmt.Println("response: ",responseData)
		var roomId string
		switch responseData.Type{
		case "createRoom":
			roomId = responseData.Data
			fmt.Println("new ro ID:",roomId)
			c := &Client{conn: conn}
			c.CreateRoom(roomId)
			fmt.Println("created room",*c)
			err = conn.WriteMessage(messageType, message)
			if err != nil {
				log.Println("Error during message writing:", err)
				break
			}
		case "joinRoom":
			roomId = responseData.Data
			c := &Client{conn: conn}
			c.AddtoRoom(roomId)
			fmt.Println("join room",*c)
			err = conn.WriteMessage(messageType, message)
			if err != nil {
				log.Println("Error during message writing:", err)
				break
			}
			fmt.Println(RoomList)
		case "chatMessage":
			messageData := ChatMessage{}
			json.Unmarshal([]byte(responseData.Data),&messageData)
			fmt.Println("messssss",messageData)
			for client, _ := range RoomList[messageData.RoomId].Clients {
				fmt.Print("chat clients",*client)
				if (client.conn!=conn){
					client.conn.WriteJSON(responseData)
					fmt.Print("success")
				}
			}
			
		}
        
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