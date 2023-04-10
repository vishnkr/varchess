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

import { SET_SERVER_STATUS } from '../store/mutation_types';
import { CHECK_SERVER_STATUS } from '../store/action_types';

export default defineComponent({
  components:{ SlideShow, QuickPlayDialog },
  setup(){
    const store = useStore();
    const errorText = ref(null);
    const router = useRouter();
    function redirectToHome(){
      router.replace({path:'/'});
      store.dispatch(`webSocket/close`)
    }

    onMounted(()=>{ 

        store.dispatch(CHECK_SERVER_STATUS);
        store.subscribe((mutation,state)=>{
        if(mutation.type=== SET_SERVER_STATUS){
          errorText.value = state.serverStatus.errorMessage ? state.serverStatus.errorMessage : null;
      }
      })
    });

    return {
      tab:ref(''),
      redirectToHome,
    };
  }
})
</script>

<style>
 .background-container {
    display:grid;
    grid-template-rows: 1fr;
    height: 100%;
    background-color: #333333;
}

.logo-font {
  font-family: 'Titillium Web', sans-serif;
}

</style>