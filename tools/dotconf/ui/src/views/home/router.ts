import Vue from 'vue'
import VueRouter, {RouteConfig} from 'vue-router'
import Home from '@/views/home/home.vue'
import FindDot from '@/views/page/findDot.vue'
import DotList from '@/views/page/dotList.vue'
import Import from '@/views/page/import.vue'
import Export from '@/views/page/export.vue'
import Config from '@/views/page/config.vue'


Vue.use(VueRouter);

const routes: Array<RouteConfig> = [
    {
        path: '/', component: Home,
        children: [
            {
                path: '/home',
                name: 'home',
                component: Home
            },
            {
                path: '/findDot',
                name: 'findDot',
                component: FindDot,
            },
            {
                path: '/dotList',
                name: 'dotList',
                component: DotList,
            },
            {
                path: '/import',
                name: 'import',
                component: Import,
            },
            {
                path: '/export',
                name: 'export',
                component: Export,
            },
            {
                path: '/Config',
                name: 'config',
                component: Config,
            },
        ]
    }
]

const router = new VueRouter({
    mode: "hash",
    base: process.env.BASE_URL,
    routes
});

export default router
