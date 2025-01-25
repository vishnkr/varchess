package wss

import (
	"encoding/json"
	"fmt"
	"sync"
)

type GameHub struct {
    games map[string]*Game
    mu    sync.Mutex
}

var hub = &GameHub{games: make(map[string]*Game)}

type Game struct {
    gameID string
    mu     sync.Mutex
    players map[string]*Client
    //state   *GameState        
}

func NewGame(gameID string) *Game {
    return &Game{
        gameID:  gameID,
        players: make(map[string]*Client),
        //state:   LoadGameState(gameID),
    }
}

func (h *GameHub) GetGame(gameID string) (*Game,error){
    h.mu.Lock()
    defer h.mu.Unlock()

    game, exists := h.games[gameID]; if !exists {
        return nil,fmt.Errorf("Game not found")
    }
    return game,nil
}

func (h *GameHub) GetOrCreateGame(gameID string) *Game {
    h.mu.Lock()
    defer h.mu.Unlock()

    game, exists := h.games[gameID]
    if !exists {
        game = NewGame(gameID)
        h.games[gameID] = game
    }
    return game
}


func (h *GameHub) RemoveGame(gameID string) {
    h.mu.Lock()
    delete(h.games, gameID)
    h.mu.Unlock()
}

func (r *Game) HandleMessage(client *Client, msg []byte) {
    var wsMsg WSMessage
    json.Unmarshal(msg, &wsMsg)

    switch wsMsg.Type {
    case "move":
        /*var move Move
        json.Unmarshal(wsMsg.Payload, &move)
        r.gameState.Moves = append(r.gameState.Moves, move) 
       
        hub.SaveGameStateToRedis(r.gameID, r.gameState)
        */
        //r.Broadcast(msg)

    case "chat":
        //r.Broadcast(msg)
    case "resign":
       //r.Broadcast(msg)
    case "draw_offer":
        //r.Broadcast(msg)
    }
}




func (r *Game) UnregisterClient(c *Client) {
    if _, ok := r.players["white"]; ok && r.players["white"] == c {
        delete(r.players, "white")
    }
    if _, ok := r.players["black"]; ok && r.players["black"] == c {
        delete(r.players, "black")
    }

    if len(r.players) == 0{
        delete(hub.games,r.gameID)
    }
}

