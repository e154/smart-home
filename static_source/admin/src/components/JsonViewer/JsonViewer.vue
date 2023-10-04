<script setup lang="ts">
import {nextTick, onMounted, PropType, ref, unref, watch} from 'vue'
import Codemirror, {CmComponentRef} from "codemirror-editor-vue3";
import type {Editor, EditorConfiguration} from "codemirror";

// codemirror
// placeholder
import "codemirror/addon/display/placeholder.js";
// language
import "codemirror/mode/javascript/javascript.js";
import "codemirror/mode/jsx/jsx.js";
import "codemirror/mode/coffeescript/coffeescript";
// theme
import "codemirror/theme/darcula.css";
import 'codemirror/theme/mdn-like.css'
import 'codemirror/addon/lint/lint.css'
import 'codemirror/addon/lint/lint'
import 'codemirror/addon/lint/coffeescript-lint'
import 'codemirror/addon/hint/show-hint';
import 'codemirror/addon/hint/javascript-hint';
import 'codemirror/addon/hint/show-hint.css';
import {useAppStore} from "@/store/modules/app";
import {bool} from "vue-types";
import {useEmitt} from "@/hooks/web/useEmitt";

const emit = defineEmits(['change', 'update:modelValue'])
const appStore = useAppStore()

const props = defineProps({
  modelValue: {
    type: Object as PropType<Object>,
    default: () => null
  }
})

const {emitter} = useEmitt()
const sourceScript = ref('')
const cmComponentRef = ref<CmComponentRef>(null);
const cminstance = ref<Editor>();

const cmOptions: EditorConfiguration = {
  mode: "application/json", // Language mode
  gutters: ['CodeMirror-lint-markers'],
  theme: appStore.getIsDark ? "darcula" : "mdn-like", // Theme
  indentWithTabs: true,
  smartIndent: true,
  lineWrapping: true,
  autofocus: true,
  hintOptions: {
    closeOnPick: true,
    completeSingle: false,
  },
}

onMounted(() => {
  cminstance.value = cmComponentRef.value?.cminstance;
  cminstance.value?.focus();
})

watch(
    () => props.modelValue,
    async (val?: any) => {
      await nextTick()
      if (val === unref(sourceScript)) return
      if (val) {
        sourceScript.value = JSON.stringify(val, null, 2);
      } else {
        sourceScript.value = ""
      }
      cminstance.value?.refresh()
    },
    {
      immediate: true
    }
)

watch(
    () => appStore.getIsDark,
    (val: bool) => {
      cminstance.value?.setOption("theme", appStore.getIsDark ? "darcula" : "mdn-like")
      cminstance.value?.refresh()
    }
)

const onChange = (val: string, cm: any) => {
  emitter.emit('updateSource', val)
}

</script>

<template>

  <Codemirror
      ref="cmComponentRef"
      v-model:value="sourceScript"
      :options="cmOptions"
      @change="onChange"
  />

</template>

<style lang="less" scoped>

</style>
