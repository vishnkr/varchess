import { get, writable, type Writable } from 'svelte/store';

import { type IPiece, type Move, type Position, type Template } from '$lib/types';

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
				content: `${username} joined`
			}
		]);
	};
	const userLeave = (username: string) => {
		update((chats) => [
			...chats,
			{
				messageType: MessageType.UserLeave,
				username: username,
				content: `${username} left`
			}
		]);
	};
	const updateChat = (username: string, content: string) => {
		update((chats) => [
			...chats,
			{ messageType: MessageType.ChatMessage, username: username, content: content }
		]);
	};

	const pushSystem = (content: string) => {
		update((chats) => [
			...chats,
			{ messageType: MessageType.ResultUpdate, content }
		]);
	};

	return {
		subscribe,
		set,
		update,
		userJoin,
		userLeave,
		updateChat,
		pushSystem,
		clear: () => set([])
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


export const src: Writable<number | null> = writable(null);
export const dest: Writable<number | null> = writable(null);
export const moveSelectorPiece: Writable<IPiece | null> = writable(null);
export const piece: Writable<IPiece | null> = writable(null);
export const legalMoves: Writable<Move[]> = writable([]);
/** Square index of the side-to-move's king when in check; null otherwise. */
export const checkedKingSquare: Writable<number | null> = writable(null);
/** Set when server hints for a from-square have arrived (for illegal-move checks). */
export const lastHintsFrom: Writable<number | null> = writable(null);
export const recentMove: Writable<Move | null> = writable(null);
/** Queued premove while waiting for our turn (at most one). */
export const premove: Writable<Move | null> = writable(null);

export const clearMoveSelectorStores = () => {
	src.set(null);
	dest.set(null);
	moveSelectorPiece.set(null);
	if (get(legalMoves).length) legalMoves.set([]);
	lastHintsFrom.set(null);
};

export const clearPremove = () => {
	premove.set(null);
};

/** Full game move list (for sidebar + post-game replay). */
export const moveHistory = writable<Move[]>([]);
/** Starting position snapshot (before any moves). */
export const startPosition = writable<Position | null>(null);
/**
 * null = live tip of the game.
 * number = show board after that many plies (0 = start position).
 */
export const viewPly = writable<number | null>(null);
export const gameResult = writable<{ winner: string; reason: string } | null>(null);
/** Last ActiveGame gameConfig — used for "replay with same setup". */
export const lastGameConfig = writable<Record<string, unknown> | null>(null);
/** Live wormhole portals for the current game (FE dense indices). */
export const wormholePlayState = writable<import('$lib/utils/wormhole').WormholePlayState | null>(
	null
);
/** userId → online presence for the current game. */
export const playerPresence = writable<Record<string, boolean>>({});

export function setPlayerPresence(userId: string, online: boolean) {
	playerPresence.update((m) => ({ ...m, [userId]: online }));
}

export function resetPlayerPresence() {
	playerPresence.set({});
}

function clonePosition(pos: Position): Position {
	return {
		piecePositions: { ...pos.piecePositions },
		walls: { ...pos.walls }
	};
}

function applyMovePure(pos: Position, move: Move): Position {
	const next = clonePosition(pos);
	const moving = next.piecePositions[move.from];
	if (!moving) return next;
	delete next.piecePositions[move.from];
	next.piecePositions[move.to] = moving;
	return next;
}

export function resetGameHistory() {
	moveHistory.set([]);
	startPosition.set(null);
	viewPly.set(null);
	gameResult.set(null);
	lastGameConfig.set(null);
	wormholePlayState.set(null);
	resetPlayerPresence();
	clearPremove();
	clearMoveSelectorStores();
	checkedKingSquare.set(null);
	lastHintsFrom.set(null);
}

export function captureStartPosition(pos: Position) {
	startPosition.set(clonePosition(pos));
	moveHistory.set([]);
	viewPly.set(null);
}

export function appendHistoryMove(move: Move) {
	moveHistory.update((h) => [...h, move]);
	if (get(viewPly) === null) {
		recentMove.set(move);
	}
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
	makeMove: (move: Move): boolean => {
		let applied = false;
		update((pos) => {
			if (!pos) return pos;

			const { from, to } = move;
			const piece = pos.piecePositions[from];

			if (!piece) {
				console.warn(`No piece at source square ${from}`);
				return pos;
			}

			const newPiecePositions = { ...pos.piecePositions };
			delete newPiecePositions[from];
			newPiecePositions[to] = piece;
			applied = true;

			return {
				...pos,
				piecePositions: newPiecePositions
			};
		});
		return applied;
	}
	};
  }
  export const positionStore = createPositionStore(null);

  export const fen = writable<string|null>(null);
  export const dimensions = writable<{ ranks: number; files: number }|null>(null);

/** Seek the displayed board to after `ply` moves (0 = start). null = live tip. */
export function seekBoardToPly(ply: number | null) {
	const hist = get(moveHistory);
	const start = get(startPosition);
	if (!start) return;

	const target = ply === null ? hist.length : Math.max(0, Math.min(ply, hist.length));
	viewPly.set(target === hist.length ? null : target);

	src.set(null);
	dest.set(null);
	moveSelectorPiece.set(null);
	legalMoves.set([]);

	let pos = clonePosition(start);
	for (let i = 0; i < target; i++) {
		pos = applyMovePure(pos, hist[i]);
	}
	positionStore.set(pos);
	recentMove.set(target > 0 ? hist[target - 1] : null);
	checkedKingSquare.set(null);
}

export function seekPrev() {
	const hist = get(moveHistory);
	const current = get(viewPly);
	const tip = hist.length;
	const at = current === null ? tip : current;
	seekBoardToPly(Math.max(0, at - 1));
}

export function seekNext() {
	const hist = get(moveHistory);
	const current = get(viewPly);
	const tip = hist.length;
	const at = current === null ? tip : current;
	seekBoardToPly(Math.min(tip, at + 1));
}

export function seekStart() {
	seekBoardToPly(0);
}

export function seekEnd() {
	seekBoardToPly(null);
}

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
