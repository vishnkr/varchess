<script>
	import Footer from '$lib/Footer.svelte';
	import Navbar from '$lib/Navbar.svelte';
	import '../app.css';
</script>

<div class="flex flex-col min-h-screen">
	<Navbar />
	<div class="font-inter dark:bg-[#0a0c13] radial-bg flex flex-grow">
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
