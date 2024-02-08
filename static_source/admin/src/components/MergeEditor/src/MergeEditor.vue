<script setup lang="ts">
import {computed, onMounted, onUnmounted, PropType, reactive, ref, unref, watch} from 'vue'
import Codemirror, {CmComponentRef} from "codemirror-editor-vue3";
import {Editor} from "codemirror";
import {ApiScript} from "@/api/stub";
import {useAppStore} from "@/store/modules/app";
import {useEmitt} from "@/hooks/web/useEmitt";
import {MergeView} from "codemirror/addon/merge/merge";

import "codemirror/mode/htmlmixed/htmlmixed.js";

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

const reloadKey = ref(0)
const reload = () => {
  reloadKey.value += 1
}

onMounted(() => {
  setTimeout(() => {
    reload()
  }, 100)
})

onUnmounted(() => {

})

const code = ref(``)
// const code = computed(()=> props.source?.source || '')
const orig2 = ref(``)
// const orig2 = computed(()=> props.destination?.source || '')
const mode = ref('application/vnd.coffeescript')

const cmOptions = reactive({
  value: code,
  orig: orig2,
  theme: appStore.getIsDark ? "darcula" : "mdn-like", // Theme
  origLeft: null,
  connect: "align",
  // mode: "application/vnd.coffeescript", // Language mode
  mode: mode.value,
  lineNumbers: true,
  collapseIdentical: true,
  highlightDifferences: true
})

watch(
    () => [props.source, props.destination],
    async (value: ApiScript[]) => {
      if (value[0]) {
        code.value = value[0].source
        switch (value[0].lang) {
          case 'coffeescript':
            mode.value = "application/vnd.coffeescript"
            break
          case 'javascript':
            mode.value = "application/vnd.javascript"
            break
          case 'typescript':
            mode.value = "text/typescript"
            break
        }
      }
      if (value[1]) {
        orig2.value = value[1].source
      }
      reload()
    },
    {
      immediate: true,
    }
)

watch(
    () => appStore.getIsDark,
    (val) => {
      cmOptions.theme = appStore.getIsDark ? "darcula" : "mdn-like";
      reload()
    }
)

const onReady = (cm: any) => {

}

const onChange = (val: string, cm: Editor) => {
  // const cmMerge = cm as MergeView
  // const cminstance: Editor = cmMerge.editor()
  emit("update:source", val);
}

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



</script>

<template>

  <Codemirror
      :key="reloadKey"
      class="codemirror-merge"
      merge
      :options="cmOptions"
      @ready="onReady"
      @change="onChange"
  />

</template>

<style lang="less">

.codemirror-merge {
  .CodeMirror {
    font-size: v-bind(fontSize);
    line-height: 1.5;
  }
}

html.dark {
  .codemirror-merge.codemirror-container.bordered {
    border-radius: 4px;
    border: 1px solid var(--el-bg-color-overlay);

    .CodeMirror-merge-gap {
      border-color: var(--el-bg-color-overlay);
      background: #313335;
    }

    .CodeMirror-merge-copy {
      color: #fff;
    }

    .CodeMirror-merge-r-chunk-end {
      border-bottom: 1px solid #ffffff30;
    }

    .CodeMirror-merge-r-chunk-start {
      border-top: 1px solid #ffffff30;
    }

    .CodeMirror-merge-r-chunk {
      background: #ffffff30;
    }
    .CodeMirror-merge-collapsed-widget {
      color: #eee;
      background: #565656;
      border-color: #565656;
    }
  }
}

</style>
