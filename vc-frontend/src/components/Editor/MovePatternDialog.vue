<template>
  <v-dialog
        v-model="dialog"
        fullscreen
        hide-overlay
        transition="dialog-bottom-transition"
        scrollable
        :retain-focus="false"
      >
        <v-card tile>
          <v-toolbar
            dense
            dark
          >
            <v-btn
              icon
              dark
              @click="closeDialog"
            >
              <v-icon>mdi-close</v-icon>
            </v-btn>
            <v-toolbar-title>Set Movement Pattern</v-toolbar-title>
            <v-spacer></v-spacer>
            <v-toolbar-items>
              <v-btn
                dark
                text
                @click="closeDialog"
              >
                Save
              </v-btn>
            </v-toolbar-items>
            
          </v-toolbar>
            <div class="board-panel">
                <editor-board 
                :isEmpty="true"
                :curPiece="pieceType"
                :curPieceColor="pieceColor"
                :cols="9" 
                :rows="9" 
                />
            </div>
          <div style="flex: 5 1 auto;"></div>
          <div class="settings-panel">
                <v-card>
                    hello
                </v-card>
            </div>
        </v-card>
      </v-dialog>
</template>

<script>
import EditorBoard from './EditorBoard.vue'
export default {
    components:{EditorBoard},
    props:['dialog','pieceType','pieceColor'],
    computed:{
        
    },
    methods:{
        closeDialog(){
            this.$emit("closeDialog")
        },
        setupMPBoard(){
        for(var col =0;col<this.cols;col++){
          this.boardState.tiles.push([])
          for(var row=0;row<this.rows;row++){
            var tile = {}
            tile.tileType = this.isLight(col,row)? 'l' : 'd';
            if(row==0){
              tile.isPiecePresent=false
            }
            }}}
    },
    data(){
        return{
            rows:9,
            cols:9,
        }
    }
}
</script>

<style>

</style>