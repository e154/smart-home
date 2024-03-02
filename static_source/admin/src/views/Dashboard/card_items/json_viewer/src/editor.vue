<script setup lang="ts">
import {computed, onMounted, PropType} from "vue";
import {CardItem, Core, Cache} from "@/views/Dashboard/core";
import {ElCol, ElDivider, ElFormItem, ElRow,} from 'element-plus'
import {CommonEditor} from "@/views/Dashboard/card_items/common";
import {KeysSearch} from "@/views/Dashboard/components";
import {ItemPayloadJsonViewer} from "@/views/Dashboard/card_items/json_viewer";

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

const initDefaultValue = () => {
  currentItem.value.payload.jsonViewer = {
    attrField: undefined,
  } as ItemPayloadJsonViewer;
}

onMounted(() => {
  if (!currentItem.value.payload?.jsonViewer) {
    initDefaultValue()
  }
})


const onChangeValue = (val) => {
  if (!currentItem.value.payload.jsonViewer) {
    initDefaultValue()
  }
  currentItem.value.payload.jsonViewer.attrField = val;
}

</script>

<template>

  <CommonEditor :item="currentItem" :core="core"/>

  <div v-if="currentItem.payload.jsonViewer">
    <ElRow class="mb-10px mt-10px">
      <ElCol>
        <ElDivider content-position="left">{{ $t('dashboard.editor.jsonViewerOptions') }}</ElDivider>
      </ElCol>
    </ElRow>

    <ElFormItem :label="$t('dashboard.editor.attrField')" prop="text">
      <KeysSearch :all-keys="true" v-model="currentItem.payload.jsonViewer.attrField" :obj="currentItem.lastEvent"
                  @change="onChangeValue"/>
    </ElFormItem>
  </div>

</template>

<style lang="less">

</style>
