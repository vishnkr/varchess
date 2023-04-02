import Vuex, { createStore } from 'vuex'
import state, { RootState } from './state'
import mutations from './mutations'
import actions from './actions'
import webSocketModule from './modules/webSocket'

const store = createStore({
    state,
    mutations,
    actions,
    modules:{
        webSocket:webSocketModule
    },
})

export default store