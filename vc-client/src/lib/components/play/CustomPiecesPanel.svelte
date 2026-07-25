<script lang="ts">
	import { lastGameConfig } from '$lib/store/stores';
	import { openCustomPiecePattern } from '$lib/store/editor';
	import type { MovePattern, PieceProps } from '$lib/types';
	import { fromWirePieceProps } from '$lib/utils/pieceProps';
	import { customPieceName, isCustomNotation } from '$lib/utils/customPieces';
	import MovePatternPreview from './MovePatternPreview.svelte';
	import { browser } from '$app/environment';

	type PieceEntry = { notation: string; name: string; pattern: MovePattern };

	$: pieceProps = ($lastGameConfig?.pieceProps ?? {}) as Record<string, PieceProps>;
	$: patterns = fromWirePieceProps(pieceProps);
	$: pieces = Object.entries(patterns)
		.filter(([n]) => isCustomNotation(n))
		.map(
			([notation, pattern]): PieceEntry => ({
				notation,
				name: customPieceName(notation),
				pattern
			})
		)
		.sort((a, b) => a.name.localeCompare(b.name));

	let pinned: PieceEntry | null = null;
	let hovered: PieceEntry | null = null;
	let leaveTimer: ReturnType<typeof setTimeout> | null = null;

	$: canHover = browser && window.matchMedia('(hover: hover) and (pointer: fine)').matches;
	/** Floating preview only while hovering (and not pinned). */
	$: floatPreview = canHover && hovered && (!pinned || hovered.notation !== pinned.notation) ? hovered : null;
	/** Pinned stays in the sidebar section. */
	$: inlinePreview = pinned;

	function clearLeaveTimer() {
		if (leaveTimer) {
			clearTimeout(leaveTimer);
			leaveTimer = null;
		}
	}

	function selectPiece(p: PieceEntry) {
		if (pinned?.notation === p.notation) {
			clearPreview();
			return;
		}
		pinned = p;
		hovered = null;
	}

	function clearPreview() {
		clearLeaveTimer();
		pinned = null;
		hovered = null;
	}

	function onEnter(p: PieceEntry) {
		clearLeaveTimer();
		if (canHover && !pinned) hovered = p;
		else if (canHover && pinned && pinned.notation !== p.notation) hovered = p;
	}

	function scheduleHoverClear() {
		clearLeaveTimer();
		leaveTimer = setTimeout(() => {
			hovered = null;
			leaveTimer = null;
		}, 160);
	}

	$: requested = $openCustomPiecePattern;
	$: if (requested) {
		const match = pieces.find((p) => p.notation === requested.toLowerCase());
		if (match) {
			pinned = match;
			hovered = null;
		}
		queueMicrotask(() => {
			if ($openCustomPiecePattern === requested) openCustomPiecePattern.set(null);
		});
	}
</script>

{#if pieces.length > 0}
	<section
		class="shrink-0 overflow-hidden rounded-md border border-border bg-card text-card-foreground shadow-sm"
	>
		<header class="border-b border-border px-3 py-2">
			<h3 class="text-sm font-medium leading-none">Custom pieces</h3>
			<p class="mt-1.5 text-xs text-muted-foreground leading-snug">
				{#if canHover}
					Hover to preview · click to pin
				{:else}
					Tap a piece to pin its move pattern
				{/if}
			</p>
		</header>

		<div class="flex flex-wrap gap-2 p-3">
			{#each pieces as piece}
				<button
					type="button"
					class="flex min-w-[3.5rem] flex-col items-center gap-1 rounded-md border border-input bg-background px-2.5 py-2 text-sm transition-colors hover:bg-accent hover:text-accent-foreground
						{(inlinePreview ?? floatPreview)?.notation === piece.notation
						? 'border-primary bg-accent'
						: ''}"
					aria-pressed={pinned?.notation === piece.notation}
					on:click={() => selectPiece(piece)}
					on:mouseenter={() => onEnter(piece)}
					on:mouseleave={scheduleHoverClear}
					on:focus={() => onEnter(piece)}
					on:blur={scheduleHoverClear}
				>
					<img
						src={`/src/lib/assets/pieces/white/${piece.notation}.svg`}
						alt={piece.name}
						class="h-9 w-9 object-contain"
					/>
					<span class="text-[11px] font-medium leading-tight">{piece.name}</span>
				</button>
			{/each}
		</div>

		{#if inlinePreview}
			<div class="border-t border-border p-3 pt-2">
				<div class="mb-2 flex items-center justify-between gap-2">
					<div class="flex min-w-0 items-center gap-2 text-sm font-medium">
						<img
							src={`/src/lib/assets/pieces/white/${inlinePreview.notation}.svg`}
							alt=""
							class="h-6 w-6"
						/>
						<span class="truncate">{inlinePreview.name}</span>
						<span class="text-xs font-normal text-muted-foreground"
							>({inlinePreview.notation.toUpperCase()})</span
						>
					</div>
					<button
						type="button"
						class="text-xs text-muted-foreground underline-offset-2 hover:text-foreground hover:underline"
						on:click={clearPreview}
					>
						Hide
					</button>
				</div>
				{#key inlinePreview.notation}
					<MovePatternPreview
						notation={inlinePreview.notation}
						pattern={inlinePreview.pattern}
						boardId={`pattern-inline-${inlinePreview.notation}`}
					/>
				{/key}
			</div>
		{/if}
	</section>

	{#if floatPreview}
		<div
			class="pattern-float"
			role="dialog"
			aria-label="{floatPreview.name} move pattern"
			on:mouseenter={clearLeaveTimer}
			on:mouseleave={scheduleHoverClear}
		>
			<div
				class="rounded-md border border-border bg-popover p-4 text-popover-foreground shadow-md"
			>
				<div class="mb-2 flex items-center gap-2 text-base font-medium">
					<img
						src={`/src/lib/assets/pieces/white/${floatPreview.notation}.svg`}
						alt=""
						class="h-8 w-8"
					/>
					{floatPreview.name}
					<span class="text-sm font-normal text-muted-foreground"
						>({floatPreview.notation.toUpperCase()})</span
					>
				</div>
				{#key floatPreview.notation}
					<MovePatternPreview
						notation={floatPreview.notation}
						pattern={floatPreview.pattern}
						boardId={`pattern-float-${floatPreview.notation}`}
						large
					/>
				{/key}
			</div>
		</div>
	{/if}
{/if}

<style>
	.pattern-float {
		position: fixed;
		z-index: 60;
		left: max(0.75rem, calc((100vw - 24rem) / 2 - min(22vw, 210px)));
		top: 50%;
		transform: translateY(-50%);
		width: min(440px, calc(100vw - 22rem), 92vw);
		pointer-events: auto;
	}
	@media (max-width: 1023px) {
		.pattern-float {
			display: none;
		}
	}
</style>
