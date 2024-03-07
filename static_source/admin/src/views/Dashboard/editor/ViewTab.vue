<script setup lang="ts">
import {computed, onMounted, onUnmounted, PropType, ref, watch} from "vue";
import {Card, Core, eventBus, Tab} from "@/views/Dashboard/core";
import Vuuri from "@/components/Vuuri"
import debounce from 'lodash.debounce'
import ViewCard from "@/views/Dashboard/editor/ViewCard.vue";
import {Frame} from "@/views/Dashboard/components";
import {loadFonts} from "@/utils/fonts";
import {useAppStore} from "@/store/modules/app";
import {ElButton, ElButtonGroup, ElIcon, ElMenu, ElMenuItem, ElTag} from "element-plus";
import {DraggableContainer} from "@/components/DraggableContainer";
import {CloseBold} from "@element-plus/icons-vue";

const appStore = useAppStore()

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
  // reloadKey.value += 1
  grid.value.update();
}, 100)

const eventHandler = (event: string, args: any[]) => {
  if (props.tab?.id === args) {
    // console.log('update tab', tabId)
    reload()
  }
}

onMounted(() => {
  eventBus.subscribe(['updateTab'], eventHandler)
})

onUnmounted(() => {
  eventBus.unsubscribe(['updateTab'], eventHandler)
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
const modalCards = computed<Card[]>(() => props.tab?.modalCards)

watch(
  () => props.tab.fonts,
  (val?: string[]) => {
    if (!val) return
    val.forEach(variableName => loadFonts(variableName))
  },
  {
    immediate: true
  }
)

const getBackground = (card: Card) => {
  let background = 'inherit'
  if (card?.background) {
    background = card.background
  } else {
    if (card?.backgroundAdaptive) {
      background = appStore.isDark ? '#232324' : '#F5F7FA'
    }
  }
  return background
}

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
      <Frame :frame="item.templateFrame" :background="getBackground(item)" v-if="item.template">
        <ViewCard :card="item" :key="item" :core="core"/>
      </Frame>
      <ViewCard v-else :card="item" :key="item" :core="core"/>
    </template>
  </Vuuri>

  <DraggableContainer
      v-for="(item, index) in modalCards"
      :key="index"
      :name="'modal-card-items'"
      :initial-width="item.width" :min-width="item.width"
      :initial-height="item.height" :min-height="item.height"
  >
    <template #header>
      <span v-html="item.title"></span>
    </template>
    <template #default>
      <Frame :frame="item.templateFrame" :background="getBackground(item)" v-if="item.template">
        <ViewCard :card="item" :key="index" :core="core"/>
      </Frame>
      <ViewCard v-else :card="item" :key="index" :core="core"/>
    </template>
  </DraggableContainer>

</template>

<style lang="less">
.gap {
  .muuri-item {
    padding: 5px;

    .muuri-item-content {
    //border: 1px solid #e9edf3;
    }
  }
}

.draggable-container.container-modal-card-items {
  .draggable-container-content {
    padding: 0;
  }
}
</style>
