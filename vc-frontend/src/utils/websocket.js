/*
* This file contains the logic for handling the sending/ receiving of websocket messages between the client 
* and server. This is present in the utils directory as it can be exported when needed to any of the components 
* and Vuex state can be modified in one place
*/
//import WebSocket from 'reconnecting-websocket';
import store from "../store"
const WS = new WebSocket('ws://localhost:5000/ws');


let msgQueue = [];

WS.onopen = function(){
    console.log("socket opened!");
    while (msgQueue.length > 0) {
        console.log('msg',msgQueue)
        WS.send(msgQueue.pop())
    }
}
WS.onmessage = function(msg){
    
        let apiMsg = JSON.parse(msg.data);
        console.log('received message',apiMsg)
        switch(apiMsg.type){
            case "chatMessage": {
                let msgData = JSON.parse(apiMsg.data);
                console.log("parsed: ",msgData)
                if (store.state.chatMessages[msgData.roomId]==undefined){
                    msgData.id=1
                } else {
                    msgData.id = (store.state.chatMessages[msgData.roomId]).length+1;
                }
                store.commit('addMessage',msgData)
                break;
            }
            case "gameInfo":{
                console.log('gameInfo',apiMsg)
                store.commit('updateGameInfo',apiMsg)
                store.commit('updateClientList',apiMsg)
                break;
            }
            case "performMove":{
                if(apiMsg.isValid){ //only if move is valid you perform commit
                    store.commit('performMove',apiMsg)
                }
                break;
            }
            case "clientList":{ //if new client joins or leaves room
                console.log('clientActivity',apiMsg)
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
export default WS;

export function sendJSONReq(socket,type,msg){
    if (!isOpen(socket)) return;
    console.log('executing 2',socket)
    console.log('executing 4',JSON.stringify(msg))
    socket.send(JSON.stringify({type:type,data:JSON.stringify(msg)})) //socket.send(json);
    
  }

export function createRoom(socket,roomId,username,standardFen,customMovePatterns){
    console.log('executing create',standardFen)
    if(Object.keys(customMovePatterns).length != 0){
        sendJSONReq(socket,'createRoom',{roomId:roomId, username:username, fen:standardFen, movePatterns:customMovePatterns});
    } else {
    sendJSONReq(socket,'createRoom',{roomId:roomId, username:username, fen:standardFen});
    }
}

export function joinRoom(socket,roomId,username){
    console.log('executing join')
    sendJSONReq(socket,'joinRoom',{roomId:roomId, username:username});
}

export function sendMessage(socket,json){
    console.log('executing chat')
    sendJSONReq(socket,'chatMessage',{message: json.message, username: json.username, roomId:json.roomId});
}
export function requestGameinfo(socket,roomId){
    console.log('executing gameinfo req')
    sendJSONReq(socket,'reqGameInfo',{roomId:roomId});
}

export function sendMoveInfo(socket,json){
    console.log('executing moveinfo send')
    sendJSONReq(socket,'performMove',{
        piece:json.piece, 
        roomId:json.roomId, 
        srcRow:json.srcRow-1,
        srcCol:json.srcCol-1,
        destRow:json.destRow-1,
        destCol:json.destCol-1,
        color: json.color,
        castle: json.castle? true : null
    });
}

function isOpen(ws) { return ws.readyState === ws.OPEN }
/*export function sendRequest(socket,type,msg){
    let json = JSON.stringify(msg);
    if (socket.readyState !== 1){
        msgQueue.push(json)
    }else{
        socket.send(json);
    }
}*/
