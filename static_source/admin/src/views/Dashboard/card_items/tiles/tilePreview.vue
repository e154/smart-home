<script setup lang="ts">
import {computed, PropType} from 'vue'
import {ItemPayloadTiles, TileProp} from "@/views/Dashboard/card_items/tiles/types";
import TileView from "@/views/Dashboard/card_items/tiles/tileView.vue";

const props = defineProps({
  tileItem: {
    type: Object as PropType<Nullable<TileProp>>,
    default: () => null
  },
  baseParams: {
    type: Object as PropType<Nullable<ItemPayloadTiles>>,
  },
})

const wrapperHeight = computed(()=> (props.baseParams?.tileHeight + 2) + 'px')
const wrapperWidth = computed(()=> (props.baseParams?.tileWidth + 2) + 'px')

</script>

<template>
  <div class="tile-preview-wrapper">
    <TileView :base-params="baseParams" :tile-item="tileItem"/>
  </div>
</template>

<style lang="less" scoped>

.tile-preview-wrapper {
  border: 1px solid;
  height: v-bind(wrapperHeight);
  width: v-bind(wrapperWidth);
}

</style>
