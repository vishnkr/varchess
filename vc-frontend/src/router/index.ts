import Vue from "vue";
import {createRouter, createWebHistory} from "vue-router";
import Editor from '../pages/Editor.vue';
import GameRoom from '../components/Room/GameRoom.vue';
import JoinRoom from '../components/Room/JoinRoom.vue';
import HomePage from '../pages/HomePage.vue';
import Login from '../pages/Login.vue';
import Signup from '../pages/Signup.vue';
import Dashboard from '../components/Dashboard.vue';

export type AppRoutes = 'Home' | 'Editor' | 'Game' | 'JoinRoom' | 'Login' | 'Signup' | 'Dashboard'
const routes = [
    { path: '/', 
      name: 'Home',
      component: HomePage
    },
    {
      path: '/editor/:username',
      name: 'Editor',
      props: true,
      component: Editor
    },
    {
      path: '/game/:username/:roomId',
      name: 'Game',
      component: GameRoom,
      props: true
    },
    {
      path: '/join/:roomId',
      component: JoinRoom
    },
    {
      path: '/login',
      name:'Login',
      component: Login
    },
    {
      path: '/signup',
      name:'Signup',
      component: Signup
    },
    {
      path: '/dashboard',
      name: 'Dashboard',
      component: Dashboard
    }

  ]
//const NotFoundComponent = { template: "<p style='text-align:'center';'> Oops! The page you're looking for doesn't exist </p>" }

const router =  createRouter({
    routes,
    history: createWebHistory(import.meta.env.BASE_URL)
})

export default router;