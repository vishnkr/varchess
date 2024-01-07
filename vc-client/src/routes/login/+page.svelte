<script lang="ts">
	import { enhance } from "$app/forms";
	import { browser } from '$app/environment';
    import type { PageData } from './$types';
	import LoginBtn from "$lib/components/LoginBtn.svelte";
	import Github from "$lib/icons/Github.svelte";
	import Google from "$lib/icons/Google.svelte";
	export let data: PageData;
	/*function gotoGithubAuthProvider() {
        if (browser) {
            document.cookie = `state=${data?.githubAuthProviderState}`;
			document.cookie = `provider=github`;
        }
        window.location.href = data.githubAuthProviderRedirect || '';
    }
	function gotoGoogleAuthProvider() {
        if (browser) {
            document.cookie = `state=${data?.googleAuthProviderState}`;
			document.cookie = `provider=google`;
        }
        window.location.href = data.googleAuthProviderRedirect || '';
    }*/
	
</script>

<div class="flex flex-col items-center h-full w-full">
	<h2 class="mt-2 text-center text-3xl text-white font-bold tracking-tight">
	  Sign in to your account
	</h2>
  
	<div class="w-full max-w-md pt-4 flex flex-col space-y-4">
		<a href="/login/github">
		  <button class="flex items-center justify-center p-2 w-full bg-white rounded-full">
			<Github class="mr-2 h-6 w-6" />
			<span>Continue with GitHub</span>
		  </button>
		</a>
		<a href="/login/google">
		  <button class="flex items-center justify-center p-2 w-full bg-white rounded-full">
			<Google class="mr-2 h-6 w-6" />
			<span>Continue with Google</span>
		  </button>
		</a>
	  </div>
	  
  
	<!-- 	<form action="?/login" method="POST" use:enhance class="flex flex-col items-center space-y-2 w-full pt-4">
		<div class="form-control w-full max-w-md">
			<label for="username" class="label font-medium pb-1">
				<span class="label-text text-white">Email</span>
			</label>
			<input type="text" name="username" class="input input-bordered w-full max-w-md" />
		</div>
		<div class="form-control w-full max-w-md">
			<label for="password" class="label font-medium pb-1">
				<span class="label-text text-white">Password</span>
			</label>
			<input type="password" name="password" class="input input-bordered w-full max-w-md" />
		</div>
		<div class="w-full max-w-md">
			<a
				href="/reset-password"
				class="font-medium hover:text-blue-800 text-orange-600 hover:cursor-pointer hover:underline"
			>
				Forgot Password?</a
			>
		</div>

		<div class="w-full max-w-md pt-2">
			<button type="submit" class="p-3 text-white hover:bg-blue-800 bg-orange-600 w-full">Login</button>
		</div>
		
	</form>
	<div class="mt-4 flex items-center text-white">
        <div class="flex-grow border-t-2"></div>
        <div class="px-2">OR</div>
        <div class="flex-grow border-t-2"></div>
		<LoginBtn on:click={gotoGoogleAuthProvider} provider="google" label="Continue with Google"/>
    </div> -->
  
  </div>