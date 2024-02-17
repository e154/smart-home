<script setup lang="ts">
import {computed, PropType} from "vue";
import {CardItem, Core} from "@/views/Dashboard/core/core";
import {ElCol, ElDivider, ElFormItem, ElRow} from 'element-plus'
import {CommonEditor} from "@/views/Dashboard/card_items/common";
import {ApiImage} from "@/api/stub";
import {ImageSearch} from "@/components/ImageSearch";
import {useI18n} from "@/hooks/web/useI18n";
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

const onSelectImage = (image: ApiImage) => {
  if (!props.item?.payload?.joystick) {
    return;
  }
  // console.log('select image', index, image);
  currentItem.value.payload.joystick.stickImage = image || undefined;
}

const changedForStartAction = async (options: EntitiesActionOptions) => {
  currentItem.value.payload.joystick.startAction = options
}


const changedForEndAction = async (options: EntitiesActionOptions) => {
  currentItem.value.payload.joystick.endAction = options
}


</script>

<template>

  <CommonEditor :item="item" :core="core"/>

  <ElRow class="mb-10px mt-10px">
    <ElCol>
      <ElDivider content-position="left">{{ $t('dashboard.editor.joystick.options') }}</ElDivider>
    </ElCol>
  </ElRow>

  <ElFormItem :label="$t('dashboard.editor.image')" prop="image">
    <ImageSearch v-model="currentItem.payload.joystick.stickImage" @change="onSelectImage"/>
  </ElFormItem>

  <ElRow class="mb-10px mt-10px">
    <ElCol>
      <ElDivider content-position="left">{{ $t('dashboard.editor.joystick.startAction') }}</ElDivider>
    </ElCol>
  </ElRow>

  <EntitiesAction :options="currentItem.payload.joystick.startAction" :entity="currentItem.entity"
                  @change="changedForStartAction($event)"/>

  <ElRow class="mb-10px mt-10px">
    <ElCol>
      <ElDivider content-position="left">{{ $t('dashboard.editor.joystick.endAction') }}</ElDivider>
    </ElCol>
  </ElRow>

  <EntitiesAction :options="currentItem.payload.joystick.endAction" :entity="currentItem.entity"
                  @change="changedForEndAction($event)"/>

  <!--    <ElRow :gutter="24">-->

  <!--      <ElCol :span="12" :xs="12">-->
  <!--        <ElFormItem :label="$t('dashboard.editor.startTimeout')" prop="startTimeout">-->
  <!--          <ElInputNumber v-model="currentItem.payload.joystick.startTimeout" :min="500"/>-->
  <!--        </ElFormItem>-->
  <!--      </ElCol>-->
  <!--      <ElCol :span="12" :xs="12">-->
  <!--        <ElFormItem :label="$t('dashboard.editor.endTimeout')" prop="endTimeout">-->
  <!--          <ElInputNumber v-model="currentItem.payload.joystick.endTimeout" :min="0"/>-->
  <!--        </ElFormItem>-->
  <!--      </ElCol>-->

  <!--    </ElRow>-->

</template>

<style lang="less">

</style>
