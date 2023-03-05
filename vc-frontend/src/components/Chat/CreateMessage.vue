<template>
  <div class="container" style="margin-bottom: 30px">
          <div class="form-group d-flex flex-row align-center">
             <v-text-field v-model="newMessage" placeholder="Type Something" @keypress.enter="createMessage"></v-text-field>
             <v-btn icon class="ml-4" @click="createMessage"><v-icon>mdi-send</v-icon></v-btn>
              <p class="error-text" v-if="errorText">{{errorText}}</p>
          </div>
          
          

  </div>
</template>

<script>
import { mapActions } from 'vuex';
import Vue from 'vue';

export default Vue.extend({
    props: ['username','roomId'],
    data(){
        return{
            newMessage: null,
            errorText: null,
        }
    },
    methods:{
        ...mapActions('webSocket',['sendMessage']),

        createMessage(){
            if(this.newMessage){
                var newmessage = {message: this.newMessage, username: this.username, roomId:this.roomId}
                this.sendMessage(newmessage);
                this.$emit('sendChatMessage',newmessage);
                this.newMessage=null;
                this.errorText=null;
            }
            else{
                this.errorText = 'A message must be entered first';
            }
        }
    }

});
</script>

<style scoped>
.form-group{
    display: flex;
    flex-direction: column;
}
.input{
    flex:3;
}
.button{
    flex:1;
}
.error-text{
    color:red;
}
</style>