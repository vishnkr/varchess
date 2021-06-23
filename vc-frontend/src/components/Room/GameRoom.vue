<template>
  <div class="container">
    <div class="columns">
        <div class="column board-panel">
          <game-board :board="boardState" ref="gameBoard" :isflipped="isFlipped" :playerColor="player1 == username ? 'w' : player2 == username ? 'b' : null"/>
        </div>
        <div class="column right-panel">
          <v-text-field v-model="shareLink" id="tocopy" readonly outlined ></v-text-field>
          <v-btn @click="copyText">Copy Link <v-icon>fas fa-link</v-icon></v-btn>
          <v-btn @click="flip">Flip Board <v-icon>fas fa-retweet</v-icon> </v-btn>
          <div>Current Players: White - {{player1? player1: null}}, Black - {{player2? player2:null}}</div>
          <chat :username="username" :roomId="roomId"/>
        </div>
    </div>
  </div>
</template>

<script>
import GameBoard from '../GameBoard.vue'
import Chat from './Chat.vue'
import WS from '../../utils/websocket';
export default {
  components: { Chat,GameBoard },
  mounted(){
    this.$store.commit("setClientInfo",{
      username:this.username,
      isPlayer: this.username==this.player1 || this.username==this.player2,
      color: this.username==this.player1 ? 'w' : this.username==this.player2 ? 'b' : null,
      ws: this.ws
      })
    this.$store.subscribe((mutation, state) => {
      if (mutation.type === 'updateGameInfo') {
        this.player1 = state.gameInfo[this.roomId].p1 ? state.gameInfo[this.roomId].p1 : null;
        this.player2 = state.gameInfo[this.roomId].p2 ? state.gameInfo[this.roomId].p2 : null;
        this.isFlippedCheck()
      }
    })
    
  },
  methods:{
    getPlayerColor(){
      return this.player1 == this.username ? 'w' : this.player2 == this.username ? 'b' : null;
    },
    flip(){
      this.isFlipped=!this.isFlipped
      this.$refs.gameBoard.updateBoardState1D(this.isFlipped);
    },
    isFlippedCheck(){
      this.isFlipped = this.player2 ? this.username === this.player2 : false;
      console.log(this.username,'will have board flipped',this.isFlipped)
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
      player1: null,
      player2: null,
      moveSrc: null,
      moveDest: null,
      isFlipped: this.isFlippedCheck(),
      turn: null,
      boardState:this.$route.params.boardState ? this.$route.params.boardState : this.$store.state.boards[this.roomId],
      ws: WS,
    }
  }
}
</script>

<style>
.columns{
  display: flex;

}

.board-panel{
  flex: 2;
}
.right-panel{

   max-width:400px;
   width: 100%;
   padding: 1em;

}

</style>