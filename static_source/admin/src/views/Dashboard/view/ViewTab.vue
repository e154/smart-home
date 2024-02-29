<script setup lang="ts">
import {computed, PropType} from "vue";
import {Card, Core, Tab} from "@/views/Dashboard/core/core";
import {Vuuri} from "@/views/Dashboard/Vuuri"
import ViewCard from "@/views/Dashboard/view/ViewCard.vue";
import {Frame} from "@/views/Dashboard/components";

// ---------------------------------
// common
// ---------------------------------

const props = defineProps({
  core: {
    type: Object as PropType<Core>,
  },
  tab: {
    type: Object as PropType<Tab>,
    default: () => null
  },
})

// ---------------------------------
// component methods
// ---------------------------------

const getItemWidth = (card: Card) => {
  if (card.width > 0) {
    return `${card.width}px`
  }
  return `${props.tab?.columnWidth}px`
}

const getItemHeight = (card: Card) => {
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
  >
    <template #item="{item}">
      <Frame :frame="item.templateFrame" v-if="item.template">
        <ViewCard :card="item" :key="item" :core="core"/>
      </Frame>
      <ViewCard v-else :card="item" :key="item" :core="core"/>
    </template>
  </Vuuri>

</template>

<style lang="less">
.gap {
  .muuri-item {
    padding: 5px;

    .muuri-item-content {
      overflow: -webkit-paged-x;
    }
  }
}
</style>
