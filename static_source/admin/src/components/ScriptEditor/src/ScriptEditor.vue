<script setup lang="ts">
import {onBeforeUnmount, onMounted, type PropType, ref, type Ref, unref, watch} from 'vue';
import {ApiScript} from "@/api/stub";
import {useAppStore} from "@/store/modules/app";
import {useEmitt} from "@/hooks/web/useEmitt";

import CodeMirror from 'vue-codemirror6';
import {esLint, javascript, javascriptLanguage} from '@codemirror/lang-javascript';
import type {LintSource,} from "@codemirror/lint";
import eslint from 'eslint-linter-browserify';
import {completeFromList} from "@codemirror/autocomplete";

import {Completions, snippets} from './completions';
import {darculaTheme} from './DarculaTheme'
import {lightTheme} from './LightTheme'

const appStore = useAppStore()

const cm: Ref<InstanceType<typeof CodeMirror> | undefined> = ref();

const props = defineProps({
  modelValue: {
    type: Object as PropType<Nullable<ApiScript>>,
    default: () => null
  }
})

const sourceScript: Ref<string> = ref('')

onMounted(() => {
  document.onkeydown = onKeydown
})

onBeforeUnmount(() => {
  document.onkeydown = null
})

// const lang = json();
let lang = javascript({
  jsx: true,
  typescript: true
});

const baseExtensions = [
  javascriptLanguage.data.of({
    autocomplete: completeFromList(snippets)
  }),
  javascriptLanguage.data.of({
    autocomplete: Completions
  }),
]

let extensions = []

const theme = ref()

watch(
  () => appStore.getIsDark,
  (val: boolean) => {
    extensions = val ? [...baseExtensions] : [...baseExtensions]
    theme.value = val ? darculaTheme : lightTheme
  },
  {
    immediate: true,
  }
)


let linter: LintSource | null = null
const focused: Ref<boolean> = ref(false);
const onFocus = (f: boolean): void => {
  focused.value = f;
};

watch(
  () => props.modelValue,
  async (val?: ApiScript) => {
    if (val?.source === unref(sourceScript)) return
    if (val) {
      sourceScript.value = val?.source || '';
      switch (val.lang) {
        case 'coffeescript':
          linter = null
          break
        case 'javascript':
          linter = esLint(
            // eslint-disable-next-line
            new eslint.Linter(),
            {
              parserOptions: {
                ecmaVersion: 2022,
                sourceType: 'module',
              },
              env: {
                browser: true,
                node: true,
              },
            }
          );
          break
        case 'typescript':
          linter = null
          break
      }
    }
  },
  {
    immediate: true
  }
)

const phrases: Ref<Record<string, string>> = ref({
  // @codemirror/view
  // 'Control character': '制御文字',
  // // @codemirror/commands
  // 'Selection deleted': '選択を削除',
  // // @codemirror/language
  // 'Folded lines': '折り畳まれた行',
  // 'Unfolded lines': '折り畳める行',
  // to: '行き先',
  // 'folded code': '折り畳まれたコード',
  // unfold: '折り畳みを解除',
  // 'Fold line': '行を折り畳む',
  // 'Unfold line': '行の折り畳む解除',
  // // @codemirror/search
  // 'Go to line': '行き先の行',
  // go: 'OK',
  // Find: '検索',
  // Replace: '置き換え',
  next: '▼',
  previous: '▲',
  // all: 'すべて',
  // 'match case': '一致条件',
  // 'by word': '全文検索',
  // regexp: '正規表現',
  // replace: '置き換え',
  // 'replace all': 'すべてを置き換え',
  // close: '閉じる',
  // 'current match': '現在の一致',
  // 'replaced $ matches': '$ 件の一致を置き換え',
  // 'replaced match on line $': '$ 行の一致を置き換え',
  // 'on line': 'した行',
  // // @codemirror/autocomplete
  // Completions: '自動補完',
  // // @codemirror/lint
  // Diagnostics: 'エラー',
  // 'No diagnostics': 'エラーなし',
})

const {emitter} = useEmitt()
const emit = defineEmits(['change', 'update:source', 'save'])
const onChange = (value: string) => {
  // console.log(value)
  emitter.emit('updateSource', value)
  emit('update:source', value)
}

const onSave = (e) => {
  e.preventDefault()
  emit('save')
}

const onKeydown = (e) => {
  const evtobj = window.event ? event : e
  // console.log(e);
  // 83 = S
  if ((navigator.platform.match("Mac") ? e.metaKey : e.ctrlKey) && evtobj.keyCode == 83) {
    onSave(e)
  }
}

</script>
<template>
  <code-mirror
    ref="cm"
    basic
    :dark="appStore.getIsDark"
    :lang="lang"
    v-model="sourceScript"
    :phrases="phrases"
    :extensions="extensions"
    :linter="linter"
    :theme="theme"
    gutter
    wrap
    tab
    :tab-size="2"
    allow-multiple-selections
    @focus="onFocus"
    @update:modelValue="onChange"
  />
</template>
