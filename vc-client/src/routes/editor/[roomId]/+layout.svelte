<script lang="ts">
	import Footer from "$lib/Footer.svelte";
import Navbar from "$lib/Navbar.svelte";
	import BoardEditor from "$lib/components/editor/BoardEditor.svelte";
	import SideBar from "$lib/components/SideBar.svelte";
    import Tabs from "$lib/components/shared/Tabs.svelte";

    let items = ['Custom','Predefined'];
    let activeItem = 'Custom';
    const tabChange=(e:CustomEvent<string>)=> activeItem = e.detail;
</script>
<Navbar />
<div class="font-inter editor-wrapper">
    <div class="left-panel overflow-y-scroll max-h-full">
        <div class="px-5 rounded-lg">
            
    <div class="border-b border-gray-200 dark:border-gray-700 flex flex-col text-center">
        <div class="flex justify-center">
            <Tabs {activeItem} {items} on:tabChange={tabChange}/>
        </div>
        {#if activeItem==="Custom"}
            <details class="bg-white shadow rounded group mb-4">
                <summary class="list-none flat flex flex-wrap items-center cursor-pointer">
                    <h3 class="flex flex-1 p-4 font-semibold">Board Editor</h3>
                <div class="w-8 flex items-center justify-center">
                    <div class="border-8 border-transparent border-l-gray-600 ml-2
                    group-open:rotate-90 transition-transform origin-left"></div>
                </div>
            </summary>
                <div><BoardEditor /></div>
            </details>
            <details class="bg-white shadow rounded group">
                <summary class="list-none flat flex flex-wrap items-center cursor-pointer">
                    <h3 class="flex flex-1 p-4 font-semibold">Piece Editor</h3>
                <div class="w-8 flex items-center justify-center">
                    <div class="border-8 border-transparent border-l-gray-600 ml-2
                    group-open:rotate-90 transition-transform origin-left"></div>
                </div>
            </summary>
            <div><p>sfdg</p></div>
            </details>
        {:else}
        <details class="bg-white shadow rounded group">
            <summary class="list-none flat flex flex-wrap items-center cursor-pointer">
                <h3 class="flex flex-1 p-4 font-semibold">Duck Chess</h3>
            <div class="w-8 flex items-center justify-center">
                <div class="border-8 border-transparent border-l-gray-600 ml-2
                group-open:rotate-90 transition-transform origin-left"></div>
            </div>
        </summary>
        <div><p>sfdg</p></div>
        </details>
        {/if}
    </div>

        </div>
    </div>
    <div class="content">Board Container</div>
    <div class="right-panel">Right Panel</div>
</div>

<Footer />