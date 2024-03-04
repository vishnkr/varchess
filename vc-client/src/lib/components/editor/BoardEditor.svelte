<script lang="ts">
	import down from '$lib/assets/svg/down.svg';
	import up from '$lib/assets/svg/up.svg';
	import right from '$lib/assets/svg/right.svg';
	import left from '$lib/assets/svg/left.svg';
	import type { Dimensions } from '$lib/board/types';
	import { createEventDispatcher, onMount } from 'svelte';
	import { COLOR_THEMES } from '$lib/utils';
	// @ts-ignore
	import Switch from 'svelte-switch';
	import { boardEditor } from '../../store/editor';
	let checkedValue = false;
	export let dimensions: Dimensions;
	const dispatch = createEventDispatcher();
	// @ts-ignore
	function handleChange(e) {
		const { checked } = e.detail;
		checkedValue = checked;
		boardEditor.update((val) => ({ ...val, isWallSelectorOn: checked }));
	}
	let maxDimension = 16;

	let boardTheme: string;
	onMount(()=>{
		boardTheme = sessionStorage.getItem("board-theme") ?? 'Default';
		updateColors();
	});
	
	const changeTheme = ()=>{sessionStorage.setItem("board-theme",boardTheme); updateColors();}

	function updateColors() {
		const light = COLOR_THEMES[boardTheme].lightColor;
		const dark = COLOR_THEMES[boardTheme].darkColor;
		document.documentElement.style.setProperty('--default-light-square', light);
		document.documentElement.style.setProperty('--default-dark-square', dark);
	}

	
	function updateBoardDimensions() {
		boardEditor.update((val) => ({
			...val,
			ranks: dimensions.ranks,
			files: dimensions.files
		}));
	}
</script>

<div>
	<div>
		<button
			class="
            bg-red-600 rounded-md p-2 text-white border-white
            transform transition duration-200 hover:scale-105"
			on:click={() => dispatch('clear')}
		>
			<span>Clear Board <i class="fa-solid fa-trash" style="color: #ffffff;" /> </span>
		</button>
	</div>
	<div class="grid grid-rows-1 md:grid-cols-2 shadow-md">
		<div>
			<div class="bg-white py-2 rounded-md">
				<h3 class="text-xl font-semibold">Board Width : {dimensions.files}</h3>
				<input
					class="cursor-pointer"
					type="range"
					min={5}
					max={maxDimension}
					bind:value={dimensions.files}
					on:input={updateBoardDimensions}
				/>
			</div>
			<div class="bg-white py-2 rounded-md">
				<h3 class="text-xl font-semibold">Board Height : {dimensions.ranks}</h3>
				<input
					class="cursor-pointer"
					type="range"
					min={5}
					max={maxDimension}
					bind:value={dimensions.ranks}
					on:input={updateBoardDimensions}
				/>
			</div>
		</div>
		<div class="flex flex-col justify-between items-center h-auto p-2">
			<h3 class="text-xl">Shift Board</h3>
			<button class="dbtn" on:click={() => dispatch('shift', 'up')}
				><img class="svg" src={up} alt="Shift up" /></button
			>
			<div class="flex-1">
				<button class="dbtn" on:click={() => dispatch('shift', 'left')}
					><img class="svg" src={left} alt="Shift left" /></button
				>
				<button class="dbtn" on:click={() => dispatch('shift', 'right')}
					><img class="svg" src={right} alt="Shift right" /></button
				>
			</div>
			<button class="dbtn" on:click={() => dispatch('shift', 'down')}
				><img class="svg" src={down} alt="Shift down" /></button
			>
		</div>
	</div>
	<!-- svelte-ignore a11y-label-has-associated-control -->
	<label class="relative inline-flex items-center cursor-pointer">
		<span class="m-3 text-md font-medium text-gray-900 dark:text-gray-300"
			>Toggle Wall Selector</span
		>
		<Switch on:change={handleChange} checked={checkedValue} />
	</label>
	<div class="flex items-center">
		<span class="p-3 text-md font-medium text-gray-900 dark:text-gray-300">Theme: </span>
		<select
			class="bg-white appearance-none cursor-pointer border rounded-md py-2 px-4 pr-8 leading-tight focus:outline-none focus:ring focus:border-blue-500"
			bind:value={boardTheme}
			on:change={changeTheme}
		>
			{#each Object.keys(COLOR_THEMES) as theme}
				<option value={theme}>{theme}</option>
			{/each}
		</select>

		{#if boardTheme}
			<div class="m-1 flex items-centerr">
				<div
					class="w-6 h-6 rounded-sm"
					style="background-color: {COLOR_THEMES[boardTheme].lightColor}"
				/>
				<div
					class="w-6 h-6 rounded-sm ml-2"
					style="background-color: {COLOR_THEMES[boardTheme].darkColor}"
				/>
			</div>
		{/if}
	</div>
</div>
