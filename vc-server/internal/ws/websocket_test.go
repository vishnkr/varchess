package ws

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestHandleWebSocket(t *testing.T) {
	mux := http.NewServeMux()
	//mux.HandleFunc("/ws", HandleWebSocket)
	server := httptest.NewServer(mux)
	//roomId:= "345345"
	//username:= "Test"
	wsUrl := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	log.Println(wsUrl)
	t.Run("can establish a connection", func(t *testing.T) {
	c, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	assert.NoError(t, err)
	assert.NotNil(t, c)
 })
}


func TestWebSocketClientConnect(t *testing.T) {
	//ws := NewWebSocket()
	mux := http.NewServeMux()
	//mux.HandleFunc("/ws", HandleWSConnection)
	server := httptest.NewServer(mux)
	wsUrl := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	_, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err!=nil{
		t.Fatalf("could not establish connection: %v",err)
	}
	//request := RequestConnectUser
}