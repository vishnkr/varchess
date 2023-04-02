<template>
    <div class="q-pa-md q-gutter-lg">
        <div class="text-h6">
            Choose Color:
        </div>
        <div class="q-gutter-md">
            <q-radio v-model="editorState.curPieceColor" label="White" val="white" color="orange" keep-color />
            <q-radio v-model="editorState.curPieceColor" label="Black" val="black" color="orange" keep-color />
        </div>
        <div class="text-h6">
            Choose Piece:
        </div>
        <div class="piece-options">
            <div class="q-gutter-md">
                <q-option-group
                    v-model="editorState.curPiece"
                    :options="pieces"
                    keep-color
                    color="orange"
                />
            </div>
            <div v-if="editorState.curPiece==='c'">
                <q-btn 
                    color="orange"
                    @click="mpDialog = true;"
                />
            </div>
        </div>
    </div>
</template>

<script lang="ts">
import { EditorState } from '@/types';
import {computed, reactive, ref,watch} from 'vue'

export default {
    emits: ['update-piece-state'],
    setup(props,{emit}){
        const customPieces = ['a','j','d','i','g','s','u','v','z'];
        const editorState = reactive({
            curPiece: 'p',
            curPieceColor: 'white',
        })
        const mpDialog = ref(false);
        watch(editorState,(newVal,oldVal)=>{
            emit('update-piece-state',newVal)
        });

        return {
            editorState,
            pieces: [
                { label: 'Rook', value: 'r' },
                { label: 'Pawn', value: 'p' },
                { label: 'Queen', value: 'q'},
                { label: 'Bishop', value: 'b'},
                { label: 'King', value: 'k'},
                { label: 'Knight', value: 'n' },
                { label: 'Custom', value: 'c'}
            ],
            getCustomPiecePaths: computed(()=>{
                return customPieces.map((pieceName)=>{
                        return new URL(`../../assets/images/pieces/${editorState.curPieceColor}/${pieceName}.svg`,import.meta.url).href
                })
            }),
            mpDialog
        }
    }
};
</script>

<style>
.piece-options{
    display:flex;
}
</style>