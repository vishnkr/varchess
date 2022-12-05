<template>
  <div class="square" :id="tileId" :style="cssVar" 
        :class="[disabled || tileType=='disabled' ? 'disabled' : tileType=='d'? 'dark':'light', 
        mpTabData? mpTabData[[this.row,this.col]]:null , 
        !mpTabData && highlight? this.getHighlightType():null,
        addColor==='jump' ? 'move-jump-pattern'  : addColor==='slide' ? 'move-slide-pattern': null]"

         @mousedown="clickSquare">
      <div v-if="isPiecePresent" >
      <piece  :color="pieceColor" :pieceType="pieceType" :row="row" :col="col"/>
      </div>
  </div>
</template>

<script>
import Piece from './Piece'

export default {
    components:{Piece},
    methods:{
      getHighlightType(){
        if (this.selectedSrc){

          if (this.highlight.from==this.tileId){
            return 'highlight-from'
          } else{
              for (var square of this.highlight.to){
                if (square[0]+1==this.row && square[1]+1==this.col){

                  return 'highlight-to'
                }
              }
          }
        }
      },
      addColorToSquare(moveType){
        this.addColor = moveType;
      },
      removeColorFromSquare(){
        this.addColor= this.addColor==='jump'? 'jump' : null;
        },
      clickSquare(){
        let pieceInfo = {id:this.tileId,row:this.row,col:this.col}
        if(this.editorMode){
          if (this.editorState.isDisableTileOn){
            this.disabled = !this.disabled
          }
          var clickType  =  this.editorState.isSetMovement ? "setPattern" : "regular"
          if (this.editorState.isSetMovement){ 
            this.addColor = this.addColor ? null : this.editorState.moveType
          }
          this.$emit("setEditorBoardState",clickType,this.row,this.col)
        } else {
          if(this.isPiecePresent & (!this.selectedSrc || this.selectedSrc && this.selectedSrc.pieceColor == (this.pieceColor=='white' ? 'w':'b') && this.selectedSrc.id!=this.tileId)){ // start pos is selected
              pieceInfo.pieceType = this.pieceColor=='w' ? this.pieceType.toUpperCase() : this.pieceType.toLowerCase();
              this.$emit("sendSelectedPiece",{...pieceInfo,pieceColor:this.pieceColor})
              // get possible move squares and highlight them
          } else { // dest pos is selected
            //stop displaying possible squares here
          if(this.$store.state.curStartPos.row == this.row && this.$store.state.curStartPos.col == this.col){ //clicking same piece as destination
              this.$store.commit('undoSelection')
              this.$emit("sendSelectedPiece",null)
            } else{
              if(this.isPiecePresent){
                this.$emit("destinationSelect",{...pieceInfo,isPiecePresent:true,pieceColor:this.pieceColor,pieceType:this.pieceType})
              } else{
                this.$emit("destinationSelect",{...pieceInfo,isPiecePresent:false})
              }
            }
        }
        }
      }
    },
    data(){
      return {
        addColor:null,
        squareColor:null,
        disabled:false
      }
    },
    props:['tileType','editorMode','editorState','mpTabData','row','col','isPiecePresent','pieceType','pieceColor','x','y','tileId','highlight','selectedSrc'],
    computed:{
        cssVar(){
        return {
        '--x': this.x,
        '--y': this.y,
        //'--color':this.squareColor
        }
      },
      
    }
}
</script>

<style scoped>
.square {
  background: transparent;
  border: 1px solid transparent;
  width: 100%;
  height: 0;
  padding-bottom: 100%;
  grid-column: var(--y);
  grid-row: var(--x);
  cursor: pointer;
  background-color: var(--color);
}

.highlight-to{
  background: rgba(217,191,119, 0.6) !important;
  border: 0.2px solid;
}

.highlight-from{
  background-color: #a97d5d !important;
}

.move-jump-pattern{
  background-color: #4056b8 !important;
  border-color: black;
}
.dark {
    background-color: #b2c85d;
  }

.light {
    background-color: #e4f5cb;
}

.disabled{
  background-color: dimgray !important;
}
.move-slide-pattern{
  background-color: #ac422a !important;
  border-color: black;
}
/*
.portal{
  animation-name: zoom;
  animation-iteration-count: infinite;
  animation-duration: 3000ms;
}
<div style="overflow:hidden;" v-else>
        <div class="portal-circle">
          <svg class="portal" height="70%" width="80%">
            <ellipse cx="40%" cy="40%" rx="15%" ry="15%" 
            style="fill:yellow;stroke:purple;stroke-width:2" />
          </svg>
        </div>
      </div>
@keyframes zoom{
  0% {
  transform: scale(1, 1);
  }
  50% {
  transform: scale(1.5, 1.5);
  }
  100% {
  transform: scale(1, 1);
  }
}
@keyframes spin {
    from {
        transform:rotate(0deg);
    }
    to {
        transform:rotate(360deg);
    }
}
*/
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