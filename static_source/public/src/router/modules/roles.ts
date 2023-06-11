import { RouteConfig } from 'vue-router'
import Develop from '@/layout/develop.vue'

const rolesRouter: RouteConfig = {
  path: '/roles',
  component: Develop,
  redirect: '/roles',
  name: 'roles',
  meta: {
    title: 'roles',
    icon: 'role',
    alwaysShow: true
  },
  children: [
    {
      path: '',
      component: () => import('@/views/roles/index.vue'),
      name: 'role list',
      meta: {
        title: 'RoleList'
      }
    }
  ]
}

export default rolesRouter
