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
import {useEmitt} from "@/hooks/web/useEmitt";
import Viewer from "@/components/JsonViewer/JsonViewer.vue";
import {copyToClipboard} from "@/utils/clipboard";
import {Form} from "@/components/Form";
import TaskForm from "@/views/Automation/components/TaskForm.vue";
import TriggersSearch from "@/views/Automation/components/TriggersSearch.vue";
import ConditionsSearch from "@/views/Automation/components/ConditionsSearch.vue";
import ActionsSearch from "@/views/Automation/components/ActionsSearch.vue";
import { Dialog } from '@/components/Dialog'
import TaskTelemetry from "@/views/Automation/components/TaskTelemetry.vue";

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
const triggerIds = ref([])
const conditionIds = ref([])
const actionIds = ref([])

const fetch = async () => {
  loading.value = true
  const res = await api.v1.automationServiceGetTask(taskId.value)
      .catch(() => {
      })
      .finally(() => {
        loading.value = false
      })
  if (res) {
    let task = res.data as ApiTask
    triggerIds.value = Object.assign([], task.triggerIds);
    conditionIds.value = Object.assign([], task.conditionIds);
    actionIds.value = Object.assign([], task.actionIds);
    currentTask.value = task
  } else {
    currentTask.value = null
  }
}

const prepareForSave = async () => {
  const write = unref(writeRef)
  const validate = await write?.elFormRef?.validate()?.catch(() => {
  })
  const con = (await write?.getFormData()) as ApiTask;
  if (validate) {
    return {
      name: con.name,
      description: con.description,
      enabled: con.enabled,
      condition: con.condition,
      triggerIds: triggerIds.value ,
      conditionIds: conditionIds.value,
      actionIds: actionIds.value,
      areaId: con.area?.id,
    }
  }
}

const save = async () => {
  const body = await prepareForSave()
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
  push('/automation/tasks')
}

const copy = async () => {
  const body = prepareForSave()
  copyToClipboard(JSON.stringify(body, null, 2))
}

const exportTask = async () => {
  const body = prepareForSave()
  dialogSource.value = body
  dialogVisible.value = true
}

const callAction = async (name: string) => {
  // todo: fix
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
  //todo: fix
  await api.v1.developerToolsServiceCallTrigger({ id: taskId.value || 0, name: name })
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
        <TaskForm ref="writeRef" :task="currentTask"/>
      </el-tab-pane>
      <!-- /main -->

      <!-- telemetry -->
      <el-tab-pane :label="$t('automation.telemetry')" name="telemetry" :lazy="true">
        <TaskTelemetry :task="currentTask"/>
      </el-tab-pane>
      <!-- /telemetry -->

      <!-- pipeline -->
      <el-tab-pane :label="$t('automation.tasks.pipeline')" name="pipeline" class="mt-20px">
        <el-timeline>
          <el-timeline-item :timestamp="$t('automation.tasks.eventStart')" placement="top" type="primary"/>

          <el-timeline-item :timestamp="$t('automation.tasks.triggers')" placement="top" type="primary" center>
            <el-card>
              <TriggersSearch v-model="triggerIds"/>
            </el-card>
          </el-timeline-item>
          <el-timeline-item :timestamp="$t('automation.tasks.conditions')" placement="top" type="primary" center>
            <el-card>
              <ConditionsSearch v-model="conditionIds"/>
            </el-card>
          </el-timeline-item>
          <el-timeline-item :timestamp="$t('automation.tasks.actions')" placement="top" type="primary" center>
            <el-card>
              <ActionsSearch v-model="actionIds"/>
            </el-card>
          </el-timeline-item>
          <el-timeline-item :timestamp="$t('automation.tasks.eventEnd')" placement="top" type="success"/>
        </el-timeline>
      </el-tab-pane>
      <!-- /pipeline -->
    </el-tabs>

    <div style="text-align: right">
      <ElButton type="primary" @click="save()">
        {{ t('main.save') }}
      </ElButton>

<!--      <ElButton type="primary" @click="exportTask()">-->
<!--        <Icon icon="uil:file-export" class="mr-5px"/>-->
<!--        {{ t('main.export') }}-->
<!--      </ElButton>-->

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
  <Dialog v-model="dialogVisible" :title="t('automation.tasks.exportDialogTitle')" :maxHeight="400">
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
