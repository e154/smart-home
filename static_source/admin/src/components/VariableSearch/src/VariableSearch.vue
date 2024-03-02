<script setup lang="ts">

import {PropType, ref, unref, watch} from "vue";
import {ApiVariable} from "@/api/stub";
import {ElAutocomplete} from 'element-plus'
import api from "@/api/api";

const currentVariable = ref<Nullable<ApiVariable>>(null)
const variableName = ref<Nullable<string>>(null)
const emit = defineEmits(['change', 'update:modelValue'])

const props = defineProps({
  modelValue: {
    type: Object as PropType<Nullable<ApiVariable>>,
    default: () => undefined
  }
})

watch(
    () => props.modelValue,
    (val?: ApiVariable) => {
      if (val === unref(currentVariable)) return
      variableName.value = val?.name || null;
      currentVariable.value = val || null;
    },
)

// 监听
watch(
    () => currentVariable.value,
    (val?: ApiVariable) => {
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
  const {data} = await api.v1.variableServiceSearchVariable(params)
  const {items} = data
  cb(items)
}

const handleChange = (val) => {
  if (val == '') {
    currentVariable.value = null
  }
}
const handleSelect = (val: ApiVariable) => {
  currentVariable.value = val
  emit('change', val)
}

</script>

<template>
  <ElAutocomplete
      class="w-[100%]"
      v-model="variableName"
      :fetch-suggestions="querySearchAsync"
      placeholder="Please input"
      value-key="name"
      @select="handleSelect"
      @change="handleChange"
  />
</template>

<style lang="less">

</style>
