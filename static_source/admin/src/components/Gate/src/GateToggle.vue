<script setup lang="ts">
import {propTypes} from "@/utils/propTypes";
import {ElPopover} from 'element-plus'
import {computed} from "vue";
import {useAppStore} from "@/store/modules/app";
import {useDesign} from "@/hooks/web/useDesign";

const {getPrefixCls} = useDesign()
const prefixCls = getPrefixCls('gate-toggle')

defineProps({
  color: propTypes.string.def('')
})

const appStore = useAppStore()

const serverId = computed(() => appStore.getServerId)

const toggle = () => {

}

</script>

<template>
  <div @click="toggle">
    <ElPopover
      placement="bottom"
      title="Server Id"
      :width="200"
      trigger="hover"
      :content="serverId">
      <template #reference>
        <span :class="prefixCls" :color="color">
          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24">
            <path fill="currentColor"
                  d="M8 6a1 1 0 1 0-1-1a1 1 0 0 0 1 1m13 13h-6.18A3 3 0 0 0 13 17.18V15h4a3 3 0 0 0 3-3v-2a3 3 0 0 0-.78-2A3 3 0 0 0 20 6V4a3 3 0 0 0-3-3H7a3 3 0 0 0-3 3v2a3 3 0 0 0 .78 2A3 3 0 0 0 4 10v2a3 3 0 0 0 3 3h4v2.18A3 3 0 0 0 9.18 19H3a1 1 0 0 0 0 2h6.18a3 3 0 0 0 5.64 0H21a1 1 0 0 0 0-2M6 4a1 1 0 0 1 1-1h10a1 1 0 0 1 1 1v2a1 1 0 0 1-1 1H7a1 1 0 0 1-1-1Zm1 9a1 1 0 0 1-1-1v-2a1 1 0 0 1 1-1h10a1 1 0 0 1 1 1v2a1 1 0 0 1-1 1Zm5 8a1 1 0 1 1 1-1a1 1 0 0 1-1 1M8 10a1 1 0 1 0 1 1a1 1 0 0 0-1-1"/>
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
