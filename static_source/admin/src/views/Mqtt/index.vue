<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {Table} from '@/components/Table'
import {computed, h, onMounted, onUnmounted, reactive, ref, watch} from 'vue'
import {useAppStore} from "@/store/modules/app";
import {Pagination, TableColumn} from '@/types/table'
import api from "@/api/api";
import {ElButton, ElMessage, ElTag} from 'element-plus'
import {ApiClient} from "@/api/stub";
import {useForm} from "@/hooks/web/useForm";
import {useRouter} from "vue-router";
import ContentWrap from "@/components/ContentWrap/src/ContentWrap.vue";
import {parseTime} from "@/utils";
import { Dialog } from '@/components/Dialog'
import {useEmitt} from "@/hooks/web/useEmitt";
import {EventStateChange, EventTriggerCompleted} from "@/api/stream_types";
import {UUID} from "uuid-generator-ts";
import stream from "@/api/stream";

const {push, currentRoute} = useRouter()
const remember = ref(false)
const {register, elFormRef, methods} = useForm()
const appStore = useAppStore()
const {t} = useI18n()
const isMobile = computed(() => appStore.getMobile)

const dialogSource = ref("")
const dialogVisible = ref(false)
const importedTask = ref("")

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
