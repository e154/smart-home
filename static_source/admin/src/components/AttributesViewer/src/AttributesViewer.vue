<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {Table} from '@/components/Table'
import {PropType, reactive, watch} from 'vue'
import {ElImage} from 'element-plus'
import {TableColumn} from '@/types/table'
import {ApiAttribute} from "@/api/stub";
import {Attribute} from "@/views/Entities/components/types";
import {parseTime} from "@/utils";
import {prepareUrl} from "@/utils/serverId";

const {t} = useI18n()

interface TableObject {
  tableList: Attribute[]
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
    field: 'name',
    label: t('attributes.name'),
    sortable: true
  },
  {
    field: 'type',
    label: t('attributes.type'),
    sortable: true
  },
  {
    field: 'value',
    label: t('attributes.value')
  },
]

const attributes = (): Attribute[] => {
  const attr: Attribute[] = [];
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

const getUrl = (imageUrl: string | undefined): string => {
  if (!imageUrl) {
    return '';
  }
  return prepareUrl(import.meta.env.VITE_API_BASEPATH as string + imageUrl);
}

const getValue = (attr: Attribute): any => {
  switch (attr.type) {
    case 'STRING':
      return attr.string;
    case 'POINT':
      return attr.point;
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
    case 'ENCRYPTED':
      return attr.encrypted;
  }
}

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
