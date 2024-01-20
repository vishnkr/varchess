<script lang="ts">
	// @ts-nocheck
	import Switch from 'svelte-switch';
	let settingsChecked = {
		disableChat: false,
		showLegalMoves: false,
		autoQueen: false
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
	<label class="relative inline-flex items-center cursor-pointer">
		<span class="m-3 text-md font-medium text-gray-900 dark:text-gray-300">Auto Queen</span>
		<Switch checked={settingsChecked.autoQueen} />
	</label>
</div>
