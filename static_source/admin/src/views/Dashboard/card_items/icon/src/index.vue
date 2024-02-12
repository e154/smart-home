<script setup lang="ts">
import {onMounted, onUnmounted, PropType, ref, watch} from "vue";
import {CardItem} from "@/views/Dashboard/core/core";
import {RenderVar} from "@/views/Dashboard/core/render";
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
<!--  <div ref="el" :class="[{'hidden': item.hidden}]" class="icon-item" v-html="icon"></div>-->
  <div ref="el" :class="[{'hidden': item.hidden}]" style="width: 100%; height: 100%">
    <Icon
        style="width: 100%; height: 100%"
        :key="reloadKey"
        :icon="icon"
        :color="iconColor"
        :size="iconSize"/>
  </div>

</template>

<style lang="less" scoped>
.icon-item {
  :deep(svg) {
    height: 100%;
    width: 100%;
    color: v-bind(iconColor)!important;
  }
}
</style>
