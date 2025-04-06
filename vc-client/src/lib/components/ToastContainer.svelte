<script lang="ts">
	import { toast, type ToastMessage } from '$lib/store/alert';
	import { fly, fade } from 'svelte/transition';

	let toasts: ToastMessage[] = [];
	$: $toast, toasts = $toast;
</script>

<div class="fixed top-4 right-4 z-[1000] flex flex-col space-y-2 w-full max-w-sm">
	{#each toasts as t (t.id)}
		<div
			in:fly={{ x: -40, duration: 200 }}
			out:fade
			class="rounded-lg shadow-md px-4 py-3 text-sm font-medium text-black border-l-4 bg-white"
			class:border-green-500={t.type === 'success'}
			class:border-red-500={t.type === 'error'}
			class:border-yellow-400={t.type === 'warning'}
		>
			{t.message}
		</div>
	{/each}
</div>
