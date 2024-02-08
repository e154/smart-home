<script setup lang="ts">

import {PropType, ref, unref, watch} from "vue";
import {ApiScript} from "@/api/stub";
import {ElAutocomplete} from 'element-plus'
import api from "@/api/api";

const currentScript = ref<Nullable<ApiScript>>(null)
const scriptName = ref<Nullable<string>>(null)
const emit = defineEmits(['change', 'update:modelValue'])

const props = defineProps({
  modelValue: {
    type: Object as PropType<Nullable<ApiScript>>,
    default: () => undefined
  }
})

watch(
    () => props.modelValue,
    (val?: ApiScript) => {
      if (val === unref(currentScript)) return
      scriptName.value = val?.name || null;
      currentScript.value = val || null;
    },
)

// 监听
watch(
    () => currentScript.value,
    (val?: ApiScript) => {
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
  const {data} = await api.v1.scriptServiceSearchScript(params)
  const {items} = data
  cb(items)
}

const handleChange = (val) => {
  if (val == '') {
    currentScript.value = null
  }
}
const handleSelect = (val: ApiScript) => {
  currentScript.value = val
  emit('change', val)
}

</script>

<template>
  <ElAutocomplete
      class="w-[100%]"
      v-model="scriptName"
      :fetch-suggestions="querySearchAsync"
      placeholder="Please input"
      value-key="name"
      @select="handleSelect"
      @change="handleChange"
      clearable
  >
    <template #default="{ item }">
      <div class="value">{{ item.name }} id:({{ item.id }})</div>
    </template>
  </ElAutocomplete>
</template>

<style lang="less">

</style>
