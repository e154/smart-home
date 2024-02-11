<script setup lang="ts">
import {computed, PropType} from "vue";
import {CardItem, Core} from "@/views/Dashboard/core/core";
import {ElCol, ElDivider, ElFormItem, ElOption, ElRow, ElSelect} from 'element-plus'
import {CommonEditor} from "@/views/Dashboard/card_items/common";
import {playerType} from "./types";
import {KeysSearch} from "@/views/Dashboard/components";

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

const onChangeValue = (val) => {
  currentItem.value.payload.video.attribute = val;
}

</script>

<template>
  <div>

    <CommonEditor :item="item" :core="core"/>

    <ElRow class="mb-10px mt-10px">
      <ElCol>
        <ElDivider content-position="left">{{ $t('dashboard.editor.video.options') }}</ElDivider>
      </ElCol>
    </ElRow>

    <ElRow :gutter="24">
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('dashboard.editor.type')" prop="type">
          <ElSelect
              v-model="currentItem.payload.video.playerType"
              :placeholder="$t('dashboard.editor.video.pleaseSelectPlayerType')"
              style="width: 100%"
          >
            <ElOption label="ONVIF MSE" value="onvifMse"/>
            <ElOption label="YOUTUBE" value="youtube"/>
          </ElSelect>
        </ElFormItem>
      </ElCol>
      <ElCol :span="12" :xs="12"/>
    </ElRow>

    <ElRow :gutter="24" v-if="currentItem.payload.video.playerType === playerType.youtube">
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('dashboard.editor.value')" prop="value">
          <KeysSearch v-model="currentItem.payload.video.attribute" :obj="currentItem.lastEvent"
                      @change="onChangeValue"/>
        </ElFormItem>
      </ElCol>
      <ElCol :span="12" :xs="12"/>
    </ElRow>

  </div>
</template>

<style lang="less">

</style>
