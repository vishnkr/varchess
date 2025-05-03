<script lang="ts">
	import { gameState } from '$lib/store/stores';
	import { dimensions, fen, positionStore } from '$lib/store/stores';
	import { authStore } from '$lib/store/auth';
	import { BoardType, type BoardConfig, type Position } from '$lib/types';
	import Board from './Board.svelte';

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
</script>

<div class="grid grid-rows-[auto_1fr_auto] w-full gap-2">
	<!-- Top Player Name -->
	<div class="flex justify-start">
		{#if players && players.playerWhite.name && players.playerBlack.name}
			<div
				class="inline-flex items-center gap-2 min-w-[10rem] justify-center px-3 py-1 rounded font-medium text-sm
    {isFlipped ? 'bg-white text-black' : 'bg-black text-white border-white border'} bg-opacity-90"
			>
				<span class="h-2 w-2 rounded-full bg-green-500" />
				{isFlipped ? players.playerWhite.name : players.playerBlack.name}
			</div>
		{/if}
	</div>

	<!-- Board Wrapper -->
	<div class="w-full">
		{#if boardConfig && position}
			<Board {boardConfig} {position} {isFlipped} />
		{:else}
			<p class="text-center text-gray-500 mt-4">Waiting for game to start...</p>
		{/if}
	</div>

	<!-- Bottom Player Name -->
	<div class="flex justify-end">
		{#if players && players.playerWhite.name && players.playerBlack.name}
			<div
				class="inline-flex items-center gap-2 min-w-[10rem] justify-center px-3 py-1 rounded font-medium text-sm
    {isFlipped ?  'bg-black text-white border-white border' : 'bg-white text-black'} bg-opacity-90"
			>
				<span class="h-2 w-2 rounded-full bg-green-500" />
				{isFlipped ? players.playerBlack.name : players.playerWhite.name}
			</div>
		{/if}
	</div>
</div>
