import { writable } from 'svelte/store';
import { Role, Status, chats, configStore, gameId, gameState, members, type ConnectParams, type ConnectType } from './store/stores';
import { EventChatMessage, EventGameDrawOffer, EventGameMakeMove, EventGameResign, EventJoinGame, EventStartGame, EventUserConnect, type EventType, type WSParams, type WebSocketMessage, EventUserDisconnect } from './store/types';
import { camelToSnake } from './utils/index';
import type { Move } from './board/types';

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

export const createWebSocket = (wsServerUrl: string) => {
	return new Promise((resolve, reject) => {
		const ws = new WebSocket(wsServerUrl);

		ws.onmessage = function (event) {
			const { type, data } = JSON.parse(event.data);
			switch (type) {
				case 'UserJoin':
					handleUserJoin(data);
					break;
				case 'UserLeave':
					chats.userLeave(data.username);
					members.update((value) => value.filter((member) => data.username !== member.username));
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
		ws.onclose = function () {
			wsStore.set(null);
			resolve(ws);
		};
	});
};

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
				const wsMessage = { event: `game.${connectType}_game`, params: params };
				const json = JSON.stringify(camelToSnake(wsMessage));
				ws.send(json);
				console.log('WebSocket connection success');
			};
			ws.onclose = () => {
				console.log('close called');
				set(null);
				gameId.set(null);
				configStore.removeConfig();
			};
			ws.onerror = (e) => {
				console.error('WebSocket connection error:', e);
				return;
			};
			ws.onmessage = (e) => {
				const data = JSON.parse(e.data);
				console.log('got message', data);
				if (data.success && data.result) {
					switch (data?.event) {
						case EventUserConnect:
							gameId.set(data.result.game_id);
							handleUserJoin(data);
							break;
						case EventUserDisconnect:
							chats.userLeave(data.username);
							members.update((value) => value.filter((member) => data.username !== member.username));
							break;
						case EventChatMessage:
							chats.updateChat(data.result.username, data.result.message);
							break;
						case EventJoinGame:
							//gameState.setGameConfig(data.result.game_config)
	
							/*gameId.set(data.result.game_id);
				gameState.setGameConfig(data.result.game_config)*/
							break;
						case EventStartGame:
							gameState.update((oldState) => {
								return {
									...oldState,
									players: {
										playerBlack: data.result.players.black,
										playerWhite: data.result.players.white
									},
									status: Status.InProgress
								};
							});
							configStore.setConfig(data.result.game_config);
							console.log('updating state', gameState);
							break;
						case EventGameMakeMove:
	
							break;
						default:
							console.log('invalid msg type');
					}
				}
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

const wsStore = createWebSocketStore(null);
export {
	createWebSocketStore,
	wsStore
}