<script setup lang="ts">
import {computed, PropType} from "vue";
import {CardItem, Core} from "@/views/Dashboard/core/core";
import {
  ElCol,
  ElDivider,
  ElFormItem,
  ElInputNumber,
  ElOption,
  ElRow,
  ElSelect,
  ElSwitch
} from 'element-plus'
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
  currentItem.value.payload.slider.attribute = val;
}

const changedForActionButton = async (options: EntitiesActionOptions) => {
  currentItem.value.payload.slider.entityId = options.entityId
  currentItem.value.payload.slider.action = options.action
  currentItem.value.payload.slider.tags = options.tags
  currentItem.value.payload.slider.areaId = options.areaId
}

</script>

<template>
  <div>
    <CommonEditor :item="currentItem" :core="core"/>

    <ElRow class="mb-10px mt-10px">
      <ElCol>
        <ElDivider content-position="left">{{ $t('dashboard.editor.slider.options') }}</ElDivider>
      </ElCol>
    </ElRow>

    <ElRow :gutter="24">
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('dashboard.editor.type')" prop="type">
          <ElSelect
              v-model="currentItem.payload.slider.orientation"
              :placeholder="$t('dashboard.editor.pleaseSelectOrientation')"
              style="width: 100%"
          >
            <ElOption label="Horizontal" value="horizontal"/>
            <ElOption label="Vertical" value="vertical"/>
            <ElOption label="Circular" value="circular"/>
          </ElSelect>
        </ElFormItem>
      </ElCol>
      <ElCol :span="12" :xs="12"/>
    </ElRow>

    <ElRow :gutter="24">
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('dashboard.editor.slider.color')" prop="color">
          <ColorPicker show-alpha v-model="currentItem.payload.slider.color"/>
        </ElFormItem>
      </ElCol>
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('dashboard.editor.slider.trackColor')" prop="trackColor">
          <ColorPicker show-alpha v-model="currentItem.payload.slider.trackColor"/>
        </ElFormItem>
      </ElCol>
    </ElRow>

    <ElRow :gutter="24">
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('dashboard.editor.slider.height')" prop="height">
          <ElInputNumber v-model="currentItem.payload.slider.height" :min="1"/>
        </ElFormItem>
      </ElCol>
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('dashboard.editor.slider.step')" prop="step">
          <ElInputNumber v-model="currentItem.payload.slider.step" :min="0"/>
        </ElFormItem>
      </ElCol>
    </ElRow>


    <ElRow :gutter="24">
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('dashboard.editor.slider.min')" prop="min">
          <ElInputNumber v-model="currentItem.payload.slider.min" :min="0"/>
        </ElFormItem>
      </ElCol>
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('dashboard.editor.slider.max')" prop="max">
          <ElInputNumber v-model="currentItem.payload.slider.max" :min="0"/>
        </ElFormItem>
      </ElCol>
    </ElRow>

    <ElRow :gutter="24">
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('dashboard.editor.tooltip')" prop="round">
          <ElSwitch v-model="currentItem.payload.slider.tooltip"/>
        </ElFormItem>
      </ElCol>
    </ElRow>

    <ElRow :gutter="24">
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('dashboard.editor.attrField')" prop="value">
          <KeysSearch v-model="currentItem.payload.slider.attribute" :obj="currentItem.lastEvent"
                      @change="onChangeValue"/>
        </ElFormItem>
      </ElCol>
    </ElRow>

    <ElDivider content-position="left">{{ $t('dashboard.editor.actionOptions') }}</ElDivider>

    <EntitiesAction :options="currentItem.payload.slider" :entity="currentItem.entity"
                    @change="changedForActionButton($event)"/>

  </div>
</template>

<style lang="less">

</style>
