<template>
  <div class="chat">
      <div class="card">
          <v-card elevation="2">
              <p class="nomessages">
                  <span>Chat</span>
              </p>
              <div class="messages" v-chat-scroll="{always:false, smooth:true}">
                  <div v-for="message in messages" :key="message.id">
                      <span class="text-info">[{{message.username}}]:</span>
                      <span>{{message.message}}</span>
                      
                  </div>
              </div>
              <div class="card-action">
                  <create-message v-on:sendChatMessage="setChatMessage" :username="username" />
              </div>
          </v-card>
      </div>
  </div>
</template>

<script>
import CreateMessage from '../Chat/CreateMessage';
//import moment from 'moment';<span class="text-secondary time">{{message.timestamp}}</span>

export default {
    components: {CreateMessage},
    props: ['username'],
    data(){
        return{
            messages: [],
        }
    },
    created(){
        /*var id=0;
        var user='Varchess bot';
        var message=" Welcome to your chatroom!";
            this.messages.push({
                id: id,
                username: user,
                message: message
            })
            */
    },
    methods:{
        setChatMessage(messageInfo){
            this.messages.push(messageInfo)
        }
    }
}
</script>

<style scoped>
.chat{
    margin: 2%;
}
.chat span{
    font-size: 1.2em;
}

.messages{
    max-height: 300px;
    overflow: auto;
    margin: 30px
}
.nomessages{
    text-align: center;
}

</style>