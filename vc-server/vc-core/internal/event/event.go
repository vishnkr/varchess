package event

import (
	"encoding/json"
	"vc-server/chess"
)

type EventType string

type Event struct {
	GameID string    `json:"gid"`
	UserID string    `json:"uid"`
	Type   EventType `json:"t"`
	Data   json.RawMessage    `json:"d"`
	Timestamp int64     `json:"ts"`
}

const SystemAction = "system_action"
const UserAction = "user_action"

const (
	Start EventType = "start" // one directional vc-core-> vc-ws

	// bi directional vc-ws <-> vc-core
	Move EventType = "move"
	Join EventType = "join"
	Resign EventType = "resign"
	DrawOffer EventType = "draw-offer"
	DrawAccept EventType = "draw-accept"
	DrawReject EventType = "draw-reject"
	GameOver EventType = "over"
)

type JoinPayload struct {
	Color string `json:"c"`
}

type MovePayload struct{
	Move chess.MoveJSON `json:"m"`
}
