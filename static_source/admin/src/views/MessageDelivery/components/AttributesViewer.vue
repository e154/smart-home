<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {Table} from '@/components/Table'
import {PropType, reactive, watch} from 'vue'
import {TableColumn} from '@/types/table'
import {ApiMessage} from "@/api/stub";

const {t} = useI18n()


export interface MessageItem {
  name: string;
  value: string;
}

interface TableObject {
  tableList: MessageItem[]
  loading: boolean
}

const props = defineProps({
  message: {
    type: Object as PropType<Nullable<ApiMessage>>,
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
    label: t('plugins.name'),
    width: "170px",
    sortable: true
  },
  {
    field: 'value',
    label: t('plugins.attrValue')
  },
]

watch(
  () => props.message,
  (message) => {
    const items: MessageItem[] = [];
    for (const key in message?.attributes) {
      items.push({name: key, value: message?.attributes[key]});
    }
    tableObject.tableList = items
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
  />

</template>

<style lang="less">

</style>
