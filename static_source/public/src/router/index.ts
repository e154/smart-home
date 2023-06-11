import Vue from 'vue';
import VueRouter, {RouteConfig} from 'vue-router';

/* Develop */
import Develop from '@/layout/develop.vue';
import Dashboard from '@/layout/dashboard.vue';

/* Router modules */
import scriptsRouter from './modules/scripts';
import entitiesRouter from '@/router/modules/entities';
import automationRouter from '@/router/modules/automation';
import zigbee2mqttRouter from '@/router/modules/zigbee2mqtt';
import logsRouter from '@/router/modules/log';
import etcRouter from '@/router/modules/etc';
import dashboardsRouter from '@/router/modules/dashboard';

Vue.use(VueRouter);

/*
  Note: sub-menu only appear when children.length>=1
  Detail see: https://panjiachen.github.io/vue-element-admin-site/guide/essentials/router-and-nav.html
*/

/*
  name:'router-name'             the name field is required when using <keep-alive>, it should also match its component's name property
                                 detail see : https://vuejs.org/v2/guide/components-dynamic-async.html#keep-alive-with-Dynamic-Components
  redirect:                      if set to 'noredirect', no redirect action will be trigger when clicking the breadcrumb
  meta: {
    roles: ['admin', 'editor']   will control the page roles (allow setting multiple roles)
    title: 'title'               the name showed in subMenu and breadcrumb (recommend set)
    icon: 'svg-name'             the icon showed in the sidebar
    hidden: true                 if true, this route will not show in the sidebar (default is false)
    alwaysShow: true             if true, will always show the root menu (default is false)
                                 if false, hide the root menu when has less or equal than one children route
    breadcrumb: false            if false, the item will be hidden in breadcrumb (default is true)
    noCache: true                if true, the page will not be cached (default is false)
    affix: true                  if true, the tag will affix in the tags-view
    activeMenu: '/example/list'  if set path, the sidebar will highlight the path you set
  }
*/

/**
 ConstantRoutes
 a base page that does not have permission requirements
 all roles can be accessed
 */
export const constantRoutes: RouteConfig[] = [
  {
    path: '/',
    component: Dashboard,
    redirect: '/board',
    children: [
      {
        path: 'board',
        component: () => import(/* webpackChunkName: "dashboard" */ '@/views/dashboard/view.vue'),
        name: 'dashboard view',
        meta: {
          title: 'dashboard',
          icon: 'dashboard',
          affix: true
        }
      }
    ]
  },
  {
    path: '/redirect',
    component: Develop,
    meta: {hidden: true},
    children: [
      {
        path: '/redirect/:path(.*)',
        component: () => import(/* webpackChunkName: "redirect" */ '@/views/redirect/index.vue')
      }
    ]
  },
  {
    path: '/login',
    component: () => import(/* webpackChunkName: "login" */ '@/views/login/index.vue'),
    meta: {hidden: true}
  },
  {
    path: '/auth-redirect',
    component: () => import(/* webpackChunkName: "auth-redirect" */ '@/views/login/auth-redirect.vue'),
    meta: {hidden: true}
  },
  {
    path: '/404',
    component: Dashboard,
    meta: {hidden: true},
    redirect: '/404',
    children: [
      {
        path: '',
        component: () => import(/* webpackChunkName: "404" */ '@/views/error-page/404.vue'),
        name: '404',
        meta: {hidden: true}
      }
    ]
  },
  {
    path: '/401',
    component: () => import(/* webpackChunkName: "401" */ '@/views/error-page/401.vue'),
    meta: {hidden: true}
  }
];

/**
 * dashboardRoutes
 * the routes that need to be dynamically loaded based on user roles
 */
export const dashboardRoutes: RouteConfig[] = [
  {
    path: '/development',
    component: Develop,
    redirect: '/development/index',
    meta: {hidden: false},
    children: [
      {
        path: 'index',
        component: () => import(/* webpackChunkName: "development" */ '@/views/development/index.vue'),
        name: 'Development',
        meta: {
          title: 'development',
          icon: 'development-kit',
          noCache: true
        }
      }
    ]
  },
  {
    path: '/profile',
    component: Dashboard,
    redirect: '/profile/index',
    meta: {hidden: true},
    children: [
      {
        path: 'index',
        component: () => import(/* webpackChunkName: "profile" */ '@/views/profile/index.vue'),
        name: 'Profile',
        meta: {
          title: 'profile',
          icon: 'user',
          noCache: true
        }
      }
    ]
  },
  // {
  //   path: '*',
  //   redirect: '/404',
  //   meta: {hidden: true}
  // }
];

/**
 * developRoutes
 * the routes that need to be dynamically loaded based on user roles
 */
export const developRoutes: RouteConfig[] = [
  {
    path: '/development',
    component: Develop,
    redirect: '/development/index',
    meta: {hidden: false},
    children: [
      {
        path: 'index',
        component: () => import(/* webpackChunkName: "development" */ '@/views/development/index.vue'),
        name: 'Development',
        meta: {
          title: 'development',
          icon: 'development-kit',
          noCache: true
        }
      }
    ]
  },
  entitiesRouter,
  scriptsRouter,
  automationRouter,
  zigbee2mqttRouter,
  dashboardsRouter,
  logsRouter,
  etcRouter,
  {
    path: '/profile',
    component: Dashboard,
    redirect: '/profile/index',
    meta: {hidden: true},
    children: [
      {
        path: 'index',
        component: () => import(/* webpackChunkName: "profile" */ '@/views/profile/index.vue'),
        name: 'Profile',
        meta: {
          title: 'profile',
          icon: 'user',
          noCache: true
        }
      }
    ]
  },
  {
    path: '*',
    redirect: '/404',
    meta: {hidden: true}
  }
];

const createRouter = () => new VueRouter({
  // mode: 'history',  // Disabled due to Github Pages doesn't support this, enable this if you need.
  scrollBehavior: (to, from, savedPosition) => {
    if (savedPosition) {
      return savedPosition;
    } else {
      return {x: 0, y: 0};
    }
  },
  base: process.env.BASE_URL,
  routes: constantRoutes
});

const router = createRouter();

// Detail see: https://github.com/vuejs/vue-router/issues/1234#issuecomment-357941465
export function resetRouter() {
  const newRouter = createRouter();
  (router as any).matcher = (newRouter as any).matcher; // reset router
}

export default router;
