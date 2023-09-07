<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {Table} from '@/components/Table'
import {computed, h, PropType, reactive, ref, unref, watch} from 'vue'
import {useAppStore} from "@/store/modules/app";
import {Pagination, TableColumn} from '@/types/table'
import api from "@/api/api";
import {ElButton, ElSwitch} from 'element-plus'
import {ApiArea, ApiZigbee2Mqtt, ApiZigbee2MqttDevice} from "@/api/stub";
import {useForm} from "@/hooks/web/useForm";
import {useRouter} from "vue-router";

const {push, currentRoute} = useRouter()
const remember = ref(false)
const {register, elFormRef, methods} = useForm()
const appStore = useAppStore()
const {t} = useI18n()

const currentBridge = ref<Nullable<ApiZigbee2Mqtt>>({})

const props = defineProps({
  modelValue: {
    type: Object as PropType<Nullable<ApiZigbee2Mqtt>>,
    default: () => undefined
  }
})

interface TableObject {
  tableList: ApiZigbee2MqttDevice[]
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
      sort: '-id'
    },
);

const columns: TableColumn[] = [
  {
    field: 'id',
    label: t('zigbee2mqtt.id'),
    width: "160px",
    sortable: true
  },
  {
    field: 'model',
    label: t('zigbee2mqtt.model'),
    sortable: true
  },
  {
    field: 'status',
    label: t('zigbee2mqtt.status'),
    width: "100px",
    sortable: true
  },
  {
    field: 'description',
    label: t('zigbee2mqtt.description'),
    sortable: true
  },
]
const paginationObj = ref<Pagination>({
  currentPage: 1,
  pageSize: 50,
  total: 0,
})
const currentID = ref('')

const getList = async () => {
  tableObject.loading = true

  let params: Params = {
    page: paginationObj.value.currentPage,
    limit: paginationObj.value.pageSize,
    sort: tableObject.sort,
  }

  if (!currentBridge.value?.id) {
    return
  }

  const res = await api.v1.zigbee2MqttServiceDeviceList(currentBridge.value?.id, params)
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
    () => props.modelValue,
    (val?: ApiZigbee2Mqtt) => {
      if (val === unref(currentBridge)) return
      currentBridge.value = val
      getList()
    },
    {
      immediate: true
    }
)

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

const selectRow = (row) => {

}

</script>

<template>

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

</template>

<style lang="less">

.el-table__row {
  cursor: pointer;
}
</style>
