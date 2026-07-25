<script lang="ts">
	import { BoardType, getSquareColor, isMovePatternEditor, type Position } from './types';
	import Square from './Square.svelte';
	import slide from '$lib/assets/svg/slide.svg';
	import jump from '$lib/assets/svg/jump.svg';
	import './board-styles.css';
	import type { BoardConfig } from './types';
	import { generateSquareMaps } from './board';
	import { MoveType } from '$lib/types';
	import { recentMove } from '$lib/store/stores';
	import { tick } from 'svelte';
	import { browser } from '$app/environment';

	export let boardConfig: BoardConfig;
	export let isFlipped: boolean = false;
	export let squares = generateSquareMaps(boardConfig.dimensions, isFlipped).squares;

	export let customBoardId = 'board';
	export let position: Position;
	export let markedTargets: number[] = [];
	export let mpSquares: Record<number, MoveType> | null = null;
	/** When true, board fills a height-constrained parent (play page). */
	export let fitContainer = false;
	export let premoveFrom: number | null = null;
	export let premoveTo: number | null = null;
	export let showCoordinates = false;
	/** Animate last move on game boards (not static previews). */
	export let animateMoves = false;
	/** Square index of king currently in check. */
	export let checkedSquare: number | null = null;
	/** When set, only this side's pieces are selectable/draggable on a game board. */
	export let selectableColor: import('$lib/types').Color | null = null;

	let lastMapKey = '';
	let wrapperEl: HTMLDivElement;
	let animClass = '';
	let animStyle = '';
	let animVisible = false;
	let animHideTo: number | null = null;
	let lastAnimKey = '';

	function isCursorBoardType(bType: BoardType): boolean {
		return bType !== BoardType.MovePatternView && bType != BoardType.View;
	}

	$: {
		const key = `${boardConfig.dimensions.files}x${boardConfig.dimensions.ranks}:${isFlipped}`;
		if (key !== lastMapKey) {
			lastMapKey = key;
			squares = generateSquareMaps(boardConfig.dimensions, isFlipped).squares;
		}
	}

	const isMPSquareOcc = (idx: number) =>
		isMovePatternEditor(boardConfig.boardType) && mpSquares && mpSquares[idx] !== undefined;

	function pieceClassFromMovePiece(pieceStr: string): string {
		const isWhite = pieceStr === pieceStr.toUpperCase();
		return (isWhite ? 'w' : 'b') + pieceStr.toLowerCase();
	}

	async function playMoveAnimation(from: number, to: number, pieceStr: string) {
		if (!browser || !wrapperEl || !animateMoves) return;
		const key = `${from}-${to}-${pieceStr}`;
		if (key === lastAnimKey) return;
		lastAnimKey = key;

		await tick();
		const fromEl = document.getElementById(`${customBoardId}-s-${from}`);
		const toEl = document.getElementById(`${customBoardId}-s-${to}`);
		if (!fromEl || !toEl) return;

		const wrap = wrapperEl.getBoundingClientRect();
		const a = fromEl.getBoundingClientRect();
		const b = toEl.getBoundingClientRect();

		animClass = pieceClassFromMovePiece(pieceStr);
		animHideTo = to;
		animStyle = `width:${a.width}px;height:${a.height}px;transform:translate(${a.left - wrap.left}px, ${a.top - wrap.top}px);`;
		animVisible = true;

		await tick();
		requestAnimationFrame(() => {
			animStyle = `width:${a.width}px;height:${a.height}px;transform:translate(${b.left - wrap.left}px, ${b.top - wrap.top}px);transition:transform 140ms ease-out;`;
		});

		window.setTimeout(() => {
			animVisible = false;
			animHideTo = null;
			animStyle = '';
		}, 160);
	}

	$: if (animateMoves && $recentMove) {
		void playMoveAnimation($recentMove.from, $recentMove.to, $recentMove.piece || 'p');
	}
</script>

<div id="wrapper" class:fit-container={fitContainer} bind:this={wrapperEl}>
	<div
		id={customBoardId}
		class={`vc-board ${isCursorBoardType(boardConfig.boardType) ? 'cursor-pointer' : ''}`}
		style={`--cols: ${boardConfig.dimensions.files}; --rows: ${boardConfig.dimensions.ranks}`}
	>
		{#each Array(boardConfig.dimensions.ranks * boardConfig.dimensions.files) as _, idx}
			<Square
				boardId={customBoardId}
				squareData={squares[idx]}
				color={getSquareColor(squares[idx]?.row, squares[idx]?.column, boardConfig.isFlipped)}
				piece={position.piecePositions[idx] ?? null}
				wall={position.walls[idx] ?? false}
				boardType={boardConfig.boardType}
				isMarkedTarget={markedTargets.includes(idx)}
				isPremoveSrc={premoveFrom === idx}
				isPremoveDest={premoveTo === idx}
				isInCheck={checkedSquare === idx}
				{selectableColor}
				mpHighlight={mpSquares?.[idx] === MoveType.Slide
					? 'slide'
					: mpSquares?.[idx] === MoveType.Jump
						? 'jump'
						: null}
				{showCoordinates}
				boardRanks={boardConfig.dimensions.ranks}
				animHidePiece={animHideTo === idx}
				nonPieceSvg={isMPSquareOcc(idx) && mpSquares
					? mpSquares[idx] === MoveType.Slide
						? slide
						: jump
					: null}
			/>
		{/each}
	</div>

	{#if animVisible}
		<div class="move-anim bg-piece {animClass}" style={animStyle} aria-hidden="true" />
	{/if}
</div>

<style>
	#wrapper {
		position: relative;
		width: 100%;
		display: flex;
		justify-content: center;
		align-items: center;
	}
	.vc-board {
		display: grid;
		width: min(100%, 80%);
		aspect-ratio: var(--cols) / var(--rows);
		justify-items: stretch;
		align-items: stretch;
		grid-template-columns: repeat(var(--cols), 1fr);
		grid-template-rows: repeat(var(--rows), 1fr);
		user-select: none;
		touch-action: none;
		box-sizing: border-box;
	}
	#wrapper.fit-container {
		width: 100%;
		height: 100%;
		min-height: 0;
		container-type: size;
	}
	#wrapper.fit-container .vc-board {
		width: min(100cqw, calc(100cqh * var(--cols) / var(--rows)));
		height: auto;
		max-width: 100%;
		max-height: 100%;
		aspect-ratio: var(--cols) / var(--rows);
	}

	.move-anim {
		position: absolute;
		left: 0;
		top: 0;
		pointer-events: none;
		z-index: 20;
		will-change: transform;
	}
</style>
