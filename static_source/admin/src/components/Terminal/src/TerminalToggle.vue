<script setup lang="ts">
import {useAppStore} from "@/store/modules/app";
import { useDesign } from '@/hooks/web/useDesign'
import {propTypes} from "@/utils/propTypes";

const appStore = useAppStore()

const { getPrefixCls } = useDesign()

defineProps({
  color: propTypes.string.def('')
})

const toggleFullscreen = () => {
  appStore.setTerminal(!appStore.getTerminal)
}

</script>

<template>
  <div @click="toggleFullscreen">
    <Icon
      :size="18"
      icon="bi:terminal"
      :color="color"
    />
  </div>
</template>
