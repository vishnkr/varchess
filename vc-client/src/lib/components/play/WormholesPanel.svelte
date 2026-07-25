<script lang="ts">
	import { dimensions, wormholePlayState } from '$lib/store/stores';
	import { wormholePairLabel } from '$lib/utils/wormhole';
	import { cn } from '$lib/utils.js';

	$: state = $wormholePlayState;
	$: files = $dimensions?.files ?? 8;
	$: ranks = $dimensions?.ranks ?? 8;
</script>

{#if state && state.pairs.length > 0}
	<section
		class="shrink-0 overflow-hidden rounded-md border border-border bg-card text-card-foreground shadow-sm"
	>
		<header class="border-b border-border px-3 py-2">
			<h3 class="text-sm font-medium leading-none">Wormholes</h3>
			<p class="mt-1.5 text-xs text-muted-foreground leading-snug">
				Move onto a portal to exit at its pair. Used pairs cool for one full turn.
			</p>
		</header>

		<ul class="divide-y divide-border">
			{#each state.pairs as pair, i}
				{@const closed = (state.cooldown[i] ?? 0) > 0}
				<li class="flex items-center gap-3 px-3 py-2">
					<span
						class="w-5 shrink-0 text-center text-xs tabular-nums text-muted-foreground"
						aria-hidden="true"
					>
						{i + 1}
					</span>
					<span class="min-w-0 flex-1 truncate font-mono text-xs tracking-tight">
						{wormholePairLabel(pair, files, ranks)}
					</span>
					<span
						class={cn(
							'inline-flex shrink-0 items-center rounded-sm border px-1.5 py-0.5 text-[10px] font-medium',
							closed
								? 'border-border bg-secondary text-secondary-foreground'
								: 'border-transparent bg-primary text-primary-foreground'
						)}
					>
						{closed ? `Cooling ${state.cooldown[i]}` : 'Open'}
					</span>
				</li>
			{/each}
		</ul>
	</section>
{/if}
