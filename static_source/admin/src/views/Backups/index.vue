<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {Table} from '@/components/Table'
import {computed, reactive, ref, h} from 'vue'
import {useAppStore} from "@/store/modules/app";
import {Pagination, TableColumn} from '@/types/table'
import api from "@/api/api";
import {ElButton, ElMessage, ElPopconfirm} from 'element-plus'
import {useForm} from "@/hooks/web/useForm";
import {useRouter} from "vue-router";
import ContentWrap from "@/components/ContentWrap/src/ContentWrap.vue";

const {push, currentRoute} = useRouter()
const remember = ref(false)
const {register, elFormRef, methods} = useForm()
const appStore = useAppStore()
const {t} = useI18n()
const isMobile = computed(() => appStore.getMobile)

interface TableObject {
  tableList: string[]
  loading: boolean
}

interface Params {
  page?: number;
  limit?: number;
  sort?: string;
}

const tableObject = reactive<TableObject>(
    {
      tableList: [],
      loading: false,
    }
);

const columns: TableColumn[] = [
  {
    field: 'name',
    label: t('backup.name'),
    sortable: true,
    formatter: (row: string) => {
      return h(
          'span',
          row
      )
    }

  },
  {
    field: 'operations',
    label: t('backup.operations'),
    width: "150px",
  },
]

const currentID = ref('')

const getList = async () => {
  tableObject.loading = true
  const res = await api.v1.backupServiceGetBackupList()
      .catch(() => {
      })
      .finally(() => {
        tableObject.loading = false
      })
  if (res) {
    const {items, meta} = res.data;
    tableObject.tableList = items;
  } else {
    tableObject.tableList = [];
  }
}

getList()

const addNew = async () => {
  api.v1.backupServiceNewBackup({})
      .catch(() => {
      })
      .finally(() => {
        getList();
      })
  ElMessage({
    title: t('Success'),
    message: t('message.createdSuccessfully'),
    type: 'success',
    duration: 2000
  });
}

const restore = async (name: string) => {
  api.v1.backupServiceRestoreBackup({name: name})
      .catch(() => {
      })
      .finally(() => {
      })

  setTimeout(async () => {
    ElMessage({
      title: t('Success'),
      message: t('message.callSuccessful'),
      type: 'success',
      duration: 2000
    });
  }, 2000)
}

</script>

<template>
  <ContentWrap>
    <ElButton class="flex mb-20px items-left" type="primary" @click="addNew()" plain>
      <Icon icon="ep:plus" class="mr-5px"/>
      {{ t('backup.addNew') }}
    </ElButton>
    <Table
        :selection="false"
        :columns="columns"
        :data="tableObject.tableList"
        :loading="tableObject.loading"
        style="width: 100%"
    >
      <template #operations="{ row }">
        <ElPopconfirm
            :confirm-button-text="$t('main.ok')"
            :cancel-button-text="$t('main.no')"
            width="250"
            style="margin-left: 10px;"
            :title="$t('main.are_you_sure_to_do_want_this?')"
            @confirm="restore(row)"
        >
          <template #reference>
            <ElButton class="flex items-right" type="danger" plain>
              {{ t('backup.restore') }}
            </ElButton>
          </template>
        </ElPopconfirm>


      </template>
    </Table>
  </ContentWrap>

</template>

<style lang="less">

.el-table__row {
  cursor: pointer;
}
</style>
