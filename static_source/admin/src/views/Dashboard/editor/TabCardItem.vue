<script setup lang="ts">
import {computed, nextTick, PropType, reactive, ref, unref, watch} from 'vue'
import {ElButton, ElCard, ElCol, ElMessage, ElPopconfirm,
  ElRow, ElSkeleton, ElMenu, ElMenuItem, ElButtonGroup, ElForm, ElFormItem, ElRadioGroup, ElCheckboxGroup,
  ElRadio, ElCheckbox, ElSwitch, ElTimePicker, ElDatePicker, ElSelect, ElOption, ElInput, ElTag} from 'element-plus'
import {useI18n} from '@/hooks/web/useI18n'
import {useValidator} from '@/hooks/web/useValidator'
import {Card, CardItem, Core, Tab} from "@/views/Dashboard/core";
import {useRouter} from "vue-router";
import {useBus} from "@/views/Dashboard/bus";
import {CardEditorName, CardItemList} from "@/views/Dashboard/card_items";

const {required} = useValidator()
const {t} = useI18n()

const {currentRoute, addRoute, push} = useRouter()
const {bus} = useBus()

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
  set(val: Card) {}
})

const currentCore = computed({
  get(): Core {
    return props.core as Core
  },
  set(val: Core) {}
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

const copyCardItem = () => {
  activeCard.value.copyItem(activeCard.value.selectedItem);
}

const menuCardItemClick = (index: number) => {
  if (currentCore.value.activeTab < 0 || currentCore.value.activeCard == undefined) {
    return;
  }

    activeCard.value.selectedItem = index;

  bus.emit('selected_card_item', index)
}

const sortCardItemUp = (item: CardItem, index: number) => {
  // console.log('up', item, index)

  if (!activeCard.value.items[index - 1]) {
    return;
  }

  let rows = [activeCard.value.items[index - 1], activeCard.value.items[index]];
  activeCard.value.items.splice(index - 1, 2, rows[1], rows[0]);

  let counter = 0
  for (const index in activeCard.value.items) {
    activeCard.value.items[index].weight = counter;
    counter++;
  }

  currentCore.value.updateCard();
}

const sortCardItemDown = (item: CardItem, index: number) => {
  // console.log('down', item, index)

  if (!activeCard.value.items[index + 1]) {
    return;
  }

  let rows = [activeCard.value.items[index], activeCard.value.items[index + 1]];
  activeCard.value.items.splice(index, 2, rows[1], rows[0]);

  let counter = 0
  for (const index in activeCard.value.items) {
    activeCard.value.items[index].weight = counter;
    counter++;
  }

  currentCore.value.updateCard();
}

const getCardEditorName = (name: string) => {
  return CardEditorName(name);
}

const cancel = () => {
 console.warn('action not implemented')
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

</script>

<template>

  <ElRow :gutter="20">
    <ElCol :span="15" :xs="15">
      <ElCard class="box-card">
        <template #header>
          <div class="card-header">
            <span>{{ $t('dashboard.editor.itemDetail') }}</span>
          </div>
        </template>

        <ElForm
            v-if="cardItem"
            :model="cardItem"
            label-position="top"
            style="width: 100%"
            ref="cardItemForm"
        >

          <ElRow :gutter="24">
            <ElCol :span="12" :xs="12">
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
            <ElCol :span="12" :xs="12">
              <ElFormItem :label="$t('dashboard.editor.title')" prop="title">
                <ElInput size="small" v-model="cardItem.title"/>
              </ElFormItem>
            </ElCol>

          </ElRow>

          <component
              :is="getCardEditorName(cardItem.type)"
              :core="core"
              :item="cardItem"
          />
      </ElForm>

        <ElSkeleton v-if="!activeCard.items.length || activeCard.selectedItem === -1" :rows="5" class="mt-20px mb-20px" />

        <div class="text-right" v-if="activeCard.selectedItem > -1">
          <ElButton type="primary" @click.prevent.stop="updateCardItem">{{ $t('main.update') }}</ElButton>
          <ElButton type="default" @click.prevent.stop="copyCardItem">{{ $t('main.copy') }}</ElButton>
          <ElPopconfirm
              :confirm-button-text="$t('main.ok')"
              :cancel-button-text="$t('main.no')"
              width="250"
              style="margin-left: 10px;"
              :title="$t('main.are_you_sure_to_do_want_this?')"
              @confirm="cancel"
          >
            <template #reference>
              <ElButton type="default" plain>{{ t('main.cancel') }}</ElButton>
            </template>
          </ElPopconfirm>
          <ElPopconfirm
              :confirm-button-text="$t('main.ok')"
              :cancel-button-text="$t('main.no')"
              width="250"
              style="margin-left: 10px;"
              :title="$t('main.are_you_sure_to_do_want_this?')"
              @confirm="removeCardItem(activeCard.selectedItem)"
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
    </ElCol>

    <ElCol :span="8" :xs="12">
      <ElCard class="box-card">
        <template #header>
          <div class="item-header">
            <span>{{ $t('dashboard.editor.itemList') }}</span>
            <ElButtonGroup>
              <ElButton @click="addCardItem()" text>
                {{ t('dashboard.editor.addNewCardItem') }}
              </ElButton>
            </ElButtonGroup>
          </div>
        </template>
        <ElMenu
            v-if="activeCard && activeCard.id"
            ref="tabMenu"
            :default-active="activeCard.selectedItem + ''"
            v-model="activeCard.selectedItem"
            class="el-menu-vertical-demo">
          <ElMenuItem
              :index="index + ''"
              :key="index"
              v-for="(item, index) in activeCard.items"
              @click="menuCardItemClick(index)">
            <div class="w-[100%] item-header">
              <span>
                <strong>{{ item.title }}</strong>
              <ElTag type="info" class="mb-18px ml-10px">
                {{ item.type }}
              </ElTag>
              </span>
              <ElButtonGroup class="hide">
                <ElButton type="default" @click.prevent.stop="sortCardItemUp(item, index)">
                  <Icon icon="teenyicons:up-solid" />
                </ElButton>
                <ElButton type="default" @click.prevent.stop="sortCardItemDown(item, index)">
                  <Icon icon="teenyicons:down-solid" />
                </ElButton>
              </ElButtonGroup>
            </div>
          </ElMenuItem>
        </ElMenu>

      </ElCard>
    </ElCol>
  </ElRow>


</template>

<style lang="less" scoped>
.item-header {
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
