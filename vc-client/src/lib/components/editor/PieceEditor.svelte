<script lang="ts">
	import { Color } from '$lib/board/types';
	import { editorSubTypeSelected, pieceEditor } from '$lib/store/editor';
	import TagInput from '../shared/TagInput.svelte';
	import { EditorSubType } from '../types';
	import * as RadioGroup from '$lib/components/ui/radio-group/index.js';
	import { Label } from '$lib/components/ui/label/index.js';

	// @ts-ignores
	import Switch from 'svelte-switch';
	import type { PieceSelection } from '$lib/types';


	const standardPieces : {name:string,notation:string}[] = [
		{ name: 'Pawn', notation: 'p' },
		{ name: 'King', notation: 'k' },
		{ name: 'Queen', notation: 'q' },
		{ name: 'Bishop', notation: 'b' },
		{ name: 'Knight', notation: 'n' },
		{ name: 'Rook', notation: 'r' }
	];
	const customPieces : {name:string,notation:string}[] = [
		{ name: 'Dolphin', notation: 'd' },
		{ name: 'Ninja', notation: 'i' },
		{ name: 'Unicorn', notation: 'u' },
		{ name: 'Tower', notation: 'a' },
		{ name: 'Giraffe', notation: 'g' },
		{ name: 'Juicer', notation: 'j' },
		{ name: 'Astronaut', notation: 's' },
		{ name: 'Phage', notation: 'v' },
		{ name: 'Zebra', notation: 'z' }
	];
	let color: Color = Color.WHITE;
	let selectedPiece: PieceSelection = {piece:{ notation: 'p',pieceType: 'pawn',color:color},group:'standard'};
	let slideDirections = {
		North: [-1, 0],
		East: [0, 1],
		South: [1, 0],
		West: [0, -1],
		'North East': [-1, 1],
		'North West': [-1, -1],
		'South East': [1, 1],
		'South West': [1, -1]
	};
	let setMovePattern = false;
	const toggleSetMP = () => {
		setMovePattern = !setMovePattern;
		editorSubTypeSelected.update((val) =>
			setMovePattern ? EditorSubType.MovePattern : EditorSubType.Piece
		);
	};

	const cancel = () => {
		pieceEditor.deletePiecePattern(selectedPiece.piece.notation);
		toggleSetMP();
	};

	const selectPiece = (pieceType:string, notation:string, group:string) => {
		
		selectedPiece = { ...selectedPiece, 
			piece: {
				notation:notation, 
				pieceType:pieceType,
				color:color
			},
			group: group 
		};
		('selected piece',selectedPiece)
		pieceEditor.update((val) => ({
			...val,
			pieceSelection: {
				piece: {
					pieceType: selectedPiece.piece.pieceType,
					color,
					notation: selectedPiece.piece.notation
				},
				group: selectedPiece.group,
			}
		}));
	};
	const updateColor = (newColor: Color) => {
		pieceEditor.updateColor(newColor);
		color = newColor;
	};
</script>

<div>
	<div class="grid grid-rows">
		<div class="grid grid-rows">
			{#if !setMovePattern}
				<div class="grid grid-cols-2 gap-2 my-2">
					<!-- svelte-ignore a11y-click-events-have-key-events -->
					<!-- svelte-ignore a11y-no-static-element-interactions -->
					<div
						class={`flex items-center justify-center rounded-md mx-2 p-4 cursor-pointer border-2 transition-all duration-150 ${
							color === Color.WHITE ? 'border-green-500' : 'border-gray-800'
						} bg-white`}
						on:click={() => updateColor(Color.WHITE)}
					>
						<span class="text-black font-semibold">White</span>
					</div>

					<!-- svelte-ignore a11y-no-static-element-interactions -->
					<!-- svelte-ignore a11y-click-events-have-key-events -->
					<div
						class={`flex items-center justify-center rounded-md mx-2 p-4 cursor-pointer border-2 transition-all duration-150 ${
							color === Color.BLACK ? 'border-green-500' : 'border-gray-800'
						} bg-black`}
						on:click={() => updateColor(Color.BLACK)}
					>
						<span class="text-white font-semibold">Black</span>
					</div>
				</div>

				<div class="grid grid-cols-2">
					{#each standardPieces as piece}
						<!-- svelte-ignore a11y-click-events-have-key-events -->
						<!-- svelte-ignore a11y-no-static-element-interactions -->
						<div
							class="flex items-center rounded-md pl-4 m-1.5 bg-gray-300 border border-gray-200 dark:border-gray-700 cursor-pointer"
							on:click={() => selectPiece(piece.name?.toLowerCase(),piece.notation,"standard")}
						>
							<input
								class="cursor-pointer text-blue-600 bg-gray-100 border-gray-300 focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
								type="radio"
								name="piece"
								value={piece.notation}
								bind:group={selectedPiece.piece.notation}
							/>
							<label
								for={piece.notation}
								class="w-full py-4 ml-2 text-md font-medium text-gray-900 cursor-pointer"
							>
								{piece.name}
							</label>
							<img src={`/src/lib/assets/pieces/${color}/${piece.notation}.svg`} alt="piece" />
						</div>
					{/each}
				</div>
				<div class="px-2 m-1.5 py-2">
					<h3>Custom Pieces</h3>
					<button
						disabled={selectedPiece.group !== 'custom'}
						class="p-2 m-2 bg-orange-500 text-white text-md rounded-md disabled:bg-slate-600"
						on:click={toggleSetMP}>Set Move Pattern</button
					>
					<div class="relative grid grid-cols-1">
						{#each customPieces as piece}
							<!-- svelte-ignore a11y-click-events-have-key-events -->
							<!-- svelte-ignore a11y-no-static-element-interactions -->
							<div
								class="flex flex-cols items-center rounded-md
							 bg-gray-300 border border-gray-200 dark:border-gray-700 cursor-pointer my-1"
								on:click={() => selectPiece(piece.name?.toLowerCase(),piece.notation,"custom")}
							>
								<div class="w-2/3">
									<input
										class="text-blue-600 bg-gray-100 border-gray-300 focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
										type="radio"
										name="piece"
										value={piece.notation}
										bind:group={selectedPiece.piece.notation}
									/>
									<label
										for={piece.notation}
										class="w-full py-4 ml-2 text-md font-medium text-gray-900 cursor-pointer"
									>
										{piece.name}
									</label>
								</div>
								<img
									class="w-1/3 max-h-12"
									src={`/src/lib/assets/pieces/${color}/${piece.notation}.svg`}
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
						<span class="w-4 h-4 inline-block bg-blue-600 rounded-sm" />
						<p class="text-lg font-semibold ml-2">Slide Pattern:</p>
					</div>
					<TagInput {slideDirections} dropDownText="Select Directions" />
					<div class="flex items-center">
						<span class="w-4 h-4 inline-block bg-red-600 rounded-sm" />
						<p class="text-lg font-semibold ml-2">Jump Pattern:</p>
					</div>
					<!-- svelte-ignore a11y-label-has-associated-control -->
					<label class="relative inline-flex items-center cursor-pointer">
						<span class="m-3 text-md font-medium text-gray-900 dark:text-gray-300"
							>Select Jump Moves</span
						>
						<!-- svelte-ignore missing-declaration -->
						<Switch checked={false} />
					</label>
					<button
						class="p-2 m-2 bg-green-500 text-white text-md rounded-md hover:bg-slate-400"
						on:click={toggleSetMP}>Save Pattern</button
					>
					<button
						class="p-2 m-2 bg-red-500 text-white text-md rounded-md hover:bg-slate-400"
						on:click={cancel}>Cancel</button
					>
				</div>
			{/if}
		</div>
	</div>
</div>
