package server

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
	"unicode"
	"varchess/internal/game"

	"github.com/gorilla/websocket"
)

type MessageStruct struct {
	Type string `json:"type"`
	Data json.RawMessage`json:"data,omitempty"`
}

type CreateRoomInfo struct {
	StartFEN           string              `json:"fen"`
	CustomMovePatterns []game.MovePatterns `json:"movePatterns,omitempty"`
}

type ChatMessage struct {
	RoomId   string `json:"roomId"`
	Message  string `json:"message"`
	Username string `json:"username"`
}

type Response struct {
	Status string `json:"status"`
}

type MoveResponse struct {
	Piece   string    `json:"piece"`
	SrcRow  int       `json:"srcRow"`
	SrcCol  int       `json:"srcCol"`
	DestRow int       `json:"destRow"`
	DestCol int       `json:"destCol"`
	Type    string    `json:"type"`
	Promote game.Type `json:"promote,omitempty"`
	Castle  bool      `json:"castle,omitempty"`
	IsValid bool      `json:"isValid,omitempty"`
	Check   bool      `json:"check,omitempty"`
	Result  string    `json:"result,omitempty"`
}

type ResultMessage struct {
	Type   string `json:"type,omitempty"`
	RoomId string `json:"roomId"`
	Color  string `json:"color"`
	Result string `json:"result,omitempty"`
}

type MembersList struct{
	P1 string `json:"p1"`
	P2 string `json:"p2"`
	Members []string `json:"members"`
}

type Client struct {
	conn     *websocket.Conn
	mu       sync.Mutex
	wsServer *WsServer
	send     chan []byte
	roomId   string
	username string
	disconnected sync.Once
}

func newClient(conn *websocket.Conn, wsServer *WsServer, roomId ,username string) *Client {
	client:= &Client{
		conn:     conn,
		wsServer: wsServer,
		send:     make(chan []byte, 256),
		roomId: roomId,
		username: username,
	}
	client.joinRoom(roomId,username)
	return client
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
        "chatMessage": c.handleChatMessage,
        "resign": c.handleResultOffer,
        "draw": c.handleResultOffer,
        "performMove": c.handlePerformMove,
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
			if err := c.conn.WriteMessage(websocket.PingMessage,nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) handleResultOffer(data *MessageStruct) {
	resultMessage := &ResultMessage{}
	json.Unmarshal([]byte(data.Data), &resultMessage)
	room:= RoomsMap[resultMessage.RoomId]
	switch data.Type {
	case "resign":
		if resultMessage.Color == "w" {
			resultMessage.Result = "black"
		} else if resultMessage.Color == "b" {
			resultMessage.Result = "white"
		}
		resultMessage.Type = "result"
		if message, err := json.Marshal(resultMessage); err == nil {
			room.BroadcastToMembers(message)
		}
	case "drawOffer":
        room.DrawOffer.IsOffered = true
        room.DrawOffer.Color = resultMessage.Color
        msg := fmt.Sprintf("%s has offered a draw", resultMessage.Color)
        if message, err := json.Marshal(MessageStruct{Type: "drawOffer", Data: json.RawMessage(msg)}); err == nil {
            room.BroadcastToMembers(message)
        }
	case "drawDecision":
		if room.DrawOffer.IsOffered{
			if resultMessage.Result == "accept"{
				resultMessage.Type = "result"
                resultMessage.Result = "draw"
                if message, err := json.Marshal(resultMessage); err == nil {
                    room.BroadcastToMembers(message)
                }
			}
		} else{
			msg := fmt.Sprintf("%s has declined the draw offer", resultMessage.Color)
			if message, err := json.Marshal(MessageStruct{Type: "drawDeclined", Data: json.RawMessage(msg)}); err == nil {
				room.BroadcastToMembers(message)
			}
		}
		room.DrawOffer.IsOffered = false
	}
}

func (c *Client) handleChatMessage(data *MessageStruct) {
	chatMessage := ChatMessage{}
	json.Unmarshal([]byte(data.Data), &chatMessage)
	if room, ok := RoomsMap[chatMessage.RoomId]; ok {
		if message, err := json.Marshal(data); err == nil {
			room.BroadcastToMembers(message)
		}
	} else {
		errorMsg := "Room does not exist, connection expired"
		errorMsgBytes, _ := json.Marshal(errorMsg)
		message := MessageStruct{Type: "error", Data: json.RawMessage(errorMsgBytes)}
		if errMessage, err := json.Marshal(message); err == nil {
			c.send <- errMessage
		}
	}
}

func (c *Client) handlePerformMove(data *MessageStruct) {
	move := game.Move{}
	json.Unmarshal([]byte(data.Data), &(move))
	moveResp := &MoveResponse{Piece: move.PieceType, SrcRow: move.SrcRow, SrcCol: move.SrcCol, DestRow: move.DestRow, DestCol: move.DestCol}
	val, ok := game.StrToTypeMap[strings.ToLower(move.PieceType)]
	piece := game.Piece{Type:val}
	if !ok {
		piece.Type = game.Custom
		piece.CustomPieceName = move.PieceType
	} 
	r := []rune(move.PieceType)
	if unicode.IsUpper(r[0]) {
		piece.Color = game.White
	} else {
		piece.Color = game.Black
	}

	room, ok := RoomsMap[move.RoomId]
	
	if ok {
		var res bool = false
		curGame := room.Game
		if curGame.Turn == piece.Color {
			validMoves := curGame.Board.GetAllValidMoves(curGame.Turn)
			for mv, movePiece := range validMoves {
				if piece == movePiece && game.IsSameMove(*mv, move) {
					res = true
				}
			}
		}
		
		if res {
			curGame.Board.PerformMove(piece, move)
			over, result := curGame.Board.IsGameOver(game.GetOpponentColor(piece.Color))
			if over {
				log.Println("game over")
				moveResp.Result = result
			} else {
				underCheck := curGame.Board.IsKingUnderCheck(game.GetOpponentColor(piece.Color))
				if underCheck {
					moveResp.Check = true
				}
			}
			moveResp.IsValid = res
			moveResp.Type = "performMove"
			if move.Castle {
				moveResp.Castle = true
			}
			if message, err := json.Marshal(moveResp); err == nil {
				room.BroadcastToMembers(message)
			}
			curGame.ChangeTurn()
		}
		response := Response{Status: "successful"}
		marshalledMessage, _ := json.Marshal(response)
		c.send <- marshalledMessage
	} else {
		errorMsg := "Room does not exist, connection expired"
		errorMsgBytes, _ := json.Marshal(errorMsg)
		message := MessageStruct{Type: "error", Data: json.RawMessage(errorMsgBytes)}
		if errMessage, err := json.Marshal(message); err == nil {
			c.send <- errMessage
		}
	}
}

func (c *Client) joinRoom(roomId, username string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	fmt.Println(roomId,username)
	curRoom, ok := RoomsMap[roomId]
	if ok {
		curRoom.Clients[c] = true
		switch len(curRoom.Clients) {
		case 1:
			curRoom.P1 = c
		case 2:
			curRoom.P2 = c
		}
		response := MembersList{
			Members: curRoom.getViewerClients(),
		}
		if curRoom.P1 != nil {
			response.P1 = curRoom.P1.username
		}
		if curRoom.P2 != nil {
			response.P2 = curRoom.P2.username
		}
		responseMsg,err:=json.Marshal(response)
		if err!=nil{
			fmt.Println(err.Error())
		}

		finalResp,_:=json.Marshal(MessageStruct{Type:"memberList",Data:json.RawMessage(responseMsg)})

		curRoom.BroadcastToMembers(finalResp)
	} else {
		log.Println("Room close")
		errorMsg := "Room does not exist, connection expired"
		errorMsgBytes, _ := json.Marshal(errorMsg)
		message := MessageStruct{Type: "error", Data: json.RawMessage(errorMsgBytes)}
		if errMessage, err := json.Marshal(message); err == nil {
			c.send <- errMessage
		}
	}
	return
}

func getMessageRawByteSlice(m interface{},mType string) ([]byte, error){
	data, err := json.Marshal(m)
	if err!=nil{
		return nil,err
	}
	message:= MessageStruct{
		Type: mType,
		Data: data,
	}
	var rawMessage json.RawMessage
	if rawMessage,err = json.Marshal(message); err == nil{
		return nil,err
	}
	return rawMessage,nil
}