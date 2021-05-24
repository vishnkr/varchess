<template>
  <div class="board-editor">
    <div class="side-panel">
      <v-card class="mx-auto" max-width="500">
        <v-list-item two-line>
        <v-list-item-content>
          <v-list-item-title class="headline">
            Game Editor
          </v-list-item-title>
          <div>
          <v-list-item-title>Board width {{cols}}</v-list-item-title>
          <v-slider
            v-model="rows"
            :max="14"
            :tick-labels="labels"
            
            ticks
          ></v-slider>
          </div>
          <div>
          <v-list-item-title>Board height {{rows}}</v-list-item-title>
          <v-slider
            v-model="cols"
            :max="14"
            :tick-labels="labels"
            
          ></v-slider>
          </div>
        </v-list-item-content>
        </v-list-item>
        <v-card-text>

        </v-card-text>
      </v-card>
    </div>
    <board class="board-panel" :cols="labels[cols]" :rows="labels[rows]" />
  </div>
</template>

<script>
import Board from './Board';
export default {
  components:{Board},
  data(){
    return{
      labels: [1,2,3,4,5,6,7,8,9,10,11,12,13,14,15],
      rows: 7,
      cols: 7,
    }
  }

}
</script>

<style scoped>
.board-editor{
  display: flex;
}
.side-panel{
  flex:2;
  height: 90vh;
  /*box-shadow: 0 14px 28px rgba(250, 0, 0, 0.25), 0 10px 10px rgba(0,0,0,0.22);*/
}
.board-panel{
  flex:3;
}
</style>