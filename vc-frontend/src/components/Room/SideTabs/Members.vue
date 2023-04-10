<template>
    <q-card dark class="members-card">
        <q-card-section>
            <p class="text-h6 tab-title">Players </p>
            <q-banner dense class="text-black bg-white"> {{ members.p1 }} - White</q-banner>
            <q-separator/>
            <q-banner dense class="text-white bg-black">{{ members.p2 }} - Black</q-banner>
        </q-card-section>
        <q-card-section>
            <p class="text-h6 tab-title">Viewers </p>
            <q-list v-for="viewer in members.viewers" bordered padding class="rounded-borders text-primary">
                <q-item clickable v-ripple>
                    <q-item-section avatar>
                    <q-avatar rounded color="primary" text-color="white">
                        {{ viewer.toUpperCase()[0] }}
                    </q-avatar>
                    </q-item-section>
                    <q-item-section>{{ viewer }}</q-item-section>
                    
                </q-item>
            </q-list>
        </q-card-section>
    </q-card>

</template>

<script lang="ts">
import { RootState } from '@/store/state';
import { SET_PLAYERS, UPDATE_MEMBERS } from '../../../store/mutation_types';
import { reactive } from 'vue';
import { useStore } from 'vuex';

export default{
    setup(){
        const store = useStore<RootState>();
        const members = reactive({
            p1:store.state.gameInfo?.players.p1, 
            p2:store.state.gameInfo?.players.p2 ?? '',
            viewers: store.state.gameInfo?.members ?? []
        });
        store.subscribe((mutation,state)=>{
            if (mutation.type==UPDATE_MEMBERS && state.gameInfo?.members){
                members.viewers = state.gameInfo?.members
            } else if (mutation.type==SET_PLAYERS && state.gameInfo?.players.p2){
                members.p1 = state.gameInfo?.players.p1;
                members.p2 = state.gameInfo?.players.p2;
            }
        })
        return{
            members
        }
    }
}


</script>

<style>
.members-card{
    height: 300px; 
    overflow-y: auto;
    margin: 30px;
    padding-right: 20px;
}
</style>