import { writable } from 'svelte/store';
import { Role, Status, chats, templateStore, gameId, gameState, members, type ConnectParams, type ConnectType } from './store/stores';
import { camelToSnake } from './utils/index';
import type { EventType, WSParams, Move} from './types';
import { EventGameDrawOffer, EventGameResign, EventUserConnect, EventUserDisconnect, EventChatMessage, EventJoinGame, EventStartGame, EventGameMakeMove } from './types';

interface UserJoin {
	id: number;
	username: string;
	isHost: boolean;
	role: string;
}

function handleUserJoin(data: UserJoin) {
	chats.userJoin(data.username);
	members.update((value) => [
		...value,
		{
			id: data.id,
			username: data.username,
			isHost: data.isHost,
			role: Role[data.role as keyof typeof Role]
		}
	]);
}


export function sendWebsocketMsg(ws:WebSocket,eventType: EventType, params: WSParams){
	ws.send(JSON.stringify(camelToSnake({eventType,params})));
}

export function sendDrawOffer(ws:WebSocket, gameId:string){
	ws.send(JSON.stringify(camelToSnake(
		{
			event: EventGameDrawOffer,
			params:{gameId}
		}
	)))
}

export function sendResign(ws:WebSocket, gameId:string){
	ws.send(JSON.stringify(camelToSnake(
		{
			event: EventGameResign,
			params:{gameId}
		}
	)))
}


function createWebSocketStore(ws: WebSocket | null) {
	const { subscribe, set, update } = writable<WebSocket | null>(ws);
	return {
		newWebSocketConnection: async (
			wsServerUrl: string,
			params: ConnectParams,
			connectType: ConnectType = 'join'
		) => {
			const ws = await new WebSocket(wsServerUrl);
	
			ws.onopen = () => {
				
				const wsMessage = {
					token: params.token,
					colorPref: params.colorPref 
				}
				//const wsMessage = { event: `game.${connectType}_game`, params: connectPayload };
				const json = JSON.stringify(wsMessage);
				ws.send(json);
				console.log('WebSocket connection success');
			};
			ws.onclose = () => {
				console.log('close called');
				set(null);
				gameId.set(null);
				templateStore.removeTemplate();
			};
			ws.onerror = (e) => {
				console.error('WebSocket connection error:', e);
				return;
			};
			ws.onmessage = (e) => {
				const data = JSON.parse(e.data);
				console.log('[WebSocket] Received:', data);
				handleMessage(data);
			};			
			set(ws);
		},
		subscribe,
		set,
		update,
		sendMove: (move:Move)=>{
			if (ws && ws.readyState === WebSocket.OPEN) {
				ws.send(JSON.stringify(move));
			}
		},
		close: () => {
			if (ws) {
				ws.close();
			}
		}
	};
}

function handleMessage(data: any) {
	console.log('got',data)
	if (!data?.d) return;
	const eventData = data.d;
	switch (data.t) {
		case EventUserConnect:
			gameId.set(data.result.game_id);
			handleUserJoin(data.result);
			break;
		case EventUserDisconnect:
			chats.userLeave(data.result.username);
			members.update((m) => m.filter((mem) => mem.username !== data.result.username));
			break;
		case EventChatMessage:
			chats.updateChat(data.result.username, data.result.message);
			break;
		case EventJoinGame:
			// gameState.setGameConfig(data.result.game_config);
			break;
		case EventStartGame:
			console.log('data',eventData);
			gameState.setPlayers(eventData.players.b,eventData.players.w);
			gameState.updateStatus(Status.InProgress)
			//templateStore.setTemplate(data.gameConfig);
			break;
		case EventGameMakeMove:
			// Handle moves
			break;
		default:
			console.warn('[WebSocket] Unknown event:', data.event);
	}
}

const wsStore = createWebSocketStore(null);
export {
	createWebSocketStore,
	wsStore
}