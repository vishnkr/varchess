package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"varchess/internal/chesscore"
	"varchess/internal/user"

	"github.com/gorilla/websocket"
)

func connectClientToHub(gh *gameHub,c *Client){
	gh.register <- c
	c.hubs[gh.gameId] = gh
}

func handleCreateGame(c *Client,req Request){
	var params ParamsUserCreate

	fmt.Println("we connected")
	if err := unmarshalParameters(req,&params, c); err!=nil{
		fmt.Println("err",err)
		return
	}
	fmt.Println(params)
	var gameID string

	gameID, _ = generateRandomString(8)
	//todo: init new game
	newGame, err := chesscore.CreateGame(params.GameConfig)
	if err!=nil{
		closeConnection(c,websocket.CloseProtocolError,"Invalid Config")
		return
	}
	c.ws.gameHubs[gameID] = NewGameHub(gameID,c.ws.destroy)
	c.ws.gameHubs[gameID].game = newGame
	c.user = &user.User{Username: params.Username, ID: params.SessionID}
	go c.ws.gameHubs[gameID].run()
	connectClientToHub(c.ws.gameHubs[gameID],c)

	if params.Color=="w"{
		c.ws.gameHubs[gameID].setPlayer(c,"w")
	} else if params.Color=="b"{
		c.ws.gameHubs[gameID].setPlayer(c,"b")
	} // else throw some error 
	//timer1 := time.NewTimer(120 * time.Second)

	response := ResponseCreate{
		Response: Response{
			Event : EventUserConnect,
			Success: true,
		},
		Result: ResultCreate{
			GameID: gameID,
		},
	}
	sendMessage(c,response)
}

func handleJoinGame(c *Client,req Request){
	
	var params ParamsUserJoin
	if err := unmarshalParameters(req,&params, c); err!=nil{
		fmt.Println("err",err)
		return
	}
	gameID := params.GameID
	fmt.Println("joining",gameID)
	gameHub,ok := c.getGame(gameID)
	if !ok{
		fmt.Println("invald id")
		closeConnection(c,websocket.CloseProtocolError,"Invalid ID")
		return
	}
	
	existing:= gameHub.getExistingUsers()
	c.user = &user.User{Username: params.Username, ID: params.SessionID}
	connectClientToHub(gameHub,c)

	gameConfig,_ := gameHub.game.GetGameConfig()
	fmt.Println(gameConfig)
	if gameHub.players.white==nil  {
		gameHub.setPlayer(c,"w")
	} else if gameHub.players.black==nil{
		gameHub.setPlayer(c,"b")
	}  /*  else user is just a viewer }*/
	
	response := ResponseUserJoin{
		Response: Response{
			Event : EventUserConnect,
			Success: true,
		},
		Result: ResultUserJoin{
			GameID: gameID,
			GameConfig: gameConfig,
			ExistingUsers: existing,
			NewUser: c.user.Username,
		},
	}
	fmt.Println(response)
	sendMessage(c,response)

	response.Response.Event = EventJoinGame
	gameHub.broadcastToMembersExceptClient(c,response)
	startGame := StartGame{
		Response{
			Event: EventStartGame, 
			Success: true,
		},
		ResultStartGame{
			Players{
				White: gameHub.players.white.user.Username,
				Black:  gameHub.players.black.user.Username,
			},
			gameConfig,
		},
	}
	gameHub.broadcastToPlayers(startGame)
}

func (c *Client) getGame(gameId string) (*gameHub,bool){
	hub,ok:= c.ws.gameHubs[gameId]
	return hub,ok
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
	response := ResponseChat{
		Response: Response{
			Event: EventChatMessage,
			Success: true,
		},
		Result: ResultChat{
			Username: user.Username,
			Message: params.Message,
		},
	}
	bytes, err := json.Marshal(response)
	if err!=nil{
		log.Printf("user connect : JSON marshal failed %v",err)
		closeConnection(c,websocket.CloseProtocolError,"Internal server error")
	}
	_,ok := c.getGame(params.GameID)
	if ok{
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