import { createApp } from 'vue'
import App from './App.vue'
import { Quasar } from 'quasar'
// Import icon libraries
import '@quasar/extras/fontawesome-v5/fontawesome-v5.css'
import router from './router'
// Import Quasar css
import 'quasar/src/css/index.sass'
import store from './store'
import { library } from '@fortawesome/fontawesome-svg-core'

import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import { faGithub } from '@fortawesome/free-brands-svg-icons'
import {faHeart, faLink, faArrowUp, faArrowDown, faArrowLeft, faArrowRight} from '@fortawesome/free-solid-svg-icons'

/* add icons to the library */
library.add(faGithub,faHeart,faLink, faArrowUp, faArrowDown, faArrowLeft, faArrowRight)
const myApp = createApp(App)

myApp.use(router)
myApp.use(store)
myApp.component('font-awesome-icon', FontAwesomeIcon)
myApp.use(Quasar, {
  plugins: {}, 
})

myApp.mount('#app')