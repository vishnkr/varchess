<script lang="ts">
	import { getSquareColor, type Position, type SquareInfo, type SquareMaps, type PiecePresentInfo } from "./types";
    import Square from "./Square.svelte";
    import "./board-styles.css"
    import type { BoardConfig } from "./types";
	import { generateSquareMaps, updatePiecePositions } from "./board";
	import type { SquareIdx } from "$lib/chessboard/square";
	import { convertFenToPosition } from "./fen";

	
    
    export let boardConfig:BoardConfig;
    export const shift = (direction:string):void=>{
        let [lastCol,lastRow,afterLastCol,afterLastRow]  = [boardConfig.dimensions.files-1,boardConfig.dimensions.ranks-1,boardConfig.dimensions.files,boardConfig.dimensions.ranks];
        let tempSquares:PiecePresentInfo[][];
        switch(direction){
            case 'right':
                maxBoardState = maxBoardState.map((row,i)=>{
                    if(i<boardConfig.dimensions.ranks){
                        return [...row.slice(lastCol,afterLastCol),...row.slice(0,lastCol),...row.slice(afterLastCol)]
                    }
                    return row;
                })
                break;
            case 'left':
                maxBoardState = maxBoardState.map((row:PiecePresentInfo[],i) => {
                    let firstSquare = row[0];
                    if(i<boardConfig.dimensions.ranks){
                        return [...row.slice(1,afterLastCol),...[firstSquare],...row.slice(afterLastCol)]
                    } 
                    return row
                });
                console.log(maxBoardState)
                break;
            case 'up':
                let firstRowSquares = maxBoardState[0]
                tempSquares = maxBoardState.slice(1,afterLastRow)
                maxBoardState = [...tempSquares,...[firstRowSquares],...maxBoardState.slice(afterLastRow)]
                break;
            case 'down':
                let lastRowSquares = maxBoardState[lastRow];
                tempSquares = maxBoardState.slice(0,lastRow)
                maxBoardState = [...[lastRowSquares],...tempSquares,...maxBoardState.slice(afterLastRow)]
                break;
        }
        position.piecePositions = updatePiecePositions(maxBoardState,boardConfig.dimensions)
        
    };

    let squareMaps:SquareMaps = generateSquareMaps(boardConfig.dimensions,boardConfig.isFlipped ?? false);
    let boardSquares: Record<SquareIdx,SquareInfo> = squareMaps.squares;
    const convertedPos = convertFenToPosition(boardConfig.fen);
    let position:Position = {piecePositions:{}};
    let maxBoardState: PiecePresentInfo[][] = [];
    if (convertedPos){
        position=convertedPos.position;
        maxBoardState = convertedPos.maxBoardState;
    }
    function updateBoardState(){
        squareMaps = generateSquareMaps(boardConfig.dimensions,boardConfig.isFlipped ?? false);
        boardSquares = squareMaps.squares;
        position.piecePositions = updatePiecePositions(maxBoardState,boardConfig.dimensions);
    }
    

    $:{
        boardConfig.dimensions=boardConfig.dimensions
        updateBoardState()
        
    }
</script>
<div id="wrapper">
    <div id="board" style={`--size: ${Math.max(boardConfig.dimensions.ranks, boardConfig.dimensions.files)}`}>
        {#each Array(boardConfig.dimensions.ranks*boardConfig.dimensions.files) as _, idx}
            <Square 
                {idx} 
                gridX={boardSquares[idx]?.gridX} 
                gridY={boardSquares[idx]?.gridY} 
                color={getSquareColor(boardSquares[idx]?.row,boardSquares[idx]?.column,boardConfig.isFlipped)}
                piece={position.piecePositions[idx] ?? null}
            >
                
            </Square>
        {/each}
    </div>
</div>

<style>


#board{
  display: grid;
  width:100%;
  max-width:700px;
  justify-items: center;
  grid-template-columns: repeat(var(--size), 1fr);
  grid-template-rows: repeat(var(--size), 1fr);
  border: 1px solid var(--default-dark-square);
  background-color: hsla(0, 10%, 94%, 0.7);
  user-select: none;
  box-shadow: 0 14px 28px rgba(0,0,0,0.25), 0 10px 10px rgba(0,0,0,0.22);
  touch-action: none;
  border-collapse: collapse;
  box-sizing: border-box;
}

  #wrapper{
    position: relative;
    width: 100%;
    padding-bottom: 100%;
  }
</style>