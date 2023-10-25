<script lang="ts">
	import { enhance } from "$app/forms";
</script>

<div class="flex flex-col items-center h-full w-full">
	<h2 class="mt-2 text-center text-3xl font-bold  text-white tracking-tight">
		Register for an account
	</h2>
	<p class="text-center mt-1 text-white">
		Or <a href="/login" class="hover:text-blue-800 text-orange-600 font-medium hover:cursor-pointer hover:underline"
			>sign in</a
		> if you already have an account.
	</p>
	<form action="?/register" method="POST" use:enhance class="flex flex-col items-center space-y-2 w-full pt-4">
		<div class="form-control w-full max-w-md">
			<label for="username" class="label font-medium pb-1">
				<span class="label-text text-white">Username</span>
			</label>
			<input type="text" name="username" class="input input-bordered w-full max-w-md" />
		</div>
		<div class="form-control w-full max-w-md">
			<label for="email" class="label font-medium pb-1">
				<span class="label-text text-white">Email</span>
			</label>
			<input type="email" name="email" class="input input-bordered w-full max-w-md" />
		</div>
		<div class="form-control w-full max-w-md">
			<label for="password" class="label font-medium pb-1">
				<span class="label-text text-white">Password</span>
			</label>
			<input type="password" name="password" class="input input-bordered w-full max-w-md" />
		</div>
		<div class="form-control w-full max-w-md">
			<label for="passwordConfirm" class="label font-medium pb-1">
				<span class="label-text text-white">Confirm Password</span>
			</label>
			<input type="password" name="passwordConfirm" class="input input-bordered w-full max-w-md" />
		</div>
		<div class="w-full max-w-md pt-2">
			<button type="submit" class="p-3 text-white hover:bg-blue-800 bg-orange-600 w-full">Register</button>
		</div>
	</form>
</div>