<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog';
	import { Button } from '$lib/components/ui/button';

	export let open = false;
	export let title = 'Variant rules';
	export let slides: { heading: string; body: string; tip?: string }[] = [];

	let index = 0;

	$: if (open) index = 0;
	$: total = slides.length;
	$: slide = slides[index] ?? null;

	function prev() {
		index = (index - 1 + total) % total;
	}
	function next() {
		index = (index + 1) % total;
	}
</script>

<Dialog.Root
	{open}
	onOpenChange={(v) => {
		open = v;
	}}
>
	<Dialog.Content class="max-w-md dark:bg-darkbg bg-white text-gray-900 dark:text-white">
		<Dialog.Header>
			<Dialog.Title>{title}</Dialog.Title>
			{#if total > 1}
				<Dialog.Description class="text-gray-500 dark:text-gray-400">
					{index + 1} / {total}
				</Dialog.Description>
			{/if}
		</Dialog.Header>

		{#if slide}
			<div class="min-h-[10rem] py-2">
				<h3 class="text-base font-semibold mb-2">{slide.heading}</h3>
				<p class="text-sm text-gray-700 dark:text-gray-300 leading-relaxed">{slide.body}</p>
				{#if slide.tip}
					<p
						class="mt-3 text-xs rounded-md border border-emerald-500/40 bg-emerald-50 dark:bg-emerald-900/20 text-emerald-800 dark:text-emerald-200 px-3 py-2"
					>
						{slide.tip}
					</p>
				{/if}
			</div>
		{/if}

		<div class="flex items-center justify-between gap-2 mt-2">
			{#if total > 1}
				<div class="flex gap-2">
					<Button variant="outline" size="sm" on:click={prev} aria-label="Previous">
						<i class="fa-solid fa-chevron-left" />
					</Button>
					<Button variant="outline" size="sm" on:click={next} aria-label="Next">
						<i class="fa-solid fa-chevron-right" />
					</Button>
				</div>
				<div class="flex gap-1.5" aria-hidden="true">
					{#each slides as _, i}
						<span
							class="h-1.5 w-1.5 rounded-full {i === index
								? 'bg-emerald-500'
								: 'bg-gray-300 dark:bg-gray-600'}"
						/>
					{/each}
				</div>
			{:else}
				<div />
			{/if}
			<Button variant="secondary" size="sm" on:click={() => (open = false)}>Close</Button>
		</div>
	</Dialog.Content>
</Dialog.Root>
