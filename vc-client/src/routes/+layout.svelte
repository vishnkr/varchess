<script lang="ts">
	import Alert from '$lib/components/Alert.svelte';
	import Footer from '$lib/components/Footer.svelte';
	import { onMount } from 'svelte';
	import '../app.css';

	export let data;

	let userId: string | undefined;
	let username: string | undefined;
	$: {
		({ userId, username } = data);
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
		<main class="flex">
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
