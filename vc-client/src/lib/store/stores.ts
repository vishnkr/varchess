import { writable, type Writable } from 'svelte/store';
import { Objective, type BoardEditorState, type RuleEditorState, VariantType } from './types';
import type { RecordAuthResponse, RecordModel } from 'pocketbase';

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
	role?: Role,
	isHost?: boolean
  userData?: RecordAuthResponse<RecordModel>
}

export interface User{
  username?:string,
  userData?: RecordAuthResponse<RecordModel>
}

const roomId = writable<string|null>(null);
const members = writable<Member[]>([]);
const me = writable<User|null>({});


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
        const { type, data } = JSON.parse(event.data);
        switch (type) {
          case 'UserJoin':
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
      ws.onclose = function(){
        wsStore.set(null);
        resolve(ws);
      }
    });
  }
  

const ruleEditor = writable<RuleEditorState>({
  variantType: VariantType.Standard,
  objective: Objective.Checkmate
})

export { 
  createWebSocket, 
  members, 
  roomId, 
  chats, 
  me, 
  ruleEditor
};
