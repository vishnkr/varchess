<script lang="ts">
	import Board from '$lib/board/Board.svelte';
	import { BoardType, Color, type BoardConfig, type Position } from '$lib/board/types';
	import { convertFenToPosition } from '$lib/board/fen';
	import type { MovePattern } from '$lib/types';
	import { buildMpSquares, customPieceName } from '$lib/utils/customPieces';

	export let notation: string;
	export let pattern: MovePattern;
	export let boardId = 'pattern-preview';
	/** Larger board for floating full-pattern preview. */
	export let large = false;

	$: n = notation.toLowerCase();
	$: fenChar = n.toUpperCase();
	$: boardConfig = {
		fen: `9/9/9/9/4${fenChar}4/9/9/9/9 w - - 0 1`,
		dimensions: { ranks: 9, files: 9 },
		isFlipped: false,
		boardType: BoardType.MovePatternView
	} satisfies BoardConfig;

	$: fenResult = convertFenToPosition(boardConfig.fen);
	$: position = (fenResult?.position ?? {
		piecePositions: {
			40: { notation: n, pieceType: customPieceName(n).toLowerCase(), color: Color.WHITE }
		},
		walls: {}
	}) as Position;

	$: mpSquares = buildMpSquares(pattern);
</script>

<div class="pattern-preview" class:large>
	<Board customBoardId={boardId} {mpSquares} {boardConfig} {position} />
	<div
		class="flex items-center justify-center gap-5 mt-2 text-sm text-gray-600 dark:text-gray-400"
	>
		<span class="inline-flex items-center gap-2">
			<span class="w-3.5 h-3.5 rounded-sm bg-blue-600/80" />
			Slide
		</span>
		<span class="inline-flex items-center gap-2">
			<span class="w-3.5 h-3.5 rounded-sm bg-red-600/80" />
			Jump
		</span>
	</div>
</div>

<style>
	.pattern-preview {
		width: 100%;
	}
	.pattern-preview :global(#wrapper) {
		min-height: 0;
		padding: 0;
		width: 100%;
	}
	.pattern-preview :global(.vc-board) {
		width: 100%;
		max-width: 100%;
	}
	.pattern-preview.large :global(.vc-board) {
		width: 100%;
		max-width: 400px;
	}
</style>
