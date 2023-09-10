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

<div id="wrapper">
	<div
		id={customBoardId}
		style={`--size: ${Math.max(boardConfig.dimensions.ranks, boardConfig.dimensions.files)}`}
		class={`${boardConfig.editable ? 'cursor-pointer' : null}`}
	>
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

<style>
	#board {
		display: grid;
		width: 100%;
		max-width: 700px;
		justify-items: center;
		grid-template-columns: repeat(var(--size), 1fr);
		grid-template-rows: repeat(var(--size), 1fr);
		border: 1px solid var(--default-dark-square);
		background-color: hsla(0, 10%, 94%, 0.7);
		user-select: none;
		box-shadow: 0 14px 28px rgba(0, 0, 0, 0.25), 0 10px 10px rgba(0, 0, 0, 0.22);
		touch-action: none;
		border-collapse: collapse;
		box-sizing: border-box;
	}
	#wrapper {
		position: relative;
		width: 100%;
		display: flex;
		justify-content: center;
	}
</style>