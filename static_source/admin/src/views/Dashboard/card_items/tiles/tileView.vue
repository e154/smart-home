<script setup lang="ts">
import {computed, PropType} from 'vue'
import {ItemPayloadTiles, TileProp} from "@/views/Dashboard/card_items/tiles/types";
import {ApiImage} from "@/api/stub";
import {prepareUrl} from "@/utils/serverId";

const props = defineProps({
  tileItem: {
    type: Object as PropType<Nullable<TileProp>>,
    default: () => null
  },
  baseParams: {
    type: Object as PropType<Nullable<ItemPayloadTiles>>,
  },
})

const tileHeight = computed(() => props.baseParams?.tileHeight + 'px')
const tileWidth = computed(() => props.baseParams?.tileWidth + 'px')

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

const getTileStyle = () => {
  let style = {}
  if (props.tileItem?.height) {
    style["height"] = props.tileItem.height + 'px';
  }
  if (props.tileItem?.width) {
    style["width"] = props.tileItem.width + 'px';
  }
  if (props.tileItem && props.tileItem?.image) {
    const background = getImage(props.tileItem);
    if (background) {
      style["background"] = background
      const {position} = props.tileItem;
      if (!position)  {
        style["background-size"] = "cover";
      }
    }
  }
  if (!style["background"] && props.baseParams?.image) {
    const background = getImage(props.baseParams)
    if (background) {
      style["background"] = background
      const {position} = props.baseParams;
      if (!position) {
        style["background-size"] = "cover";
      }
    }
  }


  return style
}
</script>

<template>
  <div class="tile-wrapper">
    <div :class="[{'positioned': tileItem?.height || tileItem?.width, 'tile-inner': true}]" :style="getTileStyle()"></div>
  </div>
</template>

<style lang="less" scoped >

.tile-wrapper {
  height: v-bind(tileHeight);
  width: v-bind(tileWidth);
  //cursor: pointer;
  overflow: hidden;
  position: relative;
}

.tile-wrapper .tile-inner {
  margin: 0 auto;
  height: inherit;
  width: inherit;
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
