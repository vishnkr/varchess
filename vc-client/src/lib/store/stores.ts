import { writable} from 'svelte/store';
import { Objective, type RuleEditorState, VariantType } from './types';

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
}

export interface User{
  username?:string,
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

const chats = newChatStore();

function createWebSocketStore(ws:WebSocket|null){
  const { subscribe, set, update } = writable<WebSocket|null>(ws);
  return{
    subscribe,
    set,
    update
  }
}

  
export const wsStore = createWebSocketStore(null);

const ruleEditor = writable<RuleEditorState>({
  variantType: VariantType.Custom,
  objective: Objective.Checkmate
})

export { 
  createWebSocketStore, 
  members, 
  roomId, 
  chats, 
  me, 
  ruleEditor
};
