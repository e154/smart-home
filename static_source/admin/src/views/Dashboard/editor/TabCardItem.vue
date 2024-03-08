<script setup lang="ts">
import {computed, onMounted, onUnmounted, PropType, ref, watch} from 'vue'
import {useI18n} from '@/hooks/web/useI18n'
import {Card, CardItem, Core, eventBus, requestCurrentState} from "@/views/Dashboard/core";
import {CardEditorName, CardItemList} from "@/views/Dashboard/card_items";
import {JsonViewer} from "@/components/JsonViewer";
import {
  ElButton,
  ElCascader,
  ElCol,
  ElCollapse,
  ElCollapseItem,
  ElDivider,
  ElEmpty,
  ElForm,
  ElFormItem,
  ElInput,
  ElMessage,
  ElPopconfirm,
  ElRow,
} from 'element-plus'
import {JsonEditor} from "@/components/JsonEditor";
import {Dialog} from "@/components/Dialog";
import {ApiDashboardCardItem} from "@/api/stub";

const {t} = useI18n()

const cardItem = ref<CardItem>(null)
// const card = ref<Card>({} as Card)
const itemTypes = CardItemList;
const itemProps = {
  expandTrigger: 'hover' as const,
}
const cardItemType = computed(() => cardItem.value?.type)
const handleTypeChanged = (value: string[]) => {
  cardItem.value.type = value[value.length - 1]
}

const props = defineProps({
  core: {
    type: Object as PropType<Core>,
  },
  card: {
    type: Object as PropType<Nullable<Card>>,
    default: () => null
  },
})

const activeCard = computed({
  get(): Card {
    return props.card as Card
  },
  set(val: Card) {
  }
})

const currentCore = computed({
  get(): Core {
    return props.core as Core
  },
  set(val: Core) {
  }
})

watch(
    () => props.card,
    (val?: Card) => {
      if (!val) return
      activeCard.value = val
      if (val?.selectedItem > -1) {
        cardItem.value = val?.items[val?.selectedItem] || null
      } else {
        cardItem.value = null
      }

    },
    {
      deep: true,
      immediate: true
    }
)

// ---------------------------------
// import/export
// ---------------------------------

const addCardItem = () => {
  currentCore.value.createCardItem(undefined, 'text');
}

const removeCardItem = (index: number) => {
  currentCore.value.removeCardItem(index);
}

const duplicate = () => {
  activeCard.value.copyItem(activeCard.value.selectedItem);
}

const getCardEditorName = (name: string) => {
  return CardEditorName(name);
}

const updateCardItem = async () => {
  const {data} = await currentCore.value.updateCard();

  if (data) {
    ElMessage({
      title: t('Success'),
      message: t('message.updatedSuccessfully'),
      type: 'success',
      duration: 2000
    });
  }
}

const updateCurrentState = () => {
  if (cardItem.value.entityId) {
    requestCurrentState(cardItem.value?.entityId)
  }
}

const eventHandler = (event: string, args: any[]) => {
  switch (event) {
    case 'showCardItemImportDialog':
      cardIdForImport.value = args
      importDialogVisible.value = true
      break;
    case 'showCardItemExportDialog':
      showExportDialog(args)
      break;
  }
}
const cardIdForImport = ref<number>()
onMounted(() => {
  eventBus.subscribe(['showCardItemImportDialog', 'showCardItemExportDialog'], eventHandler)
})

onUnmounted(() => {
  eventBus.unsubscribe(['showCardItemImportDialog', 'showCardItemExportDialog'], eventHandler)
})

// ---------------------------------
// import/export
// ---------------------------------

const dialogSource = ref({})
const importDialogVisible = ref(false)
const exportDialogVisible = ref(false)
const importedCardItem = ref(null)

const prepareForExport = (cardItemId?: number) => {
  if (currentCore.value.activeCard == undefined) {
    return;
  }
  dialogSource.value = currentCore.value.serializeCardItem(cardItemId)
}

const showExportDialog = (cardItemId?: number) => {
  prepareForExport(cardItemId)
  exportDialogVisible.value = true
}

const importHandler = (val: any) => {
  if (importedCardItem.value == val) {
    return
  }
  importedCardItem.value = val
}

const importCardItem = async () => {
  let cardItem: ApiDashboardCardItem
  try {
    if (importedCardItem.value?.json) {
      cardItem = importedCardItem.value.json as ApiDashboardCardItem;
    } else if (importedCardItem.value.text) {
      cardItem = JSON.parse(importedCardItem.value.text) as ApiDashboardCardItem;
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
  const res = await currentCore.value.importCardItem(cardIdForImport.value, cardItem);
  if (res) {
    cardIdForImport.value = undefined
    ElMessage({
      title: t('Success'),
      message: t('message.importedSuccessful'),
      type: 'success',
      duration: 2000
    })
  }
  importDialogVisible.value = false
}


</script>

<template>

  <ElRow class="mb-10px" v-if="activeCard.selectedItem !== -1">
    <ElCol>
      <ElDivider content-position="left">{{ $t('dashboard.cardItemOptions') }}</ElDivider>
    </ElCol>
  </ElRow>

  <ElForm
      v-if="cardItem"
      :model="cardItem"
      label-position="top"
      style="width: 100%"
      ref="cardItemForm"
  >

    <ElRow>
      <ElCol>

        <ElFormItem :label="$t('dashboard.editor.type')" prop="type">
          <ElCascader
              v-model="cardItemType"
              :options="itemTypes"
              :props="itemProps"
              :placeholder="$t('dashboard.editor.pleaseSelectType')"
              style="width: 100%"
              @change="handleTypeChanged"
          />
        </ElFormItem>
      </ElCol>
    </ElRow>

    <ElRow>
      <ElCol>
        <ElFormItem :label="$t('dashboard.editor.title')" prop="title">
          <ElInput v-model="cardItem.title"/>
        </ElFormItem>
      </ElCol>

    </ElRow>

    <component
        :is="getCardEditorName(cardItem.type)"
        :core="core"
        :item="cardItem"
    />
  </ElForm>

  <ElEmpty v-if="!activeCard.items.length || activeCard.selectedItem === -1" :rows="5" class="mt-20px mb-20px"
           description="Select card item or">
    <ElButton type="primary" @click="addCardItem()">
      {{ t('dashboard.editor.addNewCardItem') }}
    </ElButton>
    <ElButton type="primary" @click="importDialogVisible = true">
      {{ t('main.import') }}
    </ElButton>
  </ElEmpty>

  <ElRow class="mb-10px mt-10px" v-if="activeCard.selectedItem > -1 && cardItem.entity">
    <ElCol>
      <ElCollapse>
        <ElCollapseItem :title="$t('dashboard.editor.eventstateJSONobject')">
          <ElButton class="mb-10px w-[100%]" @click.prevent.stop="updateCurrentState()">
            <Icon icon="ep:refresh" class="mr-5px"/>
            {{ $t('dashboard.editor.getEvent') }}
          </ElButton>
          <JsonViewer v-model="cardItem.lastEvent"/>
        </ElCollapseItem>
      </ElCollapse>
    </ElCol>
  </ElRow>

  <ElRow v-if="activeCard.selectedItem > -1" class="mb-10px">
    <ElCol>
      <ElDivider class="mb-10px" content-position="left">{{ $t('main.actions') }}</ElDivider>
    </ElCol>
  </ElRow>

  <div v-if="activeCard.selectedItem > -1" class="text-right">

    <ElButton type="primary" @click.prevent.stop="updateCardItem" plain>{{
        $t('main.update')
      }}
    </ElButton>

    <ElButton @click.prevent.stop="duplicate">{{ $t('main.duplicate') }}</ElButton>

    <ElPopconfirm
        :confirm-button-text="$t('main.ok')"
        :cancel-button-text="$t('main.no')"
        width="250"
        style="margin-left: 10px;"
        :title="$t('main.are_you_sure_to_do_want_this?')"
        @confirm="removeCardItem(activeCard.selectedItem)"
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
  <Dialog v-model="exportDialogVisible" :title="t('main.dialogExportTitle')" :maxHeight="400" width="80%">
    <JsonViewer v-model="dialogSource"/>
  </Dialog>
  <!-- /export dialog -->

  <!-- import dialog -->
  <Dialog v-model="importDialogVisible" :title="t('main.dialogImportTitle')" :maxHeight="400" width="80%"
          custom-class>
    <JsonEditor @change="importHandler"/>
    <template #footer>
      <ElButton type="primary" @click="importCardItem()" plain>{{ t('main.import') }}</ElButton>
      <ElButton @click="importDialogVisible = false">{{ t('main.closeDialog') }}</ElButton>
    </template>
  </Dialog>
  <!-- /import dialog -->

</template>

<style lang="less">

</style>
