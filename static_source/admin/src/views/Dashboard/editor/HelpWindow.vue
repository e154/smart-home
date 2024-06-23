<script setup lang="ts">
import {computed, onMounted, onUnmounted, PropType, ref} from 'vue'
import {ElCollapse, ElCollapseItem, ElDescriptions, ElDescriptionsItem, ElIcon,} from 'element-plus'
import {useI18n} from '@/hooks/web/useI18n'
import {Core, eventBus, Filters, Tab} from "@/views/Dashboard/core";
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
  eventBus.subscribe('toggleHelpMenu', eventHandler)
})

onUnmounted(() => {
  eventBus.unsubscribe('toggleHelpMenu', eventHandler)
})

const filters = computed(() => Filters)

</script>

<template>
  <DraggableContainer :name="'help'" :initial-width="500" :min-width="500" v-show="showMenuWindow">
    <template #header>
      <div class="w-[100%]">
        <div style="float: left">Help</div>
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

      <ElCollapse class="ml-10px">
        <ElCollapseItem
            v-for="(filter, index) in filters"
            :name="index"
            :key="index"
        >

          <template #title>
            {{ filter.name }}
          </template>

          <ElDescriptions v-if="filter"
                          class="ml-10px mr-10px mb-20px"
                          direction="vertical"
                          :column="3"
                          border
          >
            <ElDescriptionsItem :label="$t('dashboard.editor.filter.name')">{{ filter.name }}
            </ElDescriptionsItem>
            <ElDescriptionsItem :label="$t('dashboard.editor.filter.example')">{{ filter.example }}
            </ElDescriptionsItem>
            <ElDescriptionsItem :label="$t('dashboard.editor.filter.args')">{{ filter.args || '-' }}
            </ElDescriptionsItem>
            <ElDescriptionsItem :label="$t('dashboard.editor.filter.description')">
              {{ filter.description }}
            </ElDescriptionsItem>
          </ElDescriptions>
        </ElCollapseItem>
      </ElCollapse>
    </template>
  </DraggableContainer>
</template>

<style lang="less">

</style>
