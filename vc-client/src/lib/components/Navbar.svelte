<script lang="ts">
	import { onMount } from 'svelte';
	let isLoggedIn = false;
	onMount(() => {
		const btn = document.querySelector('button.mobile-menu-btn');
		const menu = document.querySelector('.mobile-menu');
		btn?.addEventListener('click', () => {
			menu?.classList.toggle('hidden');
		});
	});
</script>

<nav class="sticky top-0 z-50 bg-[#0a0c13] dark:bg-[#1a1e25]">
		<div class="max-w-8xl mx-auto px-4">
			<div class="flex justify-between">
				<div class="flex space-x-4">
					<div>
						<a class=" flex" href="/">
							<img src="/logo.svg" alt="logo" class="md:w-40 ml-3 w-32" />
						</a>
					</div>
				</div>
				<div class="hidden md:flex items-center space-x-1" />
			</div>
		</div>
		<div class="mobile-menu text-white hidden">
			<a href="/login" class="block py-2 px-4 text-sm hover:bg-gray-800  cursor-pointer">Log in</a>
		</div>
	</nav>