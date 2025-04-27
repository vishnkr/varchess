<script lang="ts">
	export let svg: any = null;  // Accept any Svelte component passed as a prop
	export let iconClass: string | null = null;
	export let title: string;
	export let isDarkMode: boolean = false;
  </script>
  
  <details class="bg-white dark:bg-darkbg2 shadow rounded group mb-4 w-full max-w-full">
	<summary class="list-none flat flex flex-wrap items-center cursor-pointer text-black dark:text-white">
	  <h3 class="flex flex-1 p-4 md:text-lg text-md font-semibold ">
		{title}
	  </h3>
  
	  {#if svg}
		<svelte:component 
		  this={svg}
		/>
	  {:else if iconClass}
		<i class={iconClass} />
	  {/if}
  
	  <div class="w-8 flex items-center justify-center">
		<div
		  class="border-8 border-transparent border-l-gray-600 ml-2
		  group-open:rotate-90 transition-transform origin-left"
		/>
	  </div>
	</summary>
	<div><slot /></div>
  </details>
  