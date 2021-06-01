<template>
  <div :style="cssVar" id="board-container">
  <div id="board" >
      <div  v-for="col in cols" :key="col">
        <div v-for="row in rows" :key="row">
          <square v-if="boardState.tiles[row-1][col-1].isPiecePresent" v-on:childToParent="squareClicked" :isDark="!isLight(row,col)" :row="row" :col="col" :isPiecePresent="boardState.tiles[row-1][col-1].isPiecePresent" :pieceType="boardState.tiles[row-1][col-1].pieceType" :pieceColor="boardState.tiles[row-1][col-1].pieceColor" />
          <square v-else v-on:childToParent="squareClicked" :isDark="!isLight(row,col)" :row="row" :col="col" />
        </div>
      </div>
        
  </div>
  </div>
</template>

<script>

import Square from '../Square'

export default {
    name:'Board',
    methods:
    {
      squareClicked(row,col){
        
        if(this.boardState.tiles[row-1][col-1].isPiecePresent){
          this.boardState.tiles[row-1][col-1].isPiecePresent = false;
        }
        else{
          
          this.boardState.tiles[row-1][col-1].isPiecePresent = true;
          this.boardState.tiles[row-1][col-1].pieceColor = this.curPieceColor;
          console.log(this.curPiece.toLowerCase(), this.pieceList, this.curPiece)
          this.boardState.tiles[row-1][col-1].pieceType = this.pieceList[this.curPiece];
        }

      },


      isEven(val){return val%2==0},
      isLight(row,col){
        return this.isEven(row)&&this.isEven(col)|| (!this.isEven(row)&&!this.isEven(col))},

      setupDefaultBoard(){

        for(var col =0;col<this.maxCols;col++){
          this.boardState.tiles.push([])
          for(var row=0;row<this.maxRows;row++){
            var tile = {}
            if( (row===0||row==7) && (col===0||col===7)) {
                tile.pieceType='r';
                tile.isPiecePresent=true
            }
            else if((row===1||row==6) && (col===0||col===7)){
                tile.pieceType='n';
                tile.isPiecePresent=true
            }
            else if((row===2||row==5) && (col===0||col===7)){
                tile.pieceType='b';
                tile.isPiecePresent=true
            }
            else if((row===3) && (col===0||col===7)){
                tile.pieceType='q';
                tile.isPiecePresent=true
            }
            else if((row===4) && (col===0||col===7)){
                tile.pieceType='k';
                tile.isPiecePresent=true
            }
            else if((col===1||col===6) && row<8){
                tile.pieceType='p';
                tile.isPiecePresent=true
            }
            else{tile.isPiecePresent=false}
            if(col==0||col==1){tile.pieceColor='black'}
            else if(col==6||col==7){tile.pieceColor='white'}
            this.boardState.tiles[col].push(tile)
          }
        }
        
      },
      
    },
    computed:{
      cssVar(){
        return {
        '--size': Math.max(this.maxCols,this.maxRows)
        }
      }
    },
    mounted(){
      if(!this.isMounted) {this.setupDefaultBoard();}
      this.isMounted = true;
    },
    data(){
      return {
        pieceList : {'pawn':'p','king':'k','queen':'q','bishop':'b','knight':'n','custom':'custom'},
        rowMultiplier:0,
        isMounted: false,
        maxRows: 15,
        maxCols: 15,
        boardState: {
          
          tiles:[],
        }
      }
    },
    props:{
      rows: Number,
      cols: Number,
      editorMode: Boolean,
      curPiece:String,
      curPieceColor:String,
    },
    components:{
      Square,
    },
    
    created(){
        
    }
}
</script>

<style scoped>


#board-container{
}

#board{
    padding: 2%;
    display: grid;
    grid-template-columns: repeat(var(--size), 1fr);
    /*grid-template-rows: repeat(var(--size), 1fr);*/
    box-shadow: 0 14px 28px rgba(0,0,0,0.25), 0 10px 10px rgba(0,0,0,0.22);
    align-items: center;
    justify-content: center;
    
}

@media only screen and (max-device-width: 480px) {
  #board-container{
    width: 50%;
    
  }
  #board{
    
    display: grid;
    grid-template-columns: repeat(var(--size), 0fr);
    /*grid-template-rows: repeat(var(--size), 1fr);*/
    box-shadow: 0 14px 28px rgba(0,0,0,0.25), 0 10px 10px rgba(0,0,0,0.22);
    align-items: center;
    justify-content: center;
    
}
}


</style>