<script setup lang="ts">
import {nextTick, onMounted, PropType, computed, ref, unref, watch} from 'vue'
import {useAppStore} from "@/store/modules/app";
import JsonEditorVue from "json-editor-vue";

const emit = defineEmits(['change', 'update:modelValue'])
const appStore = useAppStore()

const props = defineProps({
  modelValue: {
    type: Object as PropType<Object>,
    default: () => null
  }
})

const sourceScript = ref(null)

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

onMounted(() => {

})

watch(
    () => props.modelValue,
    async (val?: any) => {
      await nextTick()
      if (val === unref(sourceScript)) return
      sourceScript.value = props.modelValue || null
    },
    {
      immediate: true
    }
)

type Mode = 'tree' | 'text' | 'table'
const mode = ref<Mode>('text')

const navigationBar = ref(false)
const statusBar = ref(false)
const readOnly = ref(false)
const mainMenuBar = ref(true)
const className = computed(() => appStore.getIsDark ? 'jse-theme-dark' : undefined)

const onChange = (val: any, cm: any) => {
  emit('change', val)
}

</script>

<template>
  <JsonEditorVue
      id="json-editor-vue"
      v-model="sourceScript"
      v-model:mode="mode"
      @change="onChange"
      :class="className"
      :read-only="readOnly"
      :main-menu-bar="mainMenuBar"
      :navigation-bar="navigationBar"
      :status-bar="statusBar"
  />
</template>

<style lang="less">
#json-editor-vue {
  height: 400px;
  min-height: 400px;
}

.jse-theme-dark {
  --jse-theme: dark;

  /* over all fonts, sizes, and colors */
  --jse-theme-color: #2f6dd0;
  --jse-theme-color-highlight: #467cd2;
  --jse-background-color: #1e1e1e;
  --jse-text-color: #d4d4d4;
  --jse-text-color-inverse: #4d4d4d;

  /* main, menu, modal */
  --jse-main-border: 1px solid #4f4f4f;
  --jse-menu-color: #fff;
  --jse-modal-background: #2f2f2f;
  --jse-modal-overlay-background: rgba(0, 0, 0, 0.5);
  --jse-modal-code-background: #2f2f2f;

  /* tooltip in text mode */
  --jse-tooltip-color: var(--jse-text-color);
  --jse-tooltip-background: #4b4b4b;
  --jse-tooltip-border: 1px solid #737373;
  --jse-tooltip-action-button-color: inherit;
  --jse-tooltip-action-button-background: #737373;

  /* panels: navigation bar, gutter, search box */
  --jse-panel-background: #333333;
  --jse-panel-background-border: 1px solid #464646;
  --jse-panel-color: var(--jse-text-color);
  --jse-panel-color-readonly: #737373;
  --jse-panel-border: 1px solid #3c3c3c;
  --jse-panel-button-color-highlight: #e5e5e5;
  --jse-panel-button-background-highlight: #464646;

  /* navigation-bar */
  --jse-navigation-bar-background: #656565;
  --jse-navigation-bar-background-highlight: #7e7e7e;
  --jse-navigation-bar-dropdown-color: var(--jse-text-color);

  /* context menu */
  --jse-context-menu-background: #4b4b4b;
  --jse-context-menu-background-highlight: #595959;
  --jse-context-menu-separator-color: #595959;
  --jse-context-menu-color: var(--jse-text-color);
  --jse-context-menu-pointer-background: #737373;
  --jse-context-menu-pointer-background-highlight: #818181;
  --jse-context-menu-pointer-color: var(--jse-context-menu-color);

  /* contents: json key and values */
  --jse-key-color: #9cdcfe;
  --jse-value-color: var(--jse-text-color);
  --jse-value-color-number: #b5cea8;
  --jse-value-color-boolean: #569cd6;
  --jse-value-color-null: #569cd6;
  --jse-value-color-string: #ce9178;
  --jse-value-color-url: #ce9178;
  --jse-delimiter-color: #949494;
  --jse-edit-outline: 2px solid var(--jse-text-color);

  /* contents: selected or hovered */
  --jse-selection-background-color: #464646;
  --jse-selection-background-inactive-color: #333333;
  --jse-hover-background-color: #343434;
  --jse-active-line-background-color: rgba(255, 255, 255, 0.06);
  --jse-search-match-background-color: #343434;

  /* contents: section of collapsed items in an array */
  --jse-collapsed-items-background-color: #333333;
  --jse-collapsed-items-selected-background-color: #565656;
  --jse-collapsed-items-link-color: #b2b2b2;
  --jse-collapsed-items-link-color-highlight: #ec8477;

  /* contents: highlighting of search results */
  --jse-search-match-color: #724c27;
  --jse-search-match-outline: 1px solid #966535;
  --jse-search-match-active-color: #9f6c39;
  --jse-search-match-active-outline: 1px solid #bb7f43;

  /* contents: inline tags inside the JSON document */
  --jse-tag-background: #444444;
  --jse-tag-color: #bdbdbd;

  /* contents: table */
  --jse-table-header-background: #333333;
  --jse-table-header-background-highlight: #424242;
  --jse-table-row-odd-background: rgba(255, 255, 255, 0.1);

  /* controls in modals: inputs, buttons, and `a` */
  --jse-input-background: #3d3d3d;
  --jse-input-border: var(--jse-main-border);
  --jse-button-background: #808080;
  --jse-button-background-highlight: #7a7a7a;
  --jse-button-color: #e0e0e0;
  --jse-button-secondary-background: #494949;
  --jse-button-secondary-background-highlight: #5d5d5d;
  --jse-button-secondary-background-disabled: #9d9d9d;
  --jse-button-secondary-color: var(--jse-text-color);
  --jse-a-color: #55abff;
  --jse-a-color-highlight: #4387c9;

  /* svelte-select */
  --jse-svelte-select-background: #3d3d3d;
  --jse-svelte-select-border: 1px solid #4f4f4f;
  --list-background: #3d3d3d;
  --item-hover-bg: #505050;
  --multi-item-bg: #5b5b5b;
  --input-color: #d4d4d4;
  --multi-clear-bg: #8a8a8a;
  --multi-item-clear-icon-color: #d4d4d4;
  --multi-item-outline: 1px solid #696969;
  --list-shadow: 0 2px 8px 0 rgba(0, 0, 0, 0.4);

  /* color picker */
  --jse-color-picker-background: #656565;
  --jse-color-picker-border-box-shadow: #8c8c8c 0 0 0 1px;
}

</style>
