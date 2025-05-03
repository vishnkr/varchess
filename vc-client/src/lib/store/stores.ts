import { writable, type Writable } from 'svelte/store';

import { type IPiece,type Move, type Position, type Template } from '$lib/types';
import { sequence } from '@sveltejs/kit/hooks';

export enum Role {
	Viewer,
	Black,
	White
}

export enum MessageType {
	UserJoin,
	UserLeave,
	ChatMessage,
	ResultUpdate
}

export interface ChatMessage {
	messageType: MessageType;
	content?: string;
	username?: string;
}

export interface Member {
	id: number;
	username: string;
	role?: Role;
	isHost?: boolean;
}

export interface User {
	username?: string;
}


const roomId = writable<string | null>(null);
const members = writable<Member[]>([]);

function newTemplateStore() {
	const { subscribe, set, update } = writable<Template | null>(null);
	const setTemplate = (template: Template) => {
		set(template);
	};

	const removeTemplate = () => {
		set(null);
	};
	return {
		setTemplate,
		removeTemplate,
		update,
		subscribe
	};
}

function newChatStore() {
	const { subscribe, set, update } = writable<ChatMessage[]>([]);
	const userJoin = (username: string) => {
		update((chats) => [
			...chats,
			{
				messageType: MessageType.UserJoin,
				username: username,
				content: `${username} has joined the Room`
			}
		]);
	};
	const userLeave = (username: string) => {
		update((chats) => [
			...chats,
			{
				messageType: MessageType.UserLeave,
				username: username,
				content: `${username} has left the Room`
			}
		]);
	};
	const updateChat = (username: string, content: string) => {
		update((chats) => [
			...chats,
			{ messageType: MessageType.ChatMessage, username: username, content: content }
		]);
	};

	return {
		subscribe,
		set,
		update,
		userJoin,
		userLeave,
		updateChat
	};
}

const templateStore = newTemplateStore();
const chats = newChatStore();

export interface ConnectParams {
	token: string;
	userId: string;
	colorPref?: string;
}

export interface CreateParams extends ConnectParams {
	gameConfig?: Template;
	color: string;
}
export type ConnectType = 'create' | 'join';

export interface MoveSelector {
	src: Writable<number | null>;
	dest: Writable<number | null>;
	piece: Writable<IPiece | null>;
	recentMove: Writable<Move | null>;
	legalMoves: Writable<Move[]>;
}


export const src:Writable<number|null> = writable(null);
export const dest:Writable<number|null> = writable(null);
export const moveSelectorPiece: Writable<IPiece | null> = writable(null);
export const  piece: Writable<IPiece | null> = writable(null);
export const legalMoves: Writable<Move[]> = writable([]);
export const recentMove: Writable<Move | null> = writable(null);
export const clearMoveSelectorStores = () =>{
	src.set(null);
	dest.set(null);
	moveSelectorPiece.set(null);
	legalMoves.set([]);
	recentMove.set(null);
}
const gameId = writable<string | null>(null);

export interface Players {
	playerWhite: Player;
	playerBlack: Player;
}

export interface Player{
	name:string;
	userId: string;
}

export enum Status {
	None,
	Waiting,
	InProgress,
	Completed
}

export interface GameState {
	status: Status;
	//chesscore : ChessCoreLib,
	turn: 'w' | 'b';
	players?: Players;
	
}

function newGameState() {
	const { subscribe, update } = writable<GameState>({
		status: Status.None,
		turn: 'w'
	});

	const updateStatus = (newStatus: Status) => {
		update((gameState) => {
			return {
				...gameState,
				status: newStatus
			};
		});
	};
	const changeTurn = () => {
		update((gameState) => {
			return {
				...gameState,
				turn: gameState.turn === 'w' ? 'b' : 'w'
			};
		});
	};

	const setPlayers = (playerBlack: Player, playerWhite: Player) => {
		update((gameState) => {
			return {
				...gameState,
				players: { playerBlack, playerWhite }
			};
		});
	};

	return {
		updateStatus,
		setPlayers,
		changeTurn,
		update,
		subscribe
	};
}

const gameState = newGameState();

export function createPositionStore(initialPosition:Position|null) {
	const { subscribe, update, set } = writable<Position|null>(null);
  
	return {
	  subscribe,
	  set,
	  update,
	  reset: () => set(initialPosition),
	  makeMove: (move: Move) => {
		update((pos) => {
			if (!pos) return pos;

			const { from, to } = move;
			const piece = pos.piecePositions[from];

			if (!piece) {
				console.warn(`No piece at source square ${src}`);
				return pos;
			}

			const newPiecePositions = { ...pos.piecePositions };

			delete newPiecePositions[from];
			newPiecePositions[to] = piece;

			return {
				...pos,
				piecePositions: newPiecePositions
			};
		});
	}
	};
  }
  export const positionStore = createPositionStore(null);

  export const fen = writable<string|null>(null);
  export const dimensions = writable<{ ranks: number; files: number }|null>(null);

/*
function createChessCoreWrapper() {
	const { subscribe, set } = writable<ChessCoreLib | null>(null);

	let coreInstance: ChessCoreLib | null = null;

	return {
		subscribe,
    initWasm: async()=>{
      await init();
    },
		loadPosition: (config_json: string) => {
			coreInstance = new ChessCoreLib(config_json);
			set(coreInstance);
		},
		getLegalMoves: () => {
			return coreInstance ? coreInstance.getLegalMoves() : [];
		},
		makeMove: () => {
			if (coreInstance) {
				//coreInstance.makeMove(move);
			}
		}
	};
}
*/

//const chessCore = createChessCoreWrapper();

export {
	members,
	roomId,
	chats,
	templateStore,
	gameId,
	gameState,
  ////chessCore
};
