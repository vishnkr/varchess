<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { setAuth } from '$lib/store/auth';
	import { toast } from '$lib/store/alert';

	onMount(async () => {
		const hash = window.location.hash.startsWith('#')
			? window.location.hash.slice(1)
			: window.location.hash;
		const params = new URLSearchParams(hash);
		const accessToken = params.get('accessToken');
		const refreshToken = params.get('refreshToken');
		const uid = params.get('uid');
		const username = params.get('username');

		// Clear sensitive tokens from the URL bar immediately.
		history.replaceState({}, '', '/login/callback');

		if (!accessToken || !refreshToken || !uid || !username) {
			toast.error('Sign-in incomplete. Please try again.');
			await goto('/login');
			return;
		}

		setAuth(username, uid, accessToken, refreshToken);
		toast.success(`Signed in as ${username}`);
		await goto('/home');
	});
</script>

<svelte:head>
	<title>Signing in — Varchess</title>
</svelte:head>

<div class="min-h-[40vh] flex flex-col items-center justify-center gap-3 text-gray-700 dark:text-gray-200">
	<div
		class="h-8 w-8 rounded-full border-2 border-sky-500 border-t-transparent animate-spin"
		aria-hidden="true"
	/>
	<p class="text-sm font-medium">Finishing sign-in…</p>
</div>
