<template>

  <div class="components-container dashboard-container" style="margin: 0" v-if="!loading">
    <split-pane
      split="horizontal"
      @resize="resize"
    >
      <template slot="paneL">
        <div class="top-container">
          <el-tabs
            v-model="board.activeTab"
            @edit="handleTabsEdit"
            @tab-click="updateCurrentTab"
          >
            <el-tab-pane
              v-for="(tab, index) in board.tabs"
              :label="tab.name"
              :key="index"
              :style="{backgroundColor: tab.background}"
              :class="[{'gap': tab.gap}]"
            >
              <editor-tab-muu :tab="tab" :bus="bus"/>

            </el-tab-pane>
          </el-tabs>

          <el-empty v-if="board.tabs.length === 0" :image-size="200"
                    :description="$t('dashboard.editor.please_add_tab')"></el-empty>

        </div>
      </template>
      <template slot="paneR">
        <div class="bottom-container">

          <el-tabs v-model="board.mainTab">
            <el-tab-pane :label="$t('dashboard.editor.dashboardSettings')" key="0">
              <tab-dashboard-settings v-if="board.current" :board="board.current" :core="board" :bus="bus"/>
            </el-tab-pane>
            <el-tab-pane :label="$t('dashboard.editor.tabList')" key="1">
              <tab-editor v-if="board.current" :board="board.current" :core="board" :bus="bus"
                          :tab="board.tabs[board.activeTab]" @update-value="fetchDashboard"/>
            </el-tab-pane>
            <el-tab-pane :label="$t('dashboard.editor.cardList')" key="2">
              <card-editor v-if="board.activeTab && board.tabs[board.activeTab]"
                           :board="board.current" :core="board" :bus="bus"
                           :tab="board.tabs[board.activeTab]"/>
            </el-tab-pane>
            <el-tab-pane :label="$t('dashboard.editor.cardItems')" key="3">
              <card-items-editor
                v-if="board.activeTab && board.tabs[board.activeTab] && board.tabs[board.activeTab].cards[board.activeCard] && board.tabs[board.activeTab].cards[board.activeCard].id"
                :board="board.current" :core="board" :bus="bus"
                :card="board.tabs[board.activeTab].cards[board.activeCard]"/>
            </el-tab-pane>
          </el-tabs>

        </div>
      </template>
    </split-pane>
  </div>

</template>

<script lang="ts">
import {Component, Prop, Vue} from 'vue-property-decorator';
import EditorTabMuu from './editor/tab-muu.vue';
import stream from '@/api/stream';
import {UUID} from 'uuid-generator-ts';
import Editor from '@/views/automation/new.vue';
import api from '@/api/api';
import {Core} from '@/views/dashboard/core';
import CardWrapper from '@/components/card-wrapper/index.vue';
import {EventStateChange} from '@/api/stream_types';
import SplitPane from 'vue-splitpane';
import TabDashboardSettings from "@/views/dashboard/editor/tab-dashboard-settings.vue";
import TabEditor from "@/views/dashboard/editor/tab-editor.vue";
import CardEditor from "@/views/dashboard/editor/card-editor.vue";
import CardItemsEditor from "@/views/dashboard/editor/card-items-editor.vue";

@Component({
  name: 'DashboardEditor',
  components: {
    CardItemsEditor,
    CardEditor,
    TabEditor,
    TabDashboardSettings,
    CardWrapper,
    Editor,
    EditorTabMuu,
    SplitPane
  }
})
export default class extends Vue {
  @Prop({required: true}) private id!: number;

  private loading = true;
  private bus: Vue = new Vue();

  private board: Core = new Core(this.bus);

  // id for streaming subscribe
  private currentID = '';

  private mounted() {

  }

  created() {
    const uuid = new UUID();
    this.currentID = uuid.getDashFreeUUID();

    this.fetchDashboard();

    // setTimeout(() => {
    stream.subscribe('state_changed', this.currentID, this.onStateChanged);

    // for (const entityId in this.board.current.entities) {
    //   requestCurrentState(entityId);
    // }
    // }, 1000);
  }

  private destroyed() {
    stream.unsubscribe('state_changed', this.currentID);
  }

  resize() {
    // Handle resize
  }

  private onStateChanged(event: EventStateChange) {
    this.bus.$emit('state_changed', event);
    this.board.onStateChanged(event);
  }

  // ---------------------------------
  // dashboard
  // ---------------------------------

  private async fetchDashboard() {
    this.loading = true;
    const {data} = await api.v1.dashboardServiceGetDashboardById(this.id);
    this.board.currentBoard(data);
    this.loading = false;
  }

  // ---------------------------------
  // tabs
  // ---------------------------------

  private handleTabsEdit(targetName: string, action: string) {
    console.log('targetName', targetName, 'action', action);
    switch (action) {
      case 'add':
        this.createTab();
        break;
      case 'remove':
    }
  }

  private updateCurrentTab(tab: any) {
    this.board.updateCurrentTab();
  }

  private async createTab() {
    await this.board.createTab();

    this.$notify({
      title: 'Success',
      message: 'tab created successfully',
      type: 'success',
      duration: 2000
    });
  }
}
</script>

<style lang="scss">

.splitter-pane-resizer.horizontal[data-v-212fa2a4] {
  border-top: 5px solid hsl(0deg 0% 100%);
  border-bottom: 5px solid hsl(0deg 0% 100%);
}


/* Track */
::-webkit-scrollbar-track {
  background: #f1f1f1;
}

.dashboard-container {
  padding-top: 20px;
  position: relative;
}

.item-card-editor {
  /*background-color: #f1f1f1;*/
}

.components-container {
  height: calc(100vh - 50px);
  /*min-height: calc(100vh - 50px);*/
  margin: 0;
  padding: 0;
}

.top-container {
  width: 100%;
  height: 100%;
  padding: 0 20px;
  overflow-y: scroll;
}

.bottom-container {
  width: 100%;
  height: 100%;
  padding: 0 20px;
  overflow-y: scroll;
}

.filter-list {

}
</style>
