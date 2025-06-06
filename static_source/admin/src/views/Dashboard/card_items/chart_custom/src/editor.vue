<script setup lang="ts">
import {computed, PropType} from "vue";
import {CardItem, Core, parsedObject} from "@/views/Dashboard/core";
import {
  ElButton,
  ElCard,
  ElCol,
  ElCollapse,
  ElCollapseItem,
  ElDivider,
  ElForm,
  ElFormItem,
  ElInput,
  ElOption,
  ElPopconfirm,
  ElRadioButton,
  ElRadioGroup,
  ElRow,
  ElSelect,
  ElSwitch
} from 'element-plus'
import {CommonEditor} from "@/views/Dashboard/card_items/common";
import {chartItemType, CustomAttribute, defaultData, FilterList, RangeList, SeriesItem} from "./types";
import {useI18n} from "@/hooks/web/useI18n";
import {Infotip} from "@/components/Infotip";
import {debounce} from "lodash-es";
import {EChartsOption} from "echarts";
import {ApiImage} from "@/api/stub";
import {JsonEditor} from "@/components/JsonEditor";
import {KeysSearch} from "@/views/Dashboard/components";

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

const currentItem = computed(() => props.item as CardItem)

// ---------------------------------
// component methods
// ---------------------------------

const getMetricList = computed(() => {
  if (props.item.entity && props.item.entity?.metrics) {
    return props.item.entity.metrics
  } else {
    return []
  }
})

const getMetricItem = (prop) => {
  if (!prop.metricIndex) {
    prop.metricIndex = 0;
  }
  if (props.item.entity && props.item.entity?.metrics) {
    return props.item.entity.metrics[prop.metricIndex].options.items
  } else {
    return []
  }
}

// series items
const addSeriesItem = () => {
  if (!currentItem.value.payload?.chartCustom?.seriesItems) {
    currentItem.value.payload.chartCustom.seriesItems = []
  }
  currentItem.value.payload.chartCustom?.seriesItems.push({
    chartType: chartItemType.CUSTOM,
  } as SeriesItem)
}

const removeSeriesItem = (index: number) => {
  if (!currentItem.value.payload.chartCustom.seriesItems) {
    return;
  }

  currentItem.value.payload.chartCustom.seriesItems.splice(index, 1);
}
// \series items

// attributes item
const addAttrItem = (prop: SeriesItem) => {
  if (!prop.customAttributes) {
    prop.customAttributes = [];
  }
  prop.customAttributes.push({
    value: 'value',
  } as CustomAttribute)
}

const removeAttrItem = (prop: SeriesItem, index: number) => {
  if (!prop.customAttributes) {
    return;
  }

  prop.customAttributes.splice(index, 1);
}

// attributes item

const editorHandler = debounce((val: any) => {
  if (!val) {
    val = {
      text: defaultData,
    }
  }

  try {

    let options: EChartsOption;

    if (val.json) {
      options = val.json as EChartsOption
    } else if (val.text) {
      options = parsedObject(val.text) as EChartsOption
    }

    if (!options.grid) {
      options['grid'] = {
        top: 10,
        left: 30,
        right: 0,
        bottom: 20
      };
    }

    if (options.responsive == undefined) {
      options['responsive'] = true;
    }

    if (options.animation == undefined) {
      options['animation'] = true;
    }

    if (options.maintainAspectRatio == undefined) {
      options['maintainAspectRatio'] = true;
    }

    if (options.tooltip == undefined) {
      options['tooltip'] = {
        trigger: 'axis',
        axisPointer: {
          type: 'cross'
        },
        padding: [5, 10]
      };
    }

    for (let series of options.series) {
      if (['line'].includes(series.type)) {
        if (options.xAxis == undefined) {
          options['xAxis'] = {
            type: 'value',
          };
        }
        if (options.yAxis == undefined) {
          options['yAxis'] = {
            type: 'value',
          };
        }
        if (options.yAxis?.scale == undefined) {
          options.yAxis['scale'] = false
        }
        if (options.yAxis?.show == undefined) {
          options.yAxis['show'] = false
        }
        if (options.xAxis?.show == undefined) {
          options.xAxis['show'] = false
        }
      }
      if (series.animation == undefined) {
        series['animation'] = false;
      }
      if (series.smooth == undefined) {
        series['smooth'] = false;
      }
      if (series.lineStyle == undefined) {
        series['lineStyle'] = 1;
      }
      if (series.showSymbol == undefined) {
        series['showSymbol'] = false;
      }
      if (series.animationDuration == undefined) {
        series['animationDuration'] = 2800;
      }
      if (series.animationEasing == undefined) {
        series['animationEasing'] = 'cubicInOut';
      }
    }

    currentItem.value.payload.chartCustom.chartOptions = options;

  } catch (e) {
    console.error(e)
  }

}, 500)

const onSelectImage = (image: ApiImage) => {
  if (!props.item?.payload?.chartCustom?.image) {
    return;
  }
  // console.log('select image', image);
  currentItem.value.payload.chartCustom.image = image || undefined;
}

const onChangePropValue = (val: string, prop: any, index: number): void => {
  prop.customAttributes[index].value = val
}

</script>

<template>

  <CommonEditor :item="item" :core="core"/>

  <Infotip
      type="warning"
      :show-index="false"
      title="WARNING"
      :schema="[
      {
        label: 'Experimental component, may change in the future',
      },
    ]"
  />

  <Infotip
      :show-index="false"
      :title="$t('dashboard.editor.chart.chartDocumentation')"
      :schema="[
      {
        label: t('dashboard.editor.chart.info1'),
      },
    ]"
  />

  <ElRow class="mb-10px mt-10px">
    <ElCol>
      <ElDivider content-position="left">{{ $t('dashboard.editor.chart.chartOptions') }}</ElDivider>
    </ElCol>
  </ElRow>

  <ElRow class="mb-10px" style="min-height: 200px">
    <ElCol>
      <JsonEditor v-model="currentItem.payload.chartCustom.chartOptions" height="auto" @change="editorHandler"/>
    </ElCol>
  </ElRow>

  <!-- chart items -->
  <ElRow class="mb-10px">
    <ElCol>
      <ElDivider content-position="left">{{ $t('dashboard.editor.chart.seriesOptions') }}</ElDivider>
    </ElCol>
  </ElRow>

  <ElRow class="mb-10px">
    <ElCol>
      <ElButton class="w-[100%]" @click.prevent.stop="addSeriesItem()">
        <Icon icon="ep:plus" class="mr-5px"/>
        {{ $t('dashboard.editor.chart.addSeriesItem') }}
      </ElButton>

      <!-- props -->
      <ElCollapse>
        <ElCollapseItem
            v-for="(prop, index) in currentItem.payload?.chartCustom?.seriesItems"
            :name="index"
            :key="index"
        >

          <template #title>
            {{ $t('dashboard.editor.chart.series') }} {{ index }}
          </template>

          <ElCard shadow="never" class="item-card-editor">

            <ElForm
                label-position="top"
                :model="prop"
                style="width: 100%"
                ref="cardItemForm">

              <ElRow class="mb-10px">
                <ElCol>
                  <ElFormItem :label="$t('dashboard.editor.chart.itemType')" prop="text">
                    <ElRadioGroup v-model="prop.chartType">
                      <ElRadioButton label="custom"/>
                      <ElRadioButton label="attr"/>
                      <ElRadioButton label="metric"/>
                    </ElRadioGroup>
                  </ElFormItem>
                </ElCol>
              </ElRow>

              <!-- attr -->
              <ElRow v-if="prop.chartType == 'attr'">
                <ElCol>
                  <ElFormItem :label="$t('dashboard.editor.chart.automatic')" prop="borderWidth">
                    <ElSwitch v-model="prop.attrAutomatic"/>
                  </ElFormItem>
                </ElCol>
              </ElRow>

              <ElRow v-if="prop.chartType == 'attr' && !prop.attrAutomatic">
                <ElCol>
                  <div class="mb-10px">
                    <ElButton class="w-[100%]" @click.prevent.stop="addAttrItem(prop)">
                      <Icon icon="ep:plus" class="mr-5px"/>
                      {{ $t('dashboard.editor.chart.addCustomAttribute') }}
                    </ElButton>
                  </div>

                  <!-- props -->
                  <ElCollapse v-if="prop.customAttributes?.length">
                    <ElCollapseItem
                        v-for="(attr, index) in prop.customAttributes"
                        :name="index"
                        :key="index"
                    >

                      <template #title>
                        {{ attr.value }}
                      </template>

                      <ElCard shadow="never" class="item-card-editor">

                        <ElForm
                            label-position="top"
                            :model="attr"
                            style="width: 100%"
                            ref="cardItemForm">

                          <ElRow>
                            <ElCol>
                              <ElFormItem :label="$t('dashboard.editor.attrField')" prop="text">
                                <!--                                <ElInput class="w-[100%]" placeholder="Please input" v-model="attr.value"/>-->
                                <KeysSearch v-model="attr.value" :obj="currentItem.lastEvent"
                                            @change="onChangePropValue($event, prop, index)"/>
                              </ElFormItem>
                            </ElCol>
                          </ElRow>

                          <ElRow class="mb-10px">
                            <ElCol>
                              <ElFormItem :label="$t('dashboard.editor.chart.itemDescription')" prop="text">
                                <ElInput class="w-[100%]" placeholder="Please input" v-model="attr.description"/>
                              </ElFormItem>
                            </ElCol>
                          </ElRow>

                          <div style="text-align: right;">
                            <ElPopconfirm
                                :confirm-button-text="$t('main.ok')"
                                :cancel-button-text="$t('main.no')"
                                width="250"
                                style="margin-left: 10px;"
                                :title="$t('main.are_you_sure_to_do_want_this?')"
                                @confirm="removeAttrItem(prop, index)"
                            >
                              <template #reference>
                                <ElButton type="danger" plain>
                                  <Icon icon="ep:delete" class="mr-5px"/>
                                  {{ t('main.remove') }}
                                </ElButton>
                              </template>
                            </ElPopconfirm>
                          </div>

                        </ElForm>

                      </ElCard>

                    </ElCollapseItem>
                  </ElCollapse>
                  <!-- /props -->

                </ElCol>
              </ElRow>
              <!-- /attr -->

              <!-- metric -->
              <ElRow v-if="prop.chartType == 'metric'">
                <ElCol>

                  <ElFormItem :label="$t('dashboard.editor.chart.entity_metric')" prop="index">
                    <ElSelect v-model="prop.metricIndex" placeholder="Select" clearable class="w-[100%]">
                      <ElOption
                          v-for="(prop, index) in getMetricList"
                          :key="index"
                          :label="prop.name"
                          :value="index"/>
                    </ElSelect>
                  </ElFormItem>
                </ElCol>
              </ElRow>

              <ElRow v-if="prop.chartType == 'metric'">
                <ElCol>
                  <ElFormItem :label="$t('dashboard.editor.chart.metric_props')" prop="index">
                    <ElSelect v-model="prop.metricProps" placeholder="Select" clearable class="w-[100%]">
                      <ElOption
                          v-for="p in getMetricItem(prop)"
                          :key="p.name"
                          :label="p.name"
                          :value="p.name"/>
                    </ElSelect>
                  </ElFormItem>
                </ElCol>
              </ElRow>

              <ElRow v-if="prop.chartType == 'metric'">
                <ElCol>
                  <ElFormItem :label="$t('dashboard.editor.chart.range')" prop="index">
                    <ElSelect v-model="prop.metricRange" placeholder="Select" clearable class="w-[100%]">
                      <ElOption
                          v-for="p in rangeList"
                          :key="p.value"
                          :label="p.label"
                          :value="p.value"/>
                    </ElSelect>
                  </ElFormItem>
                </ElCol>
              </ElRow>

              <ElRow v-if="prop.chartType == 'metric'">
                <ElCol>
                  <ElFormItem :label="$t('dashboard.editor.chart.filter')" prop="index">
                    <ElSelect v-model="prop.metricFilter" placeholder="Select" clearable class="w-[100%]">
                      <ElOption
                          v-for="p in filterList"
                          :key="p.value"
                          :label="p.label"
                          :value="p.value"/>
                    </ElSelect>
                  </ElFormItem>
                </ElCol>
              </ElRow>
              <!-- /metric -->

              <div class="mt-10px" style="text-align: right;">
                <ElPopconfirm
                    :confirm-button-text="$t('main.ok')"
                    :cancel-button-text="$t('main.no')"
                    width="250"
                    style="margin-left: 10px;"
                    :title="$t('main.are_you_sure_to_do_want_this?')"
                    @confirm="removeSeriesItem(index)"
                >
                  <template #reference>
                    <ElButton type="danger" plain>
                      <Icon icon="ep:delete" class="mr-5px"/>
                      {{ t('main.remove') }}
                    </ElButton>
                  </template>
                </ElPopconfirm>
              </div>

            </ElForm>

          </ElCard>

        </ElCollapseItem>
      </ElCollapse>
      <!-- /props -->

    </ElCol>
  </ElRow>
  <!-- /chart items -->

</template>

<style lang="less">

</style>
