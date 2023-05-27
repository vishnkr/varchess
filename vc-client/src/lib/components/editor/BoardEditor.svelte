<script lang="ts">
    
    import down from '$lib/assets/svg/down.svg';
    import up from '$lib/assets/svg/up.svg';
    import right from '$lib/assets/svg/right.svg';
    import left from '$lib/assets/svg/left.svg';
	import type { Dimensions } from '$lib/board/types';
	import { createEventDispatcher } from 'svelte';

    export let dimensions:Dimensions;
    const dispatch = createEventDispatcher();

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
        <input type="checkbox" value="" class="sr-only peer">
        <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600"></div>
        <span class="ml-3 text-md font-medium text-gray-900 dark:text-gray-300">Disable Square</span>
      </label>
</div>