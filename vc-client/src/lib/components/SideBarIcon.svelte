<script lang="ts">
    export let iconClass:string;
    export let text:string = 'tooltip 💡';
</script>

<div class="relative flex items-center justify-center h-12 w-12 mt-2 mb-2 mx-auto  
            bg-gray-400 hover:bg-red-600 dark:bg-gray-800 text-green-500 hover:text-white
            hover:rounded-xl rounded-3xl transition-all duration-300 ease-linear
            cursor-pointer shadow-lg group">
    <i class={`${iconClass} fa-lg`} />
    <span class="sidebar-tooltip group-hover:scale-100">
      {text}
    </span>
</div>