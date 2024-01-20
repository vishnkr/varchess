<script lang="ts">
	let ruleSelections = {
		variantType: 'Standard',
		objective: 'Checkmate'
	};

	const variantTypes = [
		{ name: 'Checkmate', desc: 'Check the king and attack all of its escape squares to win' },
		{ name: 'Antichess', desc: 'Sacrifice all of your pieces on the board to win' },
		{ name: 'n-Check', desc: 'Check the opponent king n-times to win' },
		{
			name: 'Target Square',
			desc: 'Move one of your selected pieces to a target square on the board to win'
		},
		{ name: 'Duck Chess', desc: 'Move a duck along with a piece' },
		{ name: 'Wormhole', desc: 'Chess with teleportation' },
		{ name: 'Archer Chess', desc: 'Chess with ranged attacks' }
	];
	export let loggedIn = false;
</script>

<div>
	<div class="flex justify-center items-center justify-start m-2 space-x-2">
		<button class="p-2 bg-orange-400 rounded-md text-white">View Variant Rules</button>
	</div>
	<div class="m-2">
		<h3 class="font-bold md:text-md sm:text-lg">Select Variant Type</h3>
		<div class="relative flex flex-col" id="goal">
			{#each variantTypes as variant, index}
				<!-- svelte-ignore a11y-click-events-have-key-events -->
				<div
					class="p-1 m-1 flex lg:flex-row flex-col items-center space-x-3 cursor-pointer
                    {variant.name === ruleSelections.objective
						? 'border-indigo-900 border-solid border-2 rounded'
						: 'border-gray-300'}"
					on:click={() => {
						ruleSelections.objective = variant.name;
					}}
				>
					<input
						type="radio"
						class="form-radio text-indigo-600 h-4 w-4"
						value={variant.name}
						checked={ruleSelections.objective === variant.name}
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
