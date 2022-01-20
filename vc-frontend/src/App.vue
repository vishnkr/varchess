<template>
  <v-app>
    <v-app-bar
      app
      color="blue lighten-2"
      dark
    >
        <div class="title" v-on:click="redirectToHome">
        <v-img src="./assets/logo.svg" max-height="50" max-width="50" contain />
        <div class="d-flex align-center">
          <h1>VarChess</h1>
        </div>
        </div>
      <v-spacer></v-spacer>
      <div v-if="!isAuthenticated">
      <v-btn rounded v-on="on" v-bind="attrs" color="primary" dark @click="signup">Sign up</v-btn>
      <v-btn rounded
              color="primary"
              v-bind="attrs"
              v-on="on"
              dark 
              @click="login">Log in</v-btn>
      </div>
      
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
          this.errorText = state.errorMessage ? state.errorMessage : null;
      }
     })
  },
  data(){
    return {
      errorText: null,
      isAuthenticated:true,
    }
  },
  methods:{
    redirectToHome(){
    this.$router.push({name:'Home'})
    location.reload()
    },
    login(){
      this.$router.push({name:'Login'})
    },
    signup(){
      this.$router.push({name:'Signup'})
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