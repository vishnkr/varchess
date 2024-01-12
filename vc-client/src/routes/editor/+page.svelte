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
	import { browser } from '$app/environment';
	import MpEditBoard from '$lib/board/MPEditBoard.svelte';
	import { beforeNavigate, goto } from '$app/navigation';
	import { boardEditor, editorSubTypeSelected, pieceEditor, ruleEditor} from '$lib/store/editor';
	import { editorMaxBoard } from '$lib/board/board';
	import { onMount } from 'svelte';
	import type { Config, ConnectParams } from '$lib/store/stores';
	import { camelToSnake } from '$lib/utils/index';
	import { configStore, gameId, wsStore } from '$lib/store/stores';
	
	let stonkfish: typeof import ('stonkfish-wasm');

	let activeItem = 'Room';
	const tabChange = (e: CustomEvent<string>) => (activeItem = e.detail);
	//export let data;
	//let username:string = $page.data.username;
	
	export let boardConfig: BoardConfig = {
		fen: 'rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR',
		dimensions: { ranks: 8, files: 8 },
		editable: true,
		interactive: false,
		isFlipped: false,
		boardType: BoardType.Editor,
	};

	//let dirty = true;

	/*const copyToClipboard = () => {
		if(roomId){
			navigator.clipboard.writeText(roomId);
		}
	};*/

	onMount(async () => {
		stonkfish = await import('stonkfish-wasm');
		await stonkfish.default();
	});

	/*beforeNavigate(({ cancel }) => {
		if (dirty) {
			const confirmMessage = "Leaving this page might result in loss. Are you sure you want to leave?";
			if (!confirm(confirmMessage)) {
			cancel();
			}
		}
	});*/
	let clearBoard: () => void;
	let shiftBoard: (direction: string) => void;
	function exitRoom(){
		if (browser) { 
			window.location.href = '/home';
		}
	}

	// Save Template 
	let isPopupVisible = false;
	let templateName = '';

	const showPopup = () => {
		isPopupVisible = true;
	};

	const hidePopup = () => {
		isPopupVisible = false;
	};
	const confirmTemplate = () => {
		// TODO: construct json and send post request 
    	templateName = '';
    	hidePopup();
  	};


	// Play Game 
	const getFEN = ()=>{
		let position = '';
		for(let i=0; i<$boardEditor.ranks;i++){
			let empty_count = 0;
			for(let j=0;j<$boardEditor.files;j++){
				let square = $editorMaxBoard[i][j]
				if (square.isPiecePresent || square.wall){
					if (empty_count>0){position+=`${empty_count}`; empty_count=0;}
					if (square.wall){
						position+="."
					} else{
						position+= square.piece?.color==Color.BLACK ? square.piece?.pieceType : square.piece?.pieceType.toUpperCase(); 
					}
				} else{
					empty_count+=1;
				}
			}
			if (i!=$boardEditor.ranks-1){
				if (empty_count>0){position+=`${empty_count}`;}
				position+="/";
			}
		}
		
		const turn = 'w';
		const castleRights = "KQkq";
		const ep="-";
		return `${position} ${turn} ${castleRights} ${ep} 0 0`
	}

	const generateGameConfigJSON = () =>{
		let boardEditorState = $boardEditor;
		let pieceEditorState = $pieceEditor;
		let ruleEditorState =  $ruleEditor;

		const fen = getFEN();
		const gameConfig:Config = {
			variantType: ruleEditorState.variantType,
			dimensions: {
				ranks : boardEditorState.ranks,
				files : boardEditorState.files
			},
			fen,
			pieceProps: pieceEditorState.movePatterns,
		}
		configStore.setConfig(gameConfig);
		return gameConfig;
	}
	const playGame = ()=>{
		const config  = generateGameConfigJSON();		
		console.log(config);
		const config_json = JSON.stringify(camelToSnake(config))
		const chesscore = new stonkfish.ChessCoreLib(config_json);
		console.log(chesscore.getLegalMoves());
		const url = `ws://${import.meta.env.VITE_WS_HOST}/ws`;
		const params:ConnectParams = {
			connectType:"create",
			sessionId: "sdf",
			gameConfig: config
		}
		wsStore.newWebSocketConnection(url,params);
	}
	$: {
		console.log('gmae id up',$gameId)
    	if ($gameId !== null) { goto(`/game/${$gameId}?waiting=true`); }
  	}

</script>

<svelte:head>
	<title>Editor - Varchess</title>
</svelte:head>
<div class="font-inter flex-grow">
	<div class="flex-1 flex m-4 lg:flex-row flex-col">
		<div class="text-black rounded-md lg:w-5/12 mx-3 p-3 max-h-[45rem] overflow-y-auto">
			<div class="border-b border-gray-200 dark:border-gray-700 flex flex-col text-center">
					<div class="flex flex-col justify-center">
						<div class="m-2 flex">
							<a href="/home">
								<button class="w-32 bg-red-600 hover:bg-red-800 text-white rounded-md font-semibold py-2 px-2 mx-3" on:click={exitRoom}>
									Exit
								</button>
							</a>
							<button class="w-32 bg-orange-600 hover:bg-orange-800 text-white rounded-md font-semibold py-2 px-2 mx-3" on:click={showPopup}>
								Save
							</button>
							<button on:click={playGame} class="w-32 bg-green-600 hover:bg-green-800 text-white font-semibold py-2 px-4 rounded">Play</button>
							<form action="/?play"></form>
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
		{#if $editorSubTypeSelected===EditorSubType.MovePattern}
			<MpEditBoard />
		{:else}
			<EditableBoard {boardConfig} bind:shift={shiftBoard} bind:clear={clearBoard} />
		{/if}
		</div>
	</div>
	{#if isPopupVisible}
      <div class="absolute inset-0 flex items-center justify-center bg-gray-800 bg-opacity-50 z-50">
        <div class="bg-white p-4 rounded-md">
          <label for="templateName">Enter Template Name:</label>
          <input type="text" id="templateName" bind:value={templateName} class="border p-1 rounded-md" />
          <button on:click={confirmTemplate} class="bg-green-500 text-white px-4 py-2 rounded-md mt-2">Confirm</button>
          <button on:click={hidePopup} class="bg-red-500 text-white px-4 py-2 rounded-md mt-2 ml-2">Cancel</button>
        </div>
      </div>
    {/if}
</div>

<style>
	::-webkit-scrollbar {
		width: 10px;
	}
</style>
