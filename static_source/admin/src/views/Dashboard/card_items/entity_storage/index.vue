<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {Table} from '@/components/Table'
import {Form} from '@/components/Form'
import {h, onMounted, onUnmounted, PropType, reactive, ref, unref, watch} from 'vue'
import {FormSchema} from "@/types/form";
import {useAppStore} from "@/store/modules/app";
import {Pagination, TableColumn} from '@/types/table'
import api from "@/api/api";
import {UUID} from 'uuid-generator-ts'
import {ElButton} from 'element-plus'
import {ApiEntity, ApiEntityStorage} from "@/api/stub";
import {useForm} from "@/hooks/web/useForm";
import {parseTime} from "@/utils";
import stream from "@/api/stream";
import { Dialog } from '@/components/Dialog'
import {EventStateChange} from "@/api/stream_types";
import AttributesViewer from "@/views/Entities/components/AttributesViewer.vue";
import {CardItem} from "@/views/Dashboard/core";
import {debounce} from "lodash-es";

const remember = ref(false)
const {register, elFormRef, methods} = useForm()
const appStore = useAppStore()
const {t} = useI18n()
const dialogVisible = ref(false)
const dialogSource = ref({})

interface TableObject {
  tableList: ApiEntityStorage[]
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

const props = defineProps({
  item: {
    type: Object as PropType<Nullable<CardItem>>,
    default: () => null
  },
})

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
    label: t('entityStorage.dateTimerange'),
    colProps: {
      span: 24
    },
    value: [],
    componentProps: {
      type: 'datetimerange',
      shortcuts: [
        {
          text: t('entityStorage.today'),
          value: () => {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - start.getHours() * 3600 * 1000 - start.getMinutes() * 60 * 1000 - start.getSeconds() * 1000)
            return [start, end]
          }
        },
        {
          text: t('entityStorage.yesterday'),
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
          text: t('entityStorage.aWeekAgo'),
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
])

const columns: TableColumn[] = [
  {
    field: 'createdAt',
    label: t('main.createdAt'),
    type: 'time',
    sortable: true,
    width: "150px",
    formatter: (row: ApiEntityStorage) => {
      return h(
          'span',
          parseTime(row.createdAt)
      )
    }
  },
  {
    field: 'state',
    label: t('entityStorage.state'),
    sortable: true,
    width: "200px",
  },
  {
    field: 'attributes',
    label: t('entityStorage.attributes'),
  },
  {
    field: 'entityId',
    label: t('entityStorage.entityId'),
    width: "200px",
  },
]
const paginationObj = ref<Pagination>({
  currentPage: 1,
  pageSize: 20,
  total: 0,
  pageSizes: [20, 50, 100, 150, 250],
})

const el = ref(null)
const currentID = ref('')
onMounted(() => {
  props.item.setTarget(el.value)

  const uuid = new UUID()
  currentID.value = uuid.getDashFreeUUID()

  setTimeout(() => {
    stream.subscribe('state_changed', currentID.value, onStateChanged)
  }, 1000)
})

onUnmounted(() => {
  stream.unsubscribe('state_changed', currentID.value)
})

const getList = debounce( async () => {

  tableObject.loading = true

  let params: Params = {
    page: paginationObj.value.currentPage,
    limit: paginationObj.value.pageSize,
    sort: tableObject.sort,
    startDate: tableObject.startDate,
    endDate: tableObject.endDate,
  }

  const res = await api.v1.entityStorageServiceGetEntityStorageList(params)
      .catch(() => {
      })
      .finally(() => {
        tableObject.loading = false
      })
  if (res) {
    const {items, meta} = res.data;
    tableObject.tableList = items;
    paginationObj.value.pageSize = meta.pagination.limit;
    paginationObj.value.currentPage = meta.pagination.page;
    paginationObj.value.total = meta.pagination.total;
  } else {
    tableObject.tableList = [];
  }
}, 100)

const onStateChanged = (event: EventStateChange) => {
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
  getList()
}

const selectRow = (row: ApiEntityStorage) => {
  if (!row) return;
  dialogSource.value = row?.attributes
  dialogVisible.value = true
}

getList()

</script>

<template>

  <div ref="el" class="w-[100%] h-[100%]" style="overflow:hidden;overflow-y: scroll">

    <Dialog v-model="dialogVisible" :maxHeight="400" width="80%">
      <div style="padding: 10px">
        <AttributesViewer v-model="dialogSource"/>
      </div>
      <template #footer>
        <ElButton @click="dialogVisible = false">{{ t('main.closeDialog') }}</ElButton>
      </template>
    </Dialog>

    <Table
        v-model:pageSize="paginationObj.pageSize"
        v-model:currentPage="paginationObj.currentPage"
        :columns="columns"
        :data="tableObject.tableList"
        :pagination="paginationObj"
        @sort-change="sortChange"
        style="width: 100%"
        class="storageTable"
        :selection="false"
        :showUpPagination="20"
        @current-change="selectRow"
    >
      <template #attributes="{row}">
        <span>{{ Object.keys(row.attributes).length || $t('entityStorage.nothing') }}</span>
      </template>
    </Table>

  </div>
</template>

<style lang="less">

</style>
