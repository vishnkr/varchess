
package main
import (
	"github.com/gorilla/websocket"
	"sync"
	//"net/http"
	"strings"
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
	CustomMovePatterns []MovePatterns `json:"movePatterns,omitempty"`
}

type ChatMessage struct{
	RoomId string `json:"roomId"`
	Message string `json:"message"`
	Username string `json:"username"`
}


type Response struct{
	Status string `json:"status"`
}

type MoveResponse struct{
	Piece string `json:"piece"`
	SrcRow int `json:"srcRow"`
	SrcCol int `json:"srcCol"`
	DestRow int `json:"destRow"`
	DestCol int `json:"destCol"`
	Type string `json:"type"`
	Promote Type `json:"promote,omitempty"`
	Castle bool `json:"castle,omitempty"`
	IsValid bool `json:"isValid,omitempty"`
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
	//var response []byte
	for {
		_, msg, err:= c.conn.ReadMessage()
		if err!=nil{
			return
		}
		reqData:= MessageStruct{}
		json.Unmarshal([]byte(msg),&reqData)
		fmt.Println("received data:",reqData)		
		switch reqData.Type{
			case "createRoom", "joinRoom":
				userInfo := UserRoomInfo{}
				json.Unmarshal([]byte(reqData.Data),&userInfo)
				c.username = userInfo.Username
				if (reqData.Type=="createRoom"){
					room:= c.CreateRoom(userInfo.RoomId,userInfo.StartFEN)
					if (len(userInfo.CustomMovePatterns)!=0){
						room.Game.Board.CustomMovePatterns = userInfo.CustomMovePatterns
					}
				} else{
					c.AddtoRoom(userInfo.RoomId)
				}
				
			case "chatMessage":
				chatMessage := ChatMessage{}
				json.Unmarshal([]byte(reqData.Data),&chatMessage)
				fmt.Println(reqData.Data,chatMessage,"reachin here",chatMessage.RoomId)
				if room, ok := RoomsMap[chatMessage.RoomId]; ok {
					for member, _ := range room.Clients {
						if (member.conn!=c.conn){
							if message,err:= json.Marshal(reqData); err==nil{
								member.send <- message
							}
						}
					}
				} else {
					fmt.Println("Room close")
					message := MessageStruct{Type:"error",Data:"Room does not exist, connection expired"}
					if errMessage,err:= json.Marshal(message); err==nil{
						c.send <- errMessage
					}
				}		
				
			case "performMove":
				move:= &Move{}
				json.Unmarshal([]byte(reqData.Data),&(move))
				fmt.Println("movedata",reqData.Data)
				moveResp:=&MoveResponse{Piece:move.PieceType,SrcRow:move.SrcRow,SrcCol:move.SrcCol,DestRow:move.DestRow,DestCol:move.DestCol}
				val,ok := strToTypeMap[strings.ToLower(move.PieceType)]
				piece:=&Piece{}
				if (!ok){
					piece.Type = Custom
					piece.CustomPiece = &CustomPiece{PieceName:move.PieceType}
				} else {
					piece.Type = val
				}
				r := []rune(move.PieceType)
				if (unicode.IsUpper(r[0])){
					piece.Color = White
				} else { piece.Color = Black }
				if room, ok := RoomsMap[move.RoomId]; ok {
					game:= room.Game
					var res bool
					var reason string
					fmt.Println(move.Color,game.Turn)
					if (game.Turn==move.Color){
						res,reason=game.Board.isValidMove(piece,move)
					} else{ res,reason = false,"wrong color"}
					if (res) {
						game.Board.performMove(piece,move)
						moveResp.IsValid = res
						moveResp.Type = "performMove"
						if (move.Castle){
							moveResp.Castle = true
						}
						for member, _ := range room.Clients {
							if message,err:= json.Marshal(moveResp); err==nil{
								member.send <- message
							}
						}
						game.Turn = changeTurn(game.Turn)
					}
					fmt.Println("move valid:",res,reason)
					response:= Response{Status:"successful"}
					marshalledMessage,_ := json.Marshal(response)
					c.send <- marshalledMessage
				} 
				if !ok {
					fmt.Println("Room close")
					message := MessageStruct{Type:"error",Data:"Room does not exist, connection expired"}
					if errMessage,err:= json.Marshal(message); err==nil{
						c.send <- errMessage
					}
				}	
				
		}
		
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
