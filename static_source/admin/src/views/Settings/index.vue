<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {h, reactive, ref, watch} from 'vue'
import {useAppStore} from "@/store/modules/app";
import api from "@/api/api";
import {ElDivider, ElCol, ElIcon, ElRow, ElFormItem, ElForm} from 'element-plus'
import {ApiDashboard, ApiDashboardShort, ApiVariable} from '@/api/stub';
import DashboardSearch from "@/views/Dashboard/components/DashboardSearch.vue";
import ContentWrap from "@/components/ContentWrap/src/ContentWrap.vue";

const appStore = useAppStore()

const {t} = useI18n()

export interface Settings {
  mainDashboard?: ApiDashboard;
  devDashboard?: ApiDashboard;
}

const settings = ref<Settings>({} as Settings)
const loading = ref(true)

const getSettings = async () => {
  loading.value = true
  api.v1.variableServiceGetVariableByName('mainDashboard').then((resp) => {
    if (!resp.data?.value || resp.data?.value == 'undefined') {
      return;
    }
   const id: number = parseInt(resp.data?.value);
   api.v1.dashboardServiceGetDashboardById(id).then((resp) => {
     settings.value.mainDashboard = resp.data
   });
  });

  api.v1.variableServiceGetVariableByName('devDashboard').then((resp) => {
    if (!resp.data?.value || resp.data?.value == 'undefined') {
      return;
    }
    const id: number = parseInt(resp.data?.value);
    api.v1.dashboardServiceGetDashboardById(id).then((resp) => {
      settings.value.devDashboard = resp.data
    });
  });

  loading.value = false
}

const updateVariable = (name: string, value: string) => {
  api.v1.variableServiceUpdateVariable(name, {value: value})
}

const changedMainDashboard = (values: ApiDashboard, event: any) => {
    updateVariable("mainDashboard", values?.id + '' || '')
}

const changedDevDashboard = (values: ApiDashboard, event: any) => {
    updateVariable("devDashboard", values?.id + '' || '')
}

getSettings()

</script>

<template>

  <ContentWrap>

    <ElDivider content-position="left">{{$t('settings.dashboardOptions')}}</ElDivider>

    <ElForm label-position="top">
      <ElRow :gutter="24">
        <ElCol :span="12" :xs="12">
          <ElFormItem :label="$t('settings.mainDashboard')" prop="mainDashboard">
            <DashboardSearch v-model="settings.mainDashboard" @update:modelValue="changedMainDashboard"/>
          </ElFormItem>
        </ElCol>
        <ElCol :span="12" :xs="12">
          <ElFormItem :label="$t('settings.devDashboard')" prop="devDashboard">
            <DashboardSearch v-model="settings.devDashboard" @update:modelValue="changedDevDashboard"/>
          </ElFormItem>
        </ElCol>
      </ElRow>
    </ElForm>

  </ContentWrap>
</template>

<style>

</style>
