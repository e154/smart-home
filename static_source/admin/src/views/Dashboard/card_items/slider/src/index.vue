<script setup lang="ts">
import {computed, onMounted, PropType, ref, unref, watch} from "vue";
import {Cache, CardItem, eventBus, RenderVar, requestCurrentState} from "@/views/Dashboard/core";
import {debounce} from "lodash-es";
import {ItemPayloadSlider, OrientationType} from "./types";
import slider from "vue3-slider"
import api from "@/api/api";
import {useI18n} from "@/hooks/web/useI18n";
import {ApiEntityCallActionRequest, ApiTypes} from "@/api/stub";
import RoundSlider from 'vue-three-round-slider'

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

const el = ref<ElRef>(null)
const radius = ref(120)

onMounted(() => {

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
const tooltip = ref(false)
const orientation = ref<OrientationType>(OrientationType.horizontal)

const currentSlider = computed<ItemPayloadSlider>(() => props.item?.payload.slider || {} as ItemPayloadSlider)

const _cache = new Cache()
const getValue = debounce(async () => {

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

  let token: string = props.item?.payload.slider?.attribute || ''
  if (token) {
    const result = await RenderVar(token, props.item?.lastEvent)
    if (result !== '[NO VALUE]') {
      const val = parseInt(result) || 0
      if (unref(value) !== val) {
        value.value = val
      }
    }
  }
})

watch(
    () => props.item,
    (val?: CardItem) => {
      radius.value = el.value?.clientWidth / 2
      if (!val) return;
      getValue()
    },
    {
      deep: true,
      immediate: true
    }
)

const callAction = debounce(async (val: number) => {
  if (currentSlider.value?.eventName) {
    eventBus.emit(currentSlider.value?.eventName, currentSlider.value?.eventArgs)
  }
  if (!currentSlider.value.action) {
    console.warn('no action')
    return;
  }
  await api.v1.interactServiceEntityCallAction({
    id: props.item?.entityId,
    name: currentSlider.value.action,
    tags: currentSlider.value.tags,
    areaId: currentSlider.value.areaId,
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

const roundSliderHandler = (val: {value: number}) => {
  callAction(val.value)
}

requestCurrentState(props.item?.entityId);

</script>

<template>
  <div ref="el" class="h-[100%] w-[100%]">
    <slider v-if="['vertical', 'horizontal'].includes(orientation)"
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

    <RoundSlider
        v-if="orientation === 'circular'"
        v-model="value"
        start-angle="315"
        end-angle="+270"
        line-cap="round"
        :radius="radius"
        :min="min"
        :max="max"
        :step="step"
        :width="height"
        :pathColor="trackColor"
        :rangeColor="color"
        :disabled="!item?.enabled"
        :valueChange="roundSliderHandler"
        :showTooltip="false"
    />

    <div class="universal" v-if="['verticalV2', 'universal'].includes(orientation)">
      <div class="wrapper">
        <input type="range" min="0" max="100" :value="value" :step="step" @input="dragEnd(parseInt($event.target.value))" />
      </div>
    </div>
  </div>
</template>

<style lang="less">

.universal {
  margin: 0 auto;
  display: grid;
  height: 100%;
  width: 100%;

  .wrapper {
    position: relative;
    width: 100%;
    height: 100%;
  }
  .wrapper::before, .wrapper::after {
    position: absolute;
    z-index: 99;
    color: v-bind(color);
    line-height: 1;
    pointer-events: none;
  }
  .wrapper::before {
    content: "+";
    color: #eee;
    top: 50%;
    right: 6%;
    transform: translate(-50%, -50%);
  }
  .wrapper::after {
    content: "−";
    color: #eee;
    top: 50%;
    left: 10%;
    transform: translate(-50%, -50%);
  }

  input[type=range] {
    -webkit-appearance: none;
    background-color: v-bind(trackColor);
    position: absolute;
    margin: 0;
    padding: 0;
    width: 100%;
    height: 100%;
    border-radius: 1rem;
    overflow: hidden;
    cursor: row-resize;
  }

  input[type=range][step] {
    background-color: v-bind(trackColor);
  }
  input[type='range']::-webkit-slider-thumb {
    width: 0;
    -webkit-appearance: none;
    cursor: ew-resize;
    box-shadow: -20rem 0 0 20rem v-bind(color);
    background: v-bind(color);
  }
  input[type=range]::-moz-range-thumb {
    border: none;
    width: 0;
    cursor: ew-resize;
    box-shadow: -20rem 0 0 20rem v-bind(color);
    background: v-bind(color);
  }

}
</style>
