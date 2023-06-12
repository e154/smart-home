import {RouteConfig} from 'vue-router';
import Develop from '@/layout/develop.vue';

const etcRouter: RouteConfig = {
  path: '/etc',
  component: Develop,
  redirect: '/etc',
  name: 'etc',
  meta: {
    icon: 'settings',
    title: 'etc',
    noCache: true,
  },
  children: [
    {
      path: 'variables',
      component: () => import('@/views/variables/index.vue'),
      name: 'variables',
      meta: {
        title: 'variables'
      }
    },
    {
      path: 'plugins',
      component: () => import('@/views/plugins/index.vue'),
      name: 'plugins',
      meta: {
        title: 'plugins'
      }
    },
    {
      path: 'swagger',
      component: () => import('@/views/swagger/index.vue'),
      name: 'swagger',
      meta: {
        title: 'swagger'
      }
    },
    {
      path: 'images',
      component: () => import('@/views/images/index.vue'),
      name: 'images',
      meta: {
        title: 'images'
      }
    },
    {
      path: 'areas',
      component: () => import('@/views/areas/index.vue'),
      name: 'area list',
      meta: {
        title: 'areaList'
      }
    },
    {
      path: 'areas/edit/:id',
      component: () => import('@/views/areas/edit.vue'),
      props: true,
      name: 'area edit',
      meta: {
        title: 'areaEdit',
        hidden: true
      }
    },
    {
      path: 'areas/new',
      component: () => import('@/views/areas/new.vue'),
      props: true,
      name: 'area new',
      meta: {
        title: 'areaNew',
        hidden: true
      }
    },
    {
      path: 'users',
      component: () => import('@/views/users/index.vue'),
      name: 'user list',
      meta: {
        title: 'userList'
      }
    },
    {
      path: 'users/edit/:id',
      component: () => import('@/views/users/edit.vue'),
      props: true,
      name: 'user edit',
      meta: {
        title: 'UserEdit',
        hidden: true
      }
    },
    {
      path: 'users/new',
      component: () => import('@/views/users/new.vue'),
      name: 'user new',
      meta: {
        hidden: true,
        title: 'UserNew'
      }
    },
    {
      path: 'backups',
      component: () => import('@/views/backup/index.vue'),
      name: 'backup list',
      meta: {
        title: 'backupList'
      }
    },
    {
      path: 'message_delivery',
      component: () => import('@/views/message_delivery/index.vue'),
      name: 'Message Delivery',
      meta: {
        title: 'messageDelivery'
      }
    }
  ]
};

export default etcRouter;
