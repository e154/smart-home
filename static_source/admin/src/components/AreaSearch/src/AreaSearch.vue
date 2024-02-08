<script setup lang="ts">

import {PropType, ref, unref, watch} from "vue";
import {ApiArea} from "@/api/stub";
import {ElAutocomplete} from 'element-plus'
import api from "@/api/api";

const currentArea = ref<Nullable<ApiArea>>(null)
const areaName = ref<Nullable<string>>(null)
const emit = defineEmits(['change', 'update:modelValue'])

const props = defineProps({
  modelValue: {
    type: Object as PropType<Nullable<ApiArea>>,
    default: () => undefined
  }
})

watch(
    () => props.modelValue,
    (val?: ApiArea) => {
      if (val === unref(currentArea)) return
      areaName.value = val?.name || null;
      currentArea.value = val || null;
    },
)

// 监听
watch(
    () => currentArea.value,
    (val?: ApiArea) => {
      if (props.modelValue == unref(val)) return;
      emit('update:modelValue', val)
      if (!val) {
        emit('change', val)
      }
    },
)

const querySearchAsync = async (queryString: string, cb: Fn) => {
  if (queryString == null || queryString == 'null') {
    queryString = ''
  }
  const params = {query: queryString, limit: 25, offset: 0}
  const {data} = await api.v1.areaServiceSearchArea(params)
  const {items} = data
  cb(items)
}

const handleChange = (val) => {
  if (val == '') {
    currentArea.value = null
  }
}
const handleSelect = (val: ApiArea) => {
  currentArea.value = val
  emit('change', val)
}

</script>

<template>
  <ElAutocomplete
      class="w-[100%]"
      v-model="areaName"
      :fetch-suggestions="querySearchAsync"
      placeholder="Please input"
      value-key="name"
      @select="handleSelect"
      @change="handleChange"
  />
</template>

<style lang="less">

</style>
