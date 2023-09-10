<script>
	import Alert from '$lib/components/Alert.svelte';
	import Footer from '$lib/components/Footer.svelte';
	import Navbar from '$lib/components/Navbar.svelte';
	import "../app.css"
	let loggedIn = false;
</script>

<div class="flex flex-col min-h-screen">
	<Navbar {loggedIn}/>
	<div class="font-inter dark:bg-[#0a0c13] radial-bg flex flex-col flex-grow">
		<Alert />
		<main class="flex flex-grow">
			<slot />
		</main>
	</div>
	<Footer />
</div>

<style>
	.radial-bg {
		background: radial-gradient(circle, hsl(240, 60%, 9%), hsl(210, 7%, 5%));
	}
</style>