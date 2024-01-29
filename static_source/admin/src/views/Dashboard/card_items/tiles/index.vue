<script setup lang="ts">
import {computed, onMounted, PropType, ref, unref, watch} from "vue";
import {CardItem, requestCurrentState} from "@/views/Dashboard/core";
import {TileProp} from "@/views/Dashboard/card_items/tiles/types";
import {Cache, GetTokens, RenderText} from "@/views/Dashboard/render";
import {debounce, map} from "lodash-es";
import TileView from "@/views/Dashboard/card_items/tiles/tileView.vue";
import api from "@/api/api";
import {ElMessage} from "element-plus";
import {useI18n} from "@/hooks/web/useI18n";
import {ApiTypes} from "@/api/stub";

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

const tileHeight = computed(()=> props.item.payload.tiles.tileHeight + 'px');
const tileWidth = computed(()=> props.item.payload.tiles.tileWidth + 'px');

const getBoard = (str: string) => {
  try {
    return JSON.parse(str);
  } catch (e) {
    return [];
  }
}

const _cache = new Cache()
const update = debounce(() => {
  let v: string = props.item?.payload.tiles?.attribute || ''
  const tokens = GetTokens(props.item?.payload.tiles?.attribute, _cache)
  if (tokens.length) {
    v = RenderText(tokens, v, props.item?.lastEvent)
  }
  board.value = getBoard(v) || []
})

const tileTemplates = ref<Map<string, TileProp>>({});
const prepareTileTemplates = () => {
  tileTemplates.value = {};
  if (!props.item?.payload?.tiles?.items) {
    return
  }
  for (const item of props.item?.payload.tiles?.items) {
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
  if (!currentItem.value.payload.tiles?.entityId ||
      !currentItem.value.payload.tiles?.actionName) {
    return;
  }
  await api.v1.interactServiceEntityCallAction({
    id: currentItem.value.payload.tiles.entityId,
    name: currentItem.value.payload.tiles?.actionName,
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
          <TileView
              :key="index"
              :base-params="props.item.payload.tiles"
              :tile-item="tileTemplates[cell]"/>
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
  height: v-bind(tileHeight) !important;
  width: v-bind(tileWidth) !important;
}
</style>
