<script setup lang="ts">
import {computed, PropType} from "vue";
import {CardItem, Core} from "@/views/Dashboard/core/core";
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
  ElPopconfirm,
  ElRow,
  ElSwitch,
  ElTag
} from 'element-plus'
import {CommonEditor} from "@/views/Dashboard/card_items/common";
import {useI18n} from "@/hooks/web/useI18n";
import {ImageSearch} from "@/components/ImageSearch";
import {ApiImage} from "@/api/stub";
import {GridProp, ItemPayloadGrid} from "./types";
import CellPreview from "./CellPreview.vue";
import {KeysSearch} from "@/views/Dashboard/components";
import {EntitiesAction, EntitiesActionOptions} from "@/components/EntitiesAction";

const {t} = useI18n()

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

const currentCore = computed(() => props.core as Core)

const currentItem = computed(() => props.item as CardItem)

// ---------------------------------
// component methods
// ---------------------------------

const initDefaultValue = () => {
  currentItem.value.payload.grid = {
    items: [],
    defaultImage: undefined,
    cellHeight: 25,
    cellWidth: 25,
    attribute: '',
    gap: false,
    gapSize: 5,
    tooltip: false,
    fontSize: 18,
  } as ItemPayloadGrid;
}

const onSelectDefaultImage = (image: ApiImage) => {
  if (!props.item.payload.grid) {
    initDefaultValue();
  }

  currentItem.value.payload.grid.image = image as ApiImage || undefined;
  props.item.update();
}

const addProp = () => {
  if (!props.item.payload.grid) {
    initDefaultValue();
  }

  let counter = 0;
  if (props.item.payload.grid.items.length) {
    counter = props.item.payload.grid.items.length;
  }

  currentItem.value.payload.grid.items.push({
    key: 'new proper' + counter,
    image: undefined,
    position: false,
    tooltip: false,
    top: 0,
    left: 0,
    height: props.item.payload.grid.cellHeight || 0,
    width: props.item.payload.grid.cellWidth || 0,
  } as GridProp);
  props.item.update();
}

const removeProp = (index: number) => {
  if (!props.item.payload.grid) {
    initDefaultValue();
  }

  props.item.payload.grid.items!.splice(index, 1);
  props.item.update();
}

const onSelectImageForState = (index: number, image: ApiImage) => {
  if (!props.item.payload.grid) {
    initDefaultValue();
  }

  currentItem.value.payload.grid.items[index].image = image as ApiImage || undefined;
  props.item.update();
}

const changedForActionButton = async (options: EntitiesActionOptions) => {
  currentItem.value.payload.grid.entityId = options.entityId
  currentItem.value.payload.grid.actionName = options.action
  currentItem.value.payload.grid.tags = options.tags
  currentItem.value.payload.grid.areaId = options.areaId
}

const onChangeValue = (val) => {
  currentItem.value.payload.grid.attribute = val;
}

</script>

<template>

  <CommonEditor :item="item" :core="core"/>

  <ElRow class="mb-10px mt-10px">
    <ElCol>
      <ElDivider content-position="left">{{ $t('dashboard.editor.gridOptions') }}</ElDivider>
    </ElCol>
  </ElRow>

  <ElRow :gutter="24">
    <ElCol :span="12" :xs="12">
      <ElFormItem :label="$t('dashboard.editor.grid.height')" prop="cellHeight">
        <ElInputNumber v-model="currentItem.payload.grid.cellHeight" :min="1" :value-on-clear="25"/>
      </ElFormItem>
    </ElCol>
    <ElCol :span="12" :xs="12">
      <ElFormItem :label="$t('dashboard.editor.grid.width')" prop="cellWidth">
        <ElInputNumber v-model="currentItem.payload.grid.cellWidth" :min="1" :value-on-clear="25"/>
      </ElFormItem>
    </ElCol>
  </ElRow>

  <ElRow :gutter="24">
    <ElCol :span="12" :xs="12">
      <ElFormItem :label="$t('dashboard.editor.grid.showCellValue')" prop="showCellValue">
        <ElSwitch v-model="currentItem.payload.grid.showCellValue"/>
      </ElFormItem>
    </ElCol>
    <ElCol :span="12" :xs="12" v-if="currentItem.payload.grid.showCellValue">
      <ElFormItem :label="$t('dashboard.editor.grid.fontSize')" prop="fontSize">
        <ElInputNumber v-model="currentItem.payload.grid.fontSize" :min="1" :value-on-clear="12"/>
      </ElFormItem>
    </ElCol>
  </ElRow>

  <ElRow :gutter="24">
    <ElCol :span="12" :xs="12">
      <ElFormItem :label="$t('dashboard.editor.grid.gap')" prop="gap">
        <ElSwitch v-model="currentItem.payload.grid.gap"/>
      </ElFormItem>
    </ElCol>
    <ElCol :span="12" :xs="12" v-if="currentItem.payload.grid.gap">
      <ElFormItem :label="$t('dashboard.editor.grid.gapSize')" prop="gapSize">
        <ElInputNumber v-model="currentItem.payload.grid.gapSize" :min="1" :value-on-clear="1"/>
      </ElFormItem>
    </ElCol>
  </ElRow>

  <ElRow :gutter="24">
    <ElCol :span="12" :xs="12">
      <ElFormItem :label="$t('dashboard.editor.grid.tooltip')" prop="tooltip">
        <ElSwitch v-model="currentItem.payload.grid.tooltip"/>
      </ElFormItem>
    </ElCol>
  </ElRow>

  <ElRow>
    <ElCol>
      <ElFormItem :label="$t('dashboard.editor.attrField')" prop="attribute">
        <KeysSearch v-model="currentItem.payload.grid.attribute" :obj="currentItem.lastEvent" @change="onChangeValue"/>
      </ElFormItem>
    </ElCol>
  </ElRow>

  <ElRow class="mb-10px">
    <ElCol>
      <ElDivider content-position="left">{{ $t('dashboard.editor.grid.items') }}</ElDivider>
    </ElCol>
  </ElRow>

  <ElRow class="mb-10px">
    <ElCol>
      <ElButton class="w-[100%]" @click.prevent.stop="addProp()">
        <Icon icon="ep:plus" class="mr-5px"/>
        {{ $t('dashboard.editor.addNewProp') }}
      </ElButton>
    </ElCol>
  </ElRow>

  <!-- props -->
  <ElCollapse v-if="item.payload.grid?.items?.length">
    <ElCollapseItem
        :name="index"
        :key="index"
        v-for="(prop, index) in item.payload.grid.items"
    >

      <template #title>
        <div style="width: 100%;height: inherit;clear: both;">
          <ElTag size="small">{{ prop.key }}</ElTag>
          <CellPreview :base-params="currentItem.payload.grid" :tile-item="prop"/>
        </div>
      </template>

      <ElCard shadow="never" class="item-card-editor">

        <ElForm
            label-position="top"
            :model="prop"
            style="width: 100%"
            ref="cardItemForm">

          <ElRow>
            <ElCol>
              <ElFormItem :label="$t('dashboard.editor.text')" prop="text">
                <ElInput placeholder="Please input" v-model="prop.key"/>
              </ElFormItem>
            </ElCol>
          </ElRow>

          <ElRow>
            <ElCol>
              <ElFormItem :label="$t('dashboard.editor.image')" prop="image">
                <ImageSearch v-model="prop.image" @change="onSelectImageForState(index, ...arguments)"/>
              </ElFormItem>
            </ElCol>
          </ElRow>

          <ElRow :gutter="24">
            <ElCol :span="12" :xs="12">
              <ElFormItem :label="$t('dashboard.editor.grid.height')" prop="cellHeight">
                <ElInputNumber v-model="prop.height" :min="0"/>
              </ElFormItem>
            </ElCol>
            <ElCol :span="12" :xs="12">
              <ElFormItem :label="$t('dashboard.editor.grid.width')" prop="cellWidth">
                <ElInputNumber v-model="prop.width" :min="0"/>
              </ElFormItem>
            </ElCol>
          </ElRow>

          <ElRow>
            <ElCol>
              <ElFormItem :label="$t('dashboard.editor.grid.position')" prop="position">
                <ElSwitch v-model="prop.position"/>
              </ElFormItem>
            </ElCol>
          </ElRow>


          <ElRow :gutter="24" v-if="prop.position">
            <ElCol :span="12" :xs="12">
              <ElFormItem :label="$t('dashboard.editor.grid.top')" prop="top">
                <ElInputNumber v-model="prop.top" :step="1"/>
              </ElFormItem>
            </ElCol>
            <ElCol :span="12" :xs="12">
              <ElFormItem :label="$t('dashboard.editor.grid.left')" prop="left">
                <ElInputNumber v-model="prop.left" :step="1"/>
              </ElFormItem>
            </ElCol>
          </ElRow>

          <ElDivider v-if="prop.position" content-position="left">{{ $t('dashboard.editor.grid.preview') }}</ElDivider>

          <ElRow v-if="prop.position">
            <ElCol>
              <CellPreview :base-params="currentItem.payload.grid" :tile-item="prop"/>
            </ElCol>
          </ElRow>

          <div style="text-align: right;">
            <ElPopconfirm
                :confirm-button-text="$t('main.ok')"
                :cancel-button-text="$t('main.no')"
                width="250"
                style="margin-left: 10px;"
                :title="$t('main.are_you_sure_to_do_want_this?')"
                @confirm="removeProp(index)"
            >
              <template #reference>
                <ElButton type="danger" plain>
                  <Icon icon="ep:delete" class="mr-5px"/>
                  {{ t('main.remove') }}
                </ElButton>
              </template>
            </ElPopconfirm>
          </div>


        </ElForm>

      </ElCard>

    </ElCollapseItem>
  </ElCollapse>
  <!-- /props -->

  <ElRow>
    <ElCol>
      <ElDivider content-position="left">{{ $t('dashboard.editor.image') }}</ElDivider>
    </ElCol>
  </ElRow>

  <ElRow>
    <ElCol>
      <ElFormItem :label="$t('dashboard.editor.image')" prop="image">
        <ImageSearch v-model="currentItem.payload.grid.image" @change="onSelectDefaultImage"/>
      </ElFormItem>
    </ElCol>
  </ElRow>

  <ElRow>
    <ElCol>
      <ElFormItem :label="$t('dashboard.editor.grid.position')" prop="position">
        <ElSwitch v-model="currentItem.payload.grid.position"/>
      </ElFormItem>
    </ElCol>
  </ElRow>

  <ElRow v-if="currentItem.payload.grid.position">
    <ElCol>
      <ElFormItem :label="$t('dashboard.editor.grid.top')" prop="top">
        <ElInputNumber v-model="currentItem.payload.grid.top" :step="1"/>
      </ElFormItem>
    </ElCol>
  </ElRow>

  <ElRow v-if="currentItem.payload.grid.position">
    <ElCol>
      <ElFormItem :label="$t('dashboard.editor.grid.left')" prop="left">
        <ElInputNumber v-model="currentItem.payload.grid.left" :step="1"/>
      </ElFormItem>
    </ElCol>
  </ElRow>

  <ElDivider v-if="currentItem.payload.grid.position" content-position="left">{{
      $t('dashboard.editor.grid.preview')
    }}
  </ElDivider>

  <ElRow v-if="currentItem.payload.grid.position" class="mb-20px">
    <ElCol>
      <CellPreview :base-params="currentItem.payload.grid"/>
    </ElCol>
  </ElRow>

  <ElRow>
    <ElCol>
      <ElFormItem :label="$t('dashboard.editor.grid.tileClick')" prop="tileClick">
        <ElSwitch v-model="currentItem.payload.grid.tileClick"/>
      </ElFormItem>
    </ElCol>
  </ElRow>

  <ElDivider v-if="currentItem.payload.grid?.tileClick" content-position="left">{{
      $t('dashboard.editor.actionOptions')
    }}
  </ElDivider>

  <EntitiesAction v-if="currentItem.payload.grid?.tileClick" :options="currentItem.payload.grid"
                  :entity="currentItem.entity" @change="changedForActionButton($event)"/>

</template>

<style lang="less">
.el-collapse-item__header {
  clear: both;
  overflow: hidden;

  .tile-preview-wrapper {
    float: right;
    margin-right: 15px;
  }
}

</style>
