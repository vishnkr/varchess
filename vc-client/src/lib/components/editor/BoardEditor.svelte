<script lang="ts">
    
    import down from '$lib/assets/svg/down.svg';
    import up from '$lib/assets/svg/up.svg';
    import right from '$lib/assets/svg/right.svg';
    import left from '$lib/assets/svg/left.svg';
	import type { Dimensions } from '$lib/board/types';
	import { createEventDispatcher } from 'svelte';

    // @ts-ignore
    import Switch from "svelte-switch";

    let checkedValue = false;
    export let dimensions:Dimensions;
    const dispatch = createEventDispatcher();
    
    // @ts-ignore
    function handleChange(e){
        const {checked} = e.detail;
        checkedValue = checked
    }
</script>
<div>
    <div>
        <button class="
            bg-red-600 rounded-md p-2 text-white border-white 
            transform transition duration-200 hover:scale-105"
            on:click={()=>dispatch('clear')}
            >
            <span>Clear Board <i class="fa-solid fa-trash" style="color: #ffffff;"></i> </span>
            
        </button>
    </div>
    <div class="bg-white py-2 rounded-md shadow-md">
        <h3 class="text-xl font-semibold">Board Width : {dimensions.files} </h3>
        <input class="cursor-pointer" type="range" min={5} max={16} bind:value={dimensions.files}/>
    </div>
    <div class="bg-white py-2 rounded-md shadow-md">
        <h3 class="text-xl font-semibold">Board Height : {dimensions.ranks}</h3>
        <input class="cursor-pointer" type="range" min={5} max={16}  bind:value={dimensions.ranks}/>
    </div>
    <div class="flex flex-col justify-between items-center h-auto p-2">
        <h3 class="text-xl">Shift Board</h3>
        <button class="dbtn" on:click={()=>dispatch('shift','up')}><img class="svg" src={up} alt="Shift up"/></button>
        <div class="flex-1">
            <button class="dbtn" on:click={()=>dispatch('shift','left')}><img class="svg" src={left} alt="Shift left"/></button>
            <button class="dbtn" on:click={()=>dispatch('shift','right')}><img class="svg" src={right} alt="Shift right"/></button>
        </div>
        <button class="dbtn" on:click={()=>dispatch('shift','down')}><img class="svg" src={down} alt="Shift down"/></button>
    </div>
    <label class="relative inline-flex items-center cursor-pointer">
        <span class="m-3 text-md font-medium text-gray-900 dark:text-gray-300">Disable Squares</span>
        <Switch on:change={handleChange} checked={checkedValue} />
    </label>
</div>