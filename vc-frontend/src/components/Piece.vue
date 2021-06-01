<template>
    <div>
        <img class="resize" :src="getPieceURL">
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
    },
    computed:{
        getPieceURL(){ return require(`../assets/images/pieces/${this.color}/${this.pieceType}.svg`)}
    }
}
</script>

<style scoped>
.resize{
    height: auto;
    width: auto;
}
</style>