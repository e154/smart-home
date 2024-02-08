<script setup lang="ts">
import {computed, PropType} from "vue";
import {CardItem, Core} from "@/views/Dashboard/core/core";
import {ElCol, ElDivider, ElFormItem, ElOption, ElRow, ElSelect} from 'element-plus'
import {CommonEditor} from "@/views/Dashboard/card_items/common";
import {ApiImage} from "@/api/stub";
import {ImageSearch} from "@/components/ImageSearch";
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

const onSelectImage = (index: number, image: ApiImage) => {
  if (!props.item?.payload?.joystick) {
    return;
  }
  // console.log('select image', index, image);
  currentItem.value.payload.joystick.stickImage = image || undefined;
}


</script>

<template>
  <div>

    <CommonEditor :item="item" :core="core"/>

    <ElDivider content-position="left">{{ $t('dashboard.editor.joystick.options') }}</ElDivider>

    <ElFormItem :label="$t('dashboard.editor.image')" prop="image">
      <ImageSearch v-model="currentItem.payload.joystick.stickImage" @change="onSelectImage(index, ...arguments)"/>
    </ElFormItem>

    <ElRow :gutter="24" v-if="currentItem.entity">

      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('dashboard.editor.joystick.startAction')" prop="startAction"
                    :aria-disabled="!currentItem.entity">
          <ElSelect
              v-model="currentItem.payload.joystick.startAction"
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
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('dashboard.editor.joystick.endAction')" prop="endAction"
                    :aria-disabled="!currentItem.entity">
          <ElSelect
              v-model="currentItem.payload.joystick.endAction"
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

    </ElRow>

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

  </div>
</template>

<style lang="less">

</style>
