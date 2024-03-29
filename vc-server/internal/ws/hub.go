package ws

import (
	"varchess/internal/chesscore"
)

type players struct{
	white *Client 
	black *Client 
}

type gameHub struct {
 gameId string
 game *chesscore.Game
 gameConfig string
 players players
 clients map[*Client]bool
 broadcast chan []byte
 register chan *Client
 unregister chan *Client
 destroy chan<- string
}

func NewGameHub(gameId string,destroy chan<- string) *gameHub{
	return &gameHub{
		gameId:    gameId,
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		destroy: 	destroy,
		players: players{},
	}
}

func (h *gameHub) getExistingUsers()[]string{
	users:=[]string{}
	for c:= range h.clients{
		users = append(users, c.user.Username)
	}
	return users
}

func (h *gameHub) setPlayer(c *Client, color string){
	if color=="w"{
		h.players.white = c
	} else{
		h.players.black = c
	}
}

func (h *gameHub) broadcastToPlayers(msg interface{}){
	sendMessage(h.players.white,msg)
	sendMessage(h.players.black,msg)
}

func (h *gameHub) broadcastToMembersExceptClient(c *Client, msg interface{}){
	for client := range h.clients{
		if client!=c{
			sendMessage(client,msg)
		}
	}
}

func (h *gameHub) run(){
	for {
		select{
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			delete(h.clients,client)
			close(client.send)
			if len(h.clients) == 0 {
				h.destroy <- h.gameId
			}
		case message := <-h.broadcast:
			for client := range h.clients{
				select {
				case client.send <- message:
				default:
				 close(client.send)
				 delete(h.clients, client)
				}
			}
		}
	}
}