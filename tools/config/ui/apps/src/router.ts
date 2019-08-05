import Vue from 'vue';
import Router from 'vue-router';
import Home from './views/Home.vue';
import FindDot from './views/FindDot.vue';
import DotList from './views/DotList.vue';
import Import from './views/Import.vue';
import Export from './views/Export.vue';

Vue.use(Router);

export default new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    {
      path:"/findDot",
      name:"findDot",
      component:FindDot
    },
    {
      path:"/dotList",
      name:"dotList",
      component:DotList
    },
    {
      path:"/import",
      name:"import",
      component:Import
    },
    {
      path:"/export",
      name:"export",
      component:Export
    },
    {
      path: '/Config',
      name: 'config',
      component: () => import('./views/Config.vue')
    }
  ],
});
