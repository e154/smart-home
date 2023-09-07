<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {Table} from '@/components/Table'
import {Form} from '@/components/Form'
import {computed, h, onDeactivated, onMounted, onUnmounted, reactive, ref, watch} from 'vue'
import {FormSchema} from "@/types/form";
import {useAppStore} from "@/store/modules/app";
import {Pagination, TableColumn} from '@/types/table'
import api from "@/api/api";
import {UUID} from 'uuid-generator-ts'
import { ElTable, ElTableColumn, ElPagination, ElDatePicker  } from 'element-plus'
import {ApiLog} from "@/api/stub";
import {useForm} from "@/hooks/web/useForm";
import {parseTime} from "@/utils";
import stream from "@/api/stream";
import ContentWrap from "@/components/ContentWrap/src/ContentWrap.vue";

const remember = ref(false)
const {register, elFormRef, methods} = useForm()
const appStore = useAppStore()
const {t} = useI18n()
const isMobile = computed(() => appStore.getMobile)

interface TableObject {
  tableList: ApiLog[]
  params?: any
  loading: boolean
  sort?: string
  query?: string
  startDate?: string;
  endDate?: string;
}

interface Params {
  page?: number;
  limit?: number;
  sort?: string;
  query?: string;
  startDate?: string;
  endDate?: string;
}

const tableObject = reactive<TableObject>(
    {
      tableList: [],
      loading: false,
    }
);

const schema = reactive<FormSchema[]>([
  {
    field: 'dateTime',
    component: 'DatePicker',
    label: t('logs.dateTimerange'),
    colProps: {
      span: 24
    },
    value: [],
    componentProps: {
      type: 'datetimerange',
      shortcuts: [
        {
          text: t('logs.today'),
          value: () => {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - start.getHours() * 3600 * 1000 - start.getMinutes() * 60 * 1000 - start.getSeconds() * 1000)
            return [start, end]
          }
        },
        {
          text: t('logs.yesterday'),
          value: () => {
            const end = new Date()
            const start = new Date()
            end.setHours(0)
            end.setMinutes(0)
            end.setSeconds(0)
            start.setTime(end.getTime() - 3600 * 1000 * 24)
            return [start, end]
          }
        },
        {
          text: t('logs.aWeekAgo'),
          value: () => {
            const end = new Date()
            const start = new Date()
            end.setHours(0)
            end.setMinutes(0)
            end.setSeconds(0)
            start.setTime(end.getTime() - 3600 * 1000 * 24 * 7)
            return [start, end]
          }
        }
      ],
      onChange: (val: Date[]) => {
        if (val && val.length > 1) {
          tableObject.startDate = val[0].toISOString()
          tableObject.endDate = val[1].toISOString()
        } else {
          tableObject.startDate = undefined
          tableObject.endDate = undefined
        }
        getList()
      }
    },
  },
  {
    field: 'levelList',
    label: t('logs.level'),
    component: 'CheckboxButton',
    value: [],
    colProps: {
      span: 24
    },
    componentProps: {
      options: [
        {
          label: 'Emergency',
          value: 'Emergency'
        },
        {
          label: 'Alert',
          value: 'Alert'
        },
        {
          label: 'Critical',
          value: 'Critical'
        },
        {
          label: 'Error',
          value: 'Error'
        },
        {
          label: 'Warning',
          value: 'Warning'
        },
        {
          label: 'Notice',
          value: 'Notice'
        },
        {
          label: 'Info',
          value: 'Info'
        },
        {
          label: 'Debug',
          value: 'Debug'
        }
      ]
    }
  },
])

const columns: TableColumn[] = [
  {
    field: 'createdAt',
    label: t('main.createdAt'),
    type: 'time',
    sortable: true,
    width: "150px",
    formatter: (row: ApiLog) => {
      return h(
          'span',
          parseTime(row.createdAt)
      )
    }
  },
  {
    field: 'level',
    label: t('logs.level'),
    sortable: true,
    width: "110px"
  },
  {
    field: 'body',
    label: t('logs.body')
  },
  {
    field: 'owner',
    label: t('logs.owner'),
    sortable: true,
    width: "150px"
  },
]
const paginationObj = ref<Pagination>({
  currentPage: 1,
  pageSize: 100,
  total: 0,
  pageSizes: [50, 100, 150, 250],
})
const currentID = ref('')

onMounted(() => {
  const uuid = new UUID()
  currentID.value = uuid.getDashFreeUUID()

  setTimeout(() => {
    stream.subscribe('log', currentID.value, onLogs)
  }, 1000)
})

onUnmounted(() => {
  stream.unsubscribe('log', currentID.value)
})

const getList = async () => {
  tableObject.loading = true

  let params: Params = {
    page: paginationObj.value.currentPage,
    limit: paginationObj.value.pageSize,
    sort: tableObject.sort,
    query: tableObject.query,
    startDate: tableObject.startDate,
    endDate: tableObject.endDate,
  }

  const res = await api.v1.logServiceGetLogList(params)
      .catch(() => {
      })
      .finally(() => {
        tableObject.loading = false
      })
  if (res) {
    const {items, meta} = res.data;
    tableObject.tableList = items;
    // paginationObj.value.pageSize = meta.limit;
    paginationObj.value.currentPage = meta.page;
    paginationObj.value.total = meta.total;
  } else {
    tableObject.tableList = [];
  }
}

const onLogs = (log: ApiLog) => {
  getList()
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

const onFormChange = async () => {
  const {getFormData} = methods
  const formData = await getFormData()
  const {levelList} = formData
  if (levelList && levelList.length > 0) {
    tableObject.query = levelList.join(',')
  } else {
    tableObject.query = undefined
  }
  getList()
}

const tableRowClassName = (data) => {
  const { row, rowIndex } = data
  let style = ''
  switch (row.level) {
    case 'Emergency':
      style = 'log-emergency'
      break
    case 'Alert':
      style = 'log-alert'
      break
    case 'Critical':
      style = 'log-critical'
      break
    case 'Error':
      style = 'log-error'
      break
    case 'Warning':
      style = 'log-warning'
      break
    case 'Notice':
      style = 'log-notice'
      break
    case 'Info':
      style = 'log-info'
      break
    case 'Debug':
      style = 'log-debug'
      break
  }
  return style
}

getList()

</script>

<template>
  <ContentWrap>
    <Form
        :schema="schema"
        label-position="top"
        label-width="auto"
        hide-required-asterisk
        @change="onFormChange"
        @register="register"
    />

    <Table
        v-model:pageSize="paginationObj.pageSize"
        v-model:currentPage="paginationObj.currentPage"
        :columns="columns"
        :data="tableObject.tableList"
        :loading="tableObject.loading"
        :pagination="paginationObj"
        @sort-change="sortChange"
        :row-class-name="tableRowClassName"
        style="width: 100%"
        class="logsTable"
        :selection="false"
        :showUpPagination="20"
    />
  </ContentWrap>

</template>

<style lang="less">

.logsTable {

.el-table__row {

td.el-table__cell {
  padding: 0;
  border-bottom: none !important;
}

&.log-emergency {
  --el-table-tr-bg-color: var(--el-color-error-light-5);
}

&.log-alert {
  --el-table-tr-bg-color: var(--el-color-error-light-5);
}

&.log-critical {
  --el-table-tr-bg-color: var(--el-color-error-light-5);
}

&.log-error {
  --el-table-tr-bg-color: var(--el-color-error-light-5);
}

&.log-warning {
  --el-table-tr-bg-color: var(--el-color-warning-light-5);
}

&.log-notice {
  --el-table-tr-bg-color: var(--el-color-success-light-5);
}

&.log-info {
  background-color: inherit;
}

&.log-debug {
  //background-color: #82aeff;
  --el-table-tr-bg-color: var(--el-color-info-light-5);
}

}
}
</style>
