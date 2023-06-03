<script lang="ts">
	  import { onMount } from "svelte";
    let loggedIn = false;
    onMount(()=>{
      const btn = document.querySelector('button.mobile-menu-btn');
      const menu = document.querySelector('.mobile-menu');
      btn?.addEventListener("click",()=>{
        menu?.classList.toggle("hidden");
      })
    })
    
</script>

<nav class="sticky top-0 z-50 bg-[#0a0c13] dark:bg-[#1a1e25]">
  <div class="max-w-8xl mx-auto px-4">
    <div class="flex justify-between">
        <div class="flex space-x-4">
          <div>
            <a class=" flex" href="/">
              <img src="/logo.svg" alt="logo" class="w-40 ml-3" >
            </a>
          </div>
        </div>
        <div class="hidden md:flex items-center space-x-1"></div>
        <div class="hidden md:flex md:gap-x-9 md:flex-end items-center space-x-1">
          <div class="font-semibold px-3 py-5 cursor-pointer text-white">
            <a href="/guide" class="group transition duration-300 ease-in-out cursor-pointer">Guide
              <span class="block max-w-0 group-hover:max-w-full transition-all duration-500 h-0.5 bg-gradient-to-r from-orange-500 to-blue-800"/>
            </a>
          </div>
          <div class="font-semibold px-3 py-5  cursor-pointer text-white">
            <a href="/login" class="group transition duration-300 ease-in-out cursor-pointer">Log in <span aria-hidden="true">→</span>
              <span class="block max-w-0 group-hover:max-w-full transition-all duration-500 h-0.5 bg-gradient-to-r from-orange-500 to-blue-800"/>
            </a>
          </div>
        </div>
        <div class="md:hidden flex items-center text-white hover:text-black">
          <button class="mobile-menu-btn">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="#ffffff" class="w-6 h-6">
              <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5" />
            </svg> 
          </button>
                   
        </div>
    </div>
  </div>
  <div class="mobile-menu text-white hidden md:hidden">
    <a href="/guide" class="block py-3 px-4 text-sm hover:bg-gray-800 cursor-pointer">Guide</a>
    <a href="/login" class="block py-2 px-4 text-sm hover:bg-gray-800  cursor-pointer">Log in</a>
  </div>
</nav>




<style>
</style>