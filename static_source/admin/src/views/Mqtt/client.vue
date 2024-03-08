<script setup lang="tsx">
import { Descriptions } from '@/components/Descriptions'
import { useI18n } from '@/hooks/web/useI18n'
import {computed, h, reactive, ref, watch} from 'vue'
import { Form } from '@/components/Form'
import { ElFormItem, ElInput, ElButton } from 'element-plus'
import {ContentWrap} from "@/components/ContentWrap";
import { useValidator } from '@/hooks/web/useValidator'
import { useForm } from '@/hooks/web/useForm'
import { DescriptionsSchema } from '@/components/Descriptions'
import {ApiClient, ApiSubscription} from "@/api/stub";
import api from "@/api/api";
import {useRoute} from "vue-router";
import {parseTime} from "@/utils";
import {Table} from '@/components/Table'
import {Pagination, TableColumn} from '@/types/table'

const route = useRoute();

const { required } = useValidator()

const { t } = useI18n()

const loading = ref(true)
const clientId = computed(() => route.params.id as string);
const client = ref<Nullable<ApiClient>>(null)

const clientSchema = reactive<DescriptionsSchema[]>([
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
    client.value = res.data
  } else {
    client.value = null
  }
}

// ------------------------------------------------
// subscriptions
// ------------------------------------------------

interface TableObject {
  tableList: ApiSubscription[]
  params?: any
  loading: boolean
  sort?: string
}

interface Params {
  page?: number;
  limit?: number;
  sort?: string;
}

const tableObject = reactive<TableObject>(
    {
      tableList: [],
      loading: false,
    }
);

const columns: TableColumn[] = [
  {
    field: 'id',
    label: t('mqtt.subscription.id'),
    width: "100px"
  },
  {
    field: 'topicName',
    label: t('mqtt.subscription.topicName'),
  },
  {
    field: 'name',
    label: t('mqtt.subscription.name'),
  },
  {
    field: 'qos',
    label: t('mqtt.subscription.qos'),
  },
  {
    field: 'noLocal',
    label: t('mqtt.subscription.noLocal'),
  },
  {
    field: 'retainAsPublished',
    label: t('mqtt.subscription.retainAsPublished'),
  },
  {
    field: 'retainHandling',
    label: t('mqtt.subscription.retainHandling'),
  },
]

const paginationObj = ref<Pagination>({
  currentPage: 1,
  pageSize: 50,
  total: 0,
  pageSizes: [50, 100, 150, 250],
})

const getSubscriptionList = async () => {
  tableObject.loading = true

  let params: Params = {
    page: paginationObj.value.currentPage,
    limit: paginationObj.value.pageSize,
    sort: tableObject.sort,
    clientId: clientId.value,
  }

  const res = await api.v1.mqttServiceGetSubscriptionList(params)
      .catch(() => {
      })
      .finally(() => {
        tableObject.loading = false
      })
  if (res) {
    const {items, meta} = res.data;
    tableObject.tableList = items;
    paginationObj.value.currentPage = meta.pagination.page;
    paginationObj.value.total = meta.pagination.total;
  } else {
    tableObject.tableList = [];
  }
}

watch(
    () => paginationObj.value.currentPage,
    () => {
      getSubscriptionList()
    }
)

watch(
    () => paginationObj.value.pageSize,
    () => {
      getSubscriptionList()
    }
)

const sortChange = (data) => {
  const {column, prop, order} = data;
  const pref: string = order === 'ascending' ? '+' : '-'
  tableObject.sort = pref + prop
  getSubscriptionList()
}


fetch()

getSubscriptionList()

</script>

<template>

  <ContentWrap v-if="client">
    <Descriptions
        :title="t('mqtt.client.client')"
        :data="client"
        :schema="clientSchema"
    />

    <h2 class="mt-20px"><strong>{{ t('mqtt.client.subscriptions') }}</strong></h2>

    <Table
        class="mt-20px"
        :selection="false"
        v-model:pageSize="paginationObj.pageSize"
        v-model:currentPage="paginationObj.currentPage"
        :columns="columns"
        :data="tableObject.tableList"
        :loading="tableObject.loading"
        :pagination="paginationObj"
        style="width: 100%"
        :showUpPagination="20"
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
