<script setup lang="ts">
import {computed, PropType, ref, unref, watch} from "vue";
import {Card, CardItem, Core, requestCurrentState, Tab} from "@/views/Dashboard/core";
import {ElDivider, ElCollapse, ElCollapseItem, ElCard, ElForm, ElFormItem, ElInputNumber, ElSwitch,
  ElRow, ElCol, ElSelect, ElOption, ElInput, ElTag, ElButton, ElColorPicker } from 'element-plus'
import {Vuuri} from "@/views/Dashboard/Vuuri"
import {useBus} from "@/views/Dashboard/bus";
import Viewer from "@/components/JsonViewer/JsonViewer.vue";
import CommonEditor from "@/views/Dashboard/card_items/common/editor.vue";
import {useI18n} from "@/hooks/web/useI18n";

const {t} = useI18n()

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

    <ElDivider content-position="left">{{ $t('dashboard.editor.colorPicker.options') }}</ElDivider>

    <ElRow :gutter="24">
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('dashboard.editor.colorPicker.defaultColor')" prop="color">
          <ElColorPicker show-alpha v-model="currentItem.payload.colorPicker.color"/>
        </ElFormItem>
      </ElCol>
    </ElRow>

    <ElRow :gutter="24" >
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('dashboard.editor.value')" prop="value">
          <ElInput size="small" v-model="currentItem.payload.colorPicker.attribute"/>
        </ElFormItem>
      </ElCol>
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('dashboard.editor.action')" prop="action" :aria-disabled="!currentItem.entity">

          <ElSelect
              v-model="currentItem.payload.colorPicker.action"
              clearable
              :placeholder="$t('dashboard.editor.selectAction')"
              style="width: 100%"
          >
            <ElOption
                v-for="p in currentItem.entityActions"
                :key="p.value"
                :label="p.label"
                :value="p.value"/>
          </ElSelect>
        </ElFormItem>
      </ElCol>
    </ElRow>

    <ElRow style="padding-bottom: 20px" v-if="currentItem.entity">
      <ElCol>
        <ElCollapse>
          <ElCollapseItem :title="$t('dashboard.editor.eventstateJSONobject')">
            <ElButton type="default" @click.prevent.stop="updateCurrentState()" style="margin-bottom: 20px">
              <Icon icon="ep:refresh" class="mr-5px"/>
              {{ $t('dashboard.editor.getEvent') }}
            </ElButton>

            <Viewer v-model="currentItem.lastEvent"/>

          </ElCollapseItem>
        </ElCollapse>
      </ElCol>
    </ElRow>

  </div>
</template>

<style lang="less" >

</style>
