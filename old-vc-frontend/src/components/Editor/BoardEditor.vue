<template>
   <div>
       <div>
            <v-list-item-title style="padding-top:10px">Board height (rows)</v-list-item-title>
                <v-slider
                v-model="labelRow"
                :max="labels.length-1"
                @input="emitBoardDimensions"
                :tick-labels="labels"
                color="orange darken-3"
                ></v-slider>
        </div>
        <div>
            <v-list-item-title>Board width (cols)</v-list-item-title>
            <v-slider
                v-model="labelCol"
                :max="labels.length-1"
                @input="emitBoardDimensions"
                :tick-labels="labels"
                color="orange darken-3"
            ></v-slider>
        </div>
        <v-checkbox
            v-model="editorState.isDisableTileOn"
            label="Toggle: Disable selected tile"
        ></v-checkbox>

    </div>
</template>

<script>

export default {
    props:['editorState'],
    methods:{
        emitBoardDimensions(){
            this.$emit('update-board-dimensions',{
                rows: this.labels[this.labelRow],
                cols: this.labels[this.labelCol],
            });
        }
    },
    data(){
      return {
        labels: [5,6,7,8,9,10,11,12,13,14,15,16],
        labelRow: 3,
        labelCol: 3,
      }
    }
}
</script>

<style>

</style>