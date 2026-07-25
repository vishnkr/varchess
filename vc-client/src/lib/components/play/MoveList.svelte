<script lang="ts">
	import type { Move } from '$lib/types';
	import { moveHistory, viewPly, seekBoardToPly, dimensions } from '$lib/store/stores';
	import { formatMoveNotation } from '$lib/utils/index';
	import { isCustomNotation, customPieceName } from '$lib/utils/customPieces';
	import { openCustomPiecePattern } from '$lib/store/editor';
	import { cn } from '$lib/utils.js';

	function labelFor(move: Move, index: number): string {
		const n = Math.floor(index / 2) + 1;
		const prefix = index % 2 === 0 ? `${n}.` : `${n}...`;
		const files = $dimensions?.files ?? 8;
		const ranks = $dimensions?.ranks ?? 8;
		const notation = formatMoveNotation(move, files, ranks);
		return `${prefix} ${notation}`;
	}

	function pieceNotation(move: Move): string | null {
		const p = (move.piece ?? '').toLowerCase();
		if (!p || !isCustomNotation(p)) return null;
		return p;
	}

	function pieceColorFolder(move: Move, index: number): 'white' | 'black' {
		const raw = move.piece ?? '';
		if (raw && raw === raw.toUpperCase()) return 'white';
		if (raw && raw === raw.toLowerCase()) return 'black';
		return index % 2 === 0 ? 'white' : 'black';
	}

	$: activePly = $viewPly === null ? $moveHistory.length : $viewPly;

	function onSelect(index: number) {
		seekBoardToPly(index + 1);
	}

	function onPieceIconClick(e: MouseEvent, notation: string) {
		e.stopPropagation();
		openCustomPiecePattern.set(notation);
	}
</script>

<section
	class="flex h-full min-h-0 flex-col overflow-hidden rounded-md border border-border bg-card text-card-foreground shadow-sm"
>
	<header class="shrink-0 border-b border-border px-3 py-2">
		<h3 class="text-sm font-medium leading-none">Moves</h3>
	</header>

	{#if $moveHistory.length === 0}
		<p class="px-3 py-3 text-sm text-muted-foreground">No moves yet</p>
	{:else}
		<ul class="min-h-0 flex-1 divide-y divide-border overflow-y-auto">
			{#each $moveHistory as move, index}
				{@const custom = pieceNotation(move)}
				<li>
					<button
						type="button"
						class={cn(
							'flex w-full items-center gap-2 px-3 py-1.5 text-left font-mono text-sm transition-colors',
							activePly === index + 1
								? 'bg-accent text-accent-foreground'
								: 'text-foreground hover:bg-muted/60'
						)}
						on:click={() => onSelect(index)}
					>
						{#if custom}
							<!-- svelte-ignore a11y-click-events-have-key-events -->
							<!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
							<img
								src={`/src/lib/assets/pieces/${pieceColorFolder(move, index)}/${custom}.svg`}
								alt={customPieceName(custom)}
								title="View {customPieceName(custom)} moves"
								class="h-5 w-5 shrink-0 cursor-pointer object-contain rounded-sm hover:ring-1 hover:ring-ring"
								on:click={(e) => onPieceIconClick(e, custom)}
							/>
						{/if}
						<span class="truncate">{labelFor(move, index)}</span>
					</button>
				</li>
			{/each}
		</ul>
	{/if}
</section>
