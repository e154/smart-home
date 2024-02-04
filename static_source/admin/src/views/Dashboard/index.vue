<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {Table} from '@/components/Table'
import {h, reactive, ref, watch} from 'vue'
import {Pagination, TableColumn} from '@/types/table'
import api from "@/api/api";
import {ElButton, ElMessage} from 'element-plus'
import {ApiDashboard} from "@/api/stub";
import {useRouter} from "vue-router";
import {parseTime} from "@/utils";
import ContentWrap from "@/components/ContentWrap/src/ContentWrap.vue";
import {Dialog} from '@/components/Dialog'
import {Core} from "@/views/Dashboard/core";
import {useCache} from "@/hooks/web/useCache";
import JsonEditor from "@/components/JsonEditor/JsonEditor.vue";
import {prepareUrl} from "@/utils/serverId";
import {Infotip} from "@/components/Infotip";

const {push} = useRouter()
const counter = ref(0);
const {t} = useI18n()
const {wsCache} = useCache()

interface TableObject {
  tableList: ApiDashboard[]
  params?: any
  loading: boolean
  sort?: string
}

interface Params {
  page?: number;
  limit?: number;
  sort?: string;
}

const cachePref = 'dashboard'
const tableObject = reactive<TableObject>(
    {
      tableList: [],
      loading: false,
      sort: wsCache.get(cachePref + 'Sort') || '-createdAt'
    },
);

const columns: TableColumn[] = [
  {
    field: 'id',
    label: t('dashboard.id'),
    sortable: true,
    width: "60px"
  },
  {
    field: 'name',
    label: t('dashboard.name'),
    sortable: true,
    width: "170px"
  },
  {
    field: 'areaId',
    label: t('dashboard.area'),
    width: "100px",
    sortable: true,
    formatter: (row: ApiDashboard) => {
      return h(
          'span',
          row.area?.name
      )
    }
  },
  {
    field: 'description',
    sortable: true,
    label: t('dashboard.description')
  },
  {
    field: 'link',
    label: t('dashboard.landing'),
    width: "120px",
  },
  {
    field: 'operations',
    label: t('dashboard.operations'),
    width: "120px",
  },
  {
    field: 'createdAt',
    label: t('main.createdAt'),
    type: 'time',
    sortable: true,
    width: "170px",
    formatter: (row: ApiDashboard) => {
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
    formatter: (row: ApiDashboard) => {
      return h(
          'span',
          parseTime(row.updatedAt)
      )
    }
  }
]
const paginationObj = ref<Pagination>({
  currentPage: wsCache.get(cachePref + 'CurrentPage') || 1,
  pageSize: wsCache.get(cachePref + 'PageSize') || 50,
  total: 0,
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

  const res = await api.v1.dashboardServiceGetDashboardList(params)
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
    counter.value = meta.pagination.total + 1
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
  Core.createNew('new' + counter.value)
      .then((dashboard: ApiDashboard) => {
        ElMessage({
          title: t('Success'),
          message: t('message.createdSuccessfully'),
          type: 'success',
          duration: 2000
        });
        push({path: `/dashboards/edit/${dashboard.id}`});
      })
      .catch((e) => {
        counter.value++
      })
}

const editDashboard = (row: ApiDashboard) => {
  push(`/dashboards/edit/${row.id}`)
}

const showDashboard = (row: ApiDashboard) => {
  push(`/dashboards/view/${row.id}`)
}

//todo: experimental function
const openLanding = (item: ApiDashboard): string => {
  const uri = window.location.origin || import.meta.env.VITE_API_BASEPATH as string;
  const accessToken = wsCache.get("accessToken")
  const url = prepareUrl(uri + '/#/landing/' + item.id + '?access_token=' + accessToken);
  console.log('url', url)
  window.open(url, '_blank', 'noreferrer');
}

const dialogVisible = ref(false)
const importValue = ref(null)

const importHandler = (val: any) => {
  if (importValue.value == val) {
    return
  }
  importValue.value = val
}

const importDashboard = async () => {
  let dashboard: ApiDashboard
  try {
    if (importValue.value?.json) {
      dashboard = importValue.value.json as ApiDashboard;
    } else if (importValue.value.text) {
      dashboard = JSON.parse(importValue.value.text) as ApiDashboard;
    }
  } catch {
    ElMessage({
      title: t('Error'),
      message: t('message.corruptedJsonFormat'),
      type: 'error',
      duration: 2000
    });
    return
  }
  const data = await Core._import(dashboard);
  if (data) {
    await getList();
    ElMessage({
      title: t('Success'),
      message: t('message.importedSuccessful'),
      type: 'success',
      duration: 2000
    });
  }
}

</script>

<template>
  <ContentWrap>
    <ElButton class="flex mb-20px items-left" type="primary" @click="addNew()" plain>
      <Icon icon="ep:plus" class="mr-5px"/>
      {{ t('dashboard.addNewDashboard') }}
    </ElButton>
    <ElButton class="flex mb-20px items-left" type="primary" @click="dialogVisible = true" plain>
      {{ t('main.import') }}
    </ElButton>

    <Infotip
        type="warning"
        :show-index="false"
        title="WARNING"
        :schema="[
      {
        label: 'The functionality is experimental and may change in the future.',
      },
      {
        label: 'Added a landing page link that provides direct access to the dashboard. The direct link allows you to embed the dashboard into your solution.',
      },
      {
        label: ' &nbsp;',
      },
      {
        label: 'Be careful, the link contains your access token. Generate a link on behalf of a non-privileged user.',
      },
    ]"
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
        @current-change="showDashboard"
    >

      <template #status="{ row }">

        <div class="w-[100%] text-center">
          <Icon icon="noto:green-circle" class="mr-5px" v-if="row?.isLoaded"/>
          <Icon icon="noto:red-circle" class="mr-5px" v-if="!row?.isLoaded"/>
        </div>

      </template>

      <template #link="{ row }">
        <div class="w-[100%] text-center">
          <ElButton :link="true" @click.prevent.stop="openLanding(row)">
            {{ $t('main.open') }} &nbsp; <Icon icon="pajamas:external-link" />
          </ElButton>
        </div>
      </template>

      <template #operations="{ row }">
        <div class="w-[100%] text-center">
          <ElButton :link="true" @click.prevent.stop="editDashboard(row)">
            {{ $t('main.edit') }}
          </ElButton>
        </div>
      </template>
    </Table>
  </ContentWrap>

  <!-- import dialog -->
  <Dialog v-model="dialogVisible" :title="t('dashboard.dialogImportTitle')" :maxHeight="400" width="80%" custom-class>
    <JsonEditor @change="importHandler"/>
    <template #footer>
      <ElButton type="primary" @click="importDashboard()" plain>{{ t('main.import') }}</ElButton>
      <ElButton @click="dialogVisible = false">{{ t('main.closeDialog') }}</ElButton>
    </template>
  </Dialog>
  <!-- /import dialog -->

</template>

<style lang="less">

.el-table__row {
  cursor: pointer;
}

</style>
