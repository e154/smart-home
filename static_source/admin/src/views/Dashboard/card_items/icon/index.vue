<script setup lang="ts">
import {computed, onMounted, onUnmounted, PropType, ref, unref, watch} from "vue";
import {ButtonAction, Card, CardItem, Core, requestCurrentState, Tab} from "@/views/Dashboard/core";
import {Cache, Compare, GetTokens, RenderText, RenderVar, Resolve} from "@/views/Dashboard/render";
import {ElImage, ElIcon} from "element-plus";
import { Picture as IconPicture } from '@element-plus/icons-vue'
import {Attribute, GetAttrValue} from "@/api/stream_types";
import {debounce} from "lodash-es";

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
const reloadKey = ref(0)
const icon = ref<Nullable<string>>(null)
const iconSize = ref<Nullable<number>>(null)
const iconColor = ref<Nullable<string>>(null)

const reload = () => {
  reloadKey.value += 1
}

const update = debounce(() => {
  if (!props.item?.payload?.icon) {
    return;
  }

  iconColor.value = props.item?.payload?.icon?.iconColor || '#000000';
  iconSize.value = props.item?.payload?.icon?.iconSize || 14;

  if (props.item?.payload.icon.attrField) {
    let token: string = props.item?.payload.icon?.attrField || ''
    icon.value = RenderVar(token, props.item?.lastEvent)
    return
  }
  icon.value = props.item?.payload.icon?.value || '';
  return;
}, 100)

watch(
    () => props.item,
    (val?: CardItem) => {
      if (!val) return;
      update()
    },
    {
      deep: true,
      immediate: true
    }
)


update();

</script>

<template>
  <div ref="el" :class="[{'hidden': item.hidden}]" style="width: 100%; height: 100%">
    <Icon
        style="width: 100%; height: 100%"
        :key="reloadKey"
        :icon="icon"
        :color="iconColor"
        :size="iconSize" />
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
