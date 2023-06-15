<script lang="ts">
	import { Color } from '$lib/board/types';
	import { editorSettings } from '../../board/stores';
	const standardPieces = [
		{ name: 'Pawn', class: 'p' },
		{ name: 'King', class: 'k' },
		{ name: 'Queen', class: 'q' },
		{ name: 'Bishop', class: 'b' },
		{ name: 'Knight', class: 'n' },
		{ name: 'Rook', class: 'r' }
	];
	let color: Color = Color.WHITE;
	let selectedPiece = 'p';

	const selectPiece = (piece: string) => {
		selectedPiece = piece;
		editorSettings.update((val) => ({
			...val,
			pieceSelection: {
				pieceType: selectedPiece,
				color
			}
		}));
	};
    const updateColor = (color:Color)=>{
        color = color;
        editorSettings.update((val) => ({
			...val,
			pieceSelection: {
				pieceType: selectedPiece,
				color
			}
		}));
    }
</script>

<div>
	<div class="flex flex-col">
		<div class="grid grid-cols-2 grid-rows-3">
            <!-- svelte-ignore a11y-click-events-have-key-events -->
            <div 
                class="flex items-center rounded-md pl-4 m-1.5 bg-white border border-gray-200 dark:border-gray-700 cursor-pointer"
			    on:click={()=>updateColor(Color.WHITE)}
            >
                <input class="cursor-pointer w-4 h-4 text-black-600 bg-white border-gray-300 focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600" 
                    type="radio"
                    value={Color.WHITE}
                    bind:group={color}
                />
                <label class="ml-2" for={color}>White</label>
            </div>
            <!-- svelte-ignore a11y-click-events-have-key-events -->
            <div 
                class="flex items-center rounded-md pl-4 m-1.5 text-white bg-black border border-gray-200 dark:border-gray-700 cursor-pointer"
			    on:click={()=>updateColor(Color.BLACK)}
            >
                <input class="cursor-pointer w-4 h-4 text-white bg-black border-gray-300 focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600" 
                    type="radio"
                    value={Color.BLACK}
                    bind:group={color}
                    />
                <label class="ml-2" for={color}>Black</label>
            </div>
			{#each standardPieces as piece}
				<!-- svelte-ignore a11y-click-events-have-key-events -->
				<div
					class="flex items-center rounded-md pl-4 m-1.5 bg-gray-300 border border-gray-200 dark:border-gray-700 cursor-pointer"
					on:click={() => selectPiece(piece.class)}
				>
					<input
						class="cursor-pointer w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
						type="radio"
						name="piece"
						value={piece.class}
						bind:group={selectedPiece}
					/>
					<label
						for={piece.class}
						class="w-full py-4 ml-2 text-md font-medium text-gray-900 dark:text-gray-300 cursor-pointer"
					>
						{piece.name}
					</label>
					<img src={`/src/lib/assets/pieces/${color}/${piece.class}.svg`} alt="piece" />
				</div>
			{/each}
		</div>
	</div>
</div>
