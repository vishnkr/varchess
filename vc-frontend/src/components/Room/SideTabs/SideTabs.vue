<template>
    <div class="q-pa-md">
        <div class="q-gutter-y-md" style="max-width: 600px">
            <q-card >
            <div class="btns bg-dark">
                <q-btn @click="emitAction('flip')" class="text-black" color="white"> Flip <font-awesome-icon :icon="['fas', 'repeat']" style="padding-left: 5px; font-size: 1.2em" />  </q-btn>
                <q-btn @click="emitAction('draw')" class="text-white bg-blue-9">  Draw <font-awesome-icon :icon="['fas', 'handshake']" style="padding-left: 5px; font-size: 1.2em"/></q-btn>
                <q-btn @click="emitAction('resign')" color="negative">  Resign <font-awesome-icon :icon="['fas', 'flag']" style="padding-left: 5px; font-size: 1.2em"/></q-btn>
            </div>

            <q-tabs v-model="tab" align="justify" class="bg-dark text-orange">
                <q-tab name="members" label="Members"> 
                    <font-awesome-icon :icon="['fas', 'user-group']" /> 
                </q-tab>
                <q-tab name="chat" @click="newMessage=false" label="Chat"> 
                    <font-awesome-icon :icon="['fas', 'message']" />
                    <q-badge v-if="newMessage" color="red" floating>
                    </q-badge>
                </q-tab>
                <q-tab name="move-pattern" label="Move Patterns"> 
                    <font-awesome-icon icon="fa-solid fa-chess" size="lg"/> 
                </q-tab>
                
            </q-tabs>
            <q-separator />
            <q-tab-panels dark v-model="tab" animated>
                <q-tab-panel name="members">
                    <div class="text-h6 tab-title" >Members</div>
                    <members />
                </q-tab-panel>
                <q-tab-panel name="chat">
                    <div class="text-h6 tab-title">Chat</div>
                    <chat :messages="messages" @send-chat-message="sendChatMessage"/>
                </q-tab-panel>
                <q-tab-panel name="move-pattern" v-if="movePatterns.movePattern">
                    <div class="text-h6 tab-title" >Move Patterns</div>
                    <move-pattern-tab :move-patterns="movePatterns.movePattern" :piece="Object.keys(movePatterns)[0]" />
                </q-tab-panel>

                </q-tab-panels>
            </q-card>

        </div>
    </div>

    
</template>

<script lang="ts">
import { reactive, ref } from 'vue';
import Chat from './Chat.vue';
import MovePatternTab from './MovePatternTab.vue';
import Members from './Members.vue';
import { RootState } from '@/store/state';
import { useStore } from 'vuex';
import { ADD_CHAT_MESSAGE } from '../../../store/mutation_types';
import {SEND_MESSAGE} from '../../../store/action_types';

export default {
    components:{Chat,MovePatternTab,Members},
    emits:['sidetab-action'],
    props:{roomId:{type:String,required:true}},
    setup(props,{emit}){
        const tab = ref('members');
        const store = useStore<RootState>();
        const newMessage = ref(false);
        const messages = ref(store.state.chatMessages[props.roomId])
        const movePatterns = reactive({movePattern:store.state.movePatterns})

        store.subscribe((mutation,state)=>{
            if (mutation.type == ADD_CHAT_MESSAGE){
                messages.value = state.chatMessages[props.roomId];
                if(tab.value!='chat'){
                    newMessage.value = true
                }
            }
        })
        const sendChatMessage = (message:string)=>{
            store.dispatch(SEND_MESSAGE,{roomId:props.roomId,message,username:store.state.userInfo?.username});
        }
        return{
            tab,
            emitAction: (type:string) => emit('sidetab-action',type),
            newMessage,
            messages,
            sendChatMessage,
            movePatterns
        }
    }
}
</script>

<style>
.tab-title{
    text-align: center;
}
.menu-item{
  color: white;
  background: #F2C037;
}
.btns{
    display: grid;
    grid-auto-flow: column;
}
.btns>*{
    margin: 5px;
    width:70px;
    max-height: 60px;
}
</style>