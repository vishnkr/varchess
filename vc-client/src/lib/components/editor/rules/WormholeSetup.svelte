<script lang="ts">
	import { boardEditor, ruleEditor, wormholeEditor } from '$lib/store/editor';
	import { Button } from '$lib/components/ui/button';
	import { smallIndexToAlgebraic } from '$lib/utils/index';

	$: pairs = (Array.isArray($ruleEditor.customData?.wormholePairs)
		? $ruleEditor.customData.wormholePairs
		: []) as number[][];
	$: ranks = $boardEditor.ranks;
	$: files = $boardEditor.files;
	$: pairing = $wormholeEditor.active;
	$: pending = $wormholeEditor.pending;

	function label(sq: number) {
		return smallIndexToAlgebraic(sq, files, ranks);
	}

	function usedSquares(exceptPair = -1): Set<number> {
		const used = new Set<number>();
		pairs.forEach((p, i) => {
			if (i === exceptPair || !Array.isArray(p) || p.length < 2) return;
			used.add(p[0]);
			used.add(p[1]);
		});
		return used;
	}

	function togglePairing() {
		if (pairing) {
			wormholeEditor.stop();
		} else {
			wormholeEditor.start();
		}
	}

	function removePair(i: number) {
		const next = pairs.filter((_, idx) => idx !== i);
		ruleEditor.updateCustomData({ wormholePairs: next });
		if (pairing) wormholeEditor.start();
	}

	function clearPairs() {
		ruleEditor.updateCustomData({ wormholePairs: [] });
		wormholeEditor.stop();
	}
</script>

<div
	class="m-2 rounded-lg border border-violet-400/50 bg-violet-50/80 dark:bg-violet-950/30 p-3 space-y-3 text-sm"
>
	<div>
		<h4 class="font-semibold text-violet-900 dark:text-violet-200">Wormhole setup (required)</h4>
		<p class="text-xs text-gray-600 dark:text-gray-400 mt-1">
			Place at least one portal pair on the board before saving or playing. Prefer mirrored pairs so
			both sides get a fair fight.
		</p>
	</div>

	<ul class="space-y-1.5 text-xs text-gray-700 dark:text-gray-300 list-disc pl-4">
		<li>Landing on a portal teleports you to its pair.</li>
		<li>Exit must be empty or hold an enemy (capture). Own piece on exit → illegal.</li>
		<li>Kings may use portals.</li>
		<li>After a teleport, that pair cools down for one full turn (both players).</li>
	</ul>

	<div class="flex flex-wrap gap-2">
		<Button
			type="button"
			size="sm"
			variant={pairing ? 'default' : 'outline'}
			on:click={togglePairing}
		>
			<i class="fa-solid fa-link mr-1.5" />
			{pairing ? 'Click two squares…' : 'Add pair on board'}
		</Button>
		{#if pairs.length > 0}
			<Button type="button" size="sm" variant="ghost" on:click={clearPairs}>Clear all</Button>
		{/if}
	</div>

	{#if pairing}
		<p class="text-xs text-violet-700 dark:text-violet-300">
			{#if pending == null}
				Select the first portal square on the board.
			{:else}
				First portal: <strong>{label(pending)}</strong> — now click its exit.
			{/if}
		</p>
	{/if}

	{#if pairs.length === 0}
		<p class="text-xs font-medium text-amber-700 dark:text-amber-300">
			No pairs yet — add at least one to unlock Play / Save for this variant.
		</p>
	{:else}
		<ul class="space-y-1">
			{#each pairs as pair, i}
				<li
					class="flex items-center justify-between gap-2 rounded border border-violet-300/40 dark:border-violet-700/50 px-2 py-1.5"
				>
					<span class="font-mono text-xs">
						Pair {i + 1}: {label(pair[0])} ↔ {label(pair[1])}
					</span>
					<button
						type="button"
						class="text-xs text-red-600 dark:text-red-400 hover:underline"
						on:click={() => removePair(i)}
					>
						Remove
					</button>
				</li>
			{/each}
		</ul>
		<p class="text-[11px] text-gray-500">
			{usedSquares().size} portal squares · cooldown fixed at 1 turn
		</p>
	{/if}
</div>
