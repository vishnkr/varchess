<template>
    <q-page>
        <div class="join">
        <div class="child-div text-white text-h4">
            Join Room: {{ roomId }}
        </div>
        <div class="q-gutter-md child-div" style="max-width: 300px">
            <q-input standout outline v-model="username" label-color="white" label="username" />
            <q-btn color="green-9" @click="checkJoinRoom">Join Room</q-btn>
        </div>
        
        </div>
    </q-page>
   
</template>

<script lang="ts">
import {  RoomState } from '@/types';
import { GameInfo } from '@/classes';
import { SET_MOVE_PATTERNS, SET_PLAYERS, SET_USER_INFO, UPDATE_BOARD_STATE } from '../../store/mutation_types';
import {CONNECT_WS, GET_ROOM_STATE} from '../../store/action_types';
import { convertFENtoBoardState } from '../../utils/fen';
import { ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useStore } from 'vuex';


export default{
    setup(){
        const router = useRouter();
        const route = useRoute();
        const store = useStore();
        const roomId = ref(route.params.roomId);
        const username = ref(null);

        store.subscribe((mutation,state)=>{
            if(mutation.type===SET_PLAYERS){
                if (store.state.gameInfo.players.p2 && roomId && username){
                    let gameInfo : GameInfo = store.state.gameInfo;
                    store.commit(SET_USER_INFO,{
                        username:username.value,
                        isAuthenticated:false, 
                        curGameRole: username.value === gameInfo.players.p1 ? 'p1' : username.value === gameInfo.players.p2 ? 'p2' : 'member'
                    })
                    router.push({name:'Game',params:{username:username.value,roomId:roomId.value}})
                }
            }
        })

        const checkJoinRoom = async () =>{
            if(username.value){
                try{
                    let roomState:RoomState= await store.dispatch(GET_ROOM_STATE,{roomId:roomId.value});
                    if (roomState.movePatterns){
                        store.commit(SET_MOVE_PATTERNS,roomState.movePatterns);
                    }
                    store.commit(SET_PLAYERS,{p1:roomState.p1,p2:roomState.p2});
                    store.commit(UPDATE_BOARD_STATE,{roomId:roomId.value,boardState:convertFENtoBoardState(roomState.fen)});
                    store.dispatch(CONNECT_WS,{roomId:roomId.value,username:username.value})
                } catch(error){
                    console.error(error)
                }
            }
        }

        return {
            roomId,
            username,
            checkJoinRoom
        }

        

    }
}

</script>

<style scoped>

.join{
    display:flex;
    flex-direction: column;
    justify-content: space-between;
    align-items:center;
    margin: 1%;
    height: 100%;
}
.child-div {
  flex: 1; 
  margin: 10px 0;
}

</style>