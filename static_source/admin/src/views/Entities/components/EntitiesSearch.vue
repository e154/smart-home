<script setup lang="ts">

import {PropType, ref, unref, watch} from "vue";
import {ApiEntityShort} from "@/api/stub";
import {ElSelect, ElOption} from 'element-plus'
import api from "@/api/api";

const options = ref<ApiEntityShort[]>([])
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

const remoteMethod = async (query: string) => {
  loading.value = true
  const params = {query: query, limit: 25, offset: 0}
  const {data} = await api.v1.entityServiceSearchEntity(params)
      .catch(() => {
      })
      .finally(() => {
        loading.value = false
      })
  const {items} = data
  options.value = items
}

remoteMethod("")

const handleSelect = (val: ApiEntityShort) => {
  console.log(val)
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
      :remote-method="remoteMethod"
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
