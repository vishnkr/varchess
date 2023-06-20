<script lang="ts">

    let ruleSelections = {
        variantType: "Standard", 
        objective: "Checkmate"
    }
    const variantTypes = [
        { name: "Standard", desc:"Pieces move like regular chess"},
        { name: "Duck Chess", desc: "Move a duck along with a piece"},
        { name: "Wormhole", desc: "Chess with teleportation"},
        { name: "Archer Chess", desc: "Chess with ranged attacks"}
    ]
	const goals = [
		{ name: 'Checkmate', desc: 'Check the king and attack all of its escape squares to win' },
		{ name: 'Antichess', desc: "Sacrifice all of your pieces on the board to win" },
		{ name: 'n-Check', desc: 'Check the opponent king n-times to win' },
		{ name: 'Target Square', desc: 'Move one of your selected pieces to a target square on the board to win'}
	];
    export let loggedIn=false;
</script>

<div>
        <div class="flex items-center justify-start m-2 space-x-2">
            <h3 class="font-bold md:text-md sm:text-lg">Select Variant Type</h3>
            <div class="relative inline-block">
                <select 
                    class="appearance-none border rounded-md py-2 px-4 pr-8 leading-tight focus:outline-none 
                    cursor-pointer focus:ring focus:border-blue-500" 
                    bind:value={ruleSelections.variantType}>
                    {#each variantTypes as variant, index (variant.name)}
                    <option value={variant.name}>{variant.name}</option>
                    {/each}
                </select>
                <div class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2 text-gray-700">
                    <svg
                    class="fill-current h-4 w-4"
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 20 20"
                    >
                    <path
                        fill-rule="evenodd"
                        d="M6.293 8.707a1 1 0 010-1.414l3-3a1 1 0 10-1.414-1.414L5 5.586 2.707 3.293a1 1 0 00-1.414 1.414l3 3a1 1 0 001.414 0z"
                        clip-rule="evenodd"
                    />
                    </svg>
                </div>
            </div>
            <button class="p-2 bg-orange-400 rounded-md text-white">View Variant Rules</button>
        </div>
        <div class=" m-2">
            <h3 class="font-bold text-start text-md md:text-lg">Select Objective:</h3>
            <div class="relative flex flex-col" id="goal">
                {#each goals as goal, index}
                    <!-- svelte-ignore a11y-click-events-have-key-events -->
                    <div
                        class="p-1 m-1 flex lg:flex-row flex-col items-center space-x-3 cursor-pointer
                    {goal.name === ruleSelections.objective
                            ? 'border-indigo-900 border-solid border-2 rounded'
                            : 'border-gray-300'}"
                        on:click={() => {ruleSelections.objective = goal.name}}
                    >
                        <input
                            type="radio"
                            class="form-radio text-indigo-600 h-4 w-4"
                            value={goal.name}
                            checked={ruleSelections.objective === goal.name}
                        />
                        <h5 class="font-semibold text-md">{goal.name}</h5>
                        <h6 class="text-sm">{goal.desc}</h6>
                    </div>
                {/each}
                {#if !loggedIn}
                <div class="absolute inset-0 flex flex-col bg-black opacity-70 rounded items-center justify-center">
                    <!-- Content for the overlay div -->
                    <i class="fa-solid fa-lock" style="color: #ffffff;"></i>
                    <p class="text-white md:text-2xl text-lg px-2">
                        Login to modify objective 
                    </p>
                </div>
                {/if}
            </div>
        </div>
        
</div>
