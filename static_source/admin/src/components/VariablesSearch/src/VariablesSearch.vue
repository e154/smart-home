<script setup lang="ts">

import {PropType, ref, unref, watch} from "vue";
import {ApiVariable} from "@/api/stub";
import {ElSelect, ElOption} from 'element-plus'
import api from "@/api/api";

const options = ref<ApiVariable[]>([])
const value = ref<string[]>([])
const loading = ref(false)
const emit = defineEmits(['change', 'update:modelValue'])

const props = defineProps({
  modelValue: {
    type: Array as PropType<string[]>,
    default: () => []
  }
})

watch(
    () => props.modelValue,
    (val?: string[]) => {
      if (val === unref(value)) return
      value.value = val || [] ;
      if (val) {
        options.value = val.map(v => {
          return {name: v}
        })
      }
    },
  {
    immediate: true
  }
)

watch(
    () => value.value,
    (val?: string[]) => {
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

const searchMethod = async (query: string) => {
  loading.value = true
  const params = {query: query, limit: 25, offset: 0}
  const {data} = await api.v1.variableServiceSearchVariable(params)
      .catch(() => {
      })
      .finally(() => {
        loading.value = false
      })
  const {items} = data
  options.value = items
}

const handleSelect = (val: ApiVariable) => {
  emit('change', val)
}

</script>

<template>
  <ElSelect
      v-model="value"
      class="w-[100%]"
      multiple
      filterable
      remote
      reserve-keyword
      placeholder="Please enter a keyword"
      :remote-method="searchMethod"
      :loading="loading"
      @select="handleSelect"
  >
    <ElOption
        v-for="item in options"
        :key="item.name"
        :label="item.name"
        :value="item.name"
    />
  </ElSelect>

</template>

<style lang="less">

</style>
