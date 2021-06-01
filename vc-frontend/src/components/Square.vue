<template>
  <div class="square" :class="isDark? 'dark':'light'"  @click="emitToBoard">
      <div v-if="isPiecePresent">
      <piece  :color="pieceColor" :pieceType="pieceType"/>
      </div>
      
      
  </div>
</template>

<script>
import Piece from './Piece'
export default {
    components:{Piece},
    methods:{
      toggleIsPiecePresent(){
        this.isPiecePresent = !this.isPiecePresent;
      },
      emitToBoard(){
          this.$emit("childToParent",this.row,this.col)
      }
    },

    props:{
        isDark: Boolean,
        row: Number,
        col: Number,
        isPiecePresent: Boolean,
        pieceType: String,
        pieceColor: String,
    }
}
</script>

<style scoped>


.square {
  background: transparent;
  border: 1px solid transparent;
  float: left;
  font-size: 24px;
  font-weight: bold;
  line-height: 34px;
  height: 48px;/*48px;*/
  margin-right: -1px;
  margin-top: -1px;
  padding: 0;
  text-align: center;
  width: 48px;
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