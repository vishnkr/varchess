<script lang="ts">
	import { page } from '$app/stores';
	import Board from '$lib/board/Board.svelte';
	import { BoardType, type BoardConfig } from '$lib/board/types';
	import Chat from '$lib/components/Chat.svelte';
	import Tabs from '$lib/components/shared/Tabs.svelte';
	import { onMount } from 'svelte';
	import { configStore, wsStore } from '$lib/store/stores';
	import { camelToSnake } from '$lib/utils';
	import { goto } from '$app/navigation';

	let stonkfish: typeof import ('stonkfish-wasm');
	let chesscore: unknown;
	let boardConfig: BoardConfig = {
		fen: 'rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR',
		dimensions: { ranks: 8, files: 8 },
		editable: false,
		interactive: true,
		isFlipped: false,
		boardType: BoardType.GameBoard
	};
	let isWaiting: boolean = $page.data.isWaiting;;
	let gameId: string = $page.data.gameId;

	let mpBoardConfig: BoardConfig = { ...boardConfig, interactive: false };
	let activeItem = 'Chat';
	
	let items = ['Chat', 'Move Patterns'];
	const tabChange = (e: CustomEvent<string>) => (activeItem = e.detail);

	// --------------  Used by Waiting Page ------------
	let copiedToClipboard = false;
	let playAsWhite = true;

	const copyToClipboard = () => {
		copiedToClipboard = true;
    	setTimeout(() => {
     	 copiedToClipboard = false;
    	}, 2000);
		navigator.clipboard.writeText(gameId);
	};
	$: {
		if (!$wsStore || ($wsStore && $wsStore.readyState!== WebSocket.OPEN)){ goto('/home');}
	}
	onMount(async () => {
		stonkfish = await import('stonkfish-wasm');
		await stonkfish.default();	
		if ($configStore) {
			const config_json = JSON.stringify(camelToSnake($configStore))
			chesscore = new stonkfish.ChessCoreLib(config_json)
		}
	});

</script>

<svelte:head>
	<title>Play - Varchess</title>
</svelte:head>

<div class="font-inter text-zinc-90 flex-grow">
	{#if !isWaiting}
	<div class="flex-1 flex m-4 lg:flex-row flex-col">
		<div class="text-white rounded-md lg:w-8/12 p-3">
			<div>
				<Board {boardConfig} />
			</div>
		</div>

		<div class="bg-zinc-700 rounded-md lg:w-4/12 p-3">
			<div class="p-2 flex flex-grow justify-between text-white">
				<button
					class="flex gap-1 items-center justify-center rounded-md bg-black text-white hover:bg-gray-400 px-4 py-2 shadow-md"
				>
					<i class="fa-solid fa-repeat fa-lg" style="color: #ffffff;" />
					<span class="hidden md:inline"> Flip </span>
				</button>
				<button
					class="flex gap-1 items-center justify-center rounded-md bg-orange-600 text-white hover:bg-gray-400 px-4 py-2 shadow-md"
				>
					<i class="fa-solid fa-right-from-bracket fa-lg" style="color: #ffffff;" />
					<span class="hidden md:inline"> Exit </span>
				</button>
				<button
					class="bg-blue-600 flex gap-1 items-center justify-center rounded-md text-white hover:bg-gray-400 px-4 py-2 shadow-md"
				>
					<i class="fa-solid fa-handshake-simple fa-lg" style="color: #ffffff;" />
					<span class="hidden md:inline"> Draw </span>
				</button>
				<button
					class="bg-red-600 flex gap-1 items-center justify-center rounded-md text-white hover:bg-gray-400 px-4 py-2 shadow-md"
				>
					<i class="fa-solid fa-flag fa-lg" style="color: #ffffff;" />
					<span class="hidden md:inline"> Resign </span>
				</button>
			</div>
			<div
				class="border-b border-gray-200 bg-black rounded-md dark:border-gray-700 flex flex-col text-center"
			>
			<div class="rounded-md bg-gray-600 p-2 mb-4">
				<h3 class="text-black bg-white rounded-sm m-1 p-1">White - Player-1</h3>
				<h3 class="text-white bg-black rounded-sm m-1 p-1">Black - Player 2</h3>
			  </div>
				<div class="flex justify-center py-4">
					
					<Tabs {activeItem} {items} on:tabChange={tabChange} />
				</div>
				<div class="p-2 mx-1">
					{#if activeItem === 'Chat'}
						<Chat />
					{:else if activeItem === 'Move Pattern'}
						<Board boardConfig={mpBoardConfig} />
					{/if}
				</div>
			</div>
		</div>
	</div>
	{:else}
	<div class="flex flex-col items-center justify-center h-screen">
		<div class="m-4">
			<h1 class="text-xl text-white">Waiting for opponent...</h1>
		</div>
		<!-- Shareable input and Copy button -->
		<div class="mb-4 flex flex-col items-center">
		  <label for="shareableUrl" class="text-white mb-2">Share Game ID</label>
		  <div class="flex items-center">
			<input type="text" id="shareableUrl" value="{gameId}" class="w-48 border p-2 rounded-md mr-2" readonly />
			<button on:click={copyToClipboard} class="bg-blue-500 text-white px-2 py-1 rounded-md">Copy</button>
		  </div>
		  {#if copiedToClipboard}
			<p class="text-green-500 mt-2">Copied to clipboard</p>
		  {/if}
		</div>
		
		<div class="grid grid-cols-2">
			<!-- svelte-ignore a11y-click-events-have-key-events -->
			<div
				class={`flex items-center rounded-md p-4 m-1.5 bg-white hover:bg-gray-500 hover:text-white  border border-gray-200 dark:border-gray-700 cursor-pointer`}
				on:click={() => playAsWhite = true}
			>
				<input
					class="cursor-pointer w-4 h-4 text-black-600 bg-white border-gray-300 focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
					type="radio"
					value="White"
					name="color"
					checked={playAsWhite}	
				/>
				<label class="ml-2 cursor-pointer" for="White">Play as White</label>
			</div>
			<!-- svelte-ignore a11y-click-events-have-key-events -->
			<div
				class="flex items-center rounded-md pl-4 m-1.5 text-white bg-black hover:bg-gray-500 border border-gray-200 dark:border-gray-700 cursor-pointer"
				on:click={() => playAsWhite = false}
			>
				<input
					class="cursor-pointer w-4 h-4 text-white bg-black border-gray-300 focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
					type="radio"
					value="Black"
					name="color"
					checked={!playAsWhite}
				/>
				<label class="ml-2 cursor-pointer" for="Black">Play as Black</label>
			</div>
		</div>
	  </div>

	{/if}
</div>

