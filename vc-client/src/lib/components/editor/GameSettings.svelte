<script lang="ts">
	import { settings } from '$lib/store/settings';
	import { Switch } from '$lib/components/ui/switch';
</script>

<div class="flex flex-col p-3 gap-3">
	<label class="relative inline-flex items-center justify-between cursor-pointer gap-3">
		<span class="text-md font-medium text-gray-900 dark:text-gray-300">Show Legal Moves</span>
		<Switch
			checked={$settings.showPossibleMoves}
			on:click={() => settings.patch({ showPossibleMoves: !$settings.showPossibleMoves })}
		/>
	</label>
	<label class="relative inline-flex items-center justify-between cursor-pointer gap-3">
		<span class="text-md font-medium text-gray-900 dark:text-gray-300">Hide Chat</span>
		<Switch
			checked={$settings.hideChat}
			on:click={() => settings.patch({ hideChat: !$settings.hideChat })}
		/>
	</label>
	<label class="relative inline-flex items-center justify-between cursor-pointer gap-3">
		<span class="text-md font-medium text-gray-900 dark:text-gray-300">Enable Premoves</span>
		<Switch
			checked={$settings.enablePremove}
			on:click={() => settings.patch({ enablePremove: !$settings.enablePremove })}
		/>
	</label>
	<p class="text-xs text-gray-500">Full settings are available under Settings in the sidebar.</p>
</div>
