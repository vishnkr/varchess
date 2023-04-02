<template>
    <div class="piece" :style="cssVars" >
       <img v-if="props.squareInfo.isPiecePresent" :src="getPieceURL">
    </div>
  
</template>

<script setup lang="ts">
import {computed, PropType, StyleValue} from 'vue'

import { SquareInfo } from '@/types';

const props = defineProps({
    squareInfo: {type: Object as PropType<SquareInfo>, required: true}
})

const getPieceURL = computed(()=>{
    const path = new URL(`../../assets/images/pieces/${props.squareInfo?.pieceColor}/${props.squareInfo?.pieceType}.svg`, import.meta.url).href
    return path
})

const cssVars = computed(()=>{
    return {
            '--x': props.squareInfo.row,
            '--y': props.squareInfo.col,
        } as any as StyleValue
})

</script>

<style scoped>

.piee {
        position: relative;
        text-align: center;
        color: white;
        font-size: 1.2em;
        -webkit-text-stroke-width: 1px;
        -webkit-text-stroke-color: black;
        width: 100%;
        height: 0;
        padding-bottom: 100%;
        pointer-events: none;
        grid-column: var(--y);
        grid-row: var(--x);
    }
    img {
        -khtml-user-select: none;
        -o-user-select: none;
        -moz-user-select: none;
        -webkit-user-select: none;
        user-select: none;
    }
</style>