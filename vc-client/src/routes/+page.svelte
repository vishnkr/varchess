<script lang="ts">
	import { handleKeyDown } from '$lib/utils';
	import Modal from '$lib/components/shared/Modal.svelte';
	import { goto } from '$app/navigation';

	function routeToPage(route: string, replaceState: boolean) {
		goto(`/${route}`, { replaceState });
	}

	let showModal = false;
	let username = '';
	const setShowModal = (val: boolean) => (showModal = val);
	const baseUrl: string = import.meta.env.VITE_VARCHESS_SERVER_BASE;
	const cardData = [
		{
			title: 'Variable Board Sizes',
			description: 'Play on variable board sizes ranging from 5x5 to 16x16',
			bg: 'bg-[#2b2b2b]'
		},
		{
			title: 'Custom Pieces',
			description: 'Create your own pieces with unique move patterns',
			bg: 'bg-blue-600'
		},
		{
			title: 'Chess with Walls',
			description: 'Ever wondered what chess would be like with walls?',
			bg: 'bg-red-600'
		},
		{
			title: 'Predefined Variants',
			description: 'Play Wormhole, SniperChess and many more variants',
			bg: 'bg-orange-600'
		},
		{
			title: 'Variant Templates',
			description: 'Save your game templates for later use and share with friends!',
			bg: 'bg-yellow-600'
		}
	];
	async function handleSubmit() {
		if (username?.length == 0) {
		}
		const response = await fetch(`${baseUrl}/health`, {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json'
			}
			//body: JSON.stringify({username})
		});
		showModal = false;
	}
</script>

<svelte:head>
	<title>Varchess - Create and Play Custom Chess Variants</title>
</svelte:head>

<section class="font-inter bg-[#0a0c13] dark:bg-[#1d2a35] dark:text-white flex-grow">
	<div class=" max-w-6xl mx-auto py-5 sm:py-24 px-4 sm:px-6 lg:px-8">
		<div class="isolate">
			<div class="relative px-6 lg:px-8 py-16 sm:py-8">
				<div class=" pb-8 ">
					<div class="flex space-x-8 items-center justify-start">
						<div class="text-left w-[38rem]">
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
								<span class="btn-custom-1" on:click={() => setShowModal(true)}>
									Try Quick Play
								</span>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
	<div class="grid grid-cols-2 gap-4 mx-auto py-5 text-white sm:py-24 px-4 sm:px-6 lg:px-8">
		{#each cardData as card}
			<div
				class="block p-6 dark:bg-white border dark:border-gray-200 rounded-lg shadow dark:hover:bg-gray-100 bg-gray-800 border-gray-700 hover:bg-gray-700"
			>
				<h5 class="mb-2 text-2xl font-bold tracking-tight dark:text-gray-900 text-white">
					{card.title}
				</h5>
				<p class="font-normal dark:text-gray-700 text-gray-400">{card.description}</p>
			</div>
		{/each}
	</div>
	<Modal isOpen={showModal} on:close={() => setShowModal(false)}>
		<div class="grid grid-cols-1 text-white">
			<div class="text-center px-3 py-3 bg-[#1d2a35] rounded-md">
				<h3>Create a room now to play with friends!</h3>
				<h5>
					NOTE: QuickPlay Variants are limited to 8x8 boards or smaller. Login to play on larger
					boards.
				</h5>
			</div>
			<form method="POST" on:submit|preventDefault={handleSubmit}>
				<div class="text-zinc-700 p-5 grid grid-cols-1 grid-rows-2 items-center">
					<label for="username">
						Enter Username
						<input
							type="text"
							bind:value={username}
							class="rounded-md border border-gray-300 px-4 py-2 focus:border-blue-300 outline-none"
							name="username"
						/>
					</label>
				</div>
				<div>
					<button class="mt-10 flex items-center justify-center gap-x-6">
						<!-- svelte-ignore a11y-click-events-have-key-events -->
						<span
							class="btn-custom-1"
							on:click={handleSubmit}
							on:keydown={(e) => handleKeyDown(e, handleSubmit)}
						>
							Create Room
						</span>
					</button>
				</div>
			</form>
		</div>
	</Modal>
</section>
