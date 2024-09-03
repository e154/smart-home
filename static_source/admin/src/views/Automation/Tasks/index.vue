<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {Table} from '@/components/Table'
import {h, onMounted, onUnmounted, reactive, ref, watch} from 'vue'
import {Pagination, TableColumn} from '@/types/table'
import api from "@/api/api";
import {ElButton, ElMessage} from 'element-plus'
import {ApiStatistics, ApiTask} from "@/api/stub";
import {useRouter} from "vue-router";
import {ContentWrap} from "@/components/ContentWrap";
import {parseTime} from "@/utils";
import {EventStateChange, EventTaskCompleted} from "@/api/types";
import {UUID} from "uuid-generator-ts";
import stream from "@/api/stream";
import {useCache} from "@/hooks/web/useCache";
import Statistics from "@/components/Statistics/Statistics.vue";

const {push} = useRouter()
const {t} = useI18n()
const {wsCache} = useCache()

interface TableObject {
  tableList: ApiTask[]
  params?: any
  loading: boolean
  sort?: string
}

interface Params {
  page?: number;
  limit?: number;
  sort?: string;
}

const cachePref = 'tasks'
const tableObject = reactive<TableObject>(
    {
      tableList: [],
      loading: false,
      sort: wsCache.get(cachePref + 'Sort') || '-createdAt'
    }
);

const currentID = ref('')

const onStateChanged = (event: EventStateChange) => {
  getList()
  getStatistic()
}

const onEventTaskActivated = (event: EventTaskCompleted) => {
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
    stream.subscribe('event_task_loaded', currentID.value, onStateChanged);
    stream.subscribe('event_task_unloaded', currentID.value, onStateChanged);
    stream.subscribe('event_task_completed', currentID.value, onEventTaskActivated);
  }, 1000)
})

onUnmounted(() => {
  stream.unsubscribe('event_task_loaded', currentID.value);
  stream.unsubscribe('event_task_unloaded', currentID.value);
  stream.unsubscribe('event_task_completed', currentID.value);
})

const columns: TableColumn[] = [
  {
    field: 'id',
    label: t('automation.id'),
    sortable: true,
    width: "60px"
  },
  {
    field: 'name',
    label: t('automation.name'),
    sortable: true,
    width: "170px"
  },
  {
    field: 'areaId',
    label: t('automation.area'),
    width: "100px",
    sortable: true,
    formatter: (row: ApiTask) => {
      return h(
          'span',
          row.area?.name
      )
    }
  },
  {
    field: 'actions',
    label: t('automation.tasks.actions'),
    width: "100px",
    formatter: (row: ApiTask) => {
      return h(
          'span',
          row?.actions?.length || t('automation.nothing')
      )
    }
  },
  {
    field: 'triggers',
    label: t('automation.tasks.triggers'),
    width: "100px",
    formatter: (row: ApiTask) => {
      return h(
          'span',
          row?.triggers?.length || t('automation.nothing')
      )
    }
  },
  {
    field: 'description',
    label: t('automation.description'),
    sortable: true,
    formatter: (row: ApiTask) => {
      return h(
          'span',
          row?.description || t('automation.nothing')
      )
    }
  },
  {
    field: 'status',
    label: t('entities.status'),
    width: "70px",
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

  const res = await api.v1.automationServiceGetTaskList(params)
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
  push('/automation/tasks/new')
}


const selectRow = (row) => {
  if (!row) {
    return
  }
  const {id} = row
  push(`/automation/tasks/edit/${id}`)
}

const enable = async (task: ApiTask) => {
  if (!task?.id) return;
  await api.v1.automationServiceEnableTask(task.id);
  ElMessage({
    title: t('Success'),
    message: t('message.requestSentSuccessfully'),
    type: 'success',
    duration: 2000
  });
}

const disable = async (task: ApiTask) => {
  if (!task?.id) return;
  await api.v1.automationServiceDisableTask(task.id);
  ElMessage({
    title: t('Success'),
    message: t('message.requestSentSuccessfully'),
    type: 'success',
    duration: 2000
  });
}

const tableRowClassName = (data) => {
  const {row, rowIndex} = data
  let style = ''
  if (row.completed) {
    style = 'completed'
  }
  return style
}

const statistic = ref<Nullable<ApiStatistics>>(null)
const getStatistic = async () => {

  const res = await api.v1.automationServiceGetStatistic()
    .catch(() => {
    })
    .finally(() => {

    })
  if (res) {
    statistic.value = res.data
  }
}

getStatistic()

</script>

<template>

  <Statistics v-model="statistic" :cols="6"/>

  <ContentWrap>
    <ElButton class="flex mb-20px items-left" type="primary" @click="addNew()" plain>
      <Icon icon="ep:plus" class="mr-5px"/>
      {{ t('automation.addNew') }}
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

      <template #status="{ row }">
        <div class="w-[100%] text-center">
          <ElButton :link="true" @click.prevent.stop="enable(row)" v-if="!row?.isLoaded">
            <Icon icon="noto:red-circle" class="mr-5px"/>
          </ElButton>
          <ElButton :link="true" @click.prevent.stop="disable(row)" v-if="row?.isLoaded">
            <Icon icon="noto:green-circle" class="mr-5px"/>
          </ElButton>
        </div>
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
