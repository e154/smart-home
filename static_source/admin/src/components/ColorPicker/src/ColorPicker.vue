<script setup lang="ts">

import {computed, ref, unref, watch} from "vue";
import {ElColorPicker} from 'element-plus'
import {propTypes} from "@/utils/propTypes";
import {useAppStore} from "@/store/modules/app";

const appStore = useAppStore()
const emit = defineEmits(['change', 'update:modelValue'])

const props = defineProps({
  modelValue: propTypes.string.def(''),
})

const currentColor = ref<string>(null)

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
  '#00adef'
]

const predefineColors = computed(() => currentColor.value? [...defaultColors, ...appStore.getLastColors] : [...defaultColors])

const updateColor = (val: string) => {
  let lastColors = appStore.getLastColors
  if (!defaultColors.includes(val) && !lastColors.includes(val)) {
    lastColors.push(val)
    if (lastColors.length > 5) {
      lastColors.shift()
    }
    appStore.setLastColors(lastColors)
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
