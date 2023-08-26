<script setup lang="ts">

import {PropType, ref, unref, watch} from "vue";
import {ApiAction} from "@/api/stub";
import {ElAutocomplete} from 'element-plus'
import api from "@/api/api";

const currentAction = ref<Nullable<ApiAction>>(null)
const actionName = ref<Nullable<string>>(null)
const emit = defineEmits(['change', 'update:modelValue'])

const props = defineProps({
  modelValue: {
    type: Object as PropType<Nullable<ApiAction>>,
    default: () => undefined
  }
})

watch(
    () => props.modelValue,
    (val?: ApiAction) => {
      if (val === unref(currentAction)) return
      actionName.value = val?.name || null;
      currentAction.value = val || null;
    },
)

// 监听
watch(
    () => currentAction.value,
    (val?: ApiAction) => {
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
  const {data} = await api.v1.actionServiceSearchAction(params)
  const {items} = data
  cb(items)
}

const handleChange = (val) => {
  if (val == '') {
    currentAction.value = null
  }
}
const handleSelect = (val: ApiAction) => {
  currentAction.value = val
  emit('change', val)
}

</script>

<template>
  <ElAutocomplete
      class="w-[100%]"
      v-model="actionName"
      :fetch-suggestions="querySearchAsync"
      placeholder="Please input"
      value-key="name"
      @select="handleSelect"
      @change="handleChange"
  />
</template>

<style lang="less">

</style>
