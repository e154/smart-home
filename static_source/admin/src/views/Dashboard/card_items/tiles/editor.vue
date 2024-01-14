<script setup lang="ts">
import {computed, PropType} from "vue";
import {CardItem, CompareProp, Core, requestCurrentState, Tab} from "@/views/Dashboard/core";
import {
  ElButton, ElCard,
  ElCol,
  ElCollapse,
  ElCollapseItem,
  ElDivider, ElForm, ElFormItem,
  ElInput,
  ElInputNumber, ElOption, ElPopconfirm,
  ElRow, ElSelect, ElSwitch, ElTag
} from 'element-plus'
import CommonEditor from "@/views/Dashboard/card_items/common/editor.vue";
import {useI18n} from "@/hooks/web/useI18n";
import JsonViewer from "@/components/JsonViewer/JsonViewer.vue";
import ImageSearch from "@/views/Images/components/ImageSearch.vue";
import {ApiEntity, ApiImage} from "@/api/stub";
import {ItemPayloadTiles} from "@/views/Dashboard/card_items/tiles/types";
import {prepareUrl} from "@/utils/serverId";
import EntitySearch from "@/views/Entities/components/EntitySearch.vue";
import TilePreview from "@/views/Dashboard/card_items/tiles/tilePreview.vue";

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
  currentItem.value.payload.tiles = {
    items: [],
    defaultImage: undefined,
    columns: 5,
    rows: 5,
    tileHeight: 25,
    tileWidth: 25,
    attribute: '',
  } as ItemPayloadTiles;
}

const updateCurrentState = () => {
  if (currentItem.value.entityId) {
    requestCurrentState(currentItem.value?.entityId)
  }
}

const reloadKeyDefaultImage = 0;

const onSelectDefaultImage = (image: ApiImage) => {
  if (!props.item.payload.tiles) {
    initDefaultValue();
  }

  currentItem.value.payload.tiles.image = image as ApiImage || undefined;
  props.item.update();
}

const addProp = () => {
  if (!props.item.payload.tiles) {
    initDefaultValue();
  }

  let counter = 0;
  if (props.item.payload.tiles.items.length) {
    counter = props.item.payload.tiles.items.length;
  }

  currentItem.value.payload.tiles.items.push({
    key: 'new proper' + counter,
    image: undefined,
    position: false,
    top: 0,
    left: 0,
    height: props.item.payload.tiles.tileHeight || 0,
    width: props.item.payload.tiles.tileWidth || 0,
  } as CompareProp);
  props.item.update();
}

const removeProp = (index: number) => {
  if (!props.item.payload.tiles) {
    initDefaultValue();
  }

  props.item.payload.tiles.items!.splice(index, 1);
  props.item.update();
}

const onSelectImageForState = (index: number, image: ApiImage) => {
  if (!props.item.payload.tiles) {
    initDefaultValue();
  }

  currentItem.value.payload.tiles.items[index].image = image as ApiImage || undefined;
  props.item.update();
}

const getUrl = (image: ApiImage): string => {
  if (!image || !image?.url) {
    return '';
  }
  return prepareUrl(import.meta.env.VITE_API_BASEPATH as string + image?.url);
}

const tilePreview = ({top, left, position, image}) => {
  let style = {
    width: `${currentItem.value.payload.tiles?.tileWidth}px`,
    height: `${currentItem.value.payload.tiles?.tileHeight}px`,
    background: `url(${getUrl(image)})`,
  }
  if (position) {
    style.background = `url(${getUrl(image)}) ${left}px ${top}px no-repeat`
  }
  return style
}

const getActionList = (entity?: ApiEntity) => {
  if (!entity) {
    return [];
  }
  return entity.actions;
}

const changedForActionButton = async (entity: ApiEntity) => {
  if (entity?.id) {
    currentItem.value.payload.tiles.entity = await currentCore.value.fetchEntity(entity.id);
    currentItem.value.payload.tiles.entityId = entity.id;
  } else {
    currentItem.value.payload.tiles.entity = undefined;
    currentItem.value.payload.tiles.entityId = '';
    currentItem.value.payload.tiles.actionName = '';
  }
}

</script>

<template>

  <CommonEditor :item="item" :core="core"/>

  <ElDivider content-position="left">{{ $t('dashboard.editor.tilesOptions') }}</ElDivider>

  <ElRow :gutter="24">
    <ElCol :span="12" :xs="12">
      <ElFormItem :label="$t('dashboard.editor.tiles.columns')" prop="columns">
        <ElInputNumber v-model="currentItem.payload.tiles.columns" :min="1" :value-on-clear="5"/>
      </ElFormItem>
    </ElCol>
    <ElCol :span="12" :xs="12">
      <ElFormItem :label="$t('dashboard.editor.tiles.rows')" prop="rows">
        <ElInputNumber v-model="currentItem.payload.tiles.rows" :min="1" :value-on-clear="5"/>
      </ElFormItem>
    </ElCol>
  </ElRow>

  <ElRow :gutter="24">
    <ElCol :span="12" :xs="12">
      <ElFormItem :label="$t('dashboard.editor.tiles.height')" prop="tileHeight">
        <ElInputNumber v-model="currentItem.payload.tiles.tileHeight" :min="1" :value-on-clear="25"/>
      </ElFormItem>
    </ElCol>
    <ElCol :span="12" :xs="12">
      <ElFormItem :label="$t('dashboard.editor.tiles.width')" prop="tileWidth">
        <ElInputNumber v-model="currentItem.payload.tiles.tileWidth" :min="1" :value-on-clear="25"/>
      </ElFormItem>
    </ElCol>
  </ElRow>

  <ElRow>
    <ElCol>
      <ElFormItem :label="$t('dashboard.editor.attrField')" prop="attribute">
        <ElInput placeholder="Please input" v-model="currentItem.payload.tiles.attribute"/>
      </ElFormItem>
    </ElCol>
  </ElRow>

  <ElDivider content-position="left">{{ $t('dashboard.editor.tiles.items') }}</ElDivider>

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
  <ElCollapse v-if="item.payload.tiles?.items?.length">
    <ElCollapseItem
        :name="index"
        :key="index"
        v-for="(prop, index) in item.payload.tiles.items"
    >

      <template #title>
        <div style="width: 100%;height: inherit;clear: both;">
          <ElTag size="small">{{ prop.key }}</ElTag>
          <TilePreview :base-params="currentItem.payload.tiles" :tile-item="prop"/>
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
              <ElFormItem :label="$t('dashboard.editor.tiles.height')" prop="tileHeight">
                <ElInputNumber v-model="prop.height" :min="0"/>
              </ElFormItem>
            </ElCol>
            <ElCol :span="12" :xs="12">
              <ElFormItem :label="$t('dashboard.editor.tiles.width')" prop="tileWidth">
                <ElInputNumber v-model="prop.width" :min="0"/>
              </ElFormItem>
            </ElCol>
          </ElRow>

          <ElRow :gutter="24">
            <ElCol :span="8" :xs="8">
              <ElFormItem :label="$t('dashboard.editor.tiles.position')" prop="position">
                <ElSwitch v-model="prop.position"/>
              </ElFormItem>
            </ElCol>
          </ElRow>


          <ElRow :gutter="24" v-if="prop.position">
            <ElCol :span="12" :xs="12">
              <ElFormItem :label="$t('dashboard.editor.tiles.top')" prop="top">
                <ElInputNumber v-model="prop.top" :step="1"/>
              </ElFormItem>
            </ElCol>
            <ElCol :span="12" :xs="12">
              <ElFormItem :label="$t('dashboard.editor.tiles.left')" prop="left">
                <ElInputNumber v-model="prop.left" :step="1"/>
              </ElFormItem>
            </ElCol>
          </ElRow>



          <ElDivider v-if="prop.position" content-position="left">{{ $t('dashboard.editor.tiles.preview') }}</ElDivider>

          <ElRow v-if="prop.position">
            <ElCol>
              <TilePreview :base-params="currentItem.payload.tiles" :tile-item="prop"/>
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
        <ImageSearch :key="reloadKeyDefaultImage" v-model="currentItem.payload.tiles.image" @change="onSelectDefaultImage"/>
      </ElFormItem>
    </ElCol>
  </ElRow>

  <ElRow :gutter="24">
    <ElCol :span="8" :xs="8">
      <ElFormItem :label="$t('dashboard.editor.tiles.position')" prop="position">
        <ElSwitch v-model="currentItem.payload.tiles.position"/>
      </ElFormItem>
    </ElCol>
  </ElRow>

  <ElRow :gutter="24" v-if="currentItem.payload.tiles.position">
    <ElCol :span="12" :xs="12">
      <ElFormItem :label="$t('dashboard.editor.tiles.top')" prop="top">
        <ElInputNumber v-model="currentItem.payload.tiles.top" :step="1"/>
      </ElFormItem>
    </ElCol>
    <ElCol :span="12" :xs="12">
      <ElFormItem :label="$t('dashboard.editor.tiles.left')" prop="left">
        <ElInputNumber v-model="currentItem.payload.tiles.left" :step="1"/>
      </ElFormItem>
    </ElCol>
  </ElRow>

  <ElDivider v-if="currentItem.payload.tiles.position" content-position="left">{{ $t('dashboard.editor.tiles.preview') }}</ElDivider>

  <ElRow v-if="currentItem.payload.tiles.position">
    <ElCol>
      <TilePreview :base-params="currentItem.payload.tiles"/>
    </ElCol>

  </ElRow>

  <ElDivider content-position="left">{{ $t('dashboard.editor.action') }}</ElDivider>

  <ElRow :gutter="24">
    <ElCol :span="12" :xs="12">
      <ElFormItem :label="$t('dashboard.editor.entity')" prop="entity">
        <EntitySearch v-model="currentItem.payload.tiles.entity" @change="changedForActionButton($event)"/>
      </ElFormItem>
    </ElCol>

    <ElCol :span="12" :xs="12">
      <ElFormItem :label="$t('dashboard.editor.action')"  prop="action" :aria-disabled="!currentItem.payload.tiles.entity">
        <ElSelect
            v-model="currentItem.payload.tiles.actionName"
            clearable
            :placeholder="$t('dashboard.editor.selectAction')"
            style="width: 100%"
        >
          <ElOption
              v-for="item in getActionList(currentItem.payload.tiles.entity)"
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

<style lang="less" >
.el-collapse-item__header {
  clear: both;
  overflow: hidden;
  .tile-preview-wrapper {
    float: right;
    margin-right: 15px;
  }
}

</style>
