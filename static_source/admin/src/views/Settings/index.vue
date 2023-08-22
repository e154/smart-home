<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {ref} from 'vue'
import {useAppStore} from "@/store/modules/app";
import api from "@/api/api";
import {ElCol, ElDivider, ElForm, ElFormItem, ElRow} from 'element-plus'
import {ApiDashboard} from '@/api/stub';
import DashboardSearch from "@/views/Dashboard/components/DashboardSearch.vue";
import ContentWrap from "@/components/ContentWrap/src/ContentWrap.vue";

const appStore = useAppStore()

const {t} = useI18n()

export interface Settings {
  mainDashboardDark?: ApiDashboard;
  mainDashboardLight?: ApiDashboard;
  devDashboardDark?: ApiDashboard;
  devDashboardLight?: ApiDashboard;
}

const settings = ref<Settings>({} as Settings)
const loading = ref(true)

const getVariableByName = async (name: string) => {
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

const getSettings = async () => {
  loading.value = true
  await Promise.all([
    getVariableByName('mainDashboardDark'),
    getVariableByName('mainDashboardLight'),
    getVariableByName('devDashboardDark'),
    getVariableByName('devDashboardLight')
  ])
  loading.value = false
}

const changedVariable = (name: string) => {
  api.v1.variableServiceUpdateVariable(name, {value: settings.value[name]?.id + '' || ''})
}

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
    </ElForm>

  </ContentWrap>
</template>

<style>

</style>
