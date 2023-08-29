<script setup lang="tsx">
import { Descriptions } from '@/components/Descriptions'
import { useI18n } from '@/hooks/web/useI18n'
import {computed, reactive, ref} from 'vue'
import { Form } from '@/components/Form'
import { ElFormItem, ElInput, ElButton } from 'element-plus'
import ContentWrap from "@/components/ContentWrap/src/ContentWrap.vue";
import { useValidator } from '@/hooks/web/useValidator'
import { useForm } from '@/hooks/web/useForm'
import { DescriptionsSchema } from '@/components/Descriptions'
import {ApiClient} from "@/api/stub";
import api from "@/api/api";
import {useRoute} from "vue-router";
const route = useRoute();

const { required } = useValidator()

const { t } = useI18n()

const loading = ref(true)
const clientId = computed(() => route.params.id as string);
const currentRow = ref<Nullable<ApiClient>>(null)

const schema = reactive<DescriptionsSchema[]>([
  {
    field: 'clientId',
    label: t('mqtt.client.clientId')
  },
  {
    field: 'username',
    label: t('mqtt.client.username')
  },
  {
    field: 'keepAlive',
    label: t('mqtt.client.keepAlive')
  },
  {
    field: 'version',
    label: t('mqtt.client.version')
  },
  {
    field: 'willRetain',
    label: t('mqtt.client.willRetain')
  },
  {
    field: 'willQos',
    label: t('mqtt.client.willQos')
  },
  {
    field: 'willTopic',
    label: t('mqtt.client.willTopic')
  },
  {
    field: 'willPayload',
    label: t('mqtt.client.willPayload')
  },
  {
    field: 'remoteAddr',
    label: t('mqtt.client.remoteAddr')
  },
  {
    field: 'localAddr',
    label: t('mqtt.client.localAddr')
  },
  {
    field: 'subscriptionsCurrent',
    label: t('mqtt.client.subscriptionsCurrent')
  },
  {
    field: 'subscriptionsTotal',
    label: t('mqtt.client.subscriptionsTotal')
  },
  {
    field: 'packetsReceivedBytes',
    label: t('mqtt.client.packetsReceivedBytes')
  },
  {
    field: 'packetsReceivedNums',
    label: t('mqtt.client.packetsReceivedNums')
  },
  {
    field: 'packetsSendBytes',
    label: t('mqtt.client.packetsSendBytes')
  },
  {
    field: 'packetsSendNums',
    label: t('mqtt.client.packetsSendNums')
  },
  {
    field: 'messageDropped',
    label: t('mqtt.client.messageDropped')
  },
  {
    field: 'inflightLen',
    label: t('mqtt.client.inflightLen')
  },
  {
    field: 'queueLen',
    label: t('mqtt.client.queueLen')
  },
  {
    field: 'connectedAt',
    label: t('mqtt.client.connectedAt')
  },
  {
    field: 'disconnectedAt',
    label: t('mqtt.client.disconnectedAt')
  }
])

const fetch = async () => {
  loading.value = true
  const res = await api.v1.mqttServiceGetClientById(clientId.value)
      .catch(() => {
      })
      .finally(() => {
        loading.value = false
      })
  if (res) {
    currentRow.value = res.data
  } else {
    currentRow.value = null
  }
}

fetch()

</script>

<template>

  <ContentWrap v-if="currentRow">
    <Descriptions
        :title="t('mqtt.client.client')"
        :data="currentRow"
        :schema="schema"
    />
  </ContentWrap>
</template>

<style lang="less" scoped>
:deep(.is-required--item) {
  position: relative;

  &::before {
    margin-right: 4px;
    color: var(--el-color-danger);
    content: '*';
  }
}
</style>
