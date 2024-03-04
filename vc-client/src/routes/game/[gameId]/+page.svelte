<script lang="ts">
	import Board from '$lib/board/Board.svelte';
	import { BoardType, type BoardConfig, type Move } from '$lib/board/types';
	import Chat from '$lib/components/Chat.svelte';
	import Tabs from '$lib/components/shared/Tabs.svelte';
	import { onDestroy, onMount } from 'svelte';
	import { configStore, wsStore, gameState, gameId, moveSelector } from '$lib/store/stores';
	import { camelToSnake } from '$lib/utils';
	import { beforeNavigate, goto } from '$app/navigation';
	import { browser } from '$app/environment';
	import init,{ChessCoreLib} from 'stonkfish';
	
	
	let boardConfig: BoardConfig;
	let mpBoardConfig: BoardConfig;
	let activeItem = 'Chat';
	let players = $gameState.players;
	let items = ['Chat', 'Move Patterns'];
	const tabChange = (e: CustomEvent<string>) => (activeItem = e.detail);
	let isFlipped = false;
	let username: string;
	let isPlayer = false;
	let isMounted = false;
	export let data;
	
	const goHome = ()=>{if (browser) { goto('/home') }};

	async function initWasm(){ await init(); isMounted = true;}
	initWasm();
	$: { if (data.username) username = data.username;}
	$: {
		if (!$wsStore) {
			goHome()
		}
	}
	let chesscore;
	const {legalMoves } = moveSelector;
	$: {
		if ($configStore) {
			const config_json = JSON.stringify(camelToSnake($configStore));
			if (isMounted){
				chesscore = new ChessCoreLib(config_json);
				let moves:Move[] = chesscore.getLegalMoves();
				
				legalMoves.set(moves)
			}
			
			boardConfig = {
				fen: $configStore.fen,
				dimensions: $configStore.dimensions,
				boardType: isPlayer ? BoardType.GameBoard : BoardType.View
			};
		}
	}
	onMount(async () => {
		initWasm();
		
		if($gameState?.players?.playerBlack === data.username){
			isPlayer = true
			isFlipped=true
		} else if ($gameState?.players?.playerWhite === data.username){
			isPlayer = true
		}
	});
	
	function clearStores(){
		wsStore.set(null);//close();
		configStore.removeConfig();
		gameId.set(null);
	}
	let dirty = true;
	
	beforeNavigate(({ cancel }) => {
		clearStores()
		if (dirty) {
			/*const confirmMessage = "Exiting this page results in loss. Are you sure you want to leave?";
			if (!confirm(confirmMessage)) {
			cancel();
			} else {
				wsStore.set(null);
				configStore.removeConfig();
				gameId.set(null);
        	}*/
		}
	});

</script>

<svelte:head>
	<title>Play - Varchess</title>
</svelte:head>

<div class="font-inter text-zinc-90 flex-grow">
	<div class="flex m-4 lg:flex-row flex-col">
		<div class="text-white rounded-md lg:w-8/12 p-3">
			<div class="max-w-[90%]">
				<Board {boardConfig} {isFlipped}/>
			</div>
		</div>

		<div class="bg-zinc-700 rounded-md lg:w-4/12 p-3">
			<div class="p-2 flex flex-grow justify-between text-white">
				<button
					on:click={()=> isFlipped = !isFlipped}
					class="flex gap-1 items-center justify-center rounded-md bg-black text-white hover:bg-gray-400 md:px-4 md:py-2 px-2 py-1 shadow-md"
				>
					<i class="fa-solid fa-repeat fa-lg" style="color: #ffffff;" />
					<span class="text-md md:text-lg"> Flip </span>
				</button>
					<button
					on:click={goHome}
					class="flex gap-1 items-center justify-center rounded-md bg-orange-600 text-white hover:bg-gray-400 md:px-4 md:py-2 px-2 py-1 shadow-md"
					>
					<i class="fa-solid fa-right-from-bracket fa-lg" style="color: #ffffff;" />
					<span class="text-md md:text-lg"> Exit </span>
					</button>
				<button
					class="bg-blue-600 flex gap-1 items-center justify-center rounded-md text-white hover:bg-gray-400 md:px-4 md:py-2 px-2 py-1 shadow-md"
				>
					<i class="fa-solid fa-handshake-simple fa-lg" style="color: #ffffff;" />
					<span class="text-md md:text-lg"> Draw </span>
				</button>
				<button
					class="bg-red-600 flex gap-1 items-center justify-center rounded-md text-white hover:bg-gray-400 md:px-4 md:py-2 px-2 py-1 shadow-md"
				>
					<i class="fa-solid fa-flag fa-lg" style="color: #ffffff;" />
					<span class="text-md md:text-lg"> Resign </span>
				</button>
			</div>
			<div
				class="border-b border-gray-200 bg-black rounded-md dark:border-gray-700 flex flex-col text-center"
			>
				<div class="rounded-md bg-gray-600 p-2 mb-4">
					<h3 class="text-black bg-white rounded-sm m-1 p-1">{players?.playerWhite}</h3>
					<h3 class="text-white bg-black rounded-sm m-1 p-1">{players?.playerBlack}</h3>
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
</div>
