<template>
  <div class="members" v-chat-scroll="{always:false, smooth:true}">
      <v-card
        class="mx-auto"
        tile
        elevation="2"
      >
        <div class="piece-buttons">
            <span v-for="pieceInfo in movePatterns" :key="pieceInfo.piece">
                <v-btn
                    medium
                    color="light-green lighten-2"
                    dark
                    class="piece-button"
                    @click="changePiece(pieceInfo.piece)"
                    >
                    <piece :pieceType="pieceInfo.piece" :color="color"/>
                </v-btn>
            </span>
        </div >
        <div class="small-board">
            <board 
            :key="mpBoardKey"
            :board="boardState"
            :mpTabData="mpTabData"
            />
        </div>
        <div>

        </div>
    </v-card>
  </div>
</template>



<script>

import Piece from '../BoardPiece.vue'
import Board from '../Board.vue'
export default {
    components:{Piece, Board},
    created(){
        this.setupMPBoard()
    },
    props:['color','movePatterns'],
    data: () => ({
      selectedItem: null,
      boardState : {tiles:[]},
      mpBoardKey:false,
      rows:9,
      cols:9,
      curPiece: null,
      curPos: [5,5],
      //contains row and col which should be colored on the smaller board
      mpTabData:{}
    }),
    methods:{
        changePiece(piece){
            this.curPiece=piece
            this.setupMPBoard()
            this.mpBoardKey = !this.mpBoardKey
            this.fillMpTabData(piece)
        },
        fillMpTabData(piece){
            console.log(this.movePatterns)
            this.mpTabData={}
            var patterns = {}
            for (var pattern of this.movePatterns){
                if (pattern.piece === piece){
                    patterns = pattern
                    break;
                }
            }
            var row,col
            var difference
            for (difference of patterns.jumpPattern){
                row = this.curPos[0]-difference[0]
                col = this.curPos[1]-difference[1]
                if(row>=0 && row<=this.rows && col>=0 && col<=this.cols){
                this.mpTabData[[row,col]] = 'move-jump-pattern'
                }
            }
            console.log('jumpa',this.mpTabData)
            row = this.curPos[0]
            col = this.curPos[1]
            console.log('pat',patterns)
            for(var direction of patterns.slidePattern){
                row = this.curPos[0]
                col = this.curPos[1]
                console.log('slides',row,col,direction)
                while(row>=0 && row<=this.rows && col>=0 && col<=this.cols){
                    console.log('slides',row,col,direction)
                    this.mpTabData[[row,col]] = 'move-slide-pattern'
                    row += direction[0]
                    col += direction[1]
                }
            }
            console.log('slidea',this.mpTabData)

        },
            
        isEven(val){return val%2==0},
        isLight(row,col){
          return this.isEven(row)&&this.isEven(col)|| (!this.isEven(row)&&!this.isEven(col))},
        setupMPBoard(){
          this.boardState = {tiles:[]};
          console.log('calledmpboard')
          for(var col =0;col<this.cols;col++){
            this.boardState.tiles.push([])
            for(var row=0;row<this.rows;row++){
              var tile = {}
              tile.tileType = this.isLight(col,row)? 'l' : 'd';
              if(row==4 && col==4){
                tile.isPiecePresent=true
                tile.pieceType = this.curPiece
                tile.pieceColor = this.color

              } else{ tile.isPiecePresent=false }
              this.boardState.tiles[col].push(tile)
            }
            
          }
          this.boardState.rows = 9
          this.boardState.cols = 9
        }
    }
}
</script>

<style scoped>
.member-container{
    height: 300px;
    overflow: auto;
    margin: 30px
}
.piece-button{
    margin: 1%;
}
</style>