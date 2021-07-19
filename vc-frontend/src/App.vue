<template>
  <v-app>
    <v-app-bar
      app
      color="dark-grey"
      dark
    >
      <div class="d-flex align-center">
        <h2>VarChess</h2>
      </div>

      <v-spacer></v-spacer>

      
    </v-app-bar>

    <v-main>
       <v-alert v-if="errorText"
                  border="right"
                  colored-border
                  type="error"
                  elevation="2"
                >
                  {{errorText}}
    </v-alert>
      <router-view/>
    </v-main>
    <v-footer padless>
    <v-col
      class="text-center"
      cols="12"
    >
      <a href="https://github.com/vishnkr/varchess" data-size="large">
       <v-icon>fab fa-github</v-icon>
      </a>
      <span> Licensed by <a href="https://www.gnu.org/licenses/gpl-3.0.en.html">GNU GPLv3</a>.</span>
      
       <strong>Varchess</strong> by vishnkr - {{ new Date().getFullYear() }} 
      
    </v-col>
  </v-footer>
  </v-app>
</template>

<script>
import WS from '../src/utils/websocket'
export default {
  name: 'App',
  mounted(){
    this.$store.subscribe((mutation, state) => {
       if(mutation.type==="websocketError"){
         console.log('rech')
        this.errorText = state.errorMessage;
      }
     })
  },
  data(){
    return {
      ws : WS,
      errorText: null,
    }
  },
};
</script>

<style scoped>
a:-webkit-any-link{
text-decoration:none !important;
}
</style>