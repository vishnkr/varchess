package server

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/olahol/melody"
)

type websocketMessageType string

const (
	ChatMessage websocketMessageType = "ChatMessage"
	UserJoin websocketMessageType = "UserJoin"
	UserLeave websocketMessageType = "UserLeave"
	GameResult websocketMessageType = "GameResult"
	Error websocketMessageType = "Error"
)

type websocketMessage struct {
    Type websocketMessageType `json:"type"`
    Data websocketData `json:"data"`
}

type websocketData interface{}

type websocketError struct{
	Message string `json:"message"`
}

func writeWebSocketResponse(s *melody.Session, msgType websocketMessageType, data interface{}) error{
	msg, _ := json.Marshal(websocketMessage{Data: data, Type: msgType})
	return s.Write(msg)
}

func createMelodyForRooms(server *server) *melody.Melody{
	m:= melody.New()
	m.HandleMessage(wsMessageHandler(m))
	m.HandleConnect(func(s *melody.Session){
		member:= member{
			Session: s,
			IsHost: false,
		}
		roomId := chi.URLParam(s.Request,"roomId")
		username := chi.URLParam(s.Request,"username")
		if room,ok := server.rooms[roomId]; ok{
			if len(room.members)==0{
				member.IsHost=true
			}
			s.Keys = make(map[string]interface{})
			s.Keys["roomId"] = roomId
			s.Keys["username"] = username 
			server.rooms[roomId].members =  append(server.rooms[roomId].members,member)
			data := map[string]interface{}{
                "username": username,
				"id": len(room.members),
				"isHost": member.IsHost,
				"role": "Viewer",
				"members":room.getRoomMembers(),
            }
            message := websocketMessage{
                Type: UserJoin,
                Data: data,
            }
            msgBytes, err := json.Marshal(message)
            if err != nil {
                log.Printf("Error marshaling WebSocket message: %v", err)
                return
            }
			
			for _,member:= range server.rooms[roomId].members{
				member.Session.Write(msgBytes)
			}
			
		} else {
			//writeWebSocketResponse(s,Error,websocketError{Message: "Room does not exist"})
			msg,_ := json.Marshal(websocketMessage{Data: websocketError{Message: "Room does not exist"}, Type: Error})
			s.CloseWithMsg(msg)
		}
	})
	disconnectHandler := func(s *melody.Session){
		if roomId,ok := s.Keys["roomId"].(string); ok{
			if room,exists := server.rooms[roomId]; exists{
				for i := len(server.rooms[roomId].members) - 1; i >= 0; i-- {
					if server.rooms[roomId].members[i].Session == s {
						server.rooms[roomId].members = append(server.rooms[roomId].members[:i], server.rooms[roomId].members[i+1:]...)
					}
				}
			
				data := map[string]interface{}{
					"username": s.Keys["username"],
				}
				message := websocketMessage{
					Type: UserLeave,
					Data: data,
				}
				msgBytes, err := json.Marshal(message)
				if err != nil {
					log.Printf("Error marshaling WebSocket message: %v", err)
					return
				}
				for _,member:= range server.rooms[roomId].members{
					member.Session.Write(msgBytes)
				}
				if len(room.members)==0{
					delete(server.rooms,roomId)
					fmt.Println("deleted room",roomId,server.rooms)
				}
			}
		}
		
	}

	m.HandleDisconnect(disconnectHandler)
	return m
}

func wsMessageHandler(m *melody.Melody) func(s *melody.Session,msg []byte){
	return func (s *melody.Session,msg []byte){
			var wsMsg websocketMessage
			if err := json.Unmarshal(msg, &wsMsg); err != nil {
				log.Printf("Error unmarshaling JSON: %v", err)
				return
			}
			log.Printf("Received WebSocket message: Type=%s, Data=%s", wsMsg.Type, wsMsg.Data)
			/*m.BroadcastFilter(msg, func(q *melody.Session) bool {
			return q.Request.URL.Path == s.Request.URL.Path
			})*/
			switch wsMsg.Type{
			case ChatMessage:
				roomID, ok := s.Keys["roomId"].(string)
				if !ok {
					log.Println("Room ID not found in session keys")
					return
				}
				m.BroadcastFilter(msg, broadcastToRoom(roomID))
			
			}
	}
}

func broadcastToRoom(roomID string) func(q *melody.Session) bool{
    return func(q *melody.Session) bool {
        otherRoomID, ok := q.Keys["roomId"].(string)
        return ok && roomID == otherRoomID
    }
}

func broadcastToRoomExceptSender(roomID string,sender *melody.Session) func(q *melody.Session) bool {
    return func(q *melody.Session) bool {
        otherRoomID, ok := q.Keys["roomId"].(string)
        return ok && roomID == otherRoomID && q.Request.URL!=sender.Request.URL
    }
}