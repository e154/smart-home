<template>
  <Bar
    :chart-options="chartOptions"
    :chart-data="chartData"
    :chart-id="chartId"
    :dataset-id-key="datasetIdKey"
    :plugins="plugins"
    :css-classes="cssClasses"
    :styles="styles"
    :width="width"
    :height="height"
    ref="bar"
  />
</template>

<script lang="ts">
import {Component, Prop, Vue} from 'vue-property-decorator';

import {Bar} from 'vue-chartjs/legacy';
import {BarElement, CategoryScale, Chart as ChartJS, Legend, LinearScale, Title, Tooltip} from 'chart.js';

ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale);


@Component({
  name: 'BarChart',
  components: {Bar}
})
export default class extends Vue {

  @Prop() private bus!: Vue;
  @Prop({default: 'bar-chart'}) private chartId!: string;
  @Prop({default: 'label'}) private datasetIdKey!: string;
  @Prop({default: '400'}) private width!: number;
  @Prop({default: '400'}) private height!: number;
  @Prop({default: 'bar-chart'}) private cssClasses!: string;
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
    if (!this.$refs.bar) {
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
