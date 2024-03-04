<script lang="ts">
	import { BoardType, getSquareColor, type Position } from './types';
	import Square from './Square.svelte';
	import slide from '$lib/assets/svg/slide.svg';
	import jump from '$lib/assets/svg/jump.svg';
	import './board-styles.css';
	import type { BoardConfig, Move } from './types';
	import { generateSquareMaps } from './board';
	import { convertFenToPosition } from './fen';
	import { MoveType } from '$lib/store/types';
	import { moveSelector } from '$lib/store/stores';

	export let boardConfig: BoardConfig;
	export let isFlipped:boolean = false;
	export let squares = generateSquareMaps(
		boardConfig.dimensions,
		isFlipped
	).squares;

	export let customBoardId = 'board';
	export let position: Position = convertFenToPosition(boardConfig.fen)?.position ?? {
		piecePositions: {},
		walls: {}
	};
	export let mpSquares: Record<number, MoveType> | null = null;
	function isCursorBoardType(bType: BoardType): boolean {
		return bType !== BoardType.MovePatternView && bType != BoardType.View;
	}
	const { src, dest, piece, recentMove, legalMoves } = moveSelector;
	let markedTargets: number[] = [];

	$: squares = generateSquareMaps(boardConfig.dimensions, isFlipped).squares;

	$: {
		if($src && $piece){
		//set marked targets
		markedTargets = $legalMoves
		.filter((move)=> move.src===$src)// && $piece?.pieceType===move.piece.pieceType})
		.map((move)=> move.dest)
		markedTargets.forEach((target)=>{ squares[target].isMarkedTarget = true;})
		} else { 
			markedTargets = []; 
			for (let square in squares){
				squares[square].isMarkedTarget = false;
			}
		}
	}

	$: if($dest){

	}
	

	const isMPSquareOcc = (idx:number)=> boardConfig.boardType === BoardType.MovePatternEditor && mpSquares && mpSquares[idx] !== undefined;
</script>

<div id="wrapper">
	<div
		id={customBoardId}
		style={`--size: ${Math.max(boardConfig.dimensions.ranks, boardConfig.dimensions.files)}`}
		class={`${isCursorBoardType(boardConfig.boardType) ? 'cursor-pointer' : null}`}
	>
		{#each Array(boardConfig.dimensions.ranks * boardConfig.dimensions.files) as _, idx}
				<Square
					boardId={customBoardId}
					squareData={squares[idx]}
					color={getSquareColor(squares[idx]?.row, squares[idx]?.column, boardConfig.isFlipped)}
					piece={position.piecePositions[idx] ?? null}
					wall={position.walls[idx] ?? false}
					boardType={boardConfig.boardType}
					nonPieceSvg={isMPSquareOcc(idx) && mpSquares ? mpSquares[idx] === MoveType.Slide ? slide : jump : null}
				/>
    	{/each}
  </div>
</div>

<style>
	#board {
		display: grid;
		width: 100%;
		max-width: 80%;
		justify-items: center;
		grid-template-columns: repeat(var(--size), 1fr);
		grid-template-rows: repeat(var(--size), 1fr);
		/*border: 1px solid var(--default-dark-square);
		background-color: #090a21;/*hsla(0, 10%, 94%, 0.7);*/
		user-select: none;
		/*box-shadow: 0 14px 28px rgba(0, 0, 0, 0.25), 0 10px 10px rgba(0, 0, 0, 0.22);*/
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
