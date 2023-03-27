import Vuex from 'vuex'
import state, { RootState } from './state'
import mutations from './mutations'
import actions from './actions'
import webSocketModule from './modules/webSocket'
import { Store } from 'vuex'
import Vue from 'vue';

Vue.use(Vuex)
const store :Store<RootState>= new Vuex.Store({
    state,
    mutations,
    actions,
    modules:{
        webSocket:webSocketModule
    },
})

export default store