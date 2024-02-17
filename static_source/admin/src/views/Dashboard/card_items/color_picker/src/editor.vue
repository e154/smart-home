<script setup lang="ts">
import {computed, PropType} from "vue";
import {CardItem, Core} from "@/views/Dashboard/core/core";
import {ElCol, ElDivider, ElFormItem, ElRow} from 'element-plus'
import {CommonEditor} from "@/views/Dashboard/card_items/common";
import {useI18n} from "@/hooks/web/useI18n";
import {KeysSearch} from "@/views/Dashboard/components";
import {EntitiesAction, EntitiesActionOptions} from "@/components/EntitiesAction";
import {ColorPicker} from "@/components/ColorPicker";

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

    <ElRow class="mb-10px mt-10px">
      <ElCol>
        <ElDivider content-position="left">{{ $t('dashboard.editor.colorPicker.options') }}</ElDivider>
      </ElCol>
    </ElRow>

    <ElRow>
      <ElCol>
        <ElFormItem :label="$t('dashboard.editor.colorPicker.defaultColor')" prop="color">
          <ColorPicker show-alpha v-model="currentItem.payload.colorPicker.color"/>
        </ElFormItem>
      </ElCol>
    </ElRow>

    <ElRow>
      <ElCol>
        <ElFormItem :label="$t('dashboard.editor.attrField')" prop="value">
          <KeysSearch v-model="currentItem.payload.colorPicker.attribute" :obj="currentItem.lastEvent"
                      @change="onChangeValue"/>
        </ElFormItem>
      </ElCol>
    </ElRow>

    <ElDivider content-position="left">{{ $t('dashboard.editor.actionOptions') }}</ElDivider>

    <EntitiesAction :options="currentItem.payload.colorPicker" :entity="currentItem.entity"
                    @change="changedForActionButton($event)"/>

  </div>
</template>

<style lang="less">

</style>
