<template>
  <div class="home">
    <div class="info-card">
      <v-card class="mx-auto"
    max-width="344">
    <v-card-text>
      <p class="display-1 text--primary">
        VarChess
      </p>
      <v-container>
        <v-row>
        <v-col
          cols="12"
        >
      <div class="headline">Create your own chess variants and play with friends! The perfect place to experiment your playing style with varying board sizes, customized piece movements and placements. More features coming soon.</div>
      </v-col>
      <v-col
          cols="12"
        >
          <v-text-field
            v-model="username"
            label="Enter Username"
            solo-inverted
            
          ></v-text-field>
        </v-col>
      </v-row>
      <v-dialog
        transition="dialog-bottom-transition"
        max-width="600"
      >
        <template v-slot:activator="{ on, attrs }">
          <div class="text-center">
            <v-btn
              rounded
              color="primary"
              v-bind="attrs"
              v-on="on"
              dark 
              @click="checkUsername"
              >
              Create Room
            </v-btn>
            <span v-if="errorText">{{errorText}}</span>
          </div>
        </template>
        <template v-slot:default="dialog">
          <v-card>
            <v-toolbar
              color="primary"
              dark
            ><h3>Choose Game Setup</h3></v-toolbar>
             
            <v-card-actions class="justify-end">
              <v-radio-group
                    v-model="mode"
                    row
                  >
                    <v-radio
                      label="Standard 8x8 board"
                      value="standard"
                    ></v-radio>
                    <v-radio
                      label="Custom Variant"
                      value="custom"
                    ></v-radio>
                  </v-radio-group>
              <v-btn
                color="error"
                @click="dialog.value=false"
              >Cancel</v-btn>
              <v-btn
                color="success"
                @click="enterRoom"
              >Enter Room</v-btn>
            </v-card-actions>
          </v-card>
        </template>
        </v-dialog>
      </v-container>
    </v-card-text>
    </v-card>
  </div>
  </div>
</template>

<script>

import {sendJSONReq} from '../utils/websocket';
import axios from 'axios';

export default {
  components:{},
  props:['shared'],
  mounted: function() {
      this.connectToWebsocket()
      },
  methods:{
    checkUsername(){
      if(!this.username || this.username==''){
        this.errorText = 'Enter Username';
        //dialog.value=false;
        return false;
      }
      else{
        this.errorText = null;
        return true;
        //dialog.value=true;
      } 
    },
    connectToWebsocket() {
      this.ws = new WebSocket( this.serverUrl );
      this.ws.addEventListener('open', (event) => { this.onWebsocketOpen(event) });
    },
    onWebsocketOpen() {
      console.log("connected to WS!");
      
    },
    async enterRoom(){
      //console.log(this.roomId);
      //sendReq(this.ws,'getRoomId');
      if(this.username){
        await axios.post('http://localhost:5000/getRoomId')
        .then((response) => {
          console.log('response http:',response);
          this.roomId = response.data.data;
        }, (error) => {
          console.log(error);
        });
        sendJSONReq(this.ws,'createRoom',this.roomId);
        this.$router.push({path: (this.mode=='custom'? `/editor/${this.username}/${this.roomId}` : `/game/${this.username}/${this.roomId}`) })
      }
    },

  },
  data:()=>{
    return {
      createClicked: false,
      errorText: null,
      mode:'standard',
      username: null,
      ws: null,
      dialog:true,
      serverUrl: "ws://localhost:5000/ws",
      roomId: null,
    }
  }
}
</script>

<style scoped>

</style>