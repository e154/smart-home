<script setup lang="ts">
import {computed, onMounted, PropType, ref, unref, watch} from "vue";
import {Card, CardItem, Core, requestCurrentState, Tab} from "@/views/Dashboard/core";
import {debounce} from "lodash-es";
import {Cache, GetTokens, RenderText} from "@/views/Dashboard/render";
import {ItemPayloadSlider, OrientationType} from "@/views/Dashboard/card_items/slider/types";
import slider from "vue3-slider"
import api from "@/api/api";
import {ElMessage} from "element-plus";
import {useI18n} from "@/hooks/web/useI18n";
import {ApiEntityCallActionRequest, ApiTypes} from "@/api/stub";
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

const el = ref(null)
onMounted(() => {
  // store dom element moveable
  props.item.setTarget(el.value)
})

// ---------------------------------
// component methods
// ---------------------------------

const value = ref(0)
const height = ref(6)
const color = ref("#4f4f4f")
const trackColor = ref("#FEFEFE")
const min = ref(0)
const max = ref(100)
const step = ref(1)
const tooltip	 = ref(false)
const orientation = ref<OrientationType>(OrientationType.horizontal)

const currentSlider = computed<ItemPayloadSlider>(() => props.item?.payload.slider || {} as ItemPayloadSlider)

const _cache = new Cache()
const getValue = debounce(() => {

  if (currentSlider.value.height != undefined) {
    height.value = currentSlider.value.height
  }
  if (currentSlider.value.color != undefined) {
    color.value = currentSlider.value.color
  }
  if (currentSlider.value.trackColor != undefined) {
    trackColor.value = currentSlider.value.trackColor
  }
  if (currentSlider.value.min != undefined) {
    min.value = currentSlider.value.min
  }
  if (currentSlider.value.max != undefined) {
    max.value = currentSlider.value.max
  }
  if (currentSlider.value.step != undefined) {
    step.value = currentSlider.value.step
  }
  if (currentSlider.value.orientation != undefined) {
    orientation.value = currentSlider.value.orientation
  }
  if (currentSlider.value.tooltip != undefined) {
    tooltip.value = currentSlider.value.tooltip
  }

  let v: string = props.item?.payload.slider?.attribute || ''
  const tokens = GetTokens(props.item?.payload.slider?.attribute, _cache)
  if (tokens.length) {
    v = RenderText(tokens, v, props.item?.lastEvent)
    if (v !== '[NO VALUE]') {
      const val = parseInt(v) || 0
      if (unref(value) !== val) {
        value.value = val
      }
    }
  }
})

watch(
    () => props.item,
    (val?: CardItem) => {
      if (!val) return;
      getValue()
    },
    {
      deep: true,
      immediate: true
    }
)

const callAction = debounce( async (val: number) => {
  if (!currentSlider.value.action) {
    console.warn('no action')
    return;
  }
  await api.v1.interactServiceEntityCallAction({
    id: props.item?.entityId,
    name: currentSlider.value.action,
    attributes: {
      "value": {
        "name": "value",
        "type": ApiTypes.INT,
        "int": val,
      }
    },
  } as ApiEntityCallActionRequest);
}, 100)

const dragEnd = (val: number) => {
  callAction(val)
}

requestCurrentState(props.item?.entityId);

</script>

<template>
  <div ref="el" class="h-[100%] w-[100%]">
    <slider
        v-model="value"
        width="100%"
        :color="color"
        :track-color="trackColor"
        :height="height"
        :min="min"
        :max="max"
        :step="step"
        :orientation="orientation"
        :tooltip="tooltip"
        v-on:drag-end="dragEnd"
    />
 </div>
</template>

<style lang="less" >

</style>
