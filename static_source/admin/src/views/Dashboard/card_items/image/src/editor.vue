<script setup lang="ts">
import {computed, PropType} from "vue";
import {CardItem, Core, requestCurrentState} from "@/views/Dashboard/core/core";
import {ElButton, ElCol, ElCollapse, ElCollapseItem, ElDivider, ElFormItem, ElRow, ElSwitch} from 'element-plus'
import {CommonEditor} from "@/views/Dashboard/card_items/common";
import {Cache} from "@/views/Dashboard/core/render";
import {JsonViewer} from "@/components/JsonViewer";
import {ApiImage} from "@/api/stub";
import {ImageSearch} from "@/components/ImageSearch";
import {KeysSearch} from "@/views/Dashboard/components";

// ---------------------------------
// common
// ---------------------------------

const _cache: Cache = new Cache();

const props = defineProps({
  core: {
    type: Object as PropType<Core>,
  },
  item: {
    type: Object as PropType<Nullable<CardItem>>,
    default: () => null
  },
})

const currentItem = computed({
  get(): CardItem {
    return props.item as CardItem
  },
  set(val: CardItem) {
  }
})

// ---------------------------------
// component methods
// ---------------------------------

const onSelectImage = (index: number, image: ApiImage) => {
  if (!props.item?.payload?.image) {
    return;
  }
  // console.log('select image', index, image);
  currentItem.value.payload.image.image = image || undefined;
}

const updateCurrentState = () => {
  if (currentItem.value.entityId) {
    requestCurrentState(currentItem.value?.entityId)
  }
}

const onChangePropValue = (val: string) => {
  currentItem.value.payload.image.attrField = val;
}

</script>

<template>

  <CommonEditor :item="currentItem" :core="core"/>

  <ElDivider content-position="left">{{ $t('dashboard.editor.imageOptions') }}</ElDivider>

  <ElRow :gutter="24">
    <ElCol :span="12" :xs="12">
      <ElFormItem :label="$t('dashboard.editor.background')" prop="background">
        <ElSwitch v-model="currentItem.payload.image.background"/>
      </ElFormItem>
    </ElCol>
  </ElRow>

  <ElFormItem :label="$t('dashboard.editor.image')" prop="image">
    <ImageSearch v-model="currentItem.payload.image.image" @change="onSelectImage"/>
  </ElFormItem>

  <ElFormItem :label="$t('dashboard.editor.attrField')" prop="text">
    <KeysSearch v-model="currentItem.payload.image.attrField" :obj="currentItem.lastEvent" @change="onChangePropValue"/>
  </ElFormItem>

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

</style>
