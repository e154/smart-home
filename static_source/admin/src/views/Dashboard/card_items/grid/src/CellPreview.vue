<script setup lang="ts">
import {computed, PropType} from 'vue'
import {GridProp, ItemPayloadGrid} from "./types";
import CellView from "./CellView.vue";

const props = defineProps({
  tileItem: {
    type: Object as PropType<Nullable<GridProp>>,
    default: () => null
  },
  baseParams: {
    type: Object as PropType<Nullable<ItemPayloadGrid>>,
  },
})

const wrapperHeight = computed(() => (props.baseParams?.cellHeight + 2) + 'px')
const wrapperWidth = computed(() => (props.baseParams?.cellWidth + 2) + 'px')

</script>

<template>
  <div class="tile-preview-wrapper">
    <CellView :base-params="baseParams" :tile-item="tileItem"/>
  </div>
</template>

<style lang="less" scoped>

.tile-preview-wrapper {
  border: 1px solid;
  height: v-bind(wrapperHeight);
  width: v-bind(wrapperWidth);
}

</style>
