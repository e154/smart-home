<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {Table} from '@/components/Table'
import {h, onMounted, onUnmounted, reactive, ref, watch} from 'vue'
import {Pagination, TableColumn} from '@/types/table'
import api from "@/api/api";
import {ElButton, ElTag} from 'element-plus'
import {ApiEntityShort, ApiVariable} from "@/api/stub";
import {useRouter} from "vue-router";
import {ContentWrap} from "@/components/ContentWrap";
import {useCache} from "@/hooks/web/useCache";
import {UUID} from "uuid-generator-ts";
import stream from "@/api/stream";
import {EventStateChange} from "@/api/types";
import {parseTime} from "@/utils";

const {push} = useRouter()
const {t} = useI18n()
const {wsCache} = useCache()

interface TableObject {
  tableList: ApiVariable[]
  params?: any
  loading: boolean
  sort?: string
}

interface Params {
  page?: number;
  limit?: number;
  sort?: string;
}

const cachePref = 'variables'
const tableObject = reactive<TableObject>(
    {
      tableList: [],
      loading: false,
      sort: wsCache.get(cachePref + 'Sort') || '-createdAt'
    }
);

const columns: TableColumn[] = [
  {
    field: 'name',
    label: t('variables.name'),
    sortable: true,
    width: "180px"
  },
  {
    field: 'value',
    label: t('variables.value')
  },
  {
    field: 'createdAt',
    label: t('main.createdAt'),
    type: 'time',
    sortable: true,
    width: "170px",
    formatter: (row: ApiEntityShort) => {
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
    formatter: (row: ApiEntityShort) => {
      return h(
          'span',
          parseTime(row.updatedAt)
      )
    }
  }
]
const paginationObj = ref<Pagination>({
  currentPage: wsCache.get(cachePref + 'CurrentPage') || 1,
  pageSize: wsCache.get(cachePref + 'PageSize') || 50,
  total: 0,
  pageSizes: [50, 100, 150, 250],
})

const getList = async () => {
  tableObject.loading = true

  wsCache.set(cachePref + 'CurrentPage', paginationObj.value.currentPage)
  wsCache.set(cachePref + 'PageSize', paginationObj.value.pageSize)
  wsCache.set(cachePref + 'Sort', tableObject.sort)

  let params: Params = {
    page: paginationObj.value.currentPage,
    limit: paginationObj.value.pageSize,
    sort: tableObject.sort,
  }

  const res = await api.v1.variableServiceGetVariableList(params)
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
  push('/etc/variables/new')
}

const selectRow = (row) => {
  if (!row) {
    return
  }
  const {name} = row
  push(`/etc/variables/edit/${name}`)
}

const currentID = ref('')

const onStateChanged = (event: EventStateChange) => {
  getList()
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

</script>

<template>
  <ContentWrap>
    <ElButton class="flex mb-20px items-left" type="primary" @click="addNew()" plain>
      <Icon icon="ep:plus" class="mr-5px"/>
      {{ t('variables.addNew') }}
    </ElButton>
    <Table
        class="variables-table"
        :expand="true"
        :selection="false"
        v-model:pageSize="paginationObj.pageSize"
        v-model:currentPage="paginationObj.currentPage"
        :columns="columns"
        :data="tableObject.tableList"
        :loading="tableObject.loading"
        :pagination="paginationObj"
        @sort-change="sortChange"
        style="width: 100%"
        @current-change="selectRow"
        :showUpPagination="20"
    >
      <template #expand="{row}">
        <div class="tag-list" v-if="row.tags">
          <ElTag v-for="tag in row.tags" type="info" :key="tag" round effect="light" size="small">
            {{ tag }}
          </ElTag>
        </div>
      </template>
    </Table>
  </ContentWrap>

</template>

<style lang="less">

.el-table__row {
  cursor: pointer;
}

.variables-table {
  .tag-list {
    .el-tag {
      margin: 0 5px;
    }
  }

  :deep(.el-table__row) {
    cursor: pointer;
  }

  tr.el-table__row [class*="el-table__cell"] {
  //background-color: green; border-top: var(--el-table-border); border-bottom: none !important;
    border-top: var(--el-table-border);
  }

  .el-table__expanded-cell {
    &.el-table__cell [class*="tag-list"] {
    //background-color: red!important; border-bottom: none !important;
    }

    &.el-table__cell:not(:has(.tag-list)) {
      display: none !important;
    //background-color: blue!important;
    }
  }

  .el-table td.el-table__cell,
  .el-table th.el-table__cell.is-leaf {
    border-bottom: none !important;
  }

  .el-table--enable-row-hover .el-table__body tr.el-table__row:hover,
  .el-table--enable-row-hover .el-table__body tr.el-table__row:hover + tr {
    & > td.el-table__cell {
      background-color: var(--el-table-row-hover-bg-color);
    }
  }
}

</style>
