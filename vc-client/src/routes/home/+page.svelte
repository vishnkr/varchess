<script lang="ts">
	import { goto } from "$app/navigation";
	import { type ConnectParams, wsStore, gameId } from "$lib/store/stores";


	let templates = [
		{num:0, name:"Template 0"},
		{num:1, name:"Template 1"},
		{num:2, name:"Template 2"},
		{num:3, name:"Template 3"},
		{num:4, name:"Template 4"},
		{num:5, name:"Template 5n"}
	]
	export let data;
	let userId;
	let gameIdInput: string;
	let username;
	$: {
		({ userId } = data);
		username = userId;
	}

	const joinGame = ()=>{
		if (gameIdInput){
			const url = `ws://${import.meta.env.VITE_WS_HOST}/ws`;
			const params:ConnectParams = {
				connectType:"join",
				sessionId: "sdf",
				gameId: gameIdInput
			};
			wsStore.newWebSocketConnection(url,params);
		}
	}
	$: {
    	if ($gameId !== null) { goto(`/game/${$gameId}`); }
  	}
</script>

<svelte:head>
	<title>Home - Varchess</title>
</svelte:head>

<div class="flex h-full w-full items-center justify-center">
	<div class="mx-auto flex items-center justify-center flex-col p-4">
		<div class="rounded text-white">
			<a href="/editor">
				<button 
					class="bg-blue-600 hover:bg-blue-700 text-white font-bold m-2 py-2 px-4 rounded"
					>
					<i class="fa-solid fa-plus"/> 
					Create New Game 
				</button>
			</a>
		</div>
		<span class="text-white">OR</span>
		<div class="my-3">
			<form method="POST" action="?/joinRoom">
			<input type="text" name="gameId" bind:value={gameIdInput} placeholder="Enter Room Code" class="border border-gray-300 px-4 py-2 rounded-l max-w-64">
			<button class="bg-blue-600 hover:bg-blue-800 text-white font-bold py-2 px-4 rounded-r" >Join Room</button>
			</form>
		</div>
		<div class="mx-3">
			<h3 class="text-center text-white my-3 text-2xl font-bold">My Templates</h3>
			<table class="border border-gray-300 w-full">
				<thead>
					<tr>
						<th class="text-white bg-gray-700 border border-gray-300 px-4 py-2">Variant Template</th>
						<th class="text-white bg-gray-700 border border-gray-300 px-4 py-2">Actions</th>
					</tr>
				</thead>
				<tbody>
					{#each templates as template}
					<tr class="bg-white">
						<td class="text-gray-900 border text-center border-gray-300 px-4 py-2">{template.name}</td>
						<td class="text-gray-900 border flex justify-center border-gray-300 px-4 py-2">
							<button class="bg-green-600 hover:bg-green-800 text-white font-bold py-2 px-4 rounded mr-2">Play</button>
							<button class="bg-blue-600 hover:bg-blue-800 text-white font-bold py-2 px-4 rounded mr-2">Edit</button>
							<button class="bg-red-600 hover:bg-red-800 text-white font-bold py-2 px-4 rounded mr-2">Delete</button>
						</td>
					</tr>
					{/each}
				</tbody>
			</table>
		</div>
		
	</div>
</div>

