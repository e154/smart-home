<script setup lang="ts">
import {useAppStore} from "@/store/modules/app";
import {propTypes} from "@/utils/propTypes";

const appStore = useAppStore()

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
