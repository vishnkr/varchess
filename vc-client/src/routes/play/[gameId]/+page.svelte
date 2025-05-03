<script lang="ts">
	import Board from '$lib/board/Board.svelte';
	import { BoardType, type BoardConfig } from '$lib/board/types';
	import Chat from '$lib/components/Chat.svelte';
	import Tabs from '$lib/components/shared/Tabs.svelte';
	import { onMount } from 'svelte';
	import { templateStore, gameState, gameId, Status, clearMoveSelectorStores } from '$lib/store/stores';
	import {sendWebsocketMsg, wsStore} from '$lib/websocket';
	import { camelToSnake } from '$lib/utils/index';
	import { beforeNavigate, goto } from '$app/navigation';
	import { browser } from '$app/environment';
	import { sendDrawOffer } from '$lib/websocket.js';
	//import chessCore from '$lib/chesscore.worker.js';
	import { get } from 'svelte/store';
	import { authStore } from '$lib/store/auth.js';
	import { EventGameResign } from '$lib/types';
	import { page } from '$app/stores';
	import { Input } from '$lib/components/ui/input';
	import { Button } from '$lib/components/ui/button';
	import GameBoard from '$lib/board/GameBoard.svelte';
	let boardConfig: BoardConfig;
	let mpBoardConfig: BoardConfig;
	let activeItem = 'Chat';
	let players = $gameState.players;
	let items = ['Chat', 'Move Patterns'];
	const tabChange = (e: CustomEvent<string>) => (activeItem = e.detail);
	let isFlipped = false;
	let username: string = "fsdf";
	let isPlayer = false;
	let isMounted = false;

	
	const goHome = ()=>{
		if (browser) { 
			goto('/home');
		}
	};
	let copiedToClipboard = false;
	let currentGameId: string|null;

	//$: players = $gameState.players;

	$: {
		if ($authStore?.username && $gameState?.players) {
			username = $authStore.username;
			const { playerWhite, playerBlack } = $gameState.players;
			isPlayer = username === playerWhite.name || username === playerBlack.name;
			isFlipped = username === playerBlack.name;
		}
	}
	$: authStore.subscribe((state) => (auth = state));
	$: {
		if (!$wsStore || $gameState.status===Status.Completed) {
			goHome()
		}
	}

	const copyToClipboard = () => {
		copiedToClipboard = true;
		setTimeout(() => copiedToClipboard = false, 2000);
		if(currentGameId)
		navigator.clipboard.writeText(currentGameId);
	};

	let auth;
	
	onMount(() => {
		currentGameId = $page.params.gameId;
		const auth = localStorage.getItem('auth');
		console.log(auth,currentGameId,'f')
		if (!auth || !currentGameId) {
			goto('/home');
		}
		clearMoveSelectorStores()
		gameId.set(currentGameId);
	});
	//let chesscore;
	//const {legalMoves } = moveSelector;

	
	function clearStores(){
		wsStore.set(null);
		templateStore.removeTemplate();
		gameId.set(null);
		gameState.updateStatus(Status.None)
	}
	let dirty = true;
	
	beforeNavigate(({ cancel }) => {
		clearStores()
		/*if (dirty) {
			const confirmMessage = "Exiting this page results in loss. Are you sure you want to leave?";
			if (!confirm(confirmMessage)) {
			cancel();
			} else {
				wsStore.set(null);
				//templateStore.removeConfig();
				gameId.set(null);
				gameState.updateStatus(Status.None);
        	}
		}*/
	});

	const handleDraw = ()=>{
		if($wsStore && $gameId) sendDrawOffer($wsStore,$gameId)
	}

	const handleResign = () =>{
		if($wsStore && $gameId) sendWebsocketMsg($wsStore,EventGameResign,{gameId:$gameId})
	}
//<Board boardConfig={mpBoardConfig} />
</script>

<svelte:head>
	<title>Play - Varchess</title>
</svelte:head>

<div class="font-inter text-zinc-90 flex-grow">
	{#if $gameState.status === Status.Waiting}
		<div class="fixed inset-0 z-10 overflow-y-auto">
			<div class="flex items-center justify-center min-h-screen">
				<div class="dark:bg-darkbg bg-lightbg2 p-8 rounded-md shadow-md">
					<h1 class="text-xl dark:text-white text-gray-800 mb-4">Waiting for opponent...</h1>
					<div class="mb-4 flex flex-col items-center">
						<label for="shareableUrl" class="dark:text-white text-gray-800 mb-2">Share Game ID</label>
						<div class="flex items-center gap-2">
							<Input id="gameId" bind:value={currentGameId} type="text" readonly />
							<Button on:click={copyToClipboard} class="text-sm px-3 py-1.5 rounded border font-medium
								bg-transparent dark:text-white dark:border-white dark:hover:bg-white/10 text-black border-black hover:bg-black/10">Copy</Button>
						</div>
						{#if copiedToClipboard}
							<p class="dark:text-white text-black mt-2">Copied to clipboard</p>
						{/if}
					</div>
				</div>
			</div>
		</div>
	{:else if $gameState.status === Status.InProgress}
	<div class="flex m-4 lg:flex-row flex-col">
		<div class="text-white rounded-md lg:w-8/12 p-3">
			<div class="max-w-[90%]">
				<GameBoard {isFlipped}/>
			</div>
		</div>
		
		
		<div class="bg-lightbg dark:bg-darkbg2 border border-black dark:border-lightbg rounded-md lg:w-4/12 p-3">
			<div class="flex flex-wrap gap-2 justify-center p-4">
				<!-- Exit Button -->
				<Button
				  on:click={clearStores}
				  class="w-[calc(50%-0.5rem)] text-lg px-3 py-2 rounded border font-medium 
					bg-transparent text-orange-600 border-orange-600 hover:bg-orange-500/10 
					dark:text-orange-400 dark:border-orange-400 dark:hover:bg-orange-500/10"
				>
				  <i class="fa-solid fa-right-from-bracket mr-2" />
				  Exit
				</Button>
			  
				<!-- Draw Button -->
				<Button
				  on:click={handleDraw}
				  class="w-[calc(50%-0.5rem)] text-lg px-3 py-2 rounded border font-medium 
					bg-transparent text-blue-600 border-blue-600 hover:bg-blue-500/10 
					dark:text-blue-400 dark:border-blue-400 dark:hover:bg-blue-500/10"
				>
				  <i class="fa-solid fa-handshake-simple mr-2" />
				  Draw
				</Button>
			  
				<!-- Resign Button -->
				<Button
				  on:click={handleResign}
				  class="w-[calc(50%-0.5rem)] text-lg px-3 py-2 rounded border font-medium 
					bg-transparent text-red-600 border-red-600 hover:bg-red-600/10 
					dark:text-red-400 dark:border-red-400 dark:hover:bg-red-600/10"
				>
				  <i class="fa-solid fa-flag mr-2" />
				  Resign
				</Button>
			  
				<!-- Flip Button (Dynamic black/white styling) -->
				<Button
				  on:click={() => isFlipped = !isFlipped}
				  class="w-[calc(50%-0.5rem)] text-lg px-3 py-2 rounded border font-medium 
					bg-transparent border-black text-black hover:bg-gray-100 
					dark:border-white dark:text-white dark:hover:bg-white/10"
				>
				  <i class="fa-solid fa-repeat mr-2" />
				  Flip
				</Button>
			  </div>
			  
			<div
				class="border-b border-gray-200 bg-black rounded-md dark:border-gray-700 flex flex-col text-center"
			>
				<div class="flex justify-center py-4">
					<Tabs {activeItem} {items} on:tabChange={tabChange} />
				</div>
				<div class="p-2 mx-1">
					{#if activeItem === 'Chat'}
						<Chat />
					{:else if activeItem === 'Move Pattern'}
						nothing
					{/if}
				</div>
			</div>
		</div>
	</div>
	{:else}
		<div>{$gameState.status}</div>
	{/if}
</div>
