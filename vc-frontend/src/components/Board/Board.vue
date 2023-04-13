<template>
    <div id="wrapper">
        <div id="board" :style="cssVars">
        <BoardSquare v-for="square in boardState1D" 
            :key="square.squareId" 
            :square="square"
            :editorState="editorState ?? null"
            @emit-square-click="handleSquareClick"
        />
    </div>
    </div>
       
    
</template>

<style scoped>

#wrapper{
    position: relative;
    width: 100%;
    height: 100%;
}
#board{
    background-color: #EAEAEA;
    grid-template-columns: repeat(var(--size), 1fr);
    grid-template-rows: repeat(var(--size), 1fr);
    box-shadow: 0 14px 28px rgba(0,0,0,0.25), 0 10px 10px rgba(0,0,0,0.22);
    display: grid;
    max-width: 1000px;
}

@media (min-width:320px)  { 
    #wrapper{
        position: relative;
        width: 100%;
    }
}
</style>

<script lang="ts">
import BoardSquare from './Square.vue'
import { ref, PropType, Ref, onMounted, computed } from 'vue';
import { IEditorState, ISquareClick, IMoveInfo } from '../../types';
import {Square, BoardState} from '../../classes';
import { isLight } from '../../utils';
import { useStore } from 'vuex';
import { RootState } from '@/store/state';
import { PERFORM_MOVE, SET_SRC_SELECTION } from '../../store/mutation_types';
export default{
    components:{BoardSquare},
    props:{
        boardState: {type: Object as PropType<BoardState>, required: true },
        isFlipped: {type:Boolean, required: true},
        boardSize: {type:Number},
        editorState:{type: Object as PropType<IEditorState>}
    },
    emits:['handle-square-click'],
    setup(props,{expose,emit}){
        
        const board: BoardState = props.boardState;
        const boardState1D : Ref<Square[]> = ref([]);
        const cssVars = computed(()=>{
            return {
            '--size': Math.max(board.dimensions.rows,board.dimensions.cols),
        }
        })
        const store = useStore<RootState>();
        store.subscribe((mutation,state)=>{
            if(mutation.type===PERFORM_MOVE && state.currentMove){
                performMove(state.currentMove)
            }
        })


        function performMove(moveInfo:IMoveInfo){
            board.squares[moveInfo.destRow][moveInfo.destCol].updateSquareInfo({
                isPiecePresent: true,
                pieceType:moveInfo.piece.toLowerCase(),
                pieceColor:moveInfo.piece === moveInfo.piece.toUpperCase()?'white' :'black'
            })

            board.squares[moveInfo.srcRow][moveInfo.srcCol].updateSquareInfo({
                isPiecePresent:false, 
                pieceType:undefined,
                pieceColor:undefined
            })
            
            updateBoardState1D(board,props.isFlipped);
            store.commit(SET_SRC_SELECTION,null) 
        }
        
        
        //Converts 2D board state into 1D 
        const updateBoardState1D = (boardState:BoardState,isFlipped:boolean)=>{
            let newBoardState = [];
            const endId = boardState.dimensions.rows * boardState.dimensions.cols - 1;
            let row, square, x=1, y=1, flipX = boardState.dimensions.rows, flipY = boardState.dimensions.cols;
            for (row of board.squares){
                for (square of row){
                    let newSquare: Square = new Square({
                        squareInfo:{
                            ...square.squareInfo,
                            row: x,
                            col: y,
                            squareColor: square.disabled ? 'disabled' : isLight(y, x) ? 'light' : 'dark',
                        },
                        disabled: square.disabled,
                        squareId: square.squareId && isFlipped ? endId - square.squareId : square.squareId || 0,
                        x: isFlipped ? flipX : x,
                        y: isFlipped ? flipY : y
                        });
                    y+=1;
                    flipY-=1;
                    newBoardState.push(newSquare)
                }
                x+=1
                flipX-=1
                flipY = boardState.dimensions.cols
                y=1
            }
            if(isFlipped){
                newBoardState = newBoardState.reverse();
            }
            boardState1D.value = newBoardState
        }

        const flipBoard = (flip:boolean)=>{
            updateBoardState1D(board,flip)
        }

        expose({updateBoardState1D,flipBoard})

        onMounted(()=>{
            updateBoardState1D(board,props.isFlipped)
        })

        
        const handleSquareClick = (payload:ISquareClick)=>{
            emit('handle-square-click',payload)
        }


        return{
            board,
            boardState1D,
            cssVars,
            handleSquareClick
        }


    }
}
</script>
