<script setup lang="ts">

import {PropType, ref, unref, watch} from "vue";
import {ApiPluginShort} from "@/api/stub";
import {ElAutocomplete} from 'element-plus'
import api from "@/api/api";

const currentPlugin = ref<Nullable<ApiPluginShort>>(null)
const pluginName = ref<Nullable<string>>(null)
const emit = defineEmits(['change', 'update:modelValue'])

const props = defineProps({
  modelValue: {
    type: Object as PropType<Nullable<ApiPluginShort>>,
    default: () => undefined
  }
})

watch(
    () => props.modelValue,
    (val?: ApiPluginShort) => {
      if (val === unref(currentPlugin)) return
      pluginName.value = val?.name || null;
      currentPlugin.value = val || null;
    },

)

// 监听
watch(
    () => currentPlugin.value,
    (val?: ApiPluginShort) => {
      if (props.modelValue == unref(val)) return;
      emit('update:modelValue', val)
      if (!val) {
        emit('change', val)
      }
    },
    {
      immediate: true
    }
)

const querySearchAsync = async (queryString: string, cb: Fn) => {
  if (queryString == null || queryString == 'null') {
    queryString = ''
  }
  const params = {query: queryString, limit: 25, offset: 0}
  const {data} = await api.v1.pluginServiceSearchPlugin(params)
  const {items} = data
  cb(items)
}

const handleChange = (val) => {
  if (val == '') {
    currentPlugin.value = null
  }
}
const handleSelect = (val: ApiPluginShort) => {
  currentPlugin.value = val
  emit('change', val)
}

</script>

<template>
  <ElAutocomplete
      class="w-[100%]"
      v-model="pluginName"
      :fetch-suggestions="querySearchAsync"
      placeholder="Please input"
      value-key="name"
      @select="handleSelect"
      @change="handleChange"
  />
</template>

<style lang="less">

</style>
