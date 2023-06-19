<script lang="ts">
	import Board from "$lib/board/Board.svelte";
	import type { BoardConfig } from "$lib/board/types";
	import Chat from "$lib/components/room/Chat.svelte";
	import Members from "$lib/components/shared/Members.svelte";
	import Tabs from "$lib/components/shared/Tabs.svelte";

    let boardConfig: BoardConfig = {
		fen: 'rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR',
		dimensions: { ranks: 8, files: 8 },
		editable: false,
		interactive: true,
		isFlipped: false
	};
	let mpBoardConfig:BoardConfig= {...boardConfig, interactive:false};
	let activeItem = "Chat";
	let actions:{ type:string; handler: () => void }[]=[{type:"Remove",handler:()=>{}}];
	let items = ["Chat","Move Pattern","Members"];
	const tabChange = (e: CustomEvent<string>) => (activeItem = e.detail);
</script>

<style>
	.grid-container {
	  display: grid;
	  grid-template-columns: 2fr 1fr; /* Set the desired column widths */
	}
  </style>
  
<svelte:head>
	<title>Play - Varchess</title>
</svelte:head>

<div class="font-inter text-zinc-90 flex-grow">
	<div class="flex-1 flex m-4 lg:flex-row flex-col">
		<div class="text-white rounded-md lg:w-8/12 p-3">
			<div>
				<Board {boardConfig}/>
			</div>
		</div>

		<div class="bg-zinc-700 rounded-md lg:w-4/12 p-3">
			<div class="p-2 flex flex-grow justify-between text-white">
				<button class="flex gap-1 items-center justify-center rounded-md bg-black text-white hover:bg-gray-400 px-4 py-2 shadow-md"> 
					<i class="fa-solid fa-repeat fa-lg" style="color: #ffffff;"></i>
					<span class="hidden md:inline">
						Flip
					</span>
				</button>
				<button class="flex gap-1 items-center justify-center rounded-md bg-orange-600 text-white hover:bg-gray-400 px-4 py-2 shadow-md"> 
					<i class="fa-solid fa-right-from-bracket fa-lg" style="color: #ffffff;"></i>
					<span class="hidden md:inline">
						Exit
					</span>
				</button>
				<button class="bg-blue-600 flex gap-1 items-center justify-center rounded-md text-white hover:bg-gray-400 px-4 py-2 shadow-md"> 
					<i class="fa-solid fa-handshake-simple fa-lg" style="color: #ffffff;"></i>
					<span class="hidden md:inline">
						Draw
					</span>
				</button>
				<button class="bg-red-600 flex gap-1 items-center justify-center rounded-md text-white hover:bg-gray-400 px-4 py-2 shadow-md"> 
					<i class="fa-solid fa-flag fa-lg" style="color: #ffffff;"></i>
					<span class="hidden md:inline">
						Resign
					</span>
				</button>
			</div>
			<div class="border-b border-gray-200 bg-black rounded-md dark:border-gray-700 flex flex-col text-center">
				<div class="flex justify-center py-4">
					<Tabs {activeItem} {items} on:tabChange={tabChange} />
				</div>
				<div class="p-2 mx-1">
				{#if activeItem === 'Chat'}
					<Chat />
				{:else if activeItem==="Move Pattern"}
					
						<Board boardConfig={mpBoardConfig} />
					
				{:else if activeItem==="Members"}
					<div class="rounded-md bg-slate-700 p-2">
						<h3 class="text-black bg-white rounded-sm m-1 p-1">White - Player-1</h3>
						<h3 class="text-white bg-black rounded-sm m-1 p-1">Black - Player 2</h3>
					</div>
					<Members {actions} />

				{/if}
				</div>
			</div>
		</div>
	</div>
</div>