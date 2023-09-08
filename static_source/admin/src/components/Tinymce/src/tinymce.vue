<template>
  <editor v-model="editorValue" :init="init" class="w-[100%]"/>
</template>

<script setup>
import { reactive, ref, watch, toRefs } from 'vue';

import tinymce from 'tinymce/tinymce.js';
import 'tinymce/models/dom'; //(TinyMCE 6)

// 外觀
import 'tinymce/skins/ui/oxide/skin.css';
// import 'tinymce/skins/ui/oxide-dark/skin.css';
// import contentUiCss from 'tinymce/skins/ui/oxide/content.css';
import contentUiCss from 'tinymce/skins/ui/oxide/content.css?inline'
import 'tinymce/themes/silver';

// Icon
import 'tinymce/icons/default';

// 用到的外掛
import 'tinymce/plugins/advlist';
import 'tinymce/plugins/lists';
import 'tinymce/plugins/anchor';
import 'tinymce/plugins/autolink';
import 'tinymce/plugins/autoresize';
import 'tinymce/plugins/autosave';
import 'tinymce/plugins/charmap';
import 'tinymce/plugins/code';
import 'tinymce/plugins/fullscreen';
import 'tinymce/plugins/image';
import 'tinymce/plugins/importcss';
import 'tinymce/plugins/table';
import 'tinymce/plugins/quickbars';

// 語言包
// import 'tinymce-i18n/langs5/zh_TW.js';
import 'tinymce-i18n/langs/en_CA.js' //(TinyMCE 6 的簡體中文)

// TinyMCE-Vue
import Editor from '@tinymce/tinymce-vue';

const props = defineProps({
  modelValue: {
    type: String,
    default: '',
  },
  plugins: {
    type: [String, Array],
    default: 'advlist table lists autolink autoresize autosave charmap code fullscreen image importcss',
  },
  toolbar1: {
    type: [String, Array],
    default: 'forecolor backcolor removeformat | table | fontfamily fontsize blocks | alignleft aligncenter alignright alignjustify | bullist numlist | outdent indent | axupimgs | bold italic underline strikethrough ',
  },
});

const emit = defineEmits(['update:modelValue']);

const useDarkMode = false;

const init = reactive({
  language: 'en_CA',
  height: 500,
  width: '100%',
  menubar: false,
  content_css: false,
  content_style: contentUiCss.toString(),
  // skin: useDarkMode ? 'oxide-dark' : 'oxide',
  skin: false,
  plugins: props.plugins,
  toolbar1: props.toolbar1,
  toolbar_mode: 'sliding',
  quickbars_insert_toolbar: true,
  branding: false,
});

const { modelValue } = toRefs(props);
const editorValue = ref(modelValue.value);

watch(modelValue, (newValue) => {
  editorValue.value = newValue;
});

watch(editorValue, (newValue) => {
  emit('update:modelValue', newValue);
});
</script>
