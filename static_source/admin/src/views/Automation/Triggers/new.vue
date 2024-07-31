<script setup lang="ts">
import {ref, unref} from 'vue'
import {useI18n} from '@/hooks/web/useI18n'
import {ElButton, ElMessage} from 'element-plus'
import {useRouter} from 'vue-router'
import {ApiAttribute, ApiNewTriggerRequest} from "@/api/stub";
import {ContentWrap} from "@/components/ContentWrap";
import TriggerForm from "@/views/Automation/components/TriggerForm.vue";
import api from "@/api/api";

const {push} = useRouter()
const {t} = useI18n()

const writeRef = ref<ComponentRef<typeof TriggerForm>>()
const loading = ref(false)
const currentRow = ref({
  name: '',
  description: '',
  script: null,
  scriptId: null,
  attributes: [],
} as ApiNewTriggerRequest)

const save = async () => {
  const write = unref(writeRef)
  const validate = await write?.form?.validate()?.catch(() => {
  })
  if (!validate) {
    return
  }

  const tr = unref(currentRow)
  let data = {
    name: tr.name,
    description: tr.description,
    entityIds: tr.entityIds || [],
    scriptId: tr.script?.id || null,
    areaId: tr.area?.id || null,
    pluginName: tr.pluginName,
    attributes: {},
    enabled: tr.enabled,
  } as ApiNewTriggerRequest

  let attributes: { [key: string]: ApiAttribute } = {};
  for (const attr of tr.attributes) {
    if (attr.name == 'notice') {
      continue
    }
    attributes[attr.name] = attr;
  }
  data.attributes = attributes

  const res = await api.v1.triggerServiceAddTrigger(data)
      .catch(() => {
      })
      .finally(() => {
        loading.value = false
      })
  if (res) {
    ElMessage({
      title: t('Success'),
      message: t('message.createdSuccessfully'),
      type: 'success',
      duration: 2000
    })

    cancel()
  }
}

const cancel = () => {
  push('/automation/triggers')
}

</script>

<template>
  <ContentWrap>
    <TriggerForm ref="writeRef" :trigger="currentRow"/>

    <div style="text-align: right">

      <ElButton type="primary" @click="save()">
        {{ t('main.save') }}
      </ElButton>

      <ElButton type="default" @click="cancel()">
        {{ t('main.cancel') }}
      </ElButton>

    </div>
  </ContentWrap>

</template>

<style lang="less" scoped>

</style>
