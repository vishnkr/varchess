<template>
  <div class="members" v-chat-scroll="{always:false, smooth:true}">
      <v-card
        class="mx-auto"
        tile
        elevation="2"
      >
        <div class="member-container">
        <v-list>
        <v-subheader>Players</v-subheader>
        <v-list-item-group
            v-model="selectedItem"
            color="primary"
        >
            <v-list-item
            v-for="(player, i) in players"
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
            <v-list-item-icon v-if="username==member && !players.includes(username)">
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

<script>
export default {
    props:['members','username','players'],
    data: () => ({
      selectedItem: null,
    }),
}
</script>

<style scoped>
.member-container{
    height: 300px;
    overflow: auto;
    margin: 30px
}
</style>