<script setup lang="ts">
import {computed, onMounted, onUnmounted, PropType, ref, unref, watch} from "vue";
import {ButtonAction, Card, CardItem, Core, requestCurrentState, Tab} from "@/views/Dashboard/core";
import {Cache, Compare, GetTokens, RenderText, Resolve} from "@/views/Dashboard/render";
import {ElImage, ElIcon} from "element-plus";
import { Picture as IconPicture } from '@element-plus/icons-vue'
import {Attribute, GetAttrValue} from "@/api/stream_types";

// ---------------------------------
// common
// ---------------------------------
const props = defineProps({
  item: {
    type: Object as PropType<Nullable<CardItem>>,
    default: () => null
  },
})

const el = ref(null)
onMounted(() => {
  // store dom element moveable
  props.item.setTarget(el.value)
})

onUnmounted(() => {

})

// ---------------------------------
// component methods
// ---------------------------------

const getUrl = (): string => {
  if (!props.item?.payload?.image) {
    return '';
  }
  if (props.item?.payload.image.attrField) {
    const imageUrl = RenderText([props.item?.payload.image.attrField], '[' + props.item?.payload.image.attrField + ']', props.item?.lastEvent);
    return import.meta.env.VITE_API_BASEPATH as string + imageUrl;
  }
  return import.meta.env.VITE_API_BASEPATH as string + props.item?.payload.image?.image?.url || '';
}

</script>

<template>
  <div ref="el" class="w-[100%] h-[100%] overflow-hidden">
    <ElImage v-show="!item.hidden" :src="getUrl()">
      <template #error>
        <div class="image-slot">
          <ElIcon><icon-picture /></ElIcon>
        </div>
      </template>
    </ElImage>
  </div>

</template>

<style lang="less" >
.el-image__error, .el-image__placeholder, .el-image__inner {
  height: auto;
}

.el-image.item-element {
  overflow: visible;
}
</style>
