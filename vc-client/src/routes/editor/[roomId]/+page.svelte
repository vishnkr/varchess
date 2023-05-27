<script lang="ts">
	//import {VcBoardElement} from '$lib/chessboard/VCBoardElement'
	import BoardEditor from "$lib/components/editor/BoardEditor.svelte";
	import pieceSvg from "$lib/assets/svg/piece.svg";
    import boardSvg from "$lib/assets/svg/board.svg";
    import Tabs from "$lib/components/shared/Tabs.svelte";
	import type { BoardConfig } from "$lib/board/types";
	import Board from "$lib/board/Board.svelte";
	import PieceEditor from "$lib/components/editor/PieceEditor.svelte";
    let items = ['Custom','Predefined'];
    let activeItem = 'Custom';
    const tabChange=(e:CustomEvent<string>)=> activeItem = e.detail;
    
    export let boardConfig:BoardConfig = {
		fen: "rnbqkbnr/pnpppppp/p6p/5.../7P/P6P/PPPPPPPP/RNBQKBNR",
		dimensions: { ranks: 8, files: 8 },
		editable: false,
		interactive: false,
        isFlipped:false,
	}
    let shiftBoard:(direction:string)=>void;
    function handleShift(event:CustomEvent<string>){
        const direction = event.detail;
        shiftBoard(direction);
    }

</script>

<svelte:head>
	<title> Editor - Varchess</title>
</svelte:head>
<div class="font-inter grid min-h-screen text-zinc-900">
    <div class="flex-1 flex m-4"> 
        <div class="bg-zinc-700 rounded-md w-3/12 mx-3 p-3">
            <div class="border-b border-gray-200 dark:border-gray-700 flex flex-col text-center">
                <div class="flex justify-center">
                    <Tabs {activeItem} {items} on:tabChange={tabChange}/>
                </div>
                {#if activeItem==="Custom"}
                    <details class="bg-white shadow rounded group mb-4">
                        <summary class="list-none flat flex flex-wrap items-center cursor-pointer">
                            <h3 class="flex flex-1 p-4 text-xl font-semibold">Board Editor  </h3>
                            <img class="w-g h-6"  src={boardSvg} alt="board editor" />
                        <div class="w-8 flex items-center justify-center">
                            <div class="border-8 border-transparent border-l-gray-600 ml-2
                            group-open:rotate-90 transition-transform origin-left" />
                        </div>
                    </summary>
                        <div><BoardEditor 
                            bind:dimensions={boardConfig.dimensions}
                            on:shift={handleShift}
                            /></div>
                    </details>
                    <details class="bg-white shadow rounded grou mb-4">
                        <summary class="list-none flat flex flex-wrap items-center cursor-pointer">
                            <h3 class="flex flex-1 p-4 text-xl font-semibold">Piece Editor </h3>
                            <img class="w-g h-6" src={pieceSvg} alt="piece editor" />
                        <div class="w-8 flex items-center justify-center">
                            <div class="border-8 border-transparent border-l-gray-600 ml-2
                            group-open:rotate-90 transition-transform origin-left" />
                        </div>
                    </summary>
                    <div><PieceEditor /></div>
                    </details>
                {:else}
                <details class="bg-white shadow rounded group">
                    <summary class="list-none flat flex flex-wrap items-center cursor-pointer">
                        <h3 class="flex flex-1 p-4 font-semibold">Duck Chess</h3>
                    <div class="w-8 flex items-center justify-center">
                        <div class="border-8 border-transparent border-l-gray-600 ml-2
                        group-open:rotate-90 transition-transform origin-left"></div>
                    </div>
                </summary>
                <div><p>sfdg</p></div>
                </details>
                {/if}
            </div>
        </div>
        <div class="bg-zinc-700 rounded-md w-6/12 mx-3 p-3">
            
            <Board 
            boardConfig={boardConfig}
            bind:shift={shiftBoard}
            />
        </div>
        <div class="bg-zinc-700 text-white rounded-md w-3/12 mx-3 p-3">
           Right Panel
        </div>
    </div>
    
</div>