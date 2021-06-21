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
        if(apiMsg.type==="chatMessage"){
            let msgData = JSON.parse(apiMsg.data);
            console.log("parsed: ",msgData)
            if (store.state.chatMessages[msgData.roomId]==undefined){
                msgData.id=1
            } else {
                msgData.id = (store.state.chatMessages[msgData.roomId]).length+1;
            }
            store.commit('addMessage',msgData)
        } else if(apiMsg.type==="gameInfo"){
            //let msgData = JSON.parse(apiMsg);
            console.log('gameInfo',apiMsg)
            store.commit('updateGameInfo',apiMsg)

        }
        
    
}
export default WS;

export function sendJSONReq(socket,type,msg){
    if (!isOpen(socket)) return;
    console.log('executing 2',socket)
    console.log('executing 4',JSON.stringify(msg))
    socket.send(JSON.stringify({type:type,data:JSON.stringify(msg)})) //socket.send(json);
    
  }

export function createRoom(socket,roomId,username,standardFen){
    console.log('executing create',standardFen)
    sendJSONReq(socket,'createRoom',{roomId:roomId, username:username, fen:standardFen});
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
function isOpen(ws) { return ws.readyState === ws.OPEN }
/*export function sendRequest(socket,type,msg){
    let json = JSON.stringify(msg);
    if (socket.readyState !== 1){
        msgQueue.push(json)
    }else{
        socket.send(json);
    }
}*/
