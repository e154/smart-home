<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {Table} from '@/components/Table'
import {h, onMounted, onUnmounted, reactive, ref, watch} from 'vue'
import {Pagination, TableColumn} from '@/types/table'
import api from "@/api/api";
import {ApiClient} from "@/api/stub";
import {ContentWrap} from "@/components/ContentWrap";
import {parseTime} from "@/utils";
import {UUID} from "uuid-generator-ts";
import stream from "@/api/stream";

const {t} = useI18n()

interface TableObject {
  tableList: ApiClient[]
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

const currentID = ref('')

const onEventMqttNewClient = () => {
  getList()
}

onMounted(() => {
  const uuid = new UUID()
  currentID.value = uuid.getDashFreeUUID()

  setTimeout(() => {
    stream.subscribe('event_mqtt_new_client', currentID.value, onEventMqttNewClient);
  }, 1000)
})

onUnmounted(() => {
  stream.unsubscribe('event_mqtt_new_client', currentID.value);
})

const columns: TableColumn[] = [
  {
    field: 'clientId',
    label: t('mqtt.client.clientId'),
    width: "200px"
  },
  {
    field: 'username',
    label: t('mqtt.client.username'),
  },
  {
    field: 'connectedAt',
    label: t('mqtt.client.connectedAt'),
    type: 'time',
    sortable: true,
    width: "170px",
    formatter: (row: ApiClient) => {
      return h(
          'span',
          parseTime(row?.connectedAt)
      )
    }
  },
]
const paginationObj = ref<Pagination>({
  currentPage: 1,
  pageSize: 50,
  total: 0,
  pageSizes: [50, 100, 150, 250],
})

const getList = async () => {
  tableObject.loading = true

  let params: Params = {
    page: paginationObj.value.currentPage,
    limit: paginationObj.value.pageSize,
    sort: tableObject.sort,
  }

  const res = await api.v1.mqttServiceGetClientList(params)
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
      getList()
    }
)

watch(
    () => paginationObj.value.pageSize,
    () => {
      getList()
    }
)

const sortChange = (data) => {
  const {column, prop, order} = data;
  const pref: string = order === 'ascending' ? '+' : '-'
  tableObject.sort = pref + prop
  getList()
}


const selectRow = (row) => {
  if (!row) {
    return
  }
  const {clientId} = row
  push(`/etc/mqtt/client/${clientId}`)
}

getList()

</script>

<template>
  <ContentWrap>

    <Table
        :selection="false"
        v-model:pageSize="paginationObj.pageSize"
        v-model:currentPage="paginationObj.currentPage"
        :columns="columns"
        :data="tableObject.tableList"
        :loading="tableObject.loading"
        :pagination="paginationObj"
        @sort-change="sortChange"
        style="width: 100%"
        :showUpPagination="20"
        @current-change="selectRow"
    />


  </ContentWrap>

</template>

<style lang="less">

</style>
