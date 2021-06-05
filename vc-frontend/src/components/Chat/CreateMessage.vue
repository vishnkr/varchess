<template>
  <div class="container" style="margin-bottom: 30px">
          <div class="form-group">
              <v-text-field v-model="newMessage" label="Enter Message..." solo></v-text-field>
              <p class="error-text" v-if="errorText">{{errorText}}</p>
          </div>
          <div><v-btn depressed color="primary" class="button" @click="createMessage">Submit</v-btn></div>
          

  </div>
</template>

<script>

export default {
    props: ['name'],
    data(){
        return{
            newMessage: null,
            errorText: null,
        }
    },
    methods:{
        createMessage(){
            console.log(this.newMessage)
            if(this.newMessage){
                var newMessage = {message: this.newMessage, name: this.name};
                this.$emit('sendChatMessage',newMessage);
                this.newMessage=null;
                this.errorText=null;
            }
            else{
                this.errorText = 'A message must be entered first';
            }
        }
    }

}
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