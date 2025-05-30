<script setup lang="ts">
import {computed, onMounted, onUnmounted, PropType, reactive, ref, unref, watch} from 'vue'
import {Form} from '@/components/Form'
import {ElButton, ElCol, ElDivider, ElMessage, ElPopconfirm, ElRow} from 'element-plus'
import {useI18n} from '@/hooks/web/useI18n'
import {useForm} from '@/hooks/web/useForm'
import {useValidator} from '@/hooks/web/useValidator'
import {FormSchema} from '@/types/form'
import {Core, Tab, eventBus} from "@/views/Dashboard/core";
import {JsonViewer} from "@/components/JsonViewer";
import {Dialog} from "@/components/Dialog";
import FontEditor from "@/views/Dashboard/components/src/FontEditor.vue";
import TabListWindow from "@/views/Dashboard/editor/TabListWindow.vue";

const {register, elFormRef, methods} = useForm()
const {required} = useValidator()
const {t} = useI18n()
const {setValues} = methods

const props = defineProps({
  core: {
    type: Object as PropType<Nullable<Core>>,
    default: () => null
  },
  tab: {
    type: Object as PropType<Nullable<Tab>>,
    default: () => null
  },
})

const currentCore = computed(() => props.core as Core)

const activeTab = computed({
  get(): Tab {
    return currentCore.value.getActiveTab
  },
  set(val: Tab) {
  }
})

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
    field: 'icon',
    label: t('dashboard.icon'),
    component: 'Input',
    colProps: {
      span: 24
    },
    componentProps: {
      placeholder: t('dashboard.icon')
    }
  },
  {
    field: 'enabled',
    label: t('dashboard.enabled'),
    component: 'Switch',
    value: false,
    colProps: {
      md: 24,
      span: 24
    },
  },
  {
    field: 'appearance',
    label: t('dashboard.editor.appearanceOptions'),
    component: 'Divider',
    colProps: {
      span: 24
    },
  },
  {
    field: 'gap',
    label: t('dashboard.gap'),
    component: 'Switch',
    value: false,
    colProps: {
      md: 12,
      span: 12
    },
  },
  {
    field: 'columnWidth',
    label: t('dashboard.columnWidth'),
    component: 'InputNumber',
    value: 300,
    colProps: {
      span: 12
    },
  },
  {
    field: 'background',
    label: t('dashboard.background'),
    component: 'ColorPicker',
    colProps: {
      span: 12
    },
    componentProps: {
      placeholder: t('dashboard.background'),
    }
  },
  {
    field: 'backgroundAdaptive',
    label: t('dashboard.editor.backgroundAdaptive'),
    component: 'Switch',
    value: true,
    colProps: {
      span: 12
    },
    componentProps: {
      placeholder: t('dashboard.backgroundAdaptive'),
    }
  },
  {
    field: 'backgroundImage',
    label: t('dashboard.editor.image'),
    component: 'Image',
    colProps: {
      span: 24
    },
    value: null,
    componentProps: {
      placeholder: t('dashboard.image')
    }
  },
])

const eventHandler = (event: string, args: any[]) => {
  showExportDialog()
}

onMounted(() => {
  eventBus.subscribe('showTabExportDialog', eventHandler)
})

onUnmounted(() => {
  eventBus.unsubscribe('showTabExportDialog', eventHandler)
})

watch(
    () => props.tab,
    (val?: Tab) => {
      if (!val) return
      setValues({
        name: val.name,
        columnWidth: val.columnWidth,
        gap: val.gap,
        background: val.background,
        icon: val.icon,
        enabled: val.enabled,
        weight: val.weight,
        backgroundImage: val.backgroundImage || undefined,
        backgroundAdaptive: val.backgroundAdaptive,
      })
    },
    {
      deep: false,
      immediate: true
    }
)

// ---------------------------------
// common
// ---------------------------------
const updateTab = async () => {
  const formRef = unref(elFormRef)
  await formRef?.validate(async (isValid) => {
    if (isValid) {
      const {getFormData} = methods
      const formData = await getFormData()

      console.log(formData.backgroundImage)

      activeTab.value.background = formData.background;
      activeTab.value.backgroundImage = formData.backgroundImage || undefined;
      activeTab.value.backgroundAdaptive = formData.backgroundAdaptive;
      activeTab.value.columnWidth = formData.columnWidth;
      activeTab.value.enabled = formData.enabled;
      activeTab.value.gap = formData.gap;
      activeTab.value.icon = formData.icon;
      activeTab.value.name = formData.name;
      activeTab.value.weight = formData.weight;

      const res = await currentCore.value?.updateTab();
      if (res) {
        ElMessage({
          title: t('Success'),
          message: t('message.updatedSuccessfully'),
          type: 'success',
          duration: 2000
        });
      }
    }
  })

}

const removeTab = async () => {
  if (!activeTab.value) return;
  await currentCore.value.removeTab();
}

const cancel = () => {
  if (!activeTab.value) return;
  setValues({
    name: activeTab.value.name,
    columnWidth: activeTab.value.columnWidth,
    gap: activeTab.value.gap,
    background: activeTab.value.background,
    icon: activeTab.value.icon,
    enabled: activeTab.value.enabled,
    weight: activeTab.value.weight,
    backgroundImage: activeTab.value?.backgroundImage || undefined,
    backgroundAdaptive: activeTab.value?.backgroundAdaptive,
  })
}

// ---------------------------------
// import/export
// ---------------------------------

const dialogSource = ref({})
const exportDialogVisible = ref(false)

const prepareForExport = () => {
  if (currentCore.value.activeTabIdx == undefined) {
    return;
  }
  dialogSource.value = activeTab.value.serialize()
}

const showExportDialog = () => {
  prepareForExport()
  exportDialogVisible.value = true
}

</script>

<template>

  <ElRow class="mb-10px">
    <ElCol>
      <ElDivider content-position="left">{{ $t('dashboard.tabOptions') }}</ElDivider>
    </ElCol>
  </ElRow>

  <Form v-if="currentCore.tabs.length"
        :schema="schema"
        :rules="rules"
        label-position="top"
        @register="register"
  />

  <FontEditor v-if="activeTab" :tab="activeTab"/>

  <ElRow class="mb-10px">
    <ElCol>
      <ElDivider content-position="left">{{ $t('main.actions') }}</ElDivider>
    </ElCol>
  </ElRow>

  <div class="text-right" v-if="currentCore.tabs.length">
    <ElButton type="primary" @click.prevent.stop="updateTab" plain>{{ $t('main.update') }}</ElButton>
    <ElButton @click.prevent.stop="cancel" plain>{{ t('main.cancel') }}</ElButton>
    <ElPopconfirm
        :confirm-button-text="$t('main.ok')"
        :cancel-button-text="$t('main.no')"
        width="250"
        style="margin-left: 10px;"
        :title="$t('main.are_you_sure_to_do_want_this?')"
        @confirm="removeTab"
    >
      <template #reference>
        <ElButton class="mr-10px" type="danger" plain>
          <Icon icon="ep:delete" class="mr-5px"/>
          {{ t('main.remove') }}
        </ElButton>
      </template>
    </ElPopconfirm>
  </div>

  <!-- tab list window -->
  <TabListWindow :core="currentCore"/>
  <!-- /tab list window -->

  <!-- export dialog -->
  <Dialog v-model="exportDialogVisible" :title="t('main.dialogExportTitle')" :maxHeight="400" width="80%">
    <JsonViewer v-model="dialogSource"/>
  </Dialog>
  <!-- /export dialog -->

</template>

<style lang="less">

</style>
