import { createRouter, createWebHistory } from 'vue-router'

import OpLogin from '@/pages/op_login.vue'
import OpLogout from '@/pages/op_logout.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/op_login',
      name: 'op_login',
      component: OpLogin,
    },
    {
      path: '/op_logout',
      name: 'op_logout',
      component: OpLogout,
    },
  ],
})

export default router
