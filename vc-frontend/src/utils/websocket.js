

const WS = new WebSocket('ws://localhost:5000/ws');
WS.onopen = function(){
    console.log("socket opened!");
}
WS.onmessage = function(msg){
    let apiMsg = JSON.parse(msg.data);
    console.log('receirved message',apiMsg)
    msgQueue.push(apiMsg);
}
export default WS;


export let msgQueue = [];

export function sendJSONReq(socket,type,msg){
    if (!isOpen(socket)) return;
    socket.send(JSON.stringify({type:type,data:msg}))
  }

function isOpen(ws) { return ws.readyState === ws.OPEN }
export function sendRequest(socket,type,msg){
    let json = JSON.stringify(msg);
    if (socket.readyState !== 1){
        msgQueue.push(json)
    }else{
        socket.send(json);
    }
}
