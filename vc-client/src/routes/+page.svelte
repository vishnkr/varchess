<script lang="ts">
	import { enhance } from '$app/forms';
	import Button from '$lib/components/shared/Button.svelte'
	import Card from '$lib/components/shared/Card.svelte'
	import Modal from '$lib/components/shared/Modal.svelte';
	import { goto } from '$app/navigation';
	import Footer from '$lib/Footer.svelte';
	import Navbar from '$lib/Navbar.svelte';

	function routeToPage(route:string, replaceState:boolean) {
		goto(`/${route}`, { replaceState }) 
	}
	function handleGetStartedClick() {
		routeToPage('login', false);
	}
	let showModal = false;
	let inputText = '';
	const setShowModal = (val:boolean) => showModal = val
	function handleSubmit() {
 
		showModal = false;
	}
	//from-violet-700 via-rose-500 to-blue-600
</script>
<svelte:head> 
	<title> Varchess - Create and Play Custom Chess Variants</title>
</svelte:head>
<Navbar/>
<section class="font-inter bg-white dark:bg-[#1d2a35] text-white">
	<div class=" max-w-6xl mx-auto py-5 sm:py-24 px-4 sm:px-6 lg:px-8">
		<div class="isolate">
			<div class="relative px-6 lg:px-8 py-16 sm:py-8">
				<div class=" pb-8 ">
					<div class="flex space-x-8 items-center justify-start">
						<div class="text-left w-[38rem]">
							<h1 class="text-center font-bold tracking-tight text-zinc-900  dark:text-white text-6xl">
								Chess Variants 
								<span class="text-transparent bg-clip-text bg-gradient-to-r from-orange-500 to-blue-800 ">
									Redefined by You
								</span>
							</h1>
							<p class="mt-8 text-2xl text-center leading-8 text-zinc-700 dark:text-white">
								Create. Customize. Play
							</p>
							<div class="mt-10 flex items-center justify-center gap-x-6">
								<span class="rounded-md bg-orange-600 px-3.5 py-1.5 text-base font-semibold leading-7 
								text-white shadow-sm hover:bg-blue-800 focus-visible:outline focus-visible:outline-2 
								focus-visible:outline-offset-2 focus-visible:outline-blue-800 cursor-pointer"
								on:click={() => setShowModal(true)}
								>
									Try Quick Play
								</span>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
	<div class="relative w-full bg-white text-black">
		<div class="relative flex px-6 lg:px-8 py-8">
			<div class="">
				
			</div>
		</div>
	</div>
	<Modal isOpen={showModal} on:close={() => setShowModal(false)}>
		<div class="w-full">
			<div class="flex flex-col text-center px-3 py-3 bg-[#1d2a35] rounded-md">
				<h3 >Create a room now to play with friends!</h3>
            	<h5>NOTE: QuickPlay Variants are limited to 8x8 boards or smaller. Login to play on larger boards.</h5>
			</div>
			<div class="text-zinc-700 p-5 flex flex-col">
				<form method="POST" action="?/createRoom">
					<label >
						Enter Username
						<input type="text" name="username" />
					</label>
					
					
					<Button transition={false} 
					> Create Room 
					</Button>
				</form>
			</div>
		</div>
		
	</Modal>
</section>
<Footer />