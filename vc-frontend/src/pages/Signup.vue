<template>
    <q-page class="bg-light-green window-height window-width row justify-center items-center">
    <div class="column">
      <div class="row">
        <h5 class="text-h5 text-white q-my-md">Create an Account</h5>
      </div>
      <div class="row">
        <q-card square bordered class="q-pa-lg shadow-1">
          <q-card-section>
            <q-form class="q-gutter-md">
              <q-input square filled clearable v-model="username" type="email" label="email" />
              <q-input square filled clearable v-model="password" type="password" label="password" />
            </q-form>
          </q-card-section>
          <q-card-actions class="q-px-md">
            <q-btn unelevated color="light-green-7" size="lg" class="full-width" label="Login" />
          </q-card-actions>
          <q-card-section class="text-center q-pa-none">
            <p class="text-grey-6">Not registered? <AppLink name="SignUp"><p style="cursor: pointer;">Create an Account</p></AppLink></p>
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
        const q = useQuasar()
        const username = ref(null)
        const password = ref(null)
        
        return {
            username,
            password,
            onSubmit(){
                q.notify({
                    color: 'green-4',
                    textColor: 'white',
                    icon: 'cloud_done',
                    message: 'Submitted'
                })
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

.card{
    
}
</style>