

const serverUrl = import.meta.env.VITE_ENVIRONMENT === 'production' ? import.meta.env.VITE_SERVER_BASE : "localhost:5000";

export function connect(){
    console.log(serverUrl)
    const ws = new WebSocket(`ws://${serverUrl}/ws`);
    ws.addEventListener("message",()=>{
        /*const data = JSON.parse(edata)
        console.log("got",data);*/
    })
    return ws;
}

export const ws = connect();
