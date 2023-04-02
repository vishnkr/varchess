<template>
    <q-page class="parent-background">
        <q-dialog v-model="isLoading">
            <LoadingScreen :shareLink="getShareLink" @update-loading="closeRoom"/>
        </q-dialog>
        <div class="board-editor">
            <q-card dark class="card q-pa-md">
                <q-card-section class="top-row">
                    <div class="text-h5" align="right"> 
                        Game Editor
                    </div>
                    <div class="top-btns">
                        <q-btn color="negative" style="margin-right:5px;" label="Clear Board"></q-btn>
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
                            <PieceEditor @update-piece-state="updatePieceState"/>
                        </q-tab-panel>
                    </q-tab-panels>
                </q-card-section>
            </q-card>
        </div>
        <Board 
            :isFlipped="true" 
            :boardState="boardState" 
            :editorState="editorState" 
            @handle-square-click="handleSquareClick"
            ref="boardRef"
        />
    </q-page>
</template>

<script lang="ts">
import { createDefaultMaxBoardStateSquares, isLight } from '../utils';
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
        const username = router.currentRoute.value.params.username;
        const roomId = ref(null);
        const editorState: EditorState = reactive({
            curPiece: 'p',
            curPieceColor: 'white',
            isDisableTileOn: false,
            piecesInPlay:{},
            editorType: 'Game'
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
                    let firstRowSquares = maxBoardState.squares[0].slice(0,afterLastCol)
                    tempSquares = maxBoardState.squares.slice(1,afterLastRow).map((row:Square[])=> row.slice(0,afterLastCol))
                    maxBoardState.squares = [...tempSquares,...[firstRowSquares],...maxBoardState.squares.slice(afterLastRow)]
                    break;
                case 'down':
                    let lastRowSquares = maxBoardState.squares[lastRow];
                    tempSquares = maxBoardState.squares.slice(0,lastRow).map((row:Square[])=> row.slice(0,afterLastCol))
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

        const updatePieceState = ({curPiece,curPieceColor}:{curPiece:string,curPieceColor:PieceColor})=>{
            editorState.curPiece = curPiece;
            editorState.curPieceColor = curPieceColor;
        }

        const setEditorDisable = (value:boolean)=>{ editorState.isDisableTileOn = value }

        const togglePieceOnSquare = (squareInfo:{row:number,col:number})=>{
            let rowIndex = squareInfo.row-1;
            let colIndex = squareInfo.col-1; 
            if (!maxBoardState.squares[rowIndex][colIndex].squareInfo.isPiecePresent){
                maxBoardState.squares[rowIndex][colIndex].squareInfo.isPiecePresent = true
                maxBoardState.squares[rowIndex][colIndex].squareInfo.pieceColor = editorState.curPieceColor
                maxBoardState.squares[rowIndex][colIndex].squareInfo.pieceType = editorState.curPiece
                
            } else{
                maxBoardState.squares[rowIndex][colIndex].squareInfo.isPiecePresent = false
            }
        }
        const setDisable = (squareInfo:{row:number,col:number})=>{
            let [rowIndex, colIndex] = [squareInfo.row-1, squareInfo.col-1];
            boardState.squares[rowIndex][colIndex].squareInfo.isPiecePresent = false
            boardState.squares[rowIndex][colIndex].disabled = !boardState.squares[rowIndex][colIndex].disabled;
        }
            
        const handleSquareClick = (payload:{clickType:string,squareInfo:{row:number,col:number}})=>{
            switch (payload.clickType){
                case 'toggle-piece':
                    togglePieceOnSquare(payload.squareInfo)
                    break;
                case 'disable':
                    setDisable(payload.squareInfo);
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
            shiftBoard
        }
    },
    
});

</script>

<style scoped>
.card{
    display:flex;
    justify-content: center;
    flex-direction: column;
}

@media screen and (max-width:990px) {
  .parent-background { flex-wrap: wrap;  }
  .parent-background:first-child { flex-basis: 100%; }
}

.parent-background{
    display: flex;
    flex-direction: row;
    margin: 1em;
}
.board-editor{
    flex:1;
    margin-right: 1em;
    margin-bottom: 1em;
}
.top-row{
    display:flex;
    flex-direction: row;
}
.top-row>.top-btns{
    display:flex;
    justify-content: end;
    flex:3;
}

.top-row>div{
    flex:1;
}

</style>