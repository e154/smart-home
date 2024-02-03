<script setup lang="ts">
import {computed, onMounted, PropType, ref, unref, watch} from "vue";
import {Card, CardItem, comparisonType, Core, requestCurrentState, Tab} from "@/views/Dashboard/core";
import {ElDivider, ElCollapse, ElCollapseItem, ElCard, ElForm, ElFormItem, ElPopconfirm, ElSwitch,
  ElRow, ElCol, ElSelect, ElOption, ElInput, ElTag, ElButton } from 'element-plus'
import CommonEditor from "@/views/Dashboard/card_items/common/editor.vue";
import {Cache, GetTokens} from "@/views/Dashboard/render";
import JsonViewer from "@/components/JsonViewer/JsonViewer.vue";
import {ApiImage} from "@/api/stub";
import ImageSearch from "@/views/Images/components/ImageSearch.vue";
import KeysSearch from "@/views/Dashboard/components/KeysSearch.vue";

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
  set(val: CardItem) {}
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
  currentItem.value.payload.image.attrField= val;
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

<style lang="less" >

</style>
