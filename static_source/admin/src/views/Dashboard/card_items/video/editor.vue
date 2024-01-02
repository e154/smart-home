<script setup lang="ts">
import {computed, PropType, ref, unref, watch} from "vue";
import {Card, CardItem, Core, requestCurrentState, Tab} from "@/views/Dashboard/core";
import {ElDivider, ElCollapse, ElCollapseItem, ElCard, ElForm, ElFormItem, ElInputNumber, ElSwitch,
  ElRow, ElCol, ElSelect, ElOption, ElInput, ElTag, ElButton, ElColorPicker } from 'element-plus'
import JsonViewer from "@/components/JsonViewer/JsonViewer.vue";
import CommonEditor from "@/views/Dashboard/card_items/common/editor.vue";
import {playerType} from "@/views/Dashboard/card_items/video/types";

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
          <ElInput v-model="currentItem.payload.video.attribute"/>
        </ElFormItem>
      </ElCol>
      <ElCol :span="12" :xs="12"/>
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

  </div>
</template>

<style lang="less" >

</style>
