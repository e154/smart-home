<template>
  <div>
    <common-editor :item="item" :board="board"></common-editor>

    <el-divider content-position="left">Chart options</el-divider>

    <el-form-item :label="$t('dashboard.editor.chart.type')" prop="type">
      <el-select
        v-model="item.payload.chart.type"
        placeholder="please select type"
        style="width: 100%"
      >
        <el-option label="linear" value="line"></el-option>
        <el-option label="bar" value="bar"></el-option>
        <el-option label="radar" value="radar"></el-option>
        <el-option label="doughnut" value="doughnut"></el-option>
      </el-select>
    </el-form-item>

    <el-form-item :label="$t('dashboard.editor.chart.entity_metric')" prop="index">
      <el-select v-model="item.payload.chart.metric_index" placeholder="Select">
        <el-option
          v-for="(prop, index) in item.entity.metrics"
          :key="index"
          :label="prop.name"
          :value="index">
        </el-option>
      </el-select>
    </el-form-item>

    <div v-if="item.entity.metrics && item.payload.chart.metric_index !== undefined">
      <el-form-item :label="$t('dashboard.editor.chart.metric_props')" prop="index">
        <el-select v-model="item.payload.chart.props" multiple placeholder="Select">
          <el-option
            v-for="(props, index) in item.entity.metrics[item.payload.chart.metric_index].options.items"
            :key="props.name"
            :label="props.name"
            :value="props.name">
          </el-option>
        </el-select>
      </el-form-item>
    </div>

    <el-form-item :label="$t('dashboard.editor.chart.range')" prop="index">
      <el-select v-model="item.payload.chart.range" placeholder="Select">
        <el-option
          v-for="(props, index) in rangeList"
          :key="props.value"
          :label="props.label"
          :value="props.value">
        </el-option>
      </el-select>
    </el-form-item>

    <el-form-item :label="$t('dashboard.editor.chart.filter')" prop="index">
      <el-select v-model="item.payload.chart.filter" placeholder="Select" clearable>
        <el-option
          v-for="(props, index) in filterList"
          :key="props.value"
          :label="props.label"
          :value="props.value">
        </el-option>
      </el-select>
    </el-form-item>

    <el-form-item :label="$t('dashboard.editor.chart.borderWidth')" prop="borderWidth">
      <el-input-number size="small"
                       v-model="item.payload.chart.borderWidth"></el-input-number>
    </el-form-item>

    <el-form-item :label="$t('dashboard.editor.chart.legend')" prop="legend">
      <el-switch
        v-model="item.payload.chart.legend"></el-switch>
    </el-form-item>

    <el-form-item :label="$t('dashboard.editor.chart.xAxis')" prop="xAxis">
      <el-switch
        v-model="item.payload.chart.xAxis"></el-switch>
    </el-form-item>

    <el-form-item :label="$t('dashboard.editor.chart.yAxis')" prop="yAxis">
      <el-switch
        v-model="item.payload.chart.yAxis"></el-switch>
    </el-form-item>

    <el-row style="padding-bottom: 20px">
      <el-col>
        <event-viewer :item="item"></event-viewer>
      </el-col>
    </el-row>

  </div>
</template>

<script lang="ts">
import {Component, Prop, Vue} from 'vue-property-decorator';
import {CardItem, Core} from '@/views/dashboard/core';
import CommonEditor from '@/views/dashboard/card_items/common/editor.vue';
import EventViewer from '@/views/dashboard/card_items/common/event_viewer.vue';
import {RangeList, FilterList} from '@/views/dashboard/card_items/chart/types';

@Component({
  name: 'IChartEditor',
  components: {
    CommonEditor,
    EventViewer
  }
})
export default class extends Vue {
  @Prop() private item!: CardItem;
  @Prop() private board!: Core;

  private rangeList = RangeList;
  private filterList = FilterList;

  private created() {
  }

  private mounted() {
  }
}
</script>

<style scoped>

</style>
