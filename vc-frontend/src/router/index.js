import Vue from "vue";
import Router from "vue-router";
import HomePage from '../views/HomePage'
import EditorDialog from '../components/Editor/EditorDialog'
import GameRoom from '../components/Room/GameRoom'
import JoinRoom from '../components/Room/JoinRoom'
Vue.use(Router);
const routes = [
    { path: '/', 
      name: 'Home',
     component: HomePage 
    },
    {
      path: '/editor/:username/:roomId',
      name: 'Editor',
      component: EditorDialog,
    },
    {
      path: '/game/:username/:roomId',
      component: GameRoom,
    },
    {
      path: '/join/:roomId',
      component: JoinRoom,
    }

  ]
//const NotFoundComponent = { template: '<p>Page not found</p>' }

const router = new Router({
    routes
})
  
export default router