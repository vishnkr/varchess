
package main
import (
	"github.com/gorilla/websocket"
	"sync"
	//"net/http"
	"unicode"
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
	StartFEN string `json:"fen,omitempty"`
}

type ChatMessage struct{
	RoomId string `json:"roomId"`
	Message string `json:"message"`
	Username string `json:"username"`
}

type MoveInfo struct{
	Type string `json:"piece"`
	RoomId string `json:"roomId"`
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

func (c *Client) disconnect() {
	//close(c.send)
	c.conn.Close()
	c.wsServer.unregister <- c
}

func (c *Client) Read(){
	defer c.disconnect()
	var response []byte
	for {
		_, msg, err:= c.conn.ReadMessage()
		if err!=nil{
			return
		}
		reqData:= MessageStruct{}
		json.Unmarshal([]byte(msg),&reqData)
		fmt.Println("REAChed",reqData)		
		switch reqData.Type{
			case "createRoom", "joinRoom":
				userInfo := UserRoomInfo{}
				json.Unmarshal([]byte(reqData.Data),&userInfo)
				c.username = userInfo.Username
				if (reqData.Type=="createRoom"){
					c.CreateRoom(userInfo.RoomId,userInfo.StartFEN)
				} else{
					c.AddtoRoom(userInfo.RoomId)
				}
				response = []byte("successful")

			case "chatMessage":
				chatMessage := ChatMessage{}
				json.Unmarshal([]byte(reqData.Data),&chatMessage)
				marshalledMessage,_ := json.Marshal(chatMessage)
				fmt.Println(reqData.Data,chatMessage,"reachin here",chatMessage.RoomId)
				for member, _ := range RoomsMap[chatMessage.RoomId].Clients {
					fmt.Println("reachin here")
					if (member.conn!=c.conn){
						member.conn.WriteJSON(reqData)
						fmt.Println("success")
					}
				}
				response = marshalledMessage

			case "performMove":
				moveInfo := MoveInfo{}
				move:= &Move{}
				fmt.Println("REACHED")
				json.Unmarshal([]byte(reqData.Data),&(move))
				json.Unmarshal([]byte(reqData.Data),&(moveInfo))
				
				piece:=&Piece{Type:strToTypeMap[moveInfo.Type]}
				r := []rune(moveInfo.Type)
				if (unicode.IsUpper(r[0])){
					piece.Color = White
				} else { piece.Color = Black }
				game:= RoomsMap[moveInfo.RoomId].Game
				res,reason:=game.Board.isValidMove(piece,move)
				fmt.Println("move valid:",res,reason)
				response = []byte("comple")
		}
		c.send <- response
	}
}


func (c *Client) Write(){
	defer c.disconnect()
	for {
		msg,ok := <- c.send
		if !ok{
			// The WsServer closed the channel.
			c.conn.WriteMessage(websocket.CloseMessage, []byte("closing"))
			return
		}
		err := c.conn.WriteMessage(websocket.TextMessage, msg) 
		if err != nil { 
		return 
		} 
	}
}

 /*
const (
	writeWait = 10 * time.Second 	// Max wait time when writing message to peer
	pongWait = 60 * time.Second // Max time till next pong from peer
	pingPeriod = (pongWait * 9) / 10 // Send ping interval, must be less then pong wait time
	maxMessageSize = 10000 // Maximum message size allowed from peer.
)


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
