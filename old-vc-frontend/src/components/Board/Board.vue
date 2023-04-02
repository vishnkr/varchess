<template>
  <div :style="cssVar" id='board-container'>
      <div id="board" style="z-index:-1;">
        <square v-for="square in this.boardState1D"  :key="square.tileId" 
        ref="squares"
        :tileType="square.tileType" 
        :isPiecePresent="square.isPiecePresent" 
        :tileId="square.tileId"
        :editorMode="editorMode"
        :editorState="editorState? editorState: null"
        :pieceType="square.pieceType" 
        :pieceColor="square.pieceColor" 
        :x="square.x" :y="square.y" 
        :row="square.row" :col="square.col"  
        :highlight="highlightData"
        :selectedSrc="selectedSrc"
        :mpTabData="mpTabData? mpTabData: null"
        v-on:sendSelectedPiece="setSelectedPiece"
        v-on:destinationSelect="emitDestinationSelect"
        v-on:setEditorBoardState="handleEditorSquareClick"
        />
        
      </div>
  </div>
</template>

<script lang="ts">
import Square from './Square.vue';
import Vue from 'vue';
import { BoardState,PiecePosition,Tile } from '../../types';
import { PERFORM_MOVE, SET_SRC_SELECTION } from '../../utils/mutation_types';
import { mapMutations } from 'vuex';

interface BoardData{
  board: BoardState,
  boardState1D: Tile[],
  rows: number,
  cols: number,
  highlightData:PiecePosition | null,
}

export default Vue.extend({
  components: { Square },
  props:['isflipped','playerColor',"editorMode","editorState","boardSize","mpTabData"],
  created(){
      console.log(this.$data);
  },
  watch: { 
    isflipped() { 
      this.updateBoardState1D(this.isflipped)
    }
  },
  mounted(){ 
      this.rows = this.board.rows
      this.cols = this.board.cols
      this.updateBoardState1D(this.isflipped)
      this.$store.subscribe((mutation,state)=>{
        if(mutation.type === PERFORM_MOVE && this.$store.state.currentMove){
          this.performMove(this.$store.state.currentMove)
        }
      })
  },
  data():BoardData{
        return {
            board: this.$store.state.board,
            boardState1D: [],
            rows: 0,
            cols:0,
            highlightData: null,
        }
    },
    methods:{
      ...mapMutations([SET_SRC_SELECTION]),
      handleEditorSquareClick(type,row,col){
        if (type=="regular"){
          if (this.editorState.isDisableTileOn){
            //perform disable
            this.boardState.tiles[row-1][col-1].disabled = !this.boardState.tiles[row-1][col-1].disabled
            return
          }
          this.editorModeSquareClicked(row,col)
        } else {
          this.$emit('setMP',row,col)
        }
      },
      setSlidePattern(slideDirections){
        let tileIDs:number[] = [];
        for(var i = 0;i<this.boardState1D.length;i++){
          this.$refs.squares[i].removeColorFromSquare()
        }
        for (let direction of slideDirections){
          let row = this.editorState.piecePos[0]+direction[0]
          let col = this.editorState.piecePos[1]+direction[1]
          while(row<this.boardState.rows && row>=0 && col<this.boardState.cols && col>=0){
            tileIDs.push(this.boardState.tiles[row][col].tileId)
            row+=direction[0]
            col+=direction[1]
          }
        } 
        for(i=0;i<tileIDs.length;i++){
          this.$refs.squares[tileIDs[i]].addColorToSquare('slide')
        }
      },
      editorModeSquareClicked(row,col){
        if(this.board.tiles[row-1][col-1].isPiecePresent){
          this.board.tiles[row-1][col-1].isPiecePresent = false;
        }
        else{
          this.boardState.tiles[row-1][col-1].isPiecePresent = true;
          this.boardState.tiles[row-1][col-1].pieceColor = this.editorState.curPieceColor;
          this.boardState.tiles[row-1][col-1].pieceType = this.editorState.curPiece =='c' ? this.editorState.customPiece : this.editorState.curPiece;
          this.editorState.curPiece =='c'? this.$emit("customPieceAdd",this.editorState.customPiece) : null;
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
            let oldRookPos = moveInfo.destCol<moveInfo.srcCol? 0 : this.board.tiles[0].length -1;
            let newRookPos = moveInfo.destCol<moveInfo.srcCol? moveInfo.srcCol-1 : moveInfo.srcCol+1
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
        async setSelectedPiece(pieceInfo){
          
          if(pieceInfo && this.playerColor == pieceInfo.pieceColor[0]){
            let tomoves = await this.$store.dispatch('getPossibleToSquares',{roomId:this.roomId,color:pieceInfo.pieceColor,srcRow:pieceInfo.row-1,srcCol:pieceInfo.col-1,piece:pieceInfo.pieceType});
            this.selectedSrc = {id:pieceInfo.id,pieceColor:pieceInfo.pieceColor[0],pieceType:pieceInfo.pieceType}
            this.SET_SRC_SELECTION({row:pieceInfo.row,col:pieceInfo.col,piece:pieceInfo.pieceType})
            this.highlightData = {from:pieceInfo.id,to:tomoves.moves};
          }
          else{ 
            this.selectedSrc = null;
            this.highlightData = null;
          }
        },
        updateBoardState1D(flipped){
          this.isFlipped = flipped
          let stack = new Array();
          this.boardState1D=[]
          let row,tile,x=1,y=1,flipX = this.rows,flipY = this.cols ;
          let tileId = flipped ? this.rows*this.cols - 1 : 0;
          for(row of this.$store.state.board.tiles){
              for(tile of row){
                let newTile: Tile = {
                  tileId: tile.tileId + (flipped ? -1 : 1),
                  x: flipped? flipX : x,
                  row: x,
                  col: y,
                  tileType : tile.disabled ? 'disabled' : this.isLight(y,x)? 'l' : 'd',
                  y: flipped? flipY : y,
                  ...tile
                };
                y+=1;
                flipY-=1;
                this.boardState1D.push(newTile);
              }
              x+=1
              flipX-=1
              flipY = this.cols
              y=1;
          }
          if(flipped){
            this.boardState1D.reverse();
          }
      },
        isEven(val){return val%2==0},
        isLight(row,col){
            return this.isEven(row)&&this.isEven(col)|| (!this.isEven(row)&&!this.isEven(col))},
    },
    computed:{
      selectedSrc(){
        return this.$store.state.curStartPos
      },
      cssVar(){
        return {
        '--container_size': this.boardSize ? `${this.boardSize}px` : `${700}px`,
        '--size': Math.max(this.rows,this.cols),
        }
      }
    },
});
</script>

<style >
#board-container{

    max-width: var(--container_size);
    width:80%;
    max-height: var(--container_size);

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