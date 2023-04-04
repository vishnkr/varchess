<script setup lang="ts">
  import Navigation from './components/Navigation.vue';
  
</script>

<template>
  <q-layout >
  <div class="parent-background">
      <q-header>
        <Navigation />
      </q-header>
      <q-page-container>
        
          <div class="content">
            <RouterView />
          </div>
      </q-page-container>
  </div>
  <footer class="footer">
        <a href="https://github.com/vishnkr/varchess" data-size="large">
          <font-awesome-icon icon="fa-brands fa-github" />
        </a>
        <span style="margin-right: 2rem;"> Licensed by <a href="https://www.gnu.org/licenses/gpl-3.0.en.html">GNU GPLv3</a>.
        <strong>Varchess</strong> by vishnkr - {{ new Date().getFullYear() }} </span>
        <span> Made with <font-awesome-icon icon="fa-solid fa-heart" style="color: #f2021a;" /> </span> 
  </footer>
</q-layout>
</template>

<style scoped>
html,body{
  height: 100vh;
}
#app{
  height:100%;
}
a{
  text-decoration: none;
}

.parent-background {
  display: flex;
  height:100%;
  flex-direction: column;
  background-color: #333333;
}
.content{
  height:100%;
  display:flex;
  flex-direction: column;
}

.footer {
  display: flex;
  background-color: #cdcaca;
  height: 20px;
  font-size: medium;
  text-align: center;
  align-items: center;
  justify-content: center;
  position: sticky;
}
</style>
