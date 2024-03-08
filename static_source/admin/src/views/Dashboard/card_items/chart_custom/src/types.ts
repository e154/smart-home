import {EChartsOption} from "echarts";
import {ApiImage} from "@/api/stub";

export enum chartItemType {
    ATTR = "attr",
    METRIC = "metric",
    CUSTOM = "custom"
}

export interface CustomAttribute {
    value: string
    description: string
}

export interface ChartData {
    labels: Array<string>
    datasets: Array<ChartDataSet>
    lastTime?: string
}

export interface SeriesItem {
    chartType: chartItemType,
    attrAutomatic: boolean,
    customAttributes: Array<CustomAttribute>,
    metricIndex: number,
    metricProps: string,
    metricRange: "6h" | "12h" | "24h" | "7d" | "30d" | "1m",
    metricFilter: string,
    chartData: ChartData,
}

export interface ItemPayloadChartCustom {
    image?: ApiImage,
    chartOptions: EChartsOption,
    seriesItems: Array<SeriesItem>,
}

export interface ChartDataSet {
    label?: string,
    data: Array<number>
}

export interface ChartDataInterface {
    labels: Array<string>;
    datasets: Array<ChartDataSet>;
    lastTime?: string;
}

export interface ItemsType {
    label: string;
    value: string;
}

export const RangeList: ItemsType[] = [
    {label: '6 Hours', value: '6h'},
    {label: '12 Hours', value: '12h'},
    {label: '24 Hours', value: '24h'},
    {label: '7 Days', value: '7d'},
    {label: '30 Days', value: '30d'},
];

export const FilterList: ItemsType[] = [
    {label: 'secToTime', value: 'secToTime'},
    {label: 'formatdate', value: 'formatdate'},
    {label: 'formatBytes', value: 'formatBytes'},
    {label: 'seconds', value: 'seconds'},
    {label: 'getDayOfWeek', value: 'getDayOfWeek'},
];

export const defaultData = `{
  "grid": {
    "top": 10,
    "left": 30,
    "right": 0,
    "bottom": 20
  },
  "xAxis": {
    "data": [],
    "show": false,
    "type": "category"
  },
  "yAxis": {
    "show": true,
    "type": "value",
    "scale": true
  },
  "series": [
    {
      "data": [
    
      ],
      "type": "line",
      "smooth": false,
      "animation": false,
      "lineStyle": 1
    },
    {
      "data": [
       
      ],
      "type": "line",
      "smooth": false,
      "animation": false,
      "lineStyle": 1
    }
  ],
  "tooltip": {
    "padding": [
      5,
      10
    ],
    "trigger": "axis",
    "axisPointer": {
      "type": "cross"
    }
  },
  "responsive": true,
  "maintainAspectRatio": false
}
`
