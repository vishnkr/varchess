package ws

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)



func handleUserConnect(c *Client, req Request){
	var params ParamsUserConnect
	if err := unmarshalParameters(req,&params, c); err!=nil{
		return
	}
	gameID := params.GameID
	if _,ok := c.ws.gameHubs[gameID]; !ok{
		c.ws.gameHubs[gameID] = NewGameHub(gameID,c.ws.destroy)
		go c.ws.gameHubs[gameID].run()
	}
	gameHub := c.ws.gameHubs[gameID]
	gameHub.register <- c

	response := ResponseUserConnect{
		Response: Response{
			Event : EventUserConnect,
			Success: true,
		},
		Result: ResultUserConnect{
			GameID: gameID,
			NewUser: c.user.ID,
		},
	}
	bytes, err := json.Marshal(response)
	if err!=nil{
		log.Printf("user connect : JSON marshal failed %v",err)
		closeConnection(c,websocket.CloseProtocolError,"Internal server error")
	}
	c.ws.gameHubs[gameID].broadcast <- bytes
}

func handleChatMessage(c *Client, req Request){
	user := c.user
	if user.ID == "" {
		closeConnection(c, websocket.ClosePolicyViolation,"Unauthenticated")
		return
	}


}

func handleMakeMove(c *Client, req Request){}
func handleOfferDraw(c *Client, req Request){}
func handleDrawResult(c *Client, req Request){}
func handleResign(c *Client, req Request){}

func unmarshalParameters(req Request, v any, c *Client) error{
	err := json.Unmarshal(req.Params, v)
	if err != nil {
		closeConnection(c, websocket.CloseInvalidFramePayloadData, "Bad Parameters")
		return err
	}
	return nil
}


func closeConnection(c *Client, status int, msg string){
	err:= c.conn.WriteControl(websocket.CloseMessage,
		[]byte(websocket.FormatCloseMessage(status,msg)),time.Now().Add(writeWait))
	if err!=nil{
		log.Printf("failed to write control: %v", err)
	}
}