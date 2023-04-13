<template>
    <q-card dark class="mp">
        <div class="piece-buttons"> 
            <span v-for="(pieceInfo, pieceKey) of movePatterns" :key="pieceKey"> 
                <button class="mp-btn" @click="changePiece(pieceKey)">
                    <img :src="getPieceURL(pieceKey)" />
                </button>
            </span>
        </div>
        <board :board-state="mpBoard" :is-flipped="false" ref="boardRef"/>
    </q-card>
</template>

<script lang="ts">
import { PropType, computed, reactive, ref } from 'vue';
import Board from '../../Board/Board.vue'
import { IMovePatterns } from '@/types';
import { setupMPBoard } from '../../../utils';
import { BoardState } from '@/classes';


export default{
    components:{Board},
    props:{
        movePatterns:{type: Object as PropType<IMovePatterns>,required:true}
    },
    setup(props){
        const selectedPiece = ref(Object.keys(props.movePatterns)[0]);
        const boardRef = ref();
        const [curRow,curCol] = [4,4];
        const mpBoard:BoardState = reactive(setupMPBoard(selectedPiece.value,'white',props.movePatterns[selectedPiece.value],curRow,curCol));
        function updateMpBoardState(){
            if(props.movePatterns[selectedPiece.value]){
               mpBoard.squares[curRow][curCol].squareInfo.pieceType = selectedPiece.value
               clearTempColors()
               updateTempColors()
            }
        }

        function clearTempColors(){
            for(let row =0; row<mpBoard.dimensions.rows;row++){
                for(let col =0; col<mpBoard.dimensions.cols;col++){
                    mpBoard.squares[row][col].updateSquareInfo({tempSquareColor: null})
                }
            }
        }

        function updateTempColors(){
            for(let [jumpX,jumpY] of props.movePatterns[selectedPiece.value].jumpPatterns){
                mpBoard.squares[curRow+jumpX][curCol+jumpY].updateSquareInfo({tempSquareColor: 'jump'});
            }
            for(let [dx,dy] of props.movePatterns[selectedPiece.value].slidePatterns){
                let [slideRow,slideCol] = [curRow+dx,curCol+dy];
                while(slideRow >= 0 && slideRow < mpBoard.dimensions.rows && slideCol >= 0 && slideCol < mpBoard.dimensions.cols){
                    mpBoard.squares[slideRow][slideCol].updateSquareInfo({tempSquareColor: 'slide'});
                    slideRow+=dx;
                    slideCol+=dy;
                }
            }
        }

        const getPieceURL = computed(()=>{
            return (piece:string) =>{
                const path = `/assets/images/pieces/white/${piece}.svg`
                return path
            }
        })

        const changePiece = (piece:string)=>{
            selectedPiece.value = piece
            updateMpBoardState()
            boardRef.value.updateBoardState1D(mpBoard)
        }


        return{
            changePiece,
            mpBoard,
            boardRef,
            getPieceURL
        }

    }
}

</script>

<style scoped>
.mp-btn{
    background-color: orange;
    width:10%;
    border-radius: 5px;
    margin: 1rem;
    cursor: pointer;
}

</style>