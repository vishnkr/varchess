import { writable, type Readable, type Writable } from "svelte/store";

const serverUrl = import.meta.env.VITE_ENVIRONMENT === 'production' ? import.meta.env.VITE_SERVER_BASE : 'localhost:5000';

enum Role{
	Viewer,
	Black,
	White
}

export enum MessageType{
	UserJoin,
	UserLeave,
	ChatMessage,
	ResultUpdate
}

export interface ChatMessage{
	messageType: MessageType,
	content? : string,
	username?: string
}

interface Member{
	username: string,
	role: Role,
	isHost: boolean
}

interface RoomState{
	roomId: string,
	members: Member[],
}

const roomStore = writable<RoomState>();



function newChatStore(){
	const { subscribe, set, update } = writable<ChatMessage[]>([]);
	const userJoin = (username:string)=>{ 
		update((chats)=>[...chats,{messageType:MessageType.UserJoin,username:username,content:`${username} has joined the Room`}])
	}
	const userLeave = (username:string)=>{ 
		update((chats)=>[...chats,{messageType:MessageType.UserLeave,username:username,content:`${username} has left the Room`}])
	}
	const updateChat = (username:string,content:string)=>{
		update((chats)=>[...chats,{messageType:MessageType.ChatMessage,username:username,content:content}])
	}
	return {
		subscribe,
		set,
		update,
		userJoin,
		userLeave,
		updateChat
	}
}

const wsStore : Writable<WebSocket|null> = writable(null);
const chatStore = newChatStore();

function createWebSocket(roomId:string, username:string) {
	  const ws = new WebSocket(`ws://${serverUrl}/ws/${roomId}/${username}`);
	  ws.onmessage = function(event) {
		const {type,data} = JSON.parse(event.data);
		switch (type) {
		  case 'UserJoin':
			chatStore.userJoin(data.username)
			break;
		  case 'UserLeave':
			chatStore.userLeave(data.username)
			break;
		  case 'ChatMessage':
			// Handle incoming chat messages here
			chatStore.updateChat(data.username,data.content)
			break;
		}
	  };
	  
  	wsStore.set(ws);
	return ws;
  }


export { createWebSocket, roomStore, chatStore};
/*
export const socketStore = readable({},set=>{
	
})
function createWebSocketStore() {
	const { subscribe, set, update } = readable<WebSocket | null>(null);
	let ws: WebSocket | null = null;
		  
	function connect(roomId:string,username:string){
		if (!ws || ws.readyState === WebSocket.CLOSED) {
			ws = new WebSocket(`ws://${serverUrl}/ws/${roomId}/${username}`);
			ws.addEventListener('message', (event) => {

				update((currentWs) => {

				if (currentWs) {
					console.log('got message3',event)
					const data = JSON.parse(event.data);
					currentWs.dispatchEvent(new CustomEvent('message', { detail: data }));
					console.log('got',data)
				}
				return currentWs;
				});
			});
		}
		set(ws);
		console.log('connected ws')
	}
	return {
		ws,
		subscribe,
		connect,
	};
}
		  
export const webSocketStore = createWebSocketStore();
		  */