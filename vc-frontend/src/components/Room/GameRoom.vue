<template>
  <div class="container">
    <div class="columns">
        <div class="column board-panel">
          <board :board="boardState" 
            ref="gameBoard" 
            :isflipped="isFlipped" 
            :playerColor="player1 == username ? 'w' : player2 == username ? 'b' : null"
            :editorMode="false"
            v-on:destinationSelect="validateMove"
            />
        </div>
        <div class="column right-panel">
          <v-row row d-flex nowrap align="center" justify="center" class="px-2">
            <v-text-field width="10px" class="centered-input" v-model="shareLink" id="tocopy" readonly  ></v-text-field>
            <v-btn width="6.5rem" class="ma-2" rounded dark color="blue" @click="copyText">Copy<v-icon>fas fa-link</v-icon></v-btn>
          </v-row>
          <v-btn class="flip" @click="flip">Flip Board <v-icon>fas fa-retweet</v-icon> </v-btn>
          <v-tabs fixed-tabs>
            <v-tab>
              Chat 
            </v-tab>
            <v-tab-item><chat :username="username" :roomId="roomId"/></v-tab-item>
            <v-tab>
              Members
            </v-tab>
            <v-tab-item><members :username="username" :members="members" :players="playerList"/></v-tab-item>
            <v-tab v-if="movePatterns.length!=0">
              Move Pattern
            </v-tab>
            
            <v-tab-item v-if="movePatterns.length!=0"><move-pattern-tab :movePatterns="movePatterns" color="white" /></v-tab-item>
          </v-tabs>
        </div>
    </div>
  </div>
</template>

<script>
import Board from '../Board.vue'
import Chat from '../Chat/Chat.vue'
import Members from './Members.vue'
import MovePatternTab from './MovePatternTab.vue'
import WS,{sendMoveInfo} from '../../utils/websocket';
export default {
  components: { Chat,Board,Members,MovePatternTab },
  computed:{
  },
  mounted(){
    this.updatePlayerList()
    this.$store.commit("setClientInfo",{
      username:this.username,
      isPlayer: this.username==this.player1 || this.username==this.player2,
      color: this.username==this.player1 ? 'w' : this.username==this.player2 ? 'b' : null,
      ws: this.ws
      })

    this.$store.subscribe((mutation, state) => {
      if (mutation.type === "updateGameInfo") {
        this.player1 = state.gameInfo[this.roomId].p1 ? state.gameInfo[this.roomId].p1 : null;
        this.player2 = state.gameInfo[this.roomId].p2 ? state.gameInfo[this.roomId].p2 : null;
        this.updatePlayerList()
        this.isFlippedCheck()
        
      }
      else if(mutation.type === "performMove"){
        this.$refs.gameBoard.performMove(this.$store.state.currentMove)
      }
      else if(mutation.type==="websocketError"){
        this.error = state.errorMessage;
      }
    })
    
  },
  methods:{
    updatePlayerList(){
      console.log('pl',this.$store.state.gameInfo[this.roomId])
      var roomInfo = this.$store.state.gameInfo[this.roomId]
      this.player1 = roomInfo && roomInfo.p1 ? roomInfo.p1 : null;
      this.player2 = roomInfo && roomInfo.p2 ? roomInfo.p2 : null;
      this.playerList = this.player1? this.player2? [this.player1,this.player2] : [this.player1] : null
      this.members = roomInfo && roomInfo.members.filter((value)=>{
        return value!=this.player1 && value!=this.player2
      })
      console.log('members',this.members)
    },
    validateMove(destInfo){
      var srcInfo = this.$store.state.curStartPos
      var piece = this.getPlayerColor()=='w' ? srcInfo.piece.toUpperCase() : this.getPlayerColor()=='b' ? srcInfo.piece.toLowerCase() : srcInfo.piece
      var info = {roomId: this.roomId,
          piece: piece,
          srcRow: srcInfo.row,
          srcCol: srcInfo.col,
          destRow: destInfo.row,
          destCol: destInfo.col,
          color:this.getPlayerColor(),
        } 
        if(info.piece=='k'||info.piece=='K' && info.srcRow==info.destRow && Math.abs(info.srcCol-info.destCol)==2){
          info.castle=true
        }
      sendMoveInfo(this.ws,info)
    },
    getPlayerColor(){
      return this.player1 == this.username ? 'w' : this.player2 == this.username ? 'b' : null;
    },
    flip(){
      this.isFlipped=!this.isFlipped
      this.$refs.gameBoard.updateBoardState1D(this.isFlipped);
    },
    isFlippedCheck(){
      this.isFlipped = this.player2 ? this.username === this.player2 : false;
      return this.isFlipped;
    },
    getShareUrl(){ return `${window.location.origin}/#/join/${this.$route.params.roomId}`},
    copyText(){
      let input=document.getElementById("tocopy");
      input.select();
      document.execCommand("copy");
    },
    
  },
  data(){
    return{
      shareLink: this.getShareUrl(),
      username: this.$route.params.username,
      roomId: this.$route.params.roomId,
      playerList:[],
      player1: null, 
      player2: null,
      moveSrc: null,
      moveDest: null,
      isFlipped: this.isFlippedCheck(),
      movePatterns: this.$store.state.movePatterns,
      turn: null,
      members: [],//this.$store.state.roomClients[roomId],
      boardState:this.$route.params.boardState ? this.$route.params.boardState : this.$store.state.boards[this.roomId],
      ws: WS,
    }
  }
}
</script>

<style scoped>
.columns{
  display: flex;
  margin-bottom: 2em;
}
.centered-input input {
  text-align: center;
}
.copylink{
  display: flex;
  align-items: center;
  justify-content: center;
}

.board-panel{
  flex: 1;
}
.flip{margin-bottom: 1em;}
.right-panel{
   max-width:400px;
   width: 100%;
   padding: 1em;
   background-color: rgb(224, 223, 223);
   border-radius: 1%;
}

</style>