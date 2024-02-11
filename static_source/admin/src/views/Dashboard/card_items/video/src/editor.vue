<script setup lang="ts">
import {computed, PropType} from "vue";
import {CardItem, Core, requestCurrentState} from "@/views/Dashboard/core/core";
import {
  ElButton,
  ElCol,
  ElCollapse,
  ElCollapseItem,
  ElDivider,
  ElFormItem,
  ElOption,
  ElRow,
  ElSelect
} from 'element-plus'
import {JsonViewer} from "@/components/JsonViewer";
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

const updateCurrentState = () => {
  if (currentItem.value.entityId) {
    requestCurrentState(currentItem.value.entityId)
  }
}

const onChangeValue = (val) => {
  currentItem.value.payload.video.attribute = val;
}

</script>

<template>
  <div>

    <CommonEditor :item="item" :core="core"/>

    <ElDivider content-position="left">{{ $t('dashboard.editor.video.options') }}</ElDivider>

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

    <ElRow class="mb-10px" v-if="currentItem.entity">
      <ElCol>
        <ElCollapse>
          <ElCollapseItem :title="$t('dashboard.editor.eventstateJSONobject')">
            <ElButton class="mb-10px w-[100%]" type="default" @click.prevent.stop="updateCurrentState()">
              <Icon icon="ep:refresh" class="mr-5px"/>
              {{ $t('dashboard.editor.getEvent') }}
            </ElButton>

            <JsonViewer v-model="currentItem.lastEvent"/>

          </ElCollapseItem>
        </ElCollapse>
      </ElCol>
    </ElRow>

  </div>
</template>

<style lang="less">

</style>
