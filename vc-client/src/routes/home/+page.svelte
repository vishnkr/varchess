<script>
	import Card from '$lib/components/shared/Card.svelte';
</script>

<svelte:head>
	<title>Home - Varchess</title>
</svelte:head>
<div class="w-full grid grid-rows-2">
	<div class="w-full flex justify-around items-center gap-4 text-white">
		<div class="w-1/3 bg-cyan-700 flex items-center justify-center rounded h-[200px] p-5 m-6">
			<h1 class="text-2xl">Play against Computer</h1>
		</div>
		<div class="w-1/3  bg-cyan-700 flex items-center justify-center rounded h-[200px] p-5 m-6">
			<h1 class="text-2xl">New Game</h1>
		</div>
		<div class="w-1/3  bg-cyan-700 flex items-center justify-center rounded h-[200px] p-5 m-6">
			<h1 class="text-2xl">Create Template</h1>
		</div>
	</div>
	<div class="px-9">
		<Card>Saved templates go here</Card>
	</div>
</div>
