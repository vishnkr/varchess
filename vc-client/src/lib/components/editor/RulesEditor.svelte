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
			desc: 'Check the king and attack all of its escape squares to win',
			ruleComponent: Checkmate
		},
		{ 	name: 'Antichess',
			type: VariantType.Antichess,
			desc: 'Sacrifice all of your pieces on the board to win',
			ruleComponent: Antichess
		},
		{ 	name: 'n-Check', 
			type: VariantType.NCheck, 
			desc: 'Check the opponent king n-times to win',
			ruleComponent: NCheck 
		},
		{
			name: 'GoalChess',
			type: VariantType.GoalChess,
			desc: 'Move one of your selected pieces to a target square on the board to win'
		},
		{ 	name: 'Duck Chess', 
			type: VariantType.DuckChess,
			desc: 'Move a duck along with a piece' 
		},
		{ 	name: 'Wormhole', 
			type: VariantType.Wormhole, 
			desc: 'Chess with teleportation' 
		},
		{ 	name: 'Archer Chess', 
			type: VariantType.ArcherChess,
			desc: 'Chess with ranged attacks' 
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
			{#each variantTypes as variant, index}
				<!-- svelte-ignore a11y-click-events-have-key-events -->
				<div
					class="p-1 m-1 flex lg:flex-row flex-col items-center space-x-3 cursor-pointer
                    {variant.type === $ruleEditor.variantType
						? 'border-indigo-900 border-solid border-2 rounded'
						: 'border-gray-300'}"
					on:click={()=>setVariant(variant.type)}
				>
					<input
						type="radio"
						class="form-radio text-indigo-600 h-4 w-4"
						value={variant.name}
						checked={$ruleEditor.variantType === variant.type}
					/>
					<h5 class="font-semibold text-md">{variant.name}</h5>
					<h6 class="text-sm">{variant.desc}</h6>
				</div>
			{/each}
			{#if loggedIn}
				<div
					class="absolute inset-0 flex flex-col bg-black opacity-70 rounded items-center justify-center"
				>
					<!-- Content for the overlay div -->
					<i class="fa-solid fa-lock" style="color: #ffffff;" />
					<p class="text-white md:text-2xl text-lg px-2">Login to modify objective</p>
				</div>
			{/if}
		</div>
	</div>
</div>
