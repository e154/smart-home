<script setup lang="ts">
import {computed, h, onMounted, onUnmounted, PropType, reactive, ref, unref, watch} from "vue";
import {Table} from '@/components/Table'
import {Card, CardItem, Core, Tab} from "@/views/Dashboard/core";
import {FormSchema} from "@/types/form";
import {Pagination, TableColumn} from "@/types/table";
import {ApiLog} from "@/api/stub";
import {parseTime} from "@/utils";
import {UUID} from "uuid-generator-ts";
import stream from "@/api/stream";
import {useI18n} from "@/hooks/web/useI18n";
import api from "@/api/api";
import {debounce} from "lodash-es";

const {t} = useI18n()

// ---------------------------------
// common
// ---------------------------------

const props = defineProps({
  item: {
    type: Object as PropType<Nullable<CardItem>>,
    default: () => null
  },
})

const el = ref(null)
onMounted(() => {
  // store dom element moveable
  props.item.setTarget(el.value)

  const uuid = new UUID()
  currentID.value = uuid.getDashFreeUUID()

  setTimeout(() => {
    stream.subscribe('log', currentID.value, onLogs)
  }, 1000)
})

onUnmounted(() => {
  stream.unsubscribe('log', currentID.value)
})


// ---------------------------------
// component methods
// ---------------------------------
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
    width: "170px",
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
    width: "170px"
  },
]
const paginationObj = ref<Pagination>({
  currentPage: 1,
  pageSize: 20,
  total: 0,
  pageSizes: [20, 50, 100, 150, 250],
})
const currentID = ref('')

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

const getList = debounce( async () => {
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
    paginationObj.value.currentPage = meta.pagination.page;
    paginationObj.value.total = meta.pagination.total;
  } else {
    tableObject.tableList = [];
  }
}, 100)


const onLogs = (log: ApiLog) => {
  getList()
}

const sortChange = (data) => {
  const {column, prop, order} = data;
  const pref: string = order === 'ascending' ? '+' : '-'
  tableObject.sort = pref + prop
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

getList()

</script>

<template>
  <div ref="el" style="overflow:hidden;overflow-y: scroll">

    <Table
        v-model:pageSize="paginationObj.pageSize"
        v-model:currentPage="paginationObj.currentPage"
        :columns="columns"
        :data="tableObject.tableList"
        :pagination="paginationObj"
        @sort-change="sortChange"
        :row-class-name="tableRowClassName"
        style="width: 100%"
        class="logsTable"
        :selection="false"
    />

  </div>
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
