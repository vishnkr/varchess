<script lang="ts">
	import { onDestroy } from 'svelte';
	import { ruleEditor, wormholeEditor } from '$lib/store/editor';
	import { VariantType } from '$lib/types';
	import WormholeSetup from './rules/WormholeSetup.svelte';
	import VariantInfoDialog from './rules/VariantInfoDialog.svelte';

	type VariantCard = {
		name: string;
		type: VariantType;
		desc: string;
		slides: { heading: string; body: string; tip?: string }[];
	};

	const variantTypes: VariantCard[] = [
		{
			name: 'Checkmate',
			type: VariantType.Checkmate,
			desc: 'You know the rules.',
			slides: [
				{
					heading: 'Goal',
					body: 'Checkmate the opponent’s king — it is under attack and has no legal escape.'
				},
				{
					heading: 'Standard play',
					body: 'All usual chess rules apply: checks, castling, en passant, and promotion.'
				}
			]
		},
		{
			name: 'Antichess',
			type: VariantType.Antichess,
			desc: 'Try to lose all your pieces faster than your opponent.',
			slides: [
				{
					heading: 'Goal',
					body: 'Lose all your pieces — or get stalemated. The player who cannot move (or has nothing left) wins.'
				},
				{
					heading: 'Captures',
					body: 'If you can capture, you must. King is a normal piece: no check, no checkmate.'
				},
				{
					heading: 'Scenario',
					body: 'Your queen can take a pawn — you have to take it, even if it walks into a fork next move.',
					tip: 'Force your opponent into capturing your last pieces.'
				}
			]
		},
		{
			name: 'n-Check',
			type: VariantType.NCheck,
			desc: 'Check the king n times. Because one check is never enough.',
			slides: [
				{
					heading: 'Goal',
					body: 'Deliver N checks against the enemy king. You can still win by normal checkmate.'
				},
				{
					heading: 'Counting',
					body: 'Each check you give increments your counter. Reach the target to win instantly.'
				},
				{
					heading: 'Scenario',
					body: 'At 2/3 checks, a quiet check that does not mate still ends the game if it hits the target.',
					tip: 'Set “Checks to win” below when this variant is selected.'
				}
			]
		},
		{
			name: 'Wormhole',
			type: VariantType.Wormhole,
			desc: 'Teleport across the board because walking is overrated.',
			slides: [
				{
					heading: 'Portals',
					body: 'Place linked wormhole pairs on the board. Landing on one portal instantly exits at its pair.'
				},
				{
					heading: 'Exit rule',
					body: 'The exit must be empty or hold an enemy piece (you capture them). If your own piece blocks the exit, that move is illegal.'
				},
				{
					heading: 'King & cooldown',
					body: 'Kings may teleport. After a pair is used, it cools down for one full turn (both players) before anyone can use it again.',
					tip: 'Mirror pairs across the board so neither color gets a private highway.'
				},
				{
					heading: 'Scenario',
					body: 'Your knight steps onto a portal on c3 and appears on f6, capturing a queen — then that pair goes dark for a turn.',
					tip: 'Add pairs in the setup panel before Play / Save.'
				}
			]
		},
		{
			name: 'Archer Chess',
			type: VariantType.ArcherChess,
			desc: 'Why move next to your enemies when you can just shoot them from across the board?',
			slides: [
				{
					heading: 'Ranged attacks',
					body: 'Archers can strike from a distance instead of only moving adjacent to enemies.'
				},
				{
					heading: 'Coming into focus',
					body: 'Full archer scenarios will expand here as the variant rules are finalized.'
				}
			]
		}
	];

	let infoOpen = false;
	let infoTitle = '';
	let infoSlides: VariantCard['slides'] = [];

	function setVariant(type: VariantType) {
		ruleEditor.updateVariantType(type);
		if (type !== VariantType.Wormhole) {
			wormholeEditor.stop();
		}
	}

	function openInfo(e: MouseEvent, variant: VariantCard) {
		e.stopPropagation();
		infoTitle = `${variant.name} rules`;
		infoSlides = variant.slides;
		infoOpen = true;
	}

	$: targetChecks = Number($ruleEditor.customData?.targetChecks) || 3;

	function onTargetChecksChange(e: Event) {
		const value = Number((e.currentTarget as HTMLInputElement).value) || 3;
		ruleEditor.updateCustomData({ targetChecks: value });
	}

	onDestroy(() => wormholeEditor.stop());
</script>

<div class="text-black dark:text-white">
	<div class="m-2">
		<h3 class="font-bold md:text-md sm:text-lg">Select Variant Type</h3>
		<div class="relative flex flex-col" id="goal">
			{#each variantTypes as variant}
				<!-- svelte-ignore a11y-click-events-have-key-events -->
				<!-- svelte-ignore a11y-no-static-element-interactions -->
				<div
					class="p-3 m-2 flex flex-col lg:flex-row items-center justify-center text-center lg:text-left space-y-2 lg:space-y-0 lg:space-x-3 cursor-pointer rounded-xl border-2 transition
		{variant.type === $ruleEditor.variantType
			? 'border-green-500 bg-green-100 dark:bg-green-800/20'
			: 'border-gray-300 dark:border-gray-600 hover:bg-gray-100 dark:hover:bg-gray-700/30'}"
					on:click={() => setVariant(variant.type)}
				>
					<div class="flex flex-col flex-1 min-w-0">
						<div class="flex items-center justify-center lg:justify-start gap-2">
							<h5 class="font-semibold text-md">{variant.name}</h5>
							<button
								type="button"
								class="inline-flex items-center justify-center h-6 w-6 rounded-full border border-gray-400 dark:border-gray-500 text-xs text-gray-600 dark:text-gray-300 hover:border-emerald-500 hover:text-emerald-600 dark:hover:text-emerald-400"
								aria-label="About {variant.name}"
								on:click={(e) => openInfo(e, variant)}
							>
								<i class="fa-solid fa-info" />
							</button>
						</div>
						<h6 class="text-sm text-gray-600 dark:text-gray-400">{variant.desc}</h6>
					</div>
				</div>
			{/each}
		</div>

		{#if $ruleEditor.variantType === VariantType.NCheck}
			<label class="flex items-center gap-2 m-2 text-sm">
				<span>Checks to win</span>
				<input
					type="number"
					min="1"
					max="10"
					class="w-16 rounded border border-gray-400 bg-transparent px-2 py-1"
					value={targetChecks}
					on:change={onTargetChecksChange}
				/>
			</label>
		{/if}

		{#if $ruleEditor.variantType === VariantType.Wormhole}
			<WormholeSetup />
		{/if}
	</div>
</div>

<VariantInfoDialog bind:open={infoOpen} title={infoTitle} slides={infoSlides} />
