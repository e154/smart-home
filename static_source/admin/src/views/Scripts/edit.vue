<script setup lang="ts">
import {computed, onMounted, onUnmounted, ref, unref} from 'vue'
import {useI18n} from '@/hooks/web/useI18n'
import {
  ElButton,
  ElCol,
  ElDescriptions,
  ElDescriptionsItem,
  ElFormItem,
  ElMessage,
  ElOption,
  ElPopconfirm,
  ElRow,
  ElSelect,
  ElTabPane,
  ElTabs
} from 'element-plus'
import {useRoute, useRouter} from 'vue-router'
import api from "@/api/api";
import Form from './components/Form.vue'
import {ApiScript} from "@/api/stub";
import {ScriptEditor} from "@/components/ScriptEditor";
import ContentWrap from "@/components/ContentWrap/src/ContentWrap.vue";
import {useEmitt} from "@/hooks/web/useEmitt";
import {Infotip} from '@/components/Infotip'
import {parseTime} from "@/utils";
import {UUID} from "uuid-generator-ts";
import stream from "@/api/stream";
import {MergeEditor} from "@/components/MergeEditor";

const {emitter} = useEmitt()
const {push} = useRouter()
const route = useRoute();
const {t} = useI18n()

const writeRef = ref<ComponentRef<typeof Form>>()
const loading = ref(false)
const scriptId = computed(() => route.params.id as number);
const currentScript = ref<Nullable<ApiScript>>(null)
const activeTab = ref('source')
const currentVersionIdx = ref(0)
const currentVersion = ref<Nullable<ApiScript>>(null)
const versions = ref<Nullable<ApiScript[]>>([])

const currentID = ref('')

onMounted(() => {
  const uuid = new UUID()
  currentID.value = uuid.getDashFreeUUID()

  setTimeout(() => {
    stream.subscribe('event_updated_script_model', currentID.value, fetch)
  }, 1000)
})

onUnmounted(() => {
  stream.unsubscribe('event_updated_script_model', currentID.value)
})

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
    if (res.data?.versions && res.data.versions.length) {
      versions.value = res.data.versions;
      if (res.data.versions.length > 1) {
        currentVersion.value = res.data.versions[1]
      }
    }
  } else {
    currentScript.value = null
  }
}

const exec = async () => {
  await api.v1.scriptServiceExecSrcScriptById({
    name: currentScript.value?.name,
    source: currentScript.value?.source,
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

const updateVersions = async () => {
  loading.value = true
  const res = await api.v1.scriptServiceGetScriptById(scriptId.value)
      .catch(() => {
      })
      .finally(() => {
        loading.value = false
      })
  if (res) {
    if (res.data?.versions && res.data.versions.length) {
      currentVersion.value = res.data.versions[0]
      versions.value = res.data.versions
    }
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
      source: currentScript.value?.source || '',
      description: data.description,
    }
    const res = await api.v1.scriptServiceUpdateScriptById(scriptId.value, body)
        .catch(() => {
        })
        .finally(async () => {
          updateVersions();
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

const updateCurrentTab = (tab: any, ev: any) => {
  const {index, paneName} = tab;
  if (paneName == 'source' || paneName == 'versions') {
    emitter.emit('updateEditor')
  }
}

const onMergeEditorChange = (val: string) => {
  currentScript.value.source = val
  reload()
}

const onScriptEditorChange = (val: string) => {
  if (currentScript.value?.source == val) {
    return
  }
  currentScript.value.source = val
}

const selectVersionHandler = () => {
  currentVersion.value = currentScript.value?.versions[currentVersionIdx.value] || null
}

const reloadKey = ref(0)
const reload = () => reloadKey.value += 1

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
        <ScriptEditor :key="reloadKey"
                      v-model="currentScript"
                      class="mb-20px"
                      @update:source="onScriptEditorChange"
                      @save="save"/>
      </ElTabPane>
      <!-- /source -->

      <!-- versions -->
      <ElTabPane :label="$t('scripts.scriptVersions')" name="versions">

        <ElRow v-if="activeTab == 'versions' && !loading && versions" class="mb-20px">
          <ElCol>
            <ElFormItem :label="$t('scripts.scriptVersions')" prop="action">
              <ElSelect
                  v-model="currentVersionIdx"
                  clearable
                  :placeholder="$t('dashboard.editor.selectAction')"
                  style="width: 100%"
                  @change="selectVersionHandler"
              >
                <ElOption
                    v-for="(p, index) in versions"
                    :key="index"
                    :label="parseTime(p.createdAt)"
                    :value="index"/>
              </ElSelect>

            </ElFormItem>
          </ElCol>
          <ElCol>
            <MergeEditor :source="currentScript"
                         :destination="currentVersion"
                         @update:source="onMergeEditorChange"/>
          </ElCol>
        </ElRow>

      </ElTabPane>
      <!-- /versions -->

      <!-- info -->
      <ElTabPane :label="$t('scripts.scriptInfo')" name="info">
        <ElDescriptions v-if="currentScript?.scriptInfo"
                        class="ml-10px mr-10px mb-20px"
                        :title="$t('scripts.scriptInfo')"
                        direction="vertical"
                        :column="2"
                        border
        >
          <ElDescriptionsItem :label="$t('scripts.alexaIntents')">{{ currentScript.scriptInfo.alexaIntents }}
          </ElDescriptionsItem>
          <ElDescriptionsItem :label="$t('scripts.entityActions')">{{ currentScript.scriptInfo.entityActions }}
          </ElDescriptionsItem>
          <ElDescriptionsItem :label="$t('scripts.entityScripts')">{{ currentScript.scriptInfo.entityScripts }}
          </ElDescriptionsItem>
          <ElDescriptionsItem :label="$t('scripts.automationTriggers')">
            {{ currentScript.scriptInfo.automationTriggers }}
          </ElDescriptionsItem>
          <ElDescriptionsItem :label="$t('scripts.automationConditions')">
            {{ currentScript.scriptInfo.automationConditions }}
          </ElDescriptionsItem>
          <ElDescriptionsItem :label="$t('scripts.automationActions')">
            {{ currentScript.scriptInfo.automationActions }}
          </ElDescriptionsItem>
        </ElDescriptions>
      </ElTabPane>
      <!-- /info -->

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
