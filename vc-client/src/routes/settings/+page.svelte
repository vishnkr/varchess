<script lang="ts">
	import { onMount } from 'svelte';
	import { COLOR_THEMES } from '$lib/utils/index';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Button } from '$lib/components/ui/button';

	let showPossibleMoves = false;
	let enablePremove = false;
	let boardTheme: string;

	onMount(() => {
		boardTheme = localStorage.getItem('board-theme') ?? 'Default';
		updateColors();
		const storedShowMoves = localStorage.getItem('show-possible-moves');
		if (storedShowMoves) showPossibleMoves = JSON.parse(storedShowMoves);
		const storedPremove = localStorage.getItem('enable-premove');
		if (storedPremove) enablePremove = JSON.parse(storedPremove);
	});

	function saveSettings() {
		localStorage.setItem('board-theme', boardTheme);
		localStorage.setItem('show-possible-moves', JSON.stringify(showPossibleMoves));
		localStorage.setItem('enable-premove', JSON.stringify(enablePremove));
	}

	const changeTheme = () => {
		sessionStorage.setItem('board-theme', boardTheme);
		updateColors();
	};

	function updateColors() {
		const light = COLOR_THEMES[boardTheme].lightColor;
		const dark = COLOR_THEMES[boardTheme].darkColor;
		document.documentElement.style.setProperty('--default-light-square', light);
		document.documentElement.style.setProperty('--default-dark-square', dark);
	}
</script>

<div class="text-gray-700 dark:text-gray-300">
	<h1 class="text-2xl font-bold text-center mb-2">Settings</h1>
	<div class="flex items-start justify-center min-h-screen">
		<div
			class="flex flex-col bg-white dark:bg-gray-800 shadow-md rounded-lg p-8 max-w-md w-full space-y-4"
		>
			<div class="flex items-center space-x-4">
				<span class="text-md font-medium dark:text-gray-300">Theme:</span>
				<select
					class="bg-white text-black appearance-none cursor-pointer border rounded-md py-1 px-4 pr-8 leading-tight focus:outline-none focus:ring focus:border-blue-500"
					bind:value={boardTheme}
					on:change={changeTheme}
				>
					{#each Object.keys(COLOR_THEMES) as theme}
						<option value={theme}>{theme}</option>
					{/each}
				</select>
				{#if boardTheme}
					<div class="m-1 flex items-center">
						<div
							class="w-6 h-6 rounded-sm border"
							style="background-color: {COLOR_THEMES[boardTheme]
								.lightColor}; border-color: rgba(0, 0, 0, 0.3);"
						/>
						<div
							class="w-6 h-6 rounded-sm ml-2 border"
							style="background-color: {COLOR_THEMES[boardTheme]
								.darkColor}; border-color: rgba(0, 0, 0, 0.3);"
						/>
					</div>
				{/if}
			</div>

			<div class="flex items-center space-x-2">
				<Checkbox id="show-possible-moves" bind:checked={showPossibleMoves} />
				<label for="show-possible-moves" class="text-md"
					>Show possible moves after clicking on a piece</label
				>
			</div>

			<div class="flex items-center space-x-2">
				<Checkbox id="enable-premove" bind:checked={enablePremove} />
				<label for="enable-premove" class="text-md">Enable premoves</label>
			</div>

			<Button on:click={saveSettings} class="mt-4">Save Settings</Button>
		</div>
	</div>
</div>
