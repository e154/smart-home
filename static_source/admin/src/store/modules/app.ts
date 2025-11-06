import { defineStore } from 'pinia'
import { store } from '../index'
import { setCssVar, humpToUnderline } from '@/utils'
import { ElMessage } from 'element-plus'
import { ElementPlusSize } from '@/types/elementPlus'
import { useCache } from '@/hooks/web/useCache'
import { LayoutType } from '@/types/layout'
import { ThemeTypes } from '@/types/theme'
import {ApiCurrentUser, ApiScript} from "@/api/stub";
import stream from "@/api/stream";
import pushService from "@/api/pushService";

const { wsCache } = useCache()

interface AppState {
  token: string;
  user: ApiCurrentUser;
  avatar: string;
  breadcrumb: boolean
  breadcrumbIcon: boolean
  collapse: boolean
  uniqueOpened: boolean
  hamburger: boolean
  screenfull: boolean
  size: boolean
  locale: boolean
  tagsView: boolean
  tagsViewIcon: boolean
  logo: boolean
  fixedHeader: boolean
  greyMode: boolean
  systemTheme: boolean
  dynamicRouter: boolean
  pageLoading: boolean
  layout: LayoutType
  title: string
  isDark: boolean
  currentSize: ElementPlusSize
  sizeMap: ElementPlusSize[]
  mobile: boolean
  footer: boolean
  theme: ThemeTypes
  fixedMenu: boolean
  terminal: boolean
  maxZIndex: number
  serverId: string
  lastColors: string[]
  onlineStatus: 'online' | 'offline'
  standalone: boolean
  activeWindow: string
  globalScript: ApiScript | undefined
}

export const useAppStore = defineStore('app', {
  state: (): AppState => {
    const mqStandAlone = '(display-mode: standalone)'
    const standalone = navigator.standalone || window.matchMedia(mqStandAlone).matches
    return {
      token: wsCache.get("accessToken") as string || '',
      user: wsCache.get("currentUser") as ApiCurrentUser,
      avatar: "",
      sizeMap: ['default', 'large', 'small'],
      mobile: false, // 是否是移动端
      title: import.meta.env.VITE_APP_TITLE, // 标题
      pageLoading: false, // 路由跳转loading

      breadcrumb: wsCache.get('breadcrumb') || true, // 面包屑
      breadcrumbIcon: wsCache.get('breadcrumbIcon') || true, // 面包屑图标
      collapse: wsCache.get('collapse') || true, // 折叠菜单
      uniqueOpened: false, // 是否只保持一个子菜单的展开
      hamburger: true, // 折叠图标
      screenfull: true, // 全屏图标
      size: true, // 尺寸图标
      locale: true, // 多语言图标
      tagsView: wsCache.get('tagsView') || true, // 标签页
      tagsViewIcon: wsCache.get('tagsViewIcon') || false, // 是否显示标签图标
      logo: false, // logo
      fixedHeader: false, // 固定toolheader
      footer: false, // 显示页脚
      greyMode: wsCache.get('greyMode') || false, // 是否开始灰色模式，用于特殊悼念日
      systemTheme: wsCache.get('systemTheme') || true,
      dynamicRouter: wsCache.get('dynamicRouter') || false, // 是否动态路由
      fixedMenu: wsCache.get('fixedMenu') || true, // 是否固定菜单
      terminal: wsCache.get('terminal') || false,
      maxZIndex: 10,
      serverId: wsCache.get('serverId') || '',
      layout: wsCache.get('layout') || 'classic', // layout布局
      isDark: wsCache.get('isDark') || false, // 是否是暗黑模式
      currentSize: wsCache.get('currentSize') || 'small', // 组件尺寸
      lastColors: wsCache.get('lastColors') || [],
      onlineStatus: 'offline',
      standalone: standalone,
      activeWindow: '',
      theme: wsCache.get('theme') || {
        // 主题色
        elColorPrimary: '#409eff',
        // 左侧菜单边框颜色
        leftMenuBorderColor: 'inherit',
        // 左侧菜单背景颜色
        leftMenuBgColor: '#001529',
        // 左侧菜单浅色背景颜色
        leftMenuBgLightColor: '#0f2438',
        // 左侧菜单选中背景颜色
        leftMenuBgActiveColor: 'var(--el-color-primary)',
        // 左侧菜单收起选中背景颜色
        leftMenuCollapseBgActiveColor: 'var(--el-color-primary)',
        // 左侧菜单字体颜色
        leftMenuTextColor: '#bfcbd9',
        // 左侧菜单选中字体颜色
        leftMenuTextActiveColor: '#fff',
        // logo字体颜色
        logoTitleTextColor: '#fff',
        // logo边框颜色
        logoBorderColor: 'inherit',
        // 头部背景颜色
        topHeaderBgColor: '#fff',
        // 头部字体颜色
        topHeaderTextColor: 'inherit',
        // 头部悬停颜色
        topHeaderHoverColor: '#f6f6f6',
        // 头部边框颜色
        topToolBorderColor: '#eee'
      },
      globalScript: wsCache.get('globalScript') as ApiScript || undefined,
    }
  },
  getters: {
    getUser() {
      return this.user;
    },
    getToken() {
      return this.token;
    },
    getAvatar() {
      return wsCache.get("avatar") || this.avatar;
    },
    getBreadcrumb(): boolean {
      return this.breadcrumb
    },
    getBreadcrumbIcon(): boolean {
      return this.breadcrumbIcon
    },
    getCollapse(): boolean {
      return this.collapse
    },
    getUniqueOpened(): boolean {
      return this.uniqueOpened
    },
    getHamburger(): boolean {
      return this.hamburger
    },
    getScreenfull(): boolean {
      return this.screenfull
    },
    getSize(): boolean {
      return this.size
    },
    getLocale(): boolean {
      return this.locale
    },
    getTagsView(): boolean {
      return this.tagsView
    },
    getTagsViewIcon(): boolean {
      return this.tagsViewIcon
    },
    getLogo(): boolean {
      return this.logo
    },
    getFixedHeader(): boolean {
      return this.fixedHeader
    },
    getGreyMode(): boolean {
      return this.greyMode
    },
    getSystemTheme(): boolean {
      return this.systemTheme
    },
    getDynamicRouter(): boolean {
      return this.dynamicRouter
    },
    getFixedMenu(): boolean {
      return this.fixedMenu
    },
    getTerminal(): boolean {
      return this.terminal
    },
    getPageLoading(): boolean {
      return this.pageLoading
    },
    getLayout(): LayoutType {
      return this.layout
    },
    getTitle(): string {
      return this.title
    },
    getIsDark(): boolean {
      return this.isDark
    },
    getCurrentSize(): ElementPlusSize {
      return this.currentSize
    },
    getSizeMap(): ElementPlusSize[] {
      return this.sizeMap
    },
    getMobile(): boolean {
      return this.mobile
    },
    getTheme(): ThemeTypes {
      return this.theme
    },
    getFooter(): boolean {
      return this.footer
    },
    getServerId(): string {
      return this.serverId
    },
    getIsGate(): boolean {
      return window?.app_settings?.run_mode == 'gate'
    },
    getLastColors(): string[] {
      return this.lastColors
    },
    getOnlineStatus(): 'online' | 'offline' {
      return this.onlineStatus
    },
    getStandalone(): boolean {
      return this.standalone
    },
    getActiveWindow(): string {
      return this.activeWindow
    },
    getGlobalScript(): ApiScript | undefined {
      return this.globalScript
    },
  },
  actions: {
    SetUser(user: ApiCurrentUser) {
      wsCache.set("currentUser", user)
      this.user = user;
    },
    SetAvatar(avatar: string) {
      wsCache.set("avatar", avatar)
      this.avatar = avatar;
    },
    SetToken(token: string) {
      wsCache.set("accessToken", token)
      this.token = token;

      pushService.shutdown()
      stream.disconnect()

      if (token) {
        // ws
        stream.connect(import.meta.env.VITE_API_BASEPATH as string || window.location.origin, token);
        // push
        pushService.start()
      }
    },
    RemoveToken() {
      stream.disconnect();
      pushService.shutdown();
      wsCache.delete('accessToken')
      wsCache.delete('currentUser')
      wsCache.delete('avatar')
      this.user = null;
      this.token = null;
    },
    setBreadcrumb(breadcrumb: boolean) {
      wsCache.set('breadcrumb', breadcrumb)
      this.breadcrumb = breadcrumb
    },
    setBreadcrumbIcon(breadcrumbIcon: boolean) {
      wsCache.set('breadcrumbIcon', breadcrumbIcon)
      this.breadcrumbIcon = breadcrumbIcon
    },
    setCollapse(collapse: boolean) {
      wsCache.set('collapse', collapse)
      this.collapse = collapse
    },
    setUniqueOpened(uniqueOpened: boolean) {
      this.uniqueOpened = uniqueOpened
    },
    setHamburger(hamburger: boolean) {
      this.hamburger = hamburger
    },
    setScreenfull(screenfull: boolean) {
      this.screenfull = screenfull
    },
    setSize(size: boolean) {
      this.size = size
    },
    setLocale(locale: boolean) {
      this.locale = locale
    },
    setTagsView(tagsView: boolean) {
      wsCache.set('tagsView', tagsView)
      this.tagsView = tagsView
    },
    setTagsViewIcon(tagsViewIcon: boolean) {
      wsCache.set('tagsViewIcon', tagsViewIcon)
      this.tagsViewIcon = tagsViewIcon
    },
    setLogo(logo: boolean) {
      this.logo = logo
    },
    setFixedHeader(fixedHeader: boolean) {
      this.fixedHeader = fixedHeader
    },
    setGreyMode(greyMode: boolean) {
      wsCache.set('greyMode', greyMode)
      this.greyMode = greyMode
    },
    setSystemTheme(systemTheme: boolean) {
      wsCache.set('systemTheme', systemTheme)
      this.systemTheme = systemTheme
    },
    setDynamicRouter(dynamicRouter: boolean) {
      wsCache.set('dynamicRouter', dynamicRouter)
      this.dynamicRouter = dynamicRouter
    },
    setFixedMenu(fixedMenu: boolean) {
      wsCache.set('fixedMenu', fixedMenu)
      this.fixedMenu = fixedMenu
    },
    setTerminal(terminal: boolean) {
      wsCache.set('terminal', terminal)
      this.terminal = terminal
    },
    setPageLoading(pageLoading: boolean) {
      this.pageLoading = pageLoading
    },
    setLayout(layout: LayoutType) {
      if (this.mobile && layout !== 'classic') {
        ElMessage.warning('移动端模式下不支持切换其他布局')
        return
      }
      this.layout = layout
      wsCache.set('layout', this.layout)
    },
    setTitle(title: string) {
      this.title = title
    },
    setIsDark(isDark: boolean) {
      this.isDark = isDark
      if (this.isDark) {
        document.documentElement.classList.add('dark')
        document.documentElement.classList.remove('light')
      } else {
        document.documentElement.classList.add('light')
        document.documentElement.classList.remove('dark')
      }
      wsCache.set('isDark', this.isDark)
    },
    setCurrentSize(currentSize: ElementPlusSize) {
      this.currentSize = currentSize
      wsCache.set('currentSize', this.currentSize)
    },
    setMobile(mobile: boolean) {
      this.mobile = mobile
    },
    setTheme(theme: ThemeTypes) {
      this.theme = Object.assign(this.theme, theme)
      wsCache.set('theme', this.theme)
    },
    setCssVarTheme() {
      for (const key in this.theme) {
        setCssVar(`--${humpToUnderline(key)}`, this.theme[key])
      }
    },
    setFooter(footer: boolean) {
      this.footer = footer
    },
    setServerId(id: string) {
      this.serverId = id
      wsCache.set('serverId', this.serverId)
    },
    getMaxZIndex(): number {
      return ++this.maxZIndex
    },
    setLastColors(list: string[]) {
      this.lastColors = list
      wsCache.set('lastColors', this.lastColors)
    },
    setOnlineStatus(status: string) {
      this.onlineStatus = status
    },
    setActiveWindow(name: string) {
      this.activeWindow = name
    },
    setGlobalScript(script: ApiScript) {
      // this.globalScript = Object.assign(this.globalScript, script)
      this.globalScript = script
      wsCache.set('globalScript', this.globalScript)
    }
  }
})

export const useAppStoreWithOut = () => {
  return useAppStore(store)
}
