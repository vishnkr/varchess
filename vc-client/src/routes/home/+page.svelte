<script>
	import { displayAlert } from '$lib/store/alert.ts';
	import { me,roomId,createWebSocket } from '$lib/store/stores.ts';
	import { goto } from '$app/navigation';
	let templates = [
		{num:0, name:"Template 0"},
		{num:1, name:"Template 1"},
		{num:2, name:"Template 2"},
		{num:3, name:"Template 3"},
		{num:4, name:"Template 4"},
		{num:5, name:"Template 5n"}
	]
	export let data;
	let roomCode;
	let user;
	let isLoggedIn;
	let username;
	const baseUrl = import.meta.env.VITE_SERVER_BASE;
	$: {
		({ user, isLoggedIn } = data);
		username = user?.username || '';
	}
	async function handleSubmit() {
		if (username?.length == 0) {
		}
		try{
			const response = await fetch(`${baseUrl}/rooms`, {
				method: 'POST',
				headers: {
				'Content-Type': 'application/json'
				},
				body: JSON.stringify({username}) 
			});
			if (response.ok){

					const data = await response.json();
					if(data.roomId) {
						let ws = await createWebSocket(data.roomId,username);
						me.set({username:username})
						roomId.set(data.roomId)
						goto(`/editor/${data.roomId}`)
					} 			
				
			} else {
						displayAlert('Unable to create room. Please try again later.','DANGER',6000)
			}
		} catch(e){
			displayAlert(e.message,'DANGER',7000)
		}
	}

	async function joinRoom(){
		try{
			let ws = await createWebSocket(roomCode,username);
            me.set({id:0,isHost:false,role:0,username:username})
			roomId.set(roomCode)
			goto(`/editor/${roomCode}`)
		} catch(e){
			displayAlert('Unable to join room','DANGER',7000)
		}
	}
</script>

<svelte:head>
	<title>Home - Varchess</title>
</svelte:head>

<div class="flex h-full w-full items-center justify-center">
	<div class="mx-auto flex items-center justify-center flex-col p-4">
		<div class="rounded text-white">
			<button class="bg-blue-600 hover:bg-blue-700 text-white font-bold m-2 py-2 px-4 rounded" on:click={handleSubmit}>
				<i class="fa-solid fa-plus"/> 
				Create New Game 
			</button>
		</div>
		<span class="text-white">OR</span>
		<div class="my-3">
			<input type="text" placeholder="Enter Room Code" class="border border-gray-300 px-4 py-2 rounded-l max-w-64" bind:value={roomCode}>
			<button class="bg-blue-600 hover:bg-blue-800 text-white font-bold py-2 px-4 rounded-r" on:click={joinRoom}>Join Room</button>
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

