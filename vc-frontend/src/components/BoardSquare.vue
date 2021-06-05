<template>
  <div class="square" :style="cssVar" :class="tileType=='d'? 'dark':'light'"  @click="emitToBoard">
      <div v-if="isPiecePresent">
      <board-piece  :color="pieceColor" :pieceType="pieceType" :row="row" :col="col"/>
      </div>

      
  </div>
</template>

<script>
import BoardPiece from './BoardPiece'
export default {
    components:{BoardPiece},
    methods:{
      toggleIsPiecePresent(){
        this.isPiecePresent = !this.isPiecePresent;
      },
      emitToBoard(){
          this.$emit("childToParent",this.row,this.col)
      }
    },

    props:{
        tileType: String,
        row: Number,
        col: Number,
        isPiecePresent: Boolean,
        pieceType: String,
        pieceColor: String,
    },
    computed:{
        cssVar(){
        return {
        '--x': this.row,
        '--y': this.col,
        }
      }
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
}

.dark {
  background-color: #b2c85d;
}

.light {
  background-color: #e4f5cb;
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