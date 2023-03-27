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

<script lang="ts">
import { convertFENtoBoardState } from '../../utils/fen';
import {SET_MOVE_PATTERNS, SET_PLAYERS, UPDATE_BOARD_STATE, UPDATE_MEMBERS} from '../../utils/mutation_types'
import { mapActions, mapMutations } from 'vuex';
import Vue from 'vue';

interface Data{
  username: string | null,
  roomId?: string,
  server_host?: string
}

import { RoomState } from '../../types';
export default Vue.extend({
    created(){
      this.roomId = this.$route.params.roomId;
    },
    mounted(){
      this.$store.subscribe((mutation,state)=>{
        if(mutation.type === SET_PLAYERS){
          if (state.gameInfo.players.p2 && this.roomId && this.username){
            this.enterRoom(this.roomId,this.username)
          }
      }
      })
    },
    methods:{
      ...mapActions('webSocket',['connect','joinRoom']),
      ...mapActions(['getRoomState']),
      ...mapMutations([SET_MOVE_PATTERNS,UPDATE_BOARD_STATE,UPDATE_MEMBERS,SET_PLAYERS]),
        
      enterRoom(roomId:string,username:string){
        this.$router.push({name:'Game', params:{username: username,roomId: roomId}})
      },

      async checkJoinRoom(){
        if (this.username){  
          try{
            
            let roomState:RoomState = await this.getRoomState({roomId:this.roomId});
            if (roomState.movePatterns){
              this.SET_MOVE_PATTERNS(roomState.movePatterns)
            }
            this.SET_PLAYERS({p1:roomState.p1,p2:roomState.p2})
      
            this.UPDATE_BOARD_STATE({roomId:roomState.roomId,boardState:convertFENtoBoardState(roomState.fen)})
            this.connect({roomId:this.roomId,username:this.username});
          } catch(error){
            console.error(error)
          }
        }
      },
    },
    data():Data{
        return {
            username: null,
            roomId:undefined,
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