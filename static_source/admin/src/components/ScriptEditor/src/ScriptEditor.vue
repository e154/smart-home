<script setup lang="ts">
import {computed, nextTick, onMounted, onUnmounted, PropType, ref, unref, watch} from 'vue'
import Codemirror, {CmComponentRef, CodeMirror} from "codemirror-editor-vue3";
import type {Editor, EditorConfiguration} from "codemirror";
import {ApiScript} from "@/api/stub";
import prettier from 'prettier/standalone';
import parserBabel from 'prettier/parser-babel';

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
import {HintDictionary, HintDictionaryCoffee} from "@/views/Scripts/components/types";
import {useAppStore} from "@/store/modules/app";
import {useEmitt} from "@/hooks/web/useEmitt";

const emit = defineEmits(['change', 'update:source', 'save'])
const appStore = useAppStore()

const props = defineProps({
  modelValue: {
    type: Object as PropType<Nullable<ApiScript>>,
    default: () => null
  }
})

const { emitter } = useEmitt()
const sourceScript = ref('')
const cmComponentRef = ref<CmComponentRef>();
const cminstance = ref<Editor>();

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

const getHint = () => {
  switch (cmOptions.mode) {
    case 'application/vnd.coffeescript':
      return HintDictionaryCoffee
    default:
      return HintDictionary
  }
}

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
        list: (!curWord ? [] : getHint().words.filter(function (item) {
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

  document.onkeydown = onKeydown
})

onUnmounted(() => {
  document.onkeydown = null
})


watch(
    () => props.modelValue,
    async (val?: ApiScript) => {
      await nextTick()
      if (val?.source === unref(sourceScript)) return
      if (val) {
        sourceScript.value = val?.source || '';
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
      cminstance.value?.setOption("theme", appStore.getIsDark ? "darcula" : "mdn-like")
      cminstance.value?.refresh()
    }
)

const showEditorHint = (e: KeyboardEvent, handle: Function) => {
  // cminstance.value?.showHint()
}

const onChange = (val: string, cm: any) => {
  // console.log(val)
  // console.log(cm.getValue())
  emitter.emit('updateSource', val)
  emit('update:source', val)
}

const onSave = (e) => {
  e.preventDefault()
  emit('save')
}

const autoFormatSelection = () => {
  let plugins = [parserBabel]
  let parser = "babel"
  if (props.modelValue?.lang == 'coffeescript') {
    console.warn("coffeescript prettier plugin not installed")
    return
  }
  const code = cminstance.value?.getValue(); // Получите весь код из редактора
  const formattedCode = prettier.format(code, {
    "singleQuote": true, // Использовать одинарные кавычки для строк
    "semi": true, // Добавлять точку с запятой в конце выражений
    "trailingComma": "none", // Не использовать запятую в конце массивов и объектов
    "tabWidth": 2, // Количество пробелов для одного уровня отступа
    "printWidth": 180, // Максимальная длина строки кода
    parser: parser,
    plugins: plugins
  });

  // Замените весь код отформатированным кодом
  cminstance.value?.setValue(formattedCode);
}

const commentSelectedText = () => {
  const fromLine = cminstance.value?.getCursor("from").line; // Начальная строка выделения
  const toLine = cminstance.value?.getCursor("to").line; // Конечная строка выделения
  let commentSymbol = "// "; // Символ комментария (можно изменить по вашему желанию)

  if (props.modelValue?.lang == 'coffeescript') {
    commentSymbol = "#"
  }

  // Перебираем строки выделенного текста и добавляем комментарии
  for (var i = fromLine; i <= toLine; i++) {
    let lineText = cminstance.value?.getLine(i);
    if (lineText?.startsWith(commentSymbol)) {
      lineText = lineText?.replace(commentSymbol, "")
      cminstance.value?.replaceRange(lineText, { line: i, ch: 0 }, { line: i, ch: lineText.length+3 }); // Добавляем комментарий
    } else {
      cminstance.value?.replaceRange(commentSymbol + lineText, { line: i, ch: 0 }, { line: i, ch: lineText.length }); // Добавляем комментарий
    }
  }
}

const onKeydown = ( e ) => {
  var evtobj = window.event? event : e
  // console.log(e);
  // 191 = /
  if ((navigator.platform.match("Mac") ? e.metaKey : e.ctrlKey) && evtobj.keyCode == 191) {
    commentSelectedText()
  }
  // 70 = F
  if ((navigator.platform.match("Mac") ? e.metaKey : e.ctrlKey) && e.shiftKey && evtobj.keyCode == 70) {
    autoFormatSelection()
  }

  // 83 = S
  if ((navigator.platform.match("Mac") ? e.metaKey : e.ctrlKey) && evtobj.keyCode == 83) {
    onSave(e)
  }
}

useEmitt({
  name: 'updateEditor',
  callback: (val: string) => {
    setTimeout(() => {
      // console.log('update editor')
      cminstance.value?.refresh()
      cminstance.value?.focus();
    }, 100)
  }
})

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
:deep(.CodeMirror) {
  font-size: v-bind(fontSize);
  line-height: 1.5;
}
</style>
