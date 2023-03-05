import Vue from 'vue'
import Vuex from 'vuex'
//import VuexPersistence from 'vuex-persist'
import state from './state'
import mutations from './mutations'
import actions from './actions'
import webSocketModule from './modules/webSocket'
Vue.use(Vuex)

/*const vuexLocal = new VuexPersistence({
    storage: window.localStorage
})*/

const store = new Vuex.Store({
    state,
    mutations,
    actions,
    modules:{
        webSocket:webSocketModule
    },
})

export default store