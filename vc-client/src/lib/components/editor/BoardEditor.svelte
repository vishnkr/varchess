<script lang="ts">
    let width:number = 8;
    let height:number = 8;
	import DPad from "../DPad.svelte";
    import Range from "../shared/Range.svelte";
</script>
<div>
    <div class="bg-white p-10 rounded-md shadow-md">
        <h3 class="text-xl font-semibold">Board Width : {width}</h3>
        <Range min={5} max={16} initialValue={width} on:change={(e) => width = e.detail.value}/>
    </div>
    <div class="bg-white p-10 rounded-md shadow-md">
        <h3 class="text-xl font-semibold">Board Height : {height}</h3>
        <Range min={5} max={16} initialValue={height} on:change={(e) => height = e.detail.value}/>
    </div>
    <DPad />
</div>

<style>

</style>