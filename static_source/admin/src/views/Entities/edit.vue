<script setup lang="ts">
import {computed, onMounted, onUnmounted, reactive, ref, unref} from 'vue'
import {useI18n} from '@/hooks/web/useI18n'
import {ElButton, ElMessage, ElPopconfirm, ElTabs, ElTabPane} from 'element-plus'
import {useRoute, useRouter} from 'vue-router'
import api from "@/api/api";
import Form from './components/Form.vue'
import ContentWrap from "@/components/ContentWrap/src/ContentWrap.vue";
import {Attribute, Entity, EntityAction, EntityState, Plugin} from "@/views/Entities/components/types";
import Actions from "@/views/Entities/components/Actions.vue";
import {useEmitt} from "@/hooks/web/useEmitt";
import {
  ApiArea,
  ApiEntityAction,
  ApiEntityState,
  ApiPlugin, ApiScript,
  ApiUpdateEntityRequestAction,
  ApiUpdateEntityRequestState
} from "@/api/stub";
import States from "@/views/Entities/components/States.vue";
import AttributesEditor from "@/views/Entities/components/AttributesEditor.vue";
import { Dialog } from '@/components/Dialog'
import JsonViewer from "@/components/JsonViewer/JsonViewer.vue";
import {copyToClipboard} from "@/utils/clipboard";
import {EventStateChange} from "@/api/stream_types";
import {UUID} from "uuid-generator-ts";
import stream from "@/api/stream";
import Storage from "@/views/Entities/components/Storage.vue";
import Metrics from "@/views/Entities/components/Metrics.vue";

const {push} = useRouter()
const route = useRoute();
const {t} = useI18n()

const writeRef = ref<ComponentRef<typeof Form>>()
const loading = ref(false)
const entityId = computed(() => route.params.id as string);
const currentEntity = ref<Nullable<Entity>>(null)
const activeTab = ref('attributes')
const dialogSource = ref({})
const dialogVisible = ref(false)
const lastEvent = ref<Nullable<EventStateChange>>(null)

interface Internal {
  attributes: Attribute[];
  settings: Attribute[];
}
const internal = reactive<Internal>(
    {
      attributes: [],
      settings: [],
    },
);

const fetch = async () => {
  loading.value = true
  const res = await api.v1.entityServiceGetEntity(entityId.value)
      .catch(() => {
      })
      .finally(() => {
        loading.value = false
      })
  if (res) {
    const entity = res.data as Entity

    // attributes
    internal.attributes = [];
    if (entity.attributes) {
      for (const key in entity.attributes) {
        internal.attributes.push(entity.attributes[key]);
      }
    }

    // settings
    internal.settings = [];
    if (entity.settings) {
      for (const key in entity.settings) {
        internal.settings.push(entity.settings[key]);
      }
    }

    currentEntity.value = entity

    await fetchPlugin()

    loading.value = false
  } else {
    currentEntity.value = null
  }
}

const fetchPlugin = async () => {
  loading.value = true
  const res = await api.v1.pluginServiceGetPlugin(currentEntity.value.pluginName)
  if (res) {
    const plugin = res.data as ApiPlugin;

    // attributes
    let actorAttrs: Attribute[] = [];
    if (plugin.options?.actorAttrs) {
      for (const key in plugin.options.actorAttrs) {
        actorAttrs.push(plugin.options.actorAttrs[key]);
      }
    }

    // actorSetts
    let actorSetts: Attribute[] = [];
    if (plugin.options?.actorSetts) {
      for (const key in plugin.options.actorSetts) {
        actorSetts.push(plugin.options.actorSetts[key]);
      }
    }

    // setts
    let setts: Attribute[] = [];
    if (plugin.options?.setts) {
      for (const key in plugin.options?.setts) {
        setts.push(plugin.options?.setts[key]);
      }
    }

    // actorActions
    let actorActions: EntityAction[] = [];
    if (plugin.options?.actorActions) {
      for (const key in plugin.options?.actorActions) {
        actorActions.push(plugin.options?.actorActions[key]);
      }
    }

    // actorStates
    let actorStates: EntityState[] = [];
    if (plugin.options?.actorStates) {
      for (const key in plugin.options?.actorStates) {
        actorStates.push(plugin.options?.actorStates[key]);
      }
    }

    currentEntity.value.plugin = {
      name: plugin.name,
      version: plugin.version,
      enabled: plugin.enabled,
      system: plugin.system,
      actor: plugin.actor,
      actorCustomAttrs: plugin.options?.actorCustomAttrs,
      actorCustomActions: plugin.options?.actorCustomActions,
      actorCustomStates: plugin.options?.actorCustomStates,
      actorCustomSetts: plugin.options?.actorCustomSetts,
      triggers: plugin.options?.triggers,
      actorAttrs: actorAttrs,
      actorSetts: actorSetts,
      setts: setts,
      actorActions: actorActions,
      actorStates: actorStates,
    } as Plugin;
    // console.log(plugin)
  }
}

const prepareForSave = async () => {
  const write = unref(writeRef)
  const validate = await write?.elFormRef?.validate()?.catch(() => {
  })
  if (validate) {
    const data = (await write?.getFormData()) as Entity
    let actions: ApiUpdateEntityRequestAction[] = []
    let states: ApiUpdateEntityRequestState[] = []
    for (const a of data?.actions) {
      actions.push({
        name: a.name,
        description: a.description,
        icon: a.icon,
        imageId: a.image?.id || a.imageId,
        scriptId: a.script?.id || a.scriptId,
        type: a.type,
      })
    }
    for (const a of data?.states) {
      states.push({
        name: a.name,
        description: a.description,
        imageId: a.image?.id || a.imageId,
      })
    }
    let attributes: { [key: string]: Attribute } = {};
    for (const index in internal.attributes) {
      attributes[internal.attributes[index].name] = internal.attributes[index];
    }
    let settings: { [key: string]: Attribute } = {};
    for (const index in internal.settings) {
      settings[internal.settings[index].name] = internal.settings[index];
    }
    const body = {
      id: data.id,
      pluginName: data.plugin?.name || data.pluginName,
      description: data.description,
      areaId: data.area?.id,
      icon: data.icon,
      imageId: data.image?.id,
      autoLoad: data.autoLoad,
      restoreState: data.restoreState,
      parentId: data.parent?.id,
      actions: actions,
      states: states,
      attributes: attributes,
      settings: settings,
      scriptIds: data.scriptIds,
      metrics: data.metrics,
      tags: data.tags,
    }
   return body
  }
  return null
}

let _scriptsPromises = []
let _scripts: Map<string, ApiScript> = []
const fetchScript = async (id: number) => {
  const res = await api.v1.scriptServiceGetScriptById(id)
  if (res) {
    _scripts[id] = res.data
  }
}

const prepareForExport = async () => {
  const write = unref(writeRef)
  const validate = await write?.elFormRef?.validate()?.catch(() => {
  })
  if (validate) {
    const data = (await write?.getFormData()) as Entity
    let actions: ApiEntityAction[] = []
    let states: ApiEntityState[] = []

    data.actions.forEach(action => {
      if (action.scriptId) {
        _scripts[action.scriptId] = null
      }
      if (action.script?.id) {
        _scripts[action.script.id] = null
      }
    })
    data.scripts.forEach(script => {
      if (script.id) {
        _scripts[script.id] = null
      }
    })

    _scripts.forEach((value: ApiScript, key: string) => {
      _scriptsPromises.push(fetchScript(key))
    })

    await Promise.all(_scriptsPromises)

    for (const a of data?.actions) {
      let script:ApiScript = null;
      if (a.script) {
        script = {
          id: a.script.id,
          lang: _scripts[a.script.id]?.lang,
          name: a.script.name,
          source: _scripts[a.script.id]?.source || '',
          description: a.script.description,
        } as ApiScript
      }
      actions.push({
        name: a.name,
        description: a.description,
        icon: a.icon,
        image: a.image,
        script: script,
        scriptId: script?.id || 0,
        type: a.type,
      } as ApiEntityAction)
    }
    for (const a of data?.states) {
      states.push({
        name: a.name,
        description: a.description,
        image: a.image,
        icon: a.icon,
      } as ApiEntityState)
    }
    let attributes: { [key: string]: Attribute } = {};
    for (const index in internal.attributes) {
      attributes[internal.attributes[index].name] = internal.attributes[index];
    }
    let settings: { [key: string]: Attribute } = {};
    for (const index in internal.settings) {
      settings[internal.settings[index].name] = internal.settings[index];
    }
    let area: ApiArea = null
    if (data.area) {
      area = {
        id: data.area.id,
        name: data.area.name,
        description: data.area.description,
        polygon: data.area.polygon,
      } as ApiArea
    }
    let scripts: ApiScript[] = [];
    if (data.scripts) {
      data.scripts.forEach((value: ApiScript) => {
        scripts.push({
          id: value.id,
          lang: _scripts[value.id]?.lang,
          name: value.name,
          source: _scripts[value.id]?.source || '',
          description: value.description,
        } as ApiScript)
      })
    }
    const body = {
      id: data.id,
      pluginName: data.plugin?.name || data.pluginName,
      description: data.description,
      area: area,
      icon: data.icon,
      image: data.image,
      autoLoad: data.autoLoad,
      restoreState: data.restoreState,
      parent: data.parent,
      actions: actions,
      states: states,
      attributes: attributes,
      settings: data.settings,
      scripts: scripts,
      metrics: data.metrics,
      tags: data.tags,
    }
   return body
  }
  return null
}

const save = async () => {
  const body = await prepareForSave()
  if (!body) {
    return
  }
  const res = await api.v1.entityServiceUpdateEntity(entityId.value, body)
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

const callAction = async (name: string) => {
  await api.v1.interactServiceEntityCallAction({id: entityId.value, name: name});
  ElMessage({
    title: t('Success'),
    message: t('message.callSuccessful'),
    type: 'success',
    duration: 2000
  });
}

const setState = async (name: string) => {
  await api.v1.developerToolsServiceEntitySetState({id: entityId.value, name: name});
  ElMessage({
    title: t('Success'),
    message: t('message.callSuccessful'),
    type: 'success',
    duration: 2000
  });
}

useEmitt({
  name: 'updateActions',
  callback: (val: EntityAction[]) => {
    const second = JSON.parse(JSON.stringify(val))
    currentEntity.value.actions = JSON.parse(JSON.stringify(second))
  }
})

useEmitt({
  name: 'updateStates',
  callback: (val: EntityState[]) => {
    const second = JSON.parse(JSON.stringify(val))
    currentEntity.value.states = JSON.parse(JSON.stringify(second))
  }
})

useEmitt({
  name: 'callAction',
  callback: (val: string) => {
    callAction(val)
  }
})

useEmitt({
  name: 'setState',
  callback: (val: string) => {
    setState(val)
  }
})

const onAttrsUpdated = (attrs: Attribute[]) => {
  const second = JSON.parse(JSON.stringify(attrs))
  internal.attributes = JSON.parse(JSON.stringify(second))
}

const onSettingsUpdated = (attrs: Attribute[]) => {
  const second = JSON.parse(JSON.stringify(attrs))
  internal.settings = JSON.parse(JSON.stringify(second))
}

const cancel = () => {
  push('/entities')
}

const copy = async () => {
  const body = await prepareForExport()
  copyToClipboard(JSON.stringify(body, null, 2))
}

const exportEntity = async () => {
  const body = await prepareForExport()
  dialogSource.value = body
  dialogVisible.value = true
}

const restart = async () => {
  await api.v1.developerToolsServiceReloadEntity({id: entityId.value});
  ElMessage({
    title: t('Success'),
    message: t('message.reloadSuccessful'),
    type: 'success',
    duration: 2000
  });
}

const remove = async () => {
  loading.value = true
  const res = await api.v1.entityServiceDeleteEntity(entityId.value)
      .catch(() => {
      })
      .finally(() => {
        loading.value = false
      })
  if (res) {
    cancel()
  }
}

const currentID = ref('')

const onStateChanged = (event: EventStateChange) => {
  if (event.entity_id != entityId.value) {
    return;
  }

  lastEvent.value = event;
}

const requestCurrentState = () => {
  stream.send({
    id: UUID.createUUID(),
    query: 'event_get_last_state',
    body: btoa(JSON.stringify({'entity_id': entityId.value}))
  });
}

const handleTabChange = (tab: any, event: any) => {
  const {paneName} = tab;
  if (paneName == 'currentState') {
    requestCurrentState();
  }
}

onMounted(() => {
  const uuid = new UUID()
  currentID.value = uuid.getDashFreeUUID()

  setTimeout(() => {
    stream.subscribe('state_changed', currentID.value, onStateChanged);
  }, 1000)
})

onUnmounted(() => {
  stream.unsubscribe('state_changed', currentID.value);
})

fetch()

</script>

<template>
  <ContentWrap>
    <el-tabs class="demo-tabs" v-model="activeTab" @tab-click="handleTabChange">
      <!-- main -->
      <el-tab-pane :label="$t('entities.main')" name="main">
        <Form ref="writeRef" :current-row="currentEntity"/>
      </el-tab-pane>
      <!-- /main -->

      <!-- actions -->
      <el-tab-pane :label="$t('entities.actions')" name="actions">
        <Actions
            :actions="currentEntity?.actions"
            :custom-actions="currentEntity?.plugin?.actorCustomActions"
            :plugin-actions="currentEntity?.plugin?.actorActions"
        />
      </el-tab-pane>
      <!-- /actions -->

      <!-- states -->
      <el-tab-pane :label="$t('entities.states')" name="states">
        <States
            :states="currentEntity?.states"
            :custom-states="currentEntity?.plugin?.actorCustomStates"
            :plugin-states="currentEntity?.plugin?.actorStates"
        />
      </el-tab-pane>
      <!-- /states -->

      <!-- attributes -->
      <el-tab-pane :label="$t('entities.attributes')" name="attributes">
        <AttributesEditor
            v-model="internal.attributes"
            :custom-attrs="currentEntity?.plugin?.actorCustomAttrs"
            :plugin-attrs="currentEntity?.plugin?.actorAttrs"
            @change="onAttrsUpdated"
        />
      </el-tab-pane>
      <!-- /attributes -->

      <!-- settings -->
      <el-tab-pane :label="$t('entities.settings')" name="settings">
        <AttributesEditor
            v-model="internal.settings"
            :custom-attrs="currentEntity?.plugin?.actorCustomSetts"
            :plugin-attrs="currentEntity?.plugin?.actorSetts"
            @change="onSettingsUpdated"
        />
      </el-tab-pane>
      <!-- /settings -->

      <!-- metrics -->
      <el-tab-pane :label="$t('entities.metrics')" name="metrics">
        <Metrics :metrics="currentEntity?.metrics"/>
      </el-tab-pane>
      <!-- /metrics -->

      <!-- storage -->
      <el-tab-pane :label="$t('entities.storage')" name="storage">
        <Storage v-model="currentEntity"/>
      </el-tab-pane>
      <!-- /storage -->

      <!-- current state -->
      <el-tab-pane :label="$t('entities.currentState')" name="currentState">
        <ElButton type="default" @click.prevent.stop="requestCurrentState()" class="mb-20px">
          <Icon icon="ep:refresh" class="mr-5px"/> {{ $t('main.currentState') }}
        </ElButton>

        <JsonViewer v-model="lastEvent"/>
      </el-tab-pane>
      <!-- /current state -->

    </el-tabs>


    <div style="text-align: right" class="mt-20px">

      <ElButton type="primary" @click="save()">
        {{ t('main.save') }}
      </ElButton>

      <ElButton type="primary" @click="restart()">
        <Icon icon="ep:refresh" class="mr-5px"/>
        {{ t('main.restart') }}
      </ElButton>

      <ElButton type="primary" @click="exportEntity()">
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
  <Dialog v-model="dialogVisible" :title="t('entities.dialogExportTitle')" :maxHeight="400" width="80%">
    <JsonViewer v-model="dialogSource"/>
<!--    <template #footer>-->
<!--      <ElButton @click="copy()">{{ t('setting.copy') }}</ElButton>-->
<!--      <ElButton @click="dialogVisible = false">{{ t('main.closeDialog') }}</ElButton>-->
<!--    </template>-->
  </Dialog>
  <!-- /export dialog -->

</template>

<style lang="less" scoped>

</style>
