<script setup lang="ts">
import {computed, PropType, ref, unref, watch} from "vue";
import {Card, Core, Tab} from "@/views/Dashboard/core";
import {Vuuri} from "@/views/Dashboard/Vuuri"
import {useBus} from "@/views/Dashboard/bus";
import debounce from 'lodash.debounce'
import ViewCard from "@/views/Dashboard/editor/ViewCard.vue";

// ---------------------------------
// common
// ---------------------------------

const grid = ref(null)
const props = defineProps({
  core: {
    type: Object as PropType<Core>,
  },
  tab: {
    type: Object as PropType<Tab>,
    default: () => null
  },
})

const reloadKey = ref(0);
const reload = debounce(() => {
  // console.log('reload tab')
  reloadKey.value += 1
}, 100)

watch(
    () => props.tab,
    (val?: Tab) => {
      if (!val) return;
      reload()
    },
    {
      deep: false,
      immediate: true
    }
)

useBus({
  name: 'update_tab',
  callback: (tabId: number) => {
    console.log('update_tab', tabId, props.tab?.id)
    if (props.tab?.id === tabId) {
      reload()
    }
  }
})

// ---------------------------------
// component methods
// ---------------------------------

const getItemWidth = (card: Card) => {
  // console.log('getItemWidth', activeTab.columnWidth)
  if (card.width > 0) {
    return `${card.width}px`
  }
  return `${props.tab?.columnWidth}px`
}

const getItemHeight = (card: Card) => {
  // console.log('getItemHeight', card.height)
  return `${card.height}px`
}

const cards = computed<Card[]>(() => props.tab?.cards2)

</script>

<template>
  <Vuuri
      v-model="cards"
      item-key="id"
      :get-item-width="getItemWidth"
      :get-item-height="getItemHeight"
      :drag-enabled="false"
      ref="grid"
      :key="reloadKey"
  >
      <template #item="{item}">
        <ViewCard :card="item" :key="item" :core="core"/>
      </template>
  </Vuuri>

</template>

<style lang="less" >
.gap {
.muuri-item {
  padding: 5px;
.muuri-item-content {
  //border: 1px solid #e9edf3;
}
}
}
</style>
