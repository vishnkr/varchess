/*
* This file contains the logic for handling the sending/ receiving of websocket messages between the client 
* and server. This is present in the utils directory as it can be exported when needed to any of the components 
* and Vuex state can be modified in one place
*/
//import WebSocket from 'reconnecting-websocket';
import store from "../store"
const server_host = process.env.VUE_APP_SERVER_WS;
const WS = new WebSocket(`${server_host}/ws`);



WS.onopen = function(){
    console.log("Socket opened");
}
WS.onerror = function(){
    console.log('got error')
    store.commit('websocketError','Connection to server could not be established! Try again soon!')
}
WS.onmessage = function(msg){
        if(msg.data=="ping"){
            WS.send("pong")
        } else {
            let apiMsg = JSON.parse(msg.data);
            switch(apiMsg.type){
                case "chatMessage": {
                    let msgData = JSON.parse(apiMsg.data);
                    if (store.state.chatMessages[msgData.roomId]==undefined){
                        msgData.id=1
                    } else {
                        msgData.id = (store.state.chatMessages[msgData.roomId]).length+1;
                    }
                    store.commit('addMessage',msgData)
                    break;
                }
                case "error":{
                    store.commit('websocketError',apiMsg.data)
                    break;
                }
                case "gameInfo":{
                    store.commit('updateGameInfo',apiMsg)
                    if(apiMsg.result){
                        store.commit("setResult",apiMsg.result)
                    }
                    break;
                }

                case "performMove":{
                    if(apiMsg.isValid){ //only if move is valid you perform commit
                        store.commit('performMove',apiMsg)
                    }
                    if(apiMsg.result){
                        store.commit('setResult',apiMsg.result)
                    } else if (apiMsg.check){
                        alert("Check!");
                    }
                    break;
                }
                case "result":{
                    if(apiMsg.result){
                        store.commit('setResult',apiMsg.result)
                    }
                    break;
                }
                case "clientList":{ //if new client joins or leaves room
                    var client
                    let msgData = JSON.parse(apiMsg.data);
                    if(msgData.activityType==="join"){
                        for(client of msgData.clientList){
                            if(!store.state.roomClients[apiMsg.roomId][client]){
                                store.commit('addClientToRoom',{roomId:msgData.roomId,username:client})
                            }
                        }
                    } else {
                        for(client of store.state.roomClients[msgData.roomId]){
                            if(!msgData.clientList.includes(client)){
                                store.commit('removeClientfromRoom',{roomId:msgData.roomId,username:client})
                            }
                        }
                    }
                    break;
                }
                default:
                    break;
            }
        }
        
    
}
export default WS;

export function sendJSONReq(socket,type,msg){
    if (!isOpen(socket)) return;
    socket.send(JSON.stringify({type:type,data:JSON.stringify(msg)})) //socket.send(json);
    
  }

export function createRoom(socket,roomId,username,standardFen){
    sendJSONReq(socket,'createRoom',{roomId:roomId, username:username, fen:standardFen});
}

export function createRoomWithCustomPatterns(socket,roomId,username,standardFen,customMovePatterns){
    sendJSONReq(socket,'createRoom',{roomId:roomId, username:username, fen:standardFen, movePatterns:customMovePatterns});
}

export function joinRoom(socket,roomId,username){
    sendJSONReq(socket,'joinRoom',{roomId:roomId, username:username});
}

export function sendMessage(socket,json){
    sendJSONReq(socket,'chatMessage',{message: json.message, username: json.username, roomId:json.roomId});
}
export function requestGameinfo(socket,roomId){
    sendJSONReq(socket,'reqGameInfo',{roomId:roomId});
}

export function sendResign(socket,data){
    sendJSONReq(socket,'resign',{roomId:data.roomId,color:data.color})
}

export function sendDrawOffer(socket,data){
    sendJSONReq(socket,'draw',{roomId:data.roomId,color:data.color})
}

export function sendMoveInfo(socket,json){
    if(!store.state.gameInfo.result){
        sendJSONReq(socket,'performMove',{
            piece:json.piece, 
            roomId:json.roomId, 
            srcRow:json.srcRow-1,
            srcCol:json.srcCol-1,
            destRow:json.destRow-1,
            destCol:json.destCol-1,
            color: json.color,
            castle: json.castle? true : false
        });
    }
}

function isOpen(ws) { return ws.readyState === ws.OPEN }

