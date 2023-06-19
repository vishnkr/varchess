<script>
	import { page } from '$app/stores';
</script>
  
<div class="w-full h-full text-white flex flex-col justify-center items-center">
	<h1 class="text-3xl block">404 Error: illegal move</h1>
	<p class="text-lg block">The requested page <code>"{$page.url.pathname.substr($page.url.pathname.lastIndexOf('/'))}"</code> does not exist. Try googling "en-passant"</p>
</div>
  

