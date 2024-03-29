<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {Table} from '@/components/Table'
import {h, onMounted, onUnmounted, reactive, ref, watch} from 'vue'
import {useAppStore} from "@/store/modules/app";
import {Pagination, TableColumn} from '@/types/table'
import api from "@/api/api";
import {ElButton, ElTag} from 'element-plus'
import {ApiScript, ApiStatistics} from "@/api/stub";
import {useForm} from "@/hooks/web/useForm";
import {useRouter} from "vue-router";
import {parseTime} from "@/utils";
import {ContentWrap} from "@/components/ContentWrap";
import Statistics from "@/components/Statistics/Statistics.vue";
import {FormSchema} from "@/types/form";
import {Form} from '@/components/Form'
import {UUID} from "uuid-generator-ts";
import stream from "@/api/stream";
import {useCache} from "@/hooks/web/useCache";
import {Infotip} from "@/components/Infotip";

const {push, currentRoute} = useRouter()
const remember = ref(false)
const statistic = ref<Nullable<ApiStatistics>>(null)
const {register, elFormRef, methods} = useForm()
const appStore = useAppStore()
const {t} = useI18n()
const {wsCache} = useCache()

const cachePref = 'scripts'

interface TableObject {
  tableList: ApiScript[]
  params?: any
  loading: boolean
  sort?: string
  query?: string
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
      sort: wsCache.get(cachePref + 'Sort') || '-createdAt',
      query: wsCache.get(cachePref + 'Query'),
    },
);

const columns: TableColumn[] = [
  {
    field: 'id',
    label: t('scripts.id'),
    width: "60px",
    sortable: true
  },
  {
    field: 'name',
    label: t('scripts.name'),
    width: "170px",
    sortable: true
  },
  {
    field: 'lang',
    label: t('scripts.lang'),
    width: "200px",
    sortable: true
  },
  {
    field: 'description',
    label: t('scripts.description'),
    sortable: true
  },
  {
    field: 'createdAt',
    label: t('main.createdAt'),
    type: 'time',
    width: "170px",
    formatter: (row: ApiScript) => {
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
    formatter: (row: ApiScript) => {
      return h(
          'span',
          parseTime(row.updatedAt)
      )
    }
  },
]
const paginationObj = ref<Pagination>({
  currentPage: wsCache.get(cachePref + 'CurrentPage') || 1,
  pageSize: wsCache.get(cachePref + 'PageSize') || 50,
  total: 0,
})
const currentID = ref('')

const getStatistic = async () => {

  const res = await api.v1.scriptServiceGetStatistic()
      .catch(() => {
      })
      .finally(() => {

      })
  if (res) {
    statistic.value = res.data
  }
}

const getList = async () => {
  getStatistic()

  tableObject.loading = true

  wsCache.set(cachePref + 'CurrentPage', paginationObj.value.currentPage)
  wsCache.set(cachePref + 'PageSize', paginationObj.value.pageSize)
  wsCache.set(cachePref + 'Sort', tableObject.sort)
  wsCache.set(cachePref + 'Query', tableObject.query)

  let params: Params = {
    page: paginationObj.value.currentPage,
    limit: paginationObj.value.pageSize,
    sort: tableObject.sort,
    query: tableObject.query || undefined,
  }

  const res = await api.v1.scriptServiceGetScriptList(params)
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

onMounted(() => {
  const uuid = new UUID()
  currentID.value = uuid.getDashFreeUUID()

  setTimeout(() => {
    stream.subscribe('event_removed_script_model', currentID.value, getList)
    stream.subscribe('event_created_script_model', currentID.value, getList)
  }, 1000)
})

onUnmounted(() => {
  stream.unsubscribe('event_removed_script_model', currentID.value)
  stream.unsubscribe('event_created_script_model', currentID.value)
})

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
  push('/scripts/new')
}

const selectRow = (row) => {
  if (!row) {
    return
  }
  const {id} = row
  push(`/scripts/edit/${id}`)
}

// search form
const schema = reactive<FormSchema[]>([
  {
    field: 'name',
    label: t('scripts.search'),
    component: 'Input',
    colProps: {
      span: 24
    },
    componentProps: {
      placeholder: t('scripts.search'),
      onChange: (val: string) => {
        tableObject.query = val || undefined
        getList()
      }
    }
  },
])


const {setValues, setSchema} = methods
if (wsCache.get(cachePref + 'Query')) {
  setValues({
    name: wsCache.get(cachePref + 'Query')
  })
}

</script>

<template>

  <Statistics v-model="statistic" :cols="6"/>

  <ContentWrap>
    <ElButton class="flex mb-20px items-left" type="primary" @click="addNew()" plain>
      <Icon icon="ep:plus" class="mr-5px"/>
      {{ t('scripts.addNew') }}
    </ElButton>

    <Infotip
        :show-index="false"
        title="INFO"
        :schema="[
      {
        label: t('scripts.info4'),
      },
    ]"
    />

    <Form
        id="search-form-scripts"
        :schema="schema"
        label-position="top"
        label-width="auto"
        hide-required-asterisk
        @register="register"
    />
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
      <template #lang="{row}">
        <ElTag>
          {{ row.lang }}
        </ElTag>
      </template>
    </Table>
  </ContentWrap>
</template>

<style lang="less" scoped>

.el-table__row {
  cursor: pointer;
}

:deep(#search-form-scripts .el-col) {
  padding-left: 0 !important;
  padding-right: 0 !important;
}

</style>
