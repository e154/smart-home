<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {Table} from '@/components/Table'
import {h, PropType, reactive, watch} from 'vue'
import {TableColumn} from '@/types/table'
import {ElImage} from 'element-plus'
import {ApiPluginOptionsResultEntityState} from "@/api/stub";
import {GetFullUrl} from "@/utils/serverId";

const {t} = useI18n()

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
      <ElImage v-if="row.imageUrl" style="width: 100px; height: 100px" :src="GetFullUrl(row.imageUrl)"/>
      <span v-else>-</span>
    </template>
  </Table>

</template>

<style lang="less">

</style>
