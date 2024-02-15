<script setup lang="ts">
import {computed, onMounted, PropType, ref, watch} from 'vue'
import {CloseBold} from "@element-plus/icons-vue";
import {useI18n} from '@/hooks/web/useI18n'
import {Card, CardItem, Core, requestCurrentState} from "@/views/Dashboard/core/core";
import {useBus} from "@/views/Dashboard/core/bus";
import {CardEditorName, CardItemList} from "@/views/Dashboard/card_items";
import {JsonViewer} from "@/components/JsonViewer";
import {DraggableContainer} from "@/components/DraggableContainer";
import {
  ElButton,
  ElButtonGroup,
  ElCol,
  ElCollapse,
  ElCollapseItem,
  ElDivider,
  ElEmpty,
  ElForm,
  ElFormItem,
  ElIcon,
  ElInput,
  ElMenu,
  ElMenuItem,
  ElMessage,
  ElOption,
  ElPopconfirm,
  ElRow,
  ElSelect,
  ElTag
} from 'element-plus'
import {JsonEditor} from "@/components/JsonEditor";
import {Dialog} from "@/components/Dialog";
import {ApiDashboardCardItem} from "@/api/stub";

const {t} = useI18n()

const {emit} = useBus()

const cardItem = ref<CardItem>(null)
// const card = ref<Card>({} as Card)
const itemTypes = CardItemList;
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
  currentCore.value.createCardItem();
}

const removeCardItem = (index: number) => {
  currentCore.value.removeCardItem(index);
}

const duplicate = () => {
  activeCard.value.copyItem(activeCard.value.selectedItem);
}

const copyToClipboard = async () => {
  cardItem.value.copyToClipboard()
}

const menuCardItemClick = (index: number) => {
  if (currentCore.value.activeTabIdx < 0 || currentCore.value.activeCard == undefined) {
    return;
  }

  activeCard.value.selectedItem = index;

  emit('selected_card_item', index)
}

const sortCardItemUp = (item: CardItem, index: number) => {
  activeCard.value.sortCardItemUp(item, index)
  currentCore.value.updateCard();
}

const sortCardItemDown = (item: CardItem, index: number) => {
  activeCard.value.sortCardItemDown(item, index)
  currentCore.value.updateCard();
}

const getCardEditorName = (name: string) => {
  return CardEditorName(name);
}

// const cancel = () => {
//   console.warn('action not implemented')
// }

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

const showMenuWindow = ref(false)
onMounted(() => {
  useBus({
    name: 'toggleMenu',
    callback: (menu: string) => {
      if (menu !== 'cardItems') {
        return
      }
      console.log("cards", menu)
      showMenuWindow.value = !showMenuWindow.value
    }
  })
})

// ---------------------------------
// import/export
// ---------------------------------

const dialogSource = ref({})
const importDialogVisible = ref(false)
const exportDialogVisible = ref(false)
const importedCardItem = ref(null)

const prepareForExport = () => {
  if (currentCore.value.activeCard == undefined) {
    return;
  }
  dialogSource.value = cardItem.value.serialize()
}

const showExportDialog = () => {
  prepareForExport()
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
  const res = await activeCard.value.importCardItem(cardItem);
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


</script>

<template>

  <ElRow :gutter="24" class="mb-10px mt-10px"  v-if="activeCard.selectedItem !== -1">
    <ElCol :span="12" :xs="12">
      <ElButton class="w-[100%]" @click="addCardItem()">
        {{ t('dashboard.editor.addNewCardItem') }}
      </ElButton>
    </ElCol>
    <ElCol :span="12" :xs="12">
      <ElButton class="w-[100%]" @click="importDialogVisible = true">
        {{ t('main.import') }}
      </ElButton>
    </ElCol>
  </ElRow>

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
          <ElSelect
            v-model="cardItem.type"
            :placeholder="$t('dashboard.editor.pleaseSelectType')"
            style="width: 100%"
          >
            <ElOption
              v-for="item in itemTypes"
              :key="item.value"
              :label="$t('dashboard.editor.'+item.label)"
              :value="item.value"
            />

          </ElSelect>
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

    <ElButton type="primary" @click.prevent.stop='showExportDialog()' plain>
      <Icon icon="uil:file-export" class="mr-5px"/>
      {{ $t('main.export') }}
    </ElButton>
    <ElButton type="primary" @click.prevent.stop="updateCardItem" plain>{{
        $t('main.update')
      }}
    </ElButton>

    <ElButton @click.prevent.stop="duplicate">{{ $t('main.duplicate') }}</ElButton>

    <!--    <ElButton @click.prevent.stop="copyToClipboard">{{ $t('main.copyToClipboard') }}</ElButton>-->

    <!--    <ElPopconfirm-->
    <!--      :confirm-button-text="$t('main.ok')"-->
    <!--      :cancel-button-text="$t('main.no')"-->
    <!--      width="250"-->
    <!--      :title="$t('main.are_you_sure_to_do_want_this?')"-->
    <!--      @confirm="cancel"-->
    <!--    >-->
    <!--      <template #reference>-->
    <!--        <ElButton plain>{{ t('main.cancel') }}</ElButton>-->
    <!--      </template>-->
    <!--    </ElPopconfirm>-->

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
  <Dialog v-model="exportDialogVisible" :title="t('entities.dialogExportTitle')" :maxHeight="400" width="80%">
    <JsonViewer v-model="dialogSource"/>
  </Dialog>
  <!-- /export dialog -->

  <!-- import dialog -->
  <Dialog v-model="importDialogVisible" :title="t('entities.dialogImportTitle')" :maxHeight="400" width="80%"
          custom-class>
    <JsonEditor @change="importHandler"/>
    <template #footer>
      <ElButton type="primary" @click="importCardItem()" plain>{{ t('main.import') }}</ElButton>
      <ElButton @click="importDialogVisible = false">{{ t('main.closeDialog') }}</ElButton>
    </template>
  </Dialog>
  <!-- /import dialog -->


  <DraggableContainer :name="'editor-card-items'" :initial-width="280" :min-width="280" v-show="showMenuWindow">
    <template #header>
      <div class="w-[100%]">
        <div style="float: left">Card Items</div>
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

      <!--      <ElRow class="mb-10px mt-10px">-->
      <!--        <ElCol>-->
      <!--          <ElDivider content-position="left">{{ $t('dashboard.editor.itemList') }}</ElDivider>-->
      <!--        </ElCol>-->
      <!--      </ElRow>-->

      <ElRow class="mb-10px mt-10px">
        <ElCol>
          <ElButton class="w-[100%]" @click="addCardItem()">
            {{ t('dashboard.editor.addNewCardItem') }}
          </ElButton>
        </ElCol>
      </ElRow>

      <ElRow class="mb-10px mt-10px">
        <ElCol>
          <ElButton class="w-[100%]" @click="importDialogVisible = true">
            {{ t('main.import') }}
          </ElButton>
        </ElCol>
      </ElRow>

      <ElMenu
        v-if="activeCard && activeCard.id"
        ref="tabMenu"
        :default-active="activeCard.selectedItem + ''"
        v-model="activeCard.selectedItem"
        class="el-menu-vertical-demo box-card">
        <ElMenuItem
          :index="index + ''"
          :key="index"
          v-for="(item, index) in activeCard.items"
          @click="menuCardItemClick(index)">
          <div class="w-[100%] menu-item">
                <span>
                  {{ item.title }}
                <ElTag type="info" size="small">
                  {{ item.type }}
                </ElTag>
                </span>
            <ElButtonGroup class="buttons">
              <ElButton @click.prevent.stop="sortCardItemUp(item, index)" text size="small">
                <Icon icon="teenyicons:up-solid"/>
              </ElButton>
              <ElButton @click.prevent.stop="sortCardItemDown(item, index)" text size="small">
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
