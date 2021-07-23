<template>
  <div class="chat">
      <div class="card">
          <v-card>
              <p class="nomessages">
              </p>
              <div class="messages" v-chat-scroll="{always:false, smooth:true}">
                  <div v-for="message in messages" :key="message.id">
                      <span v-if="message.username" class="text-info">[{{message.username}}]:</span>
                      <span>{{message.message}}</span>
                  </div>
              </div>
              <div class="card-action">
                  <create-message v-on:sendChatMessage="setChatMessage" :username="username" :roomId="roomId" :ws="ws" />
              </div>
          </v-card>
      </div>
  </div>
</template>

<script>
import CreateMessage from '../Chat/CreateMessage';
import WS from '../../utils/websocket';
export default {
    components: {CreateMessage},
    computed:{
        newMessages(){
            return this.$store.state.chatMessages[this.roomId] ? this.$store.state.chatMessages[this.roomId] : []
        }
    },
    props: ['username','roomId'],
    mounted() {
    this.$store.subscribe((mutation, state) => {
      if (mutation.type === 'addMessage' || mutation.type === 'addClientToRoom' || mutation.type === 'removeClientFromRoom') {
        this.messages = state.chatMessages[this.roomId]
      } 
    })
    },
    data(){
        return{
            ws: WS,
            messages: [],
            count: 0,
        }
    },
    methods:{
        setChatMessage(messageInfo){
            messageInfo.roomId = this.roomId;
            this.$store.commit('addMessage',messageInfo);
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
    height: 300px;
    overflow: auto;
    margin: 30px
}
.nomessages{
    text-align: center;
}

</style>