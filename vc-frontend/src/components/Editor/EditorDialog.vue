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
                      :value="`${piece.toLowerCase()}`"
                    ></v-radio>
                    
                  </v-radio-group>
                </div>
              </v-tab-item>
          </v-tabs>
        </v-list-item-content>
        </v-list-item>
      </v-card>
    </div>
    <editor-board v-on:sendBoardState="setBoardState" class="board-panel" :cols="labels[cols]" :rows="labels[rows]" :editorMode="true" :curPiece="pieceSelect" :curPieceColor="colorSelect" />
  </div>
</template>

<script>
import EditorBoard from './EditorBoard';
export default {
  components:{EditorBoard},
  methods:{
    enterRoom(){
      this.$router.push({ path: '/game/123' });
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
      this.$store.commit('updateBoardState',newBoardState);
    }
  },
  data(){
    return{
      labels: [1,2,3,4,5,6,7,8,9,10,11,12,13,14,15],
      pieceList: ['Pawn','King','Queen','Bishop','Knight','Custom'],
      rows: 7,
      cols: 7,
      colorSelect: 'white',
      pieceSelect: 'pawn',
      isDisableTileOn: false,
      boardState:{},
    }
  }

}
</script>

<style scoped>
.board-editor{
  display: flex;
}
.side-panel{
  flex:1;
  height: 90vh;
  /*box-shadow: 0 14px 28px rgba(250, 0, 0, 0.25), 0 10px 10px rgba(0,0,0,0.22);*/
}
.board-panel{
  flex:1;
}

.list-header{
  display:flex;
}

@media only screen and (max-device-width: 480px) {
  .board-editor{
    display: flex;
    flex-direction: column;
  }

}
</style>