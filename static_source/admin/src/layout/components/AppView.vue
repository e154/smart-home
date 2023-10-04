<script setup lang="ts">
import {useTagsViewStore} from '@/store/modules/tagsView'
import {useAppStore} from '@/store/modules/app'
import {Footer} from '@/components/Footer'
import {computed, onMounted, onUnmounted} from 'vue'
import {Terminal} from "@/components/Terminal";

const appStore = useAppStore()

const layout = computed(() => appStore.getLayout)

const fixedHeader = computed(() => appStore.getFixedHeader)

const footer = computed(() => appStore.getFooter)

const tagsViewStore = useTagsViewStore()

const getCaches = computed((): string[] => {
  return tagsViewStore.getCachedViews
})

const onKeydown = ( event ) => {
  // console.log(event.key);
  if (event.key === "Escape") {
    appStore.setTerminal(false)
  }
  if (event.key === "`") {
    appStore.setTerminal(!appStore.getTerminal)
  }
}

onMounted(() => {
  document.addEventListener("keydown", onKeydown)
})

onUnmounted(() => {
  document.removeEventListener("keydown", onKeydown)
})

</script>

<template>
  <section
    :class="[
      'p-[var(--app-content-padding)] w-[100%] bg-[var(--app-content-bg-color)] dark:bg-[var(--el-bg-color)]', '!min-h-[calc(100%-var(--top-tool-height))]', 'h-[100%]',
      {
        '!min-h-[calc(100%-var(--app-footer-height))]':
          ((fixedHeader && (layout === 'classic' || layout === 'topLeft')) || layout === 'top') &&
          footer,

        '!min-h-[calc(100%-var(--tags-view-height)-var(--top-tool-height)-var(--app-footer-height))]':
          !fixedHeader && layout === 'classic' && footer,

        '!min-h-[calc(100%-var(--tags-view-height)-var(--app-footer-height))]':
          !fixedHeader && (layout === 'topLeft' || layout === 'top') && footer,

        '!min-h-[calc(100%-var(--top-tool-height))]': fixedHeader && layout === 'cutMenu' && footer,

        '!min-h-[calc(100%-var(--top-tool-height)-var(--tags-view-height))]':
          !fixedHeader && layout === 'cutMenu' && footer
      }
    ]"
  >
    <router-view>
      <template #default="{ Component, route }">
        <keep-alive :include="getCaches">
          <component :is="Component" :key="route.fullPath" />
        </keep-alive>
      </template>
    </router-view>
  </section>
  <Footer v-if="footer" />
  <Terminal/>
</template>
