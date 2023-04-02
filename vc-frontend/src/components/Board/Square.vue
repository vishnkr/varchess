<template>
    <div class="square" :style="cssVars" @click="emitSquareClick">
        <div v-if="square.squareInfo.isPiecePresent">
            <Piece :squareInfo="square.squareInfo"/>
        </div>
    </div>
</template>

<script lang="ts">
import {ref, PropType, computed} from 'vue'
import Piece from './Piece.vue'
import {SquareColor, SquareInfo,EditorModeType, EditorState, Square} from '../../types'

export default {
    name: 'BoardSquare',
    components:{
        Piece
    },
    emits:['emit-square-click'],
    props:{
        square: {type: Object as PropType<Square>, required: true},
        editorState: {type: Object as PropType<EditorState | null> }
    },
    
    setup(props,{emit}){
        const colorMap : { [key in SquareColor]: string } = {
            'dark': '#b2c85d',
            'light': '#e4f5cb',
            'disabled' : '#696969'

        }
        const cssVars = computed(()=>{
            return {
                '--x': props.square.squareInfo.row,
                '--y': props.square.squareInfo.col,
                '--color': colorMap[(props.square.disabled ? 'disabled' : props.square.squareInfo.squareColor)]
            }
        })
        const emitSquareClick = ()=>{
            let squarePos = {row:props.square.squareInfo.row,col:props.square.squareInfo.col}
            let payload = {clickType:'perform-move',squareInfo:squarePos}
            if (props.editorState?.editorType){
                if(props.editorState.editorType==='Game'){
                    if (props.editorState.isDisableTileOn){
                        payload.clickType = 'disable'
                    } else { payload.clickType = 'toggle-piece'}
                    
                } else if (props.editorState.editorType==='MP'){
                    payload.clickType = 'set-jump-mp'
                }
            }
            emit('emit-square-click',payload)
        }
        return{
            cssVars,
            emitSquareClick
        }
    }
}
</script>

<style scoped>
.square {
  background: transparent;
  border: 1px solid transparent;
  width: 100%;
  height: 0;
  padding-bottom: 100%;
  grid-column: var(--y);
  grid-row: var(--x);
  cursor: pointer;
  background-color: var(--color);
}

.highlight-to{
  background: rgba(217,191,119, 0.6) !important;
  border: 0.2px solid;
}

.highlight-from{
  background-color: #a97d5d !important;
}

.move-jump-pattern{
  background-color: #4056b8 !important;
  border-color: black;
}
.dark {
    background-color: #b2c85d;
  }

.light {
    background-color: #e4f5cb;
}

.disabled{
  background-color: #696969 !important;
}
.move-slide-pattern{
  background-color: #ac422a !important;
  border-color: black;
}

@media only screen and (max-device-width: 480px) {
  .square {
  background: transparent;
  border: 1px solid transparent;
  float: left;
  font-size: 6px;
  font-weight: bold;
  line-height: 34px;
  height: 12px;/*48px;*/
  margin-right: -1px;
  margin-top: -1px;
  padding: 0;
  text-align: center;
  width: 12px;
  }
  .dark {
    background-color: #b2c85d;
  }

  .light {
    background-color: #e4f5cb;
  }
}
</style>