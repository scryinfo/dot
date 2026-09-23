import { ViteSSG } from 'vite-ssg'
import './style.css'
import 'virtual:uno.css'
import App from './App.vue'

export const createApp = ViteSSG(
  // the root component
  App,
  {
    routes: [
      {
        path: '/op_login',
        component: () => import('@/pages/op_login.vue'),
      },
      {
        path: '/op_logout',
        component: () => import('@/pages/op_logout.vue'),
      },
      {
        path: '/rp_login',
        component: () => import('@/pages/rp_login.vue'),
      },
    ],
  }, // { routes },
  // ({ app, router, routes, isClient, initialState }) => {
  //   // install plugins etc.
  // },
)
