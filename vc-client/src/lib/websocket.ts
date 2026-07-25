import { get, writable } from 'svelte/store';
import {
	Status,
	gameId,
	gameState,
	positionStore,
	fen,
	dimensions,
	recentMove,
	templateStore,
	resetGameHistory,
	captureStartPosition,
	appendHistoryMove,
	gameResult,
	startPosition,
	lastGameConfig,
	wormholePlayState,
	setPlayerPresence,
	playerPresence,
	chats,
	legalMoves,
	clearPremove,
	checkedKingSquare,
	lastHintsFrom,
	viewPly,
	seekEnd,
	type ConnectParams
} from './store/stores';
import { camelToSnake, largeToSmallIndex, smallToLargeIndex } from './utils/index';
import {
	applyWormholeVariantState,
	wormholePlayStateFromConfig
} from './utils/wormhole';
import {
	EventGameDrawAccept,
	EventGameDrawOffer,
	EventGameDrawReject,
	EventGameMakeMove,
	EventGameOver,
	EventGameResign,
	EventPresence,
	EventChat,
	EventHints,
	EventStartGame,
	Color,
	type ActiveGamePayload,
	type ChatPayload,
	type EventType,
	type GameOverPayload,
	type HintsResultPayload,
	type Move,
	type MoveWirePayload,
	type PresencePayload,
	type Position,
	type WSMessage
} from './types';
import { convertFenToPosition } from './board/fen';
import { toast } from './store/alert';
import { authStore } from './store/auth';
import { settings } from './store/settings';
import { playSfx, type SfxKind } from './utils/sfx';

const RESUME_KEY = 'vc-ws-resume';
const MAX_BACKOFF_MS = 15_000;
const MAX_RECONNECT_ATTEMPTS = 8;

export type WsConnectionStatus = 'idle' | 'connecting' | 'open' | 'reconnecting' | 'closed';

export const wsStatus = writable<WsConnectionStatus>('idle');
/** True when the opponent has offered a draw (local player should accept/decline). */
export const incomingDrawOffer = writable(false);
/** True after we send a draw-offer until the echo arrives (avoid treating it as opponent's). */
let localDrawOfferPending = false;
let pendingLocalMoveSfx: string | null = null;

interface ResumeSnapshot {
	gameId: string;
	wsUrl: string;
	colorPref?: string;
}

function saveResume(snapshot: ResumeSnapshot) {
	try {
		sessionStorage.setItem(RESUME_KEY, JSON.stringify(snapshot));
	} catch {
		/* ignore */
	}
}

function loadResume(): ResumeSnapshot | null {
	try {
		const raw = sessionStorage.getItem(RESUME_KEY);
		return raw ? (JSON.parse(raw) as ResumeSnapshot) : null;
	} catch {
		return null;
	}
}

function clearResume() {
	try {
		sessionStorage.removeItem(RESUME_KEY);
	} catch {
		/* ignore */
	}
}

function boardDims() {
	const d = get(dimensions);
	return { files: d?.files ?? 8, ranks: d?.ranks ?? 8 };
}

function toWireMove(move: Move): Move {
	const { files, ranks } = boardDims();
	return {
		...move,
		from: smallToLargeIndex(move.from, files, ranks),
		to: smallToLargeIndex(move.to, files, ranks)
	};
}

function fromWireMove(move: Move): Move {
	const { files, ranks } = boardDims();
	return {
		...move,
		from: largeToSmallIndex(move.from, files, ranks),
		to: largeToSmallIndex(move.to, files, ranks)
	};
}

function toMoveJSON(move: Move) {
	const wire = toWireMove(move);
	return camelToSnake({
		from: wire.from,
		to: wire.to,
		piece: wire.piece,
		capture: wire.capture ?? false,
		promotion: wire.promotion,
		classicMoveType: wire.classicMoveType,
		variantMoveType: wire.variantMoveType,
		additionalData: wire.additionalData
	});
}

function sendTyped(ws: WebSocket, t: EventType, p: unknown = null) {
	const msg: WSMessage = { t, p };
	ws.send(JSON.stringify(msg));
}

function hydrateFromActiveGame(payload: ActiveGamePayload) {
	resetGameHistory();
	const players = payload.players ?? {};
	const white = players.w ?? players.white;
	const black = players.b ?? players.black;
	if (white && black) {
		gameState.setPlayers(
			{ name: black.name, userId: String(black.userId) },
			{ name: white.name, userId: String(white.userId) }
		);
	}

	if (payload.state === 'ip' || payload.state === 'InProgress') {
		gameState.updateStatus(Status.InProgress);
	} else if (payload.state === 'w') {
		gameState.updateStatus(Status.Waiting);
	}

	if (payload.turn === 'w' || payload.turn === 'b') {
		gameState.update((gs) => ({ ...gs, turn: payload.turn as 'w' | 'b' }));
	}

	const config = payload.gameConfig;
	if (config) {
		lastGameConfig.set(config as Record<string, unknown>);
	}
	if (config?.fen) {
		fen.set(config.fen);
		const result = convertFenToPosition(config.fen);
		if (result?.position) {
			positionStore.set(result.position);
			captureStartPosition(result.position);
		}
		if (result?.dimensions) {
			dimensions.set(result.dimensions);
		} else if (config.dimensions) {
			dimensions.set(config.dimensions);
		}
	} else if (config?.dimensions) {
		dimensions.set(config.dimensions);
	}

	const dims = get(dimensions);
	const files = dims?.files ?? config?.dimensions?.files ?? 8;
	const ranks = dims?.ranks ?? config?.dimensions?.ranks ?? 8;
	const customData = (config?.customData ?? {}) as Record<string, unknown>;
	const variantState = (payload.variantState ?? null) as Record<string, unknown> | null;
	wormholePlayState.set(wormholePlayStateFromConfig(customData, variantState, files, ranks));

	if (white?.userId) setPlayerPresence(String(white.userId), true);
	if (black?.userId) setPlayerPresence(String(black.userId), true);

	const recorded = payload.moves ?? [];
	let lastMove: Move | null = null;
	for (const moveStr of recorded) {
		try {
			const parsed = typeof moveStr === 'string' ? JSON.parse(moveStr) : moveStr;
			const wireMove: Move = {
				from: parsed.from,
				to: parsed.to,
				piece: parsed.piece,
				capture: parsed.capture,
				promotion: parsed.promotion,
				classicMoveType: parsed.classic_move_type ?? parsed.classicMoveType,
				variantMoveType: parsed.variant_move_type ?? parsed.variantMoveType,
				additionalData: parsed.additional_data ?? parsed.additionalData
			};
			const move = fromWireMove(wireMove);
			if (!positionStore.makeMove(move)) {
				console.warn('[WebSocket] Hydrate stopped: move did not apply', move);
				break;
			}
			appendHistoryMove(move);
			lastMove = move;
		} catch (err) {
			console.warn('[WebSocket] Skipping bad recorded move', err);
		}
	}
	if (lastMove) recentMove.set(lastMove);
}

function findKingSquare(pos: Position, color: Color): number | null {
	for (const key of Object.keys(pos.piecePositions)) {
		const idx = Number(key);
		const p = pos.piecePositions[idx];
		if (p && p.notation.toLowerCase() === 'k' && p.color === color) return idx;
	}
	return null;
}

function moveSfxKind(move: Move, isCheck: boolean): SfxKind {
	if (isCheck) return 'check';
	if (move.capture) return 'capture';
	return 'move';
}

function applyMovePayload(payload: MoveWirePayload | Move | null | undefined) {
	if (!payload) return;
	const raw = 'm' in (payload as MoveWirePayload) ? (payload as MoveWirePayload).m : payload;
	const isCheck = !!(payload as MoveWirePayload).check;
	const variantState =
		(payload as MoveWirePayload).variantState ??
		((payload as Record<string, unknown>).variant_state as
			| Record<string, unknown>
			| null
			| undefined);
	if (!raw || typeof (raw as Move).from !== 'number' || typeof (raw as Move).to !== 'number') return;
	const srcMove = raw as Move & Record<string, unknown>;
	const wireMove: Move = {
		from: srcMove.from,
		to: srcMove.to,
		piece: srcMove.piece,
		capture: srcMove.capture as boolean | undefined,
		promotion: srcMove.promotion,
		classicMoveType:
			srcMove.classicMoveType ?? (srcMove.classic_move_type as Move['classicMoveType']),
		variantMoveType:
			srcMove.variantMoveType ?? (srcMove.variant_move_type as string | undefined),
		additionalData:
			srcMove.additionalData ?? (srcMove.additional_data as Record<string, unknown> | undefined)
	};
	const move = fromWireMove(wireMove);

	const current = get(positionStore);
	if (!current) return;
	if (!get(startPosition)) {
		captureStartPosition(current);
	}

	// Live moves must apply at the tip — never on a reviewed intermediate ply.
	if (get(viewPly) !== null) {
		seekEnd();
	}

	if (!positionStore.makeMove(move)) {
		// Duplicate/out-of-order move or already-desynced board — do not
		// advance history/turn/highlights (that caused board vs move-list drift).
		console.warn('[WebSocket] Ignoring move that did not apply', move);
		return;
	}
	appendHistoryMove(move);
	recentMove.set(move);
	legalMoves.set([]);
	lastHintsFrom.set(null);
	clearPremove();
	gameState.changeTurn();

	const dims = get(dimensions);
	const files = dims?.files ?? 8;
	const ranks = dims?.ranks ?? 8;
	if (variantState) {
		wormholePlayState.update((cur) =>
			applyWormholeVariantState(cur, variantState as Record<string, unknown>, files, ranks)
		);
	}

	const after = get(positionStore);
	const turn = get(gameState).turn;
	if (isCheck && after) {
		const kingColor = turn === 'w' ? Color.WHITE : Color.BLACK;
		checkedKingSquare.set(findKingSquare(after, kingColor));
	} else {
		checkedKingSquare.set(null);
	}

	const moveKey = `${move.from}:${move.to}`;
	if (pendingLocalMoveSfx === moveKey) {
		pendingLocalMoveSfx = null;
		// Local move/capture already played; upgrade to check ping if needed.
		if (isCheck && get(settings).soundEnabled) playSfx('check');
	} else if (get(settings).soundEnabled) {
		playSfx(moveSfxKind(move, isCheck));
	}
}

function playerNameFromId(userId: string): string {
	const players = get(gameState).players;
	if (!players) return 'Player';
	if (String(players.playerWhite?.userId) === userId) return players.playerWhite.name || 'White';
	if (String(players.playerBlack?.userId) === userId) return players.playerBlack.name || 'Black';
	const auth = get(authStore);
	if (auth.userId && String(auth.userId) === userId) return auth.username || 'You';
	return 'Player';
}

function handleMessage(data: WSMessage) {
	if (!data || typeof data.t !== 'string') return;
	const payload = data.p;

	switch (data.t) {
		case EventStartGame:
			hydrateFromActiveGame(payload as ActiveGamePayload);
			break;
		case EventGameMakeMove:
			applyMovePayload(payload as MoveWirePayload);
			break;
		case EventGameDrawOffer:
			if (localDrawOfferPending) {
				localDrawOfferPending = false;
				break;
			}
			incomingDrawOffer.set(true);
			chats.pushSystem('Opponent offered a draw');
			if (get(settings).soundEnabled) playSfx('offer');
			break;
		case EventGameDrawAccept:
			incomingDrawOffer.set(false);
			gameResult.set({ winner: 'draw', reason: 'agreement' });
			chats.pushSystem('Draw accepted');
			gameState.updateStatus(Status.Completed);
			if (get(settings).soundEnabled) playSfx('over');
			break;
		case EventGameDrawReject:
			incomingDrawOffer.set(false);
			chats.pushSystem('Draw declined');
			break;
		case EventGameOver: {
			const result = payload as GameOverPayload & {
				Winner?: string;
				Reason?: string;
			};
			const winner = result?.winner ?? result?.Winner ?? 'unknown';
			const reason = result?.reason ?? result?.Reason ?? 'game over';
			gameResult.set({ winner, reason });
			const overMsg =
				winner === 'draw' || winner === 'unknown'
					? `Game over${reason ? ` — ${reason}` : ''}`
					: `${winner} wins${reason ? ` — ${reason}` : ''}`;
			chats.pushSystem(overMsg);
			gameState.updateStatus(Status.Completed);
			clearResume();
			clearPremove();
			if (get(settings).soundEnabled) playSfx('over');
			break;
		}
		case EventPresence: {
			const p = payload as PresencePayload;
			const uid = p?.userId ?? (p as { UserID?: string }).UserID;
			if (uid == null) break;
			const id = String(uid);
			const online = Boolean(p.online);
			const prev = get(playerPresence)[id];
			setPlayerPresence(id, online);
			// Skip no-op updates (e.g. presence snapshot repeating current state).
			if (prev === online) break;
			const name = playerNameFromId(id);
			if (online) chats.userJoin(name);
			else chats.userLeave(name);
			if (get(settings).soundEnabled) playSfx(online ? 'join' : 'leave');
			break;
		}
		case EventChat: {
			const chat = payload as ChatPayload & { Sender?: string; Message?: string };
			const sender = chat?.sender ?? chat?.Sender ?? 'Player';
			const text = chat?.msg ?? chat?.Message ?? '';
			if (text) chats.updateChat(sender, text);
			break;
		}
		case EventHints: {
			const hints = payload as HintsResultPayload;
			if (!hints || typeof hints.from !== 'number' || !Array.isArray(hints.to)) break;
			const { files, ranks } = boardDims();
			const from = largeToSmallIndex(hints.from, files, ranks);
			legalMoves.set(
				hints.to.map((to) => ({
					from,
					to: largeToSmallIndex(to, files, ranks),
					piece: ''
				}))
			);
			lastHintsFrom.set(from);
			break;
		}
		default:
			console.warn('[WebSocket] Unknown event:', data.t);
	}
}

function createWebSocketStore() {
	const socket = writable<WebSocket | null>(null);

	let intentionalClose = false;
	let reconnectTimer: ReturnType<typeof setTimeout> | null = null;
	let reconnectAttempt = 0;
	let lastUrl: string | null = null;
	let lastParams: ConnectParams | null = null;
	let pendingMoves: Move[] = [];

	function clearReconnectTimer() {
		if (reconnectTimer) {
			clearTimeout(reconnectTimer);
			reconnectTimer = null;
		}
	}

	function flushPendingMoves(ws: WebSocket) {
		while (pendingMoves.length > 0 && ws.readyState === WebSocket.OPEN) {
			const move = pendingMoves.shift();
			if (move) {
				sendTyped(ws, EventGameMakeMove, { m: toMoveJSON(move) });
			}
		}
	}

	function attachHandlers(ws: WebSocket, params: ConnectParams) {
		ws.onopen = () => {
			reconnectAttempt = 0;
			wsStatus.set('open');
			ws.send(
				JSON.stringify({
					token: params.token,
					colorPref: params.colorPref
				})
			);
			flushPendingMoves(ws);
			console.log('WebSocket connected');
		};

		ws.onclose = (event) => {
			socket.set(null);
			if (intentionalClose) {
				wsStatus.set('closed');
				return;
			}

			console.warn('[WebSocket] closed abnormally', event.code, event.reason);
			scheduleReconnect();
		};

		ws.onerror = (e) => {
			console.error('WebSocket error:', e);
			toast.error('Websocket error', 3000);
		};

		ws.onmessage = (e) => {
			try {
				const data = JSON.parse(e.data) as WSMessage;
				console.log('[WebSocket] Received:', data);
				handleMessage(data);
			} catch (err) {
				console.error('[WebSocket] Bad message:', err);
				toast.error('Invalid server message', 3000);
			}
		};
	}

	function scheduleReconnect() {
		if (intentionalClose || !lastUrl || !lastParams) {
			wsStatus.set('closed');
			return;
		}
		if (reconnectAttempt >= MAX_RECONNECT_ATTEMPTS) {
			wsStatus.set('closed');
			toast.error('Lost connection to game', 4000);
			clearResume();
			return;
		}

		wsStatus.set('reconnecting');
		const delay = Math.min(1000 * 2 ** reconnectAttempt, MAX_BACKOFF_MS);
		reconnectAttempt += 1;
		clearReconnectTimer();
		reconnectTimer = setTimeout(() => {
			const auth = get(authStore);
			if (!auth.accessToken) {
				wsStatus.set('closed');
				toast.error('Session expired', 3000);
				return;
			}
			lastParams = {
				...lastParams!,
				token: auth.accessToken
			};
			openConnection(lastUrl!, lastParams);
		}, delay);
	}

	function openConnection(wsServerUrl: string, params: ConnectParams) {
		intentionalClose = false;
		lastUrl = wsServerUrl;
		lastParams = params;
		wsStatus.set(reconnectAttempt > 0 ? 'reconnecting' : 'connecting');

		const id = get(gameId);
		if (id) {
			saveResume({
				gameId: id,
				wsUrl: wsServerUrl,
				colorPref: params.colorPref
			});
		}

		const newWs = new WebSocket(wsServerUrl);
		attachHandlers(newWs, params);
		socket.set(newWs);
	}

	function newWebSocketConnection(wsServerUrl: string, params: ConnectParams) {
		clearReconnectTimer();
		reconnectAttempt = 0;
		const existing = get(socket);
		if (existing && existing.readyState === WebSocket.OPEN) {
			intentionalClose = true;
			existing.close();
		}
		openConnection(wsServerUrl, params);
	}

	function sendMove(move: Move) {
		const ws = get(socket);
		if (ws && ws.readyState === WebSocket.OPEN) {
			sendTyped(ws, EventGameMakeMove, { m: toMoveJSON(move) });
		} else {
			pendingMoves.push(move);
			console.warn('[WebSocket] Queued move while disconnected');
		}
		// Play in the user-gesture path; skip the echoed WS move SFX for this move.
		// If the move gives check, the echo upgrades to the check ping.
		if (get(settings).soundEnabled) {
			pendingLocalMoveSfx = `${move.from}:${move.to}`;
			playSfx(move.capture ? 'capture' : 'move');
		}
	}

	function sendResign() {
		const ws = get(socket);
		if (ws && ws.readyState === WebSocket.OPEN) {
			sendTyped(ws, EventGameResign, null);
			chats.pushSystem('You resigned');
		}
	}

	function sendDrawOffer() {
		const ws = get(socket);
		if (ws && ws.readyState === WebSocket.OPEN) {
			localDrawOfferPending = true;
			sendTyped(ws, EventGameDrawOffer, null);
			chats.pushSystem('You offered a draw');
		}
	}

	function sendDrawAccept() {
		const ws = get(socket);
		if (ws && ws.readyState === WebSocket.OPEN) {
			sendTyped(ws, EventGameDrawAccept, null);
		}
	}

	function sendDrawReject() {
		const ws = get(socket);
		if (ws && ws.readyState === WebSocket.OPEN) {
			sendTyped(ws, EventGameDrawReject, null);
		}
	}

	function sendChat(message: string) {
		const ws = get(socket);
		const trimmed = message.trim();
		if (!trimmed || !ws || ws.readyState !== WebSocket.OPEN) return;
		sendTyped(ws, EventChat, { msg: trimmed });
	}

	/** Request legal destinations for a square (wire indices). */
	function requestHints(fromSmall: number) {
		const ws = get(socket);
		if (!ws || ws.readyState !== WebSocket.OPEN) return;
		const { files, ranks } = boardDims();
		sendTyped(ws, EventHints, { from: smallToLargeIndex(fromSmall, files, ranks) });
	}

	function close(options?: { clearGame?: boolean }) {
		intentionalClose = true;
		clearReconnectTimer();
		pendingMoves = [];
		clearResume();
		wsStatus.set('closed');
		const ws = get(socket);
		if (ws) {
			ws.close();
		}
		socket.set(null);
		if (options?.clearGame !== false) {
			gameId.set(null);
			templateStore.removeTemplate();
			chats.clear();
		}
	}

	/** Attempt resume after page reload using sessionStorage + fresh token. */
	function tryResume(): boolean {
		const snapshot = loadResume();
		const auth = get(authStore);
		if (!snapshot || !auth.accessToken) return false;
		gameId.set(snapshot.gameId);
		if (get(gameState).status === Status.None) {
			gameState.updateStatus(Status.Waiting);
		}
		newWebSocketConnection(snapshot.wsUrl, {
			token: auth.accessToken,
			userId: auth.userId ?? '',
			colorPref: snapshot.colorPref
		});
		return true;
	}

	return {
		subscribe: socket.subscribe,
		newWebSocketConnection,
		sendMove,
		sendResign,
		sendDrawOffer,
		sendDrawAccept,
		sendDrawReject,
		sendChat,
		requestHints,
		close,
		tryResume,
		set: socket.set
	};
}

const wsStore = createWebSocketStore();
export { createWebSocketStore, wsStore, loadResume, clearResume, RESUME_KEY };
