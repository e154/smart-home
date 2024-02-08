<script setup lang="ts">
import {computed, PropType} from "vue";
import {CardItem, Core, requestCurrentState} from "@/views/Dashboard/core";
import {
  ElButton,
  ElCol,
  ElCollapse,
  ElCollapseItem,
  ElColorPicker,
  ElDivider,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElRow
} from 'element-plus'
import {CommonEditor} from "@/views/Dashboard/card_items/common";
import {useI18n} from "@/hooks/web/useI18n";
import {Cache} from "@/views/Dashboard/render";
import {JsonViewer} from "@/components/JsonViewer";
import {KeysSearch} from "@/views/Dashboard/components";

const {t} = useI18n()

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

const updateCurrentState = () => {
  if (currentItem.value.entityId) {
    requestCurrentState(currentItem.value?.entityId)
  }
}

const onChangeValue = (val) => {
  currentItem.value.payload.icon.attrField = val;
}

</script>

<template>

  <CommonEditor :item="currentItem" :core="core"/>

  <ElDivider content-position="left">{{ $t('dashboard.editor.iconOptions') }}</ElDivider>

  <ElRow :gutter="24">
    <ElCol :span="8" :xs="8">
      <ElFormItem :label="$t('dashboard.editor.icon')" prop="icon">
        <ElInput v-model="currentItem.payload.icon.value"/>
      </ElFormItem>
    </ElCol>
    <ElCol :span="8" :xs="8">
      <ElFormItem :label="$t('dashboard.editor.iconColor')" prop="iconColor">
        <ElColorPicker show-alpha v-model="currentItem.payload.icon.iconColor"/>
      </ElFormItem>
    </ElCol>
    <ElCol :span="8" :xs="8">
      <ElFormItem :label="$t('dashboard.editor.iconSize')" prop="iconSize">
        <ElInputNumber v-model="currentItem.payload.icon.iconSize" :min="1" :value-on-clear="12"/>
      </ElFormItem>
    </ElCol>
  </ElRow>

  <ElFormItem :label="$t('dashboard.editor.attrField')" prop="text">
    <KeysSearch v-model="currentItem.payload.icon.attrField" :obj="currentItem.lastEvent" @change="onChangeValue"/>
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
