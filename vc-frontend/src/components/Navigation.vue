<template>
    <nav class="navbar-bg">
            <AppLink 
            class="navbar-brand nav-left"
            name="Home"> 
                Varchess
            </AppLink>
            <div class="nav-right">
                <QuickPlayButton />
                <AppLink name="Login">
                    <q-btn outline class="routerlink" style="color: goldenrod; " label="Login" />
                </AppLink>
            </div>
            
    </nav>
</template>

<script lang="ts">

import AppLink from '../components/AppLink.vue'
import QuickPlayButton from './Other/QuickPlayButton.vue'
import { ref, inject } from 'vue'
export default{
    components:{
        QuickPlayButton
    },
    setup(){
        return{
            showQuickPlayDialog: ref(false)
        }
    }
}
</script>

<style>

.navbar-brand {
  font-family: 'Bangers', sans-serif;
  font-size: 3rem !important;
  padding-top: 0rem !important;
  margin: 0 0 0 2rem !important;
  color: white !important;
  text-shadow: -2px 2px black;
}

nav{
    overflow: hidden;
}
.navbar-bg{
    display: flex;
    justify-content: flex-start;
    align-items: center;
    background-color: #171717;
}
.nav-left, .nav-right{
    display: flex;
    flex:1;
}
.nav-right{
    justify-content: center;
    align-items: center;
}
.nav-right>*{
    margin: 10px;
}

</style>