import Vue from 'vue';
import Router, {RouteConfig} from 'vue-router';

Vue.use(Router);

const routeConfig: RouteConfig[] = [
    {
        path:'/',
        alias: '/home',
        component: () => import(`@/views/home/home.vue`),
    }
];
export {routeConfig};
export default new Router({
    mode: 'hash',
    base: process.env.BASE_URL,
    routes: routeConfig,
});
