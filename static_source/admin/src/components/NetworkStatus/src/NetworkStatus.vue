<script setup lang="ts">
import {propTypes} from "@/utils/propTypes";
import {ElPopover} from 'element-plus'
import {ref, unref, watch} from "vue";
import {streamStatus} from "@/api/stream";
import {useDesign} from "@/hooks/web/useDesign";

const {getPrefixCls} = useDesign()
const prefixCls = getPrefixCls('screenfull')

defineProps({
  color: propTypes.string.def('')
})

const status = ref<string>('offline')

watch(
  () => streamStatus,
  (value) => {
    status.value = unref(value)
  },
  {
    immediate: true,
    deep: true,
  }
)

const toggle = () => {

}

</script>

<template>
  <div @click="toggle">
    <ElPopover
      placement="bottom"
      :title="status"
      :width="200"
      trigger="hover"
      content="Server connection status">
      <template #reference>
        <span :class="prefixCls" :color="color">
          <svg v-if="status == 'online'" xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24">
          <path fill="currentColor" d="M21.02 13.01v-2h-4.01v-5h1v-4h-4v4h1v5h-12v2H7V18H6v4h3.99v-4h-1v-4.99z"/>
        </svg>
        <svg v-else xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24">
          <path fill="currentColor"
                d="m2.39 4.93l6.08 6.08H3.01v2H7V18H6v4h3.99v-4h-1v-4.99h1.48l8.65 8.65l1.28-1.28L3.67 3.66zm14.62 6.08v-5h1v-4h-4v4h1v5h-1.44l2 2h5.44v-2z"/>
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
