<script lang="ts">
import { pieceEditor } from '$lib/store/editor';
import { MoveType } from '$lib/store/types';
import Board from './Board.svelte';
import { BoardType, type BoardConfig, Color } from './types';

let getMovePatternBoardConfig: ()=> BoardConfig = ()=>{
		const piece = $pieceEditor.pieceSelection!;
		const pieceType = piece.color === Color.WHITE ? piece.pieceType.toUpperCase() : piece.pieceType;
		return {
			fen: `9/9/9/9/4${pieceType}4/9/9/9/9`,
			dimensions: { ranks: 9, files: 9 },
			editable: false,
			interactive: false,
			isFlipped: false,
			boardType: BoardType.MovePatternEditor
		}
}
    let center = 4;
    const piece = $pieceEditor.pieceSelection;
    let movePatterns = null;
    let mpSquares: Record<number,MoveType> = {};

    function isSquareInBounds(row:number,col:number):boolean{
        return row>=0 && row<9 && col>=0 && col<9
    }
    $: if(piece){
        movePatterns = $pieceEditor.movePatterns[piece.pieceType];
        mpSquares={}
        if(movePatterns){
            movePatterns.jumpOffsets?.forEach((offset)=>{
            let newId = (center+offset[0])*9+(center+offset[1]);
            if(newId>=0 && newId<81){
                mpSquares[newId]=MoveType.Jump
            }
            })
            movePatterns.slideOffsets?.forEach((offset)=>{
                let x = offset[0];
                let y = offset[1];
                let newId = (center+x)*9+(center+y);
                while(isSquareInBounds(center+x,center+y)){
                    mpSquares[newId]=MoveType.Slide
                    x+=offset[0];
                    y+=offset[1];
                    newId = (center+x)*9+(center+y);
                }
            })
        }
    }
        
		
</script>


<Board customBoardId="board" {mpSquares} boardConfig={getMovePatternBoardConfig()}/>