<script lang="ts">
	import { ruleEditor } from "$lib/store/editor";

	import Checkmate from "./rules/Checkmate.svelte";
	import Antichess from "./rules/Antichess.svelte";
	import NCheck from "./rules/NCheck.svelte";
	import { VariantType } from "$lib/types";
	import { Button } from "../ui/button";

	const variantTypes = [
		{ 	name: 'Checkmate', 
			type: VariantType.Checkmate, 
			desc: 'You know the rules.',
			ruleComponent: Checkmate
		},
		{ 	name: 'Antichess',
			type: VariantType.Antichess,
			desc: 'Try to lose all your pieces faster than your opponent.',
			ruleComponent: Antichess
		},
		{ 	name: 'n-Check', 
			type: VariantType.NCheck, 
			desc: "Check the king n times. Because one check is never enough.",
			ruleComponent: NCheck 
		},
		{
			name: 'GoalChess',
			type: VariantType.GoalChess,
			desc: "Forget kings, just race your piece to the goal square."
		},
		{ 	name: 'Duck Chess', 
			type: VariantType.DuckChess,
			desc: "Like classic chess, but there's an immortal duck blocking the board."
		},
		{ 	name: 'Wormhole', 
			type: VariantType.Wormhole, 
			desc: 'Teleport across the board because walking is overrated.' 
		},
		{ 	name: 'Archer Chess', 
			type: VariantType.ArcherChess,
			desc: 'Why move next to your enemies when you can just shoot them from across the board?' 
		}
	];

	export let loggedIn = false;
	function setVariant(type:VariantType){
		ruleEditor.updateVariantType(type);
	}

	function findRuleComponentByVariantType(type: VariantType) {
		const variant = variantTypes.find(variant => variant.type === type);
		return variant?.ruleComponent ?? undefined;
  	}

	function viewVariantRules(){
		ruleEditor.update((cur)=> { 
			return {
			...cur,
			isViewVariantRulesOn:true,
			ruleComponent: findRuleComponentByVariantType(cur.variantType) 
			}
		}
		)
	}
</script>

<div class="text-black dark:text-white">
	<div class="flex justify-center items-center m-2 space-x-2">
		<!-- svelte-ignore missing-declaration -->
		<Button on:click={viewVariantRules} class="p-2 bg-transparent border border-black text-black dark:border-white dark:text-white dark:hover:bg-blue-600/10 hover:bg-blue-600/10 rounded-md ">View Variant Rules</Button>
	</div>
	<div class="m-2">
		<h3 class="font-bold md:text-md sm:text-lg">Select Variant Type</h3>
		<div class="relative flex flex-col" id="goal">
			<!-- svelte-ignore a11y-click-events-have-key-events -->
			{#each variantTypes as variant, index}
	<!-- svelte-ignore a11y-click-events-have-key-events -->
	<!-- svelte-ignore a11y-no-static-element-interactions -->
	<div
		class="p-3 m-2 flex flex-col lg:flex-row items-center justify-center text-center lg:text-left space-y-2 lg:space-y-0 lg:space-x-3 cursor-pointer rounded-xl border-2 transition
		{variant.type === $ruleEditor.variantType
			? 'border-green-500 bg-green-100 dark:bg-green-800/20'
			: 'border-gray-300 dark:border-gray-600 hover:bg-gray-100 dark:hover:bg-gray-700/30'}"
		on:click={() => setVariant(variant.type)}
	>
		<div class="flex flex-col flex-1">
			<h5 class="font-semibold text-md">{variant.name}</h5>
			<h6 class="text-sm text-gray-600 dark:text-gray-400">{variant.desc}</h6>
		</div>
	</div>
{/each}

		</div>
	</div>
</div>
