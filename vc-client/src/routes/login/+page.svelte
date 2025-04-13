<script lang="ts">
	import { login, signup } from '$lib/api/auth';
	import { authStore } from '$lib/store/auth';
	import { goto } from '$app/navigation';
	import { Button, buttonVariants } from '$lib/components/ui/button';
	import * as Dialog from "$lib/components/ui/dialog";
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { toast } from '$lib/store/alert';

	let username = '';
	let password = '';
	let signupEmail = '';
	let signupUsername = '';
	let signupPassword = '';
	let auth;
	$: authStore.subscribe((state) => (auth = state));

	let dialogOpen = false;

	function isValidEmail(email: string) {
		return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
	}

	async function handleLogin(event: Event) {
		event.preventDefault();
		try{
			const success = await login(username, password);
			if (!success) {
				toast.error('Login failed. Please check your credentials.');
				return;
			}
			await goto('/home');
		} catch (error){
			toast.error('Server is not reachable. Please try again later.');
		}
		
	}


	async function handleSignup(event: Event) {
		event.preventDefault();
		if (!isValidEmail(signupEmail)) {
			toast.error('Please enter a valid email.');
			return;
		}
		const success = await signup(signupEmail, signupUsername, signupPassword);
		if (!success) {
			toast.error('Signup failed. Please try again.');
			return;
		}
		toast.success('Account created! You can now log in.');
		dialogOpen = false;
	}
</script>

<svelte:head>
	<title>Log In to Varchess</title>
</svelte:head>
<div class="flex flex-col items-center w-full px-4">
	<h2 class="mt-8 text-3xl font-bold text-gray-800 dark:text-white text-center">Sign in to your account</h2>

	<form on:submit={handleLogin} class="w-full max-w-md space-y-4 pt-6">
		<div>
			<Label for="username" class="text-gray-800 dark:text-white">Username</Label>
			<Input id="username" bind:value={username} type="text" placeholder="Enter your username" required />
		</div>
		<div>
			<Label for="password" class="text-gray-800 dark:text-white">Password</Label>
			<Input id="password" bind:value={password} type="password" placeholder="••••••••" required />
		</div>
		<Button class="w-full" type="submit">Login</Button>
	</form>


	<Dialog.Root bind:open={dialogOpen}>
		<Dialog.Trigger class="mt-4 text-blue-400 text-sm hover:underline">
			Don't have an account? Create one
		</Dialog.Trigger>

		<Dialog.Content class="sm:max-w-[425px]">
			<Dialog.Header>
				<Dialog.Title>Create an account</Dialog.Title>
				<Dialog.Description>
					Enter your details below to create your account.
				</Dialog.Description>
			</Dialog.Header>

			<form method="POST" on:submit={handleSignup} class="grid gap-4 py-4">
				<div class="grid grid-cols-4 items-center gap-4">
					<Label for="signup-email" class="text-right">Email</Label>
					<Input id="signup-email" name="signup-email" type="email" bind:value={signupEmail} class="col-span-3" required />
				</div>
				<div class="grid grid-cols-4 items-center gap-4">
					<Label for="signup-username" class="text-right">Username</Label>
					<Input id="signup-username" name="signup-username"bind:value={signupUsername} type="text" class="col-span-3" required />
				</div>
				<div class="grid grid-cols-4 items-center gap-4">
					<Label for="signup-password" class="text-right">Password</Label>
					<Input id="signup-password" name="signup-password" bind:value={signupPassword} type="password" class="col-span-3" required />
				</div>
				<Dialog.Footer>
					<Button type="submit">Sign up</Button>
				</Dialog.Footer>
			</form>
		</Dialog.Content>
	</Dialog.Root>
</div>
