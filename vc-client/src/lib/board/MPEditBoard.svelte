<script lang="ts">
	import { pieceEditor } from '$lib/store/editor';
	import type { Position } from '$lib/types';
	import Board from './Board.svelte';
	import { BoardType, type BoardConfig, Color } from './types';
	import { convertFenToPosition } from './fen';
	import { buildMpSquares } from '$lib/utils/customPieces';

	$: selection = $pieceEditor.pieceSelection;
	$: notation = selection?.piece?.notation?.toLowerCase() ?? 'd';
	$: fenChar =
		selection?.piece?.color === Color.BLACK ? notation : notation.toUpperCase();

	$: boardConfig = {
		fen: `9/9/9/9/4${fenChar}4/9/9/9/9 w - - 0 1`,
		dimensions: { ranks: 9, files: 9 },
		isFlipped: false,
		boardType: BoardType.MovePatternEditor
	} satisfies BoardConfig;

	$: fenResult = convertFenToPosition(boardConfig.fen);
	$: position = (fenResult?.position ?? {
		piecePositions: {},
		walls: {}
	}) as Position;

	$: pattern = $pieceEditor.movePatterns[notation];
	$: mpSquares = buildMpSquares(pattern);
</script>

<Board
	customBoardId="mp-board"
	{mpSquares}
	{boardConfig}
	{position}
/>

<style>
	:global(#wrapper:has(#mp-board)) {
		min-height: min(70vh, 560px);
		padding: 0.5rem;
	}
</style>
