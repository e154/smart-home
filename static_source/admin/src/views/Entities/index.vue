<script setup lang="tsx">
import {useI18n} from '@/hooks/web/useI18n'
import {Table} from '@/components/Table'
import {h, onMounted, onUnmounted, reactive, ref, watch} from 'vue'
import {Pagination, TableColumn} from '@/types/table'
import api from "@/api/api";
import {ElButton, ElCollapse, ElCollapseItem, ElMessage, ElTag} from 'element-plus'
import {ApiArea, ApiEntity, ApiEntityShort, ApiPlugin, ApiStatistics, ApiTag} from "@/api/stub";
import {useForm} from "@/hooks/web/useForm";
import {useRouter} from "vue-router";
import {parseTime} from "@/utils";
import {ContentWrap} from "@/components/ContentWrap";
import {Dialog} from '@/components/Dialog'
import {UUID} from "uuid-generator-ts";
import stream from "@/api/stream";
import {EventStateChange} from "@/api/types";
import {FormSchema} from "@/types/form";
import {Form} from '@/components/Form'
import {useCache} from "@/hooks/web/useCache";
import {JsonEditor} from "@/components/JsonEditor";
import Statistics from "@/components/Statistics/Statistics.vue";

const {push} = useRouter()
const {register, methods} = useForm()
const {t} = useI18n()
const {wsCache} = useCache()

interface TableObject {
  tableList: ApiEntityShort[]
  params?: any
  loading: boolean
  sort?: string
  query?: string
  plugin?: string
  area?: ApiArea
  tags?: string[]
}

interface Params {
  page?: number;
  limit?: number;
  sort?: string;
}

const cachePref = 'entities'
const tableObject = reactive<TableObject>(
    {
      tableList: [],
      loading: false,
      sort: wsCache.get(cachePref + 'Sort') || '-createdAt',
      query: wsCache.get(cachePref + 'Query'),
      plugin: wsCache.get(cachePref + 'Plugin')?.name,
      area: wsCache.get(cachePref + 'Area'),
      tags: wsCache.get(cachePref + 'Tags')
    },
);

const columns: TableColumn[] = [
  {
    field: 'id',
    label: t('entities.name'),
    width: "190px",
    sortable: true
  },
  {
    field: 'pluginName',
    label: t('entities.pluginName'),
    width: "140px",
    sortable: true
  },
  // {
  //   field: 'tags',
  //   label: t('main.tags'),
  //   width: "150px",
  //   sortable: false,
  //   type: 'expand'
  // },
  {
    field: 'areaId',
    label: t('entities.area'),
    width: "100px",
    sortable: true,
    formatter: (row: ApiEntityShort) => {
      return h(
          'span',
          row.area?.name
      )
    }
  },
  {
    field: 'description',
    sortable: true,
    label: t('entities.description')
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
    formatter: (row: ApiEntityShort) => {
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
    formatter: (row: ApiEntityShort) => {
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

const statistic = ref<Nullable<ApiStatistics>>(null)
const getStatistic = async () => {

  const res = await api.v1.entityServiceGetStatistic()
    .catch(() => {
    })
    .finally(() => {

    })
  if (res) {
    statistic.value = res.data
  }
}

const getList = async () => {
  getStatistic()

  tableObject.loading = true

  wsCache.set(cachePref + 'CurrentPage', paginationObj.value.currentPage)
  wsCache.set(cachePref + 'PageSize', paginationObj.value.pageSize)
  wsCache.set(cachePref + 'Sort', tableObject.sort)
  wsCache.set(cachePref + 'Query', tableObject.query)
  wsCache.set(cachePref + 'Tags', tableObject.tags)

  let params: Params = {
    page: paginationObj.value.currentPage,
    limit: paginationObj.value.pageSize,
    sort: tableObject.sort,
    query: tableObject.query || undefined,
    plugin: tableObject.plugin || undefined,
    area: tableObject?.area?.id || undefined,
    tags: tableObject?.tags || undefined,
  }

  const res = await api.v1.entityServiceGetEntityList(params)
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

const currentID = ref('')

const onStateChanged = (event: EventStateChange) => {
  getList()
}

onMounted(() => {
  const uuid = new UUID()
  currentID.value = uuid.getDashFreeUUID()

  setTimeout(() => {
    stream.subscribe('event_entity_loaded', currentID.value, onStateChanged);
    stream.subscribe('event_entity_unloaded', currentID.value, onStateChanged);
  }, 200)
})

onUnmounted(() => {
  stream.unsubscribe('event_entity_loaded', currentID.value);
  stream.unsubscribe('event_entity_unloaded', currentID.value);
})

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
  push('/entities/new')
}

const selectRow = (row) => {
  if (!row) {
    return
  }
  const {id} = row
  push(`/entities/edit/${id}`)
}

const restart = async (entity: ApiEntityShort) => {
  await api.v1.developerToolsServiceReloadEntity({id: entity.id});
  ElMessage({
    title: t('Success'),
    message: t('message.requestSentSuccessfully'),
    type: 'success',
    duration: 2000
  });
}

const enable = async (entity: ApiEntityShort) => {
  if (!entity?.id) return;
  await api.v1.entityServiceEnabledEntity(entity.id);
  ElMessage({
    title: t('Success'),
    message: t('message.requestSentSuccessfully'),
    type: 'success',
    duration: 2000
  });
}

const disable = async (entity: ApiEntityShort) => {
  if (!entity?.id) return;
  await api.v1.entityServiceDisabledEntity(entity.id);
  ElMessage({
    title: t('Success'),
    message: t('message.requestSentSuccessfully'),
    type: 'success',
    duration: 2000
  });
}

const dialogVisible = ref(false)
const importedEntity = ref(null)
const showImportDialog = () => {
  dialogVisible.value = true
}

const importHandler = (val: any) => {
  if (importedEntity.value == val) {
    return
  }
  importedEntity.value = val
}

const importEntity = async () => {
  let val: ApiEntity;
  try {
    if (importedEntity.value?.json) {
      val = importedEntity.value.json as ApiEntity;
    } else if (importedEntity.value.text) {
      val = JSON.parse(importedEntity.value.text) as ApiEntity;
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
  const entity: ApiEntity = {
    id: val.id,
    pluginName: val.pluginName,
    description: val.description,
    area: val.area,
    icon: val.icon,
    image: val.image,
    autoLoad: val.autoLoad,
    parent: val.parent || undefined,
    actions: val.actions,
    states: val.states,
    attributes: val.attributes,
    settings: val.settings,
    scripts: val.scripts,
    tags: val.tags,
    metrics: val.metrics
  }
  const res = await api.v1.entityServiceImportEntity(entity)
  if (res) {
    ElMessage({
      title: t('Success'),
      message: t('message.importedSuccessful'),
      type: 'success',
      duration: 2000
    })
    getList()
  }
}

// search form
const schema = reactive<FormSchema[]>([
  {
    field: 'name',
    label: t('entities.name'),
    component: 'Input',
    colProps: {
      span: 12
    },
    componentProps: {
      placeholder: t('entities.name'),
      onChange: (val: string) => {
        tableObject.query = val || undefined
        getList()
      }
    }
  },
  {
    field: 'plugin',
    label: t('entities.pluginName'),
    component: 'Plugin',
    value: null,
    colProps: {
      span: 12
    },
    hidden: false,
    componentProps: {
      placeholder: t('entities.pluginName'),
      onChange: (val: ApiPlugin) => {
        tableObject.plugin = val?.name || undefined
        wsCache.set(cachePref + 'Plugin', val)
        getList()
      }
    },
  },
  {
    field: 'area',
    label: t('entities.area'),
    value: null,
    component: 'Area',
    colProps: {
      span: 12
    },
    componentProps: {
      placeholder: t('entities.area'),
      onChange: (val: ApiArea) => {
        wsCache.set(cachePref + 'Area', val)
        tableObject.area = val || undefined
        getList()
      }
    },
  },
  {
    field: 'tags',
    label: t('main.tags'),
    component: 'Tags',
    colProps: {
      span: 12
    },
    value: [],
    hidden: false,
    componentProps: {
      placeholder: t('main.tags'),
      onChange: (val: ApiTag) => {
        wsCache.set(cachePref + 'Tags', val)
        tableObject.tags = val || undefined
        getList()
      }
    }
  },
])

const filterList = () => {
  let list = ''
  if (tableObject?.query) {
    list += 'name(' + tableObject.query + ') '
  }
  if (tableObject?.plugin) {
    list += 'plugin(' + tableObject.plugin + ') '
  }
  if (tableObject?.area) {
    list += 'area(' + tableObject.area.name + ') '
  }
  if (tableObject?.tags && tableObject?.tags.length) {
    list += 'tags(' + tableObject.tags + ') '
  }
  if (list != '') {
    list = ': ' + list
  }
  return list
}

const {setValues, setSchema} = methods
if (wsCache.get(cachePref + 'Query')) {
  setValues({
    name: wsCache.get(cachePref + 'Query')
  })
}
if (wsCache.get(cachePref + 'Plugin')) {
  setValues({
    plugin: wsCache.get(cachePref + 'Plugin')
  })
}
if (wsCache.get(cachePref + 'Area')) {
  setValues({
    area: wsCache.get(cachePref + 'Area')
  })
}
if (wsCache.get(cachePref + 'Tags')) {
  setValues({
    tags: wsCache.get(cachePref + 'Tags')
  })
}

</script>

<template>
  <Statistics v-model="statistic" :cols="6" />

  <ContentWrap>
    <ElButton class="flex mb-20px items-left" type="primary" @click="addNew()" plain>
      <Icon icon="ep:plus" class="mr-5px"/>
      {{ t('entities.addNew') }}
    </ElButton>
    <ElButton class="flex mb-20px items-left" type="primary" @click="showImportDialog()" plain>
      {{ t('main.import') }}
    </ElButton>
    <ElCollapse class="mb-20px">
      <ElCollapseItem :title="$t('main.filter') + filterList()">
        <Form
            class="filter-form"
            :schema="schema"
            label-position="top"
            label-width="auto"
            hide-required-asterisk
            @register="register"
        />
      </ElCollapseItem>
    </ElCollapse>
    <Table
        class="entity-table"
        :expand="true"
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
    >

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

      <template #id="{row}">
        {{ row.id.split('.')[1] }}
      </template>

      <template #pluginName="{row}">
        <ElTag>
          {{ row.pluginName }}
        </ElTag>
      </template>

      <template #expand="{row}">
        <div class="tag-list" v-if="row.tags">
          <ElTag v-for="tag in row.tags" type="info" :key="tag" round effect="light" size="small">
            {{ tag }}
          </ElTag>
        </div>
      </template>

    </Table>
  </ContentWrap>

  <!-- import dialog -->
  <Dialog v-model="dialogVisible" :title="t('entities.dialogImportTitle')" :maxHeight="400" width="80%" custom-class>
    <JsonEditor @change="importHandler"/>
    <template #footer>
      <ElButton type="primary" @click="importEntity()" plain>{{ t('main.import') }}</ElButton>
      <ElButton @click="dialogVisible = false">{{ t('main.closeDialog') }}</ElButton>
    </template>
  </Dialog>
  <!-- /import dialog -->

</template>

<style lang="less">

//:deep(.filter-form .el-col) {
//  padding-left: 0!important;
//  padding-right: 0!important;
//}

.entity-table {
  .tag-list {
    .el-tag {
      margin: 0 5px;
    }
  }

  :deep(.el-table__row) {
    cursor: pointer;
  }

  tr.el-table__row [class*="el-table__cell"] {
  //background-color: green; border-top: var(--el-table-border); border-bottom: none !important;
    border-top: var(--el-table-border);
  }

  .el-table__expanded-cell {
    &.el-table__cell [class*="tag-list"] {
    //background-color: red!important; border-bottom: none !important;
    }

    &.el-table__cell:not(:has(.tag-list)) {
      display: none !important;
    //background-color: blue!important;
    }
  }

  .el-table td.el-table__cell,
  .el-table th.el-table__cell.is-leaf {
    border-bottom: none !important;
  }

  .el-table--enable-row-hover .el-table__body tr.el-table__row:hover,
  .el-table--enable-row-hover .el-table__body tr.el-table__row:hover + tr {
    & > td.el-table__cell {
      background-color: var(--el-table-row-hover-bg-color);
    }
  }
}


</style>
