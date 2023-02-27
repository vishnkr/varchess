package server

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
	"unicode"
	"varchess/pkg/game"

	"github.com/gorilla/websocket"
)

type MessageStruct struct {
	Type string `json:"type"`
	Data string `json:"data,omitempty"`
}

type UserRoomInfo struct {
	Username           string              `json:"username"`
	RoomId             string              `json:"roomId"`
	StartFEN           string              `json:"fen,omitempty"`
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

//used when player resigns or draw agreement occurs
type ResultMessage struct {
	Type   string `json:"type,omitempty"`
	RoomId string `json:"roomId"`
	Color  string `json:"color"`
	Result string `json:"result,omitempty"`
}

type Client struct {
	conn     *websocket.Conn
	mu       sync.Mutex
	wsServer *WsServer
	send     chan []byte
	roomId   string
	username string
}

//,roomId string, username string
func newClient(conn *websocket.Conn, wsServer *WsServer) *Client {
	return &Client{
		conn:     conn,
		wsServer: wsServer,
		send:     make(chan []byte, 256),
	}
}

func (c *Client) disconnect(Type string) {
	log.Println("disconnect called from", Type)
	c.conn.Close()
	c.wsServer.unregister <- c
}

const (
	writeWait = 30 * time.Second
	pongWait  = 40 * time.Second
	pingTime  = (pongWait * 9) / 10
)

func (c *Client) Read() {
	defer c.disconnect("Read")
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		if string(msg) == "pong" {
			continue
		}

		reqData := MessageStruct{}
		json.Unmarshal([]byte(msg), &reqData)
		log.Println("received data:", reqData)
		switch reqData.Type {
		case "createRoom", "joinRoom":
			userInfo := UserRoomInfo{}
			json.Unmarshal([]byte(reqData.Data), &userInfo)
			c.username = userInfo.Username
			if reqData.Type == "createRoom" {
				room := c.CreateRoom(userInfo.RoomId, userInfo.StartFEN)
				if len(userInfo.CustomMovePatterns) != 0 {
					room.Game.Board.CustomMovePatterns = userInfo.CustomMovePatterns
				}
			} else {
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

func (c *Client) Write() {
	ticker := time.NewTicker(pingTime)
	defer ticker.Stop()
	defer c.disconnect("write")
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
			if err := c.conn.WriteMessage(websocket.TextMessage, []byte("ping")); err != nil {
				return
			}
		}
	}
}

func (c *Client) ResultOffer(data *MessageStruct) {
	resultMessage := &ResultMessage{}
	json.Unmarshal([]byte(data.Data), &resultMessage)
	if data.Type == "resign" {
		if resultMessage.Color == "w" {
			resultMessage.Result = "black"
		} else if resultMessage.Color == "b" {
			resultMessage.Result = "white"
		}
		resultMessage.Type = "result"
		if message, err := json.Marshal(resultMessage); err == nil {
			RoomsMap[resultMessage.RoomId].BroadcastToMembers(message)
		}
	} else {
		fmt.Println("offer draw")
	}
}

func (c *Client) SendChatMessage(data *MessageStruct) {
	chatMessage := ChatMessage{}
	json.Unmarshal([]byte(data.Data), &chatMessage)
	if room, ok := RoomsMap[chatMessage.RoomId]; ok {
		if message, err := json.Marshal(data); err == nil {
			room.BroadcastToMembersExceptSender(message, c)
		}
	} else {
		message := MessageStruct{Type: "error", Data: "Room does not exist, connection expired"}
		if errMessage, err := json.Marshal(message); err == nil {
			c.send <- errMessage
		}
	}
}

func (c *Client) PerformMove(data *MessageStruct) {
	move := game.Move{}
	json.Unmarshal([]byte(data.Data), &(move))
	moveResp := &MoveResponse{Piece: move.PieceType, SrcRow: move.SrcRow, SrcCol: move.SrcCol, DestRow: move.DestRow, DestCol: move.DestCol}
	val, ok := game.StrToTypeMap[strings.ToLower(move.PieceType)]
	piece := game.Piece{}
	if !ok {
		piece.Type = game.Custom
		piece.CustomPieceName = move.PieceType
	} else {
		piece.Type = val
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
			//check for checkmates/check on opponents
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
		message := MessageStruct{Type: "error", Data: "Room does not exist, connection expired"}
		if errMessage, err := json.Marshal(message); err == nil {
			c.send <- errMessage
		}
	}
}
