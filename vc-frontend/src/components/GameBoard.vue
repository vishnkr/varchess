<template>
  <div :style="cssVar" id='board-container'>
      <div id="board">
        <board-square v-for="square in this.boardState1D"  :key="square.tileId" 
        :tileType="square.tileType" 
        :isPiecePresent="square.isPiecePresent" 
        :tileId="square.tileId"
        :editorMode="editorMode"
        :editorData="editorData? editorData: null"
        :pieceType="square.pieceType" 
        :pieceColor="square.pieceColor" 
        :x="square.x" :y="square.y" 
        :row="square.row" :col="square.col"  
        :isHighlighted="selectedSrc && square.tileId==selectedSrc.id && selectedSrc.pieceColor==playerColor"
        :isSelectedSrc="selectedSrc?true:false"
        v-on:sendSelectedPiece="setSelectedPiece"
        v-on:destinationSelect="emitDestinationSelect"
        v-on:setEditorBoardState="editorModeSquareClicked"
        />
        
      </div>
  </div>
</template>

<script>
import BoardSquare from './BoardSquare.vue';
import {convertBoardStateToFEN} from '../utils/fen'
export default {
  components: { BoardSquare },
  props:['board','isflipped','playerColor',"editorMode","editorData"],
  watch: { 
    isflipped() { // watch it
          this.updateBoardState1D(this.isflipped)
        }
  },
  mounted(){
      this.boardState = this.board;
      this.rows = this.board.rows
      this.cols = this.board.cols
      this.updateBoardState1D(this.isflipped)
      //console.log("tileS",this.board.tiles,this.rows,this.cols,this.boardState1D)
      convertBoardStateToFEN(this.boardState,'w','KQkq','-')
  },
  data(){
        return {
            boardState: [],
            boardState1D: [],
            rows: 0,
            cols:0,
            selectedSrc: null,
        }
    },
    methods:{
        editorModeSquareClicked(row,col){
        if(this.boardState.tiles[row-1][col-1].isPiecePresent){
          this.boardState.tiles[row-1][col-1].isPiecePresent = false;
        }
        else{
          this.boardState.tiles[row-1][col-1].isPiecePresent = true;
          this.boardState.tiles[row-1][col-1].pieceColor = this.editorData.curPieceColor;
          this.boardState.tiles[row-1][col-1].pieceType = this.editorData.curPiece =='c' ? this.editorData.customPiece : this.editorData.curPiece;
        }
        this.$emit("sendEditorBoardState",this.boardState)
        this.updateBoardState1D(this.isflipped)
        },
        emitDestinationSelect(destInfo){
          this.$emit('destinationSelect',destInfo)
        },
        
        //changes UI of the board to match new boardstate after a valid move is made
        performMove(moveInfo){
          this.boardState.tiles[moveInfo.destRow][moveInfo.destCol] = {isPiecePresent:true, pieceType:moveInfo.piece.toLowerCase(),pieceColor:moveInfo.piece === moveInfo.piece.toUpperCase()?'white' :'black'}
          this.board.tiles[moveInfo.srcRow][moveInfo.srcCol]= {isPiecePresent:false, pieceType:null,pieceColor:null}
          if(moveInfo.castle){
            console.log('castling')
            var newRookPos,oldRookPos
            if(moveInfo.destCol<moveInfo.srcCol){
              oldRookPos = 0;
              newRookPos = moveInfo.srcCol-1
            } else {
              oldRookPos = this.board.tiles[0].length -1;
              newRookPos = moveInfo.srcCol+1
            }
            this.board.tiles[moveInfo.srcRow][newRookPos].isPiecePresent = true
            this.boardState.tiles[moveInfo.destRow][newRookPos].pieceType = 'r'
            this.boardState.tiles[moveInfo.destRow][newRookPos].pieceColor = moveInfo.piece === moveInfo.piece.toUpperCase()?'white' :'black'
            this.board.tiles[moveInfo.srcRow][oldRookPos].pieceColor = null
            this.board.tiles[moveInfo.srcRow][oldRookPos].pieceType = null
            this.board.tiles[moveInfo.srcRow][oldRookPos].isPiecePresent = false
          }
          this.updateBoardState1D(this.isflipped)
          this.selectedSrc = null
        },
        setSelectedPiece(pieceInfo){
          
          if(pieceInfo && this.playerColor == pieceInfo.pieceColor[0]){
            this.selectedSrc = {id:pieceInfo.id,pieceColor:pieceInfo.pieceColor[0],pieceType:pieceInfo.pieceType}
            this.$store.commit('setSelection',{row:pieceInfo.row,col:pieceInfo.col,piece:pieceInfo.pieceType})
          }
          else{ 
            this.selectedSrc = null

          }
        },
        updateBoardState1D(flipped){
          this.isFlipped = flipped
          var stack = new Array();
          this.boardState1D=[]
          var row,tile,x=1,y=1,flipX = this.rows,flipY = this.cols ;
          var tileId = flipped ? this.rows*this.cols - 1 : 0;
          for(row of this.boardState.tiles){
              //console.log('reachin',flipped)
              for(tile of row){
                tile.tileId = tileId;
                tileId+= flipped? -1 : 1;     
                tile.x= flipped? flipX : x;
                tile.row = x
                tile.tileType = this.isLight(y,x)? 'l' : 'd';
                tile.y= flipped? flipY : y;
                tile.col = y
                y+=1
                flipY-=1
                if(flipped){
                  stack.push(tile)
                }
                else{
                this.boardState1D.push(tile);
                
                }
              }
              x+=1
              flipX-=1
              flipY = this.cols
              y=1;
          }
          if(flipped){
            var squares = stack.length
            for(var i=0;i<squares;i++){
              this.boardState1D.push(stack.pop())
            }
          }
          console.log('1d',this.boardState1D)
      },
        isEven(val){return val%2==0},
        isLight(row,col){
            return this.isEven(row)&&this.isEven(col)|| (!this.isEven(row)&&!this.isEven(col))},
    },
    computed:{
      cssVar(){
        return {
        '--sizex': this.rows,
        '--sizey': this.cols,
        '--size': Math.max(this.rows,this.cols),
        }
      }
    },
}
</script>

<style >
#board-container{

    max-width: 700px;
    width:100%;
    max-height: 700px;

}

#board{
    background-color: #EAEAEA;
    grid-template-columns: repeat(var(--size), 1fr);
    grid-template-rows: repeat(var(--size), 1fr);
    box-shadow: 0 14px 28px rgba(0,0,0,0.25), 0 10px 10px rgba(0,0,0,0.22);
    display: grid;
    justify-items: center;
    }

</style>