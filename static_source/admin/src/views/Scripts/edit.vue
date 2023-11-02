<script setup lang="ts">
import {computed, nextTick, ref, unref} from 'vue'
import {useI18n} from '@/hooks/web/useI18n'
import {ElButton, ElPopconfirm, ElDescriptions, ElDescriptionsItem, ElTabs, ElTabPane, ElMessage,
  ElTimeline, ElTimelineItem, ElCard, ElRow, ElCol, ElScrollbar} from 'element-plus'
import {useForm} from '@/hooks/web/useForm'
import {useCache} from '@/hooks/web/useCache'
import {useAppStore} from '@/store/modules/app'
import {usePermissionStore} from '@/store/modules/permission'
import {useRoute, useRouter} from 'vue-router'
import {useValidator} from '@/hooks/web/useValidator'
import api from "@/api/api";
import Form from './components/Form.vue'
import {ApiScript, ApiScriptVersion} from "@/api/stub";
import ScriptEditor from "@/views/Scripts/components/ScriptEditor.vue";
import ContentWrap from "@/components/ContentWrap/src/ContentWrap.vue";
import {useEmitt} from "@/hooks/web/useEmitt";
import { Infotip } from '@/components/Infotip'

const {register, elFormRef, methods} = useForm()
const {required} = useValidator()
const appStore = useAppStore()
const permissionStore = usePermissionStore()
const {currentRoute, addRoute, push} = useRouter()
const route = useRoute();
const {wsCache} = useCache()
const {t} = useI18n()

const writeRef = ref<ComponentRef<typeof Form>>()
const loading = ref(false)
const scriptId = computed(() => route.params.id as number);
const currentScript = ref<Nullable<ApiScript>>(null)
const currentScriptVersion = ref<Nullable<ApiScript>>(null)
const activeTab = ref('source')
const sourceScript = ref('')
import {parseTime} from "@/utils";
const { emitter } = useEmitt()

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

const copy = async () => {
  const res = await api.v1.scriptServiceCopyScriptById(scriptId.value)
      .catch(() => {
      })
      .finally(() => {

      })
  if (res) {
    const {id} = res.data;
    push(`/scripts/edit/${id}`)
  }
}

const save = async () => {
  const write = unref(writeRef)
  const validate = await write?.elFormRef?.validate()?.catch(() => {
  })
  if (validate) {
    loading.value = true
    const data = (await write?.getFormData()) as ApiScript
    const body = {
      id: data.id,
      lang: data.lang,
      name: data.name,
      source: sourceScript.value,
      description: data.description,
    }
    const res = await api.v1.scriptServiceUpdateScriptById(scriptId.value, body)
        .catch(() => {
        })
        .finally(() => {
          loading.value = false
        })
    if (res) {
      ElMessage({
        title: t('Success'),
        message: t('message.updatedSuccessfully'),
        type: 'success',
        duration: 2000
      })
    }
  }
}

const cancel = () => {
  push('/scripts')
}

const remove = async () => {
  loading.value = true
  const res = await api.v1.scriptServiceDeleteScriptById(scriptId.value)
      .catch(() => {
      })
      .finally(() => {
        loading.value = false
      })
  if (res) {
    cancel()
  }
}

const view = async (version: ApiScriptVersion) => {
  currentScriptVersion.value = {
    lang: version.lang,
    source: version.source,
    createdAt: version.createdAt,
  } as ApiScript
}

const rollback = async () => {
  let script = unref(currentScript.value) as ApiScript;
  currentScript.value = {
    id: script.id,
    name: script.name,
    description: script.description,
    scriptInfo: script.scriptInfo,
    versions: script.versions,
    lang: currentScriptVersion.value?.lang || script.lang,
    source: currentScriptVersion.value?.source || script.source,
  }  as ApiScript;
}

useEmitt({
  name: 'updateSource',
  callback: (val: string) => {
    if (sourceScript.value == val) {
      return
    }
    sourceScript.value = val
  }
})

const updateCurrentTab = (tab: any, ev: any) => {
  const {index, paneName} = tab;
  if (paneName == 'source' || paneName == 'versions') {
    emitter.emit('updateEditor')
  }
}

// elscroll 实例
const scrollbarRef = ref<ComponentRef<typeof ElScrollbar>>()


fetch()

</script>

<template>
  <ContentWrap>

    <ElTabs class="demo-tabs" v-model="activeTab" @tab-click="updateCurrentTab">
      <!-- main -->
      <ElTabPane :label="$t('scripts.main')" name="main">

        <Form ref="writeRef" :current-row="currentScript"/>

      </ElTabPane>
      <!-- /main -->

      <!-- source -->
      <ElTabPane :label="$t('scripts.source')" name="source">
        <Infotip
            :show-index="false"
            title="INFO"
            :schema="[
      {
        label: t('scripts.info1'),
        keys: ['scripts']
      },
      {
        label: t('scripts.info2'),
        keys: ['scripts']
      },
      {
        label: t('scripts.info3'),
        keys: ['scripts']
      },
    ]"
        />
        <ScriptEditor v-model="currentScript" class="mb-20px"/>
      </ElTabPane>
      <!-- /source -->

      <!-- info -->
      <ElTabPane :label="$t('scripts.scriptInfo')" name="info">
        <ElDescriptions v-if="currentScript?.scriptInfo"
                        class="ml-10px mr-10px mb-20px"
                        :title="$t('scripts.scriptInfo')"
                        direction="vertical"
                        :column="2"
                        border
        >
          <ElDescriptionsItem :label="$t('scripts.alexaIntents')">{{currentScript.scriptInfo.alexaIntents}}</ElDescriptionsItem>
          <ElDescriptionsItem :label="$t('scripts.entityActions')">{{currentScript.scriptInfo.entityActions}}</ElDescriptionsItem>
          <ElDescriptionsItem :label="$t('scripts.entityScripts')">{{currentScript.scriptInfo.entityScripts}}</ElDescriptionsItem>
          <ElDescriptionsItem :label="$t('scripts.automationTriggers')">{{currentScript.scriptInfo.automationTriggers}}</ElDescriptionsItem>
          <ElDescriptionsItem :label="$t('scripts.automationConditions')">{{currentScript.scriptInfo.automationConditions}}</ElDescriptionsItem>
          <ElDescriptionsItem :label="$t('scripts.automationActions')">{{currentScript.scriptInfo.automationActions}}</ElDescriptionsItem>
        </ElDescriptions>
      </ElTabPane>
      <!-- /info -->

      <!-- versions -->
      <ElTabPane :label="$t('scripts.scriptVersions')" name="versions">

        <ElRow :gutter="24" class="mb-20px">
          <ElCol :span="6" :xs="6">
            <ElCard shadow="never" class="item-card-editor">
              <ElScrollbar ref="scrollbarRef" class="h-full" height="500px">
                <ElTimeline v-if="currentScript && currentScript?.versions">
                  <ElTimelineItem
                      v-for="(version, index) in currentScript?.versions"
                      :key="index"
                      :timestamp="parseTime(version.createdAt)"
                      type="primary"
                      placement="top"
                      class="cursor-pointer"
                      @click="view(version)"
                  />
                </ElTimeline>
              </ElScrollbar>
            </ElCard>
          </ElCol>
          <ElCol :span="18" style="padding-bottom: 30px">

            <div v-if="currentScriptVersion">
                <ScriptEditor v-model="currentScriptVersion"/>
            </div>

            <div v-if="currentScriptVersion">
                <ElButton class="mr-10px left" type="default" @click="rollback(version)">
                  <Icon icon="ic:baseline-restore" class="mr-5px"/>
                  {{ $t('scripts.restoreVersion') }}
                </ElButton>
            </div>

          </ElCol>
        </ElRow>

      </ElTabPane>
      <!-- /versions -->
    </ElTabs>

    <div style="text-align: right">

      <ElButton type="success" @click="exec()">{{ t('main.exec') }}</ElButton>
      <ElButton type="primary" @click="save()">{{ t('main.save') }}</ElButton>
      <ElButton type="default" @click="copy()">{{ t('main.copy') }}</ElButton>
      <ElButton type="default" @click="fetch()">{{ t('main.loadFromServer') }}</ElButton>
      <ElButton type="default" @click="cancel()">{{ t('main.return') }}</ElButton>

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

</template>

<style lang="less" scoped>

</style>
