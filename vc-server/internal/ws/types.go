package ws

import (
	"encoding/json"
	"varchess/internal/chesscore"
)



type eventType string

const (
	EventChatMessage eventType = "chat.message"
	EventUserConnect eventType = "game.connect_user"
	EventCreateGame eventType = "game.create_game"
	EventJoinGame eventType = "game.join_game"
	EventUserDisconnect eventType = "game.disconnect_user"
	EventSetPlayers eventType = "game.set_players"
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
	Params ParamsUserCreate `json:"params"`
}

type ParamsUserCreate struct{
	SessionID string `json:"session_id"`
	GameConfig chesscore.GameConfig `json:"game_config"`
	Color string `json:"color,omitempty"`
}

type ParamsUserJoin struct{
	GameID string `json:"game_id"`
}

type ParamsChat struct{
	GameID string `json:"game_id"`
	User string `json:"username"`
	Message string `json:"message"`
}

type ResultUserJoin struct {
	NewUser       string   `json:"new_user"`
	ExistingUsers []string `json:"existing_users"`
	GameID       string   `json:"game_id"`
}

type ResponseUserJoin struct {
	Response
	Result ResultUserJoin `json:"result"`
}

type ResultCreate struct {
	NewUser       string   `json:"new_user"`
	ExistingUsers []string `json:"existing_users"`
	GameID       string   `json:"game_id"`
}

type ResponseCreate struct {
	Response
	Result ResultCreate `json:"result"`
}