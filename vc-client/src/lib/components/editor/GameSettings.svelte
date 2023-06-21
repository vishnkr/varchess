<script lang="ts">
	// @ts-nocheck
	import Switch from 'svelte-switch';
	let settingsChecked = {
		disableChat: false,
		showLegalMoves: false,
		boardTheme: 'Default'
	};
	function updateColors() {
		const light = colorThemes[settingsChecked.boardTheme].lightColor;
		const dark = colorThemes[settingsChecked.boardTheme].darkColor;
		document.documentElement.style.setProperty('--default-light-square', light);
		document.documentElement.style.setProperty('--default-dark-square', dark);
	}

	const colorThemes = {
		Default: { lightColor: 'hsl(51deg 24% 84%)', darkColor: 'hsl(145deg 32% 44%)' },
		Brown: { lightColor: 'hsl(36, 81%, 84%)', darkColor: 'hsl(25, 31%, 51%)' },
		Aqua: { lightColor: 'hsl(197, 34%, 83%)', darkColor: 'hsl(217, 68%, 52%)' },
		Classic: { lightColor: 'hsl(0, 0%, 100%)', darkColor: 'hsl(0, 0%, 45%)' },
		Candy: { lightColor: 'hsl(314, 100%, 90%)', darkColor: 'hsl(328, 100%, 55%)' }
	};
</script>

<div class="flex flex-col p-3">
	<!-- svelte-ignore a11y-label-has-associated-control -->
	<label class="relative inline-flex items-center cursor-pointer">
		<span class="m-3 text-md font-medium text-gray-900 dark:text-gray-300">Show Legal Moves</span>
		<Switch checked={settingsChecked.showLegalMoves} />
	</label>
	<br />
	<!-- svelte-ignore a11y-label-has-associated-control -->
	<label class="relative inline-flex items-center cursor-pointer">
		<span class="m-3 text-md font-medium text-gray-900 dark:text-gray-300">Disable Chat</span>
		<Switch checked={settingsChecked.disableChat} />
	</label>
	<!-- svelte-ignore a11y-label-has-associated-control -->
	<div class="flex items-center">
		<span class="m-3 text-md font-medium text-gray-900 dark:text-gray-300">Board Theme: </span>
		<select
			class="appearance-none cursor-pointer border rounded-md py-2 px-4 pr-8 leading-tight focus:outline-none focus:ring focus:border-blue-500"
			bind:value={settingsChecked.boardTheme}
			on:change={updateColors}
		>
			{#each Object.keys(colorThemes) as theme}
				<option value={theme}>{theme}</option>
			{/each}
		</select>

		{#if settingsChecked.boardTheme}
			<div class="ml-4 flex items-centerr">
				<div
					class="w-6 h-6 rounded-sm"
					style="background-color: {colorThemes[settingsChecked.boardTheme].lightColor}"
				/>
				<div
					class="w-6 h-6 rounded-sm ml-2"
					style="background-color: {colorThemes[settingsChecked.boardTheme].darkColor}"
				/>
			</div>
		{/if}
	</div>
</div>
