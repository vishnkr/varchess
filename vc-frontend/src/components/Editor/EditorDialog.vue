<template>
  <div class="board-editor">
    <div class="side-panel">
      <v-card class="mx-auto" max-width="550">
        <v-list-item>
          <v-list-item-content>
            <div>
                <v-list-item-title class="headline">
                Game Editor
                </v-list-item-title>
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
                    <v-list-item-title style="padding-top:10px">Board height</v-list-item-title>
                    <v-slider
                      v-model="rows"
                      :max="14"
                      @input="updateBoardDimensions"
                      :tick-labels="labels"
                      color="orange darken-3"
                    ></v-slider>
                  </div>
                  <div>
                    <v-list-item-title>Board width</v-list-item-title>
                    <v-slider
                      v-model="cols"
                      :max="14"
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
                    v-model="colorSelect"
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
                    v-model="pieceSelect"
                    column
                  >
                    <v-radio v-for="piece in pieceList"
                      :key="`${piece.toLowerCase()}`"
                      :label="`${piece}`"
                      :value="`${pieceMap[piece].toLowerCase()}`"
                    ></v-radio>
                    
                  </v-radio-group>
                </div>
                <div v-if="pieceSelect=='c'">
                  <v-card
                    elevation="16"
                    max-width="400"
                    class="mx-auto"
                  >
                    <v-virtual-scroll
                      
                      :items="pieceFilter"
                      height="200"
                      item-height="64"
                    >
                      <template v-slot:default="{ item }">
                        <v-list-item class="scroll-item">
                          <v-list-item-action>
                            <v-btn
                              depressed
                              color="primary"
                              @click="customPieceSelect = item.piece"
                            >
                              Select
                            </v-btn>
                          </v-list-item-action>

                          <v-list-item-content>
                            <img class="resize" :src="item.src">
                          </v-list-item-content>

                          <v-list-item-action>
                            <v-btn
                              depressed
                              color="primary"
                              @click="dialog=true"
                            >
                              Set Move Pattern
                            </v-btn>
                            <move-pattern-dialog v-if="dialog" v-on:closeDialog="closeDialog" 
                              :dialog="dialog"
                              :pieceColor="curPieceColor"
                              :pieceType="customPieceSelect"/>
                          </v-list-item-action>
                        </v-list-item>

                        <v-divider></v-divider>
                      </template>
                    </v-virtual-scroll>
                  </v-card>
                </div>
              </v-tab-item>
          </v-tabs>
        </v-list-item-content>
        </v-list-item>
      </v-card>
    </div>
    <editor-board v-on:sendBoardState="setBoardState" 
      class="board-panel" 
      :cols="labels[cols]" 
      :rows="labels[rows]" 
      :editorMode="true" 
      :curPiece="pieceSelect=='c'? customPieceSelect : pieceSelect"
      :curPieceColor="colorSelect"
       />
  </div>
</template>

<script>
import { convertBoardStateToFEN } from '../../utils/fen';
import {createRoom} from '../../utils/websocket';
import EditorBoard from './EditorBoard';
import  MovePatternDialog from './MovePatternDialog.vue';
export default {
  components:{EditorBoard,MovePatternDialog},
  methods:{
    getPieceURL(piece){
      return require(`../../assets/images/pieces/${this.colorSelect}/${piece}.svg`)
    },
    closeDialog(){
      this.dialog=false
    },
    enterRoom(){
      var finalboardState = this.setBoardState(this.boardState)
      var fenString = convertBoardStateToFEN(finalboardState,'w','KQkq','-');
      createRoom(this.ws,this.roomId,this.username, fenString);
      this.$router.push({name:'Game', params:{username: this.username,roomId: this.roomId, boardState: finalboardState, ws:this.ws}})
    },
    updateBoardDimensions(){
      this.setBoardState(this.boardState);
    },
    setBoardState(boardState){
      this.boardState=boardState;
      var newBoardState={tiles:[],rows:this.rows+1,cols:this.cols+1};
      for(var row=0;row<this.rows+1;row++){
        newBoardState.tiles.push(boardState.tiles[row].slice(0,this.cols+1));
      }
      var payload = {boardState: newBoardState, roomId: this.roomId}
      this.$store.commit('updateBoardState',payload);
      return newBoardState
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
  },
  data(){
    return{
      labels: [1,2,3,4,5,6,7,8,9,10,11,12,13,14,15],
      pieceList: ['Pawn','King','Queen','Bishop','Knight','Rook','Custom'],
      customPieces:['a','j','d','i','g','s','u','v','z'],
      pieceMap: {'Pawn':'p','King':'k','Queen':'q','Bishop':'b','Knight':'n','Rook':'r','Custom':'c'},
      customPieceMap:{},
      rows: 7,
      cols: 7,
      dialog:false,
      colorSelect: 'white',
      pieceSelect: 'pawn',
      customPieceSelect: '',
      isDisableTileOn: false,
      boardState:{},
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
}
.side-panel{
  flex:3;
  height: 90vh;
  /*box-shadow: 0 14px 28px rgba(250, 0, 0, 0.25), 0 10px 10px rgba(0,0,0,0.22);*/
}
.board-panel{
  flex:6;
}

.list-header{
  display:flex;
}

.resize{
  height: 3em;
  width: 3em;
}

@media only screen and (max-device-width: 480px) {
  .board-editor{
    display: flex;
    flex-direction: column;
  }

}
</style>