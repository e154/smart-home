<script setup lang="ts">
import {computed, onMounted, PropType, ref, unref, watch} from "vue";
import {Card, CardItem, Core, requestCurrentState, Tab} from "@/views/Dashboard/core";
import {ApiMetric} from "@/api/stub";
import {ChartDataInterface, ChartDataSet} from "@/views/Dashboard/card_items/chart/types";
import {parseTime} from "@/utils";
import api from "@/api/api";
import {useBus} from "@/views/Dashboard/bus";
import { EChartsOption } from 'echarts'
import { Echart } from '@/components/Echart'
import {useI18n} from "@/hooks/web/useI18n";
import {string} from "vue-types";
import {debounce} from "lodash-es";
import {Cache, GetTokens, RenderText} from "@/views/Dashboard/render";
import {ApplyFilter} from "@/views/Dashboard/filters";

const {bus} = useBus()
const { t } = useI18n()

// ---------------------------------
// common
// ---------------------------------
const props = defineProps({
  item: {
    type: Object as PropType<CardItem>,
  },
})

const el = ref(null)
onMounted(() => {
  // store dom element moveable
  props.item.setTarget(el.value)
})

const currentItem = computed(()=> props.item as CardItem)

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

const prepareMetric = (metric: ApiMetric): ChartDataInterface => {
  // console.log(metric)
  let _chartData: ChartDataInterface = {
    labels:  chartData.value?.labels || [],
    datasets: chartData.value?.datasets || [],
  };

  // exit if no data
  if (!props.item?.entity?.metrics || !props.item.payload.chart?.props || props.item.payload.chart?.props.length == 0) {
    return _chartData;
  }

  // add time last item
  if (metric.data.length > 0) {
    _chartData.lastTime = metric.data[metric.data.length -1].time
  }

  let totalLabels: Array<string> = props.item.payload.chart?.props;

  // update only
  if (_chartData.datasets.length) {

    for (const t in metric.data) {
      _chartData.labels.push(parseTime(metric.data[t].time) as string);
      for (const l in totalLabels) {
        for (const j in _chartData.datasets) {
          if (_chartData.datasets[j].label == totalLabels[l]) {
            if (!props.item.payload.chart?.filter) {
              _chartData.datasets[j].data.push(metric.data[t].value[totalLabels[l]]);
            } else {
              const data = applyFilter(metric.data[t].value[totalLabels[l]], props.item.payload.chart?.filter);
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
      borderColor: metric.options?.items[i].color,
      backgroundColor: metric.options?.items[i].color,
      data: new Array<number>(),
      borderWidth: props.item.payload.chart?.borderWidth || 1,
      radius: 0,
    };
  }

  // add data to sets
  for (const t in metric.data) {
    _chartData.labels.push(parseTime(metric.data[t].time) as string);
    for (const l in totalLabels) {
      if (!props.item.payload.chart?.filter) {
        dataSets[totalLabels[l]].data.push(metric.data[t].value[totalLabels[l]]);
      } else {
        const data = applyFilter(metric.data[t].value[totalLabels[l]], props.item.payload.chart?.filter);
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

const fetchMetric = async (id: number, from?: string): Promise<ApiMetric> => {
  const {data} = await api.v1.metricServiceGetMetric({
    id: id,
    range: props.item?.payload?.chart?.range || '24h',
    from: from,
  });

  return data;
}

const chartData = ref<{
  labels: Array<string>
  datasets: Array<ChartDataSet>
  lastTime?: string
}>({
  labels: [],
  datasets: []
})

const chartOptions = ref<EChartsOption>({})
const loaded = ref(false)
const fistTime = ref(true)

const getLineOptions = () => {
  let legendData: string[] = [];
  let series: any[] = []
  let color: string[] = []

  for (const i in chartData.value.datasets) {
    if (!chartData.value.datasets[i]) {
      continue;
    }
    if (chartData.value.datasets[i]?.label) {
      legendData.push(chartData.value.datasets[i].label)
    }

    let lineStyle = {
      width: chartData.value.datasets[i].borderWidth || 1
    }

    if (props.item.payload.chart.color) {
      color.push(props.item.payload.chart.color /*|| chartData.value.datasets[i].borderColor*/)
      lineStyle["color"] = props.item.payload.chart.color /*|| chartData.value.datasets[i].borderColor*/
    }

    let row = {
      name: chartData.value.datasets[i].label,
      smooth: false,
      showSymbol: false,
      type: 'line',
      data: chartData.value.datasets[i].data,
      animationDuration: 2800,
      animationEasing: 'cubicInOut',
      lineStyle: lineStyle,

    }
    if (props.item.payload.chart?.backgroundColor) {
      row['areaStyle'] = {
        color: props.item.payload.chart.backgroundColor
      }
    }
    series.push(row)
  }

  let options: EChartsOption = {
    xAxis: {
      show: props.item.payload.chart?.xAxis,
      data: chartData.value.labels,
      axisLine: {
        show: false
      },
    },
    yAxis: {
      type: 'value',
      show: props.item.payload.chart?.yAxis,
      axisLine: {
        show: false
      },
      scale: props.item.payload.chart?.scale || false,
    },
    grid: {
      left: 0,
      right: 0,
      bottom: 0,
      top: 0,
      containLabel: props.item.payload.chart?.xAxis || props.item.payload.chart?.yAxis
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross'
      },
      padding: [5, 10]
    },
    series: series,
    responsive: true,
    maintainAspectRatio: false,
    animation: false,
  } as EChartsOption

  if (color.length) {
    options.color = color
  }

  if (props.item.payload.chart?.dataZoom) {
    options.grid.bottom = 40;
    options.dataZoom = [
      {
        type: 'inside',
      },
      {
        start: 0,
        end: 10
      }
    ]
  }

  if (props.item.payload.chart?.legend) {
    options.legend = {
      data: legendData,
      top: 10
    }
  }

  chartOptions.value = options as EChartsOption
}

const _cache = new Cache()
const getBarOptions = () => {

  let series: any[] = []
  let xAxisData: string[] = []
  let rowData: number[] = []

  if (props.item.payload.chart.automatic) {
    if (props.item.lastEvent.new_state.attributes) {
      for (const k in props.item.lastEvent.new_state.attributes) {
        if (['int', 'float'].includes(props.item.lastEvent.new_state.attributes[k].type)) {
          xAxisData.push(k)
          rowData.push(props.item.lastEvent.new_state.attributes[k].value)
        }
      }
    }

  } else {

    for (const i in props.item.payload.chart.items) {
      xAxisData.push(props.item.payload.chart.items[i].description || '')
      let v: string = props.item.payload.chart.items[i].value || ''

      const tokens = GetTokens(props.item.payload.chart.items[i].value, _cache)
      if (tokens.length) {
        v = RenderText(tokens, v, props.item?.lastEvent)
      }
      rowData.push(parseInt(v) || 0)
    }
  }

  let row = {
    data: rowData,
    type: 'bar'
  }

  if (props.item.payload.chart.color) {
    row['color'] = props.item.payload.chart.color
  }

  if (props.item.payload.chart.backgroundColor) {
    row['showBackground'] = true
    row['backgroundStyle'] = {
      color: props.item.payload.chart.backgroundColor
    }
  }

  series.push(row)

  chartOptions.value = {
    xAxis: {
      type: 'category',
      show: props.item.payload.chart?.xAxis,
      data: xAxisData,
      // data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'], //test data
      axisLine: {
        show: false
      },
    },
    yAxis: {
      type: 'value',
      show: props.item.payload.chart?.yAxis,
      axisLine: {
        show: false
      },
      scale: props.item.payload.chart?.scale || false,
    },
    grid: {
      left: 0,
      right: 0,
      bottom: 0,
      top: 0,
      containLabel: props.item.payload.chart?.xAxis || props.item.payload.chart?.yAxis
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross'
      },
      padding: [5, 10]
    },
    series: series,
    responsive: true,
    maintainAspectRatio: false,
    animation: false,
  } as EChartsOption

}

const getPieOptions = () => {

  let series: any[] = []
  let rowData: any[] = []

  if (props.item.payload.chart.automatic) {
    if (props.item.lastEvent.new_state.attributes) {
      for (const k in props.item.lastEvent.new_state.attributes) {
        if (['int', 'float'].includes(props.item.lastEvent.new_state.attributes[k].type)) {
          rowData.push({value: props.item.lastEvent.new_state.attributes[k].value, name: k})
        }
      }
    }

  } else {

    for (const i in props.item.payload.chart.items) {
      let v: string = props.item.payload.chart.items[i].value || ''

      const tokens = GetTokens(props.item.payload.chart.items[i].value, _cache)
      if (tokens.length) {
        v = RenderText(tokens, v, props.item?.lastEvent)
      }
      rowData.push({value: parseInt(v) || 0, name: props.item.payload.chart.items[i].description})
    }
  }

  let row = {
    data: rowData,
    type: 'pie',
    radius: '50%',
  }

  series.push(row)

  chartOptions.value = {
    grid: {
      left: 0,
      right: 0,
      bottom: 0,
      top: 0,
      containLabel: props.item.payload.chart?.xAxis || props.item.payload.chart?.yAxis
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross'
      },
      padding: [5, 10]
    },
    series: series,
    responsive: true,
    maintainAspectRatio: false,
    animation: false,
  } as EChartsOption

}

const getDoughnutOptions = () => {

  let series: any[] = []
  let rowData: any[] = []

  if (props.item.payload.chart.automatic) {
    if (props.item.lastEvent.new_state.attributes) {
      for (const k in props.item.lastEvent.new_state.attributes) {
        if (['int', 'float'].includes(props.item.lastEvent.new_state.attributes[k].type)) {
          rowData.push({value: props.item.lastEvent.new_state.attributes[k].value, name: k})
        }
      }
    }

  } else {

    for (const i in props.item.payload.chart.items) {
      let v: string = props.item.payload.chart.items[i].value || ''

      const tokens = GetTokens(props.item.payload.chart.items[i].value, _cache)
      if (tokens.length) {
        v = RenderText(tokens, v, props.item?.lastEvent)
      }
      rowData.push({value: parseInt(v) || 0, name: props.item.payload.chart.items[i].description})
    }
  }

  let row = {
    data: rowData,
    type: 'pie',
    radius: ['40%', '70%'],
    avoidLabelOverlap: false,
    itemStyle: {
      borderRadius: 10,
      borderColor: '#fff',
      borderWidth: 2
    },
  }

  series.push(row)

  chartOptions.value = {
    grid: {
      left: 0,
      right: 0,
      bottom: 0,
      top: 0,
      containLabel: props.item.payload.chart?.xAxis || props.item.payload.chart?.yAxis
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross'
      },
      padding: [5, 10]
    },
    series: series,
    responsive: true,
    maintainAspectRatio: false,
    animation: false,
  } as EChartsOption

}

const prepareData = debounce( async ()  => {

  switch (props.item.payload.chart.type) {
    case 'bar':
      getBarOptions()
      break
    case 'pie':
      getPieOptions()
      break
    case 'doughnut':
      getDoughnutOptions()
      break
  }

  switch (props.item.payload.chart.type) {
    case 'bar':
    case 'pie':
    case 'doughnut':
      loaded.value = true;
      fistTime.value = false
      bus.emit('updateChart', props.item.payload.chart.type)
      return
  }

  if (!props.item?.entity?.metrics || !props.item.payload?.chart?.type) {
    return;
  }

  let metric = props.item.entity.metrics[props.item.payload.chart?.metric_index || 0];

  if (metric?.id) {
    metric = await fetchMetric(metric.id!, chartData.value?.lastTime);
  }
  chartData.value = prepareMetric(metric);

  switch (props.item.payload.chart.type) {
    case 'line':
      getLineOptions()
      break;
    default:
      console.warn(`unknown chart type ${props.item.entity.metrics[props.item.payload.chart?.metric_index || 0].type}`);
  }

  bus.emit('updateChart', props.item.payload.chart.type)

  loaded.value = true;
  fistTime.value = false

  // console.log(lineOptions.value)
  // console.log(chartData.value)
  // console.log(metric)
}, 500)

const reloadKey = ref(0)
const reload = debounce(() => {
    reloadKey.value += 1
    requestCurrentState(props.item?.entityId);
  }, 100
)

watch(
    () => props.item.lastEvent,
    (val?: CardItem) => {
      prepareData();
    },
    {
      deep: true,
      immediate: true
    }
)

watch(
    () => [props.item.width, props.item.height, props.item.payload.chart],
    (width, height, chart) => {
      reload();
    },
    {
      deep: true,
      immediate: true
    }
)

requestCurrentState(props.item?.entityId);

</script>

<template>
  <div ref="el" class="h-[100%] w-[100%]">
    <Echart v-if="item.entity && loaded" :options="chartOptions" :key="reloadKey"/>
  </div>
</template>

<style lang="less" >

</style>
