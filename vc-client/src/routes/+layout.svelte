<script lang="ts">
	import { onMount } from 'svelte';
	import '../app.css';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { type AuthState, authStore, logout } from '$lib/store/auth';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';

	import Footer from '$lib/components/Footer.svelte';
	import ToastContainer from '$lib/components/ToastContainer.svelte';
	import ThemeToggle from '$lib/components/ThemeToggle.svelte';

	import { theme } from '$lib/store/theme';
	import { get } from 'svelte/store';

	let auth: AuthState;
	$: authStore.subscribe((value) => (auth = value));

	const publicRoutes = ['/login', '/forgot-password', '/'];
	const authRoutes = ['/home', '/editor', '/games', '/templates', '/settings', '/profile'];

	onMount(() => {
		const stored = localStorage.getItem('theme') as 'dark' | 'light' | null;
		const initial =
			stored ?? (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light');
		theme.set(initial);
		document.documentElement.classList.toggle('dark', initial === 'dark');
	});
	$: {
		const currentPath = $page.url.pathname;
		if (!auth.accessToken && authRoutes.includes(currentPath)) {
			goto('/login');
		} else if (auth.accessToken && publicRoutes.includes(currentPath)) {
			goto('/home');
		}
	}

	let isSidebarOpen = false;
	let isDropdownOpen = false;

	function handleSidebarHover() {
		isSidebarOpen = true;
	}

	function handleSidebarLeave() {
		isSidebarOpen = false;
	}

	function toggleDropdown() {
		isDropdownOpen = !isDropdownOpen;
	}

	function handleLogout() {
		logout();
		goto('/login');
	}
	let currentTheme: 'light' | 'dark' = get(theme);
	theme.subscribe((val) => (currentTheme = val));
	let sidebarItems = [
		{ route: '/home', icon: 'fas fa-home', label: 'Home' },
		{ route: '/games', icon: 'fa-solid fa-chess', label: 'Games' },
		{ route: '/templates', icon: 'fa-solid fa-rectangle-list', label: 'Templates' },
		{ route: '/settings', icon: 'fas fa-cog', label: 'Settings' }
	];
</script>

<div class="flex flex-col min-h-screen">
	<nav class="sticky top-0 z-50 bg-lightbg2 dark:bg-darkbg">
		<div class="max-w-8xl mx-auto px-4">
			<div class="flex justify-between items-center h-16">
				<!-- Logo -->
				<div class="flex items-center">
					<div class="flex cursor-pointer" on:click={() => goto(auth.accessToken ? '/home' : '/')}>
						<img
							src={currentTheme === 'light' ? '/logo-dark.svg' : '/logo.svg'}
							alt="logo"
							class="md:w-40 ml-3 w-32"
						/>
					</div>
				</div>

				<div class="flex items-center gap-3">
					<ThemeToggle />
					{#if !auth.accessToken}
						<a
							href="/login"
							class="px-4 py-2 rounded transition-colors"
							class:hover:bg-darkbg={currentTheme === 'light'}
							class:hover:text-white={currentTheme === 'light'}
							class:hover:bg-lightbg={currentTheme === 'dark'}
							class:hover:text-black={currentTheme === 'dark'}
						>
							Log in
						</a>
					{:else}
						<DropdownMenu.Root>
							<DropdownMenu.Trigger>
								<i class="fas fa-user m-2" /> {auth.username}</DropdownMenu.Trigger
							>
							<DropdownMenu.Content>
								<DropdownMenu.Group>
									<DropdownMenu.Label>My Account</DropdownMenu.Label>
									<DropdownMenu.Separator />
									<DropdownMenu.Item class="cursor-pointer" on:click={() => goto('/profile')}
										>Profile</DropdownMenu.Item
									>
									<DropdownMenu.Item class="cursor-pointer" on:click={handleLogout}
										>Logout</DropdownMenu.Item
									>
								</DropdownMenu.Group>
							</DropdownMenu.Content>
						</DropdownMenu.Root>
					{/if}
				</div>
			</div>
		</div>
	</nav>

	<div
		class="flex-grow bg-lightbg dark:bg-[#0a0c13] bg-gradient-radial dark:from-[#0a0c13] dark:to-[#060709] from-[#f8fafc] to-[#e2e8f0]"
	>
		<ToastContainer />
		<main class="flex">
			{#if auth.accessToken && authRoutes.includes($page.url.pathname)}
			<div
			on:mouseenter={handleSidebarHover}
			on:mouseleave={handleSidebarLeave}
			class="fixed top-16 left-0 h-[calc(100vh-4rem)] 
				   w-16 hover:w-48 transition-width duration-300 ease-in-out
				   border-r text-gray-900 dark:text-white border-gray-300 dark:border-gray-600
				   bg-lightbg dark:bg-darkbg z-40"
		>
					<ul class="flex flex-col space-y-4 p-2">
						{#each sidebarItems as item}
							<!-- svelte-ignore a11y-click-events-have-key-events -->
							<li
								class="group flex items-center space-x-4 p-2 hover:bg-gray-200 dark:hover:bg-gray-700 cursor-pointer"
								on:click={() => goto(item.route)}
							>
								<i class={item.icon} />
								<span
									class="whitespace-nowrap transition-opacity duration-300"
									class:opacity-0={!isSidebarOpen}
									class:pointer-events-none={!isSidebarOpen}
								>
									{item.label}
								</span>
							</li>
						{/each}
					</ul>
				</div>
			{/if}

			<div class={`flex-grow p-4 transition-all duration-300 ${isSidebarOpen ? 'pl-48' : 'pl-16'}`}>
				<slot />
			</div>
			
		</main>
	</div>
	<Footer />
</div>

<style>
	:global(.transition-width) {
		transition-property: width;
	}
</style>
