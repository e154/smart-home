<script setup lang="ts">
import {computed, onMounted, PropType, reactive, ref, unref, watch} from 'vue'
import {Form} from '@/components/Form'
import {
  ElButton,
  ElButtonGroup,
  ElCol,
  ElDivider,
  ElEmpty,
  ElForm,
  ElFormItem,
  ElIcon,
  ElMenu,
  ElMenuItem,
  ElMessage,
  ElPopconfirm,
  ElRow,
  ElSwitch,
} from 'element-plus'
import {CloseBold} from '@element-plus/icons-vue'
import {useI18n} from '@/hooks/web/useI18n'
import {useForm} from '@/hooks/web/useForm'
import {useValidator} from '@/hooks/web/useValidator'
import {FormSchema} from '@/types/form'
import {ApiDashboardCard, ApiDashboardCardItem, ApiEntity} from "@/api/stub";
import {JsonViewer} from "@/components/JsonViewer";
import {Card, Core, Tab, useBus} from "@/views/Dashboard/core";
import {Dialog} from '@/components/Dialog'
import {JsonEditor} from "@/components/JsonEditor";
import {FrameEditor, KeystrokeCapture} from "@/views/Dashboard/components";
import {DraggableContainer} from "@/components/DraggableContainer";

const {register, elFormRef, methods} = useForm()
const {required} = useValidator()
const {t} = useI18n()

const {setValues} = methods
const {emit} = useBus()

export interface DashboardCard {
  id: number;
  title: string;
  height: number;
  width: number;
  background: string;
  backgroundAdaptive: boolean;
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
  },
  tab: {
    type: Object as PropType<Nullable<Tab>>,
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
      span: 12
    },
  },
  {
    field: 'hidden',
    label: t('dashboard.hidden'),
    component: 'Switch',
    value: false,
    colProps: {
      md: 12,
      span: 12
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
    field: 'height',
    label: t('dashboard.height'),
    component: 'InputNumber',
    value: 300,
    colProps: {
      md: 12,
      span: 12
    },
  },
  {
    field: 'width',
    label: t('dashboard.width'),
    component: 'InputNumber',
    value: 300,
    colProps: {
      md: 12,
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
      md: 12,
      span: 12
    },
    componentProps: {
      placeholder: t('dashboard.backgroundAdaptive'),
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
        backgroundAdaptive: card.backgroundAdaptive,
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
  set(val: Tab) {
  }
})

// ---------------------------------
// import/export
// ---------------------------------

const dialogSource = ref({})
const importDialogVisible = ref(false)
const exportDialogVisible = ref(false)
const importedCard = ref(null)

const prepareForExport = (cardId?: number) => {
  if (currentCore.value.activeCard == undefined) {
    return;
  }
  dialogSource.value = currentCore.value.serializeCard(cardId)
}

const showImportDialog = () => {
  importDialogVisible.value = true
}

const showExportDialog = (cardId?: number) => {
  prepareForExport(cardId)
  exportDialogVisible.value = true
}

const importHandler = (val: any) => {
  if (importedCard.value == val) {
    return
  }
  importedCard.value = val
}

// const pasteCardItem = () => {
//   activeCard.value.pasteCardItem();
// }

const importCard = async () => {
  let card: ApiDashboardCard
  try {
    if (importedCard.value?.json) {
      card = importedCard.value.json as ApiDashboardCard;
    } else if (importedCard.value.text) {
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
  emit('unselected_card_item')
}

const addCard = () => {
  currentCore.value.createCard();
}

const updateCard = async () => {
  const formRef = unref(elFormRef)
  await formRef?.validate(async (isValid) => {
    if (isValid) {
      const {getFormData} = methods
      const formData = await getFormData()

      activeCard.value.title = formData.title
      activeCard.value.enabled = formData.enabled
      activeCard.value.hidden = formData.hidden
      // activeCard.value.dragEnabled = formData.dragEnabled
      activeCard.value.height = formData.height
      activeCard.value.weight = formData.weight
      activeCard.value.width = formData.width
      activeCard.value.background = formData.background
      activeCard.value.backgroundAdaptive = formData.backgroundAdaptive

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
  // emit('selected_card', card.id);
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
    backgroundAdaptive: activeCard.value.backgroundAdaptive,
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

const showMenuWindow = ref(false)
onMounted(() => {
  useBus({
    name: 'toggleMenu',
    callback: (menu: string) => {
      if (menu !== 'cards') {
        return
      }
      // console.log("cards", menu)
      showMenuWindow.value = !showMenuWindow.value
    }
  })

  useBus({
    name: 'showCardImportDialog',
    callback: () => {
      importDialogVisible.value = true
    }
  })

  useBus({
    name: 'showCardExportDialog',
    callback: (cardId?: number) => {
      showExportDialog(cardId)
    }
  })

})
</script>

<template>
  <ElEmpty v-if="!(activeCard !== undefined)" :rows="5" description="Select card or">
    <ElButton type="primary" @click="addCard()">
      {{ t('dashboard.addNewCard') }}
    </ElButton>
    <ElButton type="primary" @click="showImportDialog()">
      {{ t('main.import') }}
    </ElButton>
  </ElEmpty>

  <div v-if="activeCard !== undefined">

    <ElRow class="mb-10px">
      <ElCol>
        <ElDivider content-position="left">{{ $t('dashboard.cardOptions') }}</ElDivider>
      </ElCol>
    </ElRow>

    <Form
        :schema="schema"
        :rules="rules"
        label-position="top"
        style="width: 100%"
        @register="register"
    />

    <ElForm
        label-position="top"
        style="width: 100%"
    >
      <ElRow
      >
        <ElCol>
          <ElFormItem :label="$t('dashboard.editor.template')" prop="template">
            <ElSwitch v-model="activeCard.template"/>
          </ElFormItem>
        </ElCol>
      </ElRow>

    </ElForm>

    <FrameEditor v-if="activeCard.template" :card="activeCard" :core="core"/>

    <ElRow class="mb-10px">
      <ElCol>
        <ElDivider content-position="left">{{ $t('dashboard.editor.keystrokeCapture') }}</ElDivider>
      </ElCol>
    </ElRow>

    <ElRow class="mb-10px">
      <ElCol>
        <KeystrokeCapture :card="activeCard" :core="core"/>
      </ElCol>
    </ElRow>

    <ElRow class="mb-10px">
      <ElCol>
        <ElDivider content-position="left">{{ $t('main.actions') }}</ElDivider>
      </ElCol>
    </ElRow>

    <div class="text-right">
      <!--      <ElButton type="primary" @click.prevent.stop='showExportDialog()' plain>-->
      <!--        <Icon icon="uil:file-export" class="mr-5px"/>-->
      <!--        {{ $t('main.export') }}-->
      <!--      </ElButton>-->
      <ElButton type="primary" @click.prevent.stop="updateCard" plain>{{ $t('main.update') }}</ElButton>
      <!--    <ElButton @click.prevent.stop="pasteCardItem">{{ $t('dashboard.pasteCardItem') }}</ElButton>-->
      <ElButton @click.prevent.stop="cancel" plain>{{ t('main.cancel') }}</ElButton>
      <ElPopconfirm
          :confirm-button-text="$t('main.ok')"
          :cancel-button-text="$t('main.no')"
          width="250"
          style="margin-left: 10px;"
          :title="$t('main.are_you_sure_to_do_want_this?')"
          @confirm="removeCard"
      >
        <template #reference>
          <ElButton type="danger" plain>
            <Icon icon="ep:delete" class="mr-5px"/>
            {{ t('main.remove') }}
          </ElButton>
        </template>
      </ElPopconfirm>
    </div>
  </div>

  <!-- export dialog -->
  <Dialog v-model="exportDialogVisible" :title="t('main.dialogExportTitle')" :maxHeight="400" width="80%">
    <JsonViewer v-model="dialogSource"/>
  </Dialog>
  <!-- /export dialog -->

  <!-- import dialog -->
  <Dialog v-model="importDialogVisible" :title="t('main.dialogImportTitle')" :maxHeight="400" width="80%"
          custom-class>
    <JsonEditor @change="importHandler"/>
    <template #footer>
      <ElButton type="primary" @click="importCard()" plain>{{ t('main.import') }}</ElButton>
      <ElButton @click="importDialogVisible = false">{{ t('main.closeDialog') }}</ElButton>
    </template>
  </Dialog>
  <!-- /import dialog -->

  <DraggableContainer :name="'editor-cards'" :initial-width="280" :min-width="280" v-show="showMenuWindow">
    <template #header>
      <div class="w-[100%]">
        <div style="float: left">Cards</div>
        <div style="float: right; text-align: right">
          <a href="#" @click.prevent.stop='showMenuWindow= false'>
            <ElIcon class="mr-5px">
              <CloseBold/>
            </ElIcon>
          </a>
        </div>
      </div>
    </template>
    <template #default>

      <ElMenu v-if="currentCore.activeTabIdx > -1 && activeTab.cards.length"
              :default-active="currentCore.activeCard + ''" v-model="currentCore.activeCard"
              class="el-menu-vertical-demo">
        <ElMenuItem :index="index + ''" :key="index" v-for="(card, index) in activeTab.cards"
                    @click="menuCardsClick(card)">
          <div class="w-[100%] menu-item">
            <span>{{ card.title }}</span>
            <ElButtonGroup class="buttons">
              <ElButton @click.prevent.stop="sortCardUp(card, index)" text size="small">
                <Icon icon="teenyicons:up-solid"/>
              </ElButton>
              <ElButton @click.prevent.stop="sortCardDown(card, index)" text size="small">
                <Icon icon="teenyicons:down-solid"/>
              </ElButton>
            </ElButtonGroup>
          </div>
        </ElMenuItem>
      </ElMenu>

    </template>
  </DraggableContainer>


</template>

<style lang="less">

</style>
