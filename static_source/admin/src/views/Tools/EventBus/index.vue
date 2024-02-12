<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {Table} from '@/components/Table'
import {onMounted, onUnmounted, reactive, ref, watch} from 'vue'
import {Pagination, TableColumn} from '@/types/table'
import api from "@/api/api";
import {ApiBusStateItem} from "@/api/stub";
import {ContentWrap} from "@/components/ContentWrap";
import {EventTriggerCompleted} from "@/api/types";
import {UUID} from "uuid-generator-ts";

const {t} = useI18n()

interface TableObject {
  tableList: ApiBusStateItem[]
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
    }
);

const currentID = ref('')

const onEventTriggerActivated = (event: EventTriggerCompleted) => {

}

onMounted(() => {
  const uuid = new UUID()
  currentID.value = uuid.getDashFreeUUID()

  // setTimeout(() => {
  //   stream.subscribe('event_trigger_loaded', currentID.value, onStateChanged);
  //   stream.subscribe('event_trigger_unloaded', currentID.value, onStateChanged);
  //   stream.subscribe('event_trigger_completed', currentID.value, onEventTriggerActivated);
  // }, 1000)
})

onUnmounted(() => {
  // stream.unsubscribe('event_trigger_loaded', currentID.value);
  // stream.unsubscribe('event_trigger_unloaded', currentID.value);
  // stream.unsubscribe('event_trigger_completed', currentID.value);
})

const columns: TableColumn[] = [
  {
    field: 'topic',
    label: t('tools.eventBus.topic'),
  },
  {
    field: 'min',
    label: t('tools.eventBus.min'),
    width: "80px"
  },
  {
    field: 'avg',
    label: t('tools.eventBus.avg'),
    width: "80px"
  },
  {
    field: 'max',
    label: t('tools.eventBus.max'),
    width: "80px"
  },
  {
    field: 'rps',
    label: t('tools.eventBus.rps'),
    width: "80px"
  },
  {
    field: 'subscribers',
    label: t('tools.eventBus.subscribers'),
    width: "120px"
  },
]
const paginationObj = ref<Pagination>({
  currentPage: 1,
  pageSize: 50,
  total: 0,
  pageSizes: [50, 100, 150, 250],
})

const getList = async () => {
  tableObject.loading = true

  let params: Params = {
    page: paginationObj.value.currentPage,
    limit: paginationObj.value.pageSize,
    sort: tableObject.sort,
  }

  const res = await api.v1.developerToolsServiceGetEventBusStateList(params)
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

const myInterval = ref()
onMounted(() => {
  myInterval.value = setInterval(() => {
    getList()
  }, 2000)
})

onUnmounted(() => {
  clearInterval(myInterval.value);
})

const tableRowClassName = (data) => {
  const {row, rowIndex} = data
  let style = ''
  if (row.completed) {
    style = 'completed'
  }
  return style
}

</script>

<template>
  <ContentWrap>

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
    />


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
