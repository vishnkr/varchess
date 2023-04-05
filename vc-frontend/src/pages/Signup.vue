<template>
    <q-page class="bg-orange-9 row justify-center items-center">
      <div >
      <div class="row">
        <div class="text-h4 text-white q-my-md">Create an Account</div>
      </div>
      <div class="row">
        <q-card rounded bordered class="q-pa-md shadow-1 card-props" >
          <q-card-section>
            <q-form class="q-gutter-md">
              <q-input square filled clearable v-model="username" type="email" label="username" />
              <q-input square filled clearable v-model="email" type="email" label="email" />
              <q-input square filled clearable v-model="password" type="password" label="password" />
            </q-form>
          </q-card-section>
          <q-card-actions class="q-px-md">
            <q-btn unelevated color="light-green-8" size="lg" class="full-width" label="Create Account" />
          </q-card-actions>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script lang="ts">
import { useQuasar } from 'quasar';
import {ref} from 'vue'

export default{
    setup(){
        const q = useQuasar()
        const username = ref(null)
        const password = ref(null)
        const email = ref(null)

        return {
            username,
            password,
            email,
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
                email.value = null
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

.card-props{
    width: 25rem;
    border-radius: 20px;
}
</style>