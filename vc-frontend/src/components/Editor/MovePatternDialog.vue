<template>
    <q-dialog 
        v-model="showDialog" 
        transition-show="slide-up"
        transition-hide="slide-down"
        full-width
        persistent

    >
        <q-card class="text-white mp-card" dark>
            <div class="top-row">
                <q-btn color="negative" label="Close" @click="closeDialog"></q-btn>
                <q-btn color="positive" label="Save" @click="saveMP"></q-btn>
            </div>
            
            <div class="flex-panels">
                <Board 
                    :board-state="boardState" 
                    :isFlipped="false" 
                    @handle-square-click="updateJumpMP"
                    :editor-state="editorState"
                    ref="boardRef"
                />
                <div class="options">
                    <div class="text-white text-h6 bg-dark">Move Type</div>
                    <q-radio v-model="editorState.moveType" val="jump" label="Jump" color="blue" keep-color/>
                    <q-radio v-model="editorState.moveType" val="slide" label="Slide" color="red" keep-color/>
                    <div v-if="editorState.moveType == 'slide'" class="slide-checkboxes">
                        <q-checkbox keep-color v-model="directions" label="North" val="north" color="red" />
                        <q-checkbox keep-color v-model="directions" label="South" val="south" color="red" />
                        <q-checkbox keep-color v-model="directions" label="East" val="east" color="red" />
                        <q-checkbox keep-color v-model="directions" label="West" val="west" color="red" />
                        <q-checkbox keep-color v-model="directions" label="North East" val="northeast" color="red" />
                        <q-checkbox keep-color v-model="directions" label="South East" val="southeast" color="red" />
                        <q-checkbox keep-color v-model="directions" label="North West" val="northwest" color="red" />
                        <q-checkbox keep-color v-model="directions" label="South West" val="southwest" color="red" />
                    </div>
                </div>
                
            </div>
            
        </q-card>
    </q-dialog>
</template>

<script lang="ts">
import { EditorState,MPEditorState,MoveType, MovePattern } from '../../types';
import { convertFENtoBoardState } from '../../utils/fen';
import { PropType, reactive, ref, watch } from 'vue';
import Board from '../Board/Board.vue';

type DirOffsets = {
  [key: string]: [number, number];
};


export default {
    components:{Board},
    props:{
        editorState: { type:Object as PropType<EditorState>, required:true}
    },
    emits:["close-dialog","save-mp"],
    setup(props,{emit}){
        const showDialog = ref(true);
        const editorState : MPEditorState = reactive({
            ...props.editorState,
            editorType:"MP",
            moveType:"jump", 
            curCustomPiece: props.editorState.curCustomPiece!
        });
        const piece = editorState.curCustomPiece;
        const boardRef = ref();
        const directions = ref([]);
        let tempMovePattern:MovePattern = {piece,jumpPatterns:[],slidePatterns:[]} 
        let board = convertFENtoBoardState(`9/9/9/9/4${editorState && editorState.curPieceColor==="white" && piece ? piece.toUpperCase() : piece}4/9/9/9/9 w KQkq - 0 1`);
        const boardState = reactive(board);
        const center = {row:4, col:4};
        const dirOffsets:DirOffsets = { 'east':[0,1], 'west':[0,-1],'north':[-1,0],'south':[1,0],'northwest':[-1,-1],'southwest':[1,-1],'northeast':[-1,1],'southeast':[1,1]}
        
        watch(directions,(oldVal,newVal)=>{
            newVal.filter(val => !oldVal.includes(val)).forEach((dir)=> updateSlideMP(dir,'remove'));
            oldVal.filter(val => !newVal.includes(val)).forEach((dir)=> updateSlideMP(dir,'add'));
        })
        
        const updateJumpMP = (payload:{clickType:string,row:number,col:number})=>{
            let curcolor = boardState.squares[payload.row][payload.col].squareInfo.tempSquareColor
            if (payload.clickType === "set-jump-mp"){
                boardState.squares[payload.row][payload.col].squareInfo.tempSquareColor = curcolor==='slide' ? 'slide' : 'jump';
                tempMovePattern.jumpPatterns.push([center.row-payload.row,center.col-payload.col])
            } else if (payload.clickType ==='remove-jump-mp'){
                boardState.squares[payload.row][payload.col].squareInfo.tempSquareColor = null
                tempMovePattern.jumpPatterns = tempMovePattern.jumpPatterns.filter(pair => {
                    return !(pair[0] === center.row-payload.row && pair[1] === center.col-payload.col);
                });
            }
            boardRef.value.updateBoardState1D(boardState)
        }
        
        const updateSlideMP = (direction:string,action:string)=>{
            let [dx,dy] = dirOffsets[direction];
            let [curX,curY] = [center.row+dx, center.col+dy];
            while(curX>=0 && curX < boardState.dimensions.rows && curY>=0 && curY<boardState.dimensions.cols){
                boardState.squares[curX][curY].squareInfo.tempSquareColor = (action==='add') ? 'slide' : null;
                curX+=dx,
                curY+=dy
            }
            boardRef.value.updateBoardState1D(boardState)
        }
        const closeDialog = () => {
            showDialog.value = false; 
            emit("close-dialog")
        }
        return {
            showDialog,
            closeDialog ,
            boardState,
            editorState,
            updateJumpMP,
            saveMP:()=>{ emit('save-mp',tempMovePattern); closeDialog()},
            boardRef,
            directions
        }
    }
}
</script>

<style scoped>

.mp-card{
    display:flex;
    flex-direction: column;
    margin: 1%;
    padding: 2%;
}
.top-row{
    display:flex;
    justify-content: end;
    margin: 1em;
}

.flex-panels{
    flex-direction: row;
    display: flex;
}
.slide-checkboxes{
    display:flex;
    flex-direction:column
}
.options{
    display:flex;
    margin: 1%;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    
}
.flex-panels > *{ margin:1%}

</style>