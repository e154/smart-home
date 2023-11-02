<script setup lang="ts">
import {defineEmits, ref, unref} from 'vue'
import {useI18n} from '@/hooks/web/useI18n'
import {ElButton, ElMessage} from 'element-plus'
import {useForm} from '@/hooks/web/useForm'
import {useCache} from '@/hooks/web/useCache'
import {useRoute, useRouter} from 'vue-router'
import {useValidator} from '@/hooks/web/useValidator'
import api from "@/api/api";
import Form from './components/Form.vue'
import {ApiNewTriggerRequest, ApiTrigger} from "@/api/stub";
import ContentWrap from "@/components/ContentWrap/src/ContentWrap.vue";
import TriggerForm from "@/views/Automation/components/TriggerForm.vue";

const {register, elFormRef, methods} = useForm()
const {required} = useValidator()
const {currentRoute, addRoute, push} = useRouter()
const route = useRoute();
const {wsCache} = useCache()
const {t} = useI18n()

const writeRef = ref<ComponentRef<typeof Form>>()
const loading = ref(false)
const currentRow = ref<Nullable<ApiTrigger>>(null)

const save = async () => {
  const write = unref(writeRef)
  const validate = await write?.elFormRef?.validate()?.catch(() => {
  })
  if (validate) {
    loading.value = true
    const tr = (await write?.getFormData()) as ApiTrigger;
    let data = {
      name: tr.name,
      description: tr.description,
      entityId: tr.entity?.id || null,
      scriptId: tr.script?.id || null,
      areaId: tr.area?.id || null,
      pluginName: tr.pluginName,
      attributes: {},
      enabled: tr.enabled,
    } as ApiNewTriggerRequest
    if (tr.pluginName === 'time') {
      data.attributes['cron'] = {
        string: tr?.timePluginOptions || '',
        type: "STRING",
      }
    }
    if (tr.pluginName === 'system') {
      data.attributes['system'] = {
        string: tr?.systemPluginOptions || '',
        type: "STRING",
      }
    }
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
