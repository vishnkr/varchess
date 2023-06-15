<script lang="ts">
	// @ts-nocheck
	import Switch from 'svelte-switch';
	let checkedValue = false;
	function handleChange(e) {
		const { checked } = e.detail;
		checkedValue = checked;
	}
	let selectedGoal = 'Checkmate';
	let selectedGoalIdx = 0;
	const goals = [
		{ name: 'Checkmate', desc: 'Check the king and attack all of its escape squares' },
		{ name: 'Capture', desc: "Capture all of your opponent's pieces on the board" },
		{ name: 'n-Check', desc: 'Check the opponent king n-times' },
		{
			name: 'Target Square',
			desc: 'Move one of your selected pieces to a target square on the board'
		}
	];

	function selectGoal(selected: string, idx: number) {
		selectedGoal = selected;
		selectedGoalIdx = idx;
	}
</script>

<div class="flex flex-col p-3">
	<!-- svelte-ignore a11y-label-has-associated-control -->
	<label class="relative inline-flex items-center cursor-pointer">
		<span class="m-3 text-md font-medium text-gray-900 dark:text-gray-300">Show Legal Moves</span>
		<Switch on:change={handleChange} checked={checkedValue} />
	</label>
	<br />
	<h2 class="font-bold md:text-xl sm:text-lg">Select Objective:</h2>
	<div class="flex flex-col">
		{#each goals as goal, index}
			<!-- svelte-ignore a11y-click-events-have-key-events -->
			<div
				class="p-1 m-3 flex lg:flex-row flex-col items-center space-x-3 cursor-pointer
            {selectedGoalIdx === index
					? 'border-indigo-900 border-solid border-2 rounded'
					: 'border-gray-300'}"
				on:click={() => selectGoal(goal, index)}
			>
				<input
					type="radio"
					class="form-radio text-indigo-600 h-4 w-4"
					value={goal.name}
					checked={selectedGoalIdx === index}
				/>
				<h5 class="font-bold md:text-lg text-sm">{goal.name}</h5>
				<p>{goal.desc}</p>
			</div>
		{/each}
	</div>
</div>
