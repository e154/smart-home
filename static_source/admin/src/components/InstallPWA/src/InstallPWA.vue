<script setup lang="ts">
import {propTypes} from "@/utils/propTypes";
import {ElPopover} from 'element-plus'
import {useDesign} from "@/hooks/web/useDesign";

const {getPrefixCls} = useDesign()
const prefixCls = getPrefixCls('install-pwa')

defineProps({
  color: propTypes.string.def('')
})

let deferredPrompt = null;
window.addEventListener('beforeinstallprompt', (e) => {
  e.preventDefault();
  deferredPrompt = e;
});

const toggle = async () => {
  if (deferredPrompt) {
    deferredPrompt.prompt();
    const {outcome} = await deferredPrompt.userChoice;
    if (outcome === 'accepted') {
      deferredPrompt = null;
    }
  }
}

</script>

<template>
  <div @click="toggle">
    <ElPopover
      placement="bottom"
      :width="200"
      trigger="hover"
      content="Install application">
      <template #reference>
        <span :class="prefixCls" :color="color">
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24">
          <path fill="currentColor"
                d="M14.83 9L16 10.17zM4 17h16v-3.17l-3 3L9.17 9L13 5.17V5H4z"
                opacity="0.3"/>
          <path fill="currentColor"
                d="M20 17H4V5h9V3H4c-1.11 0-2 .89-2 2v12a2 2 0 0 0 2 2h4v2h8v-2h4c1.1 0 2-.9 2-2v-5.17l-2 2z"/>
          <path fill="currentColor" d="M18 10.17V3h-2v7.17l-2.59-2.58L12 9l5 5l5-5l-1.41-1.41z"/>
        </svg>
        </span>
      </template>

    </ElPopover>
  </div>
</template>

<style lang="less" scoped>
svg, path {
  color: var(--top-header-text-color) !important;
}
</style>
