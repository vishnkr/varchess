import { writable, type Writable } from 'svelte/store';

const serverUrl = import.meta.env.VITE_ENVIRONMENT === 'production' ? import.meta.env.VITE_SERVER_BASE : 'localhost:5000';

export enum Role{
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

export interface Member{
	id:number,
	username: string,
	role: Role,
	isHost: boolean
}

function createUserStore(){
    const {subscribe, set, update} = writable<Member>({
        username: '',
        id:0,
        role:Role.Viewer,
        isHost: false,
    })
    return {
        subscribe,
        set,
        update
    }
}

const roomId = writable<string|null>(null);
const members = writable<Member[]>([]);
const me = createUserStore();


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

export const wsStore : Writable<WebSocket|null> = writable(null);
const chats = newChatStore();

function createWebSocket(roomId:string, username:string) {
    return new Promise((resolve, reject) => {
      const ws = new WebSocket(`ws://${serverUrl}/ws/${roomId}/${username}`);
  
      ws.onmessage = function (event) {
        // Handle WebSocket messages here
        const { type, data } = JSON.parse(event.data);
        //console.log('got ws data', event.data, type, data, event);
        switch (type) {
          case 'UserJoin':
            //console.log('join', data);
            chats.userJoin(data.username);
            members.update((value) => [
              ...value,
              {
                id: data.id,
                username: data.username,
                isHost: data.isHost,
                role: Role[data.role as keyof typeof Role],
              },
            ]);
            break;
          case 'UserLeave':
            chats.userLeave(data.username);
            members.update((value) =>
              value.filter((member) => data.username !== member.username)
            );
            break;
          case 'ChatMessage':
            chats.updateChat(data.username, data.content);
            break;
        }
      };
  
      ws.onerror = function (error) {
        console.error('WebSocket connection error:', error);
        reject(error);
      };
  
      ws.onopen = function () {
        wsStore.set(ws);
        resolve(ws);
      };
    });
  }
  

export { createWebSocket, members, roomId, chats, me};
