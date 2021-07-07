<template>
    <div :style="cssVars" >
       <img v-if="pieceType" :src="getPieceURL">
    </div>
  
</template>

<script>
export default {
    data(){
        return {
            
        }
    },
    props:{
        pieceType: String,
        color: String,
        row: Number,
        col: Number,
    },
    computed:{
        getPieceURL(){ return require(`../assets/images/pieces/${this.color}/${this.pieceType}.svg`)},
        cssVars(){
            return {
            '--x': this.row,
            '--y': this.col,
        }
        }
    }
}
</script>

<style scoped>
h4{
    color: black;
}
div {
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