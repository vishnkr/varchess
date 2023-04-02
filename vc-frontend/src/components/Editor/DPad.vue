<template>
    <div class="d-pad">
        <q-btn class="dpad-btn" color="orange" @click="shiftBoard('up')"><font-awesome-icon :icon="faArrowUp" /></q-btn>
        <div class="mid-row">
            <q-btn class="dpad-btn" color="orange" @click="shiftBoard('left')"><font-awesome-icon :icon="faArrowLeft" /></q-btn>
            <q-btn class="dpad-btn" color="orange" @click="shiftBoard('right')"><font-awesome-icon :icon="faArrowRight" /></q-btn>
        </div>
      <q-btn class="dpad-btn" color="orange" @click="shiftBoard('down')"><font-awesome-icon :icon="faArrowDown" /></q-btn>
    </div>
  </template>
  
  <script lang="ts">
  import { defineComponent } from 'vue'
  import { faArrowUp,faArrowDown,faArrowRight, faArrowLeft } from '@fortawesome/free-solid-svg-icons'
  export default defineComponent({
    name: 'DPad',
    emit:['shift-board-direction'],
    setup(props,{emit}){
        return{
            shiftBoard: (direction:String)=>{
                emit('shift-board-direction',direction)
            },
            faArrowUp,
            faArrowLeft,
            faArrowRight,
            faArrowDown
        }
    }
  })
  </script>
  
  <style scoped>
  .d-pad {
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    align-items: center;
    height: 120px;
  }
  
    .dpad>q-btn{
        flex:1;
    }
  .dpad-btn{
    margin: 5px;
    width:70px;
    max-height: 60px;
  }
  
  .mid-row{
    flex:1;
  }
  </style>