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
              @click="checkJoinRoom"
              >
              Enter Room
            </v-btn>
      </div>
  </div>
</template>

<script>
import { convertFENtoBoardState } from '../../utils/fen';
import { mapActions } from 'vuex';
import axios from 'axios';
import Vue from 'vue';
export default Vue.extend({
    mounted: function() {
       this.getBoardFen();
      },
    methods:{
      ...mapActions('webSocket',['connect','joinRoom']),
      async getBoardFen() {
        await axios.get(`${this.server_host}/board-fen/${this.roomId}`).then((response)=>{
          if(response.data.type && response.data.type=="error"){
            this.$store.commit('websocketError',response.data.data)
          } else {
          if(response.data.movePatterns){
            this.$store.commit('storeMovePatterns',{movePatterns:response.data.movePatterns})
          }
          this.$store.commit('updateBoardState',{roomId:this.roomId,boardState:convertFENtoBoardState(response.data.Fen)});
          }
        });
      },

      checkJoinRoom(){
          this.connect();
          this.joinRoom({roomId:this.roomId,username:this.username});
          this.$router.push({name:'Game', params:{username: this.username,roomId: this.roomId}})
        },
    },
    data(){
        return {
            username: null,
            roomId: this.$route.params.roomId,
            server_host: process.env.VUE_APP_SERVER_HOST,
        }
    }
});
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