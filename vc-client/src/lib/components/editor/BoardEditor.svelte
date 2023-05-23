<script lang="ts">
    let width:number = 8;
    let height:number = 8;
    import down from '$lib/assets/svg/down.svg';
    import up from '$lib/assets/svg/up.svg';
    import right from '$lib/assets/svg/right.svg';
    import left from '$lib/assets/svg/left.svg';
</script>
<div>
    <div class="bg-white py-2 rounded-md shadow-md">
        <h3 class="text-xl font-semibold">Board Width : {width} </h3>
        <input class="cursor-pointer" type="range" min={5} max={16} bind:value={width}/>
    </div>
    <div class="bg-white py-2 rounded-md shadow-md">
        <h3 class="text-xl font-semibold">Board Height : {height}</h3>
        <input class="cursor-pointer" type="range" min={5} max={16}  bind:value={height}/>
    </div>
    <div class="flex flex-col justify-between items-center h-auto p-2">
        <h3 class="text-xl">Shift Board</h3>
        <button class="dbtn"><img class="dbtn-img" src={up} alt="Shift up"/></button>
        <div class="flex-1">
            <button class="dbtn"><img class="dbtn-img" src={left} alt="Shift left"/></button>
            <button class="dbtn"><img class="dbtn-img" src={right} alt="Shift right"/></button>
        </div>
        <button class="dbtn"><img class="dbtn-img" src={down} alt="Shift down"/></button>
    </div>
</div>