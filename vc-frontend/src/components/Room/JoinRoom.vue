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
        await axios.get(`http://localhost:5000/getBoardFen/${this.roomId}`).then((response)=>{
          this.boardState = convertFENtoBoardState(response.data.Fen)
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
            serverUrl: "ws://localhost:5000/ws",
            ws:null,
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