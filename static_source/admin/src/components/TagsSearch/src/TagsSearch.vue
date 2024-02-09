<script setup lang="ts">

import {PropType, ref, unref, watch} from "vue";
import {ElOption, ElSelect} from 'element-plus'
import api from "@/api/api";
import {ApiTag} from "@/api/stub";

const options = ref<ApiTag[]>([])
const value = ref<ApiTag[]>([])
const loading = ref(false)
const emit = defineEmits(['change', 'update:modelValue'])

const props = defineProps({
  modelValue: {
    type: Array as PropType<string[]>,
    default: () => []
  }
})

const getList = async (tags: string[]) => {
  let params: Params = {
    tags: tags,
  }
  const res = await api.v1.tagServiceGetTagList(params)
      .catch(() => {
      })
      .finally(() => {
      })
  if (res) {
    const {items, meta} = res.data;
    options.value = items || [];
  }
}

watch(
    () => props.modelValue,
    (val?: string[]) => {
      if (val === unref(value)) return
      value.value = val || [] ;
      if (val && val.length) {
        getList(val)
      }
    },
    {
      immediate: true
    }
)

// 监听
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
  const {data} = await api.v1.tagServiceSearchTag(params)
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

const handleSelect = (val: any) => {
  emit('change', val)
}

</script>

<template>
  <ElSelect
      v-model="value"
      class="w-[100%]"
      multiple
      filterable
      allow-create
      remote
      clearable
      default-first-option
      :reserve-keyword="false"
      placeholder="Please enter a keyword"
      :remote-method="searchMethod"
      :loading="loading"
      @change="handleSelect"
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
