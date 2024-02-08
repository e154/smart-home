<script setup lang="ts">
import {computed, PropType, ref} from 'vue'
import {useI18n} from '@/hooks/web/useI18n'
import {ElCol, ElForm, ElFormItem, ElInput, ElMessage, ElRadioButton, ElRadioGroup, ElRow} from 'element-plus'
import {ApiVariable} from "@/api/stub";
import {JsonEditor} from "@/components/JsonEditor";
import {TinycmeEditor} from "@/components/Tinymce";
import {propTypes} from "@/utils/propTypes";

const {t} = useI18n()

const props = defineProps({
  modelValue: {
    type: Object as PropType<ApiVariable>,
  },
  edit: propTypes.bool.def(false),
})

const currentVariable = computed(() => props.modelValue)
const jsonValue = ref()

const editorType = ref('Plain')

const onChangeView = (val: string): void => {
  switch (val) {
    case 'Plain':
      break
    case 'Json':
      convertToJsonValue()
      break
    case 'HTML':
      break
  }
}

const onJsonChanged = (val): void => {
  // console.log('onJsonChange', val)
  if (val.text) {
    currentVariable.value.value = val.text.replace(/\n/g, '') || ''
    return
  }
}

const onHtmlChanged = (val: string): void => {
  // console.log('onHtmlChanged', val)
  currentVariable.value.value = val
}

const convertToJsonValue = () => {
  try {
    jsonValue.value = JSON.parse(currentVariable.value?.value) || {}
  } catch (e) {
    jsonValue.value = {}
    ElMessage({
      title: t('Error'),
      message: t('message.corruptedJsonFormat') + e,
      type: 'error',
      duration: 2000
    });
    return {}
  }
}

</script>

<template>
  <ElForm
      label-position="top"
      :model="currentVariable"
      style="width: 100%"
      ref="cardItemForm"
  >

    <ElRow>
      <ElCol>
        <ElFormItem :label="$t('variables.name')" prop="text">
          <ElInput
              :disabled="edit"
              placeholder="Please input"
              v-model="currentVariable.name"
          />
        </ElFormItem>
      </ElCol>
    </ElRow>

    <ElRow>
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('variables.editor')" prop="text">
          <ElRadioGroup v-model="editorType" @change="onChangeView">
            <ElRadioButton label="Plain"/>
            <ElRadioButton label="Json"/>
            <ElRadioButton label="HTML"/>
          </ElRadioGroup>
        </ElFormItem>
      </ElCol>
    </ElRow>

    <ElRow v-if="editorType === 'Plain'">
      <ElCol>
        <ElFormItem :label="$t('variables.value')" prop="text">
          <ElInput
              type="textarea"
              :autosize="{minRows: 10}"
              placeholder="Please input"
              v-model="currentVariable.value"
          />
        </ElFormItem>
      </ElCol>
    </ElRow>

    <ElRow v-else-if="editorType === 'Json'">
      <ElCol>
        <ElFormItem :label="$t('variable.value')" prop="text">
          <JsonEditor height="auto" v-if="jsonValue" v-model="jsonValue" @change="onJsonChanged"/>
        </ElFormItem>
      </ElCol>
    </ElRow>

    <ElRow v-else-if="editorType === 'HTML'">
      <ElCol>
        <ElFormItem :label="$t('variable.value')" prop="text">
          <TinycmeEditor v-model="currentVariable.value" @update:modelValue="onHtmlChanged"/>
        </ElFormItem>
      </ElCol>
    </ElRow>


  </ElForm>
</template>
