<script lang="ts">
	import pieceSvg from '$lib/assets/svg/piece.svg';
	import boardSvg from '$lib/assets/svg/board.svg';
	import { BoardType, type BoardConfig, Color } from '$lib/board/types';
	import { EditorSubType } from '$lib/components/types';
	import EditableBoard from '$lib/board/EditableBoard.svelte';
	import PieceEditor from '$lib/components/editor/PieceEditor.svelte';
	import BoardEditor from '$lib/components/editor/BoardEditor.svelte';
	import ExpandableCard from '$lib/components/ExpandableCard.svelte';
	import RulesEditor from '$lib/components/editor/RulesEditor.svelte';
	//import { browser } from '$app/environment';
	import MpEditBoard from '$lib/board/MPEditBoard.svelte';
	import { goto } from '$app/navigation';
	import {
		boardEditor,
		editorSubTypeSelected,
		pieceEditor,
		resetEditorStores,
		ruleEditor
	} from '$lib/store/editor';
	import { editorMaxBoard } from '$lib/board/board';
	import { onMount } from 'svelte';
	import type { CreateParams } from '$lib/store/stores';
	import { templateStore, gameId } from '$lib/store/stores';
	import { wsStore } from '$lib/websocket.js';
	import { Button } from '$lib/components/ui/button';
	import {
		Dialog,
		DialogContent,
		DialogHeader,
		DialogTitle,
		DialogFooter,
		DialogClose
	} from '$lib/components/ui/dialog';
	import { Input } from '$lib/components/ui/input';
	import { type Template } from '$lib/types';
	import { get } from 'svelte/store';
	import { toast } from '$lib/store/alert';
	import { createTemplate, getTemplate, updateTemplate } from '$lib/api/template';
	import { page } from '$app/stores';

	const defaultConfig: BoardConfig = {
		fen: 'rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR',
		dimensions: { ranks: 8, files: 8 },
		isFlipped: false,
		boardType: BoardType.Editor
	};

	let isEditMode = false;
	let templateId: string | null = null;
	let boardConfig: BoardConfig;
	onMount(async () => {
		resetEditorStores();
		const url = new URL(window.location.href);
		templateId = url.searchParams.get('tid');

		if (templateId) {
			isEditMode = true;
			try {
				const template: Template = await getTemplate(templateId);
				if (template.name) templateName = template.name;
				boardConfig = {
					fen: template.fen,
					dimensions: template.dimensions,
					boardType: BoardType.Editor
				};
				boardEditor.setDimensions(template.dimensions.ranks, template.dimensions.files);
				//templateStore.set(template);
			} catch (err) {
				toast.error('Failed to load template. Please try again.');
				boardConfig = defaultConfig;
			}
		} else {
			boardConfig = defaultConfig;
		}
	});

	let boardRef: any;
	let clearBoard: () => void;
	let shiftBoard: (direction: string) => void;
	function exitRoom() {
		goto('/home');
	}

	let isVariantRulesOn: boolean;
	$: {
		isVariantRulesOn = $ruleEditor.isViewVariantRulesOn;
	}

	let isPopupVisible = false;
	let playAsWhite = true;
	let templateName = '';

	const showPopup = () => (isPopupVisible = true);
	const hidePopup = () => (isPopupVisible = false);

	const confirmTemplate = async () => {
		await saveTemplate();
		hidePopup();
	};

	const getFEN = () => {
		let position = '';
		for (let i = 0; i < $boardEditor.ranks; i++) {
			let empty_count = 0;
			for (let j = 0; j < $boardEditor.files; j++) {
				let square = $editorMaxBoard[i][j];
				if (square.isPiecePresent || square.wall) {
					if (empty_count > 0) {
						position += `${empty_count}`;
						empty_count = 0;
					}
					if (square.wall) {
						position += '.';
					} else {
						position +=
							square.piece?.color == Color.BLACK
								? square.piece?.notation
								: square.piece?.notation.toUpperCase();
					}
				} else {
					empty_count += 1;
				}
			}
			if (empty_count > 0) {
				position += `${empty_count}`;
			}
			if (i != $boardEditor.ranks - 1) {
				position += '/';
			}
		}

		const turn = 'w';
		const castleRights = 'KQkq';
		const ep = '-';
		return `${position} ${turn} ${castleRights} ${ep} 0 0`;
	};
	let rule;
	$: {
		rule = $ruleEditor;
	}
	const generateGameTemplate = () => {
		const fen = getFEN();
		const gameConfig: Template = {
			name: templateName,
			variantType: $ruleEditor.variantType,
			dimensions: {
				ranks: $boardEditor.ranks,
				files: $boardEditor.files
			},
			fen,
			pieceProps: $pieceEditor.movePatterns
		};
		templateStore.setTemplate(gameConfig);
		return gameConfig;
	};

	const playGame = () => {
		const config = generateGameTemplate();
		const url = `ws://${import.meta.env.VITE_WS_HOST}/ws`;
		const params: CreateParams = {
			color: playAsWhite ? 'w' : 'b',
			sessionId: 'sdf',
			gameConfig: config,
			username: 'test'
		};

		wsStore.newWebSocketConnection(url, params, 'create');
	};
	$: {
		if ($gameId !== null) {
			goto(`/game/${$gameId}/waiting`);
		}
	}

	let randomPlaceHolders = [
		'Chess But Make It Weird',
		'Knights Gone Wild',
		'The Great Wall',
		'ELO from the Other Side',
		'Bishop Please',
		'Takes Takes Takes',
		'Holy Hell',
		'Google EP',
		"Chessn't"
	];
	const getRandomPlaceholder = () => {
		return randomPlaceHolders[Math.floor(Math.random() * randomPlaceHolders.length)];
	};

	// save template related functions
	export function getTemplatePayload(): Omit<Template, 'createdBy'> {
		const boardData = generateGameTemplate();
		const rules = get(ruleEditor);
		return {
			name: templateName,
			variantType: rules.variantType,
			dimensions: boardData.dimensions,
			fen: boardData.fen,
			pieceProps: {},
			customData: {}
		};
	}

	const saveTemplate = async () => {
		const payload = getTemplatePayload();
		try {
			if (isEditMode && templateId) {
				const updated = await updateTemplate(templateId, payload);
				toast.success(`Template "${templateName}" updated successfully!`);
			} else {
				const created = await createTemplate(payload);
				toast.success(`Template "${templateName}" saved successfully!`);
			}
		} catch (err) {
			console.error(err);
			toast.error('Failed to save template. Please try again.');
		}
	};
</script>

<svelte:head>
	<title>Editor - Varchess</title>
</svelte:head>
<div class="font-inter flex-grow">
	{#if boardConfig}
		{#if $ruleEditor.isViewVariantRulesOn && $ruleEditor.ruleComponent}
			<svelte:component this={$ruleEditor.ruleComponent} />
		{/if}
		<div class="flex-1 flex m-4 lg:flex-row flex-col">
			<div class="text-black rounded-md lg:w-5/12 mx-3 p-3 max-h-[45rem] overflow-y-auto">
				<div class="border-b border-gray-200 dark:border-gray-700 flex flex-col text-center">
					<div class="flex flex-col justify-center">
						<div class="m-2 flex gap-2 justify-center">
							<Button
								on:click={exitRoom}
								class="flex-1 text-lg px-3 py-2 rounded border font-medium
			bg-transparent text-red-600 border-red-600 hover:bg-red-500/10"
							>
								<i class="fa-solid fa-right-from-bracket mr-2" /> Exit
							</Button>

							<Button
								on:click={showPopup}
								class="flex-1 text-lg px-3 py-2 rounded border font-medium
			bg-transparent text-yellow-600 border-yellow-600 hover:bg-yellow-500/10"
							>
								<i class="fa-solid fa-floppy-disk mr-2" /> Save
							</Button>

							<Button
								on:click={playGame}
								class="flex-1 text-lg px-3 py-2 rounded border font-medium
			bg-transparent text-green-600 border-green-600 hover:bg-green-500/10"
							>
								<i class="fa-solid fa-play mr-2" /> Play
							</Button>
						</div>
					</div>
					<div class="grid grid-cols-2">
						<!-- svelte-ignore a11y-click-events-have-key-events -->
						<!-- svelte-ignore a11y-no-static-element-interactions -->
						<div
							class={`flex items-center rounded-md p-4 m-1.5 bg-white hover:bg-gray-500 hover:text-white  border border-gray-200 dark:border-gray-700 cursor-pointer`}
							on:click={() => (playAsWhite = true)}
						>
							<input
								class="cursor-pointer w-4 h-4 text-black-600 bg-white border-gray-300 focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
								type="radio"
								value="White"
								name="color"
								checked={playAsWhite}
							/>
							<label class="ml-2 cursor-pointer" for="White">Play as White</label>
						</div>
						<!-- svelte-ignore a11y-click-events-have-key-events -->
						<!-- svelte-ignore a11y-no-static-element-interactions -->
						<div
							class="flex items-center rounded-md pl-4 m-1.5 text-white bg-black hover:bg-gray-500 border border-gray-200 dark:border-gray-700 cursor-pointer"
							on:click={() => (playAsWhite = false)}
						>
							<input
								class="cursor-pointer w-4 h-4 text-white bg-black border-gray-300 focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
								type="radio"
								value="Black"
								name="color"
								checked={!playAsWhite}
							/>
							<label class="ml-2 cursor-pointer" for="Black">Play as Black</label>
						</div>
					</div>
					<ExpandableCard svg={boardSvg} title="Board Editor">
						<BoardEditor
							bind:dimensions={boardConfig.dimensions}
							on:shift={(e) => shiftBoard(e.detail)}
							on:clear={() => clearBoard()}
						/>
					</ExpandableCard>
					<ExpandableCard svg={pieceSvg} title="Piece Editor">
						<PieceEditor />
					</ExpandableCard>
					<ExpandableCard iconClass="fa-solid fa-clipboard-list fa-lg" title="Rules Editor">
						<RulesEditor />
					</ExpandableCard>
				</div>
			</div>
			<div class="rounded-md lg:w-7/12 mx-3 my-3 p-3">
				{#if $editorSubTypeSelected === EditorSubType.MovePattern}
					<MpEditBoard />
				{:else}
					<EditableBoard
						{boardConfig}
						bind:this={boardRef}
						bind:shift={shiftBoard}
						bind:clear={clearBoard}
					/>
				{/if}
			</div>
		</div>
		<Dialog open={isPopupVisible} onOpenChange={(v) => (isPopupVisible = v)}>
			<DialogContent class="sm:max-w-md">
				<DialogHeader>
					<DialogTitle>Save Template</DialogTitle>
				</DialogHeader>
				<div class="grid gap-4 py-4">
					<label for="templateName" class="text-sm font-medium">Template Name</label>
					<Input id="templateName" placeholder={getRandomPlaceholder()} bind:value={templateName} />
				</div>
				<DialogFooter class="flex justify-end gap-2">
					<DialogClose asChild>
						<Button variant="secondary" on:click={hidePopup}>Cancel</Button>
					</DialogClose>
					<Button on:click={confirmTemplate}>Confirm</Button>
				</DialogFooter>
			</DialogContent>
		</Dialog>
	{/if}
</div>

<style>
	::-webkit-scrollbar {
		width: 10px;
	}
</style>
