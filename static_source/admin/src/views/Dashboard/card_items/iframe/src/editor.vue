<script setup lang="ts">
import {computed, onMounted, PropType, ref} from "vue";
import {CardItem, Core} from "@/views/Dashboard/core";
import {ElCol, ElDivider, ElFormItem, ElInput, ElRow} from 'element-plus'
import {CommonEditor} from "@/views/Dashboard/card_items/common";
import {KeysSearch} from "@/views/Dashboard/components";
import {ItemPayloadIFrame} from "@/views/Dashboard/card_items/iframe";

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

const currentItem = computed({
  get(): CardItem {
    return props.item as CardItem
  },
  set(val: CardItem) {
  }
})

const isInit = ref(true)

onMounted(() => {
  if (!props.item?.payload?.iframe) {
    initDefaultValue();
  }
  isInit.value = false
})

// ---------------------------------
// component methods
// ---------------------------------

const initDefaultValue = () => {
  currentItem.value.payload.iframe = {
    uri: '',
    attrField: '',
  } as ItemPayloadIFrame;
}

const onChangeValue = (val) => {
  currentItem.value.payload.iframe.attrField = val;
}

</script>

<template>
  <div v-if="!isInit">

    <CommonEditor :item="item" :core="core"/>

    <ElRow class="mb-10px mt-10px">
      <ElCol>
        <ElDivider content-position="left">{{ $t('dashboard.editor.iframeOptions') }}</ElDivider>
      </ElCol>
    </ElRow>

    <ElRow class="mb-10px mt-10px">
      <ElCol>
        <ElFormItem :label="$t('dashboard.editor.uri')" prop="value">
          <ElInput placeholder="Please input" clearable v-model="currentItem.payload.iframe.uri"/>
        </ElFormItem>
      </ElCol>
    </ElRow>

    <ElRow class="mb-10px mt-10px">
      <ElCol>
        <ElDivider content-position="left">{{ $t('main.or') }}</ElDivider>
      </ElCol>
    </ElRow>

    <ElFormItem :label="$t('dashboard.editor.attrField')" prop="text">
      <KeysSearch v-model="currentItem.payload.iframe.attrField" :obj="currentItem.lastEvent" @change="onChangeValue"/>
    </ElFormItem>

  </div>
</template>

<style lang="less">

</style>
