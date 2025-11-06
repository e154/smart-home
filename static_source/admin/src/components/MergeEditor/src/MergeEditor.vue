<script setup lang="ts">
import {computed, onMounted, onUnmounted, PropType, reactive, ref, unref, watch} from 'vue'
import {ApiScript} from "@/api/stub";
import {useAppStore} from "@/store/modules/app";
import {useEmitt} from "@/hooks/web/useEmitt";

import { DiffEditor } from 'monaco-editor-vue3';

const appStore = useAppStore()
const emit = defineEmits(['change', 'update:source'])
const {emitter} = useEmitt()

const props = defineProps({
  source: {
    type: Object as PropType<Nullable<ApiScript>>,
    default: () => null
  },
  destination: {
    type: Object as PropType<Nullable<ApiScript>>,
    default: () => null
  }
})

// const reloadKey = ref(0)
// const reload = () => {
//   reloadKey.value += 1
// }

const currentSize = computed(() => appStore.getCurrentSize as string)
const fontSize = computed(() => {
  let size = 14;
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

const editorOptions = ref({
  fontSize: fontSize,
})

const theme = computed(() => appStore.getIsDark ? "vs-dark" : "vs")

const language = ref('typescript')
watch(
    () => [props.source, props.destination],
    async (value: ApiScript[]) => {
      if (value[0]) {
        // code.value = value[0].source
        switch (value[0].lang) {
          case 'coffeescript':
            language.value = 'coffeescript'
            break
          case 'javascript':
            language.value = 'javascript'
            break
          case 'typescript':
            language.value = 'typescript'
            break
        }
      }
      if (value[1]) {
        // orig2.value = value[1].source
      }
      // reload()
    },
    {
      immediate: true,
    }
)

const onChange = (val: string) => {
  // console.log('onChange', val)
  emit("update:source", val);
}

const originalCode = computed(() => props.source?.source || '')
const modifiedCode = computed(() => props.destination?.source || '')

</script>

<template>
    <DiffEditor
      style="height: 100%"
      v-model:value="originalCode"
      :options="editorOptions"
      :original="modifiedCode"
      :language="language"
      :theme="theme"
      @change="onChange"
    />
</template>

<style lang="less">

</style>
