<template>
    <div>
        <div>
            <v-pre class="text-white text-h6 bg-dark">Board height (rows)</v-pre>
            <q-slider 
                v-model="dimensions.rows" 
                color="orange" 
                marker-labels 
                :min="5" 
                :max="isLoggedIn ? 16 : 8" 
                @change="emitBoardDimensions"
                />
            <p class="text-white text-h6 bg-dark">Board width (cols)</p>
            <q-slider 
                v-model="dimensions.cols" 
                color="orange" 
                marker-labels 
                :min="5" 
                :max="isLoggedIn ? 16 : 8" 
                @change="emitBoardDimensions"
                />
            <q-toggle
                size="lg"
                :model-value="isDisableTileOn"
                label="Disable Tile"
                @update:model-value="emitDisable"
            />
            
            <p class="text-h6 bg-dark text-white"> Shift Board</p>
            <div class="dpad-wrapper">
                <DPad @shift-board-direction="emitShiftBoard"/>
            </div>
            
        </div>
    </div>
</template>

<script lang="ts">
import { RootState } from '@/store/state';
import { Dimensions } from '@/types';
import {ref,computed,defineComponent, reactive} from 'vue'
import { useStore } from 'vuex';
import DPad from './DPad.vue';

export default defineComponent({
    components:{DPad},
    props:{
        isLoggedIn: Boolean,
        isDisableTileOn: Boolean,
        rows:{type:Number,required:true},
        cols:{type:Number,required:true},
    },
    emits:['update-board-dimensions','toggle-disable','shift-board-direction'],
    setup(props,{emit}){
        const dimensions: Dimensions = reactive({rows:props.rows, cols:props.cols})
        const isDisableTileOn = ref(props.isDisableTileOn);
        function emitBoardDimensions(){
            emit('update-board-dimensions',dimensions)
        }
        function emitDisable(){
            isDisableTileOn.value = !isDisableTileOn.value;
            emit('toggle-disable',isDisableTileOn.value)
        }
        return {
            dimensions,
            isDisableTileOn,
            emitBoardDimensions,
            emitDisable,
            emitShiftBoard: (direction:string)=>{emit('shift-board-direction',direction)}
        }
    }
});
</script>