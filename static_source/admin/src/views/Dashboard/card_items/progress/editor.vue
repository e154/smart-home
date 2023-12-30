<script setup lang="ts">
import {computed, PropType, ref, unref, watch} from "vue";
import {Card, CardItem, Core, requestCurrentState, Tab} from "@/views/Dashboard/core";
import {ElDivider, ElCollapse, ElCollapseItem, ElCard, ElForm, ElFormItem, ElInputNumber, ElSwitch,
  ElRow, ElCol, ElSelect, ElOption, ElInput, ElTag, ElButton } from 'element-plus'
import JsonViewer from "@/components/JsonViewer/JsonViewer.vue";
import CommonEditor from "@/views/Dashboard/card_items/common/editor.vue";
import {useI18n} from "@/hooks/web/useI18n";

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
    requestCurrentState(currentItem.value?.entityId)
  }
}
</script>

<template>
  <div>
    <CommonEditor :item="currentItem" :core="core"/>

    <ElDivider content-position="left">{{ $t('dashboard.editor.progressOptions') }}</ElDivider>

    <ElFormItem :label="$t('dashboard.editor.type')" prop="type">
      <ElSelect
          v-model="currentItem.payload.progress.type"
          placeholder="please select type"
          style="width: 100%"
      >
        <ElOption label="linear" value=""/>
        <ElOption label="circle" value="circle"/>
        <ElOption label="dashboard" value="dashboard"/>
      </ElSelect>
    </ElFormItem>

    <ElFormItem :label="$t('dashboard.editor.textInside')" prop="textInside">
      <ElSwitch v-model="currentItem.payload.progress.textInside"/>
    </ElFormItem>

    <ElFormItem :label="$t('dashboard.editor.strokeWidth')" prop="strokeWidth">
      <ElInputNumber v-model="currentItem.payload.progress.strokeWidth"/>
    </ElFormItem>

    <ElFormItem :label="$t('dashboard.editor.width')" prop="width">
      <ElInputNumber v-model="currentItem.payload.progress.width"/>
    </ElFormItem>

    <ElFormItem :label="$t('dashboard.editor.value')" prop="value">
      <ElInput v-model="currentItem.payload.progress.value"/>
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

  </div>
</template>

<style lang="less" >

</style>
