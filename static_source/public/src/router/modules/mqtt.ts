import { RouteConfig } from 'vue-router'
import Develop from '@/layout/develop.vue'

const mqttRouter: RouteConfig = {
  path: '/mqtt',
  component: Develop,
  redirect: '/mqtt',
  name: 'mqtt',
  meta: { hidden: false },
  children: [
    {
      path: '',
      component: () => import('@/views/mqtt/index.vue'),
      name: 'mqtt list',
      meta: {
        title: 'mqtt',
        icon: 'table',
        noCache: true
      }
    }
  ]
}

export default mqttRouter
