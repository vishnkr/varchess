package wss

type Event struct {
	GameID string    `json:"gid"`
	UserID string    `json:"uid"`
	Type   EventType `json:"t"`
	Data   interface{}    `json:"d"`
}
type EventType string

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

type WSMessage struct {
    Type    string      `json:"t"`
    Payload interface{} `json:"p"`
}

type MovePayload struct {
    Move string `json:"m"`
}

type ChatMessage struct {
	GameId string `json:"gid"`
    Sender  string `json:"sender"`
    Message string `json:"msg"`
}

type GameResult struct {
    Winner string `json:"winner"`
    Reason string `json:"reason"`
}
