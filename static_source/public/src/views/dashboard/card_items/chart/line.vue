<template>
  <LineChartGenerator
    :chart-options="chartOptions"
    :chart-data="chartData"
    :chart-id="chartId"
    :dataset-id-key="datasetIdKey"
    :plugins="plugins"
    :css-classes="cssClasses"
    :styles="styles"
    :width="width"
    :height="height"
    ref="line"
  />
</template>

<script lang="ts">
import {Component, Prop, Vue} from 'vue-property-decorator';
import {Line as LineChartGenerator} from 'vue-chartjs/legacy';
import {
  CategoryScale,
  Chart as ChartJS,
  Legend,
  LinearScale,
  LineElement,
  PointElement,
  Title,
  Tooltip
} from 'chart.js';

ChartJS.register(
  Title,
  Tooltip,
  Legend,
  LineElement,
  LinearScale,
  CategoryScale,
  PointElement
);


@Component({
  name: 'LineChart',
  components: {LineChartGenerator}
})
export default class extends Vue {

  @Prop() private bus!: Vue;
  @Prop({default: 'line-chart'}) private chartId!: string;
  @Prop({default: 'label'}) private datasetIdKey!: string;
  @Prop({default: '400'}) private width!: number;
  @Prop({default: '400'}) private height!: number;
  @Prop({default: 'line-chart'}) private cssClasses!: string;
  @Prop({
    default: () => {
    }
  }) private styles!: Object;
  @Prop({default: () => []}) private plugins!: Array<Object>;
  @Prop({
    default: () => {
    }
  }) private chartData!: Object;
  @Prop({
    default: () => {
    }
  }) private chartOptions!: Object;

  private created() {
    this.bus.$on('updateChart', (t: string) => {
      this.updateChart();
    });
  }

  private mounted() {
  }

  public updateChart() {
    // @ts-ignore
    if (!this.$refs.line) {
      return;
    }
    try {
      // @ts-ignore
      this.$refs.line.updateChart();
    } catch (e) {
    }
  }
}
</script>

<style scoped>

</style>
