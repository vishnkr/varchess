import Vue from "vue";
import Router from "vue-router";
import HomePage from '../views/HomePage'
import Login from '../views/Login.vue'
import Signup from '../views/Signup.vue'
import EditorDialog from '../components/Editor/EditorDialog'
import GameRoom from '../components/Room/GameRoom'
import JoinRoom from '../components/Room/JoinRoom'
import Dashboard from "../components/Dashboard.vue";

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
    {
      path: '/login',
      name:'Login',
      component: Login,
    },
    {
      path: '/signup',
      name:'Signup',
      component: Signup,
    },
    {
      path: '/dashboard',
      name: 'Dashboard',
      component: Dashboard,
    }

  ]
//const NotFoundComponent = { template: "<p style='text-align:'center';'> Oops! The page you're looking for doesn't exist </p>" }

const router = new Router({
    routes,
    mode: 'history'
})
  
export default router