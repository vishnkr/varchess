package ws

import (
	"fmt"
	"net/http"
	"varchess/internal/game"
	"varchess/internal/template"
	"varchess/internal/user"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

type WebSocket struct{
	gameService game.Service
	templateService template.Service
	userService user.Service
	gameHubs map[string]*gameHub
	destroy chan string
}

func NewWebSocket(gameService game.Service, templateService template.Service, userService user.Service) *WebSocket{
	gameHubs := make(map[string]*gameHub)
	destroy := make(chan string)
	go handleDestroy(destroy,gameHubs)
	return &WebSocket{
		gameService, 
		templateService, 
		userService,
		gameHubs,
		destroy,
	}
}

func (ws *WebSocket) RegisterHandlers(r chi.Router){
	r.HandleFunc("/ws", ws.HandleWSConnection)
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)
   
func (ws *WebSocket) HandleWSConnection(w http.ResponseWriter, r *http.Request){
	//ctx := r.Context()
	//logger := logger.FromContext(ctx)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
	 fmt.Printf("handler: failed to upgrade connection: %v", err)
	 return
	}
	client := Client{
		ws: ws,
		hubs: make(map[string]*gameHub),
		conn: conn,
		send: make(chan []byte, 256),
	}
	go client.readPump()
	go client.writePump()
}

func handleDestroy(destroy chan string, gameHub map[string]*gameHub){
	for{
		gameId := <-destroy
		delete(gameHub,gameId)
	}
}