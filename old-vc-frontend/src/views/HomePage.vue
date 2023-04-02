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
            ><h3>{{loading ? 'Waiting for another player to join..' : 'Choose Game Setup'}}                                                                                                                                                        </h3></v-toolbar>
            <v-card-actions class="justify-start">
              <div v-if = "!loading">
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
              </div>  
              <div v-if="loading" class="text-center">
                <h3>Share the link below to invite friends to join the room</h3>
                <v-progress-circular :size="60" indeterminate color="primary" />
                <v-spacer />
                <v-row row d-flex nowrap align="center" justify="center" class="px-2">
                  <v-text-field width="10px" class="centered-input" v-model="shareLink" id="tocopy" readonly></v-text-field>
                  <v-btn width="6.5rem" class="ma-2" rounded dark color="blue" @click="copyText">Copy<v-icon>fas fa-link</v-icon></v-btn>
                  <v-btn color="error" @click="handleCancel">Cancel</v-btn>
                </v-row>
              </div>
              <v-spacer v-if="!loading" />
              <v-btn v-if="!loading"
                color="error"
                @click="handleCancel"
              >Cancel</v-btn>
              <v-btn v-if = "!loading"
                color="success"
                @click="enterRoomWithLoading"
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

<script lang="ts">
import { GAME_MODES, GameMode } from '../utils/constants';
import { faInfoCircle } from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import Vue from 'vue';
import { mapActions, mapMutations } from 'vuex'
import { convertFENtoBoardState } from '../utils/fen';
import LoadingScreen from '../components/Other/LoadingScreen.vue';
import { SET_PLAYERS, SET_SERVER_STATUS, UPDATE_BOARD_STATE } from '../utils/mutation_types';

interface Data {
  createClicked: boolean;
  errorText: string | null;
  mode: string;
  username: string | null;
  dialog: boolean;
  standardFen: string;
  roomId: string | null;
  server_host: string | undefined;
  game_modes: GameMode[];
  loading: boolean;
}

export default Vue.extend({
  components:{LoadingScreen},
  computed:{
    didPlayerJoin(){
      return false
    },
    shareLink(): string {
      return `${window.location.origin}/join/${this.roomId}`;
    },
  },
  mounted: function() {
    this.$store.commit('resetState');
    this.$store.subscribe((mutation, state) => {
       if(mutation.type=== SET_SERVER_STATUS){
        if (state.serverStatus.errorMessage) {
          this.errorText = state.serverStatus.errorMessage;
        };
      } else if(mutation.type === SET_PLAYERS){
          if (state.gameInfo.players.p2 && this.username && this.roomId){
            this.enterRoom()
          }
          
      }
     })
  },
  
  methods:{
    ...mapActions('webSocket',['connect','close']),
    ...mapActions(['createRoom']),
    ...mapMutations([UPDATE_BOARD_STATE]),
    closeWebSocket() {
      this.close()
    },
    handleCancel(){
      if (!this.loading){
        this.dialog = false;
      } else { 
        this.loading = false;
        this.close(); 
      }
    },
    copyText(){
      let input:HTMLInputElement =document.getElementById("tocopy") as HTMLInputElement;
      input.select();
      navigator.clipboard.writeText(input.value);
    },
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

    isEven(val:number){return val%2==0},
    isLight(row:number,col:number){
        return this.isEven(row)&&this.isEven(col)|| (!this.isEven(row)&&!this.isEven(col))
    },

    enterRoom() {
      this.loading = false;
      if(this.username && this.roomId){
        this.$router.push({
          name: 'Game',
          params: {
            username: this.username,
              roomId: this.roomId,
          },
        });
      }
      
    },

    async enterRoomWithLoading(){
      if (this.username) {
        if (this.mode!=='custom'){
          try{
            this.loading = true;
            this.roomId = await this.createRoom({fen:this.standardFen});
            if (this.username && this.roomId) {
              this.connect({ roomId: this.roomId, username: this.username });
              this.UPDATE_BOARD_STATE({ roomId: this.roomId, boardState: convertFENtoBoardState(this.standardFen) });
              return
            }
          } catch(error){
              this.errorText = 'Server Not Responding'
              this.dialog = false
              console.log(error);
              return
          }
        }
        this.$router.push({
          name:'Editor',
          params: {
            username: this.username,
          },
        });
      } 
    }
  },
  data():Data{
    return {
      createClicked: false,
      errorText: null,
      mode:'standard',
      username: null,
      dialog: false,
      standardFen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR",
      roomId: null,
      server_host: process.env.VUE_APP_SERVER_HOST,
      game_modes: GAME_MODES,
      loading: false,
    };
  }
});
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