<script setup lang="ts">
import {onBeforeUnmount, onMounted, type PropType, ref, watch} from 'vue';
import {ApiScript} from "@/api/stub";
import {useAppStore} from "@/store/modules/app";

import {MergeView} from "@codemirror/merge"
import {basicSetup, Editor, EditorView} from "codemirror"
import {esLint, javascript, javascriptLanguage} from '@codemirror/lang-javascript';
import type {Transaction} from "@codemirror/state";
// import {darculaTheme, lightTheme} from '@/components/ScriptEditor';

const appStore = useAppStore()

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

const updateA = (value: string) => {
  if (!view.value) return
  view.value.a.dispatch({
    changes: {from: 0, to: view.value.a.state.doc.length, insert: value || ''},
    // selection: view.value.a.state.selection,
    // scrollIntoView: true,
  });
}

const updateB = (value: string) => {
  if (!view.value) return
  view.value.b.dispatch({
    changes: {from: 0, to: view.value.b.state.doc.length, insert: value || ''},
    // selection: view.value.b.state.selection,
    // scrollIntoView: true,
  });
}

const el = ref(null)
const view = ref(null)
const setup = () => {
  view.value = new MergeView({
    a: {
      doc: '',
      extensions: [
        basicSetup,
        javascript(),
      ],
    },
    b: {
      doc: '',
      extensions: [
        basicSetup,
        EditorView.editable.of(false),
        javascript()
      ],
    },
    orientation: 'a-b',
    revertControls: 'b-to-a',
    parent: el.value,
    // parent: el.value.childNodes[0]
  })
  console.log(view.value)

  // view.value.b.dispatch = (tr: Transaction) => {
  //   console.log('---', tr)
  // }

  if (props.source) {
    updateA(props.source.source)
  }
  if (props.destination) {
    updateB(props.destination.source)
  }

 setTimeout(() => {
   view.value.a.dispatch = (tr: Transaction) => {
     console.log(tr)
     // console.log(view.value.a.state)
     // console.log(view.value.a.viewState.state)
     // console.log(view.value.a.state.doc.length)

     // const xxx = {
     //   changes: {
     //     from: tr.changes.from,
     //     to: tr.changes.to,
     //     insert: tr.changes.insert,
     //     // insert: "qwe",
     //     startState: view.value.a.viewState.state,
     //   },
     //   userEvent: tr.userEvent,
     // }
     // console.log(xxx)


     tr.annotation = (v) => {

     }

     tr.startState = view.value.a.viewState.state

     view.value.a.update([tr]);
     // if (tr.changes.empty || !tr.docChanged) {
     //   if not change value, no fire emit event
     // return;
     // }

     // console.log('---', tr)
   }
 }, 200)
}

const destroy = () => {
  if (view.value) {
    view.value.destroy()
    view.value = null
  }
}


onMounted(() => {
  setup()
})

onBeforeUnmount(() => {
  destroy()
})

watch(
  () => [props.source, props.destination],
  async (value: ApiScript[]) => {
    if (value[0]) {
      updateA(value[0].source)
    }
    if (value[1]) {
      updateB(value[1].source)
    }
  },
  {
    immediate: false,
  }
)

const emit = defineEmits(['change', 'update:source'])
const onChange = (val: string, cm: Editor) => {
  console.log('----', val)
  // const cmMerge = cm as MergeView
  // const cminstance: Editor = cmMerge.editor()
  emit("update:source", val);
}

</script>
<template>
  <div ref="el"
       @update:modelValue="onChange"></div>
</template>
