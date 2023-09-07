<script lang="ts">
	import { getSquareColor, type Position } from './types';
	import Square from './Square.svelte';
	import './board-styles.css';
	import type { BoardConfig } from './types';
	import { generateSquareMaps } from './board';
	import { convertFenToPosition } from './fen';

	export let boardConfig: BoardConfig;
	export let squares = generateSquareMaps(
		boardConfig.dimensions,
		boardConfig.isFlipped ?? false
	).squares;

	export let customBoardId = "board";
	export let position: Position = convertFenToPosition(boardConfig.fen)?.position ?? {
		piecePositions: {},
		disabled: {}
	};
</script>

<div class="relative w-full flex justify-center">
	<div
		id={customBoardId}
		style={`--size: ${Math.max(boardConfig.dimensions.ranks, boardConfig.dimensions.files)}`}
		class={`relative w-full max-w-[700px] mx-auto justify-center border border-solid border-[var(--default-dark-square)] bg-[hsla(0, 10%, 94%, 0.7)] shadow-[0px 14px 28px rgba(0, 0, 0, 0.25), 0px 10px 10px rgba(0, 0, 0, 0.22)] select-none ${boardConfig.editable ? 'cursor-pointer' : null}`}
	>
		<div class="grid grid-cols-[var(--size),1fr] grid-rows-[var(--size),1fr]">
			{#each Array(boardConfig.dimensions.ranks * boardConfig.dimensions.files) as _, idx}
				<Square
					boardId={customBoardId}
					editable={boardConfig.editable}
					interactive={boardConfig.interactive}
					squareData={squares[idx]}
					color={getSquareColor(squares[idx]?.row, squares[idx]?.column, boardConfig.isFlipped)}
					piece={position.piecePositions[idx] ?? null}
					disabled={position.disabled[idx] ?? false}
				/>
			{/each}
		</div>
	</div>
</div>