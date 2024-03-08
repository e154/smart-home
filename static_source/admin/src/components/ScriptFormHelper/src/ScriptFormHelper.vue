<script setup lang="ts">

import {computed, PropType, ref} from "vue";
import {ApiScript} from "@/api/stub";
import {ElButton, ElMessage} from 'element-plus'
import api from "@/api/api";
import {Dialog} from '@/components/Dialog'
import {ScriptEditor} from "@/components/ScriptEditor";
import {useI18n} from "@/hooks/web/useI18n";
import {useEmitt} from "@/hooks/web/useEmitt";

const {t} = useI18n()

const currentScript = ref<Nullable<ApiScript>>(null)
const emit = defineEmits(['change', 'update:modelValue'])

const props = defineProps({
  modelValue: {
    type: Object as PropType<Nullable<ApiScript>>,
    default: () => undefined
  }
})

const scriptId = computed(() => props.modelValue?.id as number);


const loading = ref(true)
const fetch = async () => {
  loading.value = true
  const res = await api.v1.scriptServiceGetScriptById(scriptId.value)
      .catch(() => {
      })
      .finally(() => {
        loading.value = false
      })
  if (res) {
    currentScript.value = res.data
  } else {
    currentScript.value = null
  }
}

const save = async () => {
  const body = {
    lang: currentScript.value?.lang,
    name: currentScript.value?.name,
    source: sourceScript.value,
    description: currentScript.value?.description,
  }
  const res = await api.v1.scriptServiceUpdateScriptById(scriptId.value, body)
      .catch(() => {
      })
  if (res) {
    currentScript.value = res.data as ApiScript
    dialogSource.value = currentScript.value.source
    ElMessage({
      title: t('Success'),
      message: t('message.updatedSuccessfully'),
      type: 'success',
      duration: 2000
    })
  }
}

const dialogSource = ref({})
const dialogVisible = ref(false)
const showModalDialog = async () => {
  await fetch()
  dialogSource.value = currentScript.value
  dialogVisible.value = true
}


const sourceScript = ref('')
useEmitt({
  name: 'updateSource',
  callback: (val: string) => {
    if (sourceScript.value == val) {
      return
    }
    sourceScript.value = val
  }
})

const exec = async () => {
  await api.v1.scriptServiceExecSrcScriptById({
    name: currentScript.value?.name,
    source: sourceScript.value,
    lang: currentScript.value?.lang
  })
  ElMessage({
    title: t('Success'),
    message: t('message.callSuccessful'),
    type: 'success',
    duration: 2000
  })
}


</script>

<template>
  <div class="tip custom-block w-[100%]" v-if="modelValue">
    <ElButton @click="showModalDialog()" :link="true">
      <Icon icon="uil:file-export" class="mr-5px"/>
      {{ t('scripts.showModalDialog') }}
    </ElButton>
  </div>

  <!-- show dialog -->
  <Dialog v-model="dialogVisible" :title="t('scripts.modalWindow')" :maxHeight="400" width="80%">
    <ScriptEditor v-if="!loading" v-model="currentScript" @save="save"/>
    <template #footer>
      <ElButton type="success" @click="exec()">{{ t('main.exec') }}</ElButton>
      <ElButton @click="save()">{{ t('main.update') }}</ElButton>
      <ElButton type="default" @click="fetch()">{{ t('main.loadFromServer') }}</ElButton>
      <ElButton @click="dialogVisible = false">{{ t('main.closeDialog') }}</ElButton>
    </template>
  </Dialog>
  <!-- /show dialog -->

</template>

<style lang="less" scoped>
.light {
  .custom-block.tip {
    padding: 8px 16px;
    background-color: #409eff1a;
    border-radius: 4px;
    border-left: 5px solid var(--el-color-primary);
  }
}

.dark {
  .custom-block.tip {
    padding: 8px 16px;
    background-color: #409eff1a;
    border-radius: 4px;
    border-left: 5px solid var(--el-color-primary);
  }
}

</style>
