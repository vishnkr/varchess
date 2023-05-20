package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WebsocketServer struct {
	clients    map[*Client]bool //map clientId/username to Client which contains connection info
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

func NewWebsocketServer() *WebsocketServer {
	return &WebsocketServer{
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

func (ws *WebsocketServer) Run() {
	for {
		select {
		case client := <-ws.register:
			ws.registerClient(client)

		//case client := <-ws.unregister:
			//ws.unregisterClient(client)
		}
	}
}

func (server *WebsocketServer) registerClient(client *Client) {
	server.clients[client] = true
}
/*
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
*/
func ServeWsHandler(wsServer *WebsocketServer, w http.ResponseWriter, r *http.Request) {
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
	
}



type MessageStruct struct {
	Type string `json:"type"`
	Data json.RawMessage`json:"data,omitempty"`
}


type Client struct {
	conn     *websocket.Conn
	mu       sync.Mutex
	wsServer *WebsocketServer
	send     chan []byte
	roomId   string
	username string
	disconnected sync.Once
}
func (c *Client) disconnect(Type string) {
	log.Println("disconnect called from", Type, "for user",c.username)
	c.conn.Close()
	c.wsServer.unregister <- c
}

const (
	writeWait = 30 * time.Second
	pongWait  = 40 * time.Second
	pingTime  = (pongWait * 9) / 10
)

func (c *Client) Read() {
	defer c.disconnected.Do(func() { c.disconnect("Read") })
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	messageHandlers := map[string]func(*MessageStruct){
        /*"chatMessage": c.handleChatMessage,
        "resign": c.handleResultOffer,
        "draw": c.handleResultOffer,
        "performMove": c.handlePerformMove,*/
    }

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		reqData := MessageStruct{}
		json.Unmarshal([]byte(msg), &reqData)
		log.Println("received data:", reqData.Type,string(reqData.Data))
		if handler, ok := messageHandlers[reqData.Type]; ok {
            handler(&reqData)
        }
	}
}

func (c *Client) Write() {
	ticker := time.NewTicker(pingTime)
	defer ticker.Stop()
	defer c.disconnected.Do(func() { c.disconnect("Write") })
	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {

				c.conn.WriteMessage(websocket.CloseMessage, []byte("closing"))
				return
			}
			err := c.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage,nil); err != nil {
				return
			}
		}
	}
}

func newClient(conn *websocket.Conn, wsServer *WebsocketServer, roomId ,username string) *Client {
	client:= &Client{
		conn:     conn,
		wsServer: wsServer,
		send:     make(chan []byte, 256),
		roomId: roomId,
		username: username,
	}
	return client
}