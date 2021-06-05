<template>
  <div class="container">
    <div class="columns">
        <div class="column board-panel">
          <game-board />
        </div>
        <div class="column right-panel">
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