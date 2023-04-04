<template>
    <div :style="cssVars" id="board-container">
        <div id="board">
        <BoardSquare v-for="square in boardState1D" 
            :key="square.squareId" 
            :square="square"
            :editorState="editorState ?? null"
            @emit-square-click="emitSquareClick"
        />
    </div>
    </div>
    
</template>

<script lang="ts">
import BoardSquare from './Square.vue'
import { ref,reactive, PropType, Ref, onMounted, computed } from 'vue';
import { BoardState, Square, EditorState } from '../../types';
import { isLight } from '../../utils';

export default{
    components:{BoardSquare},
    props:{
        boardState: {type: Object as PropType<BoardState>, required: true },
        isFlipped: {type:Boolean, required: true},
        boardSize: {type:Number},
        editorState:{type: Object as PropType<EditorState>}
    },
    emits:['handle-square-click'],
    setup(props,{expose,emit}){
        
        const board: BoardState = props.boardState;
        const boardState1D : Ref<Square[]> = ref([]);
   
        const cssVars = computed(()=>{
            return {
            '--container_size': props.boardSize ? `${props.boardSize}px` : `${700}px`,
            '--size': Math.max(board.dimensions.rows,board.dimensions.cols),
        }
        })

        //Converts 2D board state into 1D 
        const updateBoardState1D = (boardState:BoardState)=>{
            let newBoardState = [];
            let row, square, x=1, y=1, flipX = boardState.dimensions.rows, flipY = boardState.dimensions.cols;
            for (row of board.squares){
                for (square of row){
                    let newSquare :Square = {
                        disabled: square.disabled,
                        squareId: square.squareId ? square.squareId + (props.isFlipped ? -1 : 1) : 0,
                        x: props.isFlipped ? flipX : x,
                        y: props.isFlipped ? flipY : y,
                        squareInfo: {...square.squareInfo,
                            row: x,
                            col: y,
                            squareColor: square.disabled ? 'disabled' : isLight(y,x) ? 'light' : 'dark',
                        }
                    }
                    y+=1;
                    flipY-=1;
                    newBoardState.push(newSquare)
                }
                x+=1
                flipX-=1
                flipY = boardState.dimensions.cols
                y=1
            }
            if(props.isFlipped){
                newBoardState.reverse();
            }
            boardState1D.value = newBoardState
        }
        expose({updateBoardState1D})

        onMounted(()=>{
            updateBoardState1D(board)
        })

        const emitSquareClick = (payload:{clickType:string,row:number,col:number})=>{
            emit('handle-square-click',payload)
        }
        return{
            board,
            boardState1D,
            cssVars,
            emitSquareClick
        }


    }
}
</script>

<style scoped>
#board-container{
max-width: var(--container_size);
width:80%;
max-height: var(--container_size);

}

#board{
    flex:1;
    background-color: #EAEAEA;
    grid-template-columns: repeat(var(--size), 1fr);
    grid-template-rows: repeat(var(--size), 1fr);
    box-shadow: 0 14px 28px rgba(0,0,0,0.25), 0 10px 10px rgba(0,0,0,0.22);
    display: grid;
    justify-items: center;
}

</style>