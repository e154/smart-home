<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {Table} from '@/components/Table'
import {computed, h, onMounted, onUnmounted, reactive, ref, watch} from 'vue'
import {useAppStore} from "@/store/modules/app";
import {Pagination, TableColumn} from '@/types/table'
import api from "@/api/api";
import {ElButton, ElTag} from 'element-plus'
import {ApiAction} from "@/api/stub";
import {useForm} from "@/hooks/web/useForm";
import {useRouter} from "vue-router";
import ContentWrap from "@/components/ContentWrap/src/ContentWrap.vue";
import {parseTime} from "@/utils";
import { Dialog } from '@/components/Dialog'
import Viewer from "@/components/JsonViewer/JsonViewer.vue";
import {useEmitt} from "@/hooks/web/useEmitt";
import {EventActionCompleted, EventStateChange, EventTriggerActivated} from "@/api/stream_types";
import {UUID} from "uuid-generator-ts";
import stream from "@/api/stream";
import {useCache} from "@/hooks/web/useCache";

const {push, currentRoute} = useRouter()
const remember = ref(false)
const {register, elFormRef, methods} = useForm()
const appStore = useAppStore()
const {t} = useI18n()
const { wsCache } = useCache()
const isMobile = computed(() => appStore.getMobile)

const dialogSource = ref("")
const dialogVisible = ref(false)
const importedTask = ref("")

interface TableObject {
  tableList: ApiAction[]
  params?: any
  loading: boolean
  sort?: string
}

interface Params {
  page?: number;
  limit?: number;
  sort?: string;
}

const cachePref = 'actions'
const tableObject = reactive<TableObject>(
    {
      tableList: [],
      loading: false,
      sort: wsCache.get(cachePref+'Sort') || '-createdAt'
    }
);

const currentID = ref('')

const onStateChanged = (event: EventStateChange) => {
  getList()
}

const onEventActionActivated = (event: EventActionCompleted) => {
  for (const i in tableObject.tableList) {
    if (tableObject.tableList[i].id == event.id) {
      tableObject.tableList[i].completed = true;
      setTimeout(() => {
        tableObject.tableList[i].completed = false
      }, 500)
      return
    }
  }
}

onMounted(() => {
  const uuid = new UUID()
  currentID.value = uuid.getDashFreeUUID()

  setTimeout(() => {
    stream.subscribe('event_action_completed', currentID.value, onEventActionActivated);
  }, 1000)
})

onUnmounted(() => {
  stream.unsubscribe('event_action_completed', currentID.value);
})

const columns: TableColumn[] = [
  {
    field: 'id',
    label: t('automation.actions.id'),
    sortable: true,
    width: "60px"
  },
  {
    field: 'name',
    label: t('automation.actions.name'),
    sortable: true,
  },
  {
    field: 'createdAt',
    label: t('main.createdAt'),
    type: 'time',
    sortable: true,
    width: "150px",
    formatter: (row: ApiAction) => {
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
    width: "150px",
    formatter: (row: ApiAction) => {
      return h(
          'span',
          parseTime(row.updatedAt)
      )
    }
  },
]
const paginationObj = ref<Pagination>({
  currentPage: wsCache.get(cachePref+'CurrentPage') || 1,
  pageSize: wsCache.get(cachePref+'PageSize') || 50,
  total: 0,
  pageSizes: [50, 100, 150, 250],
})

const getList = async () => {
  tableObject.loading = true

  wsCache.set(cachePref+'CurrentPage', paginationObj.value.currentPage)
  wsCache.set(cachePref+'PageSize', paginationObj.value.pageSize)
  wsCache.set(cachePref+'Sort', tableObject.sort)

  let params: Params = {
    page: paginationObj.value.currentPage,
    limit: paginationObj.value.pageSize,
    sort: tableObject.sort,
  }

  const res = await api.v1.actionServiceGetActionList(params)
      .catch(() => {
      })
      .finally(() => {
        tableObject.loading = false
      })
  if (res) {
    const {items, meta} = res.data;
    tableObject.tableList = items;
    paginationObj.value.currentPage = meta.page;
    paginationObj.value.total = meta.total;
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
  push('/automation/actions/new')
}


const selectRow = (row) => {
  if (!row) {
    return
  }
  const {id} = row
  push(`/automation/actions/edit/${id}`)
}

const showImportDialog = () => {
  dialogVisible.value = true
}

useEmitt({
  name: 'updateSource',
  callback: (val: string) => {
    if (importedTask.value == val) {
      return
    }
    importedTask.value = val
  }
})

const tableRowClassName = (data) => {
  const { row, rowIndex } = data
  let style = ''
  if (row.completed) {
    style = 'completed'
  }
  return style
}

</script>

<template>
  <ContentWrap>
    <ElButton class="flex mb-20px items-left" type="primary" @click="addNew()" plain>
      <Icon icon="ep:plus" class="mr-5px"/>
      {{ t('automation.actions.addNew') }}
    </ElButton>

    <Table
        :selection="false"
        v-model:pageSize="paginationObj.pageSize"
        v-model:currentPage="paginationObj.currentPage"
        :columns="columns"
        :data="tableObject.tableList"
        :loading="tableObject.loading"
        :pagination="paginationObj"
        @sort-change="sortChange"
        :row-class-name="tableRowClassName"
        style="width: 100%"
        :showUpPagination="20"
    >
      <template #name="{ row }">
        <span @click.prevent.stop="selectRow(row)" style="cursor: pointer">
          {{ row.name }}
        </span>
      </template>

    </Table>

  </ContentWrap>

</template>

<style lang="less">

.light {
  .el-table__row {
    &.completed {
      --el-table-tr-bg-color: var(--el-color-primary-light-7);
      -webkit-transition: background-color 200ms linear;
      -ms-transition: background-color 200ms linear;
      transition: background-color 200ms linear;

    }
  }
}

.dark {
  .el-table__row {
    &.completed {
      --el-table-tr-bg-color: var(--el-color-primary-dark-2);
      -webkit-transition: background-color 200ms linear;
      -ms-transition: background-color 200ms linear;
      transition: background-color 200ms linear;
    }
  }
}

</style>
