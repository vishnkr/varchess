


let msgQueue = [];
export function sendJSONReq(socket,type,msg){
    socket.send(JSON.stringify({type:type,data:msg}))
  }

export function sendRequest(socket,type,msg){
    let json = JSON.stringify(msg);
    if (socket.readyState !== 1){
        msgQueue.push(json)
    }else{
        socket.send(json);
    }
}