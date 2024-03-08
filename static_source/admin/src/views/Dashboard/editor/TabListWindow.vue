<script setup lang="ts">
import {computed, onMounted, onUnmounted, PropType, ref} from 'vue'
import {ElIcon, ElMenu, ElMenuItem} from 'element-plus'
import {useI18n} from '@/hooks/web/useI18n'
import {Core, Tab, eventBus} from "@/views/Dashboard/core";
import {DraggableContainer} from "@/components/DraggableContainer";
import {CloseBold} from "@element-plus/icons-vue";

const {t} = useI18n()

const props = defineProps({
  core: {
    type: Object as PropType<Nullable<Core>>,
    default: () => null
  },
})

const currentCore = computed(() => props.core as Core)

// ---------------------------------
// common
// ---------------------------------

const menuTabClick = (index: number, tab: Tab) => {
  currentCore.value.selectTabInMenu(index)
}

const eventHandler = (event: string, args: any[]) => {
  showMenuWindow.value = !showMenuWindow.value
}

const showMenuWindow = ref(false)
onMounted(() => {
  eventBus.subscribe('toggleTabsMenu', eventHandler)
})

onUnmounted(() => {
  eventBus.unsubscribe('toggleTabsMenu', eventHandler)
})

</script>

<template>
  <DraggableContainer :name="'editor-tabs'" :initial-width="280" :min-width="280" v-show="showMenuWindow">
    <template #header>
      <div class="w-[100%]">
        <div style="float: left">Tabs</div>
        <div style="float: right; text-align: right">
          <a href="#" @click.prevent.stop='showMenuWindow= false'>
            <ElIcon class="mr-5px">
              <CloseBold/>
            </ElIcon>
          </a>
        </div>
      </div>
    </template>
    <template #default>

      <ElMenu v-if="currentCore.tabs.length" :default-active="currentCore.activeTabIdx + ''"
              v-model="currentCore.activeTabIdx" class="el-menu-vertical-demo">
        <ElMenuItem :index="index + ''" :key="tab" v-for="(tab, index) in currentCore.tabs"
                    @click="menuTabClick(index, tab)">
          <div class="w-[100%] card-header">
            <span>{{ tab.name }}</span>
          </div>
        </ElMenuItem>
      </ElMenu>
    </template>
  </DraggableContainer>
</template>

<style lang="less">

</style>
