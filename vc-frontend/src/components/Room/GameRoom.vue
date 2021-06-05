<template>
  <div class="container">
    <div class="columns">
        <div class="column">
          <game-board />
        </div>
        <div class="column">
          <chat />
        </div>
    </div>
  </div>
</template>

<script>
import GameBoard from '../GameBoard.vue'
import Chat from './Chat.vue'
export default {
  components: { Chat,GameBoard },

}
</script>

<style>

</style>