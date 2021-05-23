import Vue from "vue";
import Router from "vue-router";
import HomePage from '../views/HomePage'

Vue.use(Router);
const routes = [
    { path: '/', 
      name: 'Home',
     component: HomePage 
    },
  ]
//const NotFoundComponent = { template: '<p>Page not found</p>' }

const router = new Router({
    routes
})
  
export default router