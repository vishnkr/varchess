package ws

import (
	"encoding/json"
	"varchess/internal/chesscore"
)



type eventType string

const (
	EventChatMessage eventType = "chat.message"
	EventUserConnect eventType = "game.connect_user"
	EventUserDisconnect eventType = "game.disconnect_user"
	EventGameResult eventType = "game.result"
	EventGameMakeMove eventType = "game.make_move"
	EventGameDrawOffer eventType = "game.offer_draw"
	EventGameDrawResult eventType = "game.draw_result"
	EventGameResign eventType = "game.resign"
	EventError eventType = "Error"
)

type Request struct {
	Event  string          `json:"event"`
	Params json.RawMessage `json:"params"`
   }
   
type Response struct {
	Event        eventType `json:"event"`
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error_message,omitempty"`
}

type RequestConnectUser struct {
	Event  string            `json:"event"`
	Params ParamsUserConnect `json:"params"`
}
type ParamsUserConnect struct{
	Type string `json:"connect_type"` // Create or Join
	SessionID string `json:"session_id"`
	ParamsUserConnectCreate
	ParamsUserConnectJoin
}

type ParamsUserConnectCreate struct{
	GameConfig chesscore.GameConfig `json:"game_config"`
}

type ParamsUserConnectJoin struct{
	GameID string `json:"game_id"`
}

type ParamsChat struct{
	GameID string `json:"game_id"`
	User string `json:"username"`
	Message string `json:"message"`
}

type ResultUserConnect struct {
	NewUser       string   `json:"new_user"`
	ExistingUsers []string `json:"existing_users"`
	GameID       string   `json:"game_id"`
}

type ResponseUserConnect struct {
	Response
	Result ResultUserConnect `json:"result"`
}


