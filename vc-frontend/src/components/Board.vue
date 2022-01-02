<template>
  <div :style="cssVar" id='board-container'>
      <div id="board" style="z-index:-1;">
        <board-square v-for="square in this.boardState1D"  :key="square.tileId" 
        ref="squares"
        :tileType="square.tileType" 
        :isPiecePresent="square.isPiecePresent" 
        :tileId="square.tileId"
        :editorMode="editorMode"
        :editorData="editorData? editorData: null"
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

<script>
import BoardSquare from './BoardSquare.vue';

export default {
  components: { BoardSquare },
  props:['board','isflipped','playerColor',"editorMode","editorData","boardSize","mpTabData"],
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
  },
  data(){
        return {
            boardState: [],
            boardState1D: [],
            rows: 0,
            cols:0,
            selectedSrc: null,
            highlightData:null, 
            roomId: this.$route.params.roomId,
        }
    },
    methods:{
      handleEditorSquareClick(type,row,col){
        if (type=="regular"){
          this.editorModeSquareClicked(row,col)
        } else {
          this.$emit('setMP',row,col)
        }
      },
      setSlidePattern(slideDirections){
        let tileIDs = []
        for(var i = 0;i<this.boardState1D.length;i++){
          this.$refs.squares[i].removeColorFromSquare()
        }
        for (var direction of slideDirections){
          let row = this.editorData.piecePos[0]+direction[0]
          let col = this.editorData.piecePos[1]+direction[1]
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
        if(this.boardState.tiles[row-1][col-1].isPiecePresent){
          this.boardState.tiles[row-1][col-1].isPiecePresent = false;
        }
        else{
          this.boardState.tiles[row-1][col-1].isPiecePresent = true;
          this.boardState.tiles[row-1][col-1].pieceColor = this.editorData.curPieceColor;
          this.boardState.tiles[row-1][col-1].pieceType = this.editorData.curPiece =='c' ? this.editorData.customPiece : this.editorData.curPiece;
          this.editorData.curPiece =='c'? this.$emit("customPieceAdd",this.editorData.customPiece) : null;
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
            var oldRookPos = moveInfo.destCol<moveInfo.srcCol? 0 : this.board.tiles[0].length -1;
            var newRookPos = moveInfo.destCol<moveInfo.srcCol? moveInfo.srcCol-1 : moveInfo.srcCol+1
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
            this.$store.commit('setSelection',{row:pieceInfo.row,col:pieceInfo.col,piece:pieceInfo.pieceType})
            this.highlightData = {from:pieceInfo.id,to:tomoves.moves};
          }
          else{ 
            this.selectedSrc = null;
            this.highlightData = null;
          }
        },
        updateBoardState1D(flipped){
          this.isFlipped = flipped
          var stack = new Array();
          this.boardState1D=[]
          var row,tile,x=1,y=1,flipX = this.rows,flipY = this.cols ;
          var tileId = flipped ? this.rows*this.cols - 1 : 0;
          for(row of this.boardState.tiles){
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
      },
        isEven(val){return val%2==0},
        isLight(row,col){
            return this.isEven(row)&&this.isEven(col)|| (!this.isEven(row)&&!this.isEven(col))},
    },
    computed:{
      cssVar(){
        return {
        '--container_size': this.boardSize ? `${this.boardSize}px` : `${700}px`,
        '--size': Math.max(this.rows,this.cols),
        }
      }
    },
}
</script>

<style >
#board-container{

    max-width: var(--container_size);
    width:100%;
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