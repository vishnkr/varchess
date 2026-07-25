<script lang="ts">
	import { Color } from '$lib/board/types';
	import { editorSubTypeSelected, jumpPatternEditing, pieceEditor } from '$lib/store/editor';
	import { EditorSubType } from '../types';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import type { PieceSelection } from '$lib/types';
	import { Button } from '../ui/button';
	import { Switch } from '../ui/switch';
	import { get } from 'svelte/store';

	const standardPieces: { name: string; notation: string }[] = [
		{ name: 'Pawn', notation: 'p' },
		{ name: 'King', notation: 'k' },
		{ name: 'Queen', notation: 'q' },
		{ name: 'Bishop', notation: 'b' },
		{ name: 'Knight', notation: 'n' },
		{ name: 'Rook', notation: 'r' }
	];
	const customPieces: { name: string; notation: string }[] = [
		{ name: 'Dolphin', notation: 'd' },
		{ name: 'Ninja', notation: 'i' },
		{ name: 'Unicorn', notation: 'u' },
		{ name: 'Giraffe', notation: 'g' },
		{ name: 'Juicer', notation: 'j' },
		{ name: 'Astronaut', notation: 'a' },
		{ name: 'Phage', notation: 'v' },
		{ name: 'Zebra', notation: 'z' }
	];

	/** [Δrow, Δcol] on the pattern board (row decreases toward top / "North"). */
	const slideDirections: Record<string, number[]> = {
		North: [-1, 0],
		East: [0, 1],
		South: [1, 0],
		West: [0, -1],
		'North East': [-1, 1],
		'North West': [-1, -1],
		'South East': [1, 1],
		'South West': [1, -1]
	};

	let color: Color = Color.WHITE;
	let selectedPiece: PieceSelection = {
		piece: { notation: 'p', pieceType: 'pawn', color: color },
		group: 'standard'
	};
	let setMovePattern = false;
	let selectedSlideDirections: string[] = [];

	function syncSlideSelectionFromStore() {
		const notation = selectedPiece.piece.notation.toLowerCase();
		const pat = get(pieceEditor).movePatterns[notation];
		const dirs = pat?.slideDirections ?? [];
		selectedSlideDirections = Object.entries(slideDirections)
			.filter(([, offset]) => dirs.some((d) => d[0] === offset[0] && d[1] === offset[1]))
			.map(([name]) => name);
	}

	const toggleSetMP = () => {
		setMovePattern = !setMovePattern;
		if (setMovePattern) {
			syncSlideSelectionFromStore();
			jumpPatternEditing.set(true);
		}
		editorSubTypeSelected.update(() =>
			setMovePattern ? EditorSubType.MovePattern : EditorSubType.Piece
		);
	};

	const saveMovePattern = () => {
		toggleSetMP();
	};

	const cancel = () => {
		pieceEditor.deletePiecePattern(selectedPiece.piece.notation);
		selectedSlideDirections = [];
		toggleSetMP();
	};

	const selectPiece = (pieceType: string, notation: string, group: string) => {
		selectedPiece = {
			...selectedPiece,
			piece: {
				notation,
				pieceType,
				color
			},
			group
		};
		pieceEditor.update((val) => ({
			...val,
			pieceSelection: {
				piece: {
					pieceType: selectedPiece.piece.pieceType,
					color,
					notation: selectedPiece.piece.notation
				},
				group: selectedPiece.group
			}
		}));
	};

	const updateColor = (newColor: Color) => {
		selectedPiece.piece.color = newColor;
		pieceEditor.update((val) => ({
			...val,
			pieceSelection: {
				piece: {
					pieceType: selectedPiece.piece.pieceType,
					color: newColor,
					notation: selectedPiece.piece.notation
				},
				group: selectedPiece.group
			}
		}));
		color = newColor;
	};

	function toggleDirection(direction: string) {
		const offset = slideDirections[direction];
		if (!offset) return;
		const notation = selectedPiece.piece.notation.toLowerCase();
		const on = selectedSlideDirections.includes(direction);
		if (on) {
			selectedSlideDirections = selectedSlideDirections.filter((d) => d !== direction);
			pieceEditor.removeSlidePattern(notation, offset);
		} else {
			selectedSlideDirections = [...selectedSlideDirections, direction];
			pieceEditor.addSlidePattern(notation, offset);
		}
	}
</script>

<div class="bg-white dark:bg-darkbg2 text-black dark:text-lightbg p-4 rounded-lg">
	<div class="grid grid-rows">
		<div class="mb-4 text-center">
			<p class="text-sm text-gray-600 dark:text-gray-400">
				Select a color and piece, then click on the board to place it.
			</p>
		</div>
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

					<!-- svelte-ignore a11y-click-events-have-key-events -->
					<!-- svelte-ignore a11y-no-static-element-interactions -->
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
							class={`flex items-center rounded-md pl-4 m-1.5 cursor-pointer border-2 transition dark:bg-gray-700
			  ${
					selectedPiece.piece.notation === piece.notation && selectedPiece.group === 'standard'
						? 'border-green-500 bg-green-100 dark:bg-green-800/20'
						: 'border-gray-300 dark:border-gray-600 hover:bg-gray-100 dark:hover:bg-gray-700/30'
				}
			`}
							on:click={() => selectPiece(piece.name?.toLowerCase(), piece.notation, 'standard')}
						>
							<label
								class="w-full py-4 ml-2 text-md font-medium text-gray-900 cursor-pointer dark:text-white"
							>
								{piece.name}
							</label>
							<img src={`/src/lib/assets/pieces/${color}/${piece.notation}.svg`} alt="piece" />
						</div>
					{/each}
				</div>

				<div class="px-2 m-1.5 py-2">
					<h3 class="text-xl dark:text-white">Custom Pieces</h3>
					<Button
						disabled={selectedPiece.group !== 'custom'}
						class="p-2 m-2 bg-transparent text-black border-black dark:text-white dark:border-white hover:bg-orange-700/10 dark:hover:bg-orange-400/20 text-md rounded border font-medium disabled:bg-slate-600"
						on:click={toggleSetMP}
					>
						Set Move Pattern
					</Button>
					<div class="relative grid grid-cols-2">
						{#each customPieces as piece}
							<!-- svelte-ignore a11y-click-events-have-key-events -->
							<!-- svelte-ignore a11y-no-static-element-interactions -->
							<div
								class={`flex items-center rounded-md p-2 m-1 cursor-pointer border-2 transition dark:bg-gray-700
			  ${
					selectedPiece.piece.notation === piece.notation && selectedPiece.group === 'custom'
						? 'border-green-500 bg-green-100 dark:bg-green-800/20'
						: 'border-gray-300 dark:border-gray-600 hover:bg-gray-100 dark:hover:bg-gray-700/30'
				}
			`}
								on:click={() => selectPiece(piece.name?.toLowerCase(), piece.notation, 'custom')}
							>
								<div class="w-2/3">
									<!-- svelte-ignore a11y-label-has-associated-control -->
									<label
										class="w-full py-4 ml-2 text-md font-medium text-gray-900 cursor-pointer dark:text-white"
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
				<div class="px-2 m-1.5 py-2 flex flex-col gap-4">
					<h1 class="text-xl font-bold dark:text-white">
						Set Move Pattern — {selectedPiece.piece.pieceType}
					</h1>

					<div class="flex flex-col gap-1">
						<div class="flex items-center gap-2">
							<span class="w-4 h-4 bg-blue-600 rounded-sm" />
							<p class="text-lg font-semibold dark:text-white">Slide Pattern:</p>

							<DropdownMenu.Root>
								<DropdownMenu.Trigger
									class="inline-flex items-center justify-center rounded-md bg-white px-3 py-1 text-sm font-medium text-black shadow-sm border border-gray-300 hover:bg-gray-100 dark:bg-gray-700 dark:text-white dark:border-gray-600"
								>
									{selectedSlideDirections.length
										? `${selectedSlideDirections.length} selected`
										: 'Select'}
								</DropdownMenu.Trigger>

								<DropdownMenu.Content
									class="w-56 bg-white text-black border border-gray-200 shadow-md dark:bg-gray-800 dark:text-white dark:border-gray-700"
								>
									<DropdownMenu.Label class="text-gray-700 dark:text-gray-300"
										>Slide Pattern</DropdownMenu.Label
									>
									<DropdownMenu.Separator class="bg-gray-200" />

									{#each Object.keys(slideDirections) as direction}
										<DropdownMenu.CheckboxItem
											checked={selectedSlideDirections.includes(direction)}
											on:click={() => toggleDirection(direction)}
										>
											{direction}
										</DropdownMenu.CheckboxItem>
									{/each}
								</DropdownMenu.Content>
							</DropdownMenu.Root>
						</div>
						<p class="text-sm text-gray-600 ml-6 dark:text-gray-400">
							Select directions where the piece can slide continuously.
						</p>
					</div>

					<div class="flex flex-col gap-1">
						<div class="flex items-center gap-2">
							<span class="w-4 h-4 bg-red-600 rounded-sm" />
							<div class="flex items-center space-x-2">
								<Label for="jump-pattern" class="text-lg font-semibold dark:text-white">
									Jump Pattern
								</Label>
								<Switch
									id="jump-pattern"
									checked={$jumpPatternEditing}
									on:click={(e) => {
										e.preventDefault();
										jumpPatternEditing.update((v) => !v);
									}}
								/>
							</div>
						</div>
						<p class="text-sm text-gray-600 ml-6 dark:text-gray-400">
							{#if $jumpPatternEditing}
								Click empty squares on the board to add or remove jump offsets.
							{:else}
								Enable this to click squares and define custom jump moves.
							{/if}
						</p>
					</div>

					<div class="flex gap-2 items-end justify-center mt-4">
						<Button
							on:click={saveMovePattern}
							class="p-2 bg-transparent rounded border font-medium text-md text-green-600 border-green-600 hover:bg-green-600/10"
						>
							Done
						</Button>

						<Button
							on:click={cancel}
							class="p-2 bg-transparent rounded border font-medium text-md text-red-600 border-red-600 hover:bg-red-600/10"
						>
							Clear & Cancel
						</Button>
					</div>
				</div>
			{/if}
		</div>
	</div>
</div>
