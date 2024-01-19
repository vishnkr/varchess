import { writable} from 'svelte/store';
import { type RuleEditorState, VariantType, type MovePattern } from './types';
import { camelToSnake } from '$lib/utils';
import type { ChessCoreLib } from 'stonkfish';

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

export interface Config{
	variantType: VariantType;
    dimensions: {
        ranks: number;
        files: number;
    };
    fen: string;
    pieceProps: Record<string, MovePattern>;
	additionalData?: object;
}

const roomId = writable<string|null>(null);
const members = writable<Member[]>([]);
const me = writable<User|null>({});


function newConfigStore(){
	const {subscribe, set, update} = writable<Config|null>(null);
	const setConfig = (config:Config)=>{
		set(config)
	}
	
	const removeConfig = ()=>{
		set(null)
	}
	return{
		setConfig,
		removeConfig,
		update,
		subscribe
	}
}

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

const configStore = newConfigStore();
const chats = newChatStore();

export interface ConnectParams {
	sessionId: string;
	gameId?: string;
  }

export interface CreateParams extends ConnectParams{
	gameConfig?:Config,
	color:string
}
export type ConnectType = "create"|"join";

function createWebSocketStore(ws:WebSocket|null){
  const { subscribe, set, update } = writable<WebSocket|null>(ws);

  const newWebSocketConnection = async (url:string,params:ConnectParams,connectType:ConnectType = "join")=>{
	const ws = await new WebSocket(url);

	ws.onopen = () => {
		const wsMessage = { event: `game.${connectType}_game`, params: params };
		const json = JSON.stringify(camelToSnake(wsMessage));
		//console.log(json, ws);
		ws.send(json);
		//console.log('WebSocket connection success');
	  };
	ws.onclose = ()=>{
		set(null);
	}
	ws.onerror = (e)=>{
		console.error('WebSocket connection error:', e);
		 return
	}
	ws.onmessage = (e)=>{
		const data = JSON.parse(e.data);
		switch (data?.event){
			case "game.connect_user":
				if (data.result?.game_id){
					gameId.set(data.result.game_id);
				}
				break;
			default:
				console.log('invalid msg type');
		}
	}
	set(ws);
  }
  return{
	newWebSocketConnection,
    subscribe,
    set,
    update
  }
}

  
const wsStore = createWebSocketStore(null);
const gameId = writable<string|null>(null);

export interface GameState{
	status : "Waiting" | "InProgress",
	chesscore : ChessCoreLib,
}

const gameState = writable<GameState|null>(null);

const ruleEditor = writable<RuleEditorState>({
  variantType: VariantType.Checkmate,
})

export { 
  createWebSocketStore, 
  members, 
  roomId, 
  chats, 
  wsStore,
  me, 
  ruleEditor,
  configStore,
  gameId,
  gameState
};
