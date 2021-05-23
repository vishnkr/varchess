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
      <router-view/>
    </v-main>
  </v-app>
</template>

<script>

export default {
  name: 'App',

  data: () => ({
    //
  }),
};
</script>
