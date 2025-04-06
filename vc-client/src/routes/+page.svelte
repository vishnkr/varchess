<script lang="ts">
	import { Motion } from 'svelte-motion';
	import { BoardType, type BoardConfig } from '$lib/board/types';
	import Board from '$lib/board/Board.svelte';

	export let boardConfig: BoardConfig = {
		fen: 'rdbq1bn2/pp..pkpv1/p3ppp1p/9/4P4/P2PDBR.B/R.BQ1BKN1',
		dimensions: { ranks: 7, files: 9 },
		boardType: BoardType.View
	};

	const features = [
		{
			title: 'Variable Board Sizes',
			description:
				'Play on custom boards ranging from 5x5 to 16x16. Unleash your creativity in every dimension.',
			image: '/images/variable-board.png',
			gradient: 'from-[#4facfe] to-[#00f2fe]'
		},
		{
			title: 'Custom Pieces',
			description: 'Design pieces with unique movement patterns, new tactics, and strategic flair.',
			image: '/images/custom-pieces.png',
			gradient: 'from-[#43e97b] to-[#38f9d7]'
		},
		{
			title: 'Chess with Walls',
			description:
				'Block off zones and reshape the battlefield. Walls introduce rich new tactical depth.',
			image: '/images/walls.png',
			gradient: 'from-[#ff0844] to-[#ffb199]'
		},
		{
			title: 'Original Variants',
			description:
				'Play experimental game modes like Wormhole and ArcherChess  handcrafted variants with unique twists.',
			image: '/images/variants.png',
			gradient: 'from-[#f7971e] to-[#ffd200]'
		},
		{
			title: 'Save Templates',
			description: 'Build once, play forever. Save and share your favorite setups with friends.',
			image: '/images/templates.png',
			gradient: 'from-[#a18cd1] to-[#fbc2eb]'
		},
		{
			title: 'Battle the Bot',
			description: 'Test your skills against Stonkfish, our built-in AI (coming soon!).',
			image: '/images/stonkfish.png',
			gradient: 'from-[#fd746c] to-[#ff9068]'
		}
	];
</script>

<section
	class="dark:bg-[#0a0c13] bg-lightbg text-gray-800 dark:text-white font-inter overflow-y-scroll snap-y snap-mandatory"
>
	<div class="snap-start h-screen max-w-7xl mx-auto px-4 py-24">
		<div class="grid grid-cols-1 md:grid-cols-2 items-center gap-12">
			<div>
				<h1
					class="text-5xl sm:text-6xl font-extrabold leading-tight tracking-tight text-center md:text-left"
				>
					Chess Variants,<br />
					<span
						class="text-transparent bg-clip-text bg-gradient-to-r from-blue-500 to-cyan-400 dark:from-orange-400 dark:to-yellow-500"
					>
						Redefined by You
					</span>
				</h1>
				<p class="mt-6 text-2xl text-gray-600 dark:text-gray-300 max-w-xl mx-auto md:mx-0">
					Create . Customize . Play
				</p>
				<div class="mt-8">
					<a
						href="/home"
						class="inline-block px-6 py-3 bg-black text-white dark:bg-white dark:text-black font-semibold rounded-full hover:scale-105 transition"
					>
						Play Now
					</a>
				</div>
			</div>
			<Board {boardConfig} />
		</div>
	</div>
	{#each features as feature, i}
		<div class="snap-start flex items-center justify-center px-8">
			<Motion
				initial={{ opacity: 0, y: 60 }}
				transition={{ duration: 0.6, delay: i * 0.1 }}
				class="max-w-6xl w-full grid grid-cols-1 md:grid-cols-2 gap-12 items-center"
			>
				<div class="text-left">
					<h2
						class="text-4xl font-bold mb-4 bg-clip-text text-transparent bg-gradient-to-r {feature.gradient}"
					>
						{feature.title}
					</h2>
					<p class="text-xl text-gray-600 dark:text-gray-300">
						{feature.description}
					</p>
				</div>

				<Board {boardConfig} />
			</Motion>
		</div>
	{/each}
</section>

<style>
	section::before {
		content: '';
		position: fixed;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		z-index: -1;
		background: none !important;
	}
</style>
