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
	Data json.RawMessage`json:"data,omitempty"`
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
	disconnected sync.Once
}

func newClient(conn *websocket.Conn, wsServer *WsServer) *Client {
	return &Client{
		conn:     conn,
		wsServer: wsServer,
		send:     make(chan []byte, 256),
	}
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
        "createRoom": c.handleCreateRoom,
        "joinRoom": c.handleJoinRoom,
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
		log.Println("received data:", reqData)
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
			room.BroadcastToMembersExceptSender(message, c)
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
		errorMsg := "Room does not exist, connection expired"
		errorMsgBytes, _ := json.Marshal(errorMsg)
		message := MessageStruct{Type: "error", Data: json.RawMessage(errorMsgBytes)}
		if errMessage, err := json.Marshal(message); err == nil {
			c.send <- errMessage
		}
	}
}


func (c *Client) handleCreateRoom(data *MessageStruct){
	c.mu.Lock()
	defer c.mu.Unlock()
	userInfo := &UserRoomInfo{}
	err := json.Unmarshal([]byte(data.Data), &userInfo)
    if err != nil {
        log.Println("error unmarshalling createRoom data:", err)
        return
    }
	c.username = userInfo.Username
	c.roomId = userInfo.RoomId
	RoomsMap[userInfo.RoomId] = &Room{
		Game: &game.Game{
			Board: game.ConvertFENtoBoard(userInfo.StartFEN),
			Turn:  game.White,
		},
		Clients: make(map[*Client]bool),
		Id:      c.roomId,
		P1:      c,
		DrawOffer: DrawOffer{IsOffered: false},
	}
	game.DisplayBoardState(RoomsMap[userInfo.RoomId].Game.Board)
	RoomsMap[userInfo.RoomId].Clients[c] = true
	gameInfo := GameInfo{Type: "gameInfo", P1: c.username, Turn: "w", RoomId: userInfo.RoomId, Members: []string{}}
	gameInfo.Members = append(gameInfo.Members, c.username)
	marshalledInfo, _ := json.Marshal(gameInfo)
	RoomsMap[userInfo.RoomId].BroadcastToMembers(marshalledInfo)
	if len(userInfo.CustomMovePatterns) != 0 {
		RoomsMap[userInfo.RoomId].Game.Board.CustomMovePatterns = userInfo.CustomMovePatterns
	}
	return
}

func (c *Client) handleJoinRoom(data *MessageStruct) {
	c.mu.Lock()
	defer c.mu.Unlock()
	userInfo := &UserRoomInfo{}
	json.Unmarshal([]byte(data.Data), &userInfo)
	roomId:= userInfo.RoomId
	curRoom, ok := RoomsMap[roomId]
	if ok {
		var gameInfo GameInfo
		if len(curRoom.Clients) == 1 {
			RoomsMap[roomId].P2 = c
			gameInfo = GameInfo{Type: "gameInfo", P1: curRoom.P1.username, P2: c.username, Turn: curRoom.Game.Turn.String(), RoomId: roomId, Members: RoomsMap[roomId].getClientUsernames()}
		} else {
			gameInfo = GameInfo{Type: "gameInfo", P1: curRoom.P1.username, P2: curRoom.P2.username, Turn: curRoom.Game.Turn.String(), RoomId: roomId, Members: RoomsMap[roomId].getClientUsernames()}
		}
		gameInfo.Members = append(gameInfo.Members, c.username)
		RoomsMap[roomId].Clients[c] = true
		c.roomId = roomId
		marshalledInfo, _ := json.Marshal(gameInfo)
		RoomsMap[roomId].BroadcastToMembers(marshalledInfo)
	} else {
		log.Println("Room close")
		errorMsg := "Room does not exist, connection expired"
		errorMsgBytes, _ := json.Marshal(errorMsg)
		message := MessageStruct{Type: "error", Data: json.RawMessage(errorMsgBytes)}
		if errMessage, err := json.Marshal(message); err == nil {
			c.send <- errMessage
		}
	}
}