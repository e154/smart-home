import { createRouter, createWebHashHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import type { App } from 'vue'
import { Dashboard, Develop, Landing, getParentLayout } from '@/utils/routerHelper'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n()

export const constantRouterMap: AppRouteRecordRaw[] = [
  {
    path: '/',
    component: Dashboard,
    redirect: '/board',
    meta: {
      title: t('router.Dashboard'),
      icon: 'vaadin:dashboard'
    },
    children: [
      {
        path: '/board',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard/main.vue'),
        meta: {
          hidden: true,
          title: t('router.Dashboard'),
          noTagsView: true
        }
      }
    ]
  },
  {
    path: '/landing',
    component: Landing,
    redirect: 'index',
    meta: {
      title: t('router.Landing'),
      hidden: true,
      noTagsView: true
    },
    children: [
      {
        path: '/landing/:id',
        name: 'Landing',
        component: () => import('@/views/Dashboard/landing.vue'),
        meta: {
          hidden: true,
          title: t('router.Dashboard'),
          noTagsView: true
        }
      }
    ]
  },
  {
    path: '/redirect',
    component: Dashboard,
    name: 'Redirect',
    children: [
      {
        path: '/redirect/:path(.*)',
        name: 'Redirect',
        component: () => import('@/views/Redirect/Redirect.vue'),
        meta: {}
      }
    ],
    meta: {
      hidden: true,
      noTagsView: true
    }
  },
  {
    path: '/login',
    component: () => import('@/views/Login/Login.vue'),
    name: 'Login',
    meta: {
      hidden: true,
      title: t('router.login'),
      noTagsView: true
    }
  },
  {
    path: '/password_reset',
    component: () => import('@/views/PasswordReset/PasswordReset.vue'),
    name: 'Password Reset',
    meta: {
      hidden: true,
      title: t('router.passwordReset'),
      noTagsView: true
    }
  },
  {
    path: '/404',
    component: () => import('@/views/Error/404.vue'),
    name: 'NoFind',
    meta: {
      hidden: true,
      title: '404',
      noTagsView: true
    }
  }
]

export const dashboardRouterMap: AppRouteRecordRaw[] = [
  {
    path: '/development',
    component: Develop,
    redirect: '/development/index',
    meta: {
      title: t('router.Development'),
    },
    children: [
      {
        path: 'index',
        name: 'Development',
        component: () => import('@/views/Development/index.vue'),
        meta: {
          title: t('router.Development'),
          icon: 'mdi:tools'
        }
      }
    ]
  }
]

export const developRouterMap: AppRouteRecordRaw[] = [
  {
    path: '/entities',
    component: Develop,
    redirect: '/entities/index',
    meta: {
      title: t('router.Entities'),
      icon: 'icon-park-solid:layers'
    },
    children: [
      {
        path: 'index',
        name: 'Entities',
        component: () => import('@/views/Entities/index.vue'),
        meta: {
          title: t('router.ListEntities'),
          // noTagsView: true,
          hidden: true,
          canTo: true,
          activeMenu: '/entities'
        }
      },
      {
        path: 'new',
        name: 'newEntities',
        component: () => import('@/views/Entities/new.vue'),
        meta: {
          title: t('router.New'),
          noTagsView: true,
          hidden: true,
          canTo: true,
          activeMenu: '/entities'
        }
      },
      {
        path: 'edit/:id',
        name: 'editEntities',
        component: () => import('@/views/Entities/edit.vue'),
        props: true,
        meta: {
          title: t('router.EditEntity'),
          // noTagsView: true,
          hidden: true,
          canTo: true,
          activeMenu: '/entities'
        }
      }
    ]
  },
  {
    path: '/waterfall',
    component: Develop,
    name: 'waterfall',
    redirect: '/waterfall/index',
    meta: {
      hidden: true,
      title: 'waterfall'
    },
    children: [
      {
        path: 'index',
        name: 'waterfall',
        component: () => import('@/views/components/Waterfall.vue'),
        meta: {
          title: t('router.ListScripts'),
          // noTagsView: true,
          hidden: true,
          canTo: true,
          activeMenu: '/waterfall'
        }
      },
    ]
  },
  {
    path: '/image-cropping',
    component: Develop,
    name: 'ImageCropping',
    redirect: '/image-cropping/index',
    meta: {
      hidden: true,
      title: 'ImageCropping'
    },
    children: [
      {
        path: 'index',
        name: 'ImageCropping',
        component: () => import('@/views/components/ImageCropping.vue'),
        meta: {
          title: t('router.imageCropping'),
          // noTagsView: true,
          hidden: true,
          canTo: true,
          activeMenu: '/image-cropping'
        }
      },
    ]
  },
  {
    path: '/scripts',
    component: Develop,
    redirect: '/scripts/index',
    meta: {
      title: t('router.Scripts'),
      icon: 'fluent-mdl2:coffee-script'
    },
    children: [
      {
        path: 'index',
        name: 'Scripts',
        component: () => import('@/views/Scripts/index.vue'),
        meta: {
          title: t('router.ListScripts'),
          // noTagsView: true,
          hidden: true,
          canTo: true,
          activeMenu: '/scripts'
        }
      },
      {
        path: 'new',
        name: 'newScripts',
        component: () => import('@/views/Scripts/new.vue'),
        meta: {
          title: t('router.New'),
          noTagsView: true,
          hidden: true,
          canTo: true,
          activeMenu: '/scripts'
        }
      },
      {
        path: 'edit/:id',
        name: 'editScripts',
        component: () => import('@/views/Scripts/edit.vue'),
        props: true,
        meta: {
          title: t('router.EditScript'),
          // noTagsView: true,
          hidden: true,
          canTo: true,
          activeMenu: '/scripts'
        }
      }
    ]
  },
  {
    path: '/automation',
    component: Develop,
    redirect: '/automation/tasks',
    meta: {
      title: t('router.Automation'),
      icon: 'fa6-solid:circle-nodes'
    },
    children: [
      {
        path: 'tasks',
        name: 'Tasks',
        component: () => import('@/views/Automation/Tasks/index.vue'),
        meta: {
          title: t('router.Tasks'),
        },
        children: [
          {
            path: 'new',
            name: 'newTask',
            component: () => import('@/views/Automation/Tasks/new.vue'),
            meta: {
              title: t('router.New'),
              // noTagsView: true,
              hidden: true,
              canTo: true,
              activeMenu: '/automation/tasks'
            }
          },
          {
            path: 'edit/:id',
            name: 'editTask',
            component: () => import('@/views/Automation/Tasks/edit.vue'),
            props: true,
            meta: {
              title: t('router.EditTask'),
              // noTagsView: true,
              hidden: true,
              canTo: true,
              activeMenu: '/automation/tasks'
            }
          }
        ]
      },
      {
        path: 'triggers',
        name: 'Triggers',
        component: () => import('@/views/Automation/Triggers/index.vue'),
        meta: {
          title: t('router.Triggers'),
        },
        children: [
          {
            path: 'new',
            name: 'newTrigger',
            component: () => import('@/views/Automation/Triggers/new.vue'),
            meta: {
              title: t('router.New'),
              noTagsView: true,
              hidden: true,
              canTo: true,
              activeMenu: '/automation/triggers'
            }
          },
          {
            path: 'edit/:id',
            name: 'editTrigger',
            component: () => import('@/views/Automation/Triggers/edit.vue'),
            props: true,
            meta: {
              title: t('router.EditTrigger'),
              // noTagsView: true,
              hidden: true,
              canTo: true,
              activeMenu: '/automation/triggers'
            }
          }
        ]
      },
      {
        path: 'conditions',
        name: 'Conditions',
        component: () => import('@/views/Automation/Conditions/index.vue'),
        meta: {
          title: t('router.Conditions'),
        },
        children: [
          {
            path: 'new',
            name: 'newCondition',
            component: () => import('@/views/Automation/Conditions/new.vue'),
            meta: {
              title: t('router.New'),
              noTagsView: true,
              hidden: true,
              canTo: true,
              activeMenu: '/automation/conditions'
            }
          },
          {
            path: 'edit/:id',
            name: 'editCondition',
            component: () => import('@/views/Automation/Conditions/edit.vue'),
            props: true,
            meta: {
              title: t('router.EditCondition'),
              // noTagsView: true,
              hidden: true,
              canTo: true,
              activeMenu: '/automation/conditions'
            }
          }
        ]
      },
      {
        path: 'actions',
        name: 'Actions',
        component: () => import('@/views/Automation/Actions/index.vue'),
        meta: {
          title: t('router.Actions'),
        },
        children: [
          {
            path: 'new',
            name: 'newAction',
            component: () => import('@/views/Automation/Actions/new.vue'),
            meta: {
              title: t('router.New'),
              noTagsView: true,
              hidden: true,
              canTo: true,
              activeMenu: '/automation/actions'
            }
          },
          {
            path: 'edit/:id',
            name: 'editAction',
            component: () => import('@/views/Automation/Actions/edit.vue'),
            props: true,
            meta: {
              title: t('router.EditAction'),
              // noTagsView: true,
              hidden: true,
              canTo: true,
              activeMenu: '/automation/actions'
            }
          }
        ]
      },
    ]
  },
  {
    path: '/dashboards',
    component: Develop,
    redirect: '/dashboards/index',
    meta: {
      title: t('router.Dashboards'),
      icon: 'ic:sharp-dashboard-customize'
    },
    children: [
      {
        path: 'index',
        name: 'Dashboards',
        component: () => import('@/views/Dashboard/index.vue'),
        meta: {
          title: t('router.ListDashboards'),
          // noTagsView: true,
          hidden: true,
          canTo: true,
          activeMenu: '/dashboards'
        }
      },
      {
        path: 'view/:id',
        name: 'viewDashboards',
        component: () => import('@/views/Dashboard/view.vue'),
        props: true,
        meta: {
          title: t('router.DashboardsView'),
          // noTagsView: true,
          hidden: true,
          canTo: true,
          activeMenu: '/dashboards'
        }
      },
      {
        path: 'edit/:id',
        name: 'editDashboards',
        component: () => import('@/views/Dashboard/editor/editor.vue'),
        props: true,
        meta: {
          title: t('router.EditDashboard'),
          // noTagsView: true,
          hidden: true,
          canTo: true,
          activeMenu: '/dashboards'
        }
      },
    ]
  },
  {
    path: '/logs',
    component: Develop,
    redirect: '/logs/index',
    meta: {
      title: t('router.Logs'),
    },
    children: [
      {
        path: 'index',
        name: 'Logs',
        component: () => import('@/views/Logs/index.vue'),
        meta: {
          title: t('router.Logs'),
          icon: 'icon-park-outline:upload-logs'
        }
      }
    ]
  },
  {
    path: '/etc',
    component: Develop,
    meta: {
      title: t('router.etc'),
      icon: 'mdi:cog'
    },
    children: [
      {
        path: 'settings',
        component: Develop,
        meta: {
          title: t('router.Settings'),
        },
        children: [
          {
            path: 'main',
            name: 'Main',
            component: () => import('@/views/Settings/index.vue'),
            meta: {
              title: t('router.Main'),
            },
          },
          {
            path: 'plugins',
            name: 'Plugins',
            component: () => import('@/views/Plugins/index.vue'),
            meta: {
              title: t('router.Plugins'),
            },
            children: [
              {
                path: 'edit/:name',
                name: 'viewPlugin',
                component: () => import('@/views/Plugins/edit.vue'),
                props: true,
                meta: {
                  title: t('router.EditPlugin'),
                  // noTagsView: true,
                  hidden: true,
                  canTo: true,
                }
              }
            ]
          },
          {
            path: 'backups',
            name: 'Backups',
            component: () => import('@/views/Backups/index.vue'),
            meta: {
              title: t('router.Backups'),
            },
          },
        ],
      },
      {
        path: 'develop',
        component: Develop,
        meta: {
          title: t('router.Develop'),
        },
        children: [
          {
            path: 'tools',
            name: 'Tools',
            component: () => import('@/views/Tools/EventBus/index.vue'),
            meta: {
              title: t('router.Tools'),
            },
            children: [
            ]
          },
          {
            path: 'swagger',
            name: 'Swagger',
            component: () => import('@/views/Swagger/index.vue'),
            meta: {
              title: t('router.Swagger'),
            },
          },
          {
            path: 'typedoc',
            name: 'Typedoc',
            component: () => import('@/views/Typedoc/index.vue'),
            meta: {
              title: t('router.Typedoc'),
            },
          },
        ],
      },
      {
        path: 'variables',
        name: 'Variables',
        component: () => import('@/views/Variables/index.vue'),
        meta: {
          title: t('router.Variables'),
        },
        children: [
          {
            path: 'new',
            name: 'newVariable',
            component: () => import('@/views/Variables/new.vue'),
            meta: {
              title: t('router.New'),
              noTagsView: true,
              hidden: true,
              canTo: true,
              activeMenu: '/etc/variables'
            }
          },
          {
            path: 'edit/:name',
            name: 'editVariable',
            component: () => import('@/views/Variables/edit.vue'),
            props: true,
            meta: {
              title: t('router.EditVariable'),
              // noTagsView: true,
              hidden: true,
              canTo: true,
              activeMenu: '/etc/variables'
            }
          }
        ]
      },
      {
        path: 'images',
        name: 'Image browser',
        component: () => import('@/views/Images/index.vue'),
        meta: {
          title: t('router.Imagebrowser'),
        },
      },
      {
        path: 'areas',
        name: 'Areas',
        component: () => import('@/views/Areas/index.vue'),
        meta: {
          title: t('router.Areas'),
        },
        children: [
          {
            path: 'new',
            name: 'newArea',
            component: () => import('@/views/Areas/new.vue'),
            meta: {
              title: t('router.New'),
              noTagsView: true,
              hidden: true,
              canTo: true,
              activeMenu: '/etc/areas'
            }
          },
          {
            path: 'edit/:id',
            name: 'editArea',
            component: () => import('@/views/Areas/edit.vue'),
            props: true,
            meta: {
              title: t('router.EditArea'),
              // noTagsView: true,
              hidden: true,
              canTo: true,
              activeMenu: '/etc/areas'
            }
          }
        ]
      },
      {
        path: 'mqtt',
        name: 'Mqtt',
        component: () => import('@/views/Mqtt/index.vue'),
        meta: {
          title: t('router.Mqtt'),
        },
        children: [
          {
            path: 'index',
            name: 'indexMqtt',
            component: () => import('@/views/Mqtt/index.vue'),
            meta: {
              title: t('router.View'),
              // noTagsView: true,
              hidden: true,
              canTo: true,
              activeMenu: '/etc/plugins'
            }
          },
          {
            path: 'client/:id',
            name: 'viewMqttClient',
            component: () => import('@/views/Mqtt/client.vue'),
            props: true,
            meta: {
              title: t('router.View'),
              // noTagsView: true,
              hidden: true,
              canTo: true,
              activeMenu: '/etc/mqtt'
            }
          }
        ]
      },
      {
        path: 'zigbee2mqtt',
        name: 'Zigbee2mqtt',
        component: () => import('@/views/Zigbee2mqtt/index.vue'),
        meta: {
          title: t('router.Zigbee2mqtt'),
        },
        children: [
          {
            path: 'index',
            name: 'Zigbee2mqtt',
            component: () => import('@/views/Zigbee2mqtt/index.vue'),
            meta: {
              title: t('router.ListZigbee2mqtt'),
              // noTagsView: true,
              hidden: true,
              canTo: true,
              activeMenu: '/etc/zigbee2mqtt'
            }
          },
          {
            path: 'new',
            name: 'newZigbee2mqtt',
            component: () => import('@/views/Zigbee2mqtt/new.vue'),
            meta: {
              title: t('router.New'),
              noTagsView: true,
              hidden: true,
              canTo: true,
              activeMenu: '/etc/zigbee2mqtt'
            }
          },
          {
            path: 'edit/:id',
            name: 'editZigbee2mqtt',
            component: () => import('@/views/Zigbee2mqtt/edit.vue'),
            props: true,
            meta: {
              title: t('router.EditZigbee2mqtt'),
              // noTagsView: true,
              hidden: true,
              canTo: true,
              activeMenu: '/etc/zigbee2mqtt'
            }
          }
        ]
      },
      {
        path: 'users',
        name: 'Users',
        component: () => import('@/views/Users/index.vue'),
        meta: {
          title: t('router.Users'),
        },
        children: [
          {
            path: 'new',
            name: 'newUser',
            component: () => import('@/views/Users/new.vue'),
            meta: {
              title: t('router.New'),
              noTagsView: true,
              hidden: true,
              canTo: true,
              activeMenu: '/etc/users'
            }
          },
          {
            path: 'edit/:id',
            name: 'editUser',
            component: () => import('@/views/Users/edit.vue'),
            props: true,
            meta: {
              title: t('router.EditUser'),
              // noTagsView: true,
              hidden: true,
              canTo: true,
              activeMenu: '/etc/users'
            }
          }
        ]
      },
      {
        path: 'message_delivery',
        name: 'MessageDelivery',
        component: () => import('@/views/MessageDelivery/index.vue'),
        meta: {
          title: t('router.MessageDelivery'),
        },
      }
    ]
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  strict: true,
  routes: constantRouterMap as RouteRecordRaw[],
  scrollBehavior: () => ({ left: 0, top: 0 })
})

export const resetRouter = (): void => {
  const resetWhiteNameList = ['Redirect', 'Login', 'NoFind', 'Root']
  router.getRoutes().forEach((route) => {
    const { name } = route
    if (name && !resetWhiteNameList.includes(name as string)) {
      router.hasRoute(name) && router.removeRoute(name)
    }
  })
}

export const setupRouter = (app: App<Element>) => {
  app.use(router)
}

export default router
