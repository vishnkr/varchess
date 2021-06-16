
package main
import (
	"github.com/gorilla/websocket"
	"sync"
	//"net/http"
	//"log"
	"fmt"
	"encoding/json"
)

type MessageStruct struct {
	Type string `json:"type"`
	Data string `json:"data,omitempty"`
}

type UserRoomInfo struct{
	Username string `json:"username"`
	RoomId string `json:"roomId"`
}

type ChatMessage struct{
	RoomId string `json:"roomId"`
	Message string `json:"message"`
	Username string `json:"username"`
}


type Client struct{
	conn *websocket.Conn
	mu sync.Mutex
	wsServer *WsServer
	send  chan []byte
	roomId string
	username string
}
//,roomId string, username string
func newClient(conn *websocket.Conn, wsServer *WsServer) *Client{
	return &Client{
		conn: conn,
		wsServer: wsServer,
		send: make(chan []byte, 256),
	}
}

func (c *Client) Read(){
	defer c.conn.Close()
	for {
		_, msg, err:= c.conn.ReadMessage()
		if err!=nil{
			return
		}
		fmt.Print("received a message",string(msg))
		reqData:= MessageStruct{}
		json.Unmarshal([]byte(msg),&reqData)
		fmt.Println("received:",reqData.Type,string(reqData.Data))
		
		switch reqData.Type{
			case "createRoom", "joinRoom":
				userInfo := UserRoomInfo{}
				json.Unmarshal([]byte(string(reqData.Data)),&userInfo)
				fmt.Println("room ID, uname:",userInfo.RoomId, userInfo.Username)
				c.username = userInfo.Username
				if(reqData.Type=="createRoom"){
					c.CreateRoom(userInfo.RoomId)
					
				} else{
					c.AddtoRoom(userInfo.RoomId)
				}
				response:= []byte("successful")
				c.send <- response
			case "chatMessage":
				
			case "performMove":

		}

		//Rooms[c.roomId.]
	}
}


func (c *Client) Write(){
	defer c.conn.Close()
	for msg:= range c.send{
		err := c.conn.WriteMessage(websocket.TextMessage, msg) 
		if err != nil { 
		return 
		} 
	}
}

 /*
import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	writeWait = 10 * time.Second 	// Max wait time when writing message to peer
	pongWait = 60 * time.Second // Max time till next pong from peer
	pingPeriod = (pongWait * 9) / 10 // Send ping interval, must be less then pong wait time
	maxMessageSize = 10000 // Maximum message size allowed from peer.
)


func (client *Client) disconnect() {
	client.wsServer.unregister <- client
	/*
	for room := range client.rooms {
		room.unregister <- client
	}
	close(client.send)
	client.conn.Close()
}

func ServeWs(wsServer *WsServer,w http.ResponseWriter, r *http.Request){
	conn, err:= upgrader.Upgrade(w,r,nil)
	if err!=nil{
		log.Println(err)
		return
	}
	client:= newClient(conn,wsServer)
	go client.writeMessage()
	go client.readMessage()
	wsServer.register <- client
	fmt.Println("New Client joined the hub!")
	fmt.Println(client)
}


func (client *Client) readMessage(){
	defer func(){
		client.disconnect()
	}()

	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error { client.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for{
		_, jsonMessage,err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}
		client.wsServer.broadcast <- jsonMessage
	}

}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func (client *Client) writeMessage(){
	ticker := time.NewTicker(pingPeriod)
	defer func(){
		ticker.Stop()
		client.conn.Close()
	}()
	for {
		select{
		case message, ok:= <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok{
				// The WsServer closed the channel.
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err:=client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Attach queued chat messages to the current websocket message.
			n := len(client.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-client.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}


*/
