<script setup lang="ts">

import {computed, ref, unref, watch} from "vue";
import {ElColorPicker} from 'element-plus'
import {propTypes} from "@/utils/propTypes";

const emit = defineEmits(['change', 'update:modelValue'])

const props = defineProps({
  modelValue: propTypes.string.def(''),
})

const currentColor = ref<string>(null)
const lastColors = ref<string[]>([])

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

const defaultColors = [
  '#000000',
  '#232324',
  '#F5F7FA',
  '#ffffff',
]

const predefineColors = computed(() => currentColor.value? [...defaultColors, ...lastColors.value] : [...defaultColors])

const updateColor = (val: string) => {
  if (!defaultColors.includes(val) && !lastColors.value.includes(val)) {
    lastColors.value.push(val)
    if (lastColors.value.length > 5) {
      lastColors.value.shift()
    }
  }
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
