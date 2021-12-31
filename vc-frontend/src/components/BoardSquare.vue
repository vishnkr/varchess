<template>
  <div class="square" :id="tileId" :style="cssVar" 
        :class="[tileType=='d'? 'dark':'light', mpTabData? mpTabData[[this.row,this.col]]:null , !mpTabData && isPiecePresent && isHighlighted?'highlight-from':null,
        addColor==='jump' ? 'move-jump-pattern'  : addColor==='slide' ? 'move-slide-pattern': null, ]"
         @mousedown="clickSquare">
      <div v-if="isPiecePresent" >
      <board-piece  :color="pieceColor" :pieceType="pieceType" :row="row" :col="col"/>
      </div>
  </div>
</template>

<script>
import BoardPiece from './BoardPiece'
//import colorConstants from '../utils/colorConstants';
export default {
    components:{BoardPiece},
    mounted(){
      //this.squareColor = colorConstants[this.color.highlighted? this.color.highlighted : this.color.default];
    },
    methods:{
      addColorToSquare(moveType){
        this.addColor = moveType;
      },
      removeColorFromSquare(){
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
              // get possible move squares and highlight them
          } else { // dest pos is selected
            //stop displaying possible squares here
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
        squareColor:null,
      }
    },
    props:['tileType','editorMode','editorData','mpTabData','row','col','isPiecePresent','pieceType','pieceColor','x','y','tileId','isHighlighted','selectedSrc','color'],
    computed:{
        cssVar(){
        return {
        '--x': this.x,
        '--y': this.y,
        '--color':this.squareColor
        }
      },
       //to be used later in portal mode
      getPortal(){
        return require('../assets/images/exit.svg');
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
.dark {
    background-color: #b2c85d;
  }

.light {
    background-color: #e4f5cb;
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