<script lang="ts">
	import { clearMoveSelectorStores, dest, gameState, legalMoves, moveSelectorPiece, src } from '$lib/store/stores';
	import { dimensions, fen, positionStore } from '$lib/store/stores';
	import { authStore } from '$lib/store/auth';
	import { BoardType, Color, isGameBoard, type BoardConfig, type Move, type Position } from '$lib/types';
	import Board from './Board.svelte';
	import { wsStore } from '$lib/websocket';
	import { M } from 'svelte-motion';

	export let isFlipped: boolean = false;
	let boardConfig: BoardConfig | null = null;
	let position: Position | null = null;

	$: if ($fen && $dimensions) {
		boardConfig = {
			fen: $fen,
			dimensions: $dimensions,
			boardType: BoardType.GameBoard,
			isFlipped: isFlipped
		};
	}
	$: unsubscribe = positionStore.subscribe((value) => {
		position = value;
	});
	$: players = $gameState.players;

  let markedTargets: number[] = [];
  $: {
    if ($src && $moveSelectorPiece) {
      markedTargets = $legalMoves
        .filter((move) => move.from === $src)
        .map((move) => move.to);
    } else {
      markedTargets = [];
    }
  }

	$: if($dest && $src && $moveSelectorPiece){
		console.log('dest ',$dest,$src,$moveSelectorPiece)
    const move:Move = {from:$src,to:$dest,piece:$moveSelectorPiece.color === Color.BLACK ? $moveSelectorPiece.notation : $moveSelectorPiece.notation.toUpperCase()}
   
    // {"t":type,"p":{'m':{'p':piece,'f':from,'t':to,}},}
    wsStore.sendMove(move)
    clearMoveSelectorStores()
	}
</script>

<div class="grid grid-rows-[auto_1fr_auto] w-full gap-2">
	<div class="flex justify-start">
		{#if players && players.playerWhite.name && players.playerBlack.name}
			<div
				class="inline-flex items-center gap-2 min-w-[10rem] justify-center px-3 py-1 rounded font-medium text-sm
    {isFlipped ? 'bg-white text-black' : 'bg-black text-white border-white border'} bg-opacity-90"
			>
				<span class="h-3 w-3 rounded-full bg-green-500" />
				{isFlipped ? players.playerWhite.name : players.playerBlack.name}
			</div>
		{/if}
	</div>

	<div class="w-full">
		{#if boardConfig && position}
			<Board 
        {boardConfig} 
        {position} 
        {isFlipped}
        {markedTargets}
       />
		{:else}
			<p class="text-center text-gray-500 mt-4">Waiting for game to start...</p>
		{/if}
	</div>


	<div class="flex justify-end">
		{#if players && players.playerWhite.name && players.playerBlack.name}
			<div
				class="inline-flex items-center gap-2 min-w-[10rem] justify-center px-3 py-1 rounded font-medium text-sm
    {isFlipped ?  'bg-black text-white border-white border' : 'bg-white text-black border border-gray-800'} bg-opacity-90"
			>
				<span class="h-3 w-3 rounded-full bg-green-500" />
				{isFlipped ? players.playerBlack.name : players.playerWhite.name}
			</div>
		{/if}
	</div>
</div>
