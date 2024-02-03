<script setup lang="ts">
import {computed, onMounted, PropType} from 'vue'
import {ItemPayloadGrid, GridProp} from "@/views/Dashboard/card_items/grid/types";
import {ApiImage} from "@/api/stub";
import {prepareUrl} from "@/utils/serverId";
import { ElTooltip } from 'element-plus'

const props = defineProps({
  tileItem: {
    type: Object as PropType<Nullable<GridProp>>,
    default: () => null
  },
  baseParams: {
    type: Object as PropType<Nullable<ItemPayloadGrid>>,
  },
  templates: {
    type: Object as PropType<Map<string, GridProp>>,
  },
  cell: null
})

onMounted(() => {

})

const cellHeight = computed(() => props.baseParams?.cellHeight + 'px')
const cellWidth = computed(() => props.baseParams?.cellWidth + 'px')

const getUrl = (image: ApiImage): string => {
  if (!image || !image?.url) {
    return '';
  }
  return prepareUrl(import.meta.env.VITE_API_BASEPATH as string + image?.url);
}

const getImage = ({position, top, left, image}) => {
  const uri = getUrl(image);
  image = `url(${uri})`
  if (position) {
    image = `url(${uri}) ${left}px ${top}px no-repeat`
  }
  return image;
}

const getValue = computed(() => {
  if (typeof props.cell === 'object') {
    if (props.cell.hasOwnProperty('v')) {
      return props.cell.v
    }
  }
  return props.cell;
})

const getTooltip = computed(() => {
  if (typeof props.cell === 'object') {
    if (props.cell.hasOwnProperty('t')) {
      return props.cell.t
    }
  }
  return getValue.value +'';
})

const cellTemplate = computed(() => props.templates? props.templates[getValue.value] : props.tileItem || null)

const getTileStyle = () => {
  let style = {}
  if (cellTemplate.value && cellTemplate.value?.height) {
    style["height"] = cellTemplate.value.height + 'px';
  }
  if (cellTemplate.value && cellTemplate.value?.width) {
    style["width"] = cellTemplate.value.width + 'px';
  }
  if (cellTemplate.value && cellTemplate.value?.image) {
    const background = getImage(cellTemplate.value);
    if (background) {
      // console.log(set background', background)
      style["background"] = background
      const {position} = cellTemplate.value
      if (!position)  {
        style["background-size"] = "cover";
      }
    }
  }
  if (!style["background"] && props.baseParams?.image) {
    const background = getImage(props.baseParams)
    if (background) {
      // console.log(override background', background)
      style["background"] = background
      const {position} = props.baseParams;
      if (!position) {
        style["background-size"] = "cover";
      }
    }
  }
  if (typeof props.cell === 'object') {
    if (props.cell.hasOwnProperty('b')) {
      style["background-color"] = props.cell.b
    }
    if (props.cell.hasOwnProperty('c')) {
      style["color"] = props.cell.c
    }
  }
  if (props.baseParams?.fontSize) {
    style["font-size"] = props.baseParams?.fontSize +'px'
  }

  return style
}

</script>

<template>
  <div class="tile-wrapper">
    <ElTooltip :disabled="!baseParams.tooltip" :content="getTooltip" placement="auto">
      <div :class="[{'positioned': cellTemplate?.height || cellTemplate?.width, 'tile-inner': true}]" :style="getTileStyle()">
        <span v-if="baseParams?.showCellValue" v-html="getValue"></span>
      </div>
    </ElTooltip>
  </div>
</template>

<style lang="less" scoped >

.tile-wrapper {
  height: v-bind(cellHeight);
  width: v-bind(cellWidth);
  //cursor: pointer;
  overflow: hidden;
  position: relative;
}

.tile-wrapper .tile-inner {
  margin: 0 auto;
  height: inherit;
  width: inherit;
  display: flex;
  justify-content: center;
  align-items: center;
  &.positioned {
    margin: 0;
    position: absolute;
    top: 50%;
    left: 50%;
    -ms-transform: translate(-50%, -50%);
    transform: translate(-50%, -50%);
  }
}

</style>
