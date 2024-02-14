<script setup lang="ts">
import {computed, PropType} from "vue";
import {CardItem, Core} from "@/views/Dashboard/core/core";
import {
  ElButton,
  ElCard,
  ElCol,
  ElCollapse,
  ElCollapseItem,
  ElColorPicker,
  ElDivider,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElOption,
  ElPopconfirm,
  ElRow,
  ElSelect,
  ElSwitch,
  ElTag
} from 'element-plus'
import {CommonEditor} from "@/views/Dashboard/card_items/common";
import {useI18n} from "@/hooks/web/useI18n";
import {KeysSearch} from "@/views/Dashboard/components";
import {comparisonType} from "@/views/Dashboard/core/types";

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

const addProp = () => {
  // console.log('addProp')

  if (!currentItem.value?.payload.progress?.items) {
    currentItem.value.payload.progress.items = []
  }

  let counter = 0
  if (currentItem.value?.payload.progress!.items.length) {
    counter = currentItem.value?.payload.progress!.items.length
  }

  currentItem.value?.payload.progress!.items.push({
    key: '',
    value: '',
    comparison: comparisonType.EQ,
  })
}

const removeProp = (index: number) => {
  if (!currentItem.value?.payload.progress?.items) {
    return
  }

  currentItem.value?.payload.progress?.items.splice(index, 1)
}

const onChangePropValue = (val, index) => {
  currentItem.value.payload.progress.items[index].key = val;
}

const onChangeValue = (val) => {
  currentItem.value.payload.progress.value = val;
}

</script>

<template>
  <div>
    <CommonEditor :item="currentItem" :core="core"/>

    <ElRow class="mb-10px mt-10px">
      <ElCol>
        <ElDivider content-position="left">{{ $t('dashboard.editor.progressOptions') }}</ElDivider>
      </ElCol>
    </ElRow>

    <ElRow class="mb-10px">
      <ElCol>
        <ElButton class="w-[100%]" @click.prevent.stop="addProp()">
          <Icon icon="ep:plus" class="mr-5px"/>
          {{ $t('dashboard.editor.addProp') }}
        </ElButton>
      </ElCol>
    </ElRow>

    <!-- props -->
    <ElCollapse>
      <ElCollapseItem
          :name="index"
          :key="index"
          v-for="(prop, index) in currentItem.payload.progress.items"
      >

        <template #title>
          <ElTag size="small">{{ prop.key }}</ElTag>
          +
          <ElTag size="small">{{ prop.comparison }}</ElTag>
          +
          <ElTag size="small">{{ prop.value }}</ElTag>
        </template>

        <ElCard shadow="never" class="item-card-editor">

          <ElForm
              label-position="top"
              :model="prop"
              style="width: 100%"
              ref="cardItemForm"
          >

            <ElRow>
              <ElCol>
                <ElFormItem :label="$t('dashboard.editor.attrField')" prop="text">
                  <KeysSearch v-model="prop.key" :obj="currentItem.lastEvent"
                              @change="onChangePropValue($event, index)"/>
                </ElFormItem>

              </ElCol>
            </ElRow>

            <ElRow>
              <ElCol>
                <ElFormItem :label="$t('dashboard.editor.comparison')" prop="comparison">
                  <ElSelect
                      v-model="prop.comparison"
                      placeholder="please select type"
                      style="width: 100%"
                  >
                    <ElOption label="==" value="eq"/>
                    <ElOption label="<" value="lt"/>
                    <ElOption label="<=" value="le"/>
                    <ElOption label="!=" value="ne"/>
                    <ElOption label=">=" value="ge"/>
                    <ElOption label=">" value="gt"/>
                  </ElSelect>
                </ElFormItem>

              </ElCol>
            </ElRow>

            <ElRow>
              <ElCol>
                <ElFormItem :label="$t('dashboard.editor.value')" prop="value">
                  <ElInput
                      placeholder="Please input"
                      v-model="prop.value"/>
                </ElFormItem>
              </ElCol>
            </ElRow>

            <ElRow>
              <ElCol>
                <ElFormItem :label="$t('dashboard.editor.color')" prop="background">
                  <ElColorPicker show-alpha v-model="prop.color"/>
                </ElFormItem>
              </ElCol>
            </ElRow>

            <ElRow>
              <ElCol>
                <div style="text-align: right;">
                  <ElPopconfirm
                      :confirm-button-text="$t('main.ok')"
                      :cancel-button-text="$t('main.no')"
                      width="250"
                      style="margin-left: 10px;"
                      :title="$t('main.are_you_sure_to_do_want_this?')"
                      @confirm="removeProp"
                  >
                    <template #reference>
                      <ElButton type="danger" plain>
                        <Icon icon="ep:delete" class="mr-5px"/>
                        {{ t('main.remove') }}
                      </ElButton>
                    </template>
                  </ElPopconfirm>
                </div>
              </ElCol>
            </ElRow>

          </ElForm>

        </ElCard>

      </ElCollapseItem>
    </ElCollapse>
    <!-- /props -->


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

    <ElRow :gutter="24">
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('dashboard.editor.strokeWidth')" prop="strokeWidth">
          <ElInputNumber v-model="currentItem.payload.progress.strokeWidth"/>
        </ElFormItem>
      </ElCol>
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('dashboard.editor.width')" prop="width">
          <ElInputNumber v-model="currentItem.payload.progress.width"/>
        </ElFormItem>
      </ElCol>
    </ElRow>

    <ElRow :gutter="24">
      <ElCol :span="12" :xs="12">
        <ElFormItem :label="$t('dashboard.editor.showText')" prop="showText">
          <ElSwitch v-model="currentItem.payload.progress.showText"/>
        </ElFormItem>
      </ElCol>
      <ElCol :span="12" :xs="12" v-if="currentItem.payload.progress.showText">
        <ElFormItem :label="$t('dashboard.editor.textInside')" prop="textInside">
          <ElSwitch v-model="currentItem.payload.progress.textInside"/>
        </ElFormItem>
      </ElCol>
    </ElRow>

    <ElRow>
      <ElCol>
        <ElFormItem :label="$t('dashboard.editor.color')" prop="background">
          <ElColorPicker show-alpha v-model="currentItem.payload.progress.color"/>
        </ElFormItem>
      </ElCol>
    </ElRow>

    <ElRow>
      <ElCol>
        <ElFormItem :label="$t('dashboard.editor.attrField')" prop="value">
          <KeysSearch v-model="currentItem.payload.progress.value" :obj="currentItem.lastEvent"
                      @change="onChangeValue"/>
        </ElFormItem>
      </ElCol>
    </ElRow>


  </div>
</template>

<style lang="less">

</style>
