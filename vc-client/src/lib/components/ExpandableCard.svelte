<script lang="ts">
	export let svg: string | null = null;
	export let iconClass: string | null = null;
	export let title: string;
</script>

<details class="bg-white shadow rounded group mb-4">
	<summary class="list-none flat flex flex-wrap items-center cursor-pointer">
		<h3 class="flex flex-1 p-4 lg:text-xl text-lg font-semibold">{title}</h3>
		{#if svg}
			<img class="w-g h-6" src={svg} alt="piece editor" />
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
