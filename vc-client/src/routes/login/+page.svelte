<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { login, signup, fetchOAuthProviders, oauthStartUrl, type OAuthProviders } from '$lib/api/auth';
	import { authStore } from '$lib/store/auth';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { toast } from '$lib/store/alert';
	import LoginBtn from '$lib/components/LoginBtn.svelte';
	import { Eye, EyeOff } from 'lucide-svelte';

	type Mode = 'login' | 'signup';

	let mode: Mode = 'login';
	let identifier = '';
	let password = '';
	let signupEmail = '';
	let signupUsername = '';
	let signupPassword = '';
	let showPassword = false;
	let submitting = false;
	let oauthLoading: 'github' | 'google' | null = null;
	let providers: OAuthProviders = { github: false, google: false };
	let providersLoaded = false;

	$: if ($authStore.accessToken && $page.url.pathname === '/login') {
		goto('/home');
	}

	onMount(async () => {
		const err = $page.url.searchParams.get('error');
		if (err) {
			toast.error(err);
			history.replaceState({}, '', '/login');
		}
		providers = await fetchOAuthProviders();
		providersLoaded = true;
	});

	function isValidEmail(email: string) {
		return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
	}

	async function handleLogin(event: Event) {
		event.preventDefault();
		if (submitting) return;
		submitting = true;
		try {
			const result = await login(identifier.trim(), password);
			if (!result.ok) {
				toast.error(result.error);
				return;
			}
			await goto('/home');
		} finally {
			submitting = false;
		}
	}

	async function handleSignup(event: Event) {
		event.preventDefault();
		if (submitting) return;
		if (!isValidEmail(signupEmail)) {
			toast.error('Please enter a valid email.');
			return;
		}
		if (signupUsername.trim().length < 3) {
			toast.error('Username must be at least 3 characters.');
			return;
		}
		if (signupPassword.length < 8) {
			toast.error('Password must be at least 8 characters.');
			return;
		}
		submitting = true;
		try {
			const result = await signup(signupEmail.trim(), signupUsername.trim(), signupPassword);
			if (!result.ok) {
				toast.error(result.error);
				return;
			}
			toast.success('Welcome to Varchess!');
			await goto('/home');
		} finally {
			submitting = false;
		}
	}

	function startOAuth(provider: 'github' | 'google') {
		oauthLoading = provider;
		window.location.href = oauthStartUrl(provider);
	}

	$: hasOAuth = providers.github || providers.google;
</script>

<svelte:head>
	<title>{mode === 'login' ? 'Sign in' : 'Create account'} — Varchess</title>
</svelte:head>

<div class="auth-shell min-h-[calc(100vh-8rem)] flex items-center justify-center px-4 py-8">
	<div class="auth-panel w-full max-w-md">
		<div class="mb-8 text-center">
			<h1 class="text-2xl sm:text-3xl font-bold tracking-tight text-gray-900 dark:text-white">
				{mode === 'login' ? 'Welcome back' : 'Create your account'}
			</h1>
			<p class="mt-2 text-sm text-gray-600 dark:text-gray-400">
				{mode === 'login'
					? 'Sign in to play variants, save templates, and pick up where you left off.'
					: 'One account for custom boards, games, and templates.'}
			</p>
		</div>

		{#if providersLoaded && hasOAuth}
			<div class="space-y-3">
				{#if providers.google}
					<LoginBtn
						provider="google"
						loading={oauthLoading === 'google'}
						disabled={!!oauthLoading || submitting}
						on:click={() => startOAuth('google')}
						class="!rounded-lg dark:bg-zinc-900 dark:text-white dark:border-zinc-700"
					/>
				{/if}
				{#if providers.github}
					<LoginBtn
						provider="github"
						loading={oauthLoading === 'github'}
						disabled={!!oauthLoading || submitting}
						on:click={() => startOAuth('github')}
						class="!rounded-lg dark:bg-zinc-900 dark:text-white dark:border-zinc-700"
					/>
				{/if}
			</div>

			<div class="relative my-6">
				<div class="absolute inset-0 flex items-center" aria-hidden="true">
					<div class="w-full border-t border-gray-300 dark:border-zinc-700" />
				</div>
				<div class="relative flex justify-center text-xs uppercase tracking-wider">
					<span class="bg-[var(--auth-panel-bg)] px-3 text-gray-500 dark:text-gray-400"
						>or continue with email</span
					>
				</div>
			</div>
		{/if}

		<div
			class="flex rounded-lg p-1 mb-5 bg-gray-100 dark:bg-zinc-900/80 border border-gray-200 dark:border-zinc-800"
			role="tablist"
		>
			<button
				type="button"
				role="tab"
				aria-selected={mode === 'login'}
				class="flex-1 rounded-md py-2 text-sm font-medium transition
					{mode === 'login'
					? 'bg-white dark:bg-zinc-800 text-gray-900 dark:text-white shadow-sm'
					: 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white'}"
				on:click={() => (mode = 'login')}
			>
				Sign in
			</button>
			<button
				type="button"
				role="tab"
				aria-selected={mode === 'signup'}
				class="flex-1 rounded-md py-2 text-sm font-medium transition
					{mode === 'signup'
					? 'bg-white dark:bg-zinc-800 text-gray-900 dark:text-white shadow-sm'
					: 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white'}"
				on:click={() => (mode = 'signup')}
			>
				Sign up
			</button>
		</div>

		{#if mode === 'login'}
			<form on:submit={handleLogin} class="space-y-4">
				<div class="space-y-1.5">
					<Label for="identifier">Email or username</Label>
					<Input
						id="identifier"
						bind:value={identifier}
						type="text"
						autocomplete="username"
						placeholder="you@email.com or username"
						required
						disabled={submitting}
					/>
				</div>
				<div class="space-y-1.5">
					<Label for="password">Password</Label>
					<div class="relative">
						<Input
							id="password"
							bind:value={password}
							type={showPassword ? 'text' : 'password'}
							autocomplete="current-password"
							placeholder="••••••••"
							required
							disabled={submitting}
							class="pr-10"
						/>
						<button
							type="button"
							class="absolute right-2 top-1/2 -translate-y-1/2 p-1.5 text-gray-500 hover:text-gray-800 dark:hover:text-gray-200"
							aria-label={showPassword ? 'Hide password' : 'Show password'}
							on:click={() => (showPassword = !showPassword)}
						>
							{#if showPassword}
								<EyeOff class="w-4 h-4" />
							{:else}
								<Eye class="w-4 h-4" />
							{/if}
						</button>
					</div>
				</div>
				<Button class="w-full h-11" type="submit" disabled={submitting}>
					{submitting ? 'Signing in…' : 'Sign in'}
				</Button>
			</form>
		{:else}
			<form on:submit={handleSignup} class="space-y-4">
				<div class="space-y-1.5">
					<Label for="signup-email">Email</Label>
					<Input
						id="signup-email"
						bind:value={signupEmail}
						type="email"
						autocomplete="email"
						placeholder="you@email.com"
						required
						disabled={submitting}
					/>
				</div>
				<div class="space-y-1.5">
					<Label for="signup-username">Username</Label>
					<Input
						id="signup-username"
						bind:value={signupUsername}
						type="text"
						autocomplete="username"
						placeholder="Choose a display name"
						required
						minlength={3}
						disabled={submitting}
					/>
				</div>
				<div class="space-y-1.5">
					<Label for="signup-password">Password</Label>
					<div class="relative">
						<Input
							id="signup-password"
							bind:value={signupPassword}
							type={showPassword ? 'text' : 'password'}
							autocomplete="new-password"
							placeholder="At least 8 characters"
							required
							minlength={8}
							disabled={submitting}
							class="pr-10"
						/>
						<button
							type="button"
							class="absolute right-2 top-1/2 -translate-y-1/2 p-1.5 text-gray-500 hover:text-gray-800 dark:hover:text-gray-200"
							aria-label={showPassword ? 'Hide password' : 'Show password'}
							on:click={() => (showPassword = !showPassword)}
						>
							{#if showPassword}
								<EyeOff class="w-4 h-4" />
							{:else}
								<Eye class="w-4 h-4" />
							{/if}
						</button>
					</div>
				</div>
				<Button class="w-full h-11" type="submit" disabled={submitting}>
					{submitting ? 'Creating account…' : 'Create account'}
				</Button>
			</form>
		{/if}

		<p class="mt-6 text-center text-xs text-gray-500 dark:text-gray-500">
			By continuing you agree to play fair and keep shared games civil.
		</p>
	</div>
</div>

<style>
	.auth-shell {
		--auth-panel-bg: rgba(255, 255, 255, 0.92);
		background:
			radial-gradient(ellipse 80% 60% at 20% 10%, rgba(56, 189, 248, 0.18), transparent 55%),
			radial-gradient(ellipse 70% 50% at 90% 80%, rgba(251, 146, 60, 0.14), transparent 50%);
	}
	:global(.dark) .auth-shell {
		--auth-panel-bg: rgba(12, 14, 20, 0.88);
		background:
			radial-gradient(ellipse 80% 60% at 15% 0%, rgba(251, 146, 60, 0.12), transparent 55%),
			radial-gradient(ellipse 60% 40% at 100% 100%, rgba(56, 189, 248, 0.08), transparent 45%);
	}
	.auth-panel {
		background: var(--auth-panel-bg);
		border: 1px solid rgba(0, 0, 0, 0.08);
		border-radius: 1rem;
		padding: 1.75rem 1.5rem 1.5rem;
		backdrop-filter: blur(12px);
		box-shadow: 0 18px 50px rgba(0, 0, 0, 0.08);
	}
	:global(.dark) .auth-panel {
		border-color: rgba(255, 255, 255, 0.08);
		box-shadow: 0 18px 50px rgba(0, 0, 0, 0.35);
	}
</style>
