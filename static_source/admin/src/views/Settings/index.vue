<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {onMounted, onUnmounted, ref} from 'vue'
import {useAppStore} from "@/store/modules/app";
import api from "@/api/api";
import {ElCol, ElDivider, ElForm, ElFormItem, ElInput, ElInputNumber, ElRow, ElSwitch} from 'element-plus'
import {ApiDashboard, ApiEntity} from '@/api/stub';
import {DashboardSearch} from "@/components/DashboardSearch";
import {ContentWrap} from "@/components/ContentWrap";
import {UUID} from "uuid-generator-ts";
import stream from "@/api/stream";
import {EventStateChange} from "@/api/types";
import {debounce} from "lodash-es";
import {EntitySearch} from "@/components/EntitySearch";
import {Infotip} from '@/components/Infotip'

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
  createBackupAt?: string;
  maximumNumberOfBackups?: number;
  sendbackuptoTelegramBot?: ApiEntity;
  sendTheBackupInPartsMb?: number;
  gateClientId?: string;
  gateClientSecretKey?: string;
  gateClientServerHost?: string;
  gateClientServerPort?: number;
  gateClientPoolIdleSize?: number;
  gateClientPoolMaxSize?: number;
  gateClientTLS?: boolean;
  hmacKey?: string;
  certPublic?: string;
  certKey?: string;
}

const settings = ref<Settings>({} as Settings)
const loading = ref(true)

const getDashboardVar = async (name: string) => {
  await api.v1.variableServiceGetVariableByName(name).then((resp) => {
    if (!resp.data?.value || resp.data?.value == 'undefined') {
      return;
    }
    const id: number = parseInt(resp.data.value);
    api.v1.dashboardServiceGetDashboardById(id).then((resp) => {
      settings.value[name] = resp.data
    });
  });
}

const getBooleanVar = async (name: string) => {
  await api.v1.variableServiceGetVariableByName(name).then((resp) => {
    // console.log(resp.data?.value, typeof resp.data?.value)
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
    if (!resp.data?.value || resp.data?.value == 'undefined') {
      return;
    }
    const id: number = parseInt(resp.data.value);
    settings.value[name] = id;
  });
}

const getStringVar = async (name: string) => {
  await api.v1.variableServiceGetVariableByName(name).then((resp) => {
    settings.value[name] = resp.data?.value || '';
  });
}

const getStringEntity = async (name: string) => {
  await api.v1.variableServiceGetVariableByName(name).then((resp) => {
    if (!resp.data?.value || resp.data?.value == 'undefined') {
      return;
    }
    settings.value[name] = {id: resp.data?.value} as ApiEntity || undefined;
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
    getStringVar('timezone'),
    getStringVar('createBackupAt'),
    getIntegerVar('maximumNumberOfBackups'),
    getStringEntity('sendbackuptoTelegramBot'),
    getIntegerVar('sendTheBackupInPartsMb'),
    getStringVar('gateClientId'),
    getStringVar('gateClientSecretKey'),
    getStringVar('gateClientServerHost'),
    getIntegerVar('gateClientServerPort'),
    getIntegerVar('gateClientPoolIdleSize'),
    getIntegerVar('gateClientPoolMaxSize'),
    getBooleanVar('gateClientTLS'),
    getStringVar('hmacKey'),
    getStringVar('certKey'),
    getStringVar('certPublic')
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

const onEntityChanged = async (entity: ApiEntity, name: string) => {
  api.v1.variableServiceUpdateVariable(name, {value: entity?.id || ''})
}

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

    <ElDivider content-position="left">{{ $t('settings.dashboardOptions') }}</ElDivider>

    <ElForm label-position="top">
      <ElRow :gutter="24">
        <ElCol :span="12" :xs="12">
          <ElFormItem :label="$t('settings.mainDarkDashboard')" prop="mainDashboard">
            <DashboardSearch v-model="settings.mainDashboardDark"
                             @update:modelValue="changedVariable('mainDashboardDark')"/>
          </ElFormItem>
        </ElCol>
        <ElCol :span="12" :xs="12">
          <ElFormItem :label="$t('settings.devDarkDashboard')" prop="devDashboard">
            <DashboardSearch v-model="settings.devDashboardDark"
                             @update:modelValue="changedVariable('devDashboardDark')"/>
          </ElFormItem>
        </ElCol>
      </ElRow>
      <ElRow :gutter="24">
        <ElCol :span="12" :xs="12">
          <ElFormItem :label="$t('settings.mainLightDashboard')" prop="mainDashboard">
            <DashboardSearch v-model="settings.mainDashboardLight"
                             @update:modelValue="changedVariable('mainDashboardLight')"/>
          </ElFormItem>
        </ElCol>
        <ElCol :span="12" :xs="12">
          <ElFormItem :label="$t('settings.devLightDashboard')" prop="devDashboard">
            <DashboardSearch v-model="settings.devDashboardLight"
                             @update:modelValue="changedVariable('devDashboardLight')"/>
          </ElFormItem>
        </ElCol>
      </ElRow>

      <ElDivider content-position="left">{{ $t('settings.scriptsOptions') }}</ElDivider>

      <ElRow :gutter="24">
        <ElCol :span="12" :xs="12">
          <ElFormItem :label="$t('settings.restartComponentIfScriptChanged')" prop="restartComponentIfScriptChanged">
            <ElSwitch v-model="settings.restartComponentIfScriptChanged"
                      @update:modelValue="changedVariable('restartComponentIfScriptChanged')"/>
          </ElFormItem>
        </ElCol>
        <ElCol :span="12" :xs="12"/>
      </ElRow>

      <ElDivider content-position="left">{{ $t('settings.clearHistory') }}</ElDivider>

      <ElRow :gutter="24">
        <ElCol :span="12" :xs="12">
          <ElFormItem :label="$t('settings.clearMetricsDays')" prop="clearMetricsDays">
            <ElInputNumber v-model="settings.clearMetricsDays" @update:modelValue="changedVariable('clearMetricsDays')"
                           :min="1"/>
          </ElFormItem>
        </ElCol>
        <ElCol :span="12" :xs="12">
          <ElFormItem :label="$t('settings.clearLogsDays')" prop="clearLogsDays">
            <ElInputNumber v-model="settings.clearLogsDays" @update:modelValue="changedVariable('clearLogsDays')"
                           :min="1"/>
          </ElFormItem>
        </ElCol>
      </ElRow>

      <ElRow :gutter="24">
        <ElCol :span="12" :xs="12">
          <ElFormItem :label="$t('settings.clearEntityStorageDays')" prop="clearEntityStorageDays">
            <ElInputNumber v-model="settings.clearEntityStorageDays"
                           @update:modelValue="changedVariable('clearEntityStorageDays')" :min="1"/>
          </ElFormItem>
        </ElCol>
        <ElCol :span="12" :xs="12">
          <ElFormItem :label="$t('settings.clearRunHistoryDays')" prop="clearRunHistoryDays">
            <ElInputNumber v-model="settings.clearRunHistoryDays"
                           @update:modelValue="changedVariable('clearRunHistoryDays')" :min="1"/>
          </ElFormItem>
        </ElCol>
      </ElRow>

      <ElDivider content-position="left">{{ $t('settings.time') }}</ElDivider>

      <ElRow :gutter="24">
        <ElCol :span="12" :xs="12">
          <ElFormItem :label="$t('settings.timezone')" prop="timezone">
            <ElInput v-model="settings.timezone" @update:modelValue="changedVariable('timezone')"/>
          </ElFormItem>
        </ElCol>
      </ElRow>

      <ElDivider content-position="left">{{ $t('settings.backup') }}</ElDivider>

      <Infotip
          :show-index="false"
          title="INFO"
          :schema="[
      {
        label: t('settings.info1'),
      },
    ]"
      />

      <ElRow :gutter="24">
        <ElCol :span="12" :xs="12">
          <ElFormItem :label="$t('settings.createBackupAt')" prop="createBackupAt">
            <ElInput v-model="settings.createBackupAt" @update:modelValue="changedVariable('createBackupAt')"/>
          </ElFormItem>
        </ElCol>
        <ElCol :span="12" :xs="12">
          <ElFormItem :label="$t('settings.maximumNumberOfBackups')" prop="maximumNumberOfBackups">
            <ElInputNumber v-model="settings.maximumNumberOfBackups"
                           @update:modelValue="changedVariable('maximumNumberOfBackups')" :min="1"/>
          </ElFormItem>
        </ElCol>
      </ElRow>

      <ElRow :gutter="24">
        <ElCol :span="12" :xs="12">
          <ElFormItem :label="$t('settings.sendbackuptoTelegramBot')" prop="sendbackuptoTelegramBot">
            <EntitySearch v-model="settings.sendbackuptoTelegramBot"
                          @change="onEntityChanged($event, 'sendbackuptoTelegramBot')"/>
          </ElFormItem>
        </ElCol>
        <ElCol :span="12" :xs="12">
          <ElFormItem :label="$t('settings.sendTheBackupInPartsMb')" prop="maximumNumberOfBackups">
            <ElInputNumber v-model="settings.sendTheBackupInPartsMb"
                           @update:modelValue="changedVariable('sendTheBackupInPartsMb')" :min="0"/>
          </ElFormItem>
        </ElCol>
      </ElRow>

      <ElDivider content-position="left">{{ $t('settings.gate') }}</ElDivider>

      <Infotip
          :show-index="false"
          title="INFO"
          :schema="[
      {
        label: t('settings.info2'),
      },
    ]"
      />

      <ElRow :gutter="24">
        <ElCol :span="12" :xs="12">
          <ElFormItem :label="$t('settings.gateClientServerHost')" prop="gateClientServerHost">
            <ElInput v-model="settings.gateClientServerHost"
                     @update:modelValue="changedVariable('gateClientServerHost')" clearable/>
          </ElFormItem>
          <ElFormItem :label="$t('settings.gateClientId')" prop="gateClientId">
            <ElInput type="password" v-model="settings.gateClientId"
                     @update:modelValue="changedVariable('gateClientId')" clearable show-password/>
          </ElFormItem>
          <ElFormItem :label="$t('settings.gateClientSecretKey')" prop="gateClientSecretKey">
            <ElInput type="password" v-model="settings.gateClientSecretKey"
                     @update:modelValue="changedVariable('gateClientSecretKey')" clearable show-password/>
          </ElFormItem>
        </ElCol>
        <ElCol :span="12" :xs="12">
          <ElFormItem :label="$t('settings.gateClientServerPort')" prop="gateClientServerPort">
            <ElInputNumber v-model="settings.gateClientServerPort"
                           @update:modelValue="changedVariable('gateClientServerPort')" :min="1"/>
          </ElFormItem>
          <ElFormItem :label="$t('settings.gateClientPoolIdleSize')" prop="gateClientPoolIdleSize">
            <ElInputNumber v-model="settings.gateClientPoolIdleSize"
                           @update:modelValue="changedVariable('gateClientPoolIdleSize')" :min="1" :max="100"/>
          </ElFormItem>
          <ElFormItem :label="$t('settings.gateClientPoolMaxSize')" prop="gateClientPoolMaxSize">
            <ElInputNumber v-model="settings.gateClientPoolMaxSize"
                           @update:modelValue="changedVariable('gateClientPoolMaxSize')" :min="1" :max="100"/>
          </ElFormItem>
        </ElCol>
      </ElRow>

      <ElRow :gutter="24">
        <ElCol :span="12" :xs="12">
          <ElFormItem :label="$t('settings.gateClientTLS')" prop="gateClientTLS">
            <ElSwitch v-model="settings.gateClientTLS" @update:modelValue="changedVariable('gateClientTLS')"/>
          </ElFormItem>
        </ElCol>
        <ElCol :span="12" :xs="12"/>
      </ElRow>

      <ElDivider content-position="left">{{ $t('settings.hmacKey') }}</ElDivider>

      <Infotip
          type="warning"
          :show-index="false"
          title="WARNING"
          :schema="[
      {
        label: t('settings.info3'),
      },
    ]"
      />

      <ElRow :gutter="24">
        <ElCol :span="12" :xs="12">
          <ElFormItem :label="$t('settings.hmacKey')" prop="hmacKey">
            <ElInput type="password" v-model="settings.hmacKey" @update:modelValue="changedVariable('hmacKey')"
                     show-password/>
          </ElFormItem>
        </ElCol>
        <ElCol :span="12" :xs="12"/>
      </ElRow>

      <ElDivider content-position="left">{{ $t('settings.certificates') }}</ElDivider>

      <Infotip
        :show-index="false"
        title="INFO"
        :schema="[
      {
        label: t('settings.info4'),
      },
    ]"
      />

      <ElRow :gutter="24">
        <ElCol :span="12" :xs="12">
          <ElFormItem :label="$t('settings.certPublic')" prop="certPublic">
            <ElInput type="textarea"
                     :autosize="{minRows: 2, maxRows: 10}"
                     v-model="settings.certPublic"
                     @update:modelValue="changedVariable('certPublic')"/>
          </ElFormItem>
        </ElCol>
        <ElCol :span="12" :xs="12"/>
      </ElRow>

      <ElRow :gutter="24">
        <ElCol :span="12" :xs="12">
          <ElFormItem :label="$t('settings.certKey')" prop="certKey">
            <ElInput type="textarea"
                     :autosize="{minRows: 2, maxRows: 10}"
                     v-model="settings.certKey"
                     @update:modelValue="changedVariable('certKey')"/>
          </ElFormItem>
        </ElCol>
        <ElCol :span="12" :xs="12"/>
      </ElRow>

    </ElForm>

  </ContentWrap>
</template>

<style>

</style>
