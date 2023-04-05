<template>
    <div class="q-pa-md">
        <div class="q-gutter-y-md" style="max-width: 600px">
            <q-card >
                
            <q-tabs v-model="tab" align="justify" class="bg-dark text-orange">
                <q-tab name="members" label="Members"> 
                    <font-awesome-icon :icon="['fas', 'user-group']" /> 
                </q-tab>
                <q-tab name="chat" label="Chat"> 
                    <font-awesome-icon :icon="['fas', 'message']" />
                </q-tab>
                <q-tab name="move-pattern" label="Move Patterns"> 
                    <font-awesome-icon icon="fa-solid fa-chess" size="lg"/> 
                </q-tab>
                
            </q-tabs>
            <q-separator />
            <q-tab-panels dark v-model="tab" animated>
                <q-tab-panel name="members">
                    <q-card dark>
                        <q-card-section>
                            <p class="text-h6 tab-title">Players </p>
                            <q-banner dense class="text-black bg-white" >Black</q-banner>
                            <q-separator/>
                            <q-banner dense class="text-white bg-black" >Black</q-banner>
                        </q-card-section>
                        <q-card-section>
                            <p class="text-h6 tab-title">Viewers </p>
                                <q-list bordered padding class="rounded-borders text-primary">
                                    <q-item
                                        clickable
                                        :active="true"
                                        active-class="menu-item"
                                    >
                                    <q-item-section>Inbox</q-item-section>
                                    </q-item>
                                </q-list>
                        </q-card-section>
                    </q-card>
                </q-tab-panel>
                <q-tab-panel name="chat">
                    <div class="text-h6 tabl-title">Chat</div>
                    Lorem ipsum dolor sit amet consectetur adipisicing elit.
                </q-tab-panel>
                <q-tab-panel name="move-pattern">
                    <div class="text-h6 tab-title" >Move Patterns</div>
                    Lorem ipsum dolor sit amet consectetur adipisicing elit.
                </q-tab-panel>

                </q-tab-panels>
            </q-card>

        </div>
    </div>

    
</template>

<script lang="ts">
import { ref } from 'vue';

export default {
    setup(){
        const tab = ref('members');
        return{
            tab,
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
</style>