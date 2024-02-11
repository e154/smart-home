<script setup lang="ts">
import {computed, PropType} from "vue";
import {CardItem, Core} from "@/views/Dashboard/core/core";
import {ElCol, ElDivider, ElFormItem, ElRow, ElSwitch} from 'element-plus'
import {CommonEditor} from "@/views/Dashboard/card_items/common";
import {Cache} from "@/views/Dashboard/core/render";
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

const onChangePropValue = (val: string) => {
  currentItem.value.payload.image.attrField = val;
}

</script>

<template>

  <CommonEditor :item="currentItem" :core="core"/>

  <ElRow class="mb-10px mt-10px">
    <ElCol>
      <ElDivider content-position="left">{{ $t('dashboard.editor.imageOptions') }}</ElDivider>
    </ElCol>
  </ElRow>

  <ElRow>
    <ElCol>
      <ElFormItem :label="$t('dashboard.editor.background')" prop="background">
        <ElSwitch v-model="currentItem.payload.image.background"/>
      </ElFormItem>
    </ElCol>
  </ElRow>

  <ElFormItem :label="$t('dashboard.editor.image')" prop="image">
    <ImageSearch v-model="currentItem.payload.image.image" @change="onSelectImage"/>
  </ElFormItem>

  <ElRow class="mb-10px mt-10px">
    <ElCol>
      <ElDivider content-position="left">{{ $t('main.or') }}</ElDivider>
    </ElCol>
  </ElRow>

  <ElFormItem :label="$t('dashboard.editor.attrField')" prop="text">
    <KeysSearch v-model="currentItem.payload.image.attrField" :obj="currentItem.lastEvent" @change="onChangePropValue"/>
  </ElFormItem>

</template>

<style lang="less">

</style>
