<template>
  <div class="square" :class="isDark? 'dark':'light'">
      
  </div>
</template>

<script>
export default {
    methods:{

    },
    props:{
        isDark: Boolean,
        row: Number,
        col: Number,
        /*isPiecePresent: Boolean,
        piece: String,
        pieceColor: String,*/
    }
}
</script>

<style scoped>


.square {
  background: transparent;
  border: 1px solid transparent;
  float: left;
  font-size: 24px;
  font-weight: bold;
  line-height: 34px;
  height: 48px;/*48px;*/
  margin-right: -1px;
  margin-top: -1px;
  padding: 0;
  text-align: center;
  width: 48px;
}

.dark {
  background-color: #b2c85d;
}

.light {
  background-color: #e4f5cb;
}

</style>