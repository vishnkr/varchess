<script lang="ts">
	import { login } from '$lib/api/auth';
	import { logout, authStore } from '$lib/store/auth';
	import { goto } from '$app/navigation';
	import Github from '$lib/icons/Github.svelte';
	import Google from '$lib/icons/Google.svelte';
	import Input from '$lib/components/ui/input/input.svelte';
	import Button from '$lib/components/ui/button/button.svelte';

	let username = '';
	let password = '';
	let auth;

	$: authStore.subscribe((state) => (auth = state));

	async function handleLogin(event: Event) {
		event.preventDefault();

		const success = await login(username, password);
		console.log('suc',success)
		if (!success) {
			alert('Login failed');
			return;
		}

		await goto('/home');
	}

	let isModalOpen = false;
	const openModal = () => (isModalOpen = true);
	const closeModal = () => (isModalOpen = false);
</script>

<div class="flex flex-col items-center h-full w-full">
	<h2 class="mt-2 text-center text-3xl text-white font-bold tracking-tight">
		Sign in to your account
	</h2>
	<form class="flex flex-col items-center space-y-2 w-full pt-4" on:submit={handleLogin}>
		<div class="form-control w-full max-w-md">
			<label class="label font-medium pb-1">
				<span class="label-text text-white">Username</span>
			</label>
			<Input bind:value={username} type="username" placeholder="username" class="w-full text-black" required/>
		</div>
		<div class="form-control w-full max-w-md">
			<label class="label font-medium pb-1">
				<span class="label-text text-white">Password</span>
			</label>
			<Input type="password" bind:value={password} placeholder="password" class="w-full text-black" required />
		</div>
		<Button variant="outline" type="submit" value="Login">Submit</Button>	
		<div>
			<p on:click={openModal} class="cursor-pointer text-blue-500">
				Don't have an account? Create one
			</p>
		</div>
	</form>

	{#if isModalOpen}
		<div class="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50">
			<div class="bg-white p-4 rounded-md">
				<h2 class="text-xl font-bold mb-4 text-black">Sign up</h2>
				<form method="POST" action="?/signup">
					<label class="label font-medium pb-1">
						<span class="label-text text-black">Email</span>
					</label>
					<input type="text" name="signup-email" class="text-white input input-bordered w-full mb-2" />
					<label class="label font-medium pb-1 text-white">
						<span class="label-text text-black">Username</span>
					</label>
					<input type="text" name="signup-username" class="text-white input input-bordered w-full mb-2" />
					<label class="label font-medium pb-1 text-white">
						<span class="label-text text-black">Password</span>
					</label>
					<input type="password" name="signup-password" class="text-white input input-bordered w-full mb-4" />
					<div class="flex justify-center items-center space-x-2">
						<input type="submit" value="Submit" class="bg-orange-600 text-white px-3 py-2 rounded-md cursor-pointer" />
						<button on:click={closeModal} class="text-blue-500 underline cursor-pointer">Close</button>
					</div>
				</form>
			</div>
		</div>
	{/if}

	<div class="w-full max-w-md pt-4 flex flex-col space-y-4">
		<a href="/login/github">
			<button class="flex items-center justify-center p-2 w-full bg-white text-black rounded-full">
				<Github class="mr-2 h-6 w-6" />
				<span>Continue with GitHub</span>
			</button>
		</a>
		<a href="/login/google">
			<button class="flex items-center justify-center p-2 w-full bg-white text-black rounded-full">
				<Google class="mr-2 h-6 w-6" />
				<span>Continue with Google</span>
			</button>
		</a>
	</div>
</div>
