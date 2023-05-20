<script>
	import Footer from '$lib/Footer.svelte'
	import ContentContainer from '$lib/components/ContentContainer.svelte';
	import SideBar from '$lib/components/SideBar.svelte';
	import TopNavigation from '$lib/components/TopNavigation.svelte';
    import '../../app.css'

	let sidebar_show = true;
</script>

<div class="flex">
	{#if sidebar_show}
	<SideBar bind:show={sidebar_show} />
	{/if}
	<main class="w-full">
		<div class='content-container'>
			<TopNavigation />
			
		</div>
		<slot />
	</main>
</div>
