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
import { ref } from 'vue';

export default {
  emits: ['update-loading'],
  props:{
    shareLink: {type:String,required:true}
  },
  setup(props, { emit }) {
    const loading = ref(true);
    
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