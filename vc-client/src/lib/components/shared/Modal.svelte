<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	export let isOpen = false;

	const dispatch = createEventDispatcher();

	function closeModal() {
		isOpen = false;
		dispatch('close');
	}
</script>

{#if isOpen}
	<!-- svelte-ignore a11y-click-events-have-key-events -->
	<div class="modal-overlay m-5" on:click={closeModal}>
		<!-- svelte-ignore a11y-click-events-have-key-events -->
		<div class="modal max-w-2xl" on:click={(e) => e.stopPropagation()}>
			<slot />
		</div>
	</div>
{/if}

<style>
	.modal-overlay {
		position: fixed;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		background-color: rgba(0, 0, 0, 0.5);
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.modal {
		background-color: white;
		border-radius: 0.5rem;
		box-shadow: 0 0 10px rgba(0, 0, 0, 0.5);
		display: flex;
		flex-direction: column;
		align-items: center;
	}
	
</style>
