<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {Table} from '@/components/Table'
import {PropType, reactive, ref, watch} from 'vue'
import {useAppStore} from "@/store/modules/app";
import {TableColumn} from '@/types/table'
import {ElImageViewer} from 'element-plus'
import {ApiMessage} from "@/api/stub";
import {useForm} from "@/hooks/web/useForm";
import {useRouter} from "vue-router";

const {push, currentRoute} = useRouter()
const remember = ref(false)
const {register, elFormRef, methods} = useForm()
const appStore = useAppStore()
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
    width: "150px",
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
