<script setup lang="ts">
import {computed, PropType} from "vue";
import {CardItem, Core, requestCurrentState} from "@/views/Dashboard/core";
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
  ElOption,
  ElPopconfirm,
  ElRow,
  ElSelect,
  ElSwitch,
  ElTag
} from 'element-plus'
import {CommonEditor} from "@/views/Dashboard/card_items/common";
import {useI18n} from "@/hooks/web/useI18n";
import {JsonViewer} from "@/components/JsonViewer";
import ImageSearch from "@/views/Images/components/ImageSearch.vue";
import {ApiEntity, ApiImage} from "@/api/stub";
import {GridProp, ItemPayloadGrid} from "./types";
import {prepareUrl} from "@/utils/serverId";
import EntitySearch from "@/views/Entities/components/EntitySearch.vue";
import CellPreview from "./cellPreview.vue";
import {KeysSearch} from "@/views/Dashboard/components";

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

const updateCurrentState = () => {
  if (currentItem.value.entityId) {
    requestCurrentState(currentItem.value?.entityId)
  }
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

const getUrl = (image: ApiImage): string => {
  if (!image || !image?.url) {
    return '';
  }
  return prepareUrl(import.meta.env.VITE_API_BASEPATH as string + image?.url);
}

const getActionList = (entity?: ApiEntity) => {
  if (!entity) {
    return [];
  }
  return entity.actions;
}

const changedForActionButton = async (entity: ApiEntity) => {
  if (entity?.id) {
    currentItem.value.payload.grid.entity = await currentCore.value.fetchEntity(entity.id);
    currentItem.value.payload.grid.entityId = entity.id;
  } else {
    currentItem.value.payload.grid.entity = undefined;
    currentItem.value.payload.grid.entityId = '';
    currentItem.value.payload.grid.actionName = '';
  }
}

const onChangeValue = (val) => {
  currentItem.value.payload.grid.attribute = val;
}

</script>

<template>

  <CommonEditor :item="item" :core="core"/>

  <ElDivider content-position="left">{{ $t('dashboard.editor.gridOptions') }}</ElDivider>

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

  <ElDivider content-position="left">{{ $t('dashboard.editor.grid.items') }}</ElDivider>

  <ElRow>
    <ElCol>
      <div style="padding-bottom: 20px">
        <ElButton type="default" @click.prevent.stop="addProp()">
          <Icon icon="ep:plus" class="mr-5px"/>
          {{ $t('dashboard.editor.addNewProp') }}
        </ElButton>
      </div>
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

          <ElRow :gutter="24">
            <ElCol :span="8" :xs="8">
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

          <ElRow>
            <ElCol>
              <div style="padding-bottom: 20px">
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
                      <ElButton class="mr-10px" type="danger" plain>
                        <Icon icon="ep:delete" class="mr-5px"/>
                        {{ t('main.remove') }}
                      </ElButton>
                    </template>
                  </ElPopconfirm>
                </div>
              </div>
            </ElCol>
          </ElRow>

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

  <ElRow :gutter="24">
    <ElCol :span="8" :xs="8">
      <ElFormItem :label="$t('dashboard.editor.grid.position')" prop="position">
        <ElSwitch v-model="currentItem.payload.grid.position"/>
      </ElFormItem>
    </ElCol>
  </ElRow>

  <ElRow :gutter="24" v-if="currentItem.payload.grid.position">
    <ElCol :span="12" :xs="12">
      <ElFormItem :label="$t('dashboard.editor.grid.top')" prop="top">
        <ElInputNumber v-model="currentItem.payload.grid.top" :step="1"/>
      </ElFormItem>
    </ElCol>
    <ElCol :span="12" :xs="12">
      <ElFormItem :label="$t('dashboard.editor.grid.left')" prop="left">
        <ElInputNumber v-model="currentItem.payload.grid.left" :step="1"/>
      </ElFormItem>
    </ElCol>
  </ElRow>

  <ElDivider v-if="currentItem.payload.grid.position" content-position="left">{{
      $t('dashboard.editor.grid.preview')
    }}
  </ElDivider>

  <ElRow v-if="currentItem.payload.grid.position">
    <ElCol>
      <CellPreview :base-params="currentItem.payload.grid"/>
    </ElCol>

  </ElRow>

  <ElDivider content-position="left">{{ $t('dashboard.editor.action') }}</ElDivider>

  <ElRow :gutter="24">
    <ElCol :span="12" :xs="12">
      <ElFormItem :label="$t('dashboard.editor.entity')" prop="entity">
        <EntitySearch v-model="currentItem.payload.grid.entity" @change="changedForActionButton($event)"/>
      </ElFormItem>
    </ElCol>

    <ElCol :span="12" :xs="12">
      <ElFormItem :label="$t('dashboard.editor.action')" prop="action"
                  :aria-disabled="!currentItem.payload.grid.entity">
        <ElSelect
            v-model="currentItem.payload.grid.actionName"
            clearable
            :placeholder="$t('dashboard.editor.selectAction')"
            style="width: 100%"
        >
          <ElOption
              v-for="item in getActionList(currentItem.payload.grid.entity)"
              :key="item.name"
              :label="item.name"
              :value="item.name"/>
        </ElSelect>
      </ElFormItem>
    </ElCol>
  </ElRow>

  <ElRow style="padding-bottom: 20px" v-if="currentItem.entity">
    <ElCol>
      <ElCollapse>
        <ElCollapseItem :title="$t('dashboard.editor.eventstateJSONobject')">
          <ElButton type="default" @click.prevent.stop="updateCurrentState()" style="margin-bottom: 20px">
            <Icon icon="ep:refresh" class="mr-5px"/>
            {{ $t('dashboard.editor.getEvent') }}
          </ElButton>

          <JsonViewer v-model="currentItem.lastEvent"/>

        </ElCollapseItem>
      </ElCollapse>
    </ElCol>
  </ElRow>

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
