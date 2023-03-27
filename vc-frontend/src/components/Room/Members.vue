<template>
  <div class="members" v-chat-scroll="{always:false, smooth:true}">
      <v-card
        class="mx-auto"
        tile
      >
        <div class="member-container">
        <v-list>
        <v-subheader>Players</v-subheader>
        <v-list-item-group
            v-model="selectedItem"
            color="primary"
        >
            <v-list-item
            v-for="(player, i) in getPlayersList"
            :key="i"
            >
            <v-list-item-content>
                {{player}} 
            </v-list-item-content>
            <v-list-item-icon v-if="username==player">
                <v-icon color="pink">
                    mdi-star
                </v-icon>
                (You)
             </v-list-item-icon>
            </v-list-item>
        </v-list-item-group>
        <v-subheader>Viewers</v-subheader>
        <v-list-item-group
            v-model="selectedItem"
            color="primary"
        >
            <v-list-item
            v-for="(member, i) in members"
            :key="i"
            >
            <v-list-item-content>
                {{member}}
            </v-list-item-content>
            <v-list-item-icon v-if="username==member">
                <v-icon color="pink">
                    mdi-star
                </v-icon>
                (You)
            </v-list-item-icon>
            </v-list-item>
        </v-list-item-group>

        </v-list>
        </div>
    </v-card>
  </div>
</template>

<script lang="ts">
import { stringify } from 'querystring';
import Vue from 'vue';
import { ComputedOptions } from 'vue/types/options';
import { UPDATE_MEMBERS } from '../../utils/mutation_types';

export default Vue.extend({
    props:['username'],
    mounted: function(){
        this.p1 = this.$store.state.gameInfo?.players.p1
        this.p2 = this.$store.state.gameInfo?.players.p2
        this.members = this.$store.state.gameInfo?.members
        this.$store.subscribe((mutation, state) => {
        if(mutation.type=== UPDATE_MEMBERS){
            this.members = this.$store.state.gameInfo?.members;
        }
     })
    },
    computed:{ 
            getPlayersList(): string[]{
                return (this.p1 && this.p2) ? [this.p1,this.p2] : [];
            }
    },
    data: () => ({
      selectedItem: null,
      members: [],
      p1: null,
      p2: null,
    }),
})
</script>

<style scoped>
.member-container{
    height: 300px;
    overflow: auto;
    margin: 30px
}
</style>