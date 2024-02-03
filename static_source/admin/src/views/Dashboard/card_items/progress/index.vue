<script setup lang="ts">
import {computed, onMounted, PropType, ref, watch} from "vue";
import {CardItem, requestCurrentState} from "@/views/Dashboard/core";
import {debounce} from "lodash-es";
import {Cache, Compare, GetTokens, RenderText, RenderVar, Resolve} from "@/views/Dashboard/render";
import {ElProgress} from "element-plus";
import {Attribute, GetAttrValue} from "@/api/stream_types";

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

const optionalColor = ref('')
const color = computed(() => optionalColor.value || props.item.payload.progress.color || '')

const _cache = new Cache()
const update = debounce(() => {

  if (props.item?.payload.progress?.items) {
    for (const prop of props.item?.payload.progress?.items) {

      let token: string = props.item?.payload.progress?.value || ''
      const val = RenderVar(token, props.item?.lastEvent)

      if(!val) {
        continue
      }

      const tr = Compare(val, prop.value, prop.comparison)
      if(!tr) {
        continue
      }

      optionalColor.value = prop?.color || ''
    }
  }

  let token: string = props.item?.payload.progress?.value || ''
  const result = RenderVar(token, props.item?.lastEvent)
  value.value = parseInt(result) || 0
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
        :text-inside="!item.payload.progress.textInside"
        :show-text="item.payload.progress.showText"
        :color="color"/>
    <ElProgress
        v-else
        :percentage="value"
        :width="item.payload.progress.width"
        :stroke-width="item.payload.progress.strokeWidth"
        :text-inside="!item.payload.progress.textInside"
        :show-text="item.payload.progress.showText"
        :color="color"/>
  </div>
</template>

<style lang="less">

</style>
