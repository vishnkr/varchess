<script lang="ts">
	import { dimensions, fen, positionStore } from '$lib/store/stores';
	import { BoardType, type BoardConfig, type Position } from '$lib/types';
	import Board from './Board.svelte';
	import { convertFenToPosition } from './fen';

    export let boardConfig: BoardConfig; 
	export let position: Position = convertFenToPosition(boardConfig.fen)?.position ?? {
		piecePositions: {},
		walls: {}
	};

</script>

{#if boardConfig && position}
  <Board {boardConfig} {position} />
{:else}
  <p class="text-center text-gray-500 mt-4">Waiting for game to start...</p>
{/if}
