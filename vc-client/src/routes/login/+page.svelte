<script lang="ts">
	import { enhance } from '$app/forms'
</script>
<section class="w-full dark:bg-white flex justify-center items-center py-20">
	<div
		class="w-full max-w-sm p-4 bg-white border border-gray-200 rounded-lg shadow sm:p-6 md:p-8 dark:bg-gray-800 dark:border-gray-700"
	>
		<form class="space-y-6" method="POST" use:enhance>
			<h5 class="text-xl font-medium text-gray-900 dark:text-white">Log in to Varchess</h5>
			<div>
				<label for="email" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
					>Your email/username</label
				>
				<input
					type="email"
					name="email"
					id="email"
					class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-600 dark:border-gray-500 dark:placeholder-gray-400 dark:text-white"
					placeholder="email or username"
					required
				/>
			</div>
			<div>
				<label for="password" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
					>Your password</label
				>
				<input
					type="password"
					name="password"
					id="password"
					placeholder="••••••••"
					class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-600 dark:border-gray-500 dark:placeholder-gray-400 dark:text-white"
					required
				/>
			</div>
			<div class="flex items-start">
				<div class="flex items-start">
					<div class="flex items-center h-5">
						<input
							id="remember"
							type="checkbox"
							value=""
							class="w-4 h-4 border border-gray-300 rounded bg-gray-50 focus:ring-3 focus:ring-blue-300 dark:bg-gray-700 dark:border-gray-600 dark:focus:ring-blue-600 dark:ring-offset-gray-800 dark:focus:ring-offset-gray-800"
							
						/>
					</div>
					<label for="remember" class="ml-2 text-sm font-medium text-gray-900 dark:text-gray-300"
						>Remember me</label
					>
				</div>
				<a href="/" class="ml-auto text-sm text-blue-700 hover:underline dark:text-blue-500"
					>Forgot Password?</a
				>
			</div>
			<button
				type="submit"
				class="w-full text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
				>Login to your account</button
			>
			<div class="text-sm font-medium text-gray-500 dark:text-gray-300">
				Not registered? <a href="/" class="text-blue-700 hover:underline dark:text-blue-500"
					>Create account</a
				>
			</div>
		</form>
		<span class="text-md text-center"> OR </span>
		 <form class="auth-form" method="post" action="?/OAuth2">
			<div class="google-btn">
				<div class="google-icon-wrapper">
				  <img class="google-icon" src="https://upload.wikimedia.org/wikipedia/commons/5/53/Google_%22G%22_Logo.svg"/>
				</div>
				<p class="btn-text"><b>Sign in with Google</b></p>
			  </div>
	
		</form> 
	</div>
</section>

<style>
	.google-btn {
  width: 184px;
  height: 42px;
  background-color: #4285f4;
  border-radius: 2px;
  box-shadow: 0 3px 4px 0 rgba(0,0,0,.25);
	}
	.google-btn .google-icon-wrapper {
    position: absolute;
    margin-top: 1px;
    margin-left: 1px;
    width: 40px;
    height: 40px;
    border-radius: 2px;
    background-color: #fff;
  }
  .google-icon {
    position: absolute;
    margin-top: 11px;
    margin-left: 11px;
    width: 18px;
    height: 18px;
  }
  .btn-text {
    float: right;
    margin: 11px 11px 0 0;
    color: #fff;
    font-size: 14px;
    letter-spacing: 0.2px;
    font-family: "Roboto";
  }
  .google-btn:hover {
    box-shadow: 0 0 6px #4285f4;
	cursor: pointer;
  }
  .google-btn:active {
    background: #1669F2;
  }

</style>