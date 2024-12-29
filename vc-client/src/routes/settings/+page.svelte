<script lang="ts">
	import { onMount } from 'svelte';
	import { COLOR_THEMES } from '$lib/utils';
	
	let showPossibleMoves = false;
	let enablePremove = false;
    let boardTheme: string;
	// Load existing settings from local storage
	onMount(() => {
		boardTheme = localStorage.getItem("board-theme") ?? 'Default';
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
	
	const changeTheme = ()=>{sessionStorage.setItem("board-theme",boardTheme); updateColors();}

	function updateColors() {
		const light = COLOR_THEMES[boardTheme].lightColor;
		const dark = COLOR_THEMES[boardTheme].darkColor;
		document.documentElement.style.setProperty('--default-light-square', light);
		document.documentElement.style.setProperty('--default-dark-square', dark);
	}
</script>
<div>
    <h1 class="text-2xl font-bold text-center mb-2">Settings</h1>
    <div class="flex items-start justify-center min-h-screen ">

        <div class="flex flex-col bg-gray-800 shadow-md rounded-lg p-8 max-w-md w-full">
            
            
            <div class="flex items-center">
                <span class="p-3 text-md font-medium  dark:text-gray-300">Theme: </span>
                <select
                    class="bg-white text-black appearance-none cursor-pointer border rounded-md py-2 px-4 pr-8 leading-tight focus:outline-none focus:ring focus:border-blue-500"
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
    
            <div class="mb-4">
                <div class="mt-1">
                    <input type="checkbox" id="show-possible-moves" bind:checked={showPossibleMoves} class="h-4 w-4 text-indigo-600 border-gray-300 rounded focus:ring-indigo-500">
                    <span class="ml-2 text-sm">Show possible moves after clicking a piece</span>
                </div>
            </div>
    
            <div class="mb-4">
                <div class="mt-1">
                    <input type="checkbox" id="enable-premove" bind:checked={enablePremove} class="h-4 w-4 text-indigo-600 border-gray-300 rounded focus:ring-indigo-500">
                    <span class="ml-2 text-sm">Enable premoves</span>
                </div>
            </div>
    
            <button on:click={saveSettings} class="mt-4 bg-indigo-600 justify-center text-white py-2 px-4 rounded-md hover:bg-indigo-700">Save Settings</button>
        </div>
    </div>
    
</div>
