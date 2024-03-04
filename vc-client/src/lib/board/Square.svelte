<script lang="ts">
	import './board-styles.css';
	import {
		BoardType,
		doesSupportDragDrop,
		isGameBoard,
		type IPiece,
		type SquareColor,
		type SquareInfo,
		isEditor,
	} from './types';
	import { editorMaxBoard } from './board';
	import { pieceEditor, boardEditor } from '$lib/store/editor';
	import wallSvg from '$lib/assets/svg/wall.svg';
	import { moveSelector } from '$lib/store/stores';
	export let squareData: SquareInfo;

	export let color: SquareColor;
	export let piece: IPiece | null = null;
	export let wall: boolean = false;
	export let boardId: string = 'board';
	export let boardType: BoardType = BoardType.GameBoard;

	let isDraggable: boolean = doesSupportDragDrop(boardType);
	export let isMarkedTarget = false;
	
	export let nonPieceSvg: string | null = null;
	function getPieceClass(piece: IPiece) {
		return piece.color.charAt(0).toLowerCase() + piece.pieceType.charAt(0).toLowerCase();
	}

	let pieceEl: HTMLElement;
	let squareEl: HTMLElement;
	let hover: boolean = false;
	let drag: boolean = false;
	let isMoveSrc:boolean = false;

	$: isMoveSrc = $src!=null && $src===squareData.squareIndex && boardType === BoardType.GameBoard
	$: isMarkedTarget = boardType === BoardType.GameBoard && (squareData.isMarkedTarget ?? false);
	
	const { src, dest, piece: moveSelectorPiece } = moveSelector;

	function handleDragStart(e: DragEvent) {
		if (isGameBoard(boardType)){
			let dragInfo = { idx: squareData.squareIndex, piece };
			e.dataTransfer?.setData('dragInfo', JSON.stringify(dragInfo));
			pieceEl.style.opacity = '0.4';
		}
		
	}

	function handleDragEnd(e: DragEvent) {
		e.preventDefault();
		if(isGameBoard(boardType)){
			pieceEl.style.opacity = '1';
			piece = null;
		}
		
	}

	function handleDrag(e: DragEvent) {
		e.preventDefault();
		/*const data = e.dataTransfer?.getData('dragInfo');*/
		
		if (!piece && !wall && isGameBoard(boardType)) {
			drag = true;
		}
		if (isEditor(boardType)){
				piece = $pieceEditor.pieceSelection;
				if($boardEditor.isWallSelectorOn){
					editorMaxBoard.updatePieceInfo(squareData.row, squareData.column, {isPiecePresent: false,wall: true,});
				} else {
					editorMaxBoard.updatePieceInfo(squareData.row, squareData.column, {
					isPiecePresent: true,
					piece: piece
					});
				}
				
			}
	}

	function onDrop(e: DragEvent) {
		e.preventDefault();
		const data = e.dataTransfer?.getData('dragInfo');
		if(data && !wall && isDraggable) {
			var obj = JSON.parse(data);
			piece = obj.piece;
			drag = false;
			editorMaxBoard.updatePieceInfo(squareData.row, squareData.column, {
				isPiecePresent: true,
				piece  
			});
		}
	}
	function handleEditorClick(){
		if ($boardEditor.isWallSelectorOn) {
					wall = !wall;
					editorMaxBoard.updatePieceInfo(squareData.row, squareData.column, {
						isPiecePresent: false,
						wall,
						piece: null
					});
				} else {
					if(wall){ 
						editorMaxBoard.updatePieceInfo(squareData.row, squareData.column, {
							isPiecePresent: false,
							wall:false
						});
						return;
					}
					if ($pieceEditor.pieceSelection)
						editorMaxBoard.updatePieceInfo(squareData.row, squareData.column, {
							isPiecePresent: piece ? false : true,
							piece: piece ? null : $pieceEditor.pieceSelection
						});
				}
	}
	
	function handleMPEditorClick(){
		if (piece) {return;}
		let selectedPiece = $pieceEditor.pieceSelection;
		let jumpOffset = [squareData.row - 4, squareData.column - 4];
		if (selectedPiece) {
			const isJumpOffsetPresent =
				$pieceEditor.movePatterns[selectedPiece.pieceType] &&
				$pieceEditor.movePatterns[selectedPiece.pieceType].jumpOffsets.some(
					(o) => o[0] === jumpOffset[0] && o[1] === jumpOffset[1]
				);
			if (isJumpOffsetPresent) {
				pieceEditor.removeJumpPattern(selectedPiece.pieceType, jumpOffset);
			} else {
				pieceEditor.addJumpPattern(selectedPiece.pieceType, jumpOffset);
			}
		}
	}
	function handleGameClick(){
		if ($src){
			//if same src square is clicked, then cancel src selection
			// else make move if dest and move is in legalmoves and send it to backend for validation(can be done in board component if dest is set)
			if ($src === squareData.squareIndex){
				$src = null
				$moveSelectorPiece = null
			} else{
				$dest = squareData.squareIndex
			}
		} else{
			//source is not selected
			$src = squareData.squareIndex
			$moveSelectorPiece = piece
		}
	}

	function handleClick(e: MouseEvent) {
		e.preventDefault();
		switch (boardType) {
			case BoardType.Editor:
				handleEditorClick();
				break;
			case BoardType.MovePatternEditor:
				handleMPEditorClick();
				break;
			case BoardType.GameBoard:
				handleGameClick();
				break;
			default:
				break;
		}
	}

	function handleMouseEnter(e: MouseEvent) {
		e.preventDefault();
		hover=true;
	}

	function handleMouseLeave(e: MouseEvent) {
		e.preventDefault();
		hover = false;
	}

	function handleDragLeave(e: DragEvent) {
		e.preventDefault();
		drag = false;
		hover = false;
	}
	/**
	 * on:drag={handledrag}
	on:dragstart={handleDragStart}
	on:drop={onDrop}
	*/
	
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<div
	class={`relative
	${isMarkedTarget ? 'marked-target' : ''} 
	${piece?'has-piece':''}
	${isMoveSrc? 'move-src':''}`
	}
	style="--x:{squareData.gridX}; --y:{squareData.gridY};"
	data-square-color={color}
	data-square
	id={`${boardId}-s-${squareData.squareIndex}`}
	class:hover
	
	bind:this={squareEl}
	on:mouseenter={handleMouseEnter}
	on:mouseleave={handleMouseLeave}
	on:click={handleClick}
	on:mouseup={()=> drag ? handleDrag : handleClick}
	on:mousedown={()=>drag=false}
	on:mousemove={()=>drag=true}
	on:dragover={handleDrag}
	on:dragleave={handleDragLeave}
>
	{#if piece}
		<div
			class={`absolute bg-piece ${isDraggable ? 'draggable' : null} w-full h-full ${getPieceClass(
				piece
			)}`}
			draggable={isDraggable}
			id={`${boardId}-p-${squareData.squareIndex}`}
			bind:this={pieceEl}			
			on:dragstart={handleDragStart}
			on:dragend={handleDragEnd}
		/>
	{:else if wall}
		<div class="absolute inset-0 flex items-center justify-center bg-red-400">
			<!-- svelte-ignore a11y-missing-attribute -->
			<img draggable={boardType === BoardType.Editor} src={wallSvg} class="w-full h-full" />
		</div>
	{:else if nonPieceSvg}
		<div class="absolute inset-0 flex items-center justify-center">
			<!-- svelte-ignore a11y-missing-attribute -->
			<img draggable={false} src={nonPieceSvg} class="w-full h-full" />
		</div>
		<slot />
	{:else}
		<slot />
	{/if}
</div>

<style>
	[data-square-color] {
		width: 100%;
		height: 0;
		padding-bottom: 100%;
		grid-column: var(--y);
		grid-row: var(--x);
		background-color: var(--square-color);
	}

	.hover {
		background-color: var(--default-hover-square);
	}
	.drag {
		background-color: var(--drag-piece-over-square);
	}

	.portal {
		animation: spin 3s linear infinite;
	}

	@keyframes spin {
		0% {
			transform: rotate(0deg);
		}
		100% {
			transform: rotate(360deg);
		}
	}

	[data-square-color='dark'] {
		--square-color: var(--default-dark-square);
		--label-color: var(--default-light-square);
		--square-color-hover: var(--square-color-dark-hover);
		--move-target-marker-color: var(--move-target-marker-color-dark-square);
		--square-color-active: var(--square-color-dark-active);
		--outline-color-active: var(--outline-color-dark-active);
	}

	[data-square-color='light'] {
		--square-color: var(--default-light-square);
		--label-color: var(--default-dark-square);
		--square-color-hover: var(--square-color-light-hover);
		--move-target-marker-color: var(--move-target-marker-color-light-square);
		--square-color-active: var(--square-color-light-active);
		--outline-color-active: var(--outline-color-light-active);
	}
	[data-square].marked-target {
		background: radial-gradient(
		var(--move-target-marker-color) var(--move-target-marker-radius),
		var(--square-color) calc(var(--move-target-marker-radius) + 1px)
		);
	}

	[data-square].move-src {
		--square-color: var(--square-color-active);
	}

	[data-square].has-piece.marked-target {
		background: radial-gradient(
		var(--square-color) var(--move-target-marker-radius-occupied),
		var(--move-target-marker-color)
			calc(var(--move-target-marker-radius-occupied) + 1px)
		);
	}

	.draggable {
		cursor: url('https://www.google.com/intl/en_ALL/mapfiles/closedhand.cur'), all-scroll;
		cursor: -webkit-grab;
		cursor: -moz-grab;
		cursor: -o-grab;
		cursor: -ms-grab;
		cursor: grab;
	}

	.draggable:active {
		cursor: url('https://www.google.com/intl/en_ALL/mapfiles/openhand.cur'), all-scroll;
		cursor: -webkit-grabbing;
		cursor: -moz-grabbing;
		cursor: -o-grabbing;
		cursor: -ms-grabbing;
		cursor: grabbing;
	}
</style>
