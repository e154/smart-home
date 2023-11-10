<script setup lang="ts">
import {computed, PropType, ref, unref, watch} from "vue";
import {Card, CardItem, Core, Tab} from "@/views/Dashboard/core";
import {ElDivider, ElCollapse, ElCollapseItem, ElCard, ElForm, ElFormItem, ElInputNumber, ElSwitch,
  ElRow, ElCol, ElSelect, ElOption, ElInput, ElTag, ElButton, ElPopconfirm, ElColorPicker } from 'element-plus'
import {Vuuri} from "@/views/Dashboard/Vuuri"
import {useBus} from "@/views/Dashboard/bus";
import ViewCard from "@/views/Dashboard/editor/ViewCard.vue";
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


</script>

<template>
  <div>

    <CommonEditor :item="currentItem" :core="core"/>

    <ElDivider content-position="left">{{$t('dashboard.editor.buttonOptions')}}</ElDivider>

    <ElRow :gutter="24">
      <ElCol :span="8" :xs="8">
        <ElFormItem :label="$t('dashboard.editor.icon')" prop="icon">
          <ElInput size="small" v-model="currentItem.payload.button.icon"/>
        </ElFormItem>

        <ElFormItem :label="$t('dashboard.editor.text')" prop="text">
          <ElInput size="small" v-model="currentItem.payload.button.text"/>
        </ElFormItem>

        <ElFormItem :label="$t('dashboard.editor.action')" prop="action" :aria-disabled="!currentItem.entity">

          <ElSelect
              v-model="currentItem.payload.button.action"
              clearable
              :placeholder="$t('dashboard.editor.selectAction')"
              style="width: 100%"
          >
            <ElOption
                v-for="p in currentItem.entityActions"
                :key="p.value"
                :label="p.label + ' (' +p.value +')'"
                :value="p.value"/>
          </ElSelect>

        </ElFormItem>

      </ElCol>
      <ElCol :span="8" :xs="8">

        <ElFormItem :label="$t('dashboard.editor.type')" prop="type">
          <ElSelect
              v-model="currentItem.payload.button.type"
              placeholder="please select type"
              style="width: 100%"
          >
            <ElOption label="Default" value=""/>
            <ElOption label="Primary" value="primary"/>
            <ElOption label="Success" value="success"/>
            <ElOption label="Info" value="info"/>
            <ElOption label="Warning" value="warning"/>
            <ElOption label="Danger" value="danger"/>
            <ElOption label="Text" value="text"/>
          </ElSelect>
        </ElFormItem>

        <ElFormItem :label="$t('dashboard.editor.size')" prop="size">
          <ElSelect
              v-model="currentItem.payload.button.size"
              placeholder="please select type"
              style="width: 100%"
          >
            <ElOption label="Small" value="small"/>
            <ElOption label="Large" value="large"/>
            <ElOption label="Default" value="default"/>
          </ElSelect>
        </ElFormItem>

      </ElCol>
      <ElCol :span="8" :xs="8">

        <ElFormItem :label="$t('dashboard.editor.round')" prop="round">
          <ElSwitch v-model="currentItem.payload.button.round"/>
        </ElFormItem>
      </ElCol>
    </ElRow>

  </div>
</template>

<style lang="less" >

</style>
