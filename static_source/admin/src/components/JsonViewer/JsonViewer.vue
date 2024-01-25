<script setup lang="ts">
import {onMounted, PropType, computed, unref, watch, reactive, ref} from 'vue'

import {useAppStore} from "@/store/modules/app";

const emit = defineEmits(['change', 'update:modelValue'])
const appStore = useAppStore()

const props = defineProps({
  modelValue: {
    type: Object as PropType<Object>,
    default: () => null
  }
})

const value = computed(() => props.modelValue || 'NO DATA')
const currentSize = computed(() => appStore.getCurrentSize as string)
const fontSize = computed(() => {
  let size = 16;
  switch (unref(currentSize)) {
    case "default":
      size = 14;
      break
    case "large":
      size = 16;
      break
    case "small":
      size = 12;
      break
  }
  return size + 'px'
})

const theme = computed(() => appStore.getIsDark ? "jv-darkula" : "jv-light")
const depth = ref(3)

</script>

<template>
  <json-viewer
      sort
      copyable
      expanded
      :theme="theme"
      :expand-depth="depth"
      :value="value"/>

</template>

<style lang="less">
// values are default one from jv-light template
.jv-darkula {
  //background: #282A36;
  white-space: nowrap;
  color: #F8F8F2;
  font-size: v-bind(fontSize);
  font-family: Consolas, Menlo, Courier, monospace;

  .jv-ellipsis {
    color: #F8F8F2;
    background-color: #44475A;
    display: inline-block;
    line-height: 0.9;
    font-size: 0.9em;
    padding: 0px 4px 2px 4px;
    border-radius: 3px;
    vertical-align: 2px;
    cursor: pointer;
    user-select: none;
  }
  .jv-button { color: #49b3ff }
  .jv-key { color: #F8F8F2 }
  .jv-item {
    &.jv-array { color: #F8F8F2 }
    &.jv-boolean { color: #fc1e70 }
    &.jv-function { color: #067bca }
    &.jv-number { color: #fc1e70 }
    &.jv-number-float { color: #fc1e70 }
    &.jv-number-integer { color: #fc1e70 }
    &.jv-object { color: #F8F8F2 }
    &.jv-undefined { color: #e08331 }
    &.jv-string {
      color: #42b983;
      word-break: break-word;
      white-space: normal;
    }
  }
  .jv-code {
    .jv-toggle {
      &:before {
        padding: 0px 2px;
        border-radius: 2px;
      }
      &:hover {
        &:before {
          background: #eee;
        }
      }
    }
  }
}
</style>
