import Vue from 'vue'
import Vuex from 'vuex'
//import VuexPersistence from 'vuex-persist'
import state from './state'
import mutations from './mutations'
import actions from './actions'
Vue.use(Vuex)

/*const vuexLocal = new VuexPersistence({
    storage: window.localStorage
})*/

const store = new Vuex.Store({
    state,
    mutations,
    actions,
    //getters,
    //strict: debug,
    //plugins: [vuexLocal.plugin]
})

export default store