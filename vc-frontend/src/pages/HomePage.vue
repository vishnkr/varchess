<template>
    <q-page class="background-container">
        <SlideShow />
    </q-page>
    
</template>


<script lang="ts">
import QuickPlayDialog from '../components/Other/QuickPlayButton.vue';
import {ref, onMounted, defineComponent, inject, watch} from 'vue'
import { useRouter } from 'vue-router';
import { useStore } from 'vuex';
import SlideShow from '../components/Other/SlideShow.vue';

import { SET_SERVER_STATUS } from '../utils/action_mutation_types';

export default defineComponent({
  components:{ SlideShow, QuickPlayDialog },
  setup(){
    const store = useStore();
    const errorText = ref(null);
    const router = useRouter();
    const showDialog = inject('showQuickPlayDialog');
    function redirectToHome(){
      router.replace({path:'/'});
      store.dispatch(`webSocket/close`)
    }

    onMounted(()=>{ 

        store.dispatch('checkServerStatus');
        store.subscribe((mutation,state)=>{
        if(mutation.type=== SET_SERVER_STATUS){
          errorText.value = state.serverStatus.errorMessage ? state.serverStatus.errorMessage : null;
      }
      })
    });

    return {
      tab:ref(''),
      redirectToHome,
      showDialog
    };
  }
})
</script>

<style>
 .background-container {
    display:flex;;
    height: 100%;
    background-color: #333333;
}

.logo-font {
  font-family: 'Titillium Web', sans-serif;
}

</style>