
package main
import (
	"github.com/gorilla/websocket"
	"sync"
	//"net/http"
	"strings"
	"unicode"
	"time"
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
	c.conn.Close()
	c.wsServer.unregister <- c
}


const (
	writeWait = 10*time.Second
	pongWait = 40*time.Second
	pingTime = (pongWait*9)/10
)

func (c *Client) Read(){
	defer c.disconnect()
	for {
		_, msg, err:= c.conn.ReadMessage()
		if err!=nil{
			return
		}
		if (string(msg)=="pong"){
			c.conn.SetReadDeadline(time.Now().Add(pongWait))
			c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
			continue
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
					
				room, ok := RoomsMap[move.RoomId]
				if ok {
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
				} else {
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
	ticker := time.NewTicker(pingTime)
	defer ticker.Stop()
	defer c.disconnect()
	for {
		select {
		case msg,ok := <- c.send:
				if !ok{
					// The WsServer closed the channel.
					c.conn.WriteMessage(websocket.CloseMessage, []byte("closing"))
					return
				}
				err := c.conn.WriteMessage(websocket.TextMessage, msg) 
				if err != nil { 
				return 
				} 
			case <-ticker.C:
				c.conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := c.conn.WriteMessage(websocket.TextMessage, []byte("ping")); err != nil {
					
					return
				}
		}
	}

}
