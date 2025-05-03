import { get, writable } from 'svelte/store';
import { Role, Status, chats, templateStore, gameId, gameState, members, type ConnectParams, type ConnectType, positionStore, fen, dimensions } from './store/stores';
import { camelToSnake } from './utils/index';
import type { EventType, Move, WSParams} from './types';
import { EventGameDrawOffer, EventGameResign, EventUserConnect, EventUserDisconnect, EventChatMessage, EventJoinGame, EventStartGame, EventGameMakeMove } from './types';
import { convertFenToPosition } from './board/fen';
import { toast } from './store/alert';

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


function createWebSocketStore() {
	const socket = writable<WebSocket | null>(null);

	function newWebSocketConnection(
		wsServerUrl: string,
		params: ConnectParams,
		connectType: ConnectType = 'join'
	) {
		const newWs = new WebSocket(wsServerUrl);

		newWs.onopen = () => {
			const wsMessage = {
				token: params.token,
				colorPref: params.colorPref
			};
			newWs.send(JSON.stringify(wsMessage));
			console.log('WebSocket connected');
		};

		newWs.onclose = () => {
			console.log('WebSocket closed');
			socket.set(null);
			gameId.set(null);
			templateStore.removeTemplate();
		};

		newWs.onerror = (e) => {
			toast.error('Websocket error',3000);
			console.error('WebSocket error:', e);
		};

		newWs.onmessage = (e) => {
			const data = JSON.parse(e.data);
			console.log('[WebSocket] Received:', data);
			handleMessage(data);
		};

		socket.set(newWs);
	}

	function sendMove(move: Move) {
		const ws = get(socket);
		if (ws && ws.readyState === WebSocket.OPEN) {
			const msg = { t: 'move', p: { m: move } };
			console.log('Sending move:', msg);
			ws.send(JSON.stringify(msg));
		} else {
			console.warn('WebSocket not open. Cannot send move.');
		}
	}

	function close() {
		const ws = get(socket);
		if (ws) {
			ws.close();
		}
	}

	return {
		subscribe: socket.subscribe,
		newWebSocketConnection,
		sendMove,
		close,
		set: socket.set
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
			gameState.setPlayers(eventData.players.b,eventData.players.w);
			gameState.updateStatus(Status.InProgress)
			if (eventData.gameConfig.fen) {
				fen.set(eventData.gameConfig.fen)
				const result = convertFenToPosition(eventData.gameConfig.fen);
				if (result?.position) positionStore.set(result.position)
				if (result?.dimensions) dimensions.set(result.dimensions)
			}
			//templateStore.setTemplate(data.gameConfig);
			break;
		case EventGameMakeMove:
			// Handle moves
			const move = eventData.m
			positionStore.makeMove(move);
			break;
		default:
			console.warn('[WebSocket] Unknown event:', data.event);
	}
}

const wsStore = createWebSocketStore();
export {
	createWebSocketStore,
	wsStore
}