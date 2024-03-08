<script setup lang="ts">
import {computed, nextTick, onMounted, onUnmounted, PropType, ref, watch} from "vue";
import {Card, Core, eventBus, Tab} from "@/views/Dashboard/core";
import Vuuri from "@/components/Vuuri"
import ViewCard from "@/views/Dashboard/editor/ViewCard.vue";
import {Frame} from "@/views/Dashboard/components";
import {loadFonts} from "@/utils/fonts";
import {useAppStore} from "@/store/modules/app";
import {DraggableContainer} from "@/components/DraggableContainer";

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
const eventHandler = (event: string, tabId: number) => {
  if (props.tab?.id === tabId) {
    // console.log('update tab', tabId)
    reloadKey.value += 1
  }
}

const eventUpdateGridHandler = (event: string, tabId: number) => {
  if (props.tab?.id === tabId) {
    // console.log('update grid', tabId)
    nextTick(() => {
      grid.value.update();
    })
  }
}

onMounted(() => {
  eventBus.subscribe('updateTab', eventHandler)
  eventBus.subscribe('updateGrid', eventUpdateGridHandler)
})

onUnmounted(() => {
  eventBus.unsubscribe('updateTab', eventHandler)
  eventBus.unsubscribe('updateGrid', eventUpdateGridHandler)
})

// ---------------------------------
// component methods
// ---------------------------------

const getItemWidth = (card: Card): string => {
  // console.log('getItemWidth', activeTab.columnWidth)
  if (card.width > 0) {
    return `${card.width}px`
  }
  return `${props.tab?.columnWidth}px`
}

const getItemHeight = (card: Card): string => {
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

const getModalWidth = (card: Card): number => {
  if (card.width > 0) {
    return card.width
  }
  return props.tab?.columnWidth
}

const getModalHeight = (card: Card) => {
  return card.height
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
      :key="index + item?.id || 0"
      :class-name="'dashboard-modal'"
      :name="'modal-card-items-' + item.id"
      :initial-width="getModalWidth(item)"
      :initial-height="getModalHeight(item) + 24"
      :modal="true"
      :resizeable="false"
      v-show="!item.hidden"
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

.draggable-container.dashboard-modal {
  .draggable-container-content {
    padding: 0;
  }
}
</style>
