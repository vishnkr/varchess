<template>
    <q-page class="bg-custom row justify-center items-center">
    <div >
      <div class="row">
        <div class="text-h4 text-white q-my-md">Login to Varchess</div>
      </div>
      <div class="row">
        <q-card rounded bordered class="q-pa-md shadow-1 card-props" >
          <q-card-section>
            <q-form class="q-gutter-md">
              <q-input square filled clearable v-model="username" type="email" label="username" />
              <q-input square filled clearable v-model="password" type="password" label="password" />
            </q-form>
          </q-card-section>
          <q-card-actions class="q-px-md">
            <q-btn unelevated color="light-green-8" size="lg" class="full-width" label="Login" />
          </q-card-actions>
          <q-card-section class="text-center q-pa-none">
            <p class="text-grey-6">Not registered? <AppLink name="Signup"><p style="cursor: pointer;">Create an Account</p></AppLink></p>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script lang="ts">
import { useQuasar } from 'quasar';
import {ref} from 'vue'
import AppLink from '../components/AppLink.vue'
export default{
    setup(){
        const username = ref(null)
        const password = ref(null)
        
        return {
            username,
            password,
            onSubmit(){
            },
            onReset () {
                username.value = null
                password.value = null
            }
        }
    }
}
</script>

<style scoped>
.bg-gradient{
    background: linear-gradient(
       128deg,
       rgba(64, 76, 255, 1) 0%,
       rgba(174, 63, 255, 1) 100%
     );
}

.bg-custom{
    background: #531ab4;
}
.card-props{
    width: 25rem;
    border-radius: 20px;
}
</style>