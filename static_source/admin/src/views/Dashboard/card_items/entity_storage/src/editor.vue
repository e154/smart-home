<script setup lang="ts">
import {computed, PropType, ref} from "vue";
import {CardItem, Core} from "@/views/Dashboard/core";
import {
  ElButton,
  ElCard,
  ElCol,
  ElCollapse,
  ElCollapseItem,
  ElDivider,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElRow,
  ElSwitch,
  ElTag,
  ElPopconfirm
} from 'element-plus'
import {CommonEditor} from "@/views/Dashboard/card_items/common";
import {useI18n} from "@/hooks/web/useI18n";
import {EntitiesSearch} from "@/components/EntitiesSearch";

const {t} = useI18n()

const entityIds = ref([])

// ---------------------------------
// common
// ---------------------------------

const props = defineProps({
  core: {
    type: Object as PropType<Nullable<Core>>,
    default: () => null
  },
  item: {
    type: Object as PropType<Nullable<CardItem>>,
    default: () => null
  },
})

const currentItem = computed(() => props.item as CardItem)

// ---------------------------------
// component methods
// ---------------------------------

const addColumn = () => {
  if (!currentItem.value.payload.entityStorage?.columns) {
    currentItem.value.payload.entityStorage.columns = [];
  }

  let counter = 0;
  if (currentItem.value.payload.entityStorage.columns.length) {
    counter = currentItem.value.payload.entityStorage.columns.length;
  }

  currentItem.value.payload.entityStorage.columns.push({
    name: 'column ' + counter,
    attribute: ''
  });
}

const removeColumn = (index: number) => {
  currentItem.value.payload.entityStorage.columns.splice(index, 1);
}

</script>

<template>
  <div>

    <CommonEditor :item="item" :core="core"/>

    <ElRow class="mb-10px mt-10px">
      <ElCol>
        <ElDivider content-position="left">Entity storage options</ElDivider>
      </ElCol>
    </ElRow>

    <ElRow>
      <ElCol>
        <ElFormItem :label="$t('dashboard.editor.entityStorage.entities')" prop="entityIds">
          <EntitiesSearch v-model="currentItem.payload.entityStorage.entityIds"/>
        </ElFormItem>
      </ElCol>
    </ElRow>

    <ElRow class="mb-10px">
      <ElCol>
        <ElFormItem :label="$t('dashboard.editor.entityStorage.showFilter')" prop="filter">
          <ElSwitch v-model="currentItem.payload.entityStorage.filter"/>
        </ElFormItem>
      </ElCol>
    </ElRow>

    <ElRow class="mb-10px">
      <ElCol>
        <ElFormItem :label="$t('dashboard.editor.entityStorage.showPopup')" prop="filter">
          <ElSwitch v-model="currentItem.payload.entityStorage.showPopup"/>
        </ElFormItem>
      </ElCol>
    </ElRow>

    <ElRow class="mb-10px mt-10px">
      <ElCol>
        <ElDivider content-position="left"> {{ $t('dashboard.editor.entityStorage.tableColumns') }}</ElDivider>
      </ElCol>
    </ElRow>

    <ElRow>
      <ElCol>
        <ElButton class="w-[100%]" @click.prevent.stop="addColumn()">
          <Icon icon="ep:plus" class="mr-5px"/>
          {{ $t('dashboard.editor.entityStorage.addNewColumn') }}
        </ElButton>
      </ElCol>
    </ElRow>

    <!-- props -->
    <ElCollapse v-if="currentItem.payload?.entityStorage?.columns?.length">
      <ElCollapseItem
        :name="index"
        :key="index"
        v-for="(column, index) in currentItem.payload.entityStorage.columns"
      >

        <template #title>
          <div>
            <ElTag size="small">{{ column.name }}</ElTag>
          </div>
        </template>

        <ElCard shadow="never" class="item-card-editor">

          <ElForm
            label-position="top"
            :model="column"
            style="width: 100%"
            ref="cardItemForm">

            <ElRow>
              <ElCol>
                <ElFormItem :label="$t('dashboard.editor.name')" prop="name">
                  <ElInput placeholder="Please input" v-model="column.name" clearable/>
                </ElFormItem>
              </ElCol>
            </ElRow>

            <ElRow>
              <ElCol>
                <ElFormItem :label="$t('dashboard.editor.width')" prop="width">
                  <ElInputNumber placeholder="Please input" v-model="column.width" :min="1" clearable/>
                </ElFormItem>
              </ElCol>
            </ElRow>

            <ElRow>
              <ElCol>
                <ElFormItem :label="$t('dashboard.editor.attribute')" prop="attribute">
                  <ElInput placeholder="Please input" v-model="column.attribute" clearable/>
                </ElFormItem>
              </ElCol>
            </ElRow>

            <ElRow>
              <ElCol>
                <ElFormItem :label="$t('dashboard.editor.columnFilter')" prop="filter">
                  <ElInput placeholder="Please input" v-model="column.filter" clearable/>
                </ElFormItem>
              </ElCol>
            </ElRow>

            <ElRow class="mt-10px">
              <ElCol>
                <div style="text-align: right;">
                  <ElPopconfirm
                    :confirm-button-text="$t('main.ok')"
                    :cancel-button-text="$t('main.no')"
                    width="250"
                    style="margin-left: 10px;"
                    :title="$t('main.are_you_sure_to_do_want_this?')"
                    @confirm="removeColumn(index)"
                  >
                    <template #reference>
                      <ElButton type="danger" plain>
                        <Icon icon="ep:delete" class="mr-5px"/>
                        {{ t('main.remove') }}
                      </ElButton>
                    </template>
                  </ElPopconfirm>
                </div>
              </ElCol>
            </ElRow>

          </ElForm>

        </ElCard>

      </ElCollapseItem>
    </ElCollapse>
    <!-- /props -->


  </div>
</template>

<style lang="less">

</style>
