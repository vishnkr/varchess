<template>
  <div class="signup">
       <v-card class="mx-auto"
        max-width="500">
        <v-card-title class="text-h4">Create an account</v-card-title>
        <v-form>
            <v-container>
                <v-row>
                    <v-col cols="12"><v-text-field 
                    v-model="username"
                    :rules="[rules.required]"
                    label="Username or Email" solo /></v-col>
                </v-row>
                <v-row>
                    <v-col cols="12">
                        <v-text-field :append-icon="showPass ? 'mdi-eye' : 'mdi-eye-off'" 
                        v-model="password"
                        :type="showPass ? 'text' : 'password'" 
                        class="input-group--focused"
                        hint="At least 8 characters"
                        :rules="[rules.required, rules.min]"
                        @click:append="showPass = !showPass"
                        label="Password" solo />
                        </v-col>
                </v-row>
                <v-card-actions class="justify-center">
                    <v-btn dark color="blue lighten-2" @click="createAccount">Sign up</v-btn>
                </v-card-actions>
            </v-container>
        </v-form>
        </v-card>
  </div>
</template>

<script>
import axios from 'axios';
export default {
    name: 'Signup',    
    data(){
        return{
            showPass:false,
            username:'',
            password:'',
            server_host: process.env.VUE_APP_SERVER_HOST,
            rules: {
            required: value => !!value || 'Required.',
            min: v => v.length >= 8 || 'Min 8 characters',
            //emailMatch: () => (`The email and password you entered don't match`),
            },
        }
    },
    methods:{
        async createAccount(){
            await axios.post(`${this.server_host}/signup`,JSON.stringify({
                username: this.username,
                password: this.password,
            })).then((response)=>{
                console.log(response)
                this.$router.push({name:'Login'})
        }).catch(function(error){
            console.log(error.toJSON());
        })
        }
    },
}
</script>

<style>
</style>