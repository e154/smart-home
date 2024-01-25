<script setup lang="ts">
import {useI18n} from '@/hooks/web/useI18n'
import {Table} from '@/components/Table'
import {h, onMounted, onUnmounted, reactive, ref} from 'vue'
import {useAppStore} from "@/store/modules/app";
import {TableColumn} from '@/types/table'
import api from "@/api/api";
import {ElButton, ElCol, ElMessage, ElPopconfirm, ElRow, ElUpload, UploadProps} from 'element-plus'
import {useForm} from "@/hooks/web/useForm";
import {useRouter} from "vue-router";
import ContentWrap from "@/components/ContentWrap/src/ContentWrap.vue";
import {ApiBackup} from "@/api/stub";
import {parseTime} from "@/utils";
import {formatBytes} from "@/views/Dashboard/filters";
import {useCache} from "@/hooks/web/useCache";
import {useIcon} from "@/hooks/web/useIcon";
import {UUID} from "uuid-generator-ts";
import stream from "@/api/stream";

const Sun = useIcon({icon: 'emojione-monotone:sun', color: '#fde047'})

const {push, currentRoute} = useRouter()
const remember = ref(false)
const {register, elFormRef, methods} = useForm()
const appStore = useAppStore()
const {wsCache} = useCache()
const {t} = useI18n()

interface TableObject {
  tableList: ApiBackup[]
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

const currentID = ref('')

const columns: TableColumn[] = [
  {
    field: 'name',
    label: t('backup.name'),
    sortable: true,
  },
  {
    field: 'size',
    label: t('backup.size'),
    width: "100px",
    formatter: (row: ApiBackup) => {
      return h(
          'span',
          formatBytes(row.size.toString(), 2)
      )
    }
  },
  {
    field: 'operations',
    label: t('backup.operations'),
    width: "170px",
  },
  {
    field: 'modTime',
    label: t('main.createdAt'),
    type: 'time',
    sortable: true,
    width: "170px",
    formatter: (row: ApiBackup) => {
      return h(
          'span',
          parseTime(row.modTime)
      )
    }
  },
]

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

const addNew = async () => {
  const res = await api.v1.backupServiceNewBackup({})
      .catch(() => {
      })
      .finally(() => {
      })
  if (res.status == 200) {
    ElMessage({
      title: t('Success'),
      message: t('message.createdSuccessfully'),
      type: 'success',
      duration: 2000
    });
  }
}

const restore = async (backup: ApiBackup) => {
  const res = await api.v1.backupServiceRestoreBackup(backup.name)
      .catch(() => {
      })
      .finally(() => {
      })

  if (res.status == 200) {
    ElMessage({
      title: t('Success'),
      message: t('message.callSuccessful'),
      type: 'success',
      duration: 2000
    });
  }
}

const remove = async (backup: ApiBackup) => {
  const res = await api.v1.backupServiceDeleteBackup(backup.name)
      .catch(() => {
      })
      .finally(() => {
      })

  if (res.status == 200) {
    ElMessage({
      title: t('Success'),
      message: t('message.callSuccessful'),
      type: 'success',
      duration: 2000
    });
  }
}

const getUploadURL = () => {
  let uri = import.meta.env.VITE_API_BASEPATH as string || window.location.origin;
  const accessToken = wsCache.get("accessToken")
  uri += '/v1/backup/upload?access_token=' + accessToken;
  const serverId = wsCache.get('serverId')
  if (serverId) {
    uri += '&server_id=' + serverId;
  }
  return uri;
}


const getDownloadURL = (file: ApiBackup) => {
  let uri = import.meta.env.VITE_API_BASEPATH as string || window.location.origin;
  const accessToken = wsCache.get("accessToken")
  uri += '/snapshots/' + file.name + '?access_token=' + accessToken;
  const serverId = wsCache.get('serverId')
  if (serverId) {
    uri += '&server_id=' + serverId;
  }
  return uri;
}

const onSuccess: UploadProps['onSuccess'] = (file: ApiBackup, uploadFile) => {
  ElMessage({
    message: t('message.uploadSuccessfully'),
    type: 'success',
    duration: 2000
  })
}

const onError: UploadProps['onError'] = (error, uploadFile, uploadFiles) => {
  const body = JSON.parse(error.message)
  const {message, code} = body.error;
  ElMessage({
    message: message,
    type: 'error',
    duration: 0
  })
}

const forceFileDownload = (file: ApiBackup) => {
  const link = document.createElement('a')
  link.href = getDownloadURL(file)
  link.setAttribute('download', file.name)
  document.body.appendChild(link)
  link.click()
}

const onEventhandler = () => {
  getList()
}

const onEventStartedRestoreHandler = () => {
  ElMessage({
    message: t('message.startedRestoreProcess'),
    type: 'warning',
    duration: 0
  })
}

onMounted(() => {
  const uuid = new UUID()
  currentID.value = uuid.getDashFreeUUID()

  setTimeout(() => {
    stream.subscribe('event_created_backup', currentID.value, onEventhandler);
    stream.subscribe('event_removed_backup', currentID.value, onEventhandler);
    stream.subscribe('event_uploaded_backup', currentID.value, onEventhandler);
    stream.subscribe('event_started_restore', currentID.value, onEventStartedRestoreHandler);
  }, 1000)
})

onUnmounted(() => {
  stream.unsubscribe('event_created_backup', currentID.value);
  stream.unsubscribe('event_removed_backup', currentID.value);
  stream.unsubscribe('event_uploaded_backup', currentID.value);
  stream.unsubscribe('event_started_restore', currentID.value);
})

getList()

</script>

<template>
  <ContentWrap>
    <ElRow class="file-manager-body mb-20px">
      <ElCol>
        <ElButton class="flex mb-20px items-left" type="primary" @click="addNew()" plain>
          <Icon icon="iconoir:database-restore" class="mr-5px"/>
          {{ t('backup.addNew') }}
        </ElButton>

<!--        <ElButton class="flex mb-20px items-left" type="primary" @click="addNew()" plain>-->
<!--          {{ t('backup.apply') }}-->
<!--        </ElButton>-->

<!--        <ElButton class="flex mb-20px items-left" type="primary" @click="addNew()" plain>-->
<!--          {{ t('backup.rollback') }}-->
<!--        </ElButton>-->


        <ElUpload
            class="upload-demo"
            :action="getUploadURL()"
            :multiple="true"
            :on-success="onSuccess"
            :on-error="onError"
            :auto-upload="true"
        >
          <ElButton type="primary" plain>
            <Icon icon="material-symbols:upload" class="mr-5px"/>
            {{ $t('backup.uploadDump') }}
          </ElButton>
        </ElUpload>
      </ElCol>
    </ElRow>


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
            width="auto"
            style="margin-left: 10px;"
            :title="$t('backup.restoreSnapshot')"
            @confirm="restore(row)"
        >
          <template #reference>
            <ElButton class="flex items-right" type="danger" link>
              <Icon icon="ic:baseline-restore" class="mr-5px"/>
            </ElButton>
          </template>
        </ElPopconfirm>

        <ElButton class="flex items-right" link @click="forceFileDownload(row)">
          <Icon icon="material-symbols:download" class="mr-5px" />
        </ElButton>

        <ElPopconfirm
            :confirm-button-text="$t('main.ok')"
            :cancel-button-text="$t('main.no')"
            width="auto"
            style="margin-left: 10px;"
            :title="$t('backup.removeSnapshot')"
            @confirm="remove(row)"
        >
          <template #reference>
            <ElButton class="flex items-right" link>
              <Icon icon="mdi:remove" class="mr-5px"/>
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
