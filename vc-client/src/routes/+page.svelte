<script lang="ts">

	import { BoardType, type BoardConfig } from '$lib/board/types';
	import Board from '$lib/board/Board.svelte';
	import { displayAlert, getErrorMessage } from '$lib/store/alert';
	export let boardConfig: BoardConfig = {
		fen: 'rdbq1bn2/pp..pkpv1/p3ppp1p/9/4P4/P2PDBR.B/R.BQ1BKN1',
		dimensions: { ranks: 7, files: 9 },
		boardType: BoardType.View
	};

	const cardData = [
		{
			title: 'Variable Board Sizes',
			description: 'Play on variable board sizes ranging from 5x5 to 16x16',
			bg: 'bg-gradient-to-r from-[#1CB5E0] to-[#000046]'//'bg-gradient-to-r from-[#cc5333] to-[#23074d]'
		},
		{
			title: 'Custom Pieces',
			description: 'Create your own pieces with unique move patterns',
			bg: 'bg-gradient-to-r from-[#23074d] to-[#cc5333]'//'bg-gradient-to-r from-[#734b6d] to-[#42275a]'//'bg-black-900'//'bg-gradient-to-r from-[#1D4350] to-[#A43931]'
		},
		{
			title: 'Chess with Walls',
			description: 'Use walls to remove access to certain squares on the board',
			bg: 'bg-gradient-to-r from-[#cc5333] to-[#23074d]'//'bg-gradient-to-r from-[#42275a] to-[#734b6d]'//'bg-black-900'
		},
		{
			title: 'Predefined Variants',
			description: 'Play Wormhole, ArcherChess and many more variants',
			bg: 'bg-gradient-to-r from-[#000046] to-[#1CB5E0]'
		},
		{
			title: 'Variant Templates',
			description: 'Save your game templates for later use and share with friends!',
			bg: 'bg-gradient-to-r from-[#1CB5E0] to-[#000046]'//'bg-gradient-to-r from-[#FC4A1A] to-[#F7B733]'//'bg-gradient-to-r from-[#1D4350] to-[#A43931]'
		},
		{
			title: 'Play against Stonkfish Bot',
			description: 'Coming soon!',
			bg: 'bg-gradient-to-r from-[#23074d] to-[#cc5333]'//'bg-gradient-to-r from-[#000046] to-[#1CB5E0]'//'bg-gradient-to-r from-[#1D4350] to-[#A43931]'
		}
	];
</script>

<svelte:head>
	<title>Varchess - Create and Play Custom Chess Variants</title>
</svelte:head>

<section class="font-inter dark:bg-[#0a0c13] dark:text-white flex-grow">
	<div class=" max-w-6xl mx-auto sm:py-5 px-4 sm:px-6 lg:px-8 bg-">
		<div class="isolate">
			<div class="relative px-6 lg:px-8 py-16 sm:py-8">
				<div class="pb-1">
					<div class="flex flex-col md:flex-row space-x-8 items-center justify-start">
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
								<a href="/login"><span class="btn-custom-1"> Play now! </span></a>
							</div>
						</div>
						
						<div class="flex-1 sm:order-1">
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
				class="block p-6 text-center dark:bg-white border dark:border-gray-200 rounded-lg shadow dark:hover:bg-gray-100 {card.bg} gray-800 border-gray-700 hover:bg-gray-700"
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
</section>
