<!-- src/lib/components/ThemeToggle.svelte -->
<script lang="ts">
	import { onMount } from 'svelte';
	import { theme } from '$lib/store/theme';

	onMount(() => {
		const stored = localStorage.getItem('theme') as 'light' | 'dark' | null;
		const initial = stored ?? (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light');
		theme.set(initial);
		updateDOM(initial);
	});

	function updateDOM(mode: 'light' | 'dark') {
		document.documentElement.classList.toggle('dark', mode === 'dark');
		localStorage.setItem('theme', mode);
	}

	function toggleTheme() {
		theme.update(current => {
			const next = current === 'dark' ? 'light' : 'dark';
			updateDOM(next);
			return next;
		});
	}
</script>
<button
	on:click={toggleTheme}
	aria-label="Toggle theme"
	class="p-2 rounded-full transition-colors"
	class:hover:bg-darkbg={$theme === 'light'}
	class:hover:text-white={$theme === 'light'}
	class:hover:bg-lightbg={$theme === 'dark'}
	class:hover:text-black={$theme === 'dark'}
>
	<svg
		xmlns="http://www.w3.org/2000/svg"
		class="w-6 h-6 transition-colors"
		fill="none"
		stroke="currentColor"
		viewBox="0 0 24 24"
	>
		{#if $theme === 'dark'}
			<circle cx="12" cy="12" r="5" stroke-width="2" />
			<path stroke-width="2" d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42" />
		{:else}
			<path stroke-width="2" d="M21 12.79A9 9 0 0112.21 3 7 7 0 0012 21a9 9 0 009-8.21z" />
		{/if}
	</svg>
</button>
