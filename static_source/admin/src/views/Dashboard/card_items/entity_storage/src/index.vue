<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {Table} from '@/components/Table'
import {Form} from '@/components/Form'
import {onMounted, onUnmounted, PropType, reactive, ref, unref, watch} from 'vue'
import {FormSchema} from "@/types/form";
import {Pagination, TableColumn} from '@/types/table'
import api from "@/api/api";
import {UUID} from 'uuid-generator-ts'
import {ElButton, ElCheckboxButton, ElCheckboxGroup} from 'element-plus'
import {ApiEntityStorage} from "@/api/stub";
import {useForm} from "@/hooks/web/useForm";
import stream from "@/api/stream";
import {Dialog} from '@/components/Dialog'
import {EventStateChange} from "@/api/types";
import {AttributesViewer} from "@/components/Attributes";
import {CardItem, RenderVar} from "@/views/Dashboard/core";
import {debounce} from "lodash-es";
import {Column} from "@/views/Dashboard/card_items/entity_storage";

const {register} = useForm()
const {t} = useI18n()
const dialogVisible = ref(false)
const dialogSource = ref({})
const selectedEntities = ref([])
const filterList = ref([])

interface TableObject {
  tableList: ApiEntityStorage[]
  params?: any
  loading: boolean
  sort?: string
  entityId?: string
  startDate?: string;
  endDate?: string;
}

interface Params {
  page?: number;
  limit?: number;
  sort?: string;
  entityId?: string[];
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

const columns = ref<TableColumn[]>([])

const paginationObj = ref<Pagination>({
  currentPage: 1,
  pageSize: 20,
  total: 0,
  pageSizes: [20, 50, 100, 150, 250],
})

const el = ref<ElRef>(null)
const currentID = ref('')
onMounted(() => {
  const uuid = new UUID()
  currentID.value = uuid.getDashFreeUUID()

  setTimeout(() => {
    stream.subscribe('state_changed', currentID.value, onStateChanged)
  }, 1000)
})

onUnmounted(() => {
  stream.unsubscribe('state_changed', currentID.value)
})

const rebuildTable = debounce(() => {
  tableObject.loading = true

  console.log('rebuild table')

  columns.value = [];
  for (const col of props.item?.payload?.entityStorage.columns) {
    columns.value.push({
      field: col.name,
      label: col.name,
      sortable: col.sortable || false,
      width: col.width || 'auto',
    })
  }

  getList()

}, 1000)

const getList = debounce(async () => {

  tableObject.loading = true

  let params: Params = {
    page: paginationObj.value.currentPage,
    limit: paginationObj.value.pageSize,
    sort: tableObject.sort,
    startDate: tableObject.startDate,
    endDate: tableObject.endDate,
  }

  if (props.item?.payload.entityStorage?.entityIds?.length) {

    if (selectedEntities.value.length == 0) {
      params.entityId = props.item?.payload.entityStorage?.entityIds || []
    } else {
      params.entityId = selectedEntities.value
    }
  }

  const res = await api.v1.entityStorageServiceGetEntityStorageList(params)
    .catch(() => {
    })
    .finally(() => {
      tableObject.loading = false
    })
  if (res) {
    const {items, filter, meta} = res.data;
    tableObject.tableList = items;
    filterList.value = filter
    paginationObj.value.pageSize = meta.pagination.limit;
    paginationObj.value.currentPage = meta.pagination.page;
    paginationObj.value.total = meta.pagination.total;
  } else {
    tableObject.tableList = [];
  }
}, 1000)

const onStateChanged = (event: EventStateChange) => {

  if (props.item?.payload.entityStorage?.entityIds?.length) {
    if (selectedEntities.value.length == 0) {
      if (props.item?.payload.entityStorage?.entityIds.includes(event.entity_id)) {
        getList()
      }
    } else {
      if (selectedEntities.value.includes(event.entity_id)) {
        getList()
      }
    }
  } else {
    getList()
  }
}

watch(
  [() => paginationObj.value.currentPage,
    () => selectedEntities.value,
    () => paginationObj.value.pageSize,
  ],
  () => {
    getList()
  }
)

watch(
  [() => props.item?.payload.entityStorage?.columns,
    () => props.item?.payload.entityStorage?.columns && props.item?.payload.entityStorage.columns.length,
  ],
  () => {
    rebuildTable()
  },
  {
    immediate: true,
    deep: true
  }
)


watch(
  () => props.item,
  () => {
    selectedEntities.value = props.item?.payload.entityStorage?.entityIds || []
  },
  {
    immediate: true
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
  if (props.item?.payload.entityStorage?.showPopup) {
    dialogSource.value = row?.attributes
    dialogVisible.value = true
  }
  // if (props.item?.payload.entityStorage?.eventName) {
  //todo: edit ....
  // eventBus.emit(props.item?.payload.entityStorage.eventName, '')
  // stream.send({
  //   id: UUID.createUUID(),
  //   query: 'event_get_state_by_id',
  //   body: btoa(JSON.stringify({entity_id: row.entityId, storage_id: row.id}))
  // });
  // }
}

const renderCell = (row: any, col: Column): string => {
  let token = col.attribute || ''
  token = token.replace(/\.+$/, "");
  if (col.filter && token) {
    token += '|' + col.filter
  }
  const val = RenderVar(token, unref(row))
  return val || ''
}

getList()

</script>

<template>

  <div ref="el" class="w-[100%] h-[100%]" style="overflow:hidden;overflow-y: scroll" :class="[{'hidden': item.hidden}]">

    <Dialog v-model="dialogVisible" :maxHeight="400" width="80%">
      <div style="padding: 10px">
        <AttributesViewer v-model="dialogSource"/>
      </div>
      <template #footer>
        <ElButton @click="dialogVisible = false">{{ t('main.closeDialog') }}</ElButton>
      </template>
    </Dialog>

    <Form
      v-if="item.payload.entityStorage?.filter"
      :schema="schema"
      label-position="top"
      label-width="auto"
      hide-required-asterisk
      @change="onFormChange"
      @register="register"
    />

    <div class="mb-20px" v-if="item.payload.entityStorage?.entityIds?.length && item.payload.entityStorage?.filter">
      <ElCheckboxGroup v-model="selectedEntities">
        <ElCheckboxButton
          v-for="entity in item.payload.entityStorage.entityIds"
          :key="entity"
          :label="entity">
          {{ entity }}
        </ElCheckboxButton>
      </ElCheckboxGroup>
    </div>

    <Table
      v-if="item?.payload?.entityStorage?.columns && item.payload.entityStorage.columns.length"
      v-model:pageSize="paginationObj.pageSize"
      v-model:currentPage="paginationObj.currentPage"
      :columns="columns"
      :data="tableObject.tableList"
      :pagination="paginationObj"
      :loading="tableObject.loading"
      @sort-change="sortChange"
      style="width: 100%"
      class="storageTable"
      :selection="false"
      :showUpPagination="20"
      @current-change="selectRow"
    >
      <template v-for="(column, idx) in item.payload.entityStorage.columns" :key="idx" #[column.name]="{row}">
        <span>{{ renderCell(row, column) }}</span>
      </template>
    </Table>

  </div>
</template>

<style lang="less">

</style>
