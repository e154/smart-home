<script setup lang="ts">

import {PropType, ref, unref, watch} from "vue";
import {ApiRole} from "@/api/stub";
import {ElAutocomplete} from 'element-plus'
import api from "@/api/api";

const currentRole = ref<Nullable<ApiRole>>(null)
const roleName = ref<Nullable<string>>(null)
const emit = defineEmits(['change', 'update:modelValue'])

const props = defineProps({
  modelValue: {
    type: Object as PropType<Nullable<ApiRole>>,
    default: () => undefined
  }
})

watch(
    () => props.modelValue,
    (val?: ApiRole) => {
      if (val === unref(currentRole)) return
      roleName.value = val?.name || null;
      currentRole.value = val || null;
    },
)

// 监听
watch(
    () => currentRole.value,
    (val?: ApiRole) => {
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
  const {data} = await api.v1.roleServiceSearchRoleByName(params)
  const {items} = data
  cb(items)
}

const handleChange = (val) => {
  if (val == '') {
    currentRole.value = null
  }
}
const handleSelect = (val: ApiRole) => {
  currentRole.value = val
  emit('change', val)
}

</script>

<template>
  <ElAutocomplete
      class="w-[100%]"
      v-model="roleName"
      :fetch-suggestions="querySearchAsync"
      placeholder="Please input"
      value-key="name"
      @select="handleSelect"
      @change="handleChange"
  />
</template>

<style lang="less">

</style>
