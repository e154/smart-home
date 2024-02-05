<script setup lang="ts">

import {PropType, ref, unref, watch} from "vue";
import {ApiCondition, ApiTrigger} from "@/api/stub";
import {ElSelect, ElOption} from 'element-plus'
import api from "@/api/api";

const options = ref<ApiCondition[]>([])
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
  const {data} = await api.v1.conditionServiceSearchCondition(params)
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
  const res = await api.v1.conditionServiceGetConditionList(params)
      .catch(() => {
      })
      .finally(() => {
      })
  if (res) {
    const {items, meta} = res.data;
    options.value = items || [];
  }
}

const handleSelect = (val: ApiCondition) => {
  emit('change', val)
}

const label = (item: ApiCondition): string => `${item.name} (id: ${item.id})`

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
        :label="label(item)"
        :value="item.id"
    >
      <span style="float: left">{{ item.name }} (id: {{ item.id }})</span>
    </ElOption>
  </ElSelect>

</template>

<style lang="less">

</style>
