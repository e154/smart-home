<script setup lang="ts">
import {nextTick, onMounted, PropType, ref, unref, watch} from 'vue'
import Codemirror, {CmComponentRef, CodeMirror} from "codemirror-editor-vue3";
import type {Editor, EditorConfiguration} from "codemirror";
import {ApiScript} from "@/api/stub";

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
import {HintDictionaryCoffee} from "@/views/Scripts/components/types";
import {useAppStore} from "@/store/modules/app";
import {bool} from "vue-types";
import {useEmitt} from "@/hooks/web/useEmitt";

const emit = defineEmits(['change', 'update:modelValue'])
const appStore = useAppStore()

const props = defineProps({
  modelValue: {
    type: Object as PropType<Nullable<ApiScript>>,
    default: () => null
  }
})

const { emitter } = useEmitt()
const sourceScript = ref('')
const cmComponentRef = ref<CmComponentRef>(null);
const cminstance = ref<Editor>();

const cmOptions: EditorConfiguration = {
  mode: "application/vnd.coffeescript", // Language mode
  gutters: ['CodeMirror-lint-markers'],
  theme: appStore.getIsDark ? "darcula" : "mdn-like", // Theme
  indentWithTabs: true,
  smartIndent: true,
  lineWrapping: true,
  // autoCloseBrackets: true, // Автоматическое закрытие скобок
  // matchBrackets: true, // Подсветка соответствующих скобок
  extraKeys: {
    'Ctrl-Space': 'autocomplete' // Комбинация клавиш для активации автодополнения
  },
  autofocus: true,
  hintOptions: {
    closeOnPick: true,
    completeSingle: false,
    hint: (editor: CodeMirror.Editor, options: object) => {
      var cur = editor.getCursor();
      var curLine = editor.getLine(cur.line);
      var start = cur.ch;
      var end = start;
      while (end < curLine.length && /[\w$]/.test(curLine.charAt(end))) ++end;
      while (start && /[\w$]/.test(curLine.charAt(start - 1))) --start;
      var curWord = start !== end && curLine.slice(start, end);
      var regex = new RegExp('^' + curWord, 'i');
      return {
        list: (!curWord ? [] : HintDictionaryCoffee.words.filter(function (item) {
          return item.text.match(regex);
        })).sort(),
        from: CodeMirror.Pos(cur.line, start),
        to: CodeMirror.Pos(cur.line, end)
      }
    },
  },
}

onMounted(() => {
  cminstance.value = cmComponentRef.value?.cminstance;
  cminstance.value?.focus();

})

watch(
    () => props.modelValue,
    async (val?: ApiScript) => {
      await nextTick()
      if (val?.source === unref(sourceScript)) return
      sourceScript.value = val?.source || '';
      if (val) {
        switch (val.lang) {
          case 'coffeescript':
            cmOptions.mode = "application/vnd.coffeescript"
            break
          case 'javascript':
            cmOptions.mode = "application/vnd.javascript"
            break
          case 'typescript':
            cmOptions.mode = "text/typescript"
            break
        }
      }
    },
    {
      immediate: true
    }
)

watch(
    () => appStore.getIsDark,
    (val: bool) => {
      console.log('change theme')
      cminstance.value?.setOption("theme", appStore.getIsDark ? "darcula" : "mdn-like")
      // cminstance.value?.refresh()
    }
)

const showEditorHint = (e: KeyboardEvent, handle: Function) => {
  // cminstance.value?.showHint()
}

const onChange = (val: string, cm: any) => {
  emitter.emit('updateSource', val)
}

</script>

<template>

  <Codemirror
      ref="cmComponentRef"
      v-model:value="sourceScript"
      :options="cmOptions"
      placeholder="test placeholder"
      @change="onChange"
      @keypress="showEditorHint"
  />

</template>

<style lang="less" scoped>

</style>
