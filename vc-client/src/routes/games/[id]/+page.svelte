<script lang="ts">
	import { beforeNavigate, goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { browser } from '$app/environment';
	import { onMount } from 'svelte';
	import { fetchGame } from '$lib/api/games';
	import { toast } from '$lib/store/alert';
	import GameBoard from '$lib/board/GameBoard.svelte';
	import MoveList from '$lib/components/play/MoveList.svelte';
	import CustomPiecesPanel from '$lib/components/play/CustomPiecesPanel.svelte';
	import WormholesPanel from '$lib/components/play/WormholesPanel.svelte';
	import { Button } from '$lib/components/ui/button';
	import {
		Status,
		gameState,
		gameResult,
		moveHistory,
		resetGameHistory,
		seekEnd,
		seekNext,
		seekPrev,
		seekStart,
		viewPly
	} from '$lib/store/stores';
	import { formatGameResult, gameTitle, boardSizeLabel, variantLabel } from '$lib/utils/gameLabels';
	import { hydrateCompletedGame } from '$lib/utils/hydrateGame';
	import type { Game } from '$lib/types';

	let loading = true;
	let game: Game | null = null;
	let isFlipped = false;
	let error: string | null = null;

	$: isReviewing = $viewPly !== null;
	$: resultBanner = game ? formatGameResult(game) : 'Game over';

	const goBack = () => {
		if (browser) goto('/games');
	};

	onMount(async () => {
		const id = $page.params.id;
		if (!id) {
			goto('/games');
			return;
		}
		try {
			game = await fetchGame(id);
			hydrateCompletedGame(game);
		} catch (e) {
			console.error(e);
			error = 'Could not load this game';
			toast.error('Failed to load game');
		} finally {
			loading = false;
		}
	});

	beforeNavigate(() => {
		resetGameHistory();
		gameState.updateStatus(Status.None);
		gameResult.set(null);
	});
</script>

<svelte:head>
	<title>{game ? gameTitle(game) : 'Game review'} - Varchess</title>
</svelte:head>

<div class="font-inter flex-grow px-3 py-4">
	{#if loading}
		<p class="text-center dark:text-white text-gray-600 p-8">Loading game…</p>
	{:else if error || !game}
		<div class="text-center p-8 dark:text-white">
			<p class="mb-4">{error ?? 'Game not found'}</p>
			<Button on:click={goBack}>Back to My Games</Button>
		</div>
	{:else}
		<div class="max-w-6xl mx-auto mb-3 flex flex-wrap items-center justify-between gap-2">
			<div>
				<button
					type="button"
					class="text-sm text-blue-600 dark:text-blue-400 mb-1 hover:underline"
					on:click={goBack}
				>
					← My Games
				</button>
				<h1 class="text-xl font-bold dark:text-white text-gray-900">{gameTitle(game)}</h1>
				<p class="text-sm text-gray-500">
					{variantLabel(game.config?.variantType)}
					{#if boardSizeLabel(game)}
						· {boardSizeLabel(game)}
					{/if}
				</p>
			</div>
		</div>

		<div class="flex m-2 lg:flex-row flex-col gap-3 max-w-6xl mx-auto">
			<div class="text-white rounded-md lg:w-8/12 p-3">
				<div
					class="mb-3 inline-flex items-center gap-2 rounded px-3 py-1.5 text-sm font-medium
					bg-emerald-100 text-emerald-900 dark:bg-emerald-500/20 dark:text-emerald-200"
				>
					{resultBanner}
					{#if isReviewing}
						<span class="opacity-80">· move {$viewPly} / {$moveHistory.length}</span>
					{/if}
				</div>

				<div class="max-w-[90%]">
					{#if game.config?.fen}
						<GameBoard {isFlipped} interactive={false} />
					{:else}
						<p class="dark:text-gray-300 text-gray-700 py-8">
							This game was saved before board setups were stored, so a replay board isn’t available.
							Finish a new game to see full history review.
						</p>
					{/if}
				</div>
			</div>

			<div
				class="bg-lightbg dark:bg-darkbg2 border border-black dark:border-lightbg rounded-md lg:w-4/12 p-3 flex flex-col gap-3"
			>
				<div class="flex flex-wrap gap-2 justify-center">
					<Button variant="outline" on:click={() => (isFlipped = !isFlipped)}>Flip</Button>
					<Button variant="outline" on:click={goBack}>Back</Button>
				</div>

				<div class="flex gap-1 justify-center">
					<Button variant="outline" size="sm" on:click={seekStart} aria-label="First move">
						<i class="fa-solid fa-backward-fast" />
					</Button>
					<Button variant="outline" size="sm" on:click={seekPrev} aria-label="Previous move">
						<i class="fa-solid fa-backward-step" />
					</Button>
					<Button variant="outline" size="sm" on:click={seekNext} aria-label="Next move">
						<i class="fa-solid fa-forward-step" />
					</Button>
					<Button variant="outline" size="sm" on:click={seekEnd} aria-label="Last move">
						<i class="fa-solid fa-forward-fast" />
					</Button>
				</div>

				<div class="flex-1 rounded-md bg-black/5 dark:bg-black/40 p-2 min-h-[12rem] overflow-y-auto flex flex-col gap-3">
					<CustomPiecesPanel />
					<WormholesPanel />
					<div class="min-h-[8rem] flex-1">
						<MoveList />
					</div>
				</div>
			</div>
		</div>
	{/if}
</div>
