<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {Table} from '@/components/Table'
import {Form} from '@/components/Form'
import {h, reactive, ref, watch} from 'vue'
import {FormSchema} from "@/types/form";
import {Pagination, TableColumn} from '@/types/table'
import {ElButton} from 'element-plus'
import api from "@/api/api";
import {ApiMessage, ApiMessageDelivery} from "@/api/stub";
import {useForm} from "@/hooks/web/useForm";
import {parseTime} from "@/utils";
import {ContentWrap} from "@/components/ContentWrap";
import AttributesViewer from "@/views/MessageDelivery/components/AttributesViewer.vue";
import {Dialog} from '@/components/Dialog'

const {register, methods} = useForm()
const {t} = useI18n()

interface TableObject {
  tableList: ApiMessageDelivery[]
  params?: any
  loading: boolean
  sort?: string
  messageType?: string
  startDate?: string;
  endDate?: string;
}

interface Params {
  page?: number;
  limit?: number;
  sort?: string;
  messageType?: string;
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
    label: t('messageDelivery.dateTimerange'),
    colProps: {
      span: 24
    },
    value: [],
    componentProps: {
      type: 'datetimerange',
      shortcuts: [
        {
          text: t('messageDelivery.today'),
          value: () => {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - start.getHours() * 3600 * 1000 - start.getMinutes() * 60 * 1000 - start.getSeconds() * 1000)
            return [start, end]
          }
        },
        {
          text: t('messageDelivery.yesterday'),
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
          text: t('messageDelivery.aWeekAgo'),
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
    field: 'messageTypes',
    label: t('messageDelivery.messageType'),
    component: 'CheckboxButton',
    value: [],
    colProps: {
      span: 24
    },
    componentProps: {
      options: [
        {
          label: 'web_push',
          value: 'web_push'
        },
        {
          label: 'html5_notify',
          value: 'html5_notify'
        },
        {
          label: 'email',
          value: 'email'
        },
        {
          label: 'sms',
          value: 'sms'
        },
        {
          label: 'telegram',
          value: 'telegram'
        },
        {
          label: 'slack',
          value: 'slack'
        },

      ]
    }
  },
])
const columns: TableColumn[] = [
  {
    field: 'id',
    label: t('messageDelivery.id'),
    sortable: true,
    width: "60px"
  },
  {
    field: 'messageType',
    label: t('messageDelivery.messageType'),
    sortable: true,
    width: "130px",
    formatter: (row: ApiMessageDelivery) => {
      return h(
        'span',
        row.message?.type
      )
    }
  },
  {
    field: 'attributes',
    label: t('messageDelivery.attributes'),
    sortable: true,
    formatter: (row: ApiMessageDelivery) => {
      return h(
        'span',
        Object.keys(row.message?.attributes).length || t('messageDelivery.nothing')
      )
    }
  },
  {
    field: 'address',
    label: t('messageDelivery.address'),
    sortable: true,
    formatter: (row: ApiMessageDelivery) => {
      return h(
        'span',
        row.address || t('messageDelivery.nothing')
      )
    }
  },
  {
    field: 'status',
    label: t('messageDelivery.status'),
    sortable: true,
    width: "90px"
  },
  {
    field: 'createdAt',
    label: t('main.createdAt'),
    type: 'time',
    sortable: true,
    width: "170px",
    formatter: (row: ApiMessageDelivery) => {
      return h(
        'span',
        parseTime(row.createdAt)
      )
    }
  },
  {
    field: 'createdAt',
    label: t('main.updatedAt'),
    type: 'time',
    sortable: true,
    width: "170px",
    formatter: (row: ApiMessageDelivery) => {
      return h(
        'span',
        parseTime(row.updatedAt)
      )
    }
  },
]
const paginationObj = ref<Pagination>({
  currentPage: 1,
  pageSize: 100,
  total: 0,
  pageSizes: [50, 100, 150, 250],
})
const currentID = ref('')

const getList = async () => {
  tableObject.loading = true
  let params: Params = {
    page: paginationObj.value.currentPage,
    limit: paginationObj.value.pageSize,
    sort: tableObject.sort,
    messageType: tableObject.messageType,
    startDate: tableObject.startDate,
    endDate: tableObject.endDate,
  }

  const res = await api.v1.messageDeliveryServiceGetMessageDeliveryList(params)
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
}

const onLogs = (log: ApiMessageDelivery) => {
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
  const {messageTypes} = formData
  if (messageTypes && messageTypes.length > 0) {
    tableObject.messageType = messageTypes.join(',')
  } else {
    tableObject.messageType = undefined
  }
  getList()
}

const tableRowClassName = (data) => {
  const {row, rowIndex} = data
  let style = ''
  switch (row?.status) {
    case 'in_progress':
      style = 'in_progress';
      break;
    case 'succeed':
      style = 'succeed';
      break;
    case 'error':
      style = 'error';
      break;
  }
  return style
}

const dialogVisible = ref(false)
const dialogTitle = ref('')
const dialogSelectedRow = ref<Nullable<ApiMessage>>(null)
const selectRow = (row: ApiMessageDelivery) => {
  if (!row) {
    return
  }
  const {message} = row
  dialogSelectedRow.value = message || null;
  dialogVisible.value = true
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
      @current-change="selectRow"
      style="width: 100%"
      class="messageDeliveryTable"
      :selection="false"
      :showUpPagination="20"
    />
  </ContentWrap>

  <Dialog v-model="dialogVisible" :title="t('messageDelivery.dialogTitle')">
    <AttributesViewer :message="dialogSelectedRow"/>
    <template #footer>
      <ElButton @click="dialogVisible = false">{{ t('main.closeDialog') }}</ElButton>
    </template>
  </Dialog>

</template>

<style lang="less">

.messageDeliveryTable {
  .el-table__row {
    td.el-table__cell {
      padding: 0;
      border-bottom: none !important;
    }

    &.error {
      --el-table-tr-bg-color: var(--el-color-error-light-5);
    }

    &.succeed {
      background-color: inherit;
    }

    &.in_progress {
      --el-table-tr-bg-color: var(--el-color-success-light-5);
    }
  }

  .el-table__row {
    cursor: pointer;
  }
}


</style>
