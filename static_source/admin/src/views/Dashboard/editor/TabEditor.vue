<script setup lang="ts">
import {computed, PropType, reactive, ref, unref, watch} from 'vue'
import {Form} from '@/components/Form'
import {ElButton, ElCol, ElDivider, ElEmpty, ElMenu, ElMenuItem, ElMessage, ElPopconfirm, ElRow} from 'element-plus'
import {useI18n} from '@/hooks/web/useI18n'
import {useForm} from '@/hooks/web/useForm'
import {useValidator} from '@/hooks/web/useValidator'
import {FormSchema} from '@/types/form'
import {ApiDashboardCard, ApiEntity} from "@/api/stub";
import {Core, Tab} from "@/views/Dashboard/core/core";
import {DraggableContainer} from "@/components/DraggableContainer";

const {register, elFormRef, methods} = useForm()
const {required} = useValidator()
const {t} = useI18n()
const {setValues} = methods

interface DashboardTab {
  id: number;
  name: string;
  columnWidth: number;
  gap: boolean;
  background: string;
  icon: string;
  enabled: boolean;
  weight: number;
  dashboardId: number;
  cards: ApiDashboardCard[];
  entities: Map<string, ApiEntity>;
  dragEnabled: boolean;
}

const currentCore = ref<Core>(new Core())
const props = defineProps({
  core: {
    type: Object as PropType<Nullable<Core>>,
    default: () => null
  },
  tab: {
    type: Object as PropType<Nullable<DashboardTab>>,
    default: () => null
  },
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
      md: 24,
      span: 24
    },
  },
  {
    field: 'columnWidth',
    label: t('dashboard.columnWidth'),
    component: 'InputNumber',
    value: 300,
    colProps: {
      span: 24
    },
  },
  {
    field: 'background',
    label: t('dashboard.background'),
    component: 'ColorPicker',
    colProps: {
      span: 24
    },
    componentProps: {
      placeholder: t('dashboard.background'),
    }
  },
])

watch(
    () => props.tab,
    (val?: DashboardTab) => {
      if (!val) return
      setValues({
        name: val.name,
        columnWidth: val.columnWidth,
        gap: val.gap,
        background: val.background,
        icon: val.icon,
        enabled: val.enabled,
        weight: val.weight,
        dragEnabled: val.dragEnabled,
      })
    },
    {
      deep: false,
      immediate: true
    }
)

watch(
    () => props.core,
    (val?: Core) => {
      if (!val) return
      currentCore.value = val
    },
    {
      deep: false,
      immediate: true
    }
)


const activeTab = computed({
  get(): Tab {
    return currentCore.value.getActiveTab as Tab
  },
  set(val: Tab) {
  }
})

// ---------------------------------
// common
// ---------------------------------
const updateTab = async () => {
  const formRef = unref(elFormRef)
  await formRef?.validate(async (isValid) => {
    if (isValid) {
      const {getFormData} = methods
      const formData = await getFormData<DashboardTab>()

      activeTab.value.background = formData.background;
      activeTab.value.columnWidth = formData.columnWidth;
      activeTab.value.dragEnabled = formData.dragEnabled;
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

const menuTabClick = (index: number, tab: Tab) => {
  currentCore.value.selectTabInMenu(index)
}

const createTab = async () => {
  await currentCore.value.createTab();

  ElMessage({
    title: t('Success'),
    message: t('message.createdSuccessfully'),
    type: 'success',
    duration: 2000
  });
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
    dragEnabled: activeTab.value.dragEnabled,
  })
}

const sortCardUp = (tab: Tab, index: number) => {
}
const sortCardDown = (tab: Tab, index: number) => {
}

</script>

<template>

  <ElRow class="mb-10px">
    <ElCol>
      <ElDivider content-position="left">{{ $t('dashboard.tabOptions') }}</ElDivider>
    </ElCol>
  </ElRow>


  <!--  <ElContainer>-->
  <!--    <ElMain>-->
  <!--      <ElCard class="box-card">-->
  <!--        <template #header>-->
  <!--          <div class="card-header">-->
  <!--            <span>{{ $t('dashboard.tabDetail') }}</span>-->
  <!--          </div>-->
  <!--        </template>-->

  <Form v-if="currentCore.tabs.length"
        :schema="schema"
        :rules="rules"
        label-position="top"
        @register="register"
  />

  <ElEmpty v-if="!currentCore.tabs.length" :rows="5">
    <ElButton type="primary" @click="createTab()">
      {{ t('dashboard.addNewTab') }}
    </ElButton>
  </ElEmpty>

  <ElRow class="mb-10px">
    <ElCol>
      <ElDivider content-position="left">{{ $t('main.actions') }}</ElDivider>
    </ElCol>
  </ElRow>

  <div class="text-right" v-if="currentCore.tabs.length">

    <ElButton type="primary" @click.prevent.stop="updateTab">{{ $t('main.update') }}</ElButton>


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

  <DraggableContainer :name="'editor-tabs'" :initial-width="280" :min-width="280">
    <template #header>
      <span>Tabs</span>
    </template>
    <template #default>

<!--      <ElRow class="mb-10px mt-10px">-->
<!--        <ElCol>-->
<!--          <ElDivider content-position="left">{{ $t('dashboard.tabList') }}</ElDivider>-->
<!--        </ElCol>-->
<!--      </ElRow>-->

      <ElRow class="mb-10px mt-10px">
        <ElCol>
          <ElButton class="w-[100%]" @click="createTab()" size="small">
            {{ t('dashboard.addNewTab') }}
          </ElButton>
        </ElCol>
      </ElRow>

      <ElMenu v-if="currentCore.tabs.length" :default-active="currentCore.activeTabIdx + ''"
              v-model="currentCore.activeTabIdx" class="el-menu-vertical-demo">
        <ElMenuItem :index="index + ''" :key="tab" v-for="(tab, index) in currentCore.tabs"
                    @click="menuTabClick(index, tab)">
          <div class="w-[100%] card-header">
            <span>{{ tab.name }}</span>
          </div>
        </ElMenuItem>
      </ElMenu>
    </template>
  </DraggableContainer>

</template>

<style lang="less" scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.hide {
  display: none;
}

.el-menu-item:hover .hide {
  display: block;
  color: red;
}
</style>
