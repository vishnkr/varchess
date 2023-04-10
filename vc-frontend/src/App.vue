<script setup lang="ts">
  import { ref } from 'vue';
import { Store, useStore } from 'vuex';
import Navigation from './components/Navigation.vue';
import { RootState } from './store/state';
import { SET_SERVER_STATUS } from './store/mutation_types';
  
  const store: Store<RootState> = useStore();
  const errorText = ref<String|null>(null);
  const clearError = ()=>{
    errorText.value = null
  }
  store.subscribe((mutation,state)=>{
    if(mutation.type===SET_SERVER_STATUS && state.serverStatus.errorMessage){
      errorText.value = state.serverStatus.errorMessage
    }
  });
</script>

<template>
  <q-layout view="hHh lpr fFf">
    <q-header class="navbar-bg">
        <Navigation />
      </q-header>
      
      <q-page-container class="bg-dark">
        <q-banner v-if="errorText" inline-actions class="text-white bg-red">
          {{errorText}}
          <template v-slot:action>
            <q-btn flat color="white" label="Dismiss" @click="clearError" />
          </template>
        </q-banner>
            <RouterView />
      </q-page-container>
    <q-footer class="footer">
        <a href="https://github.com/vishnkr/varchess" data-size="large">
          <font-awesome-icon icon="fa-brands fa-github" />
        </a>
        <span> Licensed by <a href="https://www.gnu.org/licenses/gpl-3.0.en.html">GNU GPLv3</a>.
        <strong>Varchess</strong> by vishnkr - {{ new Date().getFullYear() }} </span>
        <span> Made with <font-awesome-icon icon="fa-solid fa-heart" style="color: #f2021a;" /> </span> 
  </q-footer>
</q-layout>
</template>

<style scoped>


a{
  text-decoration: none;
}
.navbar-bg{
    display: flex;
    justify-content: flex-start;
    align-items: center;
    background-color: #171717 !important;
}

.footer {
  background-color: #171717;
  font-size: medium;
  text-align: center;
  align-items: center;
  justify-content: center;
  position: relative;
}

.footer>*{
  margin-right: 20px;
}
</style>
