<script setup lang="ts">
import {computed, onMounted, onUnmounted, PropType, ref, watch} from "vue";
import {CardItem, parsedObject, serializedObject} from "@/views/Dashboard/core";
import {
  ChartDataInterface,
  ChartDataSet,
  chartItemType,
  SeriesItem
} from "@/views/Dashboard/card_items/chart_custom/types";
import {parseTime} from "@/utils";
import api from "@/api/api";
import {EChartsOption} from 'echarts'
import {Echart} from '@/components/Echart'
import {debounce} from "lodash-es";
import {UUID} from "uuid-generator-ts";
import stream from "@/api/stream";
import {Cache, GetTokens, RenderText} from "@/views/Dashboard/render";
import {ApiMetric} from "@/api/stub";

// ---------------------------------
// common
// ---------------------------------
const props = defineProps({
  item: {
    type: Object as PropType<CardItem>,
  },
})

const currentItem = computed(() => props.item as CardItem)

const onEventhandler = (event) => {
  const {entity_id} = event;

  if (entity_id != currentItem.value.entityId) {
    return
  }

  prepareData();
}

const currentID = ref('')
const el = ref(null)
onMounted(() => {
  // store dom element moveable
  props.item.setTarget(el.value)

  const uuid = new UUID()
  currentID.value = uuid.getDashFreeUUID()

  setTimeout(() => {
    stream.subscribe('event_updated_metric', currentID.value, onEventhandler);
  }, 1000)
})

onUnmounted(() => {
  stream.unsubscribe('event_updated_metric', currentID.value);
})

// ---------------------------------
// component methods
// ---------------------------------

const applyFilter = (value: any, filter: string): any => {
  switch (filter) {
    case 'formatBytes':
      const bytes = parseInt(value);
      if (bytes === 0) {
        return '0 Bytes';
      }
      const decimals = 2;
      const k = 1024;
      const dm = decimals < 0 ? 0 : decimals;
      // const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

      const i = Math.floor(Math.log(bytes) / Math.log(k));

      return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) /*+ ' ' + sizes[i]*/;
    default:
      console.warn(`unknown filter "${filter}"!`);
      return value;
  }
}

const prepareMetric = (metric: ApiMetric, series: SeriesItem): ChartDataInterface => {
  // console.log(metric)
  let _chartData: ChartDataInterface = {
    labels: series.chartData?.labels || [],
    datasets: series.chartData?.datasets || [],
  };

  // exit if no data
  if (!props.item?.entity?.metrics || !series.metricProps) {
    return _chartData;
  }

  // add time last item
  if (metric.data.length > 0) {
    _chartData.lastTime = metric.data[metric.data.length - 1].time
  }

  let totalLabels: Array<string> = [series.metricProps];

  // update only
  if (_chartData.datasets.length) {

    for (const t in metric.data) {
      _chartData.labels.push(parseTime(metric.data[t].time) as string);
      for (const l in totalLabels) {
        for (const j in _chartData.datasets) {
          if (_chartData.datasets[j].label == totalLabels[l]) {
            if (!series.metricFilter) {
              _chartData.datasets[j].data.push(metric.data[t].value[totalLabels[l]]);
            } else {
              const data = applyFilter(metric.data[t].value[totalLabels[l]], series.metricFilter);
              _chartData.datasets[j].data.push(data);
            }
          }
        }
      }
    }

    // 3600 max item per data
    const diff = _chartData.datasets.length - 3600;
    if (diff > 0) {
      _chartData.labels.slice(diff)
      _chartData.datasets.slice(diff)
    }

    return _chartData
  } // \update only

  // create full data
  let dataSets = new Map<string, ChartDataSet>();

  // create data sets
  for (const i in metric.options?.items) {
    // totalLabels.push(metric.options?.items[i].name);
    dataSets[metric.options?.items[i].name] = {
      label: metric.options?.items[i].name,
      data: new Array<number>(),
    };
  }

  // add data to sets
  for (const t in metric.data) {
    _chartData.labels.push(parseTime(metric.data[t].time) as string);
    for (const l in totalLabels) {
      if (!series.metricFilter) {
        dataSets[totalLabels[l]].data.push(metric.data[t].value[totalLabels[l]]);
      } else {
        const data = applyFilter(metric.data[t].value[totalLabels[l]], series.metricFilter);
        dataSets[totalLabels[l]].data.push(data);
      }
    }
  }

  for (const l in totalLabels) {
    _chartData.datasets.push(dataSets[totalLabels[l]]);
  }

  // console.log(_chartData);
  return _chartData;
}

const fetchMetric = async (id: number, series: SeriesItem): Promise<ApiMetric> => {
  const {data} = await api.v1.metricServiceGetMetric({
    id: id,
    range: series.metricRange || '24h',
    startDate: series.chartData?.lastTime,
  });

  return data;
}

const chartOptions = ref<EChartsOption>({})
const loaded = ref(false)

const _cache = new Cache()
var _chartOptions: EChartsOption = null;

const prepareData = debounce(async () => {
  // loaded.value = false

  const chart = props.item.payload.chartCustom;

  if (!chart?.chartOptions) {
    chartOptions.value = {} as EChartsOption;
    loaded.value = true
    return;
  }

  let firstTime = false;
  if (!_chartOptions) {
    firstTime = true;
    // let _chartOptions = Object.assign({}, chart?.chartOptions) as EChartsOption;
    // for testing !!!
    _chartOptions = parsedObject(serializedObject(chart?.chartOptions)) as EChartsOption;
    // console.log(_chartOptions)
  }

  if (chart?.seriesItems) {
    // console.log(chart.seriesItems)
    for (const i in chart.seriesItems) {

      let series = chart.seriesItems[i];
      // console.log(series)

      let rowData: any[] = []

      if (series.chartType == chartItemType.CUSTOM) {
        continue;
      } // CUSTOM

      if (series.chartType == chartItemType.ATTR) {

        if (series.attrAutomatic) {
          // console.log(props.item.lastEvent.new_state.attributes)
          if (props.item.lastEvent.new_state.attributes) {
            for (const k in props.item.lastEvent.new_state.attributes) {
              if (['int', 'float'].includes(props.item.lastEvent.new_state.attributes[k].type)) {
                rowData.push({value: props.item.lastEvent.new_state.attributes[k].value, name: k})
              }
            }
          }

        } else {
          for (const i in series.customAttributes) {
            let v: string = series.customAttributes[i].value || ''

            const tokens = GetTokens(series.customAttributes[i].value, _cache)
            if (tokens.length) {
              v = RenderText(tokens, v, props.item?.lastEvent)
            }
            rowData.push({value: parseInt(v) || 0, name: series.customAttributes[i].description})
          }
        }

      } // ATTR

      if (series.chartType == chartItemType.METRIC) {

        if (props.item.entity.metrics && series?.metricIndex != undefined) {
          let metric = props.item.entity.metrics[series?.metricIndex || 0];
          if (metric?.id) {
            metric = await fetchMetric(metric.id!, series);
          }
          // console.log(metric)
          series.chartData = prepareMetric(metric, series);
        }

        // console.log(series.chartData)


      } // METRIC

      // console.log(series.chartData)

      if (_chartOptions?.series[i]) {

        switch (_chartOptions.series[i].type) {
          case 'line':
            // matrics 1+ !xAxis|!yAxis
            if (!_chartOptions.xAxis) {
              _chartOptions["xAxis"] = {
                data: [],
              }
            }
            if (!_chartOptions.yAxis) {
              _chartOptions["yAxis"] = {
                type: 'value'
              }
            }
            _chartOptions.xAxis.data = series?.chartData?.labels || [];

            if (series?.chartData?.datasets) {
              _chartOptions.series[i].data = series?.chartData?.datasets[0].data || [];
              _chartOptions.series[i].name = series?.chartData?.datasets[0].label || '';

              if (_chartOptions.legend?.data) {
                _chartOptions.legend?.data.push(series?.chartData?.datasets[0].label)
              }
            }
            break;
          case 'bar':
            // attrs (custom|auto) !xAxis|!yAxis
            if (!_chartOptions.xAxis) {
              _chartOptions["xAxis"] = {
                type: 'category',
                data: [],
              }
            }
            if (!_chartOptions.yAxis) {
              _chartOptions["yAxis"] = {
                type: 'value'
              }
            }
            _chartOptions.xAxis.data = rowData.map(function (item) {
              return item.name;
            });
            _chartOptions.series[i].data = rowData.map(function (item) {
              return item.value;
            });
            break;
          case 'doughnut':
          case 'pie':
          case 'gauge':
            // attrs (custom|auto)
            if (firstTime) {
              _chartOptions.series[i].data = rowData;
            } else {
              chartOptions.value = {
                series: [{
                  data: rowData
                }]
              }
              return
            }
            break;
        }
      }
    }
  }

  // console.log(_chartOptions)

  chartOptions.value = _chartOptions;

  loaded.value = true
}, 500)

const clear = () => {
  _chartOptions = null;
  const chart = props.item.payload.chartCustom;
  if (chart?.seriesItems) {
    for (const i in chart.seriesItems) {
      chart.seriesItems[i].chartData = {}
    }
  }
}

const reloadKey = ref(0)
const reload = debounce(() => {
  reloadKey.value += 1
  // console.log('reload')
  // requestCurrentState(props.item?.entityId);
}, 100)

watch(
    () => props.item.lastEvent,
    (val?: CardItem) => {
      prepareData();
      // reload();
    },
    {
      deep: true,
      immediate: true
    }
)

watch(
    () => [props.item.width, props.item.height, props.item.payload?.chartCustom?.chartOptions],
    (width, height, chartOptions) => {
      clear();
      prepareData();
      reload();
    },
    {
      deep: true,
      immediate: true
    }
)

watch(
    () => props.item.hidden,
    (val?: boolean) => {
      if (!val) {
        reload();
      }
    }
)

// requestCurrentState(props.item?.entityId);
prepareData();

</script>

<template>
  <div ref="el" :class="[{'hidden': item.hidden}]" class="h-[100%] w-[100%]">
    <Echart v-if="item.enabled && item.entity && loaded" :options="chartOptions" :key="reloadKey"/>
  </div>
</template>

<style lang="less">

</style>
