<script setup lang="ts">
import {defineEmits, ref, unref} from 'vue'
import {useI18n} from '@/hooks/web/useI18n'
import {ElButton} from 'element-plus'
import {useForm} from '@/hooks/web/useForm'
import {useCache} from '@/hooks/web/useCache'
import {useAppStore} from '@/store/modules/app'
import {usePermissionStore} from '@/store/modules/permission'
import {useRoute, useRouter} from 'vue-router'
import {useValidator} from '@/hooks/web/useValidator'
import api from "@/api/api";
import {ApiNewTaskRequest, ApiTask} from "@/api/stub";
import ContentWrap from "@/components/ContentWrap/src/ContentWrap.vue";
import TaskForm from "@/views/Automation/components/TaskForm.vue";
import {Form} from "@/components/Form";

const {register, elFormRef, methods} = useForm()
const {required} = useValidator()
const emit = defineEmits(['to-restore'])
const appStore = useAppStore()
const permissionStore = usePermissionStore()
const {currentRoute, addRoute, push} = useRouter()
const route = useRoute();
const {wsCache} = useCache()
const {t} = useI18n()

const writeRef = ref<ComponentRef<typeof Form>>()
const loading = ref(false)
const currentRow = ref<Nullable<ApiTask>>(null)

const save = async () => {
  const write = unref(writeRef)
  const validate = await write?.elFormRef?.validate()?.catch(() => {
  })
  if (validate) {
    loading.value = true
    const data = (await write?.getFormData())
    console.log(data)
    const body = {
      name: data.name,
      description: data.description,
      enabled: data.enabled,
      condition: data.condition,
      triggers: data.triggers,
      conditions: data.conditions,
      actions: data.actions,
      areaId: data.area?.id || null,
    } as ApiNewTaskRequest
    const res = await api.v1.automationServiceAddTask(body)
        .catch(() => {
        })
        .finally(() => {
          loading.value = false
        })
    if (res) {
      const {id} = res.data
      push(`/automation/edit/${id}`)
    }
  }
}

const cancel = () => {
  push('/automation')
}

</script>

<template>
  <ContentWrap>
    <TaskForm ref="writeRef" :current-row="currentRow"/>

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
