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
        path: '/login',
        component: import('@/pages/login.vue'),
      },
    ],
  }, // { routes },
  // ({ app, router, routes, isClient, initialState }) => {
  //   // install plugins etc.
  // },
)
