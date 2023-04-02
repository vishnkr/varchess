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
                ref="MPBoard" 
                :board="boardState"
                :editorMode="true"
                :editorState="editData"
                v-on:setMP="addPattern"
                />
            </div>
          <div style="padding:10px" class="settings-panel">
                <v-card class="mx-auto" max-width="550">
                  <div style="padding:10px;">
                    <v-list-item-title >Choose Move Type</v-list-item-title>
                    <div style="display:flex;">
                    <span style="margin-right:10px;"><div class='box move-jump' /> Jump</span>
                    <span style="margin-right:10px;"><div class='box move-slide' /> Slide</span>
                    </div>
                    <v-radio-group
                      row
                      v-model="moveType"
                    >
                      <v-radio
                        label="Jump"
                        value="jump"
                      />
                      <v-radio
                        label="Slide"
                        value="slide" 
                        />
                    </v-radio-group>
                    <v-divider></v-divider>
                    <div v-if="moveType=='slide'">
                      <v-checkbox v-model="slideDirections" v-for="(val,key) in directions" :key="key" :value="val" :label="key" />
                    </div>

                </div>
                </v-card>
            </div>
          </div>
        </v-card>
      </v-dialog>
</template>

<script>
import Board from '../Board/Board.vue'
export default {
    components:{Board},
    props:['dialog','pieceType','pieceColor','editorState'],
    watch:{
      moveType(){
        this.editData.moveType = this.moveType
      },
      slideDirections(){
        this.$refs.MPBoard.setSlidePattern(this.slideDirections)
      }
    },
    computed:{
    
    },
    created(){
      this.setupMPBoard()
      this.editData.moveType = this.moveType
      this.editData.isSetMovement = true
      this.editData.piecePos = this.piecePos
    },
    methods:{
        addPattern(row,col){
          if(this.moveType=="jump"){
            if(this.piecePos[0]!=row-1 || this.piecePos[1]!=col-1){
              this.jumpPattern.push([row-this.piecePos[0]-1,col-this.piecePos[1]-1])
            }
          }
        },
        closeDialog(){
            this.$emit("close-dialog")
        },
        savePattern(){
          this.$emit("emit-move-pattern",this.editorState.customPiece,this.jumpPattern,this.slideDirections)
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
                tile.pieceType = this.editorState.customPiece
                tile.pieceColor = this.editorState.curPieceColor
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
            piecePos:[4,4],
            boardState: {tiles:[]},
            moveType: 'jump',
            slideDirections:[],
            directions: {'North':[-1,0],'South':[1,0],'East':[0,1],'West':[0,-1],'North East':[-1,1], 'North West':[-1,-1], 'South East':[1,1], 'South West':[1,-1]},
            editData: {...this.editorState},
            jumpPattern: [],
        }
    },
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

.box {
  float: left;
  height: 20px;
  width: 20px;
  margin-right:10px;
  margin-bottom: 15px;
  border: 1px solid black;
  clear: both;
}

.move-jump {
  background-color: #4056b8;
}

.move-slide {
  background-color: #ac422a;
}

</style>