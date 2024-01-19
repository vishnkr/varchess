<script lang="ts">
	import { enhance } from "$app/forms";
    import type { PageData } from './$types';
	import Github from "$lib/icons/Github.svelte";
	import Google from "$lib/icons/Google.svelte";
	export let data: PageData;

	let isModalOpen = false;
	const openModal = () => {
		isModalOpen = true;
	};

	const closeModal = () => {
		isModalOpen = false;
	};
</script>

<div class="flex flex-col items-center h-full w-full">
	<h2 class="mt-2 text-center text-3xl text-white font-bold tracking-tight">
	  Sign in to your account
	</h2>
	<form  method="POST" action="?/loginUser" class="flex flex-col items-center space-y-2 w-full pt-4">
		<div class="form-control w-full max-w-md">
			<label for="username" class="label font-medium pb-1">
				<span class="label-text text-white">Username</span>
			</label>
			<input type="text" name="username" class="input input-bordered w-full max-w-md" />
		</div>
		<div class="form-control w-full max-w-md">
			<label for="password" class="label font-medium pb-1">
				<span class="label-text text-white">Password</span>
			</label>
			<input type="password" name="password" class="input input-bordered w-full max-w-md" />
		</div>
		<input type="submit" value="Login"class="bg-orange-600 text-white px-3 py-2 rounded-md cursor-pointer" />
		<div>
			<p on:click={openModal} class="cursor-pointer text-blue-500">Don't have an account, Create one</p>
		</div>
	</form>
	{#if isModalOpen}
		<div class="fixed inset-0 flex items-center justify-center  bg-black bg-opacity-50">
		<div class="bg-white p-4 rounded-md">
			<!-- Modal content goes here -->
			<h2 class="text-xl font-bold mb-4 text-black">Sign up</h2>
			<form method="POST" action="?/signup">
				<label for="signup-email" class="label font-medium pb-1">
				<span class="label-text text-black">Email</span>
				</label>
				<input type="text" name="signup-email" class="text-white input input-bordered w-full mb-2" />
				<label for="signup-username" class="label font-medium pb-1 text-white">
					<span class="label-text text-black">Username</span>
				</label>
				<input type="text" name="signup-username" class="text-white input input-bordered w-full mb-2" />
				<label for="signup-password" class="label font-medium pb-1 text-white">
				<span class="label-text text-black">Password</span>
				</label>
				<input type="password" name="signup-password" class="text-white input input-bordered w-full mb-4" />
				<div class="flex justify-center items-center space-x-2">
					<input type="submit" value="Submit" class="bg-orange-600 text-white px-3 py-2 rounded-md cursor-pointer" />
					<button on:click={closeModal} class="text-blue-500 underline cursor-pointer">
						Close
					</button>
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