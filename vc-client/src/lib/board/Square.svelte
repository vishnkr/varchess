<script lang="ts">
	import './board-styles.css';
	import type { BoardType, IPiece, SquareColor, SquareInfo } from './types';
	import wall from '$lib/assets/svg/wall.svg';
	import { editorSettings, editorMaxBoard } from './stores';

	export let squareData: SquareInfo;
	export let editable: boolean;
	export let interactive: boolean;

	export let color: SquareColor;
	export let piece: IPiece | null = null;
	export let disabled: boolean;
	export const mode: BoardType = 'game';

	function getPieceClass(piece: IPiece) {
		return piece.color.charAt(0).toLowerCase() + piece.pieceType.charAt(0).toLowerCase();
	}

	let pieceEl: HTMLElement;
	let squareEl: HTMLElement;
  	let hover: boolean = false;
  	let dragOver:boolean=false;
	function handleDragStart(e: DragEvent) {
		  let dragInfo = { idx: squareData.squareIndex, piece };
		  e.dataTransfer?.setData('dragInfo', JSON.stringify(dragInfo));
		  pieceEl.style.opacity = '0.4';
	}

	function handleDragEnd(e: DragEvent) {
		e.preventDefault();
		pieceEl.style.opacity = '1';
		piece = null;
	}

	function handleDragOver(e: DragEvent) {
		e.preventDefault();
    if(!piece && !disabled && interactive){
      dragOver=true;
    }
	}

	function onDrop(e: DragEvent) {
    e.preventDefault();
		const data = e.dataTransfer?.getData('dragInfo');
		if (data && !disabled && interactive) {
			var obj = JSON.parse(data);
			piece = obj.piece;
			dragOver = false;
			editorMaxBoard.updatePieceInfo(squareData.row, squareData.column, {
						isPiecePresent: true,
						piece
			});
		}
		
	}

	function handleClick(e: MouseEvent) {
		e.preventDefault();
		if (editable) {
			if ($editorSettings.disableSelected) {
				disabled = !disabled;
				editorMaxBoard.updatePieceInfo(squareData.row, squareData.column, {
					isPiecePresent: false,
					disabled,
					piece: null
				});
			} else {
				if ($editorSettings.pieceSelection)
        editorMaxBoard.updatePieceInfo(squareData.row, squareData.column, {
						isPiecePresent: piece ? false : true,
						piece: piece ? null : $editorSettings.pieceSelection
			  });
			}
		}
	}

  const isViewableOnly = () => !editable && !interactive;

  function handleMouseEnter(e:MouseEvent){
    e.preventDefault()
    if(piece && !isViewableOnly()){
      hover=true;
    }
  }

  function handleMouseLeave(e:MouseEvent){
    e.preventDefault()
    hover=false;
  }

  function handleDragLeave(e:DragEvent){
    e.preventDefault()
    dragOver = false;
    hover=false;
  }

</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<div
	class={`relative`}
	style="--x:{squareData.gridX}; --y:{squareData.gridY};"
	data-square-color={color}
	id={`s-${squareData.squareIndex}`}
	bind:this={squareEl}
	on:dragover={handleDragOver}
	on:drop={onDrop}
	on:click={handleClick}
	class:hover 
	class:dragOver
	on:mouseenter={handleMouseEnter}
	on:mouseleave={handleMouseLeave}
	on:dragleave={handleDragLeave}
>
	{#if piece}
		<div
			class={`absolute bg-piece ${interactive ? 'draggable' : null} w-full h-full ${piece ? getPieceClass(piece) : ''}`}
			draggable={interactive}
			id={`p-${squareData.squareIndex}`}
			bind:this={pieceEl}
			on:dragstart={handleDragStart}
			on:dragend={handleDragEnd}
		/>
	{:else if disabled}
		<div class="absolute inset-0 flex items-center justify-center bg-red-400">
			<img draggable={false} src={wall} alt="disabled" class="w-full h-full" />
		</div>
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

  .hover{
    background-color: var(--default-hover-square);
  }
  .dragOver{
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
		--p-label-color: var(--default-light-square);
		--p-square-color-hover: var(--square-color-dark-hover);
		--p-move-target-marker-color: var(--move-target-marker-color-dark-square);
		--p-square-color-active: var(--square-color-dark-active);
		--p-outline-color-active: var(--outline-color-dark-active);
	}

	[data-square-color='light'] {
		--square-color: var(--default-light-square);
		--p-label-color: var(--default-dark-square);
		--p-square-color-hover: var(--square-color-light-hover);
		--p-move-target-marker-color: var(--move-target-marker-color-light-square);
		--p-square-color-active: var(--square-color-light-active);
		--p-outline-color-active: var(--outline-color-light-active);
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
