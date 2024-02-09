<script setup lang="ts">
import {computed, PropType} from "vue";
import {CardItem, Core, requestCurrentState} from "@/views/Dashboard/core/core";
import {
  ElButton,
  ElCol,
  ElCollapse,
  ElCollapseItem,
  ElColorPicker,
  ElDivider,
  ElFormItem,
  ElOption,
  ElRow,
  ElSelect
} from 'element-plus'
import {JsonViewer} from "@/components/JsonViewer";
import {CommonEditor} from "@/views/Dashboard/card_items/common";
import {useI18n} from "@/hooks/web/useI18n";
import {KeysSearch} from "@/views/Dashboard/components";
import {EntitiesAction, EntitiesActionOptions} from "@/components/EntitiesAction";

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

const onChangeValue = (val) => {
  currentItem.value.payload.colorPicker.attribute = val;
}

const changedForActionButton = async (options: EntitiesActionOptions) => {
  currentItem.value.payload.colorPicker.entityId = options.entityId
  currentItem.value.payload.colorPicker.action = options.action
  currentItem.value.payload.colorPicker.tags = options.tags
  currentItem.value.payload.colorPicker.areaId = options.areaId
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

    <ElRow :gutter="24">
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('dashboard.editor.value')" prop="value">
          <KeysSearch v-model="currentItem.payload.colorPicker.attribute" :obj="currentItem.lastEvent"
                      @change="onChangeValue"/>
        </ElFormItem>
      </ElCol>
    </ElRow>

    <ElDivider content-position="left">{{ $t('dashboard.editor.actionOptions') }}</ElDivider>

    <EntitiesAction :options="currentItem.payload.colorPicker" :entity="currentItem.entity" @change="changedForActionButton($event)"/>

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

<style lang="less">

</style>
