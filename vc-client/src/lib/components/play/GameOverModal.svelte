<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog';
	import { Button } from '$lib/components/ui/button';
	import {
		gameResult,
		gameState,
		lastGameConfig,
		moveHistory,
		seekEnd,
		seekNext,
		seekPrev,
		seekStart
	} from '$lib/store/stores';
	import MoveList from './MoveList.svelte';

	export let open = false;
	export let onHome: () => void = () => {};
	export let onClose: () => void = () => {};
	export let onReplaySetup: () => void = () => {};
	export let replaying = false;

	$: result = $gameResult;
	$: players = $gameState.players;
	$: canReplaySetup = !!$lastGameConfig && !replaying;

	function winnerLabel(): string {
		if (!result) return 'Game over';
		const w = (result.winner || '').toLowerCase();
		if (w === 'draw' || w === 'none') return 'Draw';
		if (w === 'white' || w === 'w') {
			return players?.playerWhite?.name ? `${players.playerWhite.name} (White) wins` : 'White wins';
		}
		if (w === 'black' || w === 'b') {
			return players?.playerBlack?.name ? `${players.playerBlack.name} (Black) wins` : 'Black wins';
		}
		if (w === 'unknown') return 'Game over';
		return result.winner;
	}

	function reasonLabel(): string {
		if (!result?.reason) return '';
		const r = result.reason.replace(/_/g, ' ');
		if (r.toLowerCase() === 'game over') return '';
		return r;
	}
</script>

<Dialog.Root
	open={open}
	onOpenChange={(v) => {
		open = v;
		if (!v) onClose();
	}}
>
	<Dialog.Content class="max-w-lg dark:bg-darkbg bg-white text-gray-900 dark:text-white">
		<Dialog.Header>
			<Dialog.Title class="text-2xl">Game over</Dialog.Title>
			<Dialog.Description class="text-base text-gray-600 dark:text-gray-300">
				{winnerLabel()}
				{#if reasonLabel()}
					<span class="block text-sm mt-1 capitalize opacity-80">{reasonLabel()}</span>
				{/if}
			</Dialog.Description>
		</Dialog.Header>

		<div class="grid grid-cols-3 gap-3 my-4 text-center">
			<div class="rounded border border-gray-200 dark:border-gray-600 p-3">
				<p class="text-xs uppercase text-gray-500">Moves</p>
				<p class="text-xl font-semibold">{$moveHistory.length}</p>
			</div>
			<div class="rounded border border-gray-200 dark:border-gray-600 p-3">
				<p class="text-xs uppercase text-gray-500">White</p>
				<p class="text-sm font-medium truncate">{players?.playerWhite?.name ?? '—'}</p>
			</div>
			<div class="rounded border border-gray-200 dark:border-gray-600 p-3">
				<p class="text-xs uppercase text-gray-500">Black</p>
				<p class="text-sm font-medium truncate">{players?.playerBlack?.name ?? '—'}</p>
			</div>
		</div>

		<div class="mb-3">
			<p class="text-sm font-medium mb-2">Replay</p>
			<div class="flex gap-2 justify-center mb-3">
				<Button variant="outline" size="sm" on:click={seekStart} aria-label="Start">
					<i class="fa-solid fa-backward-fast" />
				</Button>
				<Button variant="outline" size="sm" on:click={seekPrev} aria-label="Previous">
					<i class="fa-solid fa-backward-step" />
				</Button>
				<Button variant="outline" size="sm" on:click={seekNext} aria-label="Next">
					<i class="fa-solid fa-forward-step" />
				</Button>
				<Button variant="outline" size="sm" on:click={seekEnd} aria-label="End">
					<i class="fa-solid fa-forward-fast" />
				</Button>
			</div>
			<MoveList />
		</div>

		<Dialog.Footer class="flex flex-col sm:flex-row gap-2 sm:justify-end">
			<Button variant="outline" on:click={onClose}>Keep reviewing</Button>
			<Button variant="secondary" disabled={!canReplaySetup} on:click={onReplaySetup}>
				{#if replaying}
					Starting…
				{:else}
					Replay with same setup
				{/if}
			</Button>
			<Button on:click={onHome}>Back to home</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
