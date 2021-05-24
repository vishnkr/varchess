<template>
  <div id="board-container">
  <div id="board" :style="cssVar">
      <div class="board-row" v-for="col in cols" :key="col">
        <div v-for="row in rows" :key="row">
          
        <square :isDark="!isLight(row,col)" :row="row" :col="col" />
        </div>
      </div>
        
  </div>
  </div>
</template>

<script>

import Square from './Square'
export default {
    name:'Board',
    methods:{
      isEven(val){return val%2==0},
      isLight(row,col){
        return this.isEven(row)&&this.isEven(col)|| (!this.isEven(row)&&!this.isEven(col))},
      
    },
    computed:{
      cssVar(){
        return {
        '--size': Math.max(this.boardState.cols, this.boardState.rows),
        }
      }
    },
    data(){
      return {
        rowMultiplier:0,
        
        boardState: {
          rows: 15,
          cols: 15,
          
          tiles: [{x:0, y:0, color:'dark'},
                  {x:0, y:1, color:'light'},
                  {x:0, y:2, color:'dark'},
                  {x:0, y:3, color:'light'},
                  {x:0, y:4, color:'dark'},
                  {x:1, y:0, color:'light'},
                  {x:1, y:1, color:'dark'},
                  {x:1, y:2, color:'light'},
                  {x:1, y:3, color:'dark'},
                  {x:1, y:4, color:'light'},
                  {x:2, y:0, color:'light'},
                  {x:2, y:1, color:'dark'},
                  {x:2, y:2, color:'light'},
                  {x:2, y:3, color:'dark'},
                  {x:2, y:4, color:'light'},
                  {x:3, y:0, color:'dark'},
                  {x:3, y:1, color:'light'},
                  {x:3, y:2, color:'dark'},
                  {x:3, y:3, color:'light'},
                  {x:3, y:4, color:'light'},
                  {x:4, y:0, color:'dark'},
                  {x:4, y:1, color:'light'},
                  {x:4, y:2, color:'dark'},
                  {x:4, y:3, color:'light'},
                  {x:4, y:4, color:'dark'},    
          ]
        }
      }
    },
    props:{
      rows: Number,
      cols: Number,
    },
    components:{
      Square,
    },
    mounted(){

    },
    created(){
        
    }
}
</script>

<style scoped>


#board-container{
  width: 50%;
  height: 100vh;
}

#board{
    height: 90vh;
    display: grid;
    grid-template-columns: repeat(var(--size), 0fr);
    /*grid-template-rows: repeat(var(--size), 1fr);*/
    box-shadow: 0 14px 28px rgba(0,0,0,0.25), 0 10px 10px rgba(0,0,0,0.22);
    align-items: center;
    justify-content: center;
    
}



</style>