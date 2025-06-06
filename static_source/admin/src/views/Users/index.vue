<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {Table} from '@/components/Table'
import {h, reactive, ref, watch} from 'vue'
import {useAppStore} from "@/store/modules/app";
import {Pagination, TableColumn} from '@/types/table'
import api from "@/api/api";
import {ElButton} from 'element-plus'
import {ApiUserShot, ApiLog} from "@/api/stub";
import {useForm} from "@/hooks/web/useForm";
import {useRouter} from "vue-router";
import {parseTime} from "@/utils";
import {ContentWrap} from "@/components/ContentWrap";
import {useCache} from "@/hooks/web/useCache";

const {push, currentRoute} = useRouter()
const remember = ref(false)
const {register, elFormRef, methods} = useForm()
const appStore = useAppStore()
const {t} = useI18n()
const { wsCache } = useCache()

interface TableObject {
  tableList: ApiUserShot[]
  params?: any
  loading: boolean
  sort?: string
}

interface Params {
  page?: number;
  limit?: number;
  sort?: string;
}

const cachePref = 'users'
const tableObject = reactive<TableObject>(
    {
      tableList: [],
      loading: false,
      sort: wsCache.get(cachePref+'Sort') || '-id'
    },
);

const columns: TableColumn[] = [
  {
    field: 'id',
    label: t('users.id'),
    width: "60px",
    sortable: true
  },
  {
    field: 'nickname',
    label: t('users.nickname'),
    sortable: true
  },
  {
    field: 'role',
    label: t('users.role'),
    sortable: true,
    formatter: (row: ApiUserShot) => {
      return h(
          'span',
          row.roleName
      )
    }
  },
  {
    field: 'email',
    label: t('users.email'),
    sortable: true
  },
  {
    field: 'status',
    label: t('users.status'),
    sortable: true
  },
  {
    field: 'createdAt',
    label: t('main.createdAt'),
    type: 'time',
    sortable: true,
    width: "170px",
    formatter: (row: ApiUserShot) => {
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
    formatter: (row: ApiUserShot) => {
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
})
const currentID = ref('')

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

  const res = await api.v1.userServiceGetUserList(params)
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
  push('/etc/users/new')
}

const selectRow = (row) => {
  if (!row) {
    return
  }
  const {id} = row
  push(`/etc/users/edit/${id}`)
}

</script>

<template>

  <ContentWrap>
  <ElButton class="flex mb-20px items-left" type="primary" @click="addNew()" plain>
    <Icon icon="ep:plus" class="mr-5px"/>
    {{ t('users.addNew') }}
  </ElButton>
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
      @current-change="selectRow"
  />
  </ContentWrap>

</template>

<style lang="less">

.el-table__row {
  cursor: pointer;
}
</style>
