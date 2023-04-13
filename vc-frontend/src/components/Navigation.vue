<template>
        <div 
            class="navbar-brand nav-left"
            :class="[!showButton ? 'lite' : null]"
            @click="redirectToHome"> 
                Varchess
            </div>
            <div class="nav-right">
                <quick-play-button v-if="showButton"/>
                <div name="Login" v-if="showButton">
                    <q-btn outline class="routerlink" @click="redirectToLogin" style="color: goldenrod; " label="Login" />
                </div>
            </div>
</template>

<script lang="ts">

import QuickPlayButton from './Other/QuickPlayButton.vue'
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router';
export default{
    components:{
        QuickPlayButton
    },    
    setup(){
        const router = useRouter();
        const route = useRoute();
        const showButton = computed(()=> route.name === 'Home');
        const redirectToHome = ()=> router.push({name:'Home'})
        const redirectToLogin = ()=> router.push({name:'Login'})
        return{
            showQuickPlayDialog: ref(false),
            redirectToHome,
            redirectToLogin,
            showButton,
        }
    }
}
</script>

<style>

.navbar-brand {
  font-family: 'Bangers', sans-serif;
  font-size: 3rem;
  margin: 0 0 0 2rem;
  color: white;
  text-shadow: -2px 2px black;
  cursor: pointer;
}

.lite{
    font-size: 2rem !important;
    margin: 0 0 0 0 2rem !important;
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