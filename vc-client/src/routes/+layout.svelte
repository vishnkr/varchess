<script>
	import Footer from '$lib/Footer.svelte';
	import Navbar from '$lib/Navbar.svelte';
	import '../app.css';
</script>

<div class="flex flex-col min-h-screen">
	<Navbar />
	<div class="font-inter bg-[#0a0c13] dark:bg-[#ebebe8] flex flex-grow">
		<main class="flex flex-grow">
			<slot />
		</main>
	</div>
	<Footer />
</div>
