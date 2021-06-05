<template>
  <div :style="cssVar" id='board-container'>
      <div id="board">
        <board-square v-for="square in boardState1D"  :key="square.tileId" :tileType="square.tileType" :isPiecePresent="square.isPiecePresent" :pieceType="square.pieceType" :pieceColor="square.pieceColor" :row="square.x" :col="square.y"  />
      </div>
  </div>
</template>

<script>
import BoardSquare from './BoardSquare.vue';
export default {
  components: { BoardSquare },
  created(){
      this.boardState = this.$store.state.boardState;
      var row,tile,x=1,y=1, tileId=0;
      //console.log(this.boardState)
      this.rows = this.boardState.tiles.length
      
      for(row of this.boardState.tiles){
          for(tile of row){
            tile.tileId = tileId;
            tileId+=1;
            
            tile.x=x;
            tile.y=y;
            //col.row=x-1;
            //col.col=y-1;
            y+=1
            this.boardState1D.push(tile);
          }
          x+=1
          y=1;
      }
      this.cols = this.boardState.tiles[0].length
      /*console.log("1d",this.boardState1D)
      console.log("cols",this.cols, "rows:",this.rows)
      console.log('ogbs',this.boardState)*/
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
    /*padding-bottom: 100%;*/

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