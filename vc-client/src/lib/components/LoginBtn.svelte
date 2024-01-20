<script lang="ts">
	import type { SvelteComponent } from 'svelte';

	import Google from '$lib/icons/Google.svelte';
	import Github from '$lib/icons/Github.svelte';
	import { type ClassValue, clsx } from 'clsx';
	import { twMerge } from 'tailwind-merge';
	import Pulse from './Pulse.svelte';

	export const Provider = ['google', 'github'] as const;
	type ProviderType = (typeof Provider)[number];
	export let provider: ProviderType;

	export let icon: typeof SvelteComponent | undefined | null = undefined;
	export let label: string | undefined | null = undefined;

	export let loading = false;
	export let disabled = false;

	export let withLoader = true;

	const busyClass = 'opacity-75 cursor-not-allowed';
	function cn(...inputs: ClassValue[]) {
		return twMerge(clsx(inputs));
	}
	$: mappedLabel = {
		google: 'Continue with Google',
		github: 'Continue with GitHub'
	}[provider];

	$: mappedIcon = {
		google: Google,
		github: Github
	}[provider];
</script>

<button
	type="button"
	{...$$restProps}
	class={cn(
		'inline-flex items-center bg-white justify-center w-full font-semibold tracking-wide h-12 rounded-full text-center text-gray-700 border border-gray-200 px-4 py-2 transition duration-30 text-sm sm:text-base hover:border-gray-300 focus:bg-gray-50 active:bg-gray-100 outline-none focus:outline-none focus:ring-2 focus:ring-gray-300 focus:ring-offset-2',
		loading || disabled ? busyClass : '',
		$$props.class
	)}
	disabled={disabled || loading}
	on:click
	on:change
	on:keydown
	on:keyup
	on:mouseenter
	on:mouseleave
>
	{#if icon !== null}
		{#if loading && withLoader}
			<div />
		{:else}
			<svelte:component this={icon ? icon : mappedIcon} class="w-5 h-5 mr-2" aria-hidden="true" />
		{/if}
	{/if}

	{#if loading && withLoader}
		<div class="h-5 w-5 flex items-center justify-center mr-2 overflow-visible">
			<Pulse class="w-2 h-2" aria-hidden="true" role="status" />
		</div>
	{/if}

	{#if label !== null}
		<span class="select-none">{label ? label : mappedLabel}</span>
	{/if}
</button>
