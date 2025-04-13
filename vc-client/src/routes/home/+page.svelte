<script lang="ts">
	import { goto } from '$app/navigation';
	import { type ConnectParams, gameId } from '$lib/store/stores';
	import { wsStore } from '$lib/websocket';
	import PieChart from '$lib/components/charts/PieChart.svelte';
	import HalfDoughnut from '$lib/components/charts/HalfDoughnut.svelte';
	import { authStore } from '$lib/store/auth';
	import { onDestroy, onMount } from 'svelte';
	import {Button} from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	let username: string = '';
	const unsubscribe = authStore.subscribe((state) => {
		username = state.username || '';
	});

	onDestroy(() => {
		unsubscribe();
	});

	let theme = 'light';

	if (typeof window !== 'undefined') {
		const updateTheme = () => {
			theme = document.documentElement.classList.contains('dark') ? 'dark' : 'light';
		};

		updateTheme();
		const observer = new MutationObserver(updateTheme);
		observer.observe(document.documentElement, { attributes: true, attributeFilter: ['class'] });
	}

	let gameIdInput: string;

	const createGame = () => {
		goto('/editor');
	};
	const joinRoom = () => {
		if (gameIdInput.trim()) {
			gameId.set(gameIdInput.trim());
			goto('/game');
		}
	};
</script>

<svelte:head>
	<title>Home - Varchess</title>
</svelte:head>

<div
	class="flex flex-col md:flex-row h-screen bg-lightbg dark:bg-darkbg text-gray-900 dark:text-white"
>
	<!-- Left Section - Dashboard -->
	<div class="w-full md:w-2/3 p-4">
		<div class="bg-white dark:bg-gray-900 p-6 rounded-lg shadow-lg">
			<h2 class="text-xl font-bold text-center mb-4">Dashboard</h2>
			<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
				<div class="bg-gray-100 dark:bg-gray-800 p-4 rounded-lg">
					<PieChart {theme} />
				</div>
				<div class="bg-gray-100 dark:bg-gray-800 p-4 rounded-lg">
					<HalfDoughnut {theme} />
				</div>
			</div>
		</div>
	</div>

	<div class="w-full md:w-1/3 p-4">
		<div
			class="bg-white dark:bg-gray-800 p-6 rounded-lg flex flex-col items-center space-y-4 shadow-lg"
		>
			<h2 class="text-xl font-bold mb-2">Play a Game</h2>

			<Button class="w-full" on:click={createGame}>
				<i class="fa-solid fa-plus mr-2" /> Create New Game
			</Button>

			<span class="text-gray-700 dark:text-gray-300 text-sm font-medium">OR</span>

			<div class="w-full space-y-2">
				<Label for="gameId">Room Code</Label>
				<Input id="gameId" placeholder="Enter Room Code" bind:value={gameIdInput} class="w-full" />
			</div>

			<Button class="w-full" variant="secondary" on:click={joinRoom}>Join Room</Button>
		</div>
	</div>
</div>
