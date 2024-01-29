<script setup lang="ts">
import {computed, onMounted, PropType, ref, watch} from "vue";
import {CardItem, requestCurrentState} from "@/views/Dashboard/core";
import {GridProp} from "@/views/Dashboard/card_items/grid/types";
import {Cache, GetTokens, RenderText} from "@/views/Dashboard/render";
import {debounce} from "lodash-es";
import api from "@/api/api";
import {ElMessage} from "element-plus";
import {useI18n} from "@/hooks/web/useI18n";
import {ApiTypes} from "@/api/stub";
import CellView from "@/views/Dashboard/card_items/grid/cellView.vue";

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

const el = ref(null)
onMounted(() => {
  // store dom element moveable
  props.item.setTarget(el.value)
})

// ---------------------------------
// component methods
// ---------------------------------
const board = ref([])

const cellHeight = computed(() => props.item.payload.grid.cellHeight + 'px');
const cellWidth = computed(() => props.item.payload.grid.cellWidth + 'px');
const gapSize = computed(() => (props.item.payload.grid.gap? props.item.payload.grid.gapSize : 0) + 'px'  );

const getBoard = (str: string): any[] => {
  try {
    return JSON.parse(str);
  } catch (e) {
    return [];
  }
}

const _cache = new Cache()
const update = debounce(() => {
  let v: string = props.item?.payload.grid?.attribute || ''
  const tokens = GetTokens(props.item?.payload.grid?.attribute, _cache)
  if (tokens.length) {
    v = RenderText(tokens, v, props.item?.lastEvent)
  }
  board.value = getBoard(v) || []
})

const tileTemplates = ref<Map<string, GridProp>>({});
const prepareTileTemplates = () => {
  tileTemplates.value = {};
  if (!props.item?.payload?.grid?.items) {
    return
  }
  for (const item of props.item?.payload.grid?.items) {
    tileTemplates.value[item.key] = item;
  }
}

watch(
    () => props.item,
    (val?: CardItem) => {
      if (!val) return;
      update()
      prepareTileTemplates()
    },
    {
      deep: true,
      immediate: true
    }
)

const callAction = async (index: number) => {
  if (!currentItem.value.payload.grid?.entityId ||
      !currentItem.value.payload.grid?.actionName) {
    return;
  }
  await api.v1.interactServiceEntityCallAction({
    id: currentItem.value.payload.grid.entityId,
    name: currentItem.value.payload.grid?.actionName,
    attributes: {
      "tile": {
        "name": "tile",
        "type": ApiTypes.INT,
        "int": index,
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
          :key="index"
          v-for="(row, index) in board">
        <div
            class="grid-cell"
            :key="index"
            v-for="(cell, index) in row"
            @click.prevent.stop="callAction(index)">
          <CellView
              :key="index"
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
