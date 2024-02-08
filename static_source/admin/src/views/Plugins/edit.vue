<script setup lang="ts">
import {ElButton, ElTabs, ElTabPane, ElRow, ElCol, ElEmpty} from 'element-plus'
import api from "@/api/api";
import {useRoute, useRouter} from "vue-router";
import {computed, ref, unref} from "vue";
import {
  ApiAttribute,
  ApiPlugin,
  ApiPluginOptionsResultEntityAction,
  ApiPluginOptionsResultEntityState
} from "@/api/stub";
import Form from './components/Form.vue'
import {useI18n} from "@/hooks/web/useI18n";
import {Plugin} from './components/Types.ts'
import {AttributesViewer} from "@/components/AttributesViewer";
import ActionsViewer from "@/views/Plugins/components/ActionsViewer.vue";
import StatesViewer from "@/views/Plugins/components/StatesViewer.vue";
import SettingsEditor from "@/views/Plugins/components/SettingsEditor.vue";
import {useEmitt} from "@/hooks/web/useEmitt";
import ContentWrap from "@/components/ContentWrap/src/ContentWrap.vue";
import {useCache} from "@/hooks/web/useCache";

const {t} = useI18n()
const writeRef = ref<ComponentRef<typeof Form>>()
const route = useRoute();
const {currentRoute, addRoute, push} = useRouter()
const pluginName = computed<string>(() => route.params.name as string);

const currentPlugin = ref<Nullable<Plugin>>(null)
const loading = ref(false)
const activeTab = ref('main')
const lastState = ref<boolean>(false)

const fetch = async () => {
  loading.value = true
  const res = await api.v1.pluginServiceGetPlugin(pluginName.value as string)
      .catch(() => {
      })
      .finally(() => {
        loading.value = false
      })
  if (res) {
    const plugin: ApiPlugin = res.data;


    // setts
    let settings: ApiAttribute[] = [];
    if (res.data.options?.setts) {
      for (const key in res.data.options.setts) {
        let st = res.data.options.setts[key]
        if (res.data.settings[key]) {
          st = res.data.settings[key]
        }
        settings.push(st)
      }
    }

    // actor states
    let actorStates: ApiPluginOptionsResultEntityState[] = [];
    if (res.data.options?.actorStates) {
      for (const key in res.data.options.actorStates) {
        actorStates.push(res.data.options.actorStates[key])
      }
    }

    // actor actions
    let actorActions: ApiPluginOptionsResultEntityAction[] = [];
    if (res.data.options?.actorActions) {
      for (const key in res.data.options.actorActions) {
        actorActions.push(res.data.options.actorActions[key])
      }
    }

    // actor attributes
    // let actorAttrs: ApiAttribute[] = []
    // if (res.data.options?.actorAttrs) {
    //   for (const key in res.data.options.actorAttrs) {
    //     actorAttrs.push(res.data.options.actorAttrs[key])
    //   }
    // }

    // actor attributes
    // let actorSetts: ApiAttribute[] = []
    // if (res.data.options?.actorSetts) {
    //   for (const key in res.data.options.actorSetts) {
    //     actorSetts.push(res.data.options.actorSetts[key])
    //   }
    // }

    currentPlugin.value = {
      name: plugin.name,
      version: plugin.version,
      enabled: plugin.enabled,
      system: plugin.system,
      actor: plugin.actor,
      triggers: plugin.options?.triggers,
      actors: plugin.options?.actors,
      actorCustomAttrs: plugin.options?.actorCustomAttrs,
      actorCustomActions: plugin.options?.actorCustomActions,
      actorCustomStates: plugin.options?.actorCustomStates,
      actorCustomSetts: plugin.options?.actorCustomSetts,
      setts: settings,
      actorStates: actorStates,
      actorActions: actorActions,
      actorAttrs: Object.assign({}, res.data.options?.actorAttrs),
      actorSetts: Object.assign({}, res.data.options?.actorSetts),
    } as Plugin
    console.log(currentPlugin.value)
    const {enabled} = res.data;
    lastState.value = enabled
  } else {
    currentPlugin.value = null
  }
}

const save = async () => {
  const write = unref(writeRef)
  const validate = await write?.elFormRef?.validate()?.catch(() => {
  })
  if (validate) {
    const data = (await write?.getFormData()) as ApiPlugin
    if (data.enabled === lastState.value) {
      return
    }
    lastState.value = data.enabled || false
    if (data.enabled) {
      await api.v1.pluginServiceEnablePlugin(pluginName.value)
    } else {
      await api.v1.pluginServiceDisablePlugin(pluginName.value)
    }
    fetch()
  }
}

const cancel = () => {
  push('/etc/plugins')
}

const showActorTabIf = (): boolean => {
  if (Object.keys(currentPlugin.value?.actorAttrs || {}).length ||
      Object.keys(currentPlugin.value?.actorActions || {}).length ||
      Object.keys(currentPlugin.value?.actorStates || {}).length ||
      Object.keys(currentPlugin.value?.actorSetts || {}).length) {
    return true
  }
  return false
}

const showSettingsTabIf = (): boolean => {
  if (currentPlugin.value?.setts && Object.keys(currentPlugin.value?.setts || {}).length) {
    return true
  }
  return false
}

const saveSetting = async (val: ApiAttribute[]) => {
  let settings: { [key: string]: ApiAttribute } = {};
  for (const index in val) {
    settings[val[index].name] = val[index];
  }
  await api.v1.pluginServiceUpdatePluginSettings(pluginName.value, {settings: settings})
  fetch()
}

const readme = ref('');
const { wsCache } = useCache()
const getReadme = async () => {
  const lang = wsCache.get('lang') || 'en';
  const res = await api.v1.pluginServiceGetPluginReadme(pluginName.value, {lang: lang})
      .catch(() => {
      })
      .finally(() => {
        loading.value = false
      })
  if (res && res.data) {
    readme.value = res.data;
  }
}

const tabHandler = (tab: any, ev: any) => {
  const {props} = tab;
  if (props.name != 'readme') return;
  getReadme()
}

useEmitt({
  name: 'settingsUpdated',
  callback: (settings: ApiAttribute[]) => {
    saveSetting(settings)
  }
})

fetch()

</script>

<template>
  <ContentWrap>
    <el-tabs class="demo-tabs" v-model="activeTab" @tab-click="tabHandler">
      <el-tab-pane :label="$t('plugins.main')" name="main">
        <Form ref="writeRef" :current-row="currentPlugin"/>
        <div style="text-align: right">
          <ElButton type="primary" @click="save()">
            {{ t('main.save') }}
          </ElButton>
        </div>
      </el-tab-pane>
      <!--  /Main  -->
      <el-tab-pane
          :label="$t('plugins.actor')"
          :disabled="!showActorTabIf()"
          name="actor">

        <!-- attributes-->
        <el-row class="mt-10px"
                v-if="currentPlugin?.actorAttrs && Object.keys(currentPlugin?.actorAttrs || {}).length">
          <el-col>
            {{ $t('plugins.actorAttrs') }}
          </el-col>
        </el-row>
        <el-row class="mt-20px" v-if="currentPlugin?.actorAttrs && Object.keys(currentPlugin?.actorAttrs || {}).length">
          <el-col>
            <AttributesViewer v-model="currentPlugin.actorAttrs"/>
          </el-col>
        </el-row>
        <!-- /attributes-->
        <!-- actions-->
        <el-row class="mt-20px"
                v-if="currentPlugin?.actorActions && Object.keys(currentPlugin?.actorActions).length">
          <el-col>
            {{ $t('plugins.actorActions') }}
          </el-col>
        </el-row>
        <el-row class="mt-20px"
                v-if="currentPlugin?.actorActions && Object.keys(currentPlugin?.actorActions).length">
          <el-col>
            <ActionsViewer :actions="currentPlugin.actorActions"/>
          </el-col>
        </el-row>
        <!-- /actions-->
        <!-- states-->
        <el-row class="mt-20px"
                v-if="currentPlugin?.actorStates && Object.keys(currentPlugin?.actorStates).length">
          <el-col>
            {{ $t('plugins.actorStates') }}
          </el-col>
        </el-row>
        <el-row class="mt-20px"
                v-if="currentPlugin?.actorStates && Object.keys(currentPlugin?.actorStates).length">
          <el-col>
            <StatesViewer :states="currentPlugin.actorStates"/>
          </el-col>
        </el-row>
        <!-- /states-->
        <!-- settings-->
        <el-row class="mt-10px"
                v-if="currentPlugin?.actorSetts && Object.keys(currentPlugin?.actorSetts || {}).length">
          <el-col>
            {{ $t('plugins.actorSettings') }}
          </el-col>
        </el-row>
        <el-row class="mt-20px"
                v-if="currentPlugin?.actorSetts && Object.keys(currentPlugin?.actorSetts || {}).length">
          <el-col>
            <AttributesViewer v-model="currentPlugin.actorSetts"/>
          </el-col>
        </el-row>
        <!-- /settings-->
      </el-tab-pane>
      <!--  /Actor  -->

      <el-tab-pane
          :label="$t('plugins.settings')"
          name="settings"
          :disabled="!showSettingsTabIf()">
        <div v-if="currentPlugin?.setts && Object.keys(currentPlugin?.setts || {}).length">
          <SettingsEditor :attrs="currentPlugin.setts"/>
        </div>
      </el-tab-pane>
      <!--  /Settings  -->

      <el-tab-pane :label="$t('plugins.readme')" lazy name="readme">
        <div v-html="readme"></div>
        <ElEmpty v-if="readme == ''" description="no info"/>
      </el-tab-pane>
    </el-tabs>
  </ContentWrap>
</template>

<style lang="less" scoped>

</style>
