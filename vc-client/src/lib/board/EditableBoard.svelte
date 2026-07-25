<script lang="ts">
	/* Board wrapper component to handle edit state and abstract board edit logic from board component*/
	import type { Position, PiecePresentInfo } from './types';
	import './board-styles.css';
	import type { BoardConfig } from './types';
	import { generateSquareMaps, updatePiecePositionsFromMaxBoard } from './board';
	import { convertFenToPosition, createEmptyMaxBoardState } from './fen';
	import Board from './Board.svelte';
	import { editorMaxBoard } from './board';
	import { onDestroy } from 'svelte';
	import { get } from 'svelte/store';

	export let boardConfig: BoardConfig;

	let { squares } = generateSquareMaps(boardConfig.dimensions, boardConfig.isFlipped ?? false);
	let position: Position = { piecePositions: {}, walls: {} };
	let maxBoardState: PiecePresentInfo[][] = get(editorMaxBoard);

	/** True once we have a real 2D board (not the initial `[[]]` placeholder). */
	function isHydratedMaxBoard(state: PiecePresentInfo[][] | null | undefined): boolean {
		return !!state && state.length >= 2 && Array.isArray(state[0]) && state[0].length >= 2;
	}

	// Prefer the live editor store so remounting (e.g. after move-pattern mode)
	// does not wipe placements by re-applying the original FEN.
	if (isHydratedMaxBoard(maxBoardState)) {
		position = updatePiecePositionsFromMaxBoard(maxBoardState, boardConfig.dimensions);
	} else {
		const convertedPos = convertFenToPosition(boardConfig.fen);
		if (convertedPos) {
			({ position, maxBoardState } = convertedPos);
			editorMaxBoard.set(maxBoardState);
		}
	}

	export let customBoardId = 'board';
	function updateBoardState() {
		boardConfig.dimensions = boardConfig.dimensions;
		let squareMaps = generateSquareMaps(boardConfig.dimensions, boardConfig.isFlipped ?? false);
		squares = squareMaps.squares;
	}

	$: updateBoardState();

	$: {
		maxBoardState;
		position = updatePiecePositionsFromMaxBoard(maxBoardState, boardConfig.dimensions);
	}

	const unsubscribe = editorMaxBoard.subscribe((value) => (maxBoardState = value));
	onDestroy(unsubscribe);

	export const clear = (): void => {
		maxBoardState = createEmptyMaxBoardState();
		position = updatePiecePositionsFromMaxBoard(maxBoardState, boardConfig.dimensions);
		$editorMaxBoard = maxBoardState;
	};

	export function getBoardData() {
		return {
			position,
			maxBoardState,
			dimensions: boardConfig.dimensions
		};
	}

	export const shift = (direction: string): void => {
		let [lastCol, lastRow, afterLastCol, afterLastRow] = [
			boardConfig.dimensions.files - 1,
			boardConfig.dimensions.ranks - 1,
			boardConfig.dimensions.files,
			boardConfig.dimensions.ranks
		];
		let tempSquares: PiecePresentInfo[][];
		switch (direction) {
			case 'right':
				maxBoardState = maxBoardState.map((row, i) => {
					if (i < boardConfig.dimensions.ranks) {
						return [
							...row.slice(lastCol, afterLastCol),
							...row.slice(0, lastCol),
							...row.slice(afterLastCol)
						];
					}
					return row;
				});
				break;
			case 'left':
				maxBoardState = maxBoardState.map((row: PiecePresentInfo[], i) => {
					let firstSquare = row[0];
					if (i < boardConfig.dimensions.ranks) {
						return [...row.slice(1, afterLastCol), ...[firstSquare], ...row.slice(afterLastCol)];
					}
					return row;
				});
				break;
			case 'up':
				let firstRowSquares = maxBoardState[0];
				tempSquares = maxBoardState.slice(1, afterLastRow);
				maxBoardState = [
					...tempSquares,
					...[firstRowSquares],
					...maxBoardState.slice(afterLastRow)
				];
				break;
			case 'down':
				let lastRowSquares = maxBoardState[lastRow];
				tempSquares = maxBoardState.slice(0, lastRow);
				maxBoardState = [...[lastRowSquares], ...tempSquares, ...maxBoardState.slice(afterLastRow)];
				break;
		}
		position = updatePiecePositionsFromMaxBoard(maxBoardState, boardConfig.dimensions);
		$editorMaxBoard = maxBoardState;
	};
</script>

<Board customBoardId="board" {boardConfig} {position} {squares} />
