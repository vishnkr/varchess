<script lang="ts">
	
	import pieceSvg from '$lib/assets/svg/piece.svg';
	import boardSvg from '$lib/assets/svg/board.svg';
	import { BoardType, type BoardConfig, Color } from '$lib/board/types';
	import { EditorSubType } from '$lib/components/types';
	import { onMount } from 'svelte';
	import { editorSettings } from '$lib/board/stores';
	import {me,roomId,createWebSocket} from '$lib/store/stores';

	import Tabs from '$lib/components/shared/Tabs.svelte';
	import EditableBoard from '$lib/board/EditableBoard.svelte';
	import PieceEditor from '$lib/components/editor/PieceEditor.svelte';
	import BoardEditor from '$lib/components/editor/BoardEditor.svelte';
	import GameSettings from '$lib/components/editor/GameSettings.svelte';
	import ExpandableCard from '$lib/components/ExpandableCard.svelte';
	import Members from '$lib/components/shared/Members.svelte';
	import RulesEditor from '$lib/components/editor/RulesEditor.svelte';
	import Chat from '$lib/components/Chat.svelte';
	import { page } from '$app/stores';
	import { displayAlert } from '$lib/store/alert';

	let items = ['Room','Editor'];
	let activeItem = 'Room';
	let shareUrl = `${$page.url.origin}/join/${$roomId}`;
	const tabChange = (e: CustomEvent<string>) => (activeItem = e.detail);

	export let boardConfig: BoardConfig = {
		fen: 'rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR',
		dimensions: { ranks: 8, files: 8 },
		editable: true,
		interactive: false,
		isFlipped: false,
		boardType: BoardType.Editor,
	};

	let getMovePatternBoardConfig: ()=> BoardConfig = ()=>{
		const piece = $editorSettings.pieceSelection!;
		const pieceType = piece.color === Color.WHITE ? piece.pieceType.toUpperCase() : piece.pieceType;
		return {
			fen: `9/9/9/9/4${pieceType}4/9/9/9/9`,
			dimensions: { ranks: 9, files: 9 },
			editable: true,
			interactive: false,
			isFlipped: false,
			boardType: BoardType.MovePatternEditor
		}
	}

	const copyToClipboard = () => {
		navigator.clipboard.writeText(shareUrl);
	};

	let clearBoard: () => void;
	let shiftBoard: (direction: string) => void;

	let actions: { type: string; handler: () => void }[] = [
		{ type: 'Remove', handler: () => {} },
		{ type: 'Set as White', handler: () => {} },
		{ type: 'Set as Black', handler: () => {} }
	];

	onMount(()=>{
		if ($roomId && $me.username){
			try {
				createWebSocket($roomId,$me.username)
			} catch (e:any){
				displayAlert(e.message || 'An error occurred','DANGER',6000)
			}
			
		}
		
	})
</script>

<svelte:head>
	<title>Editor - Varchess</title>
</svelte:head>
<div class="font-inter text-zinc-90 flex-grow">
	<div class="flex-1 flex m-4 lg:flex-row flex-col">
		<div class="bg-zinc-700 rounded-md lg:w-5/12 mx-3 p-3 max-h-[45rem] overflow-y-auto">
			<div class="border-b border-gray-200 dark:border-gray-700 flex flex-col text-center">
				<div class="flex justify-center">
					<Tabs {activeItem} {items} on:tabChange={tabChange} />
				</div>
				{#if activeItem==="Editor"}
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
					<ExpandableCard iconClass="fa-solid fa-gear fa-lg" title="Game Settings">
						<GameSettings />
					</ExpandableCard>
				{:else if activeItem === "Room"}
				<div class="flex flex-col justify-center">
					<div class="flex mb-5 justify-center">
						<input
							class="border border-gray-300 px-4 py-2 text-white rounded-l max-w-64"
							bind:value={shareUrl}
							disabled
						/>
						
						<button
							class="bg-blue-600 hover:bg-blue-800 text-white font-semibold py-2 px-2"
							on:click={copyToClipboard}
						>
							Copy <i class="fa-solid fa-link" />
						</button>
						
					</div>
					<div>
						<button class="bg-green-600 hover:bg-green-800 text-white font-semibold py-2 px-4 rounded">Play</button>
					</div>
					
					<Chat />
					<Members />
				</div>
				{/if}
			</div>
		</div>
		<div class="rounded-md lg:w-7/12 mx-3 my-3 p-3">
		{#if $editorSettings.editorSubTypeSelected===EditorSubType.MovePattern}
			<EditableBoard boardConfig={getMovePatternBoardConfig()} />
		{:else}
			<EditableBoard {boardConfig} bind:shift={shiftBoard} bind:clear={clearBoard} />
		{/if}
		</div>
	</div>
</div>

<style>
	::-webkit-scrollbar {
		width: 10px;
	}
</style>
