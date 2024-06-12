<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {Table} from '@/components/Table'
import {PropType, reactive, watch} from 'vue'
import {ElImage} from 'element-plus'
import {TableColumn} from '@/types/table'
import {ApiAttribute} from "@/api/stub";
import {GetFullUrl} from "@/utils/serverId";
import {GetApiAttributeValue} from "@/components/Attributes";

const {t} = useI18n()

interface TableObject {
  tableList: ApiAttribute[]
  loading: boolean
}

const props = defineProps({
  modelValue: {
    type: Object as PropType<Record<string, ApiAttribute>>,
    default: () => null
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
    field: 'type',
    label: t('attributes.type'),
    sortable: true,
    width: '100px'
  },
  {
    field: 'name',
    label: t('attributes.name'),
    sortable: true,
  },
  {
    field: 'value',
    label: t('attributes.value')
  },
]

const attributes = (): ApiAttribute[] => {
  const attr: ApiAttribute[] = [];
  if (props.modelValue) {
    for (const key in props.modelValue) {
      attr.push(props.modelValue[key]);
    }
  }
  return attr;
}

watch(
  () => props.modelValue,
  (message: Record<string, ApiAttribute>) => {
    tableObject.tableList = attributes()
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
  >

    <template #value="{ row }">
      <div v-if="row.type === 'IMAGE'">
        <ElImage style="width: 100px; height: 100px" :src="GetFullUrl(row.imageUrl)"/>
      </div>
      <div v-else-if="row.type === 'ICON'">
        <Icon
          :icon="row.icon"
          :size="15"/>
      </div>
      <div v-else>
        <span>{{ GetApiAttributeValue(row) }}</span>
      </div>
    </template>

  </Table>

</template>

<style lang="less">

</style>
