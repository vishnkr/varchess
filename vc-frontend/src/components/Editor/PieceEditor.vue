<template>
    <div>
        <h6>
            Choose Color:
        </h6>
        <div>
            <q-radio v-model="editorState.curPieceColor" label="White" val="white" color="orange" keep-color />
            <q-radio v-model="editorState.curPieceColor" label="Black" val="black" color="orange" keep-color />
        </div>
        <h6>
            Choose Piece:
        </h6>
        <div class="piece-options">
                <q-option-group
                    v-model="editorState.curPiece"
                    :options="pieces"
                    keep-color
                    color="orange"
                />
                <div v-if="isCustomPiece" class="custom-piece-edit">
                <q-virtual-scroll 
                :items="customPieces"
                separator
                v-slot="{item,index}"
                style="max-height: 300px;"
                >          
                    <q-item :key="index" >
                        <div class="scroll-item" :class="[editorState.piecesInPlay[item]?.isMPDefined ? 'defined-mp' : null]">
                            <q-radio v-model="editorState.curCustomPiece" :val="item" color="orange" keep-color/>
                            <q-img :src="getCustomPiecePath(item)" width="4em" />
                        </div>
                    </q-item>

                </q-virtual-scroll>
                    <q-btn style="display: flex; align-items: center;" v-if="isCustomPiece && editorState.curCustomPiece"
                    color="green-9"
                    @click="dialogOpened = true;"
                    label="Set Move Pattern"
                    />
                </div>
                
                
                <div v-if="dialogOpened" style="max-width:700px;">
                    <MovePatternDialog 
                    @close-dialog="closeDialog"
                    @save-mp="saveMP"
                    :editorState="editorState" 
                    />
                </div>
            
        </div>
    </div>
</template>

<script lang="ts">
import { EditorState, MovePattern } from '@/types';
import MovePatternDialog from '../Editor/MovePatternDialog.vue'
import {computed, PropType, reactive, ref,watch} from 'vue'

export default {
    components:{MovePatternDialog},
    emits: ['update-piece-state','save-mp'],
    props:{
        editorState: {type: Object as PropType<EditorState>, required:true}
    },
    setup(props,{emit}){
        const customPieces = ['a','j','d','i','g','s','u','v','z'];
        const editorState = reactive(props.editorState)
        const dialogOpened = ref(false);
        watch(editorState,(newVal,oldVal)=>{
            emit('update-piece-state',newVal)
        });
        return {
            editorState,
            customPieces,
            pieces: [
                { label: 'Rook', value: 'r' },
                { label: 'Pawn', value: 'p' },
                { label: 'Queen', value: 'q'},
                { label: 'Bishop', value: 'b'},
                { label: 'King', value: 'k'},
                { label: 'Knight', value: 'n' },
                { label: 'Custom', value: 'c'}
            ],
            getCustomPiecePath: (pieceName:string)=>{
                return new URL(`../../assets/images/pieces/white/${pieceName}.svg`,import.meta.url).href
            },
            dialogOpened,
            closeDialog:()=>{
                editorState.curCustomPiece = null;
                dialogOpened.value=false
            },
            isCustomPiece:computed(()=>{
                return editorState.curPiece==='c' 
            }),
            saveMP:(movePattern:MovePattern)=>{
                if(movePattern.piece in editorState.piecesInPlay){
                    editorState.piecesInPlay[movePattern.piece].isMPDefined = true;
                } else { editorState.piecesInPlay[movePattern.piece] = {isMPDefined:true,isAddedToBoard:false}}
                emit('save-mp',movePattern)
            }
        }
    }
};
</script>

<style scoped>
.piece-options{
    display:flex;
    
}

.custom-piece-edit{
    display: flex;
    flex:1;
    padding: 1em;
    flex-direction: row;
    justify-content: center;
    width:50%;
    align-items: flex-start;
    overflow: scroll;
}
.custom-piece-edit> button{
    margin:1%;
}
.defined-mp{
    background-color: #2e7d32 !important;
}
h6{
    margin:0;
}

.scroll-container {
    display:flex;
    flex:1;
    max-height: 200px;
    overflow-y: scroll;
}
.scroll-container > .scroll-item{
    display:flex;
}
</style>