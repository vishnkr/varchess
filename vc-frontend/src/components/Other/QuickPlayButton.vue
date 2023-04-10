<template>
    <q-btn class="play-btn" label="Quick Play" @click="showQuickPlayDialog = true" />
    <div class="q-pa-md q-gutter-sm">
        <q-dialog v-model="showQuickPlayDialog">
            <q-card dark rounded style="border-radius: 10px;">
                <q-card-section align="center">
                    <div class="text-h4">Quick Play</div>
                </q-card-section>
                <q-card-section class="q-pt-none">
                  <div class="text-h6">Create a room now to play with friends!</div>
                  <div class="text-h6">NOTE: QuickPlay Variants are limited to 8x8 boards or smaller. Login to play on larger boards.</div>
                </q-card-section>
                <q-card-section class="bg-white">
                  <q-input standout="bg-black text-white" v-model="username" label="username" />
                  <q-card-actions align="right" >
                    <q-btn flat class="bg-green text-white" label="Create Room" @click="enterEditor"></q-btn>
                  </q-card-actions>
                </q-card-section>
            </q-card>
        </q-dialog>
    </div>

</template>
<script lang="ts">
import { RootState } from '@/store/state';
import { SET_USER_INFO } from '../../store/mutation_types';
import { Ref, ref } from 'vue';
import { useRouter } from 'vue-router';
import { useStore } from 'vuex';

export default{
  setup(){
    const showQuickPlayDialog = ref(false);
    const username : Ref<string|null> = ref(null);
    const errorText : Ref<string|null> = ref(null);
    const router = useRouter();
    
    const store = useStore<RootState>();
    const checkUsername = () =>{
        if (!username || username.value === ''){
          errorText.value = 'Enter Username'
          showQuickPlayDialog.value = false
          return false
        } 
          errorText.value=null
          return true
    };
    return{
      showQuickPlayDialog,
      username,
      checkUsername ,
      enterEditor: ()=>{
        if (checkUsername()){
          store.commit(SET_USER_INFO,{ username:username.value,isAuthenticated:false, curGameRole: 'p1'})
          router.push({name:'Editor', params: {username:username.value}})
        }
      },
      errorText
    }
  }
}
</script>

<style scoped>
.play-btn{
  background: green; 
  color: white;
}

.dialog-card{
  width:600px; 
  max-width:60vw;
  border-radius: 20px;
}
</style>