<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {Table} from '@/components/Table'
import {computed, h, onMounted, onUnmounted, reactive, ref, watch} from 'vue'
import {Pagination, TableColumn} from '@/types/table'
import api from "@/api/api";
import {ElButton, ElTag} from 'element-plus'
import {ApiCondition} from "@/api/stub";
import {useRouter} from "vue-router";
import ContentWrap from "@/components/ContentWrap/src/ContentWrap.vue";
import {parseTime} from "@/utils";
import {EventStateChange} from "@/api/stream_types";
import {UUID} from "uuid-generator-ts";
import {useCache} from "@/hooks/web/useCache";

const {push} = useRouter()
const {t} = useI18n()
const { wsCache } = useCache()

interface TableObject {
  tableList: ApiCondition[]
  params?: any
  loading: boolean
  sort?: string
}

interface Params {
  page?: number;
  limit?: number;
  sort?: string;
}

const cachePref = 'conditions'
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

onMounted(() => {
  const uuid = new UUID()
  currentID.value = uuid.getDashFreeUUID()

})

onUnmounted(() => {

})

const columns: TableColumn[] = [
  {
    field: 'id',
    label: t('automation.conditions.id'),
    sortable: true,
    width: "60px"
  },
  {
    field: 'name',
    label: t('automation.conditions.name'),
    sortable: true,
    width: "170px"
  },
  {
    field: 'areaId',
    label: t('automation.area'),
    width: "100px",
    sortable: true,
    formatter: (row: ApiCondition) => {
      return h(
          'span',
          row.area?.name
      )
    }
  },
  {
    field: 'description',
    label: t('automation.description'),
    sortable: true,
    formatter: (row: ApiCondition) => {
      return h(
          'span',
          row?.description || t('automation.nothing')
      )
    }
  },
  {
    field: 'createdAt',
    label: t('main.createdAt'),
    type: 'time',
    sortable: true,
    width: "170px",
    formatter: (row: ApiCondition) => {
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
    formatter: (row: ApiCondition) => {
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

  const res = await api.v1.conditionServiceGetConditionList(params)
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
  push('/automation/conditions/new')
}


const selectRow = (row) => {
  if (!row) {
    return
  }
  const {id} = row
  push(`/automation/conditions/edit/${id}`)
}

</script>

<template>
  <ContentWrap>
    <ElButton class="flex mb-20px items-left" type="primary" @click="addNew()" plain>
      <Icon icon="ep:plus" class="mr-5px"/>
      {{ t('automation.conditions.addNew') }}
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
        style="width: 100%"
        :showUpPagination="20"
    >
      <template #name="{ row }">
        <span @click.prevent.stop="selectRow(row)" style="cursor: pointer">
          {{ row.name }}
        </span>
      </template>

      <template #pluginName="{row}">
        <ElTag>
          {{ row.pluginName }}
        </ElTag>
      </template>
    </Table>

  </ContentWrap>

</template>

<style lang="less">

</style>
