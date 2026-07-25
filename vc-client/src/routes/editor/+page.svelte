<script lang="ts">
	import PieceIcon from '$lib/assets/svg/PieceIcon.svelte';
	import BoardIcon from '$lib/assets/svg/BoardIcon.svelte';
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
		ruleEditor,
		wormholeEditor
	} from '$lib/store/editor';
	import { editorMaxBoard } from '$lib/board/board';
	import { onMount } from 'svelte';
	import type { CreateParams } from '$lib/store/stores';
	import { templateStore, gameId, Status, gameState } from '$lib/store/stores';
	import { authStore } from '$lib/store/auth.js';
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
	import { type Template, VariantType } from '$lib/types';
	import { get } from 'svelte/store';
	import { toast } from '$lib/store/alert';
	import { createTemplate, getTemplate, updateTemplate } from '$lib/api/template';
	import { page } from '$app/stores';
	import { createGame } from '$lib/api/games';
	import { toWirePieceProps, fromWirePieceProps } from '$lib/utils/pieceProps';
	import {
		parseWormholePairs,
		validateVariantConfig,
		wormholePairsFromWire,
		wormholePairsToWire
	} from '$lib/utils/wormhole';
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
		const auth = localStorage.getItem('auth');
		if (!auth) {
			goto('/login');
			return;
		}
		resetEditorStores();
		// Placeholder so EditableBoard still hydrates once from this page's FEN.
		editorMaxBoard.set([[]]);
		const url = new URL(window.location.href);
		templateId = url.searchParams.get('tid');

		if (templateId) {
			isEditMode = true;
			try {
				const template: Template = await getTemplate(templateId);
				const fen = template.fen ?? template.position?.fen ?? defaultConfig.fen;
				const dimensions =
					template.dimensions ?? template.position?.dimensions ?? defaultConfig.dimensions;
				const pieceProps = template.pieceProps ?? template.position?.pieceProps ?? {};
				if (template.name) templateName = template.name;
				boardConfig = {
					fen,
					dimensions,
					boardType: BoardType.Editor
				};
				boardEditor.setDimensions(dimensions.ranks, dimensions.files);
				if (template.variantType) {
					ruleEditor.updateVariantType(template.variantType as VariantType);
				}
				if (template.customData) {
					const custom = { ...template.customData };
					if (template.variantType === VariantType.Wormhole) {
						const wirePairs = parseWormholePairs(custom.wormholePairs);
						custom.wormholePairs = wormholePairsFromWire(
							wirePairs,
							dimensions.files,
							dimensions.ranks
						);
						custom.cooldownTurns = 1;
					}
					ruleEditor.updateCustomData(custom);
				}
				pieceEditor.setMovePatterns(fromWirePieceProps(pieceProps));
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
		wormholeEditor.stop();
		goto('/home');
	}

	let isPopupVisible = false;
	let playAsWhite = true;
	let templateName = '';

	const showPopup = () => {
		if (!assertVariantReady()) return;
		isPopupVisible = true;
	};
	const hidePopup = () => (isPopupVisible = false);

	function assertVariantReady(): boolean {
		const err = validateVariantConfig($ruleEditor.variantType, $ruleEditor.customData);
		if (err) {
			toast.error(err);
			return false;
		}
		return true;
	}

	const confirmTemplate = async () => {
		if (!assertVariantReady()) return;
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
			pieceProps: toWirePieceProps($pieceEditor.movePatterns),
			customData: $ruleEditor.customData ?? {}
		};
		templateStore.setTemplate(gameConfig);
		return gameConfig;
	};
	/*const playGame = () => {
		const config = generateGameTemplate();
		const url = `ws://${import.meta.env.VITE_WS_HOST}/ws`;
		const params: CreateParams = {
			color: playAsWhite ? 'w' : 'b',
			sessionId: 'sdf',
			gameConfig: config,
			username: 'test'
		};

		wsStore.newWebSocketConnection(url, params, 'create');
	};*/
	/*$: {
		if ($gameId !== null) {
			goto(`/play/${$gameId}`);
		}
	}*/

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
		const customData = { ...(rules.customData ?? {}) };
		if (rules.variantType === VariantType.Wormhole) {
			const pairs = parseWormholePairs(customData.wormholePairs);
			customData.wormholePairs = wormholePairsToWire(
				pairs,
				boardData.dimensions.files,
				boardData.dimensions.ranks
			);
			customData.cooldownTurns = 1;
		}
		return {
			name: templateName,
			variantType: rules.variantType,
			dimensions: boardData.dimensions,
			fen: boardData.fen,
			pieceProps: boardData.pieceProps,
			customData
		};
	}

	const saveTemplate = async () => {
		if (templateName.length == 0) {
			toast.error('Template name cannot be empty.');
			return;
		}
		if (!assertVariantReady()) return;
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

	const playGame = async () => {
		const { userId, accessToken } = get(authStore);
		if (!userId || !accessToken) {
			goto('/login');
			return;
		}
		if (!assertVariantReady()) return;
		const gameConfig = getTemplatePayload();
		// Prefer the live editor config (includes pieceProps). templateId alone
		// would ignore unsaved pattern edits.
		const payload = { gc: gameConfig };
		try {
			console.log('cr payload', payload);
			let { gameId: newGameId } = await createGame(payload);
			console.log('ggame', newGameId);
			let colorPref = playAsWhite ? 'w' : 'b';
			gameState.updateStatus(Status.Waiting);
			gameId.set(newGameId);
			const url = `ws://${import.meta.env.VITE_WS_HOST}/play/${newGameId}`;

			const connectPayload = {
				token: accessToken,
				userId: userId,
				colorPref: colorPref
			};
			wsStore.newWebSocketConnection(url, connectPayload);
			goto(`/play/${newGameId}`);
		} catch (err) {
			console.log(err);
			toast.error('Unable to start game. Please try again later.');
		}
	};
</script>

<svelte:head>
	<title>Editor - Varchess</title>
</svelte:head>
<div class="font-inter flex-grow">
	{#if boardConfig}
		<div class="flex-1 flex m-4 lg:flex-row flex-col">
			<div class="text-black rounded-md lg:w-5/12 mx-3 p-3 max-h-[45rem] overflow-y-auto">
				<div class="border-b border-gray-200 dark:border-gray-700 flex flex-col text-center">
					<div class="flex flex-col justify-center">
						<div class="m-2 flex gap-2 justify-center">
							<Button
								on:click={exitRoom}
								class="flex-1 text-lg px-3 py-2 rounded border font-medium
			bg-transparent text-red-600 border-red-600 hover:bg-red-600/10"
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
					<!-- "Play as" Section -->
					<div class="flex items-center justify-center my-4 gap-4">
						<!-- Play as Text -->
						<div class="text-md dark:text-white font-semibold">Play as</div>

						<!-- Play as Buttons -->
						<!-- svelte-ignore a11y-click-events-have-key-events -->
						<!-- svelte-ignore a11y-no-static-element-interactions -->
						<div class="flex gap-2">
							<!-- Play as White -->
							<div
								class={`flex items-center justify-center rounded-md px-4 py-2 border cursor-pointer transition text-sm
									${
										playAsWhite
											? 'border-green-500 bg-green-100 dark:bg-green-800/20 text-black dark:text-white'
											: 'border-gray-300 dark:border-gray-600 bg-white text-black dark:hover:text-white hover:bg-gray-100 dark:hover:bg-gray-700/30'
									}
								`}
								on:click={() => (playAsWhite = true)}
							>
								<label class="cursor-pointer" for="White">White</label>
							</div>

							<!-- Play as Black -->
							<!-- svelte-ignore a11y-click-events-have-key-events -->
							<div
								class={`flex items-center justify-center rounded-md px-4 py-2 border cursor-pointer transition text-sm
									${
										!playAsWhite
											? 'border-green-500 bg-green-100 dark:bg-green-800/20 text-black dark:text-white'
											: 'border-gray-300 dark:border-gray-600 bg-black text-white hover:text-black dark:hover:text-white hover:bg-gray-100 dark:hover:bg-gray-700/30'
									}
								`}
								on:click={() => (playAsWhite = false)}
							>
								<label class="cursor-pointer" for="Black">Black</label>
							</div>
						</div>
					</div>

					<ExpandableCard svg={BoardIcon} title="Board Editor">
						<BoardEditor
							bind:dimensions={boardConfig.dimensions}
							on:shift={(e) => shiftBoard(e.detail)}
							on:clear={() => clearBoard()}
						/>
					</ExpandableCard>
					<ExpandableCard svg={PieceIcon} title="Piece Editor">
						<PieceEditor />
					</ExpandableCard>
					<ExpandableCard iconClass="fa-solid fa-clipboard-list fa-lg" title="Rule Editor">
						<RulesEditor />
					</ExpandableCard>
				</div>
			</div>
			<div class="rounded-md lg:w-7/12 mx-3 my-3 p-3">
				{#if $editorSubTypeSelected === EditorSubType.MovePattern}
					<MpEditBoard />
				{/if}
				<!-- Keep mounted while editing move patterns so placements aren't wiped on remount. -->
				<div
					class={$editorSubTypeSelected === EditorSubType.MovePattern ? 'hidden' : ''}
				>
					<EditableBoard
						{boardConfig}
						bind:this={boardRef}
						bind:shift={shiftBoard}
						bind:clear={clearBoard}
					/>
				</div>
			</div>
		</div>
		<Dialog open={isPopupVisible} onOpenChange={(v) => (isPopupVisible = v)}>
			<DialogContent class="sm:max-w-md">
				<DialogHeader>
					<DialogTitle>Save Template</DialogTitle>
				</DialogHeader>
				<div class="grid gap-4 py-4">
					<label for="templateName" class="text-sm font-medium">Template Name</label>
					<Input
						id="templateName"
						required
						placeholder={getRandomPlaceholder()}
						bind:value={templateName}
					/>
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
