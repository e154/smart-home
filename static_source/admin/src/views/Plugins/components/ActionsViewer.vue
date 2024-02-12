<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {Table} from '@/components/Table'
import {h, PropType, reactive, watch} from 'vue'
import {TableColumn} from '@/types/table'
import {ElImageViewer, ElTableColumn, ElTag} from 'element-plus'
import {ApiPluginOptionsResultEntityAction} from "@/api/stub";
import {GetApiAttrValue} from "@/components/Attributes";
import {GetFullUrl} from "@/utils/serverId";

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


</script>

<template>
  <Table
      :selection="false"
      :columns="columns"
      :data="tableObject.tableList"
      :loading="tableObject.loading"
      style="width: 100%"
  >
    <ElTableColumn prop="system" label="System">
      <template #default="scope">
        <ElTag
            :type="scope.row.tag === 'Home' ? '' : 'success'"
            disable-transitions
        >{{ scope.row.tag }}
        </ElTag>
      </template>
    </ElTableColumn>

    <template #value="{ row }">
      <div v-if="row.type === 'IMAGE'">
        <ElImageViewer style="width: 100px; height: 100px" v-bind="GetFullUrl(row.imageUrl)"/>
      </div>
      <div v-else>
        <span>{{ GetApiAttrValue(row) }}</span>
      </div>
    </template>
  </Table>

</template>

<style lang="less">

</style>
