import Vue from 'vue'
import Vuex from 'vuex'
//import VuexPersistence from 'vuex-persist'
import state from './state.js'
import mutations from './mutations.js'
//import actions from './actions.js'
Vue.use(Vuex)

/*const vuexLocal = new VuexPersistence({
    storage: window.localStorage
})*/

const store = new Vuex.Store({
    state,
    mutations,
    //getters,
    //strict: debug,
    //plugins: [vuexLocal.plugin]
})

export default store