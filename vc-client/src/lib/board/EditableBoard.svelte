<script lang="ts">
	/* Board wrapper component to handle edit state and abstract board edit logic from board component*/
	import type { Position, SquareInfo, SquareMaps, PiecePresentInfo, SquareIdx } from './types';
	import './board-styles.css';
	import type { BoardConfig } from './types';
	import { generateSquareMaps, updatePiecePositionsFromMaxBoard } from './board';
	import { convertFenToPosition, createEmptyMaxBoardState } from './fen';
	import Board from './Board.svelte';
	import { editorMaxBoard } from './board';
	import { onDestroy } from 'svelte';

	export let boardConfig: BoardConfig;
	boardConfig = { ...boardConfig, interactive: false, editable: true };
	let { squares } = generateSquareMaps(boardConfig.dimensions, boardConfig.isFlipped ?? false);
	const convertedPos = convertFenToPosition(boardConfig.fen);
	let position: Position = { piecePositions: {}, disabled: {} };
	let maxBoardState: PiecePresentInfo[][] = $editorMaxBoard;
	if (convertedPos) {
		({ position, maxBoardState } = convertedPos);
	}
	editorMaxBoard.set(maxBoardState);

	export let customBoardId = "board";
	function updateBoardState() {
		boardConfig.dimensions = boardConfig.dimensions;
		let squareMaps = generateSquareMaps(boardConfig.dimensions, boardConfig.isFlipped ?? false);
		squares = squareMaps.squares;
	}

	$: {
		boardConfig.dimensions;
		updateBoardState();
	}

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
