import { ViteSSG } from 'vite-ssg'
import './style.css'
import 'virtual:uno.css'
import App from './App.vue'

import OpLogin from '@/pages/op_login.vue'
import OpLogout from '@/pages/op_logout.vue'
import RpLogin from '@/pages/rp_login.vue'

export const createApp = ViteSSG(
  // the root component
  App,
  {
    routes: [
      {
        path: '/op_login',
        component: OpLogin,
      },
      {
        path: '/op_logout',
        component: OpLogout,
      },
      {
        path: '/rp_login',
        component: RpLogin,
      },
    ],
  },
)
