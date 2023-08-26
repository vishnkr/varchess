<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	export let activeItem: string;
	export let items: string[];

	const dispatch = createEventDispatcher();
</script>

<div class="tabs text-white">
	<ul>
		{#each items as item}
			<!-- svelte-ignore a11y-click-events-have-key-events -->
			<li on:click={() => dispatch('tabChange', item)}>
				<div class:active={item == activeItem}>
					{item}
				</div>
			</li>
		{/each}
	</ul>
</div>

<style>
	.tabs {
		margin-bottom: 40px;
	}

	ul {
		display: flex;
		justify-self: center;
		padding: 0;
		list-style-type: none;
	}
	li {
		margin: 0 16px;
		font-size: 18px;
		cursor: pointer;
	}
</style>
