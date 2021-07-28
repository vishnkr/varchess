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
      props: true,
      component: EditorDialog,
    },
    {
      path: '/game/:username/:roomId',
      name: 'Game',
      component: GameRoom,
    },
    {
      path: '/join/:roomId',
      component: JoinRoom,
    },

  ]
//const NotFoundComponent = { template: "<p style='text-align:'center';'> Oops! The page you're looking for doesn't exist </p>" }

const router = new Router({
    routes,
    mode: 'history'
})
  
export default router