<script setup lang="ts">
import {computed, onMounted, PropType, ref, unref, watch} from "vue";
import {CardItem, requestCurrentState} from "@/views/Dashboard/core";
import {debounce} from "lodash-es";
import {Cache, RenderVar} from "@/views/Dashboard/render";
import api from "@/api/api";
import {ElColorPicker} from 'element-plus'
import {useI18n} from "@/hooks/web/useI18n";
import {ApiEntityCallActionRequest, ApiTypes} from "@/api/stub";
import {ItemPayloadColorPicker} from "./types";

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

const value = ref("")

const currentColorPicker = computed<ItemPayloadColorPicker>(() => props.item?.payload.colorPicker || {} as ItemPayloadColorPicker)

const _cache = new Cache()
const getValue = debounce(() => {
  if (!value.value && currentColorPicker.value.color != undefined) {
    value.value = currentColorPicker.value.color
  }

  let token: string = props.item?.payload.colorPicker?.attribute || ''
  if (token) {
    const result = RenderVar(token, props.item?.lastEvent)
    if (result !== '[NO VALUE]') {
      if (unref(value) !== result) {
        value.value = result
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

const callAction = debounce(async (val: string) => {
  if (!currentColorPicker.value.action) {
    console.warn('no action')
    return;
  }
  await api.v1.interactServiceEntityCallAction({
    id: props.item?.entityId,
    name: currentColorPicker.value.action,
    attributes: {
      "color": {
        "name": "color",
        "type": ApiTypes.STRING,
        "string": val,
      }
    },
  } as ApiEntityCallActionRequest);
}, 100)

const updateColor = (val: string) => {
  callAction(val)
}

requestCurrentState(props.item?.entityId);

</script>

<template>
  <div ref="el" class="h-[100%] w-[100%]">
    <ElColorPicker v-model="value" v-on:change="updateColor"/>
  </div>
</template>

<style lang="less">

</style>
