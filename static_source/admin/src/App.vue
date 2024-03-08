<script setup lang="ts">
import {computed} from 'vue'
import { useAppStore } from '@/store/modules/app'
import { ConfigGlobal } from '@/components/ConfigGlobal'
import { isDark } from '@/utils/is'
import { useDesign } from '@/hooks/web/useDesign'
import { useCache } from '@/hooks/web/useCache'
import {ReloadPrompt} from '@/components/ReloadPrompt'
import stream from "@/api/stream";
import pushService from "@/api/pushService";

const { getPrefixCls } = useDesign()

const prefixCls = getPrefixCls('app')

const appStore = useAppStore()

const currentSize = computed(() => appStore.getCurrentSize)

const greyMode = computed(() => appStore.getGreyMode)

const { wsCache } = useCache()

const accessToken = wsCache.get("accessToken") as string || '';
if (accessToken) {
  // ws
  stream.connect(import.meta.env.VITE_API_BASEPATH as string || window.location.origin, accessToken);
  // push
  pushService.start()
}

// 根据浏览器当前主题设置系统主题色
const setDefaultTheme = () => {
  if (wsCache.get('isDark') !== null) {
    appStore.setIsDark(wsCache.get('isDark'))
    return
  }
  const isDarkTheme = isDark()
  appStore.setIsDark(isDarkTheme)
}

const consoleBanner = () => {
  var i, url;
  if (window.console && 'undefined' !== typeof console.log) {
    url = 'https://github.com/e154/smart-home';
    i = `Software package for automation - ${url}`;
    console.log('%c Smart home %c Copyright © 2014-%s', 'font-family: "Helvetica Neue", Helvetica, Arial, sans-serif;font-size:62px;color:#303E4D;-webkit-text-fill-color:#303E4D;-webkit-text-stroke: 1px #303E4D;', 'font-size:12px;color:#a9a9a9;', (new Date()).getFullYear());
    return console.log('%c ' + i, 'color:#333;');
  }
}

setDefaultTheme()
consoleBanner()
</script>

<template>
  <ReloadPrompt />
  <ConfigGlobal :size="currentSize">
    <RouterView :class="greyMode ? `${prefixCls}-grey-mode` : ''" />
  </ConfigGlobal>
</template>

<style lang="less">
@prefix-cls: ~'@{namespace}-app';

.size {
  width: 100%;
  height: 100%;
}

html,
body {
  padding: 0 !important;
  margin: 0;
  overflow: hidden;
  .size;

  #app {
    .size;
  }
}

.@{prefix-cls}-grey-mode {
  filter: grayscale(100%);
}
</style>
