<script setup lang="ts">
import {defineEmits, ref, unref} from 'vue'
import {useI18n} from '@/hooks/web/useI18n'
import {ElButton, ElMessage} from 'element-plus'
import {useForm} from '@/hooks/web/useForm'
import {useCache} from '@/hooks/web/useCache'
import {useAppStore} from '@/store/modules/app'
import {usePermissionStore} from '@/store/modules/permission'
import {useRoute, useRouter} from 'vue-router'
import {useValidator} from '@/hooks/web/useValidator'
import api from "@/api/api";
import Form from './components/Form.vue'
import {ApiNewActionRequest, ApiAction} from "@/api/stub";
import ContentWrap from "@/components/ContentWrap/src/ContentWrap.vue";
import ActionForm from "@/views/Automation/components/ActionForm.vue";

const {register, elFormRef, methods} = useForm()
const {required} = useValidator()
const {currentRoute, addRoute, push} = useRouter()
const route = useRoute();
const {wsCache} = useCache()
const {t} = useI18n()

const writeRef = ref<ComponentRef<typeof Form>>()
const loading = ref(false)
const currentRow = ref<Nullable<ApiAction>>(null)

const save = async () => {
  const write = unref(writeRef)
  const validate = await write?.elFormRef?.validate()?.catch(() => {
  })
  if (validate) {
    loading.value = true
    const act = (await write?.getFormData()) as ApiAction;
    const data = {
      name: act.name,
      scriptId: act.script?.id,
      entityId: act.entity?.id,
      entityActionName: act.entityActionName,
    } as ApiNewActionRequest
    const res = await api.v1.actionServiceAddAction(data)
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
  push('/automation/actions')
}

</script>

<template>
  <ContentWrap>
    <ActionForm ref="writeRef" :action="currentRow"/>

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
