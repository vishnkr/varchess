package ws

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"varchess/internal/game"
	"varchess/internal/template"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

type WebSocket struct{
	gameService game.Service
	templateService template.Service
	//userService user.Service
	gameHubs map[string]*gameHub
	destroy chan string
}

func NewWebSocket(gameService game.Service, templateService template.Service) *WebSocket{
	gameHubs := make(map[string]*gameHub)
	destroy := make(chan string)
	go handleDestroy(destroy,gameHubs)
	return &WebSocket{
		gameService, 
		templateService, 
		//userService,
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

func sendErrorMessage(c *Client, response interface{}) {
	responseBytes, err := json.Marshal(response)
	if err != nil {
		log.Printf("Failed to marshal response into JSON: %v", err)
		closeConnection(c, websocket.CloseProtocolError, "internal server error")
		return
	}
	c.send <- responseBytes
}

func generateRandomString(length int) (string, error) {
    bytes := make([]byte, length)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

func handleDestroy(destroy chan string, gameHub map[string]*gameHub){
	for{
		gameId := <-destroy
		delete(gameHub,gameId)
	}
}