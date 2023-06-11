import {RouteConfig} from 'vue-router';
import Develop from '@/layout/develop.vue';

const entitiesRouter: RouteConfig = {
  path: '/entities',
  component: Develop,
  redirect: '/entities',
  name: 'entities',
  meta: { hidden: false },
  children: [
    {
      path: '',
      component: () => import('@/views/entities/index.vue'),
      name: 'entity list',
      meta: {
        icon: 'entity2',
        title: 'entityList'
      }
    },
    {
      path: 'edit/:id',
      component: () => import('@/views/entities/edit.vue'),
      props: true,
      name: 'entity edit',
      meta: {
        title: 'entityEdit',
        hidden: true
      }
    },
    {
      path: 'new',
      component: () => import('@/views/entities/new.vue'),
      props: true,
      name: 'entity new',
      meta: {
        title: 'entityNew',
        hidden: true
      }
    }
  ]
};

export default entitiesRouter;
