<template>
  <div class="join">
      <div class="inp">
      Join Room: {{roomId}}
      <v-text-field
            v-model="username"
            label="Enter Username"
            solo-inverted
            
          ></v-text-field>
          <v-btn
              rounded
              color="primary"
              dark 
              @click="checkRoom"
              >
              Enter Room
            </v-btn>
      </div>
  </div>
</template>

<script>
import { convertFENtoBoardState } from '../../utils/fen';
import WS, {joinRoom} from '../../utils/websocket';
import axios from 'axios';
export default {
    mounted: function() {
       this.getBoardFen();
       this.connectToWebsocket()
      },
    methods:{
      async getBoardFen() {
        await axios.get(`${this.server_host}/getBoardFen/${this.roomId}`).then((response)=>{
          if(response.data.type && response.data.type=="error"){
            this.$store.commit('websocketError',response.data.data)
          } else {
          if(response.data.movePatterns){
            this.$store.commit('storeMovePatterns',{movePatterns:response.data.movePatterns})
          }

          this.boardState = convertFENtoBoardState(response.data.Fen)
          this.boardState.rows  = this.boardState.tiles.length
          this.boardState.cols = this.boardState.tiles[0].length
          var id =0;
          for(var row=0;row<this.boardState.rows;row++){
            for(var col=0;col<this.boardState.cols;col++){
              this.boardState.tiles[row][col].tileId = id;
              id+=1;
            }
          }
          }
        });
      },

        checkRoom(){
            this.connectToWebsocket()
            joinRoom(this.ws,this.roomId,this.username);
            this.$router.push({name:'Game', params:{username: this.username,roomId: this.roomId, boardState: this.boardState, ws:this.ws}})
            //this.$router.push({path:`/game/${this.username}/${this.roomId}`})
        },
        connectToWebsocket() {
        this.ws = WS
        
        },
    },
    data(){
        return {
            boardState: null,
            username: null,
            roomId: this.$route.params.roomId,
            ws:null,
            server_host: process.env.VUE_APP_SERVER_HOST,
        }
    }
}
</script>

<style>
.inp{
  position: fixed;
  top: 50%;
  left: 50%;
  /* bring your own prefixes */
  transform: translate(-50%, -50%);

}
</style>