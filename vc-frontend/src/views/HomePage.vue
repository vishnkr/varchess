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
        v-model="dialog"
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
            
          </div>
        </template>
        <template >
          <v-card>
            <v-toolbar
              color="primary"
              dark
            ><h3>Choose Game Setup</h3></v-toolbar>
              <span class="errorText" v-if="errorText">{{errorText}}</span>
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
                @click="dialog=false"
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
    <div class="portal red"></div>
  </div>
  <div class="portal blue"></div>
  
  </div>
</template>

<script>
import WS,{createRoom} from '../utils/websocket';
import axios from 'axios';
import { convertFENtoBoardState } from '../utils/fen';


export default {
  components:{},
  props:['shared'],
  mounted: function() {
    this.$store.commit('resetState');
      
      },
  methods:{
    checkUsername(){
      if(!this.username || this.username==''){
        this.errorText = 'Enter Username';
        this.dialog=false;
        return false;
      }
      else{
        this.errorText = null;
        this.dialog=true;
        return true;
      } 
    },
    connectToWebsocket() {
      this.ws = WS;
      
    },
    isEven(val){return val%2==0},
    isLight(row,col){
        return this.isEven(row)&&this.isEven(col)|| (!this.isEven(row)&&!this.isEven(col))},

    getStandardBoard(){
      return {...convertFENtoBoardState("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w - - 0 1"), rows:8,cols:8}
    },
    async enterRoom(){
      if(this.username){
        await axios.post('http://localhost:5000/getRoomId')
        .then((response) => {
          this.roomId = response.data.data;
          this.connectToWebsocket()
          if(this.mode=='custom'){
            this.$router.push({name:'Editor',params:{username: this.username,roomId: this.roomId, ws: this.ws}})
          }else{
            createRoom(this.ws,this.roomId,this.username, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w - - 0 1");
            this.$router.push({name:'Game', params:{username: this.username,roomId: this.roomId, ws: this.ws, boardState: this.getStandardBoard()}})
          }
        }, (error) => {
          this.errorText = 'Server Not Responding'
          this.dialog=false
          console.log(error);
        });
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
      dialog: false,
      serverUrl: "ws://localhost:5000/ws",
      roomId: null,
    }
  }
}
</script>

<style scoped>
/*
.portal {
  background-color: black;
  border-radius: 44px/62px;
  box-shadow: 0 0 15px 4px white;
  height: 72px;
  width: 48px;
}

.portal.red {
  background: radial-gradient(#000000, #000000 50%, #ff4640 70%);
  border: 5px solid #ff4640;
  transform: translate3D(586px, 25px, 4px) skewX(-15deg);
}
.portal.blue {
  background: radial-gradient(#000000, #000000 50%, #258aff 70%);
  border: 5px solid #258aff;
  transform: translate3D(586px, 25px, 4px) skewX(-15deg);
}*/
</style>