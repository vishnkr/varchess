<template>
  <v-app>
    <v-app-bar
      app
      color="green lighten-2"
      dark
    >
        <div class="title" v-on:click="redirectToHome">
        <v-img src="./assets/logo.svg" max-height="40" max-width="40" contain />
        <div class="d-flex align-center">
          <h1>VarChess</h1>
        </div>
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
export default {
  name: 'App',
  mounted(){
    
    this.$store.subscribe((mutation, state) => {
       if(mutation.type==="websocketError"){
        this.errorText = state.errorMessage;
      }
     })
  },
  data(){
    return {
      errorText: null,
    }
  },
  methods:{
    redirectToHome(){
    this.$router.push({name:'Home'})
    location.reload()
    }
  }
};
</script>

<style scoped>
a:-webkit-any-link{
text-decoration:none !important;
}

.title{
  cursor: pointer;
  display: inline-flex;
}
h1{
  text-shadow: -2px 0 yellow, 0 2px yellow, 2px 0 yellow, 0 -2px yellow;
  font-family: 'Bangers'; color: black;
  font-size:42px;
}
</style>