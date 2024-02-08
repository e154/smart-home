<script setup lang="ts">

import {PropType, ref, unref, watch} from "vue";
import {ApiDashboard} from "@/api/stub";
import {ElAutocomplete} from 'element-plus'
import api from "@/api/api";

const currentDashboard = ref<Nullable<ApiDashboard>>(null)
const boardName = ref<Nullable<string>>(null)
const emit = defineEmits(['change', 'update:modelValue'])

const props = defineProps({
  modelValue: {
    type: Object as PropType<Nullable<ApiDashboard>>,
    default: () => undefined
  }
})

watch(
    () => props.modelValue,
    (val?: ApiDashboard) => {
      if (val === unref(currentDashboard)) return
      boardName.value = val?.name || null;
      currentDashboard.value = val || null;
    },
    {
      immediate: true
    }
)

// 监听
watch(
    () => currentDashboard.value,
    (val?: ApiDashboard) => {
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
  const {data} = await api.v1.dashboardServiceSearchDashboard(params)
  const {items} = data
  cb(items)
}

const handleChange = (val) => {
  if (val == '') {
    currentDashboard.value = null
  }
}
const handleSelect = (val: ApiDashboard) => {
  currentDashboard.value = val
  emit('change', val)
}

</script>

<template>
  <ElAutocomplete
      class="w-[100%]"
      v-model="boardName"
      :fetch-suggestions="querySearchAsync"
      placeholder="Please input"
      value-key="name"
      @select="handleSelect"
      @change="handleChange"
      :clearable="true"
  />
</template>

<style lang="less">

</style>
