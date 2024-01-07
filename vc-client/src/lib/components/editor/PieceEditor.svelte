<script lang="ts">
	import { Color } from '$lib/board/types';
	import { editorSubTypeSelected, pieceEditor } from '$lib/store/editor';
	import TagInput from '../shared/TagInput.svelte';
	import { EditorSubType } from '../types';
	// @ts-ignores
	import Switch from 'svelte-switch';

	const standardPieces = [
		{ name: 'Pawn', class: 'p' },
		{ name: 'King', class: 'k' },
		{ name: 'Queen', class: 'q' },
		{ name: 'Bishop', class: 'b' },
		{ name: 'Knight', class: 'n' },
		{ name: 'Rook', class: 'r' }
	];
	const customPieces = [
		{ name: 'Dolphin', class: 'd' },
		{ name: 'Ninja', class: 'i' },
		{ name: 'Unicorn', class: 'u' },
		{ name: 'Tower', class: 'a' },
		{ name: 'Giraffe', class: 'g' },
		{ name: 'Juicer', class: 'j' },
		{ name: 'Astronaut', class: 's' },
		{ name: 'Phage', class: 'v' },
		{ name: 'Zebra', class: 'z' }
	];
	let color: Color = Color.WHITE;
	let selectedPiece = { class: 'p', group: 'standard' };
	let slideDirections = {"North":[-1,0],"East":[0,1],"South":[1,0],"West":[0,-1],"North East":[-1,1],"North West":[-1,-1],"South East":[1,1],"South West":[1,-1]}
	let setMovePattern = false;
	const toggleSetMP = () => {
		setMovePattern = !setMovePattern;
		editorSubTypeSelected.update((val) => setMovePattern ? EditorSubType.MovePattern : EditorSubType.Piece);
	};

	const cancel = () =>{
		pieceEditor.deletePiecePattern(selectedPiece.class)
		toggleSetMP()
	}

	const selectPiece = (pieceClass: string, group: string) => {
		selectedPiece = { ...selectedPiece, class: pieceClass, group };
		pieceEditor.update((val) => ({
			...val,
			pieceSelection: {
				pieceType: pieceClass,
				color,
				group
			}
		}));
	};
	const updateColor = (newColor: Color) => pieceEditor.updateColor(newColor);
</script>

<div>
	<div class="grid grid-rows">
		<div class="grid grid-rows">
			{#if !setMovePattern}
			<div class="grid grid-cols-2">
				<!-- svelte-ignore a11y-click-events-have-key-events -->
				<div
					class="flex items-center rounded-md p-4 m-1.5 bg-white border border-gray-200 dark:border-gray-700 cursor-pointer"
					on:click={() => updateColor(Color.WHITE)}
				>
					<input
						class="cursor-pointer w-4 h-4 text-black-600 bg-white border-gray-300 focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
						type="radio"
						value={Color.WHITE}
						bind:group={color}
					/>
					<label class="ml-2 cursor-pointer" for={color}>White</label>
				</div>
				<!-- svelte-ignore a11y-click-events-have-key-events -->
				<div
					class="flex items-center rounded-md pl-4 m-1.5 text-white bg-black border border-gray-200 dark:border-gray-700 cursor-pointer"
					on:click={() => updateColor(Color.BLACK)}
				>
					<input
						class="cursor-pointer w-4 h-4 text-white bg-black border-gray-300 focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
						type="radio"
						value={Color.BLACK}
						bind:group={color}
					/>
					<label class="ml-2 cursor-pointer" for={color}>Black</label>
				</div>
			</div>
			<div class="grid grid-cols-2">
				{#each standardPieces as piece}
					<!-- svelte-ignore a11y-click-events-have-key-events -->
					<div
						class="flex items-center rounded-md pl-4 m-1.5 bg-gray-300 border border-gray-200 dark:border-gray-700 cursor-pointer"
						on:click={() => selectPiece(piece.class, 'standard')}
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
				<h3>Custom Pieces</h3>
				<button
					disabled={selectedPiece.group !== 'custom'}
					class="p-2 m-2 bg-orange-500 text-white text-md rounded-md disabled:bg-slate-600"
					on:click={toggleSetMP} >Set Move Pattern</button
				>
				<div class="relative grid grid-cols-1">
					{#each customPieces as piece}
						<!-- svelte-ignore a11y-click-events-have-key-events -->
						<div
							class="flex flex-cols items-center rounded-md
							 bg-gray-300 border border-gray-200 dark:border-gray-700 cursor-pointer my-1"
							on:click={() => selectPiece(piece.class, 'custom')}
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
							<img
								class="w-1/3 max-h-12"
								src={`/src/lib/assets/pieces/${color}/${piece.class}.svg`}
								alt="piece"
							/>
						</div>
					{/each}
				</div>
			</div>
			{:else}
			<div class="px-2 m-1.5 py-2 flex flex-col">
				<h1 class="text-xl font-bold">Set Move Pattern</h1>
				<div class="flex items-center">
					<span class="w-4 h-4 inline-block bg-blue-600 rounded-sm"></span>
					<p class="text-lg font-semibold ml-2">Slide Pattern:</p>
				</div>
				<TagInput 
					{slideDirections} 
					dropDownText="Select Directions"
				/>
				<div class="flex items-center">
					<span class="w-4 h-4 inline-block bg-red-600 rounded-sm"></span>
					<p class="text-lg font-semibold ml-2">Jump Pattern:</p>
				</div>
				<!-- svelte-ignore a11y-label-has-associated-control -->
				<label class="relative inline-flex items-center cursor-pointer">
					<span class="m-3 text-md font-medium text-gray-900 dark:text-gray-300">Select Jump Moves</span>
					<!-- svelte-ignore missing-declaration -->
					<Switch checked={false} />
				</label>
				<button
					class="p-2 m-2 bg-green-500 text-white text-md rounded-md hover:bg-slate-400"
					on:click={toggleSetMP} >Save Pattern</button>
				<button
					class="p-2 m-2 bg-red-500 text-white text-md rounded-md hover:bg-slate-400"
					on:click={cancel}>Cancel</button>
			</div>
			
			{/if}
		</div>
	</div>	
</div>
