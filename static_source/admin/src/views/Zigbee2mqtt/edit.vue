<script setup lang="ts">
import {computed, ref, unref} from 'vue'
import {useI18n} from '@/hooks/web/useI18n'
import {ElButton, ElPopconfirm, ElTabs, ElTabPane} from 'element-plus'
import {useForm} from '@/hooks/web/useForm'
import {useCache} from '@/hooks/web/useCache'
import {useAppStore} from '@/store/modules/app'
import {usePermissionStore} from '@/store/modules/permission'
import {useRoute, useRouter} from 'vue-router'
import {useValidator} from '@/hooks/web/useValidator'
import api from "@/api/api";
import Form from './components/Form.vue'
import {ApiZigbee2Mqtt} from "@/api/stub";
import {ContentWrap} from "@/components/ContentWrap";
import Devices from "@/views/Zigbee2mqtt/components/Devices.vue";

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
const bridgeId = computed(() => route.params.id as number);
const currentBridge = ref<Nullable<ApiZigbee2Mqtt>>(null)
const activeTab = ref('devices')

const fetch = async () => {
  loading.value = true
  const res = await api.v1.zigbee2MqttServiceGetZigbee2MqttBridge(bridgeId.value)
      .catch(() => {
      })
      .finally(() => {
        loading.value = false
      })
  if (res) {
    currentBridge.value = res.data
  } else {
    currentBridge.value = null
  }
}

const save = async () => {
  const write = unref(writeRef)
  const validate = await write?.elFormRef?.validate()?.catch(() => {
  })
  if (validate) {
    loading.value = true
    const data = (await write?.getFormData())
    const body = {
      name: data.name,
      login: data.login,
      password: data.password,
      permitJoin: data.permitJoin,
      baseTopic: data.baseTopic,
    }
    const res = await api.v1.zigbee2MqttServiceUpdateBridgeById(bridgeId.value as string, body)
        .catch(() => {
        })
        .finally(() => {
          loading.value = false
        })
    if (res) {
      cancel()
    }
  }
}

const cancel = () => {
  push('/etc/zigbee2mqtt')
}

const remove = async () => {
  loading.value = true
  const res = await api.v1.zigbee2MqttServiceDeleteBridgeById(bridgeId.value as string)
      .catch(() => {
      })
      .finally(() => {
        loading.value = false
      })
  if (res) {
    cancel()
  }
}
fetch()

</script>

<template>
  <ContentWrap>

    <el-tabs class="demo-tabs" v-model="activeTab">
      <!-- main -->
      <el-tab-pane :label="$t('zigbee2mqtt.main')" name="main">
        <Form ref="writeRef" :current-row="currentBridge"/>
      </el-tab-pane>
      <!-- /main -->

      <!-- devices -->
      <el-tab-pane :label="$t('zigbee2mqtt.devices')" name="devices">
        <Devices v-model="currentBridge"/>
      </el-tab-pane>
      <!-- /devices -->
    </el-tabs>

    <div style="text-align: right">

      <ElButton type="primary" @click="save()">
        {{ t('main.save') }}
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

</template>

<style lang="less" scoped>

</style>
