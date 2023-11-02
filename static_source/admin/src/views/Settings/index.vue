<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {onMounted, onUnmounted, ref} from 'vue'
import {useAppStore} from "@/store/modules/app";
import api from "@/api/api";
import {ElCol, ElDivider, ElForm, ElFormItem, ElRow, ElSwitch, ElInputNumber, ElInput} from 'element-plus'
import {ApiDashboard} from '@/api/stub';
import DashboardSearch from "@/views/Dashboard/components/DashboardSearch.vue";
import ContentWrap from "@/components/ContentWrap/src/ContentWrap.vue";
import {UUID} from "uuid-generator-ts";
import stream from "@/api/stream";
import {EventStateChange} from "@/api/stream_types";
import {debounce} from "lodash-es";

const appStore = useAppStore()

const {t} = useI18n()

export interface Settings {
  mainDashboardDark?: ApiDashboard;
  mainDashboardLight?: ApiDashboard;
  devDashboardDark?: ApiDashboard;
  devDashboardLight?: ApiDashboard;
  restartComponentIfScriptChanged?: boolean;
  clearMetricsDays?: number;
  clearLogsDays?: number;
  clearEntityStorageDays?: number;
  clearRunHistoryDays?: number;
  timezone?: string;
}

const settings = ref<Settings>({} as Settings)
const loading = ref(true)

const getDashboardVar = async (name: string) => {
  await api.v1.variableServiceGetVariableByName(name).then((resp) => {
    if (!resp.data?.value || resp.data?.value == 'undefined') {
      return;
    }
    const id: number = parseInt(resp.data?.value);
    api.v1.dashboardServiceGetDashboardById(id).then((resp) => {
      settings.value[name] = resp.data
    });
  });
}

const getBooleanVar = async (name: string) => {
  await api.v1.variableServiceGetVariableByName(name).then((resp) => {
    console.log(resp.data?.value, typeof resp.data?.value)
    if ((typeof resp.data?.value == 'string' && resp.data?.value == 'false') ||
        !resp.data?.value || resp.data?.value == 'undefined') {
      settings.value[name] = false;
      return;
    }
    settings.value[name] = true;
  });
}

const getIntegerVar = async (name: string) => {
  await api.v1.variableServiceGetVariableByName(name).then((resp) => {
    const id: number = parseInt(resp.data?.value);
    settings.value[name] = id;
  });
}

const getStringVar = async (name: string) => {
  await api.v1.variableServiceGetVariableByName(name).then((resp) => {
    settings.value[name] = resp.data?.value;
  });
}

const getSettings = async () => {
  loading.value = true
  await Promise.all([
    getDashboardVar('mainDashboardDark'),
    getDashboardVar('mainDashboardLight'),
    getDashboardVar('devDashboardDark'),
    getDashboardVar('devDashboardLight'),
    getBooleanVar('restartComponentIfScriptChanged'),
    getIntegerVar('clearMetricsDays'),
    getIntegerVar('clearLogsDays'),
    getIntegerVar('clearEntityStorageDays'),
    getIntegerVar('clearRunHistoryDays'),
    getStringVar('timezone')
  ])
  loading.value = false
}

const changedVariable = debounce((name: string) => {
  let value = ''
  if (settings.value[name] != undefined) {
    value = settings.value[name] + ''
  }
  if (settings.value[name]?.id) {
    value = settings.value[name]?.id + ''
  }
  api.v1.variableServiceUpdateVariable(name, {value: value})
}, 500)

const currentID = ref('')

const onStateChanged = (event: EventStateChange) => {
  getSettings()
}

onMounted(() => {
  const uuid = new UUID()
  currentID.value = uuid.getDashFreeUUID()

  setTimeout(() => {
    stream.subscribe('event_removed_variable_model', currentID.value, onStateChanged);
    stream.subscribe('event_updated_variable_model', currentID.value, onStateChanged);
  }, 200)
})

onUnmounted(() => {
  stream.unsubscribe('event_removed_variable_model', currentID.value);
  stream.unsubscribe('event_updated_variable_model', currentID.value);
})


getSettings()

</script>

<template>

  <ContentWrap>

    <ElDivider content-position="left">{{$t('settings.dashboardOptions')}}</ElDivider>

    <ElForm label-position="top">
      <ElRow :gutter="24">
        <ElCol :span="12" :xs="12">
          <ElFormItem :label="$t('settings.mainDarkDashboard')" prop="mainDashboard">
            <DashboardSearch v-model="settings.mainDashboardDark" @update:modelValue="changedVariable('mainDashboardDark')"/>
          </ElFormItem>
        </ElCol>
        <ElCol :span="12" :xs="12">
          <ElFormItem :label="$t('settings.devDarkDashboard')" prop="devDashboard">
            <DashboardSearch v-model="settings.devDashboardDark" @update:modelValue="changedVariable('devDashboardDark')"/>
          </ElFormItem>
        </ElCol>
      </ElRow>
      <ElRow :gutter="24">
        <ElCol :span="12" :xs="12">
          <ElFormItem :label="$t('settings.mainLightDashboard')" prop="mainDashboard">
            <DashboardSearch v-model="settings.mainDashboardLight" @update:modelValue="changedVariable('mainDashboardLight')"/>
          </ElFormItem>
        </ElCol>
        <ElCol :span="12" :xs="12">
          <ElFormItem :label="$t('settings.devLightDashboard')" prop="devDashboard">
            <DashboardSearch v-model="settings.devDashboardLight" @update:modelValue="changedVariable('devDashboardLight')"/>
          </ElFormItem>
        </ElCol>
      </ElRow>

      <ElDivider content-position="left">{{$t('settings.scriptsOptions')}}</ElDivider>

      <ElRow :gutter="24">
        <ElCol :span="12" :xs="12">
          <ElFormItem :label="$t('settings.restartComponentIfScriptChanged')" prop="restartComponentIfScriptChanged">
            <ElSwitch v-model="settings.restartComponentIfScriptChanged" @update:modelValue="changedVariable('restartComponentIfScriptChanged')"/>
          </ElFormItem>
        </ElCol>
        <ElCol :span="12" :xs="12"/>
      </ElRow>

    <ElDivider content-position="left">{{$t('settings.clearHistory')}}</ElDivider>

    <ElRow :gutter="24">
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('settings.clearMetricsDays')" prop="clearMetricsDays">
          <ElInputNumber size="small" v-model="settings.clearMetricsDays" @update:modelValue="changedVariable('clearMetricsDays')"/>
        </ElFormItem>
      </ElCol>
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('settings.clearLogsDays')" prop="clearLogsDays">
          <ElInputNumber size="small" v-model="settings.clearLogsDays" @update:modelValue="changedVariable('clearLogsDays')"/>
        </ElFormItem>
      </ElCol>
    </ElRow>

    <ElRow :gutter="24">
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('settings.clearEntityStorageDays')" prop="clearEntityStorageDays">
          <ElInputNumber size="small" v-model="settings.clearEntityStorageDays" @update:modelValue="changedVariable('clearEntityStorageDays')"/>
        </ElFormItem>
      </ElCol>
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('settings.clearRunHistoryDays')" prop="clearRunHistoryDays">
          <ElInputNumber size="small" v-model="settings.clearRunHistoryDays"  @update:modelValue="changedVariable('clearRunHistoryDays')"/>
        </ElFormItem>
      </ElCol>
    </ElRow>

    <ElDivider content-position="left">{{$t('settings.time')}}</ElDivider>

    <ElRow :gutter="24">
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('settings.timezone')" prop="timezone">
          <ElInput size="small" v-model="settings.timezone" @update:modelValue="changedVariable('timezone')"/>
        </ElFormItem>
      </ElCol>
    </ElRow>

    </ElForm>

  </ContentWrap>
</template>

<style>

</style>
