<template>
    <q-page class="editor-page bg-dark">
        <q-dialog v-model="isLoading" v-if="roomId">
            <LoadingScreen :roomId="roomId" :username="username" :shareLink="getShareLink" @update-loading="closeRoom"/>
        </q-dialog>
        <div class="column board-editor">
            <q-card dark class="card q-pa-md">
                <q-card-section class="top-row">
                    <p class="text-h5"> 
                        Game Editor
                    </p>
                    <div class="top-btns">
                        <q-btn color="negative" style="margin-right:5px;" label="Clear Board" @click="clearBoard"></q-btn>
                        <q-btn class="bg-green-9" label="Create Room" @click="enterRoomWithLoading"></q-btn>
                    </div>
                </q-card-section>
                <q-card-section>
                    <q-tabs v-model="tab" dense align="justify" :breakpoint="0">
                        <q-tab name="board-editor" label="Board Editor"></q-tab>
                        <q-tab name="piece-editor" label="Piece Editor"></q-tab>
                    </q-tabs>
                </q-card-section>
                <q-card-section >
                    <q-tab-panels dark v-model="tab" animated>
                        <q-tab-panel name="board-editor">
                            <BoardEditor 
                                @update-board-dimensions="updateBoardDimensions" 
                                @toggle-disable="setEditorDisable"
                                @shift-board-direction="shiftBoard"
                                :isLoggedIn="true"
                                :isDisableTileOn="editorState.isDisableTileOn"
                                :rows="boardState.dimensions.rows"
                                :cols="boardState.dimensions.cols"
                            />
                        </q-tab-panel>
                        <q-tab-panel name="piece-editor">
                            <PieceEditor :editorState="editorState" @update-piece-state="updatePieceState"/>
                        </q-tab-panel>
                    </q-tab-panels>
                </q-card-section>
            </q-card>
        </div>
        <div class="board-panel column bg-dark">
            <Board 
            :isFlipped="false" 
            :boardState="boardState" 
            :editorState="editorState" 
            @handle-square-click="handleSquareClick"
            ref="boardRef"
        />
        </div>
    </q-page>
</template>

<script lang="ts">
import { createDefaultMaxBoardStateSquares, setupEmptyMaxSizeBoard } from '../utils';
import { reactive, ref,Ref, defineComponent, computed, watch } from 'vue';
import BoardEditor from '../components/Editor/BoardEditor.vue'
import Board from '../components/Board/Board.vue'
import PieceEditor from '../components/Editor/PieceEditor.vue'
import { BoardState, EditorState, PieceColor, MovePattern, Square} from '@/types';
import { convertBoardStateToFEN, convertFENtoBoardState } from '../utils/fen';
import { STANDARD_FEN } from '../utils/constants';
import { useStore } from 'vuex';
import { RootState } from '../store/state';
import { UPDATE_BOARD_STATE,CREATE_ROOM, SET_MOVE_PATTERNS, CONNECT_WS, SET_SERVER_STATUS, CLOSE_WS } from '../utils/action_mutation_types';
import { validateStartSetup } from '../utils/validator';
import { useRouter } from 'vue-router';
import LoadingScreen from '../components/Other/LoadingScreen.vue';

type EditorType = 'board-editor' | 'piece-editor'

export default defineComponent({
    components:{
        BoardEditor,
        PieceEditor,
        Board,
        LoadingScreen
    },
    setup(){
        const tab : Ref<EditorType> = ref('board-editor');
        const store = useStore<RootState>();
        const isLoading = ref(false);
        const maxBoardState = reactive({squares:createDefaultMaxBoardStateSquares()});
        const boardState :BoardState = reactive(convertFENtoBoardState(STANDARD_FEN));
        const boardRef = ref();
        const router = useRouter();
        const username = router.currentRoute.value.params.username.toString();
        const roomId = ref(null);
        const editorState: EditorState = reactive({
            curPiece: 'p',
            curPieceColor: 'white',
            isDisableTileOn: false,
            piecesInPlay:{},
            editorType: 'Game',
            curCustomPiece: null
        })

        watch(maxBoardState,(newmaxBoardState)=>{
            updateBoardStateFromMaxState()
        })

        const enterRoomWithLoading = async ()=>{
            let fen = convertBoardStateToFEN(boardState);
            let movePatterns = getMovePatterns();
            
            try{
                roomId.value = await store.dispatch(CREATE_ROOM,{fen,movePatterns})
                if(roomId.value){
                    if(validateStartSetup(boardState)){
                        isLoading.value = true;
                        if(movePatterns){
                            store.commit(SET_MOVE_PATTERNS,movePatterns)
                        }
                        store.commit(UPDATE_BOARD_STATE,{roomId:roomId.value,boardState })
                        store.dispatch(CONNECT_WS,{roomId:roomId.value,username})
                        store.commit(SET_SERVER_STATUS,{isOnline:true,errorMessage:null})
                    } else{
                        store.commit(SET_SERVER_STATUS,{isOnline:true,errorMessage:'Board state not valid: must contain 1 king for each color & not under check'})
                    }
                }
            } catch(error){
                console.error(error)
            }
        }

        function getMovePatterns(){
            return Object.entries(editorState.piecesInPlay).map(([piece,pieceData])=>{
                if(pieceData.isAddedToBoard){
                    return pieceData.movePattern
                }
            })
        }

        function shiftBoard(direction:string){
            let [lastCol,lastRow,afterLastCol,afterLastRow]  = [boardState.dimensions.cols-1,boardState.dimensions.rows-1,boardState.dimensions.cols,boardState.dimensions.rows];
            let tempSquares:Square[][];
            switch (direction){
                case 'right':
                    maxBoardState.squares = maxBoardState.squares.map((row,i) => {
                        if(i < boardState.dimensions.rows){
                            return [...row.slice(lastCol,afterLastCol),...row.slice(0,lastCol),...row.slice(afterLastCol)]
                        } 
                        return row
                    });
                    break;
                case 'left':
                    maxBoardState.squares = maxBoardState.squares.map((row:Square[],i) => {
                        let firstSquare = row[0];
                        if(i<boardState.dimensions.rows){
                            return [...row.slice(1,afterLastCol),...[firstSquare],...row.slice(afterLastCol)]
                        } 
                        return row
                    });
                    break;
                case 'up':
                    let firstRowSquares = maxBoardState.squares[0]
                    tempSquares = maxBoardState.squares.slice(1,afterLastRow)
                    maxBoardState.squares = [...tempSquares,...[firstRowSquares],...maxBoardState.squares.slice(afterLastRow)]
                    break;
                case 'down':
                    let lastRowSquares = maxBoardState.squares[lastRow];
                    tempSquares = maxBoardState.squares.slice(0,lastRow)
                    maxBoardState.squares = [...[lastRowSquares],...tempSquares,...maxBoardState.squares.slice(afterLastRow)]
                    break;
            }
        }

        function updateBoardDimensions(dimensions:{rows:number,cols:number}){
            boardState.dimensions = dimensions
            updateBoardStateFromMaxState()
        }

        function updateBoardStateFromMaxState(){
            let newSquares:Square[][] = [];
            for(let row=0;row<boardState.dimensions.rows;row++){
                newSquares.push(maxBoardState.squares[row].slice(0,boardState.dimensions.cols));
            }
            boardState.squares = newSquares;
            boardRef.value.updateBoardState1D(boardState)
        }

        const updatePieceState = (newEditorState:EditorState)=>{
            editorState.curPiece = newEditorState.curPiece;
            editorState.curPieceColor = newEditorState.curPieceColor;
            editorState.curCustomPiece = newEditorState.curCustomPiece;
            editorState.piecesInPlay = newEditorState.piecesInPlay;
        }

        const setEditorDisable = (value:boolean)=>{ 
            editorState.isDisableTileOn = value 
        }

        const togglePieceOnSquare = (squareInfo:{row:number,col:number})=>{
            if (!maxBoardState.squares[squareInfo.row][squareInfo.col].squareInfo.isPiecePresent){
                maxBoardState.squares[squareInfo.row][squareInfo.col].squareInfo.isPiecePresent = true
                maxBoardState.squares[squareInfo.row][squareInfo.col].squareInfo.pieceColor = editorState.curPieceColor
                maxBoardState.squares[squareInfo.row][squareInfo.col].squareInfo.pieceType = editorState.curPiece === 'c' && editorState.curCustomPiece ? editorState.curCustomPiece : editorState.curPiece
                
            } else{
                maxBoardState.squares[squareInfo.row][squareInfo.col].squareInfo.isPiecePresent = false
            }
        }
        const setDisable = (squareInfo:{row:number,col:number})=>{
            maxBoardState.squares[squareInfo.row][squareInfo.col].squareInfo.isPiecePresent = false
            maxBoardState.squares[squareInfo.row][squareInfo.col].disabled = !maxBoardState.squares[squareInfo.row][squareInfo.col].disabled;
        }
            
        const handleSquareClick = (payload:{clickType:string,row:number,col:number})=>{
            switch (payload.clickType){
                case 'toggle-piece':
                    togglePieceOnSquare({row:payload.row,col:payload.col})
                    break;
                case 'disable':
                    setDisable({row:payload.row,col:payload.col});
                    break;
            }
            boardRef.value.updateBoardState1D(boardState)
        }

        const getShareLink = computed(()=>{ return `${window.location.origin}/join/${roomId.value}` })
        
        const closeRoom = () =>{
            store.dispatch(CLOSE_WS)
            isLoading.value = false;
        }

        return {
            tab,
            boardState,
            roomId,
            username,
            updateBoardDimensions,
            boardRef,
            editorState,
            isLoading,
            updatePieceState,
            handleSquareClick,
            setEditorDisable,
            enterRoomWithLoading,
            getShareLink,
            closeRoom,
            shiftBoard,
            clearBoard: ()=>{maxBoardState.squares = setupEmptyMaxSizeBoard()}
        }
    },
    
});

</script>

<style scoped>
@media (min-width:320px)and(min-width:641px)  { /* smartphones, iPhone, portrait 480x320 phones */
    .editor-page{
    display: flex;
    flex-direction: column;
    max-height: 100%;
    max-width: 100%;
    align-items: center;
    }
    .board-panel{
        width:100%;
    }
}
/*
@media (min-width:481px)  {  }
@media (min-width:641px)  { 
    .editor-page{
    display: flex;
    flex-direction: column;
    max-height: 100%;
    max-width: 100%;
    align-items: center;
    }
@media (min-width:961px)  { /* tablet, landscape iPad, lo-res laptops ands desktops  }
@media (min-width:1025px) { /* big landscape tablets, laptops, and desktops  }
@media (min-width:1281px) { /* hi-res laptops and desktops  }*/
.editor-page{
    display: flex;
    flex-direction: row;
    max-height: 100%;
    max-width: 100%;
    align-items: center;
    padding-bottom: 2%;
}
.board-editor{
    width:100%;
    padding: 1%;
}

.board-panel{
    max-width: 800px;
    width:100%;
    padding:1%;
}
.board-editor> .card{
    display:flex;
    justify-content: center;
    flex-direction: column;
    flex:1;
}

.top-row{
    display:flex;
    flex-direction: row;
}
.top-row>.top-btns{
    display:flex;
    justify-content: end;
    padding: 1em;
}

.top-row>div{
    flex:1;
}

</style>