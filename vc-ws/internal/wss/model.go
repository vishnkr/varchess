package wss

type WSMessage struct {
    Type    string      `json:"t"`
    Payload interface{} `json:"p"`
}

type MoveMessage struct {
    Move   string `json:"m"`
}

type ChatMessage struct {
    Sender  string `json:"sender"`
    Message string `json:"msg"`
}

type GameResult struct {
    Winner string `json:"winner"`
    Reason string `json:"reason"`
}
