<script lang="ts">
	import { Color } from '$lib/board/types';
	import { editorSettings } from '../../board/stores';
	const standardPieces = [
		{ name: 'Pawn', class: 'p' },
		{ name: 'King', class: 'k'},
		{ name: 'Queen', class: 'q' },
		{ name: 'Bishop', class: 'b'},
		{ name: 'Knight', class: 'n'},
		{ name: 'Rook', class: 'r'}
	];
	const customPieces=[
		{name: 'Dolphin', class:'d'},
		{name: 'Ninja', class:'i'},
		{name: 'Unicorn', class:'u'},
		{name: 'Tower', class:'a'},
		{name: 'Giraffe', class:'g'},
		{name: 'Juicer', class:'j'},
		{name: 'Astronaut', class:'s'},
		{name: 'Phage', class:'v'},
		{name: 'Zebra', class:'z'},
	]
	let color: Color = Color.WHITE;
	let selectedPiece = {class:'p',group:'standard'};
	export let loggedIn=false;
	const selectPiece = (pieceClass: string,group:string) => {
		selectedPiece = {...selectedPiece,class:pieceClass,group};
		editorSettings.update((val) => ({
			...val,
			pieceSelection: {
				pieceType: pieceClass,
				color
			}
		}));
	};
    const updateColor = (color:Color)=>{
        color = color;
        editorSettings.update((val) => ({
			...val,
			pieceSelection: {
				pieceType: selectedPiece.class,
				color
			}
		}));
    }
</script>

<div>
	<div class="grid grid-rows">
		<div class="grid grid-rows">
			<div class="grid grid-cols-2">
            <!-- svelte-ignore a11y-click-events-have-key-events -->
				<div 
					class="flex items-center rounded-md p-4 m-1.5 bg-white border border-gray-200 dark:border-gray-700 cursor-pointer"
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
							
			</div>
			<div class="grid grid-cols-2">
				{#each standardPieces as piece}
				<!-- svelte-ignore a11y-click-events-have-key-events -->
				<div
					class="flex items-center rounded-md pl-4 m-1.5 bg-gray-300 border border-gray-200 dark:border-gray-700 cursor-pointer"
					on:click={() => selectPiece(piece.class,"standard")}
				>
					<input
						class="cursor-pointer text-blue-600 bg-gray-100 border-gray-300 focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
						type="radio"
						name="piece"
						value={piece.class}
						bind:group={selectedPiece.class}
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
			<div class="px-2 m-1.5 py-2">
				<h3> Custom Pieces</h3>
				<button disabled={selectedPiece.group !== "custom"} class="p-2 m-2 bg-orange-500 text-white text-md rounded-md disabled:bg-slate-600">Set Move Pattern</button>
				<div class="relative grid grid-cols-1">
					{#each customPieces as piece}
					<!-- svelte-ignore a11y-click-events-have-key-events -->
					<div
						class="flex flex-cols items-center rounded-md bg-gray-300 border border-gray-200 dark:border-gray-700 
						cursor-pointer my-1 
						"
						on:click={() => selectPiece(piece.class,"custom")}
					>
							<div class="w-2/3">
							<input
								class="text-blue-600 bg-gray-100 border-gray-300 focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
								type="radio"
								name="piece"
								value={piece.class}
								bind:group={selectedPiece.class}
							/>
							<label
								for={piece.class}
								class="w-full py-4 ml-2 text-md font-medium text-gray-900 dark:text-gray-300 cursor-pointer"
							>
								{piece.name}
							</label>
						</div>
						<img class="w-1/3 max-h-12" src={`/src/lib/assets/pieces/${color}/${piece.class}.svg`} alt="piece" />
					</div>
					{/each}
					{#if !loggedIn}
					<div class="absolute inset-0 flex flex-col bg-black opacity-70 rounded items-center justify-center">
						<!-- Content for the overlay div -->
						<i class="fa-solid fa-lock" style="color: #ffffff;"></i>
						<p class="text-white md:text-2xl text-lg px-2">
							Login to add custom pieces 
						</p>
					</div>
					{/if}
				</div>
			</div>
		</div>
	</div>
</div>