import { RouteConfig } from 'vue-router'
import Develop from '@/layout/develop.vue'

const automationRouter: RouteConfig = {
  path: '/automation',
  component: Develop,
  redirect: '/automation',
  name: 'automation',
  meta: { hidden: false },
  children: [
    {
      path: '',
      component: () => import('@/views/automation/index.vue'),
      name: 'task list',
      meta: {
        icon: 'automation',
        title: 'taskList'
      }
    },
    {
      path: 'edit/:id',
      component: () => import('@/views/automation/edit.vue'),
      props: true,
      name: 'task edit',
      meta: {
        title: 'taskEdit',
        hidden: true
      }
    },
    {
      path: 'new',
      component: () => import('@/views/automation/new.vue'),
      props: true,
      name: 'task new',
      meta: {
        title: 'taskNew',
        hidden: true
      }
    }
  ]
}

export default automationRouter
