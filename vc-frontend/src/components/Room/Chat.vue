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
    props: ['username','roomId'],
    data(){
        return{
            ws: WS,
            messages: [
            ],
            count: 0,
        }
    },
    
    mounted(){
           this.ws.onmessage = (msg)=>{
               console.log(msg)
                let apiMsg = JSON.parse(msg.data);
                if (apiMsg.type==="chatMessage" ){
                    
                    let msgData = apiMsg.data//JSON.parse(apiMsg.data);
                    console.log("parsed: ",msgData)
                    msgData.id = this.count;
                    this.count+=1;
                    //console.log('messages',this.messages,msgData);
                    this.messages.push(msgData);
                    
                }
                //let msgData = JSON.parse(apiMsg);
                //this.messages.push({username:msgData.username, message:msgData.message})
           }
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