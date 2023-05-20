<script>
    import { fly } from 'svelte/transition';
	import SideBarIcon from "./SideBarIcon.svelte";
    export let show = true;
</script>


<div class="h-screen w-16 flex flex-col
bg-white dark:bg-gray-900 shadow-lg text-white" transition:fly="{{ x: -200, duration: 150 }}">
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div class="text-black group" on:click={() => {show = false}}>
        <a href="/home" class="flex align-center">
            <img src="./favicon.png" alt="logo" class="h-12 w-12 mt-2 mb-2 mx-auto cursor-pointer shadow-lg text-center rounded-xl" />
        </a> 
    </div>
    <SideBarIcon iconClass="fa-solid fa-plus" text="Create New Room"/>
    <SideBarIcon iconClass="fa-solid fa-pencil" text="Create/ Edit Templates"/>
    <hr class="sidebar-hr" />
    <SideBarIcon iconClass="fa-solid fa-robot" text="Play against Stonkfish"/>
    <SideBarIcon iconClass="fa-solid fa-chess-board" text=""/>
    <SideBarIcon iconClass="fa-solid fa-info" text="Guide"/>
    <hr class="sidebar-hr" />
    
    <SideBarIcon iconClass="fa-solid fa-right-from-bracket" text="Log out"/>
</div>

