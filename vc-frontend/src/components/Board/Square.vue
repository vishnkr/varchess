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
import {SquareColor, SquareInfo,EditorModeType, EditorState, Square, MPEditorState, isMPEditor, SquareClick} from '../../types'
import { RootState } from '@/store/state'
import { useStore } from 'vuex'

export default {
    name: 'BoardSquare',
    components:{
        Piece
    },
    emits:['emit-square-click'],
    props:{
        square: {type: Object as PropType<Square>, required: true},
        editorState: {type: Object as PropType<EditorState | MPEditorState | null> }
    },
    
    setup(props,{emit}){
        const colorMap : { [key in SquareColor]: string } = {
            'dark': '#b2c85d',
            'light': '#e4f5cb',
            'disabled' : '#696969',
            'jump': '#4056b8',
            'slide': '#ac422a',
            'to': '#d9bf7799',
            'from': '#a97d5d'
        }
        const getColor = ()=>{
          let color : SquareColor = props.square.squareInfo.tempSquareColor ?? props.square.squareInfo.squareColor;
          if (props.editorState){
            if(props.editorState.editorType === "Game" && props.square.disabled){ color = 'disabled'}
          } 
          return colorMap[color];
        }

        const store = useStore<RootState>();
        const cssVars = computed(()=>{
            return {
                '--x': props.square.x,
                '--y': props.square.y,
                '--color': getColor()
            }
        })
        const emitSquareClick = ()=>{
            let squarePos = {row:props.square.squareInfo.row-1,col:props.square.squareInfo.col-1}
            let payload:SquareClick = {clickType:'select-mv-square',...squarePos}
            if (props.editorState?.editorType){
              if(props.editorState.editorType==='Game'){
                  if (props.editorState.isDisableTileOn){
                      payload.clickType = 'disable'
                  } else { payload.clickType = 'toggle-piece'}
                    
              } else if (isMPEditor(props.editorState)){
                  payload.clickType = props.editorState.moveType==='jump' ? 
                  props.square.squareInfo.tempSquareColor==='jump' ? 'remove-jump-mp' : 'set-jump-mp' 
                  : 'set-slide-mp'
              }
            } else{
              if(props.square.squareInfo.isPiecePresent){
                payload.piece = props.square.squareInfo.pieceType
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
</style>