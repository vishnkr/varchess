<script lang="ts">
	import { getSquareColor, type Position, type SquareInfo, type SquareMaps, type PiecePresentInfo, type SquareIdx } from "./types";
    import Square from "./Square.svelte";
    import "./board-styles.css"
    import type { BoardConfig } from "./types";
	import { generateSquareMaps, updatePiecePositions } from "./board";

    export let boardConfig:BoardConfig;
    export let position: Position;
    
    export let boardSquares: Record<SquareIdx,SquareInfo>;
</script>

<div id="wrapper">

    <div id="board" style={`--size: ${Math.max(boardConfig.dimensions.ranks, boardConfig.dimensions.files)}`}>
        {#each Array(boardConfig.dimensions.ranks*boardConfig.dimensions.files) as _, idx}
            <Square 
                gridX={boardSquares[idx]?.gridX} 
                gridY={boardSquares[idx]?.gridY} 
                color={getSquareColor(boardSquares[idx]?.row,boardSquares[idx]?.column,boardConfig.isFlipped)}
                piece={position.piecePositions[idx] ?? null}
                disabled={position.disabled[idx] ?? false}
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
  }

  #border {
        display: grid;
        grid-template-columns: 1fr;
        grid-template-rows: repeat(var(--size), 1fr);
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        z-index:2;
    }

    .notation {
        display: flex;
        justify-content: center;
        align-items: center;
        font-size: 14px;
        font-weight: bold;
        color: black;
        pointer-events: none;
    }
</style>