<script setup lang="ts">
import {computed, PropType, ref} from 'vue'
import {useI18n} from '@/hooks/web/useI18n'
import {
  ElCol,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElRadioButton,
  ElRadioGroup,
  ElRow,
  ElUpload,
  UploadFile
} from 'element-plus'
import {ApiVariable} from "@/api/stub";
import {JsonEditor} from "@/components/JsonEditor";
import {TinycmeEditor} from "@/components/Tinymce";
import {propTypes} from "@/utils/propTypes";
import {BaseButton} from "@/components/Button";
import {TagsSearch} from "@/components/TagsSearch";

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

const uploadChange = (uploadFile: UploadFile) => {
  if (!uploadFile.raw) return

  const reader = new FileReader();
  reader.onloadend = () => {
    // Use a regex to remove data url part
    const base64String = reader.result
      .replace('data:', '')
      .replace(/^.+,/, '');

    // console.log(base64String);
    currentVariable.value.value = base64String
  };
  reader.readAsDataURL(uploadFile.raw);
}

const changedTags = async (tags: string[]) => {
  // console.log(tags)
  currentVariable.value.tags = tags
}

</script>

<template>
  <ElForm
    label-position="top"
    :model="currentVariable"
    style="width: 100%"
    ref="cardItemForm"
  >

    <ElRow :gutter="24">
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('variables.name')" prop="text">
          <ElInput
            :disabled="edit"
            placeholder="Please input"
            v-model="currentVariable.name"
          />
        </ElFormItem>
      </ElCol>
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('entityAction.tags')" prop="tags">
          <TagsSearch v-model="currentVariable.tags" @change="changedTags($event)"/>
        </ElFormItem>
      </ElCol>
    </ElRow>

    <ElRow>
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('main.upload')" prop="text">
          <ElUpload
            action="''"
            :auto-upload="false"
            :show-file-list="false"
            :on-change="uploadChange"
          >
            <BaseButton size="small" type="primary">
              <Icon icon="ep:upload-filled"/>
            </BaseButton>
          </ElUpload>
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
            :autosize="{minRows: 10, maxRows: 10}"
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
