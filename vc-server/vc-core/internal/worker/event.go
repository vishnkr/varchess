package worker

import (
	"encoding/json"
	"fmt"
	"vc-server/vc-core/internal/chess"
)

type EventType string

type Event struct {
	GameID string    `json:"gid"`
	UserID string    `json:"uid"`
	Type   EventType `json:"t"`
	Data   json.RawMessage    `json:"d"`
	Timestamp int64     `json:"ts"`
}

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
	Move chess.Move `json:"m"`
}

func (w *Worker) processEvent(event Event) {
	switch event.Type {
	case Move:
		// move validation
		w.processMove(event)
	case Join:
		w.processJoin(event)
	case DrawAccept:
		// pdateGameState(event)
	case DrawOffer:
		// pdateGameState(event)
	case DrawReject:
		// pdateGameState(event)
	case Resign:
		// pdateGameState(event)
	case GameOver:
		// pdateGameState(event)
	default:
		fmt.Println("Unknown event type:", event.Type)
	}
}