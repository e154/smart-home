<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {Table} from '@/components/Table'
import {computed, PropType, reactive, ref, watch} from 'vue'
import {TableColumn} from '@/types/table'
import {ElImage} from 'element-plus'
import {ApiAttribute} from "@/api/stub";
import {getUrl, getValue} from "@/views/Plugins/components/Types";

const {t} = useI18n()

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
