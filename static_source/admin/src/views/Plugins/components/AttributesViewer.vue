<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {Table} from '@/components/Table'
import {computed, PropType, reactive, ref, watch} from 'vue'
import {useAppStore} from "@/store/modules/app";
import {TableColumn} from '@/types/table'
import {ElButton, ElTableColumn, ElSwitch, ElImage, ElTag} from 'element-plus'
import {ApiAttribute} from "@/api/stub";
import {useForm} from "@/hooks/web/useForm";
import {useRouter} from "vue-router";
import {Plugin} from "@/views/Plugins/components/Types";
import {parseTime} from "@/utils";
import {useCache} from "@/hooks/web/useCache";
import {prepareUrl} from "@/utils/serverId";

const {push, currentRoute} = useRouter()
const remember = ref(false)
const {register, elFormRef, methods} = useForm()
const appStore = useAppStore()
const {t} = useI18n()
const {wsCache} = useCache()

interface TableObject {
  tableList: ApiAttribute[]
  loading: boolean
}

const props = defineProps({
  attrs: {
    type: Array as PropType<ApiAttribute[]>,
    default: () => []
  }
})


const tableObject = reactive<TableObject>(
    {
      tableList: [],
      loading: false,
    },
);

const columns: TableColumn[] = [
  {
    field: 'name',
    label: t('plugins.name'),
    sortable: true
  },
  {
    field: 'type',
    label: t('plugins.attrType'),
    sortable: true
  },
  {
    field: 'value',
    label: t('plugins.attrValue')
  },
]


watch(
    () => props.attrs,
    (currentRow) => {
      if (!currentRow) return
      tableObject.tableList = currentRow
    },
    {
      deep: true,
      immediate: true
    }
)

const getUrl = (imageUrl: string | undefined): string => {
  if (!imageUrl) {
    return '';
  }
  return  prepareUrl(import.meta.env.VITE_API_BASEPATH as string + imageUrl);
}

const getValue = (attr: ApiAttribute): any => {
  switch (attr.type) {
    case 'STRING':
      return attr.string;
    case 'INT':
      return attr.int;
    case 'FLOAT':
      return attr.float;
    case 'ARRAY':
      return attr.array;
    case 'BOOL':
      return attr.bool;
    case 'TIME':
      return parseTime(attr.time);
    case 'MAP':
      return attr.map;
    case 'IMAGE':
      return getUrl(attr.imageUrl);
  }
}

</script>

<template>
  <Table
      :selection="false"
      :columns="columns"
      :data="tableObject.tableList"
      :loading="tableObject.loading"
      style="width: 100%"
  >
    <template #value="{ row }">
      <div v-if="row.type === 'IMAGE'">
        <ElImage style="width: 100px; height: 100px" :src="getUrl(row.imageUrl)"/>
      </div>
      <div v-else>
        <span>{{ getValue(row) }}</span>
      </div>
    </template>
  </Table>

</template>

<style lang="less">

</style>
