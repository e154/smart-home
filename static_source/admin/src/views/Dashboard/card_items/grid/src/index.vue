<script setup lang="ts">
import {computed, onMounted, PropType, ref, unref, watch} from "vue";
import {CardItem, requestCurrentState, RenderVar, Cache} from "@/views/Dashboard/core";
import {GridProp} from "./types";
import {debounce} from "lodash-es";
import api from "@/api/api";
import {ElMessage} from "element-plus";
import {useI18n} from "@/hooks/web/useI18n";
import {ApiTypes} from "@/api/stub";
import CellView from "./CellView.vue";

const {t} = useI18n()

// ---------------------------------
// common
// ---------------------------------

const props = defineProps({
  item: {
    type: Object as PropType<Nullable<CardItem>>,
    default: () => null
  },
})

const currentItem = computed(() => props.item as CardItem);

const el = ref<ElRef>(null)
onMounted(() => {

})

// ---------------------------------
// component methods
// ---------------------------------
const board = ref([])

const cellHeight = computed(() => props.item.payload.grid.cellHeight + 'px');
const cellWidth = computed(() => props.item.payload.grid.cellWidth + 'px');
const gapSize = computed(() => (props.item.payload.grid.gap ? props.item.payload.grid.gapSize : 0) + 'px');

const getBoard = (str: string): any[] => {
  try {
    return JSON.parse(str);
  } catch (e) {
    return [];
  }
}

const _cache = new Cache()
const update = debounce( async () => {
  let token: string = props.item?.payload.grid?.attribute || ''
  const result = await RenderVar(token, props.item?.lastEvent)
  board.value = getBoard(result) || []
}, 100)

const tileTemplates = ref<Map<string, GridProp>>({});
const prepareTileTemplates = debounce(() => {
  tileTemplates.value = {};
  if (!props.item?.payload?.grid?.items) {
    return
  }
  for (const item of props.item?.payload.grid?.items) {
    tileTemplates.value[item.key] = item;
  }
}, 100)

watch(
    () => props.item?.uuid,
    (val?: string) => {
      if (!val) return;
      update()
      prepareTileTemplates()
    },
)

const callAction = async (row: number, cell: number) => {
  if (!currentItem.value.payload.grid?.tileClick || !currentItem.value.payload.grid?.actionName) {
    return;
  }
  await api.v1.interactServiceEntityCallAction({
    id: currentItem.value.payload.grid.entityId,
    name: currentItem.value.payload.grid?.actionName,
    tags: currentItem.value.payload.grid?.tags || [],
    areaId: currentItem.value.payload.grid?.areaId,
    attributes: {
      "row": {
        "name": "row",
        "type": ApiTypes.INT,
        "int": row,
      },
      "cell": {
        "name": "cell",
        "type": ApiTypes.INT,
        "int": cell,
      },
    },
  });
  ElMessage({
    title: t('Success'),
    message: t('message.callSuccessful'),
    type: 'success',
    duration: 2000
  });
}

requestCurrentState(props.item?.entityId);

update()
prepareTileTemplates();

</script>

<template>
  <div ref="el">
    <div class="grid-container">
      <div
          class="grid-row"
          :key="rowIndex"
          v-for="(row, rowIndex) in board">
        <div
            class="grid-cell"
            :key="cellIndex"
            v-for="(cell, cellIndex) in row"
            @click.prevent.stop="callAction(rowIndex, cellIndex)">
          <CellView
              :key="cellIndex"
              :base-params="props.item.payload.grid"
              :tile-item="tileTemplates[cell]"
              :templates="tileTemplates"
              :cell="cell"/>
        </div>
      </div>
    </div>

  </div>
</template>

<style lang="less" scoped>
.grid-container {
}

.grid-row {
  clear: both;
}

.grid-cell {
  float: left;
  height: v-bind(cellHeight) !important;
  width: v-bind(cellWidth) !important;
  margin: v-bind(gapSize) !important;
}

</style>
