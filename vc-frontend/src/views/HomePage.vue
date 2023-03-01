<template>
  <div class="home">
    <div class="info-card">
      <v-card class="mx-auto"
        max-width="500">
        <v-card-text>
          <p class="display-1 text--primary">
            VarChess
          </p>
          <div class="text-h5 mb-2">
             Create your own chess variants and play with friends!
          </div>
          <ul>
            
           <li> Feeling creative? Play on a 12x7 board with a "dolphin" piece that can move like a knight and half a bishop.</li>
           <li> Are you an aggressive player? Fill the board with as many queens as you wish</li>
           
          </ul>
          </v-card-text>
          <v-img src="../assets/game.png" max-height="" />
      </v-card>
    </div>
    <div class="create-room">
      <v-card class="mx-auto"
    max-width="344">
    <v-card-text>
      
      <v-container>
        <v-row>
        <v-col
          cols="12"
        >
        <div class="headline text--primary">Create a room now to play!</div>
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
            <v-card-actions class="justify-start">
              <v-radio-group
                    v-model="mode"
                    row
                    dense
                  >
                    <v-radio v-for = "gameMode in game_modes" :key="gameMode.key" :label="gameMode.name" :value="gameMode.key">
                      <template #prepend-inner>
                        <font-awesome-icon icon="info-circle" />
                      </template>
                    </v-radio>
                  </v-radio-group>
                <v-spacer></v-spacer>
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
  </div>
  
  </div>
</template>

<script>
import WS,{createRoom} from '../utils/websocket';
import axios from 'axios';
import { convertFENtoBoardState } from '../utils/fen';
import { GAME_MODES } from '../utils/constants';
import { faInfoCircle } from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

export default {
  components:{},
  props:['shared'],
  mounted: function() {
    this.$store.commit('resetState');
    this.$store.subscribe((mutation, state) => {
       if(mutation.type==="websocketError"){
        this.errorText = state.errorMessage;
      }
     })
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
      return {...convertFENtoBoardState(this.standardFen), rows:8,cols:8}
    },
    async enterRoom(){
      if(this.username){
        await axios.post(`${this.server_host}/room-id`)
        .then((response) => {
          this.roomId = response.data.data;
          this.connectToWebsocket()
          if(this.mode=='custom'){
            this.$router.push({name:'Editor',params:{username: this.username,roomId: this.roomId, ws: this.ws}})
          }else{
            createRoom(this.ws,this.roomId,this.username, this.standardFen);
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
      standardFen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR",
      roomId: null,
      server_host: process.env.VUE_APP_SERVER_HOST,
      game_modes: GAME_MODES
      
    }
  }
}
</script>

<style scoped>
ul{
  color:rgb(6, 6, 180); 
  font-family:Arial; 
  font-size: 18px;
}
.home{
  display:flex;
  justify-content: center;
  align-items: center;
}
.info-card{
  margin: 2em;
}

</style>