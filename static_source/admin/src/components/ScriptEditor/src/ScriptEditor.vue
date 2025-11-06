<script setup lang="ts">
import {computed, nextTick, onMounted, onUnmounted, PropType, ref, unref, watch} from 'vue'
import {ApiScript} from "@/api/stub";
import {useAppStore} from "@/store/modules/app";
import {useEmitt} from "@/hooks/web/useEmitt";

import * as monaco from 'monaco-editor';
import {CodeEditor} from 'monaco-editor-vue3';
import api from "@/api/api";

const emit = defineEmits(['change', 'update:source', 'save'])
const appStore = useAppStore()

const props = defineProps({
  modelValue: {
    type: Object as PropType<Nullable<ApiScript>>,
    default: () => null
  }
})

const {emitter} = useEmitt()
const sourceScript = ref('')

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

const onKeydown = (e) => {
  const evtobj = window.event ? event : e
  // console.log(e);
  // 191 = /
  // if ((navigator.platform.match("Mac") ? e.metaKey : e.ctrlKey) && evtobj.keyCode == 191) {
  //   commentSelectedText()
  // }
  // // 70 = F
  // if ((navigator.platform.match("Mac") ? e.metaKey : e.ctrlKey) && e.shiftKey && evtobj.keyCode == 70) {
  //   autoFormatSelection()
  // }

  // 83 = S
  if ((navigator.platform.match("Mac") ? e.metaKey : e.ctrlKey) && evtobj.keyCode == 83) {
    onSave(e)
  }
}


onMounted(() => {
  document.onkeydown = onKeydown
})

onUnmounted(() => {
  document.onkeydown = null
})

const theme = computed(() => appStore.getIsDark ? "vs-dark" : "vs")

const language = ref('typescript')
watch(
  () => props.modelValue,
  async (val?: ApiScript) => {
    await nextTick()
    if (val?.source === unref(sourceScript)) return
    if (val) {
      sourceScript.value = val?.source || '';
      switch (val.lang) {
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
  },
  {
    immediate: true
  }
)

const editorOptions = ref({
  fontSize: fontSize,
  minimap: {enabled: true},
  automaticLayout: true,
  theme: theme,
  language: language,
})

const onChange = (val: string) => {
  // console.log(val)
  emitter.emit('updateSource', val)
  emit('update:source', val)
}

const editorInstance = ref()

const onMount = (editor) => {
  editorInstance.value = editor
  editor.focus()
}

const onSave = (e) => {
  e.preventDefault()
  emit('save')
}

useEmitt({
  name: 'updateEditor',
  callback: (val: string) => {
    setTimeout(() => {
      // console.log('update editor')

    }, 100)
  }
})

const registerGLobalScrit = (script: ApiScript) => {

  const disposable = monaco.languages.typescript.javascriptDefaults.addExtraLib(
    script.source,
    'filename.d.ts' // Имя файла должно заканчиваться на .d.ts
  );

  monaco.languages.typescript.typescriptDefaults.addExtraLib(script.source, 'filename.d.ts');

}

const fetchGlobalScript = async () => {
  if (!props.modelValue || props.modelValue?.name == "global.d") {
    return
  }

  let globalScript = appStore.getGlobalScript

  // check
  var {data} = await api.v1.scriptServiceSearchScript({query: "global.d", limit: 1, offset: 0})
  const {items} = data
  if (items.length == 0) {
    registerGLobalScrit(globalScript)
    return
  }

  // check storage
  if (globalScript) {

    // check update_at
    const fromServer = Date.parse(items[0].updatedAt)
    const fromCache = Date.parse(globalScript!.updatedAt)

    if (fromServer <= fromCache) {
      registerGLobalScrit(globalScript)
      return
    }
  }

  // fetch script
  var {data} = await api.v1.scriptServiceGetScriptById(items[0].id);
  if (!data) {
    registerGLobalScrit(data)
    return
  }

  console.info('update global script cache')
  appStore.setGlobalScript(data)
}

watch(
  () => props.modelValue,
  (newValue, oldValue) => {
    fetchGlobalScript()
  },
  { once: true }
);

</script>

<template>

  <CodeEditor
    style="height: 100%"
    v-model:value="sourceScript"
    :options="editorOptions"
    @change="onChange"
    @editorDidMount="onMount"
  >
    <template #loading="{ progress }">
      <div>Loading... {{ progress }}%</div>
    </template>
  </CodeEditor>

</template>

<style lang="less" scoped>

</style>
