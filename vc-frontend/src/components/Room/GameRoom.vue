<template>
  <div class="container">
    <div class="columns">
        <div class="column board-panel">
          <game-board />
        </div>
        <div class="column right-panel">
          <v-text-field v-model="shareLink" id="tocopy" readonly outlined ></v-text-field>
          <v-btn @click="copyText">Copy Link</v-btn>
          <chat :username="username" :roomId="roomId"/>
        </div>
    </div>
  </div>
</template>

<script>
import GameBoard from '../GameBoard.vue'
import Chat from './Chat.vue'

export default {
  components: { Chat,GameBoard },
 
  methods:{
    
    getShareUrl(){
      
      return `${window.location.origin}/#/join/${this.$route.params.roomId}`
    },
    copyText(){
    let input=document.getElementById("tocopy");
   input.select();
         document.execCommand("copy");
    },
    
  },
  data(){
    return{
      shareLink: this.getShareUrl(),
      username: this.$route.params.username,
      roomId: this.$route.params.roomId,
    }
  }
}
</script>

<style>
.columns{
  display: flex;

}

.board-panel{
  flex: 2;
}
.right-panel{

   max-width:400px;
   width: 100%;
   padding: 1em;

}

</style>