<script lang="ts">
	import { generateUsername, handleKeyDown } from '$lib/utils';
	import Modal from '$lib/components/shared/Modal.svelte';
	import { goto } from '$app/navigation';
	import { BoardType, type BoardConfig } from '$lib/board/types';
	import Board from '$lib/board/Board.svelte';
	import { me, roomId } from '$lib/store/stores';
	import { displayAlert, getErrorMessage } from '$lib/store/alert';

	let showModal = false;
	let username = generateUsername();
	const setShowModal = (val: boolean) => {
		showModal = val;
		username = generateUsername();
	};
	const baseUrl: string = import.meta.env.VITE_SERVER_BASE;
	export let boardConfig: BoardConfig = {
		fen: 'rdbq1bn2/pp..pkpv1/p3ppp1p/9/4P4/P2PDBR.B/R.BQ1BKN1',
		dimensions: { ranks: 7, files: 9 },
		editable: false,
		interactive: true,
		boardType: BoardType.View
	};
	const cardData = [
		{
			title: 'Variable Board Sizes',
			description: 'Play on variable board sizes ranging from 5x5 to 16x16',
			bg: 'bg-green-600'
		},
		{
			title: 'Custom Pieces',
			description: 'Create your own pieces with unique move patterns',
			bg: 'bg-blue-600'
		},
		{
			title: 'Chess with Walls',
			description: 'Use walls to remove access to certain squares on the board',
			bg: 'bg-red-600'
		},
		{
			title: 'Predefined Variants',
			description: 'Play Wormhole, ArcherChess and many more variants',
			bg: 'bg-orange-600'
		},
		{
			title: 'Variant Templates',
			description: 'Save your game templates for later use and share with friends!',
			bg: 'bg-purple-600'
		}
	];

	async function handleSubmit() {
		if (username?.length == 0) {
		}
		showModal = false;
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
						me.set({username:username})
						roomId.set(data.roomId)
						goto(`/editor/${data.roomId}`)
					} 			
				
			} else {
						displayAlert('Unable to create room. Please try again later.','DANGER',6000)
			}
		} catch(e){
			displayAlert(getErrorMessage(e),'DANGER',7000)
			//displayAlert(e.message,'DANGER',7000)
		}
		
	}
</script>

<svelte:head>
	<title>Varchess - Create and Play Custom Chess Variants</title>
</svelte:head>

<section class="font-inter dark:bg-[#0a0c13] dark:text-white flex-grow">
	<div class=" max-w-6xl mx-auto sm:py-5 px-4 sm:px-6 lg:px-8">
		<div class="isolate">
			<div class="relative px-6 lg:px-8 py-16 sm:py-8">
				<div class="pb-1 ">
					<div class="flex flex-cols space-x-8 items-center justify-start">
						<div class="flex-1 text-left w-[38rem]">
							<h1 class="text-center font-bold tracking-tight text-white text-6xl">
								Chess Variants
								<span
									class="text-transparent bg-clip-text bg-gradient-to-r from-orange-500 to-blue-800"
								>
									Redefined by You
								</span>
							</h1>
							<p class="mt-8 text-2xl text-center leading-8 text-white">Create. Customize. Play</p>
								<div class="mt-10 flex items-center justify-center gap-x-6">
									<a href="/login"><span class="btn-custom-1">
										Play now!
									</span></a>
								</div>
						</div>
						<div class="flex-1">
							<Board {boardConfig} />
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
	<div
		class="grid grid-cols-1 md:grid-cols-2 gap-2 mx-auto mb-5 text-white sm:py-2 px-4 sm:px-6 lg:px-8"
	>
		{#each cardData as card}
			<div
				class="block p-6 dark:bg-white border dark:border-gray-200 rounded-lg shadow dark:hover:bg-gray-100 {card.bg} gray-800 border-gray-700 hover:bg-gray-700"
			>
				<h5 class="mb-2 text-xl md:text-2xl font-bold tracking-tight dark:text-gray-900 text-white">
					{card.title}
				</h5>
				<p class="font-normal text-white dark:text-gray-400 text-sm md:text-base">
					{card.description}
				</p>
			</div>
		{/each}
	</div>
	<Modal isOpen={showModal} on:close={() => setShowModal(false)}>
		<div class="text-white bg-[#1d2a35] py-3">
			<div class="text-center px-3 py-1">
				<h3>Create a room now to play with friends!</h3>
				<p>
					NOTE: Game templates can be downloaded for future use
				</p>
			</div>
			<form on:submit|preventDefault={handleSubmit}>
				<div class="grid grid-cols-1 grid-rows-2">
					<div class="flex justify-center items-center">
						<div class="p-1 items-center">
							<label for="username">
								Your username
								<input
									type="text"
									bind:value={username}
									class="rounded-md border bg-slate-600 text-white border-gray-300 px-4 py-2 focus:border-blue-300 outline-none"
									name="username"
								/>
							</label>
						</div>
					</div>
					<div class="flex justify-center items-center">
						<div>
							<button type="submit" class="flex items-center justify-center gap-x-6">
								<!-- svelte-ignore a11y-click-events-have-key-events -->
								<span
									class="btn-custom-1"
									on:keydown={(e) => handleKeyDown(e, handleSubmit)}
								>
									Create Room
								</span>
							</button>
						</div>
					</div>
				</div>
			</form>
		</div>
	</Modal>
</section>
