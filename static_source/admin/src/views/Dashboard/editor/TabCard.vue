<script setup lang="ts">
import {computed, PropType, reactive, ref, unref, watch} from 'vue'
import {Form} from '@/components/Form'
import {
  ElButton,
  ElCard,
  ElMessage,
  ElPopconfirm,
  ElSkeleton,
  ElMenu,
  ElMenuItem,
  ElButtonGroup,
  ElContainer,
  ElAside,
  ElMain,
  ElScrollbar,
  ElEmpty,
  ElDivider,
  ElCol,
  ElRow
} from 'element-plus'
import {useI18n} from '@/hooks/web/useI18n'
import {useForm} from '@/hooks/web/useForm'
import {useValidator} from '@/hooks/web/useValidator'
import {FormSchema} from '@/types/form'
import {ApiDashboard, ApiDashboardCard, ApiDashboardCardItem, ApiEntity} from "@/api/stub";
import {copyToClipboard} from "@/utils/clipboard";
import JsonViewer from "@/components/JsonViewer/JsonViewer.vue";
import {Card, Core, Tab} from "@/views/Dashboard/core";
import {useBus} from "@/views/Dashboard/bus";
import { Dialog } from '@/components/Dialog'
import JsonEditor from "@/components/JsonEditor/JsonEditor.vue";
import KeystrokeCapture from "@/views/Dashboard/components/KeystrokeCapture.vue";

const {register, elFormRef, methods} = useForm()
const {required} = useValidator()
const {t} = useI18n()

const {setValues} = methods
const {bus} = useBus()

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

export interface DashboardCard {
  id: number;
  title: string;
  height: number;
  width: number;
  background: string;
  weight: number;
  enabled: boolean;
  dashboardTabId: number;
  payload?: string;
  items: ApiDashboardCardItem[];
  entities: Map<string, ApiEntity>;
  hidden: boolean;
  entityId: string;
}

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
    field: 'title',
    label: t('dashboard.editor.name'),
    component: 'Input',
    colProps: {
      span: 24
    },
    componentProps: {
      placeholder: t('dashboard.editor.name')
    }
  },
  {
    field: 'enabled',
    label: t('dashboard.enabled'),
    component: 'Switch',
    value: false,
    colProps: {
      md: 12,
      span: 24
    },
  },
  {
    field: 'hidden',
    label: t('dashboard.hidden'),
    component: 'Switch',
    value: false,
    colProps: {
      md: 12,
      span: 24
    },
  },
  {
    field: 'cardSize',
    label: t('dashboard.editor.size'),
    component: 'Divider',
    colProps: {
      span: 24
    },
  },
  {
    field: 'height',
    label: t('dashboard.height'),
    component: 'InputNumber',
    value: 300,
    colProps: {
      md: 12,
      span: 24
    },
  },
  {
    field: 'width',
    label: t('dashboard.width'),
    component: 'InputNumber',
    value: 300,
    colProps: {
      md: 12,
      span: 24
    },
  },
  {
    field: 'cardSize',
    label: t('dashboard.editor.color'),
    component: 'Divider',
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

const currentCore = computed(() => props.core as Core)
const activeCard = computed(() => props.core?.getActiveTab?.cards[props.core?.activeCard] as Card)

//todo: optimize
watch(
    () => currentCore.value.activeCard,
    (val?: number) => {
      if (!(val >= 0)) return
      const card = props.tab?.cards[val] as DashboardCard
      setValues({
        title: card.title,
        enabled: card.enabled,
        hidden: card.hidden,
        // dragEnabled: card.dragEnabled, // todo: fix
        height: card.height,
        weight: card.weight,
        width: card.width,
        background: card.background,
      })
    },
    {
      deep: true,
      immediate: true
    }
)


const activeTab = computed({
  get(): Tab {
    return currentCore.value.getActiveTab as Tab
  },
  set(val: Tab) {}
})

// ---------------------------------
// import/export
// ---------------------------------

const dialogSource = ref({})
const importDialogVisible = ref(false)
const exportDialogVisible = ref(false)
const importedCard = ref(null)

const prepareForExport = () => {
  if (currentCore.value.activeCard == undefined) {
    return;
  }
  dialogSource.value = activeTab.value.cards[currentCore.value.activeCard].serialize()
}

const showImportDialog = () => {
  importDialogVisible.value = true
}

const showExportDialog = () => {
  prepareForExport()
  exportDialogVisible.value = true
}

const importHandler = (val: any) => {
  if (importedCard.value == val) {
    return
  }
  importedCard.value = val
}

const copy = async () => {
  prepareForExport()
  copyToClipboard(JSON.stringify(dialogSource.value, null, 2))
}

const importCard = async () => {
  let card: ApiDashboardCard
  try {
    if (importedCard.value?.json) {
      card = importedCard.value.json as ApiDashboardCard;
    } else if(importedCard.value.text) {
      card = JSON.parse(importedCard.value.text) as ApiDashboardCard;
    }
  } catch {
    ElMessage({
      title: t('Error'),
      message: t('message.corruptedJsonFormat'),
      type: 'error',
      duration: 2000
    });
    return
  }
  const res = await currentCore.value.importCard(card);
  if (res) {
    ElMessage({
      title: t('Success'),
      message: t('message.importedSuccessful'),
      type: 'success',
      duration: 2000
    })
  }
  importDialogVisible.value = false
}

// ---------------------------------
// common
// ---------------------------------

const onSelectedCard = (id: number) => {
  currentCore.value.onSelectedCard(id);
  bus.emit('unselected_card_item')
}

const addCard = () => {
  currentCore.value.createCard();
}

const updateCard = async () => {
  const formRef = unref(elFormRef)
  await formRef?.validate(async (isValid) => {
    if (isValid) {
      const {getFormData} = methods
      const formData = await getFormData<DashboardCard>()

      activeCard.value.title = formData.title
      activeCard.value.enabled = formData.enabled
      activeCard.value.hidden = formData.hidden
      // activeCard.value.dragEnabled = formData.dragEnabled
      activeCard.value.height = formData.height
      activeCard.value.weight = formData.weight
      activeCard.value.width = formData.width
      activeCard.value.background = formData.background

      const res = await currentCore.value.updateCard();
      currentCore.value.updateCurrentTab();
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

const removeCard = async () => {
  await currentCore.value.removeCard();
}

useBus({
  name: 'selected_card',
  callback: (id: number) => onSelectedCard(id)
})

const menuCardsClick = (card: DashboardCard) => {
  // bus.emit('selected_card', card.id);
  onSelectedCard(card.id)
}

const cancel = () => {
  if (!activeCard.value) return;
  setValues({
    title: activeCard.value.title,
    enabled: activeCard.value.enabled,
    hidden: activeCard.value.hidden,
    // dragEnabled: activeCard.value.dragEnabled,
    height: activeCard.value.height,
    weight: activeCard.value.weight,
    width: activeCard.value.width,
    background: activeCard.value.background,
  })
}

const sortCardUp = (card: Card, index: number) => {
  activeTab.value.sortCardUp(card, index)
  currentCore.value.updateCurrentTab();
}

const sortCardDown = (card: Card, index: number) => {
  activeTab.value.sortCardDown(card, index)
  currentCore.value.updateCurrentTab();
}


</script>

<template>

<!--  <ElContainer style="height: 500px">-->
  <ElContainer>
    <ElMain>
      <ElScrollbar>
        <ElCard class="box-card">
          <template #header>
            <div class="card-header">
              <span>{{ $t('dashboard.cardDetail') }}</span>
            </div>
          </template>

          <Form
              v-if="core.activeCard >= 0"
              :schema="schema"
              :rules="rules"
              label-position="top"
              style="width: 100%"
              @register="register"
          />

          <ElRow v-if="core.activeCard >= 0">
            <ElCol>
              <ElDivider content-position="left">{{ $t('dashboard.editor.image') }}</ElDivider>
            </ElCol>
            <ElCol>
              <KeystrokeCapture/>
            </ElCol>
          </ElRow>

          <ElEmpty v-if="!(core.activeCard >= 0)" :rows="5">
            <ElButton type="primary" @click="addCard()">
              {{ t('dashboard.addNewCard') }}
            </ElButton>
          </ElEmpty>

          <div class="text-right" v-if="core.activeCard >= 0">
            <ElButton type="primary" @click.prevent.stop='showExportDialog()'>
              <Icon icon="uil:file-export" class="mr-5px"/>
              {{ $t('main.export') }}
            </ElButton>
            <ElButton type="primary" @click.prevent.stop="updateCard">{{ $t('main.update') }}</ElButton>
            <ElButton type="default" @click.prevent.stop="cancel" plain>{{ t('main.cancel') }}</ElButton>
            <ElPopconfirm
                :confirm-button-text="$t('main.ok')"
                :cancel-button-text="$t('main.no')"
                width="250"
                style="margin-left: 10px;"
                :title="$t('main.are_you_sure_to_do_want_this?')"
                @confirm="removeCard"
            >
              <template #reference>
                <ElButton class="mr-10px" type="danger" plain>
                  <Icon icon="ep:delete" class="mr-5px"/>
                  {{ t('main.remove') }}
                </ElButton>
              </template>
            </ElPopconfirm>
          </div>

        </ElCard>
      </ElScrollbar>
    </ElMain>
    <ElAside width="400px">
      <ElScrollbar>
        <ElCard class="box-card">
          <template #header>
            <div class="card-header">
              <span>{{ $t('dashboard.cardList') }}</span>
              <ElButtonGroup>
                <ElButton @click="addCard()" text size="small">
                  {{ t('dashboard.addNew') }}
                </ElButton>
                <ElButton @click="showImportDialog()" text size="small">
                  {{ t('dashboard.importCard') }}
                </ElButton>
              </ElButtonGroup>
            </div>
          </template>
          <ElMenu v-if="currentCore.activeTabIdx > -1 && activeTab.cards.length" :default-active="currentCore.activeCard + ''" v-model="currentCore.activeCard" class="el-menu-vertical-demo">
            <ElMenuItem :index="index + ''" :key="index" v-for="(card, index) in activeTab.cards" @click="menuCardsClick(card)">
              <div class="w-[100%] card-header">
                <span>{{ card.title }}</span>
                <ElButtonGroup class="hide">
                  <ElButton type="default" @click.prevent.stop="sortCardUp(card, index)">
                    <Icon icon="teenyicons:up-solid" />
                  </ElButton>
                  <ElButton type="default" @click.prevent.stop="sortCardDown(card, index)">
                    <Icon icon="teenyicons:down-solid" />
                  </ElButton>
                </ElButtonGroup>
              </div>
            </ElMenuItem>
          </ElMenu>

        </ElCard>
      </ElScrollbar>
    </ElAside>
  </ElContainer>

  <!-- export dialog -->
  <Dialog v-model="exportDialogVisible" :title="t('entities.dialogExportTitle')" :maxHeight="400" width="80%">
    <JsonViewer v-model="dialogSource"/>
    <template #footer>
      <ElButton @click="copy()">{{ t('setting.copy') }}</ElButton>
      <ElButton @click="exportDialogVisible = false">{{ t('main.closeDialog') }}</ElButton>
    </template>
  </Dialog>
  <!-- /export dialog -->

  <!-- import dialog -->
  <Dialog v-model="importDialogVisible" :title="t('entities.dialogImportTitle')" :maxHeight="400" width="80%" custom-class>
    <JsonEditor @change="importHandler"/>
    <template #footer>
      <ElButton type="primary" @click="importCard()" plain>{{ t('main.import') }}</ElButton>
      <ElButton @click="importDialogVisible = false">{{ t('main.closeDialog') }}</ElButton>
    </template>
  </Dialog>
  <!-- /import dialog -->
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
