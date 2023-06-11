import { RouteConfig } from 'vue-router'
import Develop from '@/layout/develop.vue'

const scriptsRouter: RouteConfig = {
  path: '/scripts',
  component: Develop,
  redirect: '/scripts/',
  name: 'scripts',
  meta: { hidden: false },
  children: [
    {
      path: '',
      component: () => import('@/views/scripts/index.vue'),
      name: 'script list',
      meta: {
        icon: 'script-13',
        title: 'scriptList'
      }
    },
    {
      path: 'edit/:id',
      component: () => import('@/views/scripts/edit.vue'),
      props: true,
      name: 'script edit',
      meta: {
        title: 'scriptEdit',
        hidden: true
      }
    },
    {
      path: 'new',
      component: () => import('@/views/scripts/new.vue'),
      props: true,
      name: 'script new',
      meta: {
        title: 'scriptNew',
        hidden: true
      }
    }

  ]
}

export default scriptsRouter
