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
               <board-editor 
               :editorState="editorState"
               @update-board-dimensions="updateBoardDimensions"
               />
              </v-tab-item>

              <v-tab>
                Piece Placement
              </v-tab>
              <v-tab-item style="padding-left:10px"  key="piece-place">
                <piece-editor 
                :editorState="editorState"
                @update-piece-state="updatePieceState"
                @set-move-pattern="setMovePattern"
                />
              </v-tab-item>
          </v-tabs>
        </v-list-item-content>
        </v-list-item>
      </v-card>
      
    </div>
    
    <board :board="boardState" :isflipped="false" :editorMode="true" 
    :editorState="editorState" :key="change" 
    v-on:sendEditorboardState="formatBoardState"
    v-on:customPieceAdd="customPieceAdd"
    />
  </div>
</template>

<script >
import { convertBoardStateToFEN } from '../../utils/fen';
import {validateStartSetup} from '../../utils/validator';
import Board from '../Board/Board.vue';
import BoardEditor from './BoardEditor.vue';
import PieceEditor from './PieceEditor.vue';
import { mapActions, mapMutations } from 'vuex';
import Vue from 'vue';
import { UPDATE_BOARD_STATE } from '../../utils/mutation_types';

export default Vue.extend({
  components:{BoardEditor,PieceEditor,Board},
  created(){
    this.setupDefaultBoardMaxSize()
  },
  methods:{
    ...mapActions('webSocket',['connect']),
    ...mapActions(['createRoom']),
    ...mapMutations([UPDATE_BOARD_STATE]),
    customPieceAdd(piece){
      this.editorState.added[piece]=true
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
    async enterRoom(){
      var finalboardState = this.boardState
      var fenString = convertBoardStateToFEN(finalboardState,'w','KQkq','-');
      try {
        this.roomId = await this.createRoom({fen:fenString,movePatterns:this.customMovePatterns});
          if(this.roomId){
            if(validateStartSetup(fenString)){
              if(this.customMovePatterns!=[]){
                this.SET_MOVE_PATTERNS({movePatterns: this.customMovePatterns})
              }
              this.connect({ roomId: this.roomId, username: this.username });
              this.SET_SERVER_STATUS({isOnline:true,errorMessage:null})
              this.$router.push({
                    name: 'Game',
                    params: {
                      username: this.username,
                      roomId: this.roomId,
                    },
                  });
          } else {
            this.SET_SERVER_STATUS({errorMessage:'Board state not valid: must contain 1 king for each color & not under check'})
          }
        }
        
      } catch (error){
        console.error(error)
      }
      
   },

    setMovePattern(piece,jumpPattern,slidePattern){
      for (let pattern of this.customMovePatterns){
        if (pattern.piece === piece) {
          pattern.jumpPattern = jumpPattern;
          pattern.slidePattern = slidePattern;
        }
        return
      }
      this.customMovePatterns.push({piece:piece,jumpPattern:jumpPattern,slidePattern:slidePattern});
      this.editorState.defined[piece] = true;
    },
    updateBoardDimensions(dimensions){
      this.rows = dimensions.rows;
      this.cols = dimensions.cols;
      this.formatBoardState(this.maxBoardState);
    },

    updatePieceState(editorState){
      this.editorState = editorState;
    },

    
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

  },
  data(){
    return{
      editorState: {
        curPieceColor:'white',
        curPiece:'p',
        isDisableTileOn:false,
        added:{},
        defined:{}
      },
      customMovePatterns:[],
      change:0,
      rows:8,
      cols:8,
      pieceSelect: 'pawn',
      boardState:{tiles:[]},
      maxBoardState:{tiles:[],rows:16,cols:16},
      username: this.$route.params.username,
      roomId:null,
    }
  }

});
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