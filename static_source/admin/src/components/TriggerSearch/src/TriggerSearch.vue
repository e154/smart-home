<script setup lang="ts">

import {PropType, ref, unref, watch} from "vue";
import {ApiTrigger} from "@/api/stub";
import {ElAutocomplete} from 'element-plus'
import api from "@/api/api";

const currentTrigger = ref<Nullable<ApiTrigger>>(null)
const triggerName = ref<Nullable<string>>(null)
const emit = defineEmits(['change', 'update:modelValue'])

const props = defineProps({
  modelValue: {
    type: Object as PropType<Nullable<ApiTrigger>>,
    default: () => undefined
  }
})

watch(
    () => props.modelValue,
    (val?: ApiTrigger) => {
      if (val === unref(currentTrigger)) return
      triggerName.value = val?.name || null;
      currentTrigger.value = val || null;
    },
)

// 监听
watch(
    () => currentTrigger.value,
    (val?: ApiTrigger) => {
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
  const {data} = await api.v1.triggerServiceSearchTrigger(params)
  const {items} = data
  cb(items)
}

const handleChange = (val) => {
  if (val == '') {
    currentTrigger.value = null
  }
}
const handleSelect = (val: ApiTrigger) => {
  currentTrigger.value = val
  emit('change', val)
}

</script>

<template>
  <ElAutocomplete
      class="w-[100%]"
      v-model="triggerName"
      :fetch-suggestions="querySearchAsync"
      placeholder="Please input"
      value-key="name"
      @select="handleSelect"
      @change="handleChange"
  />
</template>

<style lang="less">

</style>
