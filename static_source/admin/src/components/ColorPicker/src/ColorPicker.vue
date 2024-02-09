<script setup lang="ts">

import {ref, unref, watch} from "vue";
import {ElColorPicker} from 'element-plus'
import {propTypes} from "@/utils/propTypes";

const emit = defineEmits(['change', 'update:modelValue'])

const props = defineProps({
  modelValue: propTypes.string.def(''),
})

const currentColor = ref(null)

watch(
    () => props.modelValue,
    (val?: string) => {
      if (val === unref(currentColor)) return
      currentColor.value = val
    },
    {
      immediate: true
    }
)

const predefineColors = ref([
  '#000',
  '#232324',
  '#F5F7FA',
  '#fff',
])

const updateColor = (val: string) => {
  emit('change', val)
  emit('update:modelValue', val)
}

</script>

<template>
  <ElColorPicker
      v-model="currentColor"
      v-on:change="updateColor"
      :predefine="predefineColors"
  />
</template>

<style lang="less">

</style>
