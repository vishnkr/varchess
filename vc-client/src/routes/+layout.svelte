<script>
	import Footer from '$lib/Footer.svelte'
	import Navbar from '$lib/Navbar.svelte'
    import '../app.css'

</script>

<Navbar/>
<div class="font-inter bg-[#0a0c13] dark:bg-[#ebebe8]">
	<main >
		<slot />
	</main>
</div>
<Footer />