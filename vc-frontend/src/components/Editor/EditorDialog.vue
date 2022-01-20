<template>
  <div class="board-editor">
    <div class="side-panel">
      <v-card class="mx-auto" max-width="550">
        <v-list-item>
          <v-list-item-content>
            <div class="card-top">
                <v-list-item-title class="headline">
                Game Editor
                </v-list-item-title>
              <v-btn color="red" dark style="margin-right:5px;" @click="clearBoard"><v-icon>mdi-delete</v-icon> Clear</v-btn>
              <v-btn color="success" depressed @click="enterRoom">
                  Save Setup
                </v-btn>
            </div>
            <v-divider>Board settings</v-divider>
            <v-spacer /><v-spacer />
            <v-tabs
              fixed-tabs
            >
              <v-tab>
                Board Editor
              </v-tab>
              <v-tab-item style="padding-left:10px" key="board-edit">
                <div>
                  <div>
                    <v-list-item-title style="padding-top:10px">Board height (rows)</v-list-item-title>
                    <v-slider
                      v-model="labelRow"
                      :max="labels.length-1"
                      @input="updateBoardDimensions"
                      :tick-labels="labels"
                      color="orange darken-3"
                    ></v-slider>
                  </div>
                  <div>
                    <v-list-item-title>Board width (cols)</v-list-item-title>
                    <v-slider
                      v-model="labelCol"
                      :max="labels.length-1"
                      @input="updateBoardDimensions"
                      :tick-labels="labels"
                      color="orange darken-3"
                    ></v-slider>
                  </div>
                  <v-checkbox
                    v-model="isDisableTileOn"
                    label="Toggle: Disable selected tile"
                  ></v-checkbox>

                </div>
              </v-tab-item>

              <v-tab>
                Piece Placement
              </v-tab>
              <v-tab-item style="padding-left:10px"  key="piece-place">
                <div>
                  <v-list-item-title style="padding-top:10px">Choose color</v-list-item-title>
                  <v-radio-group
                    v-model="editorData.curPieceColor"
                    row
                  >
                    <v-radio
                      label="Black"
                      value="black"
                    ></v-radio>
                    <v-radio
                      label="White"
                      value="white"
                    ></v-radio>
                  </v-radio-group>
                </div>
                <div>
                  <v-list-item-title style="padding-top:10px">Choose piece</v-list-item-title>
                  <v-radio-group
                    v-model="editorData.curPiece"
                    column
                  >
                    <v-radio v-for="piece in pieceList"
                      :key="`${piece.toLowerCase()}`"
                      :label="`${piece}`"
                      :value="`${pieceMap[piece].toLowerCase()}`"
                    ></v-radio>
                    
                  </v-radio-group>
                </div>
                <div v-if="editorData.curPiece=='c'">
                  <div class="custom-pieces">
                    <v-btn
                      depressed
                      color="primary"
                      @click="dialog=true"
                      class="move-button"
                    >
                      Set Move Pattern
                    </v-btn>
                    <div class="piece-scroll">
                      <v-card
                        elevation="16"
                        max-width="150"
                        class="mx-auto"
                      >
                        <v-virtual-scroll
                          
                          :items="pieceFilter"
                          height="200"
                          item-height="64"
                        >
                          <template v-slot:default="{ item }">
                            <v-radio-group v-model="editorData.customPiece">
                            <v-list-item class="scroll-item" :class="{added : editorData.added[item.piece], defined: editorData.defined[item.piece]}">
                                <v-radio
                                  :key="item.piece"
                                  :value="item.piece"
                                  @click="editorData.customPiece = item.piece"
                                >
                                  Select
                                </v-radio>
                              <v-list-item-content>
                                <img class="resize" :src="item.src">
                              </v-list-item-content>

                              <v-list-item-action>
                                
                                <move-pattern-dialog v-if="dialog" 
                                  v-on:closeDialog="closeDialog" 
                                  v-on:movePatterns="setMovePattern"
                                  :dialog="dialog"
                                  :editorData="editorData"
                                  :pieceColor="editorData.curPieceColor"
                                  :pieceType="editorData.customPiece"
                                  :ws="ws"
                                  />
                              </v-list-item-action>
                            </v-list-item>
                            <v-divider></v-divider>
                            </v-radio-group>
                          </template>
                        </v-virtual-scroll>
                      </v-card>
                    </div>
                  </div>
                </div>
              </v-tab-item>
          </v-tabs>
        </v-list-item-content>
        </v-list-item>
      </v-card>
      
    </div>
    
    <board :board="boardState" :isflipped="false" :editorMode="true" 
    :editorData="editorData" :key="change" 
    v-on:sendEditorboardState="formatBoardState"
    v-on:customPieceAdd="customPieceAdd"
    />
  </div>
</template>

<script>
import { convertBoardStateToFEN } from '../../utils/fen';
import {createRoomWithCustomPatterns} from '../../utils/websocket';
import {validateStartSetup} from '../../utils/validator';
import Board from '../Board/Board.vue';
import  MovePatternDialog from './MovePatternDialog.vue';
export default {
  components:{MovePatternDialog,Board},
  created(){
    this.setupDefaultBoardMaxSize()
  },
  methods:{
    getPieceURL(piece){
      return require(`../../assets/images/pieces/${this.colorSelect}/${piece}.svg`)
    },
    customPieceAdd(piece){
      this.editorData.added[piece]=true
    },
    closeDialog(){
      this.dialog=false
    },
    clearBoard(){
      for(var row of this.boardState.tiles){
        for(var tile of row){
          if(tile.isPiecePresent){
            tile.isPiecePresent = false;
            tile.pieceType = null;
            tile.pieceColor = null;
          }
        }
      }
      
    },
    enterRoom(){
      var finalboardState = this.boardState
      var fenString = convertBoardStateToFEN(finalboardState,'w','KQkq','-');
      createRoomWithCustomPatterns(this.ws,this.roomId,this.username, fenString,this.customMovePatterns);
      if(this.customMovePatterns!=[]){
        this.$store.commit('storeMovePatterns',{movePatterns: this.customMovePatterns})
      }
      if(validateStartSetup(fenString)){
        this.$store.commit('websocketError',null)
        this.$router.push({name:'Game', params:{username: this.username,roomId: this.roomId, boardState: finalboardState, ws:this.ws}})
      } else {
        this.$store.commit('websocketError','Board state not valid: must contain 1 king for each color & not under check')
      }
   },

    setMovePattern(piece,jumpPattern,slidePattern){
      this.customMovePatterns.push({piece:piece,jumpPattern:jumpPattern,slidePattern:slidePattern});
      this.editorData.defined[piece] = true;
    },
    updateBoardDimensions(){
      this.formatBoardState(this.maxBoardState);
    },

    isEven(val){return val%2==0},
    isLight(row,col){
      return this.isEven(row)&&this.isEven(col)|| (!this.isEven(row)&&!this.isEven(col))},

    formatBoardState(boardState){
      this.boardState={tiles:[],rows:this.rows,cols:this.cols};
      for(var row=0;row<this.rows;row++){
        this.boardState.tiles.push(boardState.tiles[row].slice(0,this.cols));
      }
      var payload = {boardState: this.boardState, roomId: this.roomId}
      //trigger re-render of editor board
      this.change = this.change ? 0 : 1; 
      this.$store.commit('updateBoardState',payload);
      
    },
    setupDefaultBoardMaxSize(){
        for(let row =0;row<this.maxBoardState.rows;row++){
          this.maxBoardState.tiles.push([])
          for(let col=0;col<this.maxBoardState.cols;col++){
            var tile = {}
            tile.tileType = this.isLight(row,col)? 'l' : 'd';
            if(this.isEmpty){
              tile.isPiecePresent=false
            }
            else if( (col===0||col==7) && (row===0||row===7)) {
                tile.pieceType='r';
                tile.isPiecePresent=true
            }
            else if((col===1||col==6) && (row===0||row===7)){
                tile.pieceType='n';
                tile.isPiecePresent=true
            }
            else if((col===2||col==5) && (row===0||row===7)){
                tile.pieceType='b';
                tile.isPiecePresent=true
            }
            else if((col===3) && (row===0||row===7)){
                tile.pieceType='q';
                tile.isPiecePresent=true
            }
            else if((col===4) && (row===0||row===7)){
                tile.pieceType='k';
                tile.isPiecePresent=true
            }
            else if((row===1||row===6) && col<8){
                tile.pieceType='p';
                tile.isPiecePresent=true
            }
            else{tile.isPiecePresent=false}
            if(row==0||row==1){tile.pieceColor='black'}
            else if(row==6||row==7){tile.pieceColor='white'}
            this.maxBoardState.tiles[row].push(tile)
          }
        }
        this.boardState.rows = this.rows
        this.boardState.cols = this.cols
        this.formatBoardState(this.maxBoardState)
      }
  },
  computed: {
      items () {
        return Array.from({ length: this.pieces.length }, (k, v) => v + 1)
      },
      pieceFilter(){
        var pieceUrls = []
        for(var piece of this.customPieces){
          pieceUrls.push({piece:piece,src :this.getPieceURL(piece)})
        }        
        return pieceUrls
      },
      rows(){
        return this.labels[this.labelRow]
      },
      cols(){
        return this.labels[this.labelCol]
      }

  },
  data(){
    return{
      editorData: {curPieceColor:'white',curPiece:'p',added:{},defined:{}},
      labels: [5,6,7,8,9,10,11,12,13,14,15,16],
      pieceList: ['Pawn','King','Queen','Bishop','Knight','Rook','Custom'],
      customPieces:['a','j','d','i','g','s','u','v','z'],
      pieceMap: {'Pawn':'p','King':'k','Queen':'q','Bishop':'b','Knight':'n','Rook':'r','Custom':'c'},
      customMovePatterns:[],
      labelRow: 3,
      labelCol: 3,
      change:0,
      dialog:false,
      colorSelect: 'white',
      pieceSelect: 'pawn',
      isDisableTileOn: false,
      boardState:{tiles:[]},
      maxBoardState:{tiles:[],rows:16,cols:16},
      username: this.$route.params.username,
      roomId: this.$route.params.roomId,
      ws: this.$route.params.ws,
    }
  }

}
</script>

<style scoped>
.board-editor{
  display: flex;
  margin: 1em;
}
.side-panel{
  flex:3;
  height: 90vh;
}
.board-panel{
  flex:6;
}
.card-top{
  display: flex;
}
.list-header{
  display:flex;
}

.custom-pieces{
  display: flex;
}
.move-button{
  flex:1;
}
.piece-scroll{
  flex:2;
}
.added{
  background-color: rgb(236, 111, 111);
}
.defined{
  background-color: #21af49 !important;
}

@media only screen and (max-device-width: 480px) {
  .board-editor{
    display: flex;
    flex-direction: column;
  }

}
</style>