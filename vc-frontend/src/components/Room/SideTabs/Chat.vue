<template>
    <q-card dark class="chat">
            <p class="nomessages">
              </p>
              <div class="messages">
                  <div v-for="message in messages" :key="message.id">
                      <span v-if="message.username" class="text-info">{{message.username}}: </span>
                      <span class="message">{{message.message}}</span>
                  </div>
              </div>
              <div class="card-action q-pa-md">
                    <q-input dark filled v-model="curMessage" label="enter message"/>
                    <q-btn dark color="orange" @click="sendMessage">Send</q-btn>
              </div>
    </q-card>
</template>

<script lang="ts">
import {ref } from 'vue';
import { ChatMessage } from '@/types';

export default{
    props:{messages:{type: Array<ChatMessage>,required:true}},
    emits:['send-chat-message'],
    setup(props,{emit}){
        const curMessage = ref(null);
        const sendMessage = ()=> {
            emit('send-chat-message',curMessage.value)
            curMessage.value = null;
        }
        return {
            sendMessage,
            curMessage
        }
    }
}
</script>


<style>
.chat{
    margin: 2%;
}
.chat span{
    font-size: 1.2em;
}

.messages{
    height: 300px; /* set the fixed height of the messages section */
    overflow-y: auto; /* enable vertical scrolling if the content overflows the height */
    margin: 30px;
    padding-right: 20px; 
}
.nomessages{
    text-align: center;
}
.card-action{
    display: grid;
    grid-template-columns: 6fr 1fr;
}
</style>