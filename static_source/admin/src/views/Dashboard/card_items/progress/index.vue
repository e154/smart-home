<script setup lang="ts">
import {computed, onMounted, PropType, ref, unref, watch} from "vue";
import {Card, CardItem, Core, requestCurrentState, Tab} from "@/views/Dashboard/core";
import {debounce} from "lodash-es";
import {Cache, GetTokens, RenderText} from "@/views/Dashboard/render";
import {ElProgress} from "element-plus";

// ---------------------------------
// common
// ---------------------------------

const value = ref(0)

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

// ---------------------------------
// component methods
// ---------------------------------

const _cache = new Cache()
const update = debounce(() => {
  let v: string = props.item?.payload.progress?.value || ''
  const tokens = GetTokens(props.item?.payload.progress?.value, _cache)
  if (tokens.length) {
    v = RenderText(tokens, v, props.item?.lastEvent)
  }
  value.value = parseInt(v) || 0
})

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

requestCurrentState(props.item?.entityId);

update()

</script>

<template>
  <div ref="el" v-if="item.entity" class="h-[100%] w-[100%]">
   <ElProgress
       v-if="item.payload.progress.type"
       :type="item.payload.progress.type"
       :percentage="value"
       :width="item.payload.progress.width"
       :stroke-width="item.payload.progress.strokeWidth"
       :text-inside="!item.payload.progress.textInside"/>
   <ElProgress
       v-else
       :percentage="value"
       :width="item.payload.progress.width"
       :stroke-width="item.payload.progress.strokeWidth"
       :text-inside="!item.payload.progress.textInside"/>
 </div>
</template>

<style lang="less" >

</style>
