
<script lang="ts">
	
    const standardPieces=[{name:'Pawn',class:'p'},{name:'King',class:'k'},{name:'Queen',class:'q'},{name:'Bishop',class:'b'},{name:'Knight',class:'n'},{name:'Rook',class:'r'}];
    let color:string='white';
    let selectedPiece = 'p';
</script>


<div>
    <div class="flex flex-col">
        <div class="grid grid-cols-2 grid-rows-3">
            {#each standardPieces as piece}
                <div class="flex items-center rounded-md pl-4 m-1.5 bg-gray-300 border border-gray-200 dark:border-gray-700 cursor-pointer" on:click={()=>selectedPiece = piece.class}>   
                    <input class="cursor-pointer w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600" type="radio" name="piece" value={piece.class} bind:group={selectedPiece}/>
                    <label for={piece.class} class="w-full py-4 ml-2 text-md font-medium text-gray-900 dark:text-gray-300 cursor-pointer">
                        {piece.name}
                    </label>
                    <img src={`/src/lib/assets/pieces/${color}/${piece.class}.svg`} alt="piece" />
                    
                </div>
            {/each}
        </div>
    </div>
    
</div>