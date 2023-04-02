<template>
    <div>
        <div>
            <v-list-item-title style="padding-top:10px">Choose color</v-list-item-title>
                <v-radio-group
                v-model="modifiedEditorState.curPieceColor"
                row
                >
                    <v-radio
                        label="Black"
                        value="black"
                    ></v-radio>
                    <v-radio
                        label="White"
                        value="white"
                    ></v-radio>
                    </v-radio-group>
        </div>
        <div>
            <v-list-item-title style="padding-top:10px">Choose piece</v-list-item-title>
            <v-radio-group
            v-model="modifiedEditorState.curPiece"
            column
            >
            <v-radio v-for="piece in pieceList"
                :key="`${piece.toLowerCase()}`"
                :label="`${piece}`"
                :value="`${pieceMap[piece].toLowerCase()}`"
            ></v-radio>
                    
            </v-radio-group>
        </div>
        <div v-if="modifiedEditorState.curPiece=='c'">
            <div class="custom-pieces">
                <v-btn
                depressed
                color="primary"
                @click="dialog=true"
                class="move-button"
                >
                    Set Move Pattern
                </v-btn>
                <div class="piece-scroll">
                    <v-card
                    elevation="16"
                    max-width="150"
                    class="mx-auto"
                    >
                        <v-virtual-scroll
                                
                        :items="pieceFilter"
                        height="200"
                        item-height="64"
                        >
                            <template v-slot:default="{ item }">
                                <v-radio-group v-model="modifiedEditorState.customPiece">
                                    <v-list-item class="scroll-item" :class="{added : modifiedEditorState.added[item.piece], defined: modifiedEditorState.defined[item.piece]}">
                                        <v-radio
                                        :key="item.piece"
                                        :value="item.piece"
                                        @click="modifiedEditorState.customPiece = item.piece"
                                        >
                                            Select
                                        </v-radio>
                                        <v-list-item-content>
                                            <img class="resize" :src="item.src">
                                        </v-list-item-content>

                                        <v-list-item-action>
                                            <move-pattern-dialog v-if="dialog" 
                                            @close-dialog="closeDialog" 
                                            @emit-move-pattern="emitMovePattern"
                                            :dialog="dialog"
                                            :editorState="modifiedEditorState"
                                            :pieceColor="modifiedEditorState.curPieceColor"
                                            :pieceType="modifiedEditorState.customPiece"
                                            />
                                        </v-list-item-action>
                                        </v-list-item>
                                        <v-divider></v-divider>
                                    </v-radio-group>
                            </template>
                        </v-virtual-scroll>
                    </v-card>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import  MovePatternDialog from './MovePatternDialog.vue';
export default {
    props:['editorState'],
    methods:{
        emitUpdatePieceState(){
            this.$emit('update-piece-state',this.modifiedEditorState)
        },
        getPiecePath(piece){
            return require(`../../assets/images/pieces/${this.colorSelect}/${piece}.svg`)
        },
        closeDialog(){
            this.dialog=false
        },
        emitMovePattern(pieceName,jumpPattern,slidePattern){
            this.modifiedEditorState.defined[pieceName] = true;
            this.$emit("set-move-pattern",pieceName,jumpPattern,slidePattern)
        },
    },
    watch:{
        'modifiedEditorState.curPiece': function(){this.emitUpdatePieceState()},
        'modifiedEditorState.curPieceColor': function(){this.emitUpdatePieceState()},
    },
    computed:{
        pieceFilter(){
            let piecePaths = []
            for(let pieceName of this.customPieces){
            piecePaths.push({piece:pieceName,src :this.getPiecePath(pieceName)})
            }        
            return piecePaths
        },
    },
    data(){
        return {
            modifiedEditorState: this.editorState,
            colorSelect: 'white',
            pieceMap: {'Pawn':'p','King':'k','Queen':'q','Bishop':'b','Knight':'n','Rook':'r','Custom':'c'},
            customPieces:['a','j','d','i','g','s','u','v','z'],
            pieceList: ['Pawn','King','Queen','Bishop','Knight','Rook','Custom'],
            dialog:false,
        }
    },
    components: {MovePatternDialog},
}
</script>

<style>

</style>