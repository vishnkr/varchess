<script lang="ts">
    import "./board-styles.css"
	import type { Dimensions, IPiece, PiecePositions, SquareColor, SquareIdx } from "./types";

    export let idx:SquareIdx;
    export let gridX:number;
    export let gridY:number;
    export let color:SquareColor;
    export let piece:IPiece|null = null;

    function getPieceClass(piece:IPiece){
        return piece.color.charAt(0).toLowerCase() + piece.pieceType.charAt(0).toLowerCase();
    }
</script>


<div style="--x:{ gridX }; --y:{gridY};" 
    data-square-color={color}
    class={`${piece ? getPieceClass(piece): null} bg-piece`}
    >
        <slot />
</div>

<style>
    [data-square-color] {
        width: 100%;
        height: 0;
        padding-bottom: 100%;
        grid-column: var(--y);
        grid-row: var(--x);
        background-color: var(--square-color);
    }

    [data-square-color="dark"] {
    --square-color: var(--default-dark-square);
    --p-label-color: var(--default-light-square);
    --p-square-color-hover: var(--square-color-dark-hover);
    --p-move-target-marker-color: var(--move-target-marker-color-dark-square);
    --p-square-color-active: var(--square-color-dark-active);
    --p-outline-color-active: var(--outline-color-dark-active);
  }
  
  [data-square-color="light"] {
    --square-color: var(--default-light-square);
    --p-label-color: var(--default-dark-square);
    --p-square-color-hover: var(--square-color-light-hover);
    --p-move-target-marker-color: var(--move-target-marker-color-light-square);
    --p-square-color-active: var(--square-color-light-active);
    --p-outline-color-active: var(--outline-color-light-active);
  }

</style>
