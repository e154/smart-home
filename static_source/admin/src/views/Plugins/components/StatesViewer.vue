<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {Table} from '@/components/Table'
import {computed, h, PropType, reactive, ref, watch} from 'vue'
import {useAppStore} from "@/store/modules/app";
import {Pagination, TableColumn} from '@/types/table'
import {ElButton, ElTableColumn, ElSwitch, ElImageViewer, ElTag, ElImage} from 'element-plus'
import {ApiPluginOptionsResultEntityState} from "@/api/stub";
import {useForm} from "@/hooks/web/useForm";
import {useRouter} from "vue-router";
import {useCache} from "@/hooks/web/useCache";
import {prepareUrl} from "@/utils/serverId";

const {push, currentRoute} = useRouter()
const remember = ref(false)
const {register, elFormRef, methods} = useForm()
const appStore = useAppStore()
const {t} = useI18n()
const {wsCache} = useCache()

interface TableObject {
  tableList: ApiPluginOptionsResultEntityState[]
  loading: boolean
}

const props = defineProps({
  states: {
    type: Array as PropType<ApiPluginOptionsResultEntityState[]>,
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
    sortable: true,
  },
  {
    field: 'imageUrl',
    label: t('plugins.actionImage'),
    sortable: true,
  },
  {
    field: 'icon',
    label: t('plugins.actionIcon'),
    formatter: (row: ApiPluginOptionsResultEntityState) => {
      return h(
          'span',
          row.icon ? '+' : '-'
      )
    }
  },
  {
    field: 'description',
    label: t('plugins.actionDescription')
  },
]


watch(
    () => props.states,
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

</script>

<template>
  <Table
      :selection="false"
      :columns="columns"
      :data="tableObject.tableList"
      :loading="tableObject.loading"
      style="width: 100%"
  >

    <template #imageUrl="{ row }">
      <ElImage v-if="row.imageUrl" style="width: 100px; height: 100px" :src="getUrl(row.imageUrl)"/>
      <span v-else>-</span>
    </template>
  </Table>

</template>

<style lang="less">

</style>
