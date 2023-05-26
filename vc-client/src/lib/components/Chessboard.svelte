<script lang="ts">
    import {BoardElement} from "../chessboard/BoardElement";
    import Board from "../board/Board.svelte";
	import type { BoardConfig } from "$lib/board/types";

    export let boardConfig:BoardConfig;
</script>

<Board boardConfig={boardConfig}/>
<!--<board-element id="board" fen="rrbb/4/4/4/4/RRBB" ranks="6" files="4"/>-->
