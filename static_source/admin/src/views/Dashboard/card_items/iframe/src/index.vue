<script setup lang="ts">
import {onMounted, PropType, ref, watch} from "vue";
import {CardItem, RenderVar, requestCurrentState} from "@/views/Dashboard/core";
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

const el = ref<ElRef>(null)
onMounted(() => {

})

// ---------------------------------
// component methods
// ---------------------------------
const src = ref<string>(props.item?.payload?.iframe?.uri || '');

const update = debounce(async (item?: CardItem) => {
  let value = item?.payload?.iframe?.uri || '';
  if (item?.payload?.iframe?.attrField) {
    let token: string = item?.payload.iframe?.attrField || ''
    value = await RenderVar(token, item?.lastEvent)
  }
  if (src.value == value) {
    return
  }
  src.value = value
}, 100)

watch(
  () => props.item,
  (val?: CardItem) => {
    if (!val) return;
    update(val)
  },
  {
    deep: true,
    immediate: true
  }
)

// ---------------------------------
// run
// ---------------------------------

requestCurrentState(props.item.entityId!);

update()

</script>

<template>
  <div ref="el" :class="[{'hidden': item.hidden}]">
    <iframe
      :src="src"
      style="position: absolute; border: none; width: 100%; height: inherit"
      width="100%"
      height="100%"
      frameborder="0"
    >
    </iframe>
  </div>
</template>

<style lang="less">

</style>
