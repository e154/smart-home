import { RouteConfig } from 'vue-router'
import Develop from '@/layout/develop.vue'

const logsRouter: RouteConfig = {
  path: '/logs',
  component: Develop,
  redirect: '/logs',
  name: 'logs',
  meta: { hidden: false },
  children: [
    {
      path: '',
      component: () => import('@/views/log/index.vue'),
      name: 'Logs',
      meta: {
        title: 'logs',
        icon: 'log',
        noCache: true
      }
    }
  ]
}

export default logsRouter
