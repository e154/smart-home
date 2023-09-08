<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {Table} from '@/components/Table'
import {computed, h, PropType, reactive, ref, watch} from 'vue'
import {useAppStore} from "@/store/modules/app";
import {Pagination, TableColumn} from '@/types/table'
import api from "@/api/api";
import {ElButton, ElTableColumn, ElSwitch, ElImageViewer, ElTag} from 'element-plus'
import {ApiPluginOptionsResultEntityAction} from "@/api/stub";
import {useForm} from "@/hooks/web/useForm";
import {useRouter} from "vue-router";
import {Plugin} from "@/views/Plugins/components/Types";
import {parseTime} from "@/utils";
import {PATH_URL} from "@/config/axios/service";

const {push, currentRoute} = useRouter()
const remember = ref(false)
const {register, elFormRef, methods} = useForm()
const appStore = useAppStore()
const {t} = useI18n()


interface TableObject {
  tableList: ApiPluginOptionsResultEntityAction[]
  loading: boolean
}

const props = defineProps({
  actions: {
    type: Array as PropType<ApiPluginOptionsResultEntityAction[]>,
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
    formatter: (row: ApiPluginOptionsResultEntityAction) => {
      return h(
          'span',
          row.imageUrl ? '+' : '-'
      )
    }
  },
  {
    field: 'icon',
    label: t('plugins.actionIcon'),
    formatter: (row: ApiPluginOptionsResultEntityAction) => {
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
    () => props.actions,
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
  return  PATH_URL + imageUrl;
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
    <ElTableColumn prop="system"  label="System">
      <template #default="scope">
        <ElTag
            :type="scope.row.tag === 'Home' ? '' : 'success'"
            disable-transitions
        >{{ scope.row.tag }}</ElTag>
      </template>
    </ElTableColumn>

    <template #value="{ row }">
      <div v-if="row.type === 'IMAGE'">
        <ElImageViewer style="width: 100px; height: 100px" v-bind="getUrl(row.imageUrl)" />
      </div>
      <div v-else>
        <span>{{ getValue(row) }}</span>
      </div>
    </template>
  </Table>

</template>

<style lang="less">

</style>
