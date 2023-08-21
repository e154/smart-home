<script setup lang="ts">
import {computed, PropType, ref, unref, watch} from "vue";
import {Card, CardItem, Core, requestCurrentState, Tab} from "@/views/Dashboard/core";
import {ElDivider, ElCollapse, ElCollapseItem, ElCard, ElForm, ElFormItem, ElInputNumber, ElSwitch,
  ElRow, ElCol, ElSelect, ElOption, ElInput, ElTag, ElButton, ElPopconfirm, ElColorPicker } from 'element-plus'
import CommonEditor from "@/views/Dashboard/card_items/common/editor.vue";
import {FilterList, RangeList} from "@/views/Dashboard/card_items/chart/types";
import Viewer from "@/components/JsonViewer/JsonViewer.vue";
import {useI18n} from "@/hooks/web/useI18n";

const {t} = useI18n()

// ---------------------------------
// common
// ---------------------------------
const rangeList = RangeList;
const filterList = FilterList;

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

const currentItem = computed({
  get(): CardItem {
    return props.item as CardItem
  },
  set(val: CardItem) {}
})

// ---------------------------------
// component methods
// ---------------------------------

const updateCurrentState = () => {
  if (currentItem.value.entityId) {
    requestCurrentState(currentItem.value?.entityId)
  }
}

const getMetricList = computed(() => {
  if (props.item.entity && props.item.entity?.metrics && props.item.payload.chart.metric_index !== undefined && props.item.entity.metrics[props.item.payload.chart.metric_index]) {
    return props.item.entity.metrics
  } else {
    return []
  }
})

const getMetricItem = computed(() => {
  if (props.item.entity && props.item.entity?.metrics && props.item.payload.chart.metric_index !== undefined && props.item.entity.metrics[props.item.payload.chart.metric_index]) {
    return props.item.entity.metrics[props.item.payload.chart.metric_index].options.items
  } else {
    return []
  }
})

const addChartItem = () => {
  if (!currentItem.value.payload.chart.items) {
    currentItem.value.payload.chart.items = []
  }
  currentItem.value.payload.chart.items.push({
    value: 'value'
  })
}

const removeChartItem = (index: number) => {
  if (!currentItem.value.payload.chart.items) {
    return;
  }

  currentItem.value.payload.chart.items.splice(index, 1);
}

</script>

<template>

  <CommonEditor :item="item" :core="core"/>

  <ElDivider content-position="left">{{ $t('dashboard.editor.chart.chartOptions') }}</ElDivider>

  <ElFormItem :label="$t('dashboard.editor.chart.type')" prop="type">
    <ElSelect
        v-model="currentItem.payload.chart.type"
        placeholder="please select type"
        style="width: 100%"
    >
      <ElOption label="Linear" value="line"/>
      <ElOption label="Bar" value="bar"/>
      <ElOption label="Pie" value="pie"/>
      <ElOption label="Doughnut" value="doughnut"/>
    </ElSelect>
  </ElFormItem>

  <ElRow :gutter="24">
    <ElCol :span="12" :xs="12">
      <ElFormItem :label="$t('dashboard.editor.chart.xAxis')" prop="xAxis">
        <ElSwitch v-model="currentItem.payload.chart.xAxis"/>
      </ElFormItem>

      <ElFormItem :label="$t('dashboard.editor.chart.yAxis')" prop="yAxis">
        <ElSwitch v-model="currentItem.payload.chart.yAxis"/>
      </ElFormItem>
    </ElCol>
    <ElCol :span="12" :xs="12">
      <ElFormItem :label="$t('dashboard.editor.chart.legend')" prop="legend">
        <ElSwitch v-model="currentItem.payload.chart.legend"/>
      </ElFormItem>

      <ElFormItem :label="$t('dashboard.editor.chart.scale')" prop="scale">
        <ElSwitch v-model="currentItem.payload.chart.scale"/>
      </ElFormItem>
    </ElCol>
  </ElRow>

  <ElRow :gutter="24">
    <ElCol :span="12" :xs="12">
      <ElFormItem :label="$t('dashboard.editor.chart.color')" prop="background">
        <ElColorPicker show-alpha v-model="currentItem.payload.chart.color"/>
      </ElFormItem>
    </ElCol>
    <ElCol :span="12" :xs="12">
      <ElFormItem :label="$t('dashboard.editor.chart.backgroundColor')" prop="background">
        <ElColorPicker show-alpha v-model="currentItem.payload.chart.backgroundColor"/>
      </ElFormItem>
    </ElCol>
  </ElRow>

  <ElDivider content-position="left" v-if="currentItem.payload.chart.type === 'line'">{{$t('dashboard.editor.chart.lineOptions') }}</ElDivider>

  <ElRow :gutter="24" v-if="currentItem.payload.chart.type === 'line'">
    <ElCol :span="12" :xs="12">

    <ElFormItem :label="$t('dashboard.editor.chart.entity_metric')" prop="index">
      <ElSelect v-model="currentItem.payload.chart.metric_index" placeholder="Select" class="w-[100%]">
        <ElOption
            v-for="(prop, index) in getMetricList"
            :key="index"
            :label="prop.name"
            :value="index"/>
      </ElSelect>
    </ElFormItem>

    <ElFormItem :label="$t('dashboard.editor.chart.metric_props')" prop="index">
      <ElSelect v-model="currentItem.payload.chart.props" multiple placeholder="Select" clearable class="w-[100%]">
        <ElOption
            v-for="p in getMetricItem"
            :key="p.name"
            :label="p.name"
            :value="p.name"/>
      </ElSelect>
    </ElFormItem>

    </ElCol>
    <ElCol :span="12" :xs="12">
      <ElFormItem :label="$t('dashboard.editor.chart.range')" prop="index">
        <ElSelect v-model="currentItem.payload.chart.range" placeholder="Select" clearable class="w-[100%]">
          <ElOption
              v-for="p in rangeList"
              :key="p.value"
              :label="p.label"
              :value="p.value"/>
        </ElSelect>
      </ElFormItem>

      <ElFormItem :label="$t('dashboard.editor.chart.filter')" prop="index">
        <ElSelect v-model="currentItem.payload.chart.filter" placeholder="Select" clearable class="w-[100%]">
          <ElOption
              v-for="p in filterList"
              :key="p.value"
              :label="p.label"
              :value="p.value"/>
        </ElSelect>
      </ElFormItem>
    </ElCol>
  </ElRow>

  <ElRow :gutter="24" v-if="currentItem.payload.chart.type === 'line'">
    <ElCol :span="12" :xs="12">
      <ElFormItem :label="$t('dashboard.editor.chart.borderWidth')" prop="borderWidth">
        <ElInputNumber size="small" v-model="currentItem.payload.chart.borderWidth" :min="1" :max="10"/>
      </ElFormItem>
    </ElCol>

  </ElRow>

  <!-- chart items -->
  <ElDivider content-position="left" v-if="['bar', 'pie', 'doughnut'].includes(currentItem.payload.chart.type)">{{$t('dashboard.editor.chart.barOptions') }}</ElDivider>

  <ElRow :gutter="24" v-if="['bar', 'pie', 'doughnut'].includes(currentItem.payload.chart.type)">
    <ElCol :span="12" :xs="12">
      <ElFormItem :label="$t('dashboard.editor.chart.automatic')" prop="borderWidth">
        <ElSwitch v-model="currentItem.payload.chart.automatic"/>
      </ElFormItem>
    </ElCol>


  </ElRow>

  <ElRow v-if="['bar', 'pie', 'doughnut'].includes(currentItem.payload.chart.type)">
    <ElCol>
      <div style="padding-bottom: 20px">
        <ElButton type="default" @click.prevent.stop="addChartItem()">
          <Icon icon="ep:plus" class="mr-5px"/>
          {{ $t('dashboard.editor.chart.addChartItem') }}
        </ElButton>
      </div>

      <!-- props -->
      <ElCollapse>
        <ElCollapseItem
            v-for="(prop, index) in currentItem.payload.chart.items"
            :name="index"
            :key="index"
        >

          <template #title>
            {{ prop.value }}
          </template>

          <ElCard shadow="never" class="item-card-editor">

            <ElForm
                label-position="top"
                :model="prop"
                style="width: 100%"
                ref="cardItemForm">

              <ElRow :gutter="24">
                <ElCol :span="12" :xs="12">
                  <ElFormItem :label="$t('dashboard.editor.chart.itemValue')" prop="text">
                    <ElInput class="w-[100%]" placeholder="Please input" v-model="prop.value"/>
                  </ElFormItem>
                </ElCol>
                <ElCol :span="12" :xs="12">
                  <ElFormItem :label="$t('dashboard.editor.chart.itemDescription')" prop="text">
                    <ElInput class="w-[100%]" placeholder="Please input" v-model="prop.description"/>
                  </ElFormItem>
                </ElCol>
              </ElRow>


              <div class="mb-20px">
                <div style="text-align: right;">
                  <ElPopconfirm
                      :confirm-button-text="$t('main.ok')"
                      :cancel-button-text="$t('main.no')"
                      width="250"
                      style="margin-left: 10px;"
                      :title="$t('main.are_you_sure_to_do_want_this?')"
                      @confirm="removeChartItem(index)"
                  >
                    <template #reference>
                      <ElButton class="mr-10px" type="danger" plain>
                        <Icon icon="ep:delete" class="mr-5px"/>
                        {{ t('main.remove') }}
                      </ElButton>
                    </template>
                  </ElPopconfirm>
                </div>
              </div>

            </ElForm>

          </ElCard>

        </ElCollapseItem>
      </ElCollapse>
      <!-- /props -->

    </ElCol>
  </ElRow>
  <!-- /chart items -->

  <ElRow class="mt-20px mb-20px" v-if="currentItem.entity">
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

</template>

<style lang="less" >

</style>
