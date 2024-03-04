<script setup lang="ts">
import {computed, nextTick, PropType, reactive, ref, unref, watch} from 'vue'
import {Form} from '@/components/Form'
import {ElButton, ElCol, ElDivider, ElMessage, ElPopconfirm, ElRow} from 'element-plus'
import {useI18n} from '@/hooks/web/useI18n'
import {useForm} from '@/hooks/web/useForm'
import {useValidator} from '@/hooks/web/useValidator'
import {FormSchema} from '@/types/form'
import {ApiArea, ApiDashboard} from "@/api/stub";
import {copyToClipboard} from "@/utils/clipboard";
import {JsonViewer} from "@/components/JsonViewer";
import {Core, useBus} from "@/views/Dashboard/core";
import {useRouter} from "vue-router";
import {Dialog} from '@/components/Dialog'

const {register, elFormRef, methods} = useForm()
const {required} = useValidator()
const {t} = useI18n()
const dialogSource = ref({})
const dialogVisible = ref(false)
const {setValues, setSchema} = methods
const {currentRoute, addRoute, push} = useRouter()
const {emit} = useBus()

interface DashboardForm {
  name?: string;
  description?: string;
  enabled?: boolean;
  area?: ApiArea;
  areaId?: number;
}

const props = defineProps({
  core: {
    type: Object as PropType<Nullable<Core>>,
    default: () => null
  },
})

const currentCore = computed(() => props.core as Core)

const rules = {
  name: [required()],
}

const schema = reactive<FormSchema[]>([
  {
    field: 'name',
    label: t('dashboard.name'),
    component: 'Input',
    colProps: {
      span: 24
    },
    componentProps: {
      placeholder: t('dashboard.name')
    }
  },
  {
    field: 'description',
    label: t('dashboard.description'),
    component: 'Input',
    colProps: {
      span: 24
    },
    componentProps: {
      placeholder: t('dashboard.description')
    }
  },
  // {
  //   field: 'enabled',
  //   label: t('dashboard.enabled'),
  //   component: 'Switch',
  //   value: false,
  //   colProps: {
  //     span: 24
  //   },
  // },
  {
    field: 'area',
    label: t('dashboard.area'),
    value: null,
    component: 'Area',
    colProps: {
      span: 24
    },
    componentProps: {
      placeholder: t('dashboard.area'),
    }
  },
])

watch(
    () => props.core?.current,
    (val?: ApiDashboard) => {
      if (!val) return
      setValues({
        name: val.name,
        description: val.description,
        enabled: val.enabled,
        area: val.area,
        areaId: val.areaId,
      })
    },
    {
      deep: false,
      immediate: true
    }
)

const prepareForExport = async () => {
  return ""
}

const copy = async () => {
  const body = await prepareForExport()
  copyToClipboard(JSON.stringify(body, null, 2))
}


const exportDashbord = () => {
  dialogSource.value = currentCore.value.serialize()
  dialogVisible.value = true
}

const updateBoard = async () => {
  const formRef = unref(elFormRef)
  await formRef?.validate(async (isValid) => {
    if (isValid) {
      const {getFormData} = methods
      const formData = await getFormData<DashboardForm>()
      const board = currentCore.value.current;
      board.areaId = formData.area?.id
      board.area = formData.area
      board.name = formData.name
      board.description = formData.description
      board.enabled = formData.enabled
      nextTick()
      const res = await currentCore.value?.update()
          .catch(() => {
          })
          .finally(() => {
            // fetchDashboard()
          })
      ElMessage({
        title: t('Success'),
        message: t('message.updatedSuccessfully'),
        type: 'success',
        duration: 2000
      });
    }
  })

}

const fetchDashboard = () => {
  emit('fetchDashboard')
}

const cancel = () => {
  push(`/dashboards`)
}

const removeBoard = async () => {
  if (!currentCore.value) return;
  await currentCore.value.removeBoard()
      .catch(() => {
      })
      .finally(() => {
      })
  cancel()
}

</script>

<template>

  <ElRow class="mb-10px">
    <ElCol>
      <ElDivider content-position="left">{{ $t('dashboard.mainTab') }}</ElDivider>
    </ElCol>
  </ElRow>

  <Form
      :schema="schema"
      :rules="rules"
      label-position="top"
      @register="register"
      class="mb-10px"
  />


  <ElRow class="mb-10px">
    <ElCol>
      <ElDivider class="mb-10px" content-position="left">{{ $t('main.actions') }}</ElDivider>
    </ElCol>
  </ElRow>


  <div class="text-right">

    <ElButton type="primary" @click.prevent.stop='exportDashbord' plain>
      <Icon icon="uil:file-export" class="mr-5px"/>
      {{ $t('main.export') }}
    </ElButton>


    <ElButton type="primary" @click.prevent.stop="updateBoard" plain>
      {{ $t('main.update') }}
    </ElButton>


    <ElButton @click.prevent.stop="fetchDashboard" plain>{{
        $t('main.loadFromServer')
      }}
    </ElButton>

<!--    <ElPopconfirm-->
<!--        :confirm-button-text="$t('main.ok')"-->
<!--        :cancel-button-text="$t('main.no')"-->
<!--        width="250"-->
<!--        style="margin-left: 10px;"-->
<!--        :title="$t('main.are_you_sure_to_do_want_this?')"-->
<!--        @confirm="cancel"-->
<!--    >-->
<!--      <template #reference>-->
<!--        <ElButton plain>-->
<!--          {{ t('main.cancel') }}-->
<!--        </ElButton>-->
<!--      </template>-->
<!--    </ElPopconfirm>-->

    <ElPopconfirm
        :confirm-button-text="$t('main.ok')"
        :cancel-button-text="$t('main.no')"
        width="250"
        style="margin-left: 10px;"
        :title="$t('main.are_you_sure_to_do_want_this?')"
        @confirm="removeBoard"
    >
      <template #reference>
        <ElButton type="danger" plain>
          <Icon icon="ep:delete" class="mr-5px"/>
          {{ t('main.remove') }}
        </ElButton>
      </template>
    </ElPopconfirm>

  </div>

  <!-- export dialog -->
  <Dialog v-model="dialogVisible" :title="t('main.dialogExportTitle')" :maxHeight="400" width="80%">
    <JsonViewer v-model="dialogSource"/>
    <template #footer>
      <ElButton @click="copy()">{{ t('setting.copy') }}</ElButton>
      <ElButton @click="dialogVisible = false">{{ t('main.closeDialog') }}</ElButton>
    </template>
  </Dialog>
  <!-- /export dialog -->

</template>
