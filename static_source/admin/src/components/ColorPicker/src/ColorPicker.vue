<script setup lang="ts">

import {computed, ref, unref, watch} from "vue";
import {ElColorPicker} from 'element-plus'
import {propTypes} from "@/utils/propTypes";
import {useAppStore} from "@/store/modules/app";

const appStore = useAppStore()
const emit = defineEmits(['change', 'update:modelValue'])

const props = defineProps({
  modelValue: propTypes.string.def(''),
  showAlpha: propTypes.bool.def(false),
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
  'rgba(0, 0, 0, 1)', //#000000
  'rgba(35, 35, 36, 1)', //#232324
  'rgba(245, 247, 250, 1)', //#F5F7FA
  'rgba(255, 255, 255, 1)', //#FFFFFF
  'rgba(0, 173, 239, 1)' //#00ADEF
]

const predefineColors = computed(() => [...defaultColors, ...appStore.getLastColors])

const updateColor = (val: string) => {
  let lastColors = appStore.getLastColors
  if (!predefineColors.value.includes(val)) {
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
      :show-alpha="showAlpha"
  />
</template>

<style lang="less">

</style>
