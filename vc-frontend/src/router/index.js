import Vue from "vue";
import Router from "vue-router";
import HomePage from '../views/HomePage'
import EditorDialog from '../components/Editor/EditorDialog'

Vue.use(Router);
const routes = [
    { path: '/', 
      name: 'Home',
     component: HomePage 
    },
    {
      path: '/editor',
      name: 'Editor',
      component: EditorDialog,
    }
  ]
//const NotFoundComponent = { template: '<p>Page not found</p>' }

const router = new Router({
    routes
})
  
export default router