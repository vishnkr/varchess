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
import WS, {sendJSONReq} from '../../utils/websocket';
export default {
    mounted: function() {
      
      },
    methods:{
        checkRoom(){
            this.connectToWebsocket()
            console.log(this.ws)
            sendJSONReq(this.ws,'joinRoom',this.roomId);
            this.$router.push({path:`/game/${this.username}/${this.roomId}`})
        },
        connectToWebsocket() {
        this.ws = WS
        this.ws.addEventListener('open', (event) => { this.onWebsocketOpen(event) });
        },
        onWebsocketOpen() {
        console.log("connected to WS!");
        
        },
    },
    data(){
        return {
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