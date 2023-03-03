import Vue from "vue";
import Router, { Route } from "vue-router";

Vue.use(Router);
const routes = [
    { path: '/', 
      name: 'Home',
      component: () => import('@/views/HomePage.vue')
    },
    {
      path: '/editor/:username/:roomId',
      name: 'Editor',
      props: true,
      component: () => import('@/components/Editor/EditorDialog.vue'),
    },
    {
      path: '/game/:username/:roomId',
      name: 'Game',
      component: () => import('@/components/Room/GameRoom.vue'),
      props: (route:Route) => ({
        username: route.params.username,
        roomId: route.params.roomId,
        boardState: route.params.boardState
      })
    },
    {
      path: '/join/:roomId',
      component: () => import('@/components/Room/JoinRoom.vue'),
    },
    {
      path: '/login',
      name:'Login',
      component: () => import('@/views/Login.vue'),
    },
    {
      path: '/signup',
      name:'Signup',
      component: () => import('@/views/Signup.vue'),
    },
    {
      path: '/dashboard',
      name: 'Dashboard',
      component: () => import('@/components/Dashboard.vue'),
    }

  ]
//const NotFoundComponent = { template: "<p style='text-align:'center';'> Oops! The page you're looking for doesn't exist </p>" }

const router =  new Router({
    routes,
    mode: 'history'
})

export default router;