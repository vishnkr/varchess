import Vue from 'vue'
import Vuex from 'vuex'
import VuexPersistence from 'vuex-persist'

Vue.use(Vuex)

const vuexLocal = new VuexPersistence({
    storage: window.localStorage
})

const store = new Vuex.Store({
    //state,
    //mutations,
    //actions,
    //getters,
    //strict: debug,
    plugins: [vuexLocal.plugin]
})

export default store