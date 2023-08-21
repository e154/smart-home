<script setup lang="ts">
import {computed, reactive, ref, unref} from 'vue'
import {useI18n} from '@/hooks/web/useI18n'
import {ElButton, ElPopconfirm, ElTimeline, ElTimelineItem, ElCard, ElTabs, ElTabPane, ElMessage} from 'element-plus'
import {useForm} from '@/hooks/web/useForm'
import {useCache} from '@/hooks/web/useCache'
import {useAppStore} from '@/store/modules/app'
import {useRoute, useRouter} from 'vue-router'
import {useValidator} from '@/hooks/web/useValidator'
import api from "@/api/api";
import {ApiAction, ApiCondition, ApiTask, ApiTrigger} from "@/api/stub";
import ContentWrap from "@/components/ContentWrap/src/ContentWrap.vue";
import Form from "@/views/Automation/components/Form.vue";
import Triggers from "@/views/Automation/components/Triggers.vue";
import Conditions from "@/views/Automation/components/Conditions.vue";
import Actions from "@/views/Automation/components/Actions.vue";
import {useEmitt} from "@/hooks/web/useEmitt";
import { Dialog } from '@/components/Dialog'
import Viewer from "@/components/JsonViewer/JsonViewer.vue";
import {copyToClipboard} from "@/utils/clipboard";

const {register, elFormRef, methods} = useForm()
const {required} = useValidator()
const appStore = useAppStore()
const {currentRoute, addRoute, push} = useRouter()
const route = useRoute();
const {wsCache} = useCache()
const {t} = useI18n()

const writeRef = ref<ComponentRef<typeof Form>>()
const loading = ref(false)
const taskId = computed(() => route.params.id as number);
const currentTask = ref<Nullable<ApiTask>>(null)
const activeTab = ref('pipeline')
const dialogSource = ref({})
const dialogVisible = ref(false)

const fetch = async () => {
  loading.value = true
  const res = await api.v1.automationServiceGetTask(taskId.value)
      .catch(() => {
      })
      .finally(() => {
        loading.value = false
      })
  if (res) {
    currentTask.value = res.data
  } else {
    currentTask.value = null
  }
}

const prepareForSave = async () => {
  const write = unref(writeRef)
  const validate = await write?.elFormRef?.validate()?.catch((err) => {
  })
  if (validate) {
    const data = (await write?.getFormData()) as ApiTask
    let triggers: ApiTrigger[] = []
    let conditions: ApiCondition[] = []
    let actions: ApiAction[] = []
    for (const tr of data.triggers) {
      triggers.push({
        id: tr.id,
        name: tr.name,
        entityId: tr.entityId,
        scriptId: tr.scriptId,
        pluginName: tr.pluginName,
        attributes: tr.attributes,
      })
    }
    for (const cond of data.conditions) {
      conditions.push({
        id: cond.id,
        name: cond.name,
        scriptId: cond.scriptId,
      })
    }
    for (const action of data.actions) {
      actions.push({
        id: action.id,
        name: action.name,
        entityId: action.entityId,
        scriptId: action.scriptId,
        entityActionName: action.entityActionName,
      })
    }
    return {
      name: data.name,
      description: data.description,
      enabled: data.enabled,
      condition: data.condition,
      triggers: triggers,
      conditions: conditions,
      actions: actions,
      areaId: data.area?.id || null,
    }
  }
  return null
}

const save = async () => {
  const body = await prepareForSave()
  if (!body) {
    return
  }
  const res = await api.v1.automationServiceUpdateTask(taskId.value, body)
      .catch(() => {
      })
      .finally(() => {

      })
  if (res) {
    fetch()
    ElMessage({
      title: t('Success'),
      message: t('message.uploadSuccessfully'),
      type: 'success',
      duration: 2000
    })
  }
}

const remove = async () => {
  loading.value = true
  const res = await api.v1.automationServiceDeleteTask(taskId.value)
      .catch(() => {
      })
      .finally(() => {
        loading.value = false
      })
  if (res) {
    cancel()
  }
}

const cancel = () => {
  push('/automation')
}

const copy = async () => {
  const body = await prepareForSave()
  copyToClipboard(JSON.stringify(body, null, 2))
}

const exportTask = async () => {
  const body = await prepareForSave()
  dialogSource.value = body
  dialogVisible.value = true
}

const callAction = async (name: string) => {
  await api.v1.developerToolsServiceTaskCallAction({ id: taskId.value || 0, name: name })
      .catch(() => {
      })
      .finally(() => {
        ElMessage({
          title: t('Success'),
          message: t('message.callSuccessful'),
          type: 'success',
          duration: 2000
        })
      })
}

useEmitt({
  name: 'callAction',
  callback: (name: string) => {
    callTrigger(name)
  }
})

useEmitt({
  name: 'updateTriggers',
  callback: (val: ApiTrigger[]) => {
  }
})

useEmitt({
  name: 'updateConditions',
  callback: (val: ApiCondition[]) => {
  }
})

const callTrigger = async (name: string) => {
  await api.v1.developerToolsServiceTaskCallTrigger({ id: taskId.value || 0, name: name })
      .catch(() => {
      })
      .finally(() => {
        ElMessage({
          title: t('Success'),
          message: t('message.callSuccessful'),
          type: 'success',
          duration: 2000
        })
      })
}

useEmitt({
  name: 'callTrigger',
  callback: (name: string) => {
    callTrigger(name)
  }
})



fetch()

</script>

<template>
  <ContentWrap>

    <el-tabs class="demo-tabs" v-model="activeTab">
      <!-- main -->
      <el-tab-pane :label="$t('automation.main')" name="main">
        <Form ref="writeRef" :current-task="currentTask"/>
      </el-tab-pane>
      <!-- /main -->

      <!-- pipeline -->
      <el-tab-pane :label="$t('automation.pipeline')" name="pipeline" class="mt-20px">
        <el-timeline>
          <el-timeline-item :timestamp="$t('automation.eventStart')" placement="top" type="primary"/>

          <el-timeline-item :timestamp="$t('automation.triggers')" placement="top" type="primary" center>
            <el-card>
              <Triggers :triggers="currentTask?.triggers"/>
            </el-card>
          </el-timeline-item>
          <el-timeline-item :timestamp="$t('automation.conditions')" placement="top" type="primary" center>
            <el-card>
              <Conditions :conditions="currentTask?.conditions"/>
            </el-card>
          </el-timeline-item>
          <el-timeline-item :timestamp="$t('automation.actions')" placement="top" type="primary" center>
            <el-card>
              <Actions :actions="currentTask?.actions"/>
            </el-card>
          </el-timeline-item>
          <el-timeline-item :timestamp="$t('automation.eventEnd')" placement="top" type="success"/>
        </el-timeline>
      </el-tab-pane>
      <!-- /pipeline -->
    </el-tabs>

    <div style="text-align: right">
      <ElButton type="primary" @click="save()">
        {{ t('main.save') }}
      </ElButton>

      <ElButton type="primary" @click="exportTask()">
        <Icon icon="uil:file-export" class="mr-5px"/>
        {{ t('main.export') }}
      </ElButton>

      <ElButton type="default" @click="fetch()">
        {{ t('main.loadFromServer') }}
      </ElButton>

      <ElButton type="default" @click="cancel()">
        {{ t('main.return') }}
      </ElButton>

      <ElPopconfirm
          :confirm-button-text="$t('main.ok')"
          :cancel-button-text="$t('main.no')"
          width="250"
          style="margin-left: 10px;"
          :title="$t('main.are_you_sure_to_do_want_this?')"
          @confirm="remove"
      >
        <template #reference>
          <ElButton class="mr-10px" type="danger" plain>
            <Icon icon="ep:delete" class="mr-5px"/>
            {{ t('main.remove') }}
          </ElButton>
        </template>
      </ElPopconfirm>

    </div>

  </ContentWrap>

  <!-- export dialog -->
  <Dialog v-model="dialogVisible" :title="t('automation.dialogTitle')" :maxHeight="400">
    <Viewer v-model="dialogSource"/>
    <template #footer>
      <ElButton @click="copy()">{{ t('setting.copy') }}</ElButton>
      <ElButton @click="dialogVisible = false">{{ t('main.closeDialog') }}</ElButton>
    </template>
  </Dialog>
  <!-- /export dialog -->

</template>

<style lang="less" scoped>

</style>
