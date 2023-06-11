import Vue, { DirectiveOptions } from 'vue'

import 'normalize.css'
import ElementUI from 'element-ui'
import SvgIcon from 'vue-svgicon'

import '@/styles/element-variables.scss'
import '@/styles/index.scss'

import App from '@/App.vue'
import store from '@/store'
import { AppModule } from '@/store/modules/app'
import router from '@/router'
import i18n from '@/lang'
import '@/icons/components'
import '@/permission'
import '@/utils/error-log'
import '@/pwa/register-service-worker'
import * as directives from '@/directives'
import * as filters from '@/filters'

(function() {
  if (window.console && typeof console.log !== 'undefined') {
    const url = 'https://github.com/e154/smart-home'
    const i = `Modern smart systems for life - ${url}`
    console.log('%c Smart home %c Copyright © 2014-%s',
      'font-family: "Helvetica Neue", Helvetica, Arial, sans-serif;font-size:62px;color:#303E4D;-webkit-text-fill-color:#303E4D;-webkit-text-stroke: 1px #303E4D;',
      'font-size:12px;color:#a9a9a9;',
      (new Date()).getFullYear())
    console.log(`%c ${i}`, 'color:#333;')
  }
})()

Vue.use(ElementUI, {
  size: AppModule.size, // Set element-ui default size
  i18n: (key: string, value: string) => i18n.t(key, value)
})

Vue.use(SvgIcon, {
  tagName: 'svg-icon',
  defaultWidth: '1em',
  defaultHeight: '1em'
})

// Register global directives
Object.keys(directives).forEach(key => {
  Vue.directive(key, (directives as { [key: string ]: DirectiveOptions })[key])
})

// Register global filter functions
Object.keys(filters).forEach(key => {
  Vue.filter(key, (filters as { [key: string ]: Function })[key])
})

Vue.config.productionTip = false
Vue.config.devtools = true

new Vue({
  router,
  store,
  i18n,
  render: (h) => h(App)
}).$mount('#app')
