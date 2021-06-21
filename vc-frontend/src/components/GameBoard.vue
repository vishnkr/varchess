<template>
  <div :style="cssVar" id='board-container'>
      <div id="board">
        <board-square v-for="square in this.boardState1D"  :key="square.tileId" :tileType="square.tileType" :isPiecePresent="square.isPiecePresent" :pieceType="square.pieceType" :pieceColor="square.pieceColor" :x="square.x" :y="square.y" :row="square.row" :col="square.col"  />
        
      </div>
  </div>
</template>

<script>
import BoardSquare from './BoardSquare.vue';
import {convertBoardStateToFEN} from '../utils/fen'
export default {
  components: { BoardSquare },
  props:['board','isflipped'],
  mounted(){
      this.boardState = this.board;
      this.rows = this.boardState.tiles.length
      this.cols = this.boardState.tiles[0].length
      this.updateBoardState1D(this.isflipped)
      convertBoardStateToFEN(this.boardState,'w','KQkq','-')
  },
  data(){
        return {
            boardState: [],
            boardState1D: [],
            rows: 0,
            cols:0,
            
        }
    },
    methods:{
        updateBoardState1D(flipped){
          this.isFlipped = flipped
          var stack = new Array();
          this.boardState1D=[]
          var row,tile,x=1,y=1,flipX = this.rows,flipY = this.cols, tileId=0;
          for(row of this.boardState.tiles){
              for(tile of row){
                tile.tileId = tileId;
                tileId+=1;     
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