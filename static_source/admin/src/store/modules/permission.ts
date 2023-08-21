import { defineStore } from 'pinia'
import { dashboardRouterMap, developRouterMap, constantRouterMap } from '@/router'
import { generateRoutesFn1, generateRoutesFn2, flatMultiLevelRoutes } from '@/utils/routerHelper'
import { store } from '../index'
import { cloneDeep } from 'lodash-es'

export interface PermissionState {
  routers: AppRouteRecordRaw[]
  addRouters: AppRouteRecordRaw[]
  isAddRouters: boolean
  menuTabRouters: AppRouteRecordRaw[]
}

export const usePermissionStore = defineStore('permission', {
  state: (): PermissionState => ({
    routers: [],
    addRouters: [],
    isAddRouters: false,
    menuTabRouters: []
  }),
  getters: {
    getRouters(): AppRouteRecordRaw[] {
      return this.routers
    },
    getAddRouters(): AppRouteRecordRaw[] {
      return flatMultiLevelRoutes(cloneDeep(this.addRouters))
    },
    getIsAddRouters(): boolean {
      return this.isAddRouters
    },
    getMenuTabRouters(): AppRouteRecordRaw[] {
      return this.menuTabRouters
    },
    getDashboardRouters(): AppRouteRecordRaw[] {
      return this.dashboardTabRouters
    },
    getDevelopmentRouters(): AppRouteRecordRaw[] {
      return this.routers
    }
  },
  actions: {
    generateRoutes(
      type: 'admin' | 'user' | 'none',
      routers?: AppCustomRouteRecordRaw[] | string[]
    ): Promise<unknown> {
      return new Promise<void>((resolve) => {
        let routerMap: AppRouteRecordRaw[] = []
        if (type === 'admin') {
          // 模拟后端过滤菜单
          routerMap = generateRoutesFn2(routers as AppCustomRouteRecordRaw[])
        } else if (type === 'user') {
          // 模拟前端过滤菜单
          routerMap = generateRoutesFn1(cloneDeep(dashboardRouterMap), routers as string[])
        } else {
          // 直接读取静态路由表
          routerMap = cloneDeep(dashboardRouterMap)
        }
        routerMap = routerMap.concat(cloneDeep(developRouterMap))

        // 动态路由，404一定要放到最后面
        const err404: AppRouteRecordRaw[] = [
          {
            path: '/:path(.*)*',
            redirect: '/404',
            name: '404Page',
            meta: {
              hidden: true,
              breadcrumb: false
            }
          } as AppRouteRecordRaw
        ]
        this.addRouters = routerMap.concat(err404)

        // 渲染菜单的所有路由
        this.routers = cloneDeep(constantRouterMap).concat(routerMap)

        this.dashboardTabRouters = cloneDeep(constantRouterMap)
        this.dashboardTabRouters = this.dashboardTabRouters.concat(err404, cloneDeep(dashboardRouterMap))
        resolve()
      })
    },
    setIsAddRouters(state: boolean): void {
      this.isAddRouters = state
    },
    setMenuTabRouters(routers: AppRouteRecordRaw[]): void {
      this.menuTabRouters = routers
    }
  }
})

export const usePermissionStoreWithOut = () => {
  return usePermissionStore(store)
}
