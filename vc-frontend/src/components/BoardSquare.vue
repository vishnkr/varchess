<template>
  <div class="square" :id="tileId" :style="cssVar" 
        :class="[tileType=='d'? 'dark':'light', isPiecePresent && isHighlighted?'highlight-from':null,
        addColor==='jump' ? 'move-jump-pattern'  : addColor==='slide' ? 'move-slide-pattern': null, ]"
        @click="clickSquare">
      <div v-if="isPiecePresent" >
      <board-piece  :color="pieceColor" :pieceType="pieceType" :row="row" :col="col"/>
      </div>

      
  </div>
</template>

<script>
import BoardPiece from './BoardPiece'
export default {
    components:{BoardPiece},
    methods:{
      addColorToSquare(moveType){
        this.addColor = moveType
      },
      removeColorFromSquare(){
        console.log('col',this.addColor)
        this.addColor= this.addColor==='jump'? 'jump' : null;
        },
      clickSquare(){
        var pieceType;
        if(this.editorMode){
          var clickType  =  this.editorData.isSetMovement ? "setPattern" : "regular"
          if (this.editorData.isSetMovement){ 
            this.addColor = this.addColor ? null : this.editorData.moveType
          }
          this.$emit("setEditorBoardState",clickType,this.row,this.col)
        } else {
          if(this.isPiecePresent & (!this.selectedSrc || this.selectedSrc && this.selectedSrc.pieceColor == (this.pieceColor=='white' ? 'w':'b') && this.selectedSrc.id!=this.tileId)){ // start pos is selected
              pieceType = this.pieceColor=='w' ? this.pieceType.toUpperCase() : this.pieceType.toLowerCase()
              this.$emit("sendSelectedPiece",{id:this.tileId,pieceType:pieceType,pieceColor:this.pieceColor,row:this.row,col:this.col})
          } else { // dest pos is selected
          if(this.$store.state.curStartPos.row == this.row && this.$store.state.curStartPos.col == this.col){ //clicking same piece as destination
              this.$store.commit('undoSelection')
              this.$emit("sendSelectedPiece",null)
            } else{
              if(this.isPiecePresent){
                this.$emit("destinationSelect",{id:this.tileId,isPiecePresent:true,pieceColor:this.pieceColor,pieceType:this.pieceType,row:this.row,col:this.col})
              } else{
                this.$emit("destinationSelect",{id:this.tileId,isPiecePresent:false,row:this.row,col:this.col})
              }
            }
        }
        }
      }
    },
    data(){
      return {
        addColor:null,
      }
    },
    props:['tileType','editorMode','editorData','row','col','isPiecePresent','pieceType','pieceColor','x','y','tileId','isHighlighted','selectedSrc'],
    computed:{
        cssVar(){
        return {
        '--x': this.x,
        '--y': this.y,
        }
      }
    }
}
</script>

<style>


.square {
  background: transparent;
  border: 1px solid transparent;
   width: 100%;
    height: 0;
    padding-bottom: 100%;
    grid-column: var(--y);
    grid-row: var(--x);
}

.dark {
  background-color: #b2c85d;
}

.light {
  background-color: #e4f5cb;
}

.highlight-possible{
  background-color: #d9bf77 !important;
}

.highlight-from{
  background-color: #a97d5d !important;
}

.move-jump-pattern{
  background-color: #4056b8 !important;
  border-color: black;
}

.move-slide-pattern{
  background-color: #ac422a !important;
  border-color: black;
}

@media only screen and (max-device-width: 480px) {
  .square {
  background: transparent;
  border: 1px solid transparent;
  float: left;
  font-size: 6px;
  font-weight: bold;
  line-height: 34px;
  height: 12px;/*48px;*/
  margin-right: -1px;
  margin-top: -1px;
  padding: 0;
  text-align: center;
  width: 12px;
  }
  .dark {
    background-color: #b2c85d;
  }

  .light {
    background-color: #e4f5cb;
  }
}
</style>