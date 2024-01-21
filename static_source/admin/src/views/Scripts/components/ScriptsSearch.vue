<script setup lang="ts">

import {PropType, ref, unref, watch} from "vue";
import {ApiScript} from "@/api/stub";
import {ElSelect, ElOption} from 'element-plus'
import api from "@/api/api";

const options = ref<ApiScript[]>([])
const value = ref<number[]>([])
const loading = ref(false)
const emit = defineEmits(['change', 'update:modelValue'])

const props = defineProps({
  modelValue: {
    type: Array as PropType<number[]>,
    default: () => []
  }
})

watch(
    () => props.modelValue,
    (val?: number[]) => {
      if (val === unref(value)) return
      value.value = val || [] ;
      if (val) {
        getList(val)
      }
    },
)

// 监听
watch(
    () => value.value,
    (val?: number[]) => {
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
  const {data} = await api.v1.scriptServiceSearchScript(params)
      .catch(() => {
      })
      .finally(() => {
        loading.value = false
      })
  const {items} = data
  options.value = items
}

interface Params {
  page?: number;
  limit?: number;
  sort?: string;
  ids?: [];
}

const getList = async (ids:  number[]) => {
  let params: Params = {
    ids: ids,
  }
  const res = await api.v1.scriptServiceGetScriptList(params)
      .catch(() => {
      })
      .finally(() => {
      })
  if (res) {
    const {items, meta} = res.data;
    options.value = items || [];
  }
}

const handleSelect = (val: ApiScript) => {
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
        :key="item.id"
        :label="item.name"
        :value="item.id"
    />
  </ElSelect>

</template>

<style lang="less">

</style>
