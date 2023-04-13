<template>
    <q-page class="content-panel bg-dark">

            <div class="board-panel">
                 <board :board-state="boardState" :is-flipped="isFlipped" ref="boardRef" @handle-square-click="handleMove"/>
            </div>
            
            <div class="column right-panel">
                <side-tabs @sidetab-action="performSideTabAction" :room-id="roomId" />
            </div>
        
    </q-page>

</template>


<script lang="ts">
import { RootState } from '../../store/state';
import { reactive, ref, toRef } from 'vue';
import { useStore } from 'vuex';
import Board from '../Board/Board.vue';
import SideTabs from './SideTabs/SideTabs.vue';
import { useRoute } from 'vue-router';
import { VALIDATE_MOVE } from '../../store/action_types';
import { SET_SRC_SELECTION } from '../../store/mutation_types';
import { IMoveInfoPayload, ISquareClick } from '@/types';

export default{
    components:{Board, SideTabs},
    setup(){
        const route = useRoute();
        const roomId = route.params.roomId.toString();
        const store = useStore<RootState>();
        const boardState = reactive(store.state.board);
        const isFlipped = ref(store.state.userInfo.curGameRole === 'p2');
        const boardRef = ref();
        
        function validateMove(destPos:{row:number,col: number}){
            let srcInfo = store.state.curStartPos
            if(srcInfo){
                let piece = store.state.userInfo.curGameRole === 'p1' ? 
                    srcInfo.piece.toUpperCase() : 
                    srcInfo.piece;
                let mvInfo: IMoveInfoPayload ={
                    roomId,
                    piece,
                    srcRow: srcInfo.row,
                    srcCol: srcInfo.col,
                    destRow: destPos.row,
                    destCol: destPos.col,
                }
                store.dispatch(VALIDATE_MOVE,mvInfo);
            }  
        }

        const handleMove = (payload:ISquareClick)=>{
            if(payload.clickType==='select-mv-square'){
                if(store.state.curStartPos){
                    if (store.state.curStartPos.row !== payload.row || store.state.curStartPos.col !== payload.col){
                        validateMove({row:payload.row,col:payload.col})
                    }
                    else{ store.commit(SET_SRC_SELECTION,null) }
                } else if(payload.piece){
                    store.commit(SET_SRC_SELECTION,{row:payload.row,col:payload.col,piece:payload.piece})
                }
                return
            } 
        }
        function performSideTabAction(type:string){
            if (type==="flip"){
                isFlipped.value = !isFlipped.value
                boardRef.value.flipBoard(isFlipped.value)
            }
        }
        return{
            boardState,
            isFlipped,
            roomId,
            boardRef,
            performSideTabAction,
            handleMove
        }
    }
}
</script>

<style scoped>
.content-panel{
    display:grid;
    gap: 1rem;
    grid-auto-flow: row;
    grid-template-columns: 1fr 1fr;
    padding: 1rem;
    
}
.board-panel{
    display: grid;
    grid-auto-flow: row;
    gap: 1rem;
    aspect-ratio: 1/1;
}
</style>