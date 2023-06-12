<template>
  <div class="app-container dashboard-container" v-if="!loading">

    <el-tabs type="card" v-model="board.activeTab" v-if="board.tabs.length > 1">
      <el-tab-pane
        v-for="(tab, index) in board.tabs"
        :label="tab.name"
        :key="index"
        :style="{backgroundColor: tab.background}"
        :class="[{'gap': tab.gap}]"
      >
        <dashboard-tab-muu :tab="tab" :bus="bus"/>

      </el-tab-pane>
    </el-tabs>

    <div v-if="board.tabs.length == 1" :class="[{'gap': board.tabs[0].gap}]">
      <dashboard-tab-muu :tab="board.tabs[0]" :bus="bus"/>
    </div>

    <el-empty v-if="board.tabs.length === 0" :image-size="200" description="no data"></el-empty>

  </div>
</template>

<script lang="ts">
import {Component, Prop, Vue} from 'vue-property-decorator';
import api from '@/api/api';
import stream from '@/api/stream';
import {UUID} from 'uuid-generator-ts';
import {Core, requestCurrentState} from '@/views/dashboard/core';
import {EventStateChange} from '@/api/stream_types';
import DashboardTabMuu from '@/views/dashboard/view/tab-muu.vue';

@Component({
  name: 'Dashboard',
  components: {
    DashboardTabMuu
  }
})
export default class extends Vue {
  @Prop({required: false}) private id!: number;

  private loading = true;
  private bus: Vue = new Vue();

  private board: Core = new Core(this.bus);

  // id for streaming subscribe
  private currentID = '';

  private created() {
    const uuid = new UUID();
    this.currentID = uuid.getDashFreeUUID();

    this.fetchDashboard();

    setTimeout(() => {
      stream.subscribe('state_changed', this.currentID, this.onStateChanged);

      for (const entityId in this.board.current.entities) {
        requestCurrentState(entityId);
      }
    }, 1000);
  }

  private destroyed() {
    stream.unsubscribe('state_changed', this.currentID);
  }

  private mounted() {

  }

  private async fetchDashboard() {
    this.loading = true;

    api.v1.variableServiceGetVariableByName('devDashboard').then((resp) => {

      if (!resp?.data?.value) {
        this.loading = false;
        return;
      }

      const id = parseInt(resp?.data?.value);

      api.v1.dashboardServiceGetDashboardById(id).then((resp) => {
        this.loading = false;

        this.board.currentBoard(resp.data);
      });
    });
  }

  private onStateChanged(event: EventStateChange) {
    this.bus.$emit('state_changed', event);
    this.board.onStateChanged(event);
  }
}
</script>
