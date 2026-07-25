<script lang="ts">
	import { goto } from '$app/navigation';
	import { gameId, gameState, Status } from '$lib/store/stores';
	import { wsStore } from '$lib/websocket';
	import PieChart from '$lib/components/charts/PieChart.svelte';
	import HalfDoughnut from '$lib/components/charts/HalfDoughnut.svelte';
	import { authStore } from '$lib/store/auth';
	import { onMount } from 'svelte';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { toast } from '$lib/store/alert';
	import { get } from 'svelte/store';
	import { theme as themeStore } from '$lib/store/theme';
	import { fetchGames } from '$lib/api/games';
	import type { Game } from '$lib/types';
	import {
		buildDashboardStats,
		formatGameResult,
		gameIdOf,
		gameTitle,
		type DashboardStats
	} from '$lib/utils/gameLabels';

	let joining = false;
	let loadingStats = true;
	let gameIdInput = '';
	let games: Game[] = [];
	let stats: DashboardStats = {
		total: 0,
		byVariant: [],
		byResult: [
			{ name: 'Win', value: 0, color: '#22c55e' },
			{ name: 'Draw', value: 0, color: '#94a3b8' },
			{ name: 'Loss', value: 0, color: '#ef4444' }
		],
		recent: []
	};

	$: chartTheme = $themeStore;
	$: userId = $authStore.userId;
	$: stats = buildDashboardStats(games, userId);

	async function loadStats() {
		loadingStats = true;
		try {
			games = await fetchGames({ page: 1, pageSize: 100 });
		} catch (e) {
			console.error(e);
			toast.error('Could not load dashboard stats');
			games = [];
		} finally {
			loadingStats = false;
		}
	}

	onMount(() => {
		void loadStats();
	});

	const createGame = () => {
		goto('/editor');
	};

	const joinRoom = async () => {
		if (joining) return;
		const { userId: uid, accessToken } = get(authStore);
		if (!uid || !accessToken) {
			goto('/login');
			return;
		}
		const validGameId = gameIdInput.trim().toUpperCase();
		if (!validGameId) {
			toast.error('Enter a room code');
			return;
		}

		joining = true;
		try {
			gameId.set(validGameId);
			gameState.updateStatus(Status.Waiting);
			wsStore.newWebSocketConnection(`ws://${import.meta.env.VITE_WS_HOST}/play/${validGameId}`, {
				token: accessToken,
				userId: uid
			});
			await goto(`/play/${validGameId}`);
		} catch (err) {
			console.error(err);
			toast.error('Unable to join game.');
			joining = false;
		}
	};

	function openGame(g: Game) {
		const id = gameIdOf(g);
		if (id) goto(`/games/${id}`);
	}
</script>

<svelte:head>
	<title>Home - Varchess</title>
</svelte:head>

<div
	class="flex flex-col lg:flex-row min-h-[calc(100vh-8rem)] bg-lightbg dark:bg-darkbg text-gray-900 dark:text-white"
>
	<div class="w-full lg:w-2/3 p-4 space-y-4">
		<div class="bg-white dark:bg-gray-900 p-6 rounded-lg shadow-lg">
			<div class="flex items-center justify-between gap-3 mb-4">
				<h2 class="text-xl font-bold">Dashboard</h2>
				{#if !loadingStats}
					<span class="text-sm text-gray-500 dark:text-gray-400">{stats.total} games</span>
				{/if}
			</div>

			{#if loadingStats}
				<p class="text-center text-gray-500 py-16">Loading your stats…</p>
			{:else if stats.total === 0}
				<div class="text-center py-10 space-y-3">
					<p class="text-gray-600 dark:text-gray-300">No finished games yet.</p>
					<p class="text-sm text-gray-500">Create a game from the editor or join a room to get started.</p>
					<Button on:click={createGame}>Create a game</Button>
				</div>
			{:else}
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div class="bg-gray-100 dark:bg-gray-800 p-4 rounded-lg">
						<PieChart theme={chartTheme} data={stats.byVariant} />
					</div>
					<div class="bg-gray-100 dark:bg-gray-800 p-4 rounded-lg">
						<HalfDoughnut theme={chartTheme} data={stats.byResult} />
					</div>
				</div>
			{/if}
		</div>

		{#if !loadingStats && stats.recent.length > 0}
			<div class="bg-white dark:bg-gray-900 p-6 rounded-lg shadow-lg">
				<div class="flex items-center justify-between mb-3">
					<h3 class="font-semibold">Recent games</h3>
					<button
						type="button"
						class="text-sm text-sky-600 dark:text-sky-400 hover:underline"
						on:click={() => goto('/games')}>View all</button
					>
				</div>
				<ul class="divide-y divide-gray-200 dark:divide-gray-700">
					{#each stats.recent as g (gameIdOf(g))}
						<li>
							<button
								type="button"
								class="w-full text-left py-3 flex items-center justify-between gap-3 hover:bg-black/5 dark:hover:bg-white/5 rounded px-1"
								on:click={() => openGame(g)}
							>
								<span class="font-medium truncate">{gameTitle(g)}</span>
								<span class="text-sm text-gray-500 dark:text-gray-400 shrink-0"
									>{formatGameResult(g)}</span
								>
							</button>
						</li>
					{/each}
				</ul>
			</div>
		{/if}
	</div>

	<div class="w-full lg:w-1/3 p-4">
		<div
			class="bg-white dark:bg-gray-800 p-6 rounded-lg flex flex-col items-center space-y-4 shadow-lg sticky top-20"
		>
			<h2 class="text-xl font-bold mb-2">Play a Game</h2>

			<Button class="w-full" on:click={createGame}>
				<i class="fa-solid fa-plus mr-2" /> Create New Game
			</Button>

			<span class="text-gray-700 dark:text-gray-300 text-sm font-medium">OR</span>

			<div class="w-full space-y-2">
				<Label for="gameId">Room Code</Label>
				<Input
					id="gameId"
					placeholder="Enter Room Code"
					bind:value={gameIdInput}
					class="w-full uppercase"
					disabled={joining}
					on:keydown={(e) => e.key === 'Enter' && joinRoom()}
				/>
			</div>

			<Button class="w-full" on:click={joinRoom} disabled={joining || !gameIdInput.trim()}>
				{joining ? 'Joining…' : 'Join Room'}
			</Button>
		</div>
	</div>
</div>
