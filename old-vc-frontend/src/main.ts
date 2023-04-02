import Vue from 'vue'
import App from './App.vue'
import vuetify from './plugins/vuetify';
import router from './router';
import store from './store';
import {convertBoardStateToFEN} from './utils/fen'
import VueChatScroll from 'vue-chat-scroll';
import '@fortawesome/fontawesome-free/css/all.css'

Vue.config.productionTip = false;
Vue.use(VueChatScroll);

new Vue({
  render: h => h(App),
  store,
  router,
  vuetify
}).$mount('#app')