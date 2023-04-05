<template>
    <q-page class="content-panel bg-dark">
            <div class="board-wrapper">
                 <board :board-state="boardState" :is-flipped="isFlipped" />
            </div>
            <div class="column right-panel">
                <side-tabs />
            </div>
        
    </q-page>

</template>


<script lang="ts">
import { RootState } from '../../store/state';
import { reactive, ref } from 'vue';
import { useStore } from 'vuex';
import Board from '../Board/Board.vue';
import SideTabs from '../Room/SideTabs.vue';

export default{
    components:{Board, SideTabs},
    setup(){
        const store = useStore<RootState>();
        const boardState = reactive(store.state.board);
        const isFlipped = ref(store.state.userInfo.curGameRole === 'p2');
        return{
            boardState,
            isFlipped
        }
    }
}
</script>

<style scoped>

.banner{
    padding: 0 !important;
}
.content-panel{
    display: flex;
    flex-direction: row;
    justify-content: center;
}
.content-panel>.board-wrapper{
    display: flex;
    align-items: center;
    flex:1;
    padding:1%;
}
.board-wrapper{
    flex:4;
    max-width: 900px;
    width:100%;
    
}
.left-panel{
    display:flex;
    flex:1;
    flex-direction: column;
}
.left-panel > .banner{
    flex:1;
}
.right-panel{
    display: flex;
    flex:1;
    max-width:50%;
}
</style>