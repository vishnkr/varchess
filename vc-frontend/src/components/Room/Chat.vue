<template>
  <div class="chat">
      <div class="card">
          <v-card elevation="2">
              <p class="text secondary nomessages" v-if="messages.length==0">
                  [No messages yet!]
              </p>
              <div class="messages" v-chat-scroll="{always:false, smooth:true}">
                  <div v-for="message in messages" :key="message.id">
                      <span class="text-info">[{{message.name}}]:</span>
                      <span>{{message.message}}</span>
                      
                  </div>
              </div>
              <div class="card-action">
                  <create-message v-on:sendChatMessage="setChatMessage" :name="uname" />
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
    props: ['name'],
    data(){
        return{
            messages: [],
            uname: 'hello',
        }
    },
    created(){
        var id=0;
        var user='a';
        var message="hi1 this";
            this.messages.push({
                id: id,
                name: user,
                message: message
            })
    },
    methods:{
        setChatMessage(messageInfo){
            this.messages.push(messageInfo)
        }
    }
}
</script>

<style scoped>
.chat span{
    font-size: 1.2em;
}

.messages{
    max-height: 300px;
    overflow: auto;
    margin: 30px
}

</style>