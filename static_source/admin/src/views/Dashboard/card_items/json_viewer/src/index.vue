<script setup lang="ts">
import {onMounted, onUnmounted, PropType, ref, watch} from "vue";
import {CardItem} from "@/views/Dashboard/core/core";
import {RenderVar} from "@/views/Dashboard/core/render";
import {debounce} from "lodash-es";
import {JsonViewer} from "@/components/JsonViewer";

// ---------------------------------
// common
// ---------------------------------
const props = defineProps({
  item: {
    type: Object as PropType<Nullable<CardItem>>,
    default: () => null
  },
})

const el = ref<ElRef>(null)
onMounted(() => {
  // store dom element moveable
  props.item.setTarget(el.value)
})

onUnmounted(() => {

})

// ---------------------------------
// component methods
// ---------------------------------
const currentValue = ref<Nullable<any>>(null)

const update = debounce(async () => {
  if (props.item?.payload.jsonViewer.attrField) {
    let token: string = props.item?.payload.jsonViewer?.attrField || ''
    const value = RenderVar(token, props.item?.lastEvent)
    if (typeof value === 'string') {
      try {
        currentValue.value = JSON.parse(value)
      } catch (e) {
        currentValue.value = value
      }
      return
    }
    currentValue.value = value
  }
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

// ---------------------------------
// run
// ---------------------------------

update();

</script>

<template>
  <div ref="el" :class="[{'hidden': item.hidden}]" class="overflow-auto">
    <JsonViewer v-model="currentValue"/>
  </div>
</template>

<style lang="less" scoped>

.hidden {
  z-index: -99999;
}


</style>
