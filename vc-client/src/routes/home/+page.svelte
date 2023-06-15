<script>
	import Card from '$lib/components/shared/Card.svelte';
</script>

<svelte:head>
	<title>Home - Varchess</title>
</svelte:head>
<div class="container grid grid-rows-2 gap-1">
	<div class="container grid grid-cols-2 gap-1">
		<div class="flex items-center justify-center">
			<Card>Play against Computer</Card>
		</div>
		<div class="flex items-center justify-center">
			<Card>Create Room</Card>
		</div>
	</div>
	<div class="">
		<Card>Saved templates go here</Card>
	</div>
</div>
