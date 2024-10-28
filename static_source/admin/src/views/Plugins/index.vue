<script lang="ts" setup>
import {useI18n} from '@/hooks/web/useI18n'
import {Table} from '@/components/Table'
import {onMounted, onUnmounted, reactive, ref, watch} from 'vue'
import {Pagination, TableColumn} from '@/types/table'
import api from "@/api/api";
import {ElButton, ElMessage, ElRow, ElUpload, UploadProps, ElTag} from 'element-plus'
import {ApiBackup, ApiPlugin} from "@/api/stub";
import {useRouter} from "vue-router";
import {ContentWrap} from "@/components/ContentWrap";
import {EventStateChange} from "@/api/types";
import {UUID} from "uuid-generator-ts";
import stream from "@/api/stream";
import {useCache} from "@/hooks/web/useCache";

const {push} = useRouter()
const {wsCache} = useCache()
const {t} = useI18n()

interface TableObject {
  tableList: ApiPlugin[]
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
    sort: '+name'
  },
);

const columns: TableColumn[] = [
  {
    field: 'name',
    label: t('plugins.name'),
    sortable: true
  },
  {
    field: 'external',
    label: t('plugins.external'),
    sortable: true,
    width: "100px"
  },
  {
    field: 'version',
    label: t('plugins.version'),
    width: "100px"
  },
  {
    field: 'status',
    label: t('entities.status'),
    width: "70px",
  },
]
const paginationObj = ref<Pagination>({
  currentPage: 1,
  pageSize: 50,
  total: 0,
})
const currentID = ref('')

const onStateChanged = (event: EventStateChange) => {
  getList()
}

onMounted(() => {
  const uuid = new UUID()
  currentID.value = uuid.getDashFreeUUID()

  setTimeout(() => {
    stream.subscribe('event_plugin_loaded', currentID.value, onStateChanged);
    stream.subscribe('event_plugin_unloaded', currentID.value, onStateChanged);
  }, 200)
})

onUnmounted(() => {
  stream.unsubscribe('event_plugin_loaded', currentID.value);
  stream.unsubscribe('event_plugin_unloaded', currentID.value);
})

const getList = async () => {
  tableObject.loading = true

  let params: Params = {
    page: paginationObj.value.currentPage,
    limit: paginationObj.value.pageSize,
    sort: tableObject.sort,
  }

  const res = await api.v1.pluginServiceGetPluginList(params)
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
  () => [paginationObj.value.currentPage, paginationObj.value.pageSize],
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

const selectRow = (row) => {
  if (!row) {
    return
  }
  const {name} = row
  push(`/etc/settings/plugins/edit/${name}`)
}

const enable = async (plugin: ApiPlugin) => {
  if (!plugin?.name) return;
  await api.v1.pluginServiceEnablePlugin(plugin.name);
  ElMessage({
    title: t('Success'),
    message: t('message.requestSentSuccessfully'),
    type: 'success',
    duration: 2000
  });
}

const disable = async (plugin: ApiPlugin) => {
  if (!plugin?.name) return;
  await api.v1.pluginServiceDisablePlugin(plugin.name);
  ElMessage({
    title: t('Success'),
    message: t('message.requestSentSuccessfully'),
    type: 'success',
    duration: 2000
  });
}

const getUploadURL = () => {
  let uri = import.meta.env.VITE_API_BASEPATH as string || window.location.origin;
  const accessToken = wsCache.get("accessToken")
  uri += '/v1/plugins/upload?access_token=' + accessToken;
  const serverId = wsCache.get('serverId')
  if (serverId) {
    uri += '&server_id=' + serverId;
  }
  return uri;
}

const onSuccess: UploadProps['onSuccess'] = (file: ApiBackup, uploadFile) => {
  getList()
  ElMessage({
    message: t('message.uploadSuccessfully'),
    type: 'success',
    duration: 2000
  })
}

const onError: UploadProps['onError'] = (error, uploadFile, uploadFiles) => {
  getList()
  const body = JSON.parse(error.message)
  const {message, code} = body.error;
  ElMessage({
    message: message,
    type: 'error',
    duration: 5000
  })
}

</script>

<template>
  <ContentWrap>

    <ElRow class=" mb-20px">
      <ElUpload
        :action="getUploadURL()"
        :auto-upload="true"
        :multiple="false"
        :on-error="onError"
        :on-success="onSuccess"
        class="upload-demo"
      >
        <ElButton plain type="primary">
          <Icon class="mr-5px" icon="material-symbols:upload"/>
          {{ $t('plugins.uploadPlugin') }}
        </ElButton>
      </ElUpload>
    </ElRow>

    <Table
      v-model:currentPage="paginationObj.currentPage"
      v-model:pageSize="paginationObj.pageSize"
      :columns="columns"
      :data="tableObject.tableList"
      :loading="tableObject.loading"
      :pagination="paginationObj"
      :selection="false"
      :showUpPagination="20"
      style="width: 100%"
      @sort-change="sortChange"
      @current-change="selectRow"
    >
      <template #external="{ row }">
        <div class="w-[100%]" v-if="row">
          <ElTag v-if="row.external">
            {{ t('plugins.external') }}
          </ElTag>
        </div>
      </template>

      <template #status="{ row }">
        <div class="w-[100%] text-center">
          <ElButton v-if="!row?.isLoaded" :link="true" @click.prevent.stop="enable(row)">
            <Icon class="mr-5px" icon="noto:red-circle"/>
          </ElButton>
          <ElButton v-if="row?.isLoaded" :link="true" @click.prevent.stop="disable(row)">
            <Icon class="mr-5px" icon="noto:green-circle"/>
          </ElButton>
        </div>
      </template>
    </Table>
  </ContentWrap>
</template>

<style lang="less">

.el-table__row {
  cursor: pointer;
}
</style>
