<template>
  <div class="container">
    <div class="columns">
        <div class="column board-panel">
          <svg v-if="result && result=='draw'" style="position: absolute; left:0; top:0; width: 100%; z-index:1;" viewBox="0 0 250 100">
                <rect x="0" y="35%" width="100%" height="40%" fill="rgba(0,0,0)" fill-opacity="0.7"/>
                <text x="35%" y="55%"  
                      font-family="Arial, Helvetica, sans-serif"
                      dominant-baseline="central" text-anchor="middle" fill="white">
                    DRAW!
                </text>
          </svg>
          <svg v-else-if="result" style="position: absolute; left:0; top:0; width: 100%; z-index:1;" viewBox="0 0 250 100">
                <rect x="0" y="35%" width="100%" height="40%" fill="#65b5f5" fill-opacity="0.7"/>
                <text x="35%" y="55%"  
                      font-family="Arial, Helvetica, sans-serif"
                      dominant-baseline="central" text-anchor="middle" :fill="result">
                    {{result.toUpperCase()}} WINS!
                </text>
          </svg>
          
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
          <v-row row d-flex nowrap align="center" justify="center" class="px-2 flip">
            <div v-for="button in buttons" :key="button.text">
              <v-tooltip bottom>
                <template v-slot:activator="{ on, attrs }">
                  <v-btn
                    style="margin-right:5px;"
                    :color="button.color"
                    dark
                    v-bind="attrs"
                    v-on="on"
                    @click="button.onclick"
                  >
                    <v-icon>fas {{button.icon}}</v-icon>
                  </v-btn>
                </template>
                <span>{{button.text}}</span>
              </v-tooltip>
            </div>
          </v-row>
          
          <v-tabs fixed-tabs>
            <v-tab>
              Chat 
            </v-tab>
            <v-tab-item><chat :username="username" :roomId="roomId"/></v-tab-item>
            <v-tab>
              Members
            </v-tab>
            <v-tab-item><members :username="username" :members="membersList" :players="playerList"/></v-tab-item>
            <v-tab v-if="movePatterns && movePatterns.length!=0">
              Move Pattern
            </v-tab>
            
            <v-tab-item v-if="movePatterns && movePatterns.length!=0"><move-pattern-tab :movePatterns="movePatterns" color="white" /></v-tab-item>
          </v-tabs>
        </div>
    </div>
  </div>
</template>

<script lang="ts">
import Board from '../Board/Board.vue'
import Chat from '../Chat/Chat.vue'
import Members from './Members.vue'
import MovePatternTab from './MovePatternTab.vue'
import { BoardState, GameInfo, MoveInfoPayload, MovePatterns } from '../../types';
import { Component, Vue, Prop } from 'vue-property-decorator'
import { mapActions, Action } from 'vuex';
import {namespace} from 'vuex-class';

type Button = {
  text: string;
  icon: string;
  color: string;
  onclick: () => void;
};

const webSocketModule = namespace('webSocket');

@Component({
  components:{
    Chat, 
    Board,
    Members, 
    MovePatternTab,
  },
  methods: {
    ...mapActions('webSocket', ['sendResign', 'sendDrawOffer', 'sendMoveInfo']),
  },
})
export default class GameRoom extends Vue{
  @Prop() username!: string;
  @Prop() roomId!: string;

  @webSocketModule.Action('sendMoveInfo') private sendMoveInfo!: (payload: any) => Promise<void>;
  @webSocketModule.Action('sendDrawOffer') private sendDrawOffer!: (payload: any) => Promise<void>;
  @webSocketModule.Action('sendResign') private sendResign!: (payload: any) => Promise<void>;

  private gameBoard!: Vue & { updateBoardState1D(isFlipped: boolean): void };
  shareLink: string =  this.getShareUrl();

  playerList:string[] = [];
  player1: string  = ""; 
  player2: string |null = null;
  moveSrc: number |null = null;
  moveDest: number |null = null;
  isFlipped: boolean = this.isFlippedCheck();
  movePatterns: MovePatterns = this.$store.state.movePatterns;
  turn: string | null = null;
  gameInfo: GameInfo = this.$store.state.gameInfo;
  result:string | null = null;
  buttons:Button[] = [
    {text:'Flip',icon:'fa-retweet',color:'black',onclick: this.flip},
    {text:'Resign',icon:'fa-flag',color:'red darken-1',onclick:this.resign},
    {text:'Offer Draw',icon:'fa-handshake',color:'blue darken-2',onclick:this.offerDraw}
  ];
  membersList: string[] = [];
  boardState:BoardState | null=  null;

  created() {
    this.player1 = this.$store.state.gameInfo[this.roomId];
    this.boardState = this.$store.state.boards[this.roomId];
  }

  getShareUrl(){ return `${window.location.origin}/join/${this.roomId}`}

  updatePlayerList(){
    var roomInfo = this.$store.state.gameInfo[this.roomId]
    this.player1 = roomInfo.p1;
    this.player2 = roomInfo && roomInfo.p2 ? roomInfo.p2 : null;
    this.playerList = this.player1? this.player2? [this.player1,this.player2] : [this.player1] :[];
    this.membersList = roomInfo && roomInfo.members.filter((value:string)=>{
      return value!=this.player1 && value!=this.player2
    })
  }

  validateMove(destInfo:{row:number, col:number}){
    let srcInfo = this.$store.state.curStartPos;
    let piece = this.getPlayerColor()=='white' ? srcInfo.piece.toUpperCase() : this.getPlayerColor()=='black' ? srcInfo.piece.toLowerCase() : srcInfo.piece
    let info:MoveInfoPayload = {roomId: this.roomId,
        piece: piece,
        srcRow: srcInfo.row,
        srcCol: srcInfo.col,
        destRow: destInfo.row,
        destCol: destInfo.col,
        color:this.getPlayerColor()![0],
      } 
    if((info.piece=='k'||info.piece=='K') && info.srcRow==info.destRow && Math.abs(info.srcCol-info.destCol)==2){     
      info.castle=true
    } 
    this.sendMoveInfo(info)
  }

  getPlayerColor(){
    return this.player1 == this.username ? 'white' : this.player2 == this.username ? 'black' : null;
  }

  getOpponentColor(){
    var plColor = this.getPlayerColor();
    return plColor == 'white' ? "black" : plColor!=null ? "white": null;
  }

  flip(){
    this.isFlipped=!this.isFlipped
    this.gameBoard.updateBoardState1D(this.isFlipped);
  }

  resign(){
    this.sendResign({roomId:this.roomId,color: this.getPlayerColor()![0]})
  }

  offerDraw(){
    this.sendDrawOffer({roomId:this.roomId,color: this.getPlayerColor()![0]})
  }

  isFlippedCheck(){
    this.isFlipped = this.player2 ? this.username === this.player2 : false;
    return this.isFlipped;
  }

  copyText(){
    let input:HTMLInputElement =document.getElementById("tocopy") as HTMLInputElement;
    input.select();
    navigator.clipboard.writeText(input.value);
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
  flex:2;
  padding: 1em;
}
.flip{margin-bottom: 1em;}
.right-panel{
   flex:1;
   padding: 1em;
   background-color: rgb(224, 223, 223);
   border-radius: 1%;
}

@media only screen and (max-width: 700px) {
  .columns{
    display: flex;
    flex-flow: column;
    margin-bottom: 2em;
  }
}

</style>