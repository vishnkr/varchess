import { convertFenToPosition } from '$lib/board/fen';
	import {
		appendHistoryMove,
		captureStartPosition,
		dimensions,
		fen,
		gameResult,
		gameState,
		lastGameConfig,
		positionStore,
		recentMove,
		resetGameHistory,
		Status,
		wormholePlayState
	} from '$lib/store/stores';
	import type { Game, Move } from '$lib/types';
	import { largeToSmallIndex } from '$lib/utils/index';
	import { wormholePlayStateFromConfig } from '$lib/utils/wormhole';
	import { get } from 'svelte/store';

function fromWireMove(move: Move): Move {
	const dims = get(dimensions) ?? { files: 8, ranks: 8 };
	const { files, ranks } = dims;
	return {
		...move,
		from: largeToSmallIndex(move.from, files, ranks),
		to: largeToSmallIndex(move.to, files, ranks)
	};
}

export function parseRecordedMove(moveStr: string | object): Move {
	const parsed = typeof moveStr === 'string' ? JSON.parse(moveStr) : moveStr;
	const raw = parsed as Record<string, unknown>;
	return {
		from: raw.from as number,
		to: raw.to as number,
		piece: (raw.piece as string) ?? '',
		capture: raw.capture as boolean | undefined,
		promotion: raw.promotion as string | undefined,
		classicMoveType: (raw.classic_move_type ?? raw.classicMoveType) as Move['classicMoveType'],
		variantMoveType: (raw.variant_move_type ?? raw.variantMoveType) as string | undefined,
		additionalData: (raw.additional_data ?? raw.additionalData) as Record<string, unknown> | undefined
	};
}

/** Load a completed persisted game into play stores for history review. */
export function hydrateCompletedGame(game: Game) {
	resetGameHistory();

	const names = game.playerNames ?? {};
	const whiteName = names.w || names.white || 'White';
	const blackName = names.b || names.black || 'Black';
	const players = game.players ?? {};
	gameState.setPlayers(
		{ name: blackName, userId: String(players.b ?? players.black ?? '') },
		{ name: whiteName, userId: String(players.w ?? players.white ?? '') }
	);

	const config = game.config;
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
		const empty = { piecePositions: {}, walls: {} };
		positionStore.set(empty);
		captureStartPosition(empty);
	}

	const dims = get(dimensions);
	const files = dims?.files ?? config?.dimensions?.files ?? 8;
	const ranks = dims?.ranks ?? config?.dimensions?.ranks ?? 8;
	wormholePlayState.set(
		wormholePlayStateFromConfig(
			(config?.customData ?? {}) as Record<string, unknown>,
			null,
			files,
			ranks
		)
	);

	let lastMove: Move | null = null;
	for (const moveStr of game.moves ?? []) {
		try {
			const move = fromWireMove(parseRecordedMove(moveStr));
			if (!positionStore.makeMove(move)) {
				console.warn('Hydrate stopped: move did not apply', move);
				break;
			}
			appendHistoryMove(move);
			lastMove = move;
		} catch (err) {
			console.warn('Skipping bad recorded move', err);
		}
	}
	if (lastMove) recentMove.set(lastMove);

	if (game.result) {
		gameResult.set({
			winner: game.result.winner ?? '',
			reason: game.result.reason ?? ''
		});
	}

	gameState.updateStatus(Status.Completed);
}
