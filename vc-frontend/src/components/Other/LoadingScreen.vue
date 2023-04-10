<template>
    <q-dialog persistent v-model="loading">
    <q-card dark >
      <q-card-section class="text-center">
            <div class="text-h5">Share the link below to invite friends to join the room</div>
            <div class="q-pa-md">
              <div class="copy-link q-gutter-y-md row " style="max-width: 400px">
                <q-input dark filled readonly :model-value="shareLink" />
                  <q-field borderless class="no-padded-control">
                    <q-btn color="blue" @click="copyToClipboard">
                        <font-awesome-icon icon="fa-link" />
                    </q-btn>
                  </q-field>
              </div>
            </div>
            <p class="mt-2 text-h6">Waiting for opponent...</p>

            <q-spinner-cube
            color="orange"
            size="lg"
            />
        </q-card-section>
        
      <q-card-actions align="right">
        <q-btn  @click="closePopup" color="negative">Cancel</q-btn>
      </q-card-actions>
    </q-card>
  </q-dialog>
        
</template>

<script lang="ts">
import { RootState } from '@/store/state';
import { SET_PLAYERS } from '../../store/mutation_types';
import { ref } from 'vue';
import { useRoute,useRouter } from 'vue-router';
import { useStore } from 'vuex';

export default {
  emits: ['update-loading'],
  props:{
    shareLink: {type:String,required:true},
    roomId: {type:String,required:true},
    username: {type:String,required:true}
  },
  setup(props, { emit }) {
    const loading = ref(true);
    const router = useRouter();
    
    const store = useStore<RootState>();
    store.subscribe((mutation,state)=>{
      if(mutation.type === SET_PLAYERS){
        if(store.state.gameInfo?.players.p2){
          router.push({
          name: 'Game',
          params: {
            username: props.username,
              roomId: props.roomId,
          },
        });
        }
      }

    })
    return {
      loading, 
      closePopup: ()=>{
        emit('update-loading')
        loading.value = false
      },
      copyToClipboard : ()=>{
        navigator.clipboard.writeText(props.shareLink);
      }     
    };
  },
  
};
</script>

<style scoped>
.copy-link{
    display:flex;   
    justify-content: center;
}
</style>