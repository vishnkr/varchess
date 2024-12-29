<script lang="ts">
	import Alert from '$lib/components/Alert.svelte';
	import Footer from '$lib/components/Footer.svelte';
	import { onMount } from 'svelte';
	import '../app.css';
	import { goto } from '$app/navigation';
	export let data;

	let userId: string | undefined;
	let username: string | undefined;
	let isLoggedIn: boolean = false;
	$: {
		({ userId, username } = data);
		isLoggedIn = userId && username ? true : false;
	}
	let isOpen = false;

	function toggleDropdown() {
		isOpen = !isOpen;
	}

	function closeDropdown() {
		isOpen = false;
	}

	onMount(() => {
		const btn = document.querySelector('button.mobile-menu-btn');
		const menu = document.querySelector('.mobile-menu');
		btn?.addEventListener('click', () => {
			menu?.classList.toggle('hidden');
		});
	});

	let isSidebarOpen = false;

	function handleSidebarHover() {
		isSidebarOpen = true;
	}

	function handleSidebarLeave() {
		isSidebarOpen = false;
	}
</script>

<div class="flex flex-col min-h-screen">
	<nav class="sticky top-0 z-50 bg-[#0a0c13] dark:bg-[#1a1e25]">
		<div class="max-w-8xl mx-auto px-4">
			<div class="flex justify-between">
				<div class="flex space-x-4">
					<div>
						<a class=" flex" target="_self" href="/home">
							<img src="/logo.svg" alt="logo" class="md:w-40 ml-3 w-32" />
						</a>
					</div>
				</div>
				<div class="hidden md:flex items-center space-x-1" />
				{#if !userId}
					<div class="text-white md:flex items-center justify-end">
						<a
							href="/login"
							class="block py-2 px-4 md:text-lg text-sm hover:bg-gray-800 cursor-pointer">Log in</a
						>
					</div>
				{:else}
					<div class="text-white md:flex items-center justify-end relative">
						<!-- svelte-ignore a11y-click-events-have-key-events -->
						<div class="dropdown inline-block relative" on:click={toggleDropdown}>
							<button
								class="bg-transparent text-white py-2 px-4 md:text-lg text-sm hover:bg-gray-800 cursor-pointer"
							>
								{username} <i class="fa-solid fa-user" style="color: #feffff;" />
							</button>
							{#if isOpen}
								<ul class="dropdown-menu absolute pt-2 right-0">
									<li>
										<a
											href="/profile"
											class="block px-4 py-2 hover:bg-gray-800"
											on:click={closeDropdown}>Profile</a
										>
									</li>
									<li>
										<a
											href="/logout"
											class="block px-4 py-2 hover:bg-gray-800"
											on:click={closeDropdown}>Logout</a
										>
									</li>
								</ul>
							{/if}
						</div>
					</div>
				{/if}
			</div>
		</div>
	</nav>
	<div class="font-inter dark:bg-[#0a0c13] radial-bg flex flex-col flex-grow">
		<Alert />
		<main class="flex text-white">
			{#if isLoggedIn}
			<!-- Side Menu -->
			<div on:mouseenter={handleSidebarHover} on:mouseleave={handleSidebarLeave} class="flex-shrink-0 w-16 hover:w-48 transition-width duration-300 ease-in-out border-r border-gray-600 text-white h-screen">
					<!-- svelte-ignore a11y-click-events-have-key-events -->
					<ul class="flex flex-col space-y-4 p-2 cursor-pointer">
						<li class="group flex items-center space-x-4 p-2 hover:bg-gray-700" on:click={()=>goto("/home")}>
							<i class="fas fa-home"></i>
							<span class="whitespace-nowrap transition-opacity duration-300 ease-in-out" class:opacity-0={!isSidebarOpen}>Home</span>
								
						</li>
						<li class="group flex items-center space-x-4 p-2 hover:bg-gray-700" on:click={()=>goto("/games")}>
							<i class="fa-solid fa-chess"></i>
							<span class="whitespace-nowrap transition-opacity duration-300 ease-in-out" class:opacity-0={!isSidebarOpen}>Games</span>
						</li>
						<li class="group flex items-center space-x-4 p-2 hover:bg-gray-700" on:click={()=>goto("/templates")}>
							<i class="fa-solid fa-rectangle-list"></i>
							<span class="whitespace-nowrap transition-opacity duration-300 ease-in-out" class:opacity-0={!isSidebarOpen}>Templates</span>
						</li>
						<li class="group flex items-center space-x-4 p-2 hover:bg-gray-700" on:click={()=>goto("/settings")}>
							<i class="fas fa-cog"></i>
							<span class="whitespace-nowrap transition-opacity duration-300 ease-in-out" class:opacity-0={!isSidebarOpen}>Settings</span>
						</li>
						<li class="group flex items-center space-x-4 p-2 hover:bg-gray-700">
							<i class="fa-solid fa-robot"></i>
							<div class="group flex flex-col items-center hover:bg-gray-700">
								<span class="whitespace-nowrap transition-opacity duration-300 ease-in-out" class:opacity-0={!isSidebarOpen}>Stonkfish</span>
								<span class="whitespace-nowrap transition-opacity duration-300 ease-in-out rounded-md bg-red-500 px-1 text-sm text-white"  class:opacity-0={!isSidebarOpen}>Coming soon</span>
							</div>
							
						</li>
					</ul>
				</div>
			{/if}
			<!-- Main Content -->
			<div class="flex-grow">
				<slot />
			</div>
		</main>
	</div>
	<Footer />
</div>

<style>
	.radial-bg {
		background: radial-gradient(circle, hsl(240, 60%, 9%), hsl(210, 7%, 5%));
	}
</style>
