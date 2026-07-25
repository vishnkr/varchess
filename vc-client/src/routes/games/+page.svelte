<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { fetchGames } from '$lib/api/games';
	import type { Game } from '$lib/types';
	import { BoardType } from '$lib/types';
	import { toast } from '$lib/store/alert';
	import ViewBoard from '$lib/board/ViewBoard.svelte';
	import {
		boardSizeLabel,
		formatGameResult,
		gameIdOf,
		gameTitle,
		variantLabel
	} from '$lib/utils/gameLabels';

	let games: Game[] = [];
	let loading = true;

	onMount(async () => {
		try {
			games = await fetchGames();
		} catch (e) {
			console.error(e);
			toast.error('Failed to load games');
		} finally {
			loading = false;
		}
	});

	function formatDate(g: Game) {
		const raw = g.createdAt ?? g.created_at ?? g.endedAt;
		if (!raw) return '—';
		try {
			return new Date(raw).toLocaleString();
		} catch {
			return String(raw);
		}
	}

	function boardConfig(g: Game) {
		const fen = g.config?.fen ?? '8/8/8/8/8/8/8/8 w - - 0 1';
		const dims = g.config?.dimensions ?? { ranks: 8, files: 8 };
		return {
			fen,
			dimensions: dims,
			boardType: BoardType.View
		};
	}

	function openGame(g: Game) {
		const id = gameIdOf(g);
		if (!id) return;
		goto(`/games/${id}`);
	}
</script>

<svelte:head>
	<title>My Games - Varchess</title>
</svelte:head>

<div class="my-6 px-4 max-w-4xl mx-auto">
	<h3 class="text-center dark:text-white text-gray-900 mt-1 text-2xl font-bold mb-6">My Games</h3>

	{#if loading}
		<p class="text-center text-gray-500">Loading…</p>
	{:else if games.length === 0}
		<p class="text-center text-gray-500">No games yet. Create one from the editor.</p>
	{:else}
		<ul class="space-y-3">
			{#each games as game (gameIdOf(game))}
				<li>
					<button
						type="button"
						class="w-full text-left flex gap-4 items-center rounded-lg border border-gray-300 dark:border-gray-600 px-3 py-3
							dark:text-white hover:bg-black/5 dark:hover:bg-white/5 transition cursor-pointer"
						on:click={() => openGame(game)}
					>
						<div
							class="shrink-0 w-24 h-24 sm:w-28 sm:h-28 rounded overflow-hidden bg-gray-700 flex items-center justify-center pointer-events-none"
						>
							{#if game.config?.fen}
								<div class="w-full scale-90">
									<ViewBoard boardConfig={boardConfig(game)} />
								</div>
							{:else}
								<span class="text-xs text-gray-400 px-2 text-center">No board preview</span>
							{/if}
						</div>

						<div class="min-w-0 flex-1">
							<p class="font-semibold truncate">{gameTitle(game)}</p>
							<p class="text-sm text-gray-500 dark:text-gray-400 mt-0.5">
								{variantLabel(game.config?.variantType)}
								{#if boardSizeLabel(game)}
									· {boardSizeLabel(game)}
								{/if}
							</p>
							<p class="text-sm text-gray-600 dark:text-gray-300 mt-1">{formatGameResult(game)}</p>
							{#if game.playerNames}
								<p class="text-xs text-gray-500 mt-1 truncate">
									{game.playerNames.w || game.playerNames.white || 'White'}
									vs
									{game.playerNames.b || game.playerNames.black || 'Black'}
								</p>
							{/if}
						</div>

						<div class="shrink-0 text-right">
							<span class="text-xs sm:text-sm text-gray-500 block">{formatDate(game)}</span>
							<span class="text-xs text-blue-600 dark:text-blue-400 mt-2 inline-block">Review</span>
						</div>
					</button>
				</li>
			{/each}
		</ul>
	{/if}
</div>
