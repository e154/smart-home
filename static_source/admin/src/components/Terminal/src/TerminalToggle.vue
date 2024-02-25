<script setup lang="ts">
import {useAppStore} from "@/store/modules/app";
import {propTypes} from "@/utils/propTypes";
import {useDesign} from "@/hooks/web/useDesign";

const {getPrefixCls} = useDesign()
const prefixCls = getPrefixCls('terminal')

const appStore = useAppStore()

defineProps({
  color: propTypes.string.def('')
})

const toggleFullscreen = () => {
  appStore.setTerminal(!appStore.getTerminal)
}

</script>

<template>
  <span @click="toggleFullscreen" :class="prefixCls" :color="color">
    <svg xmlns="http://www.w3.org/2000/svg" role="img" width="18" height="18" viewBox="0 0 16 16">
      <g>
        <path fill="currentColor"
          d="M6 9a.5.5 0 0 1 .5-.5h3a.5.5 0 0 1 0 1h-3A.5.5 0 0 1 6 9M3.854 4.146a.5.5 0 1 0-.708.708L4.793 6.5L3.146 8.146a.5.5 0 1 0 .708.708l2-2a.5.5 0 0 0 0-.708z"/>
        <path fill="currentColor"
          d="M2 1a2 2 0 0 0-2 2v10a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V3a2 2 0 0 0-2-2zm12 1a1 1 0 0 1 1 1v10a1 1 0 0 1-1 1H2a1 1 0 0 1-1-1V3a1 1 0 0 1 1-1z"/>
      </g>
    </svg>
  </span>
</template>

<style lang="less" scoped>
svg, path {
  color: var(--top-header-text-color) !important;
}
</style>
