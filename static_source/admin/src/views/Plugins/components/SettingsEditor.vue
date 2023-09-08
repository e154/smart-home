<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {Table} from '@/components/Table'
import {computed, PropType, reactive, ref, unref, watch} from 'vue'
import {useAppStore} from "@/store/modules/app";
import {Pagination, TableColumn} from '@/types/table'
import api from "@/api/api";
import {ElButton, ElTableColumn, ElSwitch, ElImageViewer, ElTag, ElInput, ElSelect, ElOption} from 'element-plus'
import {ApiAttribute} from "@/api/stub";
import {useForm} from "@/hooks/web/useForm";
import {useRouter} from "vue-router";
import {Plugin} from "@/views/Plugins/components/Types";
import {parseTime} from "@/utils";
import {PATH_URL} from "@/config/axios/service";
import {useEmitt} from "@/hooks/web/useEmitt";

const { emitter } = useEmitt()
const {push, currentRoute} = useRouter()
const remember = ref(false)
const {register, elFormRef, methods} = useForm()
const appStore = useAppStore()
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
    width: "150px"
  },
  {
    field: 'type',
    label: t('plugins.attrType'),
    width: "100px"
  },
  {
    field: 'value',
    label: t('plugins.attrValue'),
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
  return  PATH_URL + imageUrl;
}

const boolOptions = [
  {
    value: true,
    label: 'TRUE',
  },
  {
    value: false,
    label: 'FALSE',
  },
]

const save = async () => {
  const settings = unref(tableObject.tableList)
  emitter.emit('settingsUpdated', settings)
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
      <div v-if="row.type === 'STRING'">
        <el-input type="string" v-model="row.string"/>
      </div>
      <div v-if="row.type === 'IMAGE'">
        <el-input type="string" v-model="row.imageUrl"/>
      </div>
      <div v-if="row.type === 'INT'">
        <el-input type="number" v-model="row.int"/>
      </div>
      <div v-if="row.type === 'FLOAT'">
        <el-input type="number" v-model="row.float"/>
      </div>
      <div v-if="row.type === 'ARRAY'">
        <el-input type="string" v-model="row.array"/>
      </div>
      <el-select v-model="row.bool" placeholder="please select value" v-if="row.type === 'BOOL'">
        <el-option
            v-for="item in boolOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value"
        />
      </el-select>

      <div v-if="row.type === 'TIME'">
        <el-input type="string" v-model="row.time"/>
      </div>

      <div v-if="row.type === 'MAP'">
        <el-input type="string" v-model="row.map"/>
      </div>
    </template>
  </Table>

  <div style="text-align: right" class="mt-20px">
    <ElButton type="primary" @click="save()">
      {{ t('main.save') }}
    </ElButton>
  </div>

</template>

<style lang="less">

</style>
