import {ItemsType} from '@/views/Dashboard/card_items';

export interface ItemPayloadChartItem {
  value: string
  description: string
}

export interface ItemPayloadChart {
  type: string,
  props: Array<string>,
  metric_index: number,
  width: number,
  borderWidth: number,
  xAxis: boolean,
  yAxis: boolean,
  legend: boolean,
  range: string,
  filter: string,
  color: string,
  backgroundColor: string,
  scale: boolean,
  automatic: boolean,
  dataZoom: boolean,
  items: Array<ItemPayloadChartItem>,
}

export interface ChartDataSet {
  label?: string,
  borderColor?: string,
  backgroundColor?: string,
  radius?: number,
  borderWidth?: number,
  data: Array<number>
}

export interface ChartDataInterface {
  labels: Array<string>;
  datasets: Array<ChartDataSet>;
  lastTime?: string;
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
