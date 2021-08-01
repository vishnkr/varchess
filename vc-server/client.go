
package main
import (
	"github.com/gorilla/websocket"
	"sync"
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
	Check bool `json:"check,omitempty"`
	Result string `json:"result,omitempty"`
}
//used when player resigns or draw agreement occurs
type ResultMessage struct{
	Type string `json:"type,omitempty"`
	RoomId string `json:"roomId"`
	Color string `json:"color"`
	Result string `json:"result,omitempty"`
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

func (c *Client) disconnect(Type string) {
	fmt.Println("disconnect called from",Type)
	//close(c.send)
	c.conn.Close()
	c.wsServer.unregister <- c
}


const (
	writeWait = 30*time.Second
	pongWait = 40*time.Second
	pingTime = (pongWait*9)/10
)

func (c *Client) Read(){
	defer c.disconnect("Read")
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err:= c.conn.ReadMessage()
		if err!=nil{
			return
		}
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		if (string(msg)=="pong"){continue}
		
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
				c.SendChatMessage(&reqData)

			case "resign", "draw":
				c.ResultOffer(&reqData)

			case "performMove":
				c.PerformMove(&reqData)	
		}
	}
}


func (c *Client) Write(){
	ticker := time.NewTicker(pingTime)
	defer ticker.Stop()
	defer c.disconnect("write")
	for {
		select {
			case msg,ok := <- c.send:
				c.conn.SetWriteDeadline(time.Now().Add(writeWait))
				if !ok{
					// The WsServer closed the channel.
					c.conn.WriteMessage(websocket.CloseMessage, []byte("closing"))
					return
				}
				err := c.conn.WriteMessage(websocket.TextMessage, msg) 
				if err != nil { 
					fmt.Println("err2",err)
					return 
				} 
			case <-ticker.C:
				c.conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := c.conn.WriteMessage(websocket.TextMessage, []byte("ping")); err != nil {
					fmt.Println("err3",err)
					return
				}
		}
	}
}

func (c *Client) ResultOffer(data *MessageStruct){
	resultMessage:=&ResultMessage{}
	json.Unmarshal([]byte(data.Data),&resultMessage)
	if (data.Type=="resign"){
		if (resultMessage.Color=="w"){ resultMessage.Result="black"} else if (resultMessage.Color=="b"){ resultMessage.Result="white"}
		resultMessage.Type="result"
		if message,err:= json.Marshal(resultMessage); err==nil{
			RoomsMap[resultMessage.RoomId].BroadcastToMembers(message)
		}
	} else {
		fmt.Println("offer draw")
	}
}

func (c *Client) SendChatMessage(data *MessageStruct){
	chatMessage := ChatMessage{}
	json.Unmarshal([]byte(data.Data),&chatMessage)
	if room, ok := RoomsMap[chatMessage.RoomId]; ok {
		if message,err:= json.Marshal(data); err==nil{
			room.BroadcastToMembersExceptSender(message,c)
		}
	} else {
		message := MessageStruct{Type:"error",Data:"Room does not exist, connection expired"}
		if errMessage,err:= json.Marshal(message); err==nil{
			c.send <- errMessage
		}
	}		
}

func (c *Client) PerformMove(data *MessageStruct){
	move:= &Move{}
	json.Unmarshal([]byte(data.Data),&(move))
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
		var res bool
		var reason string
		game:= room.Game
		if (game.Turn==move.Color){
			res,reason =game.Board.isValidMove(piece,move)
		} else{ res,reason = false,"wrong color"}
		if (res) {
			game.Board.performMove(piece,move)
			//check for checkmates/check on opponents
			over,result:= game.Board.isGameOver(getOpponentColor(piece.Color))
			if over{
				moveResp.Result = result
			} else{ 
				underCheck,_ :=game.Board.isKingUnderCheck(getOpponentColor(piece.Color))
				if underCheck{
					moveResp.Check = true
				}								
			}
			moveResp.IsValid = res
			moveResp.Type = "performMove"
			if (move.Castle){
				moveResp.Castle = true
			}
			if message,err:= json.Marshal(moveResp); err==nil{
				room.BroadcastToMembers(message)
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