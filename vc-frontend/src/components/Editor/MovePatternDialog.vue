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
          <div>
          <v-toolbar
            dark
            dense
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
              <v-btn
                dark
                text
                @click="savePattern"
              >
                Save
              </v-btn>

            
          </v-toolbar>
          </div>
          <div class="flex-panels">
            <div style="padding:10px" class="board-panel">
                <board 
                :board="boardState"
                :editorMode="true"
                :editorData="editData"
                />
            </div>
          <div style="padding:10px" class="settings-panel">
                <v-card class="mx-auto" max-width="550">
                  <div style="padding:10px">
                    <v-list-item-title >Choose Move Type</v-list-item-title>
                    <v-radio-group
                      row
                      v-model="moveType"
                    >
                      <v-radio
                        label="Jump"
                        value="jump"
                      ></v-radio>
                      <v-radio
                        label="Slide"
                        value="slide"
                      ></v-radio>
                    </v-radio-group>
                </div>
                </v-card>
            </div>
          </div>
        </v-card>
      </v-dialog>
</template>

<script>
import Board from '../Board.vue'
export default {
    components:{Board},
    props:['dialog','pieceType','pieceColor','editorData'],
    watch:{
      moveType(){
        console.log('got',this)
        this.editData.moveType = this.moveType
      }
    },
    computed:{
    
    },
    created(){
      this.setupMPBoard()
      this.editData.moveType = this.moveType
      this.editData.isSetMovement = true
    },
    methods:{
        closeDialog(){
            this.$emit("closeDialog")
            
        },
        savePattern(){
          //must save move pattern definitions later and then close
          this.closeDialog()
        },
        isEven(val){return val%2==0},
        isLight(row,col){
          return this.isEven(row)&&this.isEven(col)|| (!this.isEven(row)&&!this.isEven(col))},
        setupMPBoard(){
          for(var col =0;col<this.cols;col++){
            this.boardState.tiles.push([])
            for(var row=0;row<this.rows;row++){
              var tile = {}
              tile.tileType = this.isLight(col,row)? 'l' : 'd';
              if(row==4 && col==4){
                tile.isPiecePresent=true
                tile.pieceType = this.editorData.customPiece
                tile.pieceColor = this.editorData.curPieceColor
              } else{ tile.isPiecePresent=false }
              this.boardState.tiles[col].push(tile)
            }
            
          }
          this.boardState.rows = this.rows
          this.boardState.cols = this.cols
        }
    },
    data(){
        return{
            rows:9,
            cols:9,
            boardState: {tiles:[]},
            moveType: 'jump',
            editData: {...this.editorData}
        }
    }
}
</script>

<style scoped>
.flex-panels{
  display: flex;
}
.board-panel{
  flex:1;
}

.settings-panel{
  flex:1;
}
</style>