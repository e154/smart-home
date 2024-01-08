<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {Table} from '@/components/Table'
import {h, reactive, ref, watch} from 'vue'
import {useAppStore} from "@/store/modules/app";
import {Pagination, TableColumn} from '@/types/table'
import api from "@/api/api";
import {ElButton, ElSwitch} from 'element-plus'
import {ApiTask, ApiZigbee2Mqtt} from "@/api/stub";
import {useForm} from "@/hooks/web/useForm";
import {useRouter} from "vue-router";
import {parseTime} from "@/utils";
import ContentWrap from "@/components/ContentWrap/src/ContentWrap.vue";

const {push, currentRoute} = useRouter()
const remember = ref(false)
const {register, elFormRef, methods} = useForm()
const appStore = useAppStore()
const {t} = useI18n()

interface TableObject {
  tableList: ApiZigbee2Mqtt[]
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
      sort: '-id'
    },
);

const columns: TableColumn[] = [
  {
    field: 'id',
    label: t('zigbee2mqtt.id'),
    width: "60px",
    sortable: true
  },
  {
    field: 'name',
    label: t('zigbee2mqtt.name'),
    width: "170px",
    sortable: true
  },
  {
    field: 'login',
    label: t('zigbee2mqtt.login'),
    sortable: true
  },
  {
    field: 'permitJoin',
    label: t('zigbee2mqtt.permitJoin'),
    sortable: true,
    width: "170px"
  },
  {
    field: 'createdAt',
    label: t('main.createdAt'),
    type: 'time',
    sortable: true,
    width: "170px",
    formatter: (row: ApiTask) => {
      return h(
          'span',
          parseTime(row.createdAt)
      )
    }
  },
  {
    field: 'updatedAt',
    label: t('main.updatedAt'),
    type: 'time',
    sortable: true,
    width: "170px",
    formatter: (row: ApiTask) => {
      return h(
          'span',
          parseTime(row.updatedAt)
      )
    }
  },
]
const paginationObj = ref<Pagination>({
  currentPage: 1,
  pageSize: 50,
  total: 0,
})
const currentID = ref('')

const getList = async () => {
  tableObject.loading = true

  let params: Params = {
    page: paginationObj.value.currentPage,
    limit: paginationObj.value.pageSize,
    sort: tableObject.sort,
  }

  const res = await api.v1.zigbee2MqttServiceGetBridgeList(params)
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

getList()

const addNew = () => {
  push('/etc/zigbee2mqtt/new')
}

const selectRow = (row) => {
  if (!row) {
    return
  }
  const {id} = row
  push(`/etc/zigbee2mqtt/edit/${id}`)
}

</script>

<template>
  <ContentWrap>
    <ElButton class="flex mb-20px items-left" type="primary" @click="addNew()" plain>
      <Icon icon="ep:plus" class="mr-5px"/>
      {{ t('zigbee2mqtt.addNew') }}
    </ElButton>
    <Table
        :selection="false"
        v-model:pageSize="paginationObj.pageSize"
        v-model:currentPage="paginationObj.currentPage"
        :showUpPagination="20"
        :columns="columns"
        :data="tableObject.tableList"
        :loading="tableObject.loading"
        :pagination="paginationObj"
        @sort-change="sortChange"
        style="width: 100%"
        @current-change="selectRow"
    >
      <template #permitJoin="{ row }">
        <ElSwitch v-model="row.permitJoin" disabled/>
      </template>
    </Table>
  </ContentWrap>
</template>

<style lang="less">

.el-table__row {
  cursor: pointer;
}
</style>
