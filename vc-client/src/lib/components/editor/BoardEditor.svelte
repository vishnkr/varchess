<script lang="ts">
	import DownIcon from '$lib/assets/svg/DownIcon.svelte';
	import UpIcon from '$lib/assets/svg/UpIcon.svelte';
	import RightIcon from '$lib/assets/svg/RightIcon.svelte';
	import LeftIcon from '$lib/assets/svg/LeftIcon.svelte';
	import type { Dimensions } from '$lib/board/types';
	import { createEventDispatcher } from 'svelte';
	import { COLOR_THEMES } from '$lib/utils/index';
	// @ts-ignore
	import { boardEditor } from '../../store/editor';
	import { Button } from '../ui/button';
	import { Switch } from '../ui/switch';
	import { settings } from '$lib/store/settings';
	let checkedValue = false;
	export let dimensions: Dimensions;
	const dispatch = createEventDispatcher();
	// @ts-ignore

	let maxDimension = 16;

	$:boardEditor.update((val) => ({ ...val, isWallSelectorOn: checkedValue }));

	function changeTheme(e: Event) {
		const value = (e.currentTarget as HTMLSelectElement).value;
		settings.patch({ boardTheme: value });
	}

	function updateBoardDimensions() {
		boardEditor.update((val) => ({
			...val,
			ranks: dimensions.ranks,
			files: dimensions.files
		}));
	}
</script>

<div class="p-4 bg-white dark:bg-darkbg2 rounded-md shadow-lg">
	<div>
		<Button
			class="bg-red-600 font-medium text-lg px-3 py-2 rounded border bg-transparent text-red-600 border-red-600 hover:bg-red-500/10 transform transition duration-200 hover:scale-105"
			on:click={() => dispatch('clear')}
		>
			<span>Clear Board <i class="fa-solid fa-trash" /> </span>
		</Button>
	</div>
	<div class="grid grid-rows-1 md:grid-cols-2 gap-4">
		<div>
			<div class="bg-white dark:bg-darkbg2 py-2 rounded-md">
				<h3 class="text-md md:text-xl font-semibold text-gray-900 dark:text-white">Board Width : {dimensions.files}</h3>
				<input
					class="cursor-pointer bg-gray-200 dark:bg-darkbg2 dark:text-white rounded-md"
					type="range"
					min={5}
					max={maxDimension}
					bind:value={dimensions.files}
					on:input={updateBoardDimensions}
				/>
			</div>
			<div class="bg-white dark:bg-darkbg2 py-2 rounded-md">
				<h3 class="text-md md:text-xl font-semibold text-gray-900 dark:text-white">Board Height : {dimensions.ranks}</h3>
				<input
					class="cursor-pointer bg-gray-200 dark:bg-darkbg2 dark:text-white rounded-md"
					type="range"
					min={5}
					max={maxDimension}
					bind:value={dimensions.ranks}
					on:input={updateBoardDimensions}
				/>
			</div>
		</div>
		<!-- Shift Board Section -->
		<div class="flex flex-col justify-between items-center h-auto p-2">
			<h3 class="text-md md:text-xl font-semibold text-gray-900 dark:text-white mb-2">Shift Board</h3>

			<Button
				class="bg-transparent text-green-600 border border-green-600 hover:bg-green-500/10 px-3 py-2 rounded my-1 transition duration-200 transform hover:scale-105"
				on:click={() => dispatch('shift', 'up')}
			>
			<svelte:component 
		  this={UpIcon}
		/>
			</Button>

			<div class="flex gap-2">
				<Button
					class="bg-transparent text-green-600 border border-green-600 hover:bg-green-500/10 px-3 py-2 rounded my-1 transition duration-200 transform hover:scale-105"
					on:click={() => dispatch('shift', 'left')}
				>
				<svelte:component 
		  this={LeftIcon}
		/>
				</Button>
				<Button
					class="bg-transparent text-green-600 border border-green-600 hover:bg-green-500/10 px-3 py-2 rounded my-1 transition duration-200 transform hover:scale-105"
					on:click={() => dispatch('shift', 'right')}
				>
				<svelte:component 
		  this={RightIcon}
		/>
			</Button>
			</div>

			<Button
				class="bg-transparent text-green-600 border border-green-600 hover:bg-green-500/10 px-3 py-2 rounded my-1 transition duration-200 transform hover:scale-105"
				on:click={() => dispatch('shift', 'down')}
			>
			<svelte:component 
		  this={DownIcon}
		/>

			</Button>
		</div>
	</div>
	<div class="mt-4">
		<!-- svelte-ignore a11y-label-has-associated-control -->
		<label class="relative inline-flex items-center cursor-pointer flex-row">
			
			<span class="m-2 text-md font-semibold text-gray-900 dark:text-white">
				Wall Selector
			</span>
			<Switch
				bind:checked={checkedValue}
			/>
		</label>
		<p class="text-sm text-gray-600 mt-1 dark:text-gray-400">
			Enable this to place/remove walls by clicking on squares.
		</p>
	</div>
	
	

	<div class="flex items-center mt-4">
		<span class="p-3 text-md font-medium text-gray-900 dark:text-white">Theme: </span>
		<select
			class="bg-white dark:bg-darkbg2 dark:text-white appearance-none cursor-pointer border rounded-md py-2 px-4 pr-8 leading-tight focus:outline-none focus:ring focus:border-blue-500"
			value={$settings.boardTheme}
			on:change={changeTheme}
		>
			{#each Object.keys(COLOR_THEMES) as theme}
				<option value={theme}>{theme}</option>
			{/each}
		</select>

		{#if $settings.boardTheme}
			<div class="m-1 flex items-center">
				<div
					class="w-6 h-6 rounded-sm"
					style="background-color: {COLOR_THEMES[$settings.boardTheme].lightColor}"
				/>
				<div
					class="w-6 h-6 rounded-sm ml-2"
					style="background-color: {COLOR_THEMES[$settings.boardTheme].darkColor}"
				/>
			</div>
		{/if}
	</div>
</div>
