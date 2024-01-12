package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"varchess/internal/chesscore"

	"github.com/gorilla/websocket"
)



func handleUserConnect(c *Client, req Request){
	var params ParamsUserConnect
	var response ResponseUserConnect = ResponseUserConnect{}
	fmt.Println("we connected")
	if err := unmarshalParameters(req,&params, c); err!=nil{
		fmt.Println("err",err)
		return
	}
	fmt.Println(params)
	var gameID string
	if params.Type == "create" {
		gameID, _ = generateRandomString(8)
		newGame, err := chesscore.CreateGame(params.GameConfig)
		if err!=nil{
			closeConnection(c,websocket.CloseProtocolError,"Invalid Config")
			return
		} 
		c.ws.gameHubs[gameID] = NewGameHub(gameID,c.ws.destroy)
		c.ws.gameHubs[gameID].game = newGame
		go c.ws.gameHubs[gameID].run()
		response = ResponseUserConnect{
			Response: Response{
				Event : EventUserConnect,
				Success: true,
			},
			Result: ResultUserConnect{
				GameID: gameID,
				//NewUser: c.user.Username,
			},
		}
	} else if params.Type == "join"{
		gameID = params.GameID
		if !c.isValidGame(gameID){
			closeConnection(c,websocket.CloseProtocolError,"Invalid ID")
			return
		}
		response = ResponseUserConnect{
			Response: Response{
				Event : EventUserConnect,
				Success: true,
			},
			Result: ResultUserConnect{
				GameID: gameID,
				//NewUser: c.user.Username,
			},
		}
	} else{
		closeConnection(c,websocket.CloseProtocolError,"Invalid Connect Type")
		return
	}
	
	gameHub := c.ws.gameHubs[gameID]
	gameHub.register <- c
	c.hubs[gameID] = gameHub
	bytes, err := json.Marshal(response)
	if err!=nil{
		log.Printf("user connect : JSON marshal failed %v",err)
		closeConnection(c,websocket.CloseProtocolError,"Internal server error")
	}
	c.broadcastToMembers(gameID,bytes)
}

func (c *Client) isValidGame(gameId string) bool{
	_,ok:= c.ws.gameHubs[gameId]
	return ok
}

func handleChatMessage(c *Client, req Request){
	user := c.user
	if user.ID == "" {
		closeConnection(c, websocket.ClosePolicyViolation,"Unauthenticated")
		return
	}
	var params ParamsChat
	if err := unmarshalParameters(req,&params, c); err!=nil{
		return
	}
	bytes, err := json.Marshal(req)
	if err!=nil{
		log.Printf("user connect : JSON marshal failed %v",err)
		closeConnection(c,websocket.CloseProtocolError,"Internal server error")
	}
	if c.isValidGame(params.GameID){
		c.broadcastToMembers(params.GameID,bytes)
	} else {
		closeConnection(c,websocket.CloseProtocolError,"Invalid ID")
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