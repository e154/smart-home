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
              <el-row :gutter="20">
                <el-col :span="12" :xs="12">
                  <el-card>
                    <div slot="header" class="clearfix">
                      <span>Main settings</span>
                    </div>
                    <el-form label-position="top"
                             :model="board.current"
                             style="width: 100%"
                             :rules="boardRules"
                             ref="tabForm">

                      <el-form-item :label="$t('dashboard.editor.name')" prop="name">
                        <el-input size="small" v-model="board.current.name"></el-input>
                      </el-form-item>

                      <el-form-item :label="$t('dashboard.editor.description')" prop="description">
                        <el-input size="small" v-model="board.current.description"></el-input>
                      </el-form-item>

                      <el-form-item :label="$t('dashboard.editor.enabled')" prop="enabled">
                        <el-switch v-model="board.current.enabled"></el-switch>
                      </el-form-item>

                      <el-form-item :label="$t('entities.table.area')" prop="area">
                        <area-search
                          :multiple=false
                          v-model="board.current.area"
                          @update-value="changedArea"
                        />
                      </el-form-item>

                    </el-form>

                    <div style="text-align: right">
                      <export-tool
                        :title="$t('main.export')"
                        :visible="showExport"
                        :value="exportValue"
                        @on-close="showExport=false"/>
                      <el-button type="primary" icon="el-icon-document" @click.prevent.stop='_exportDashbord'>{{
                          $t('main.export')
                        }}
                      </el-button>

                      <el-button type="primary" @click.prevent.stop="updateBoard">{{ $t('main.update') }}</el-button>

                      <el-button @click.prevent.stop="fetchDashboard">{{ $t('main.load_from_server') }}</el-button>

                      <el-popconfirm
                        :confirm-button-text="$t('main.ok')"
                        :cancel-button-text="$t('main.no')"
                        icon="el-icon-info"
                        icon-color="red"
                        style="margin-left: 10px;"
                        :title="$t('main.are_you_sure_to_do_want_this?')"
                        v-on:confirm="cancel()"
                      >
                        <el-button type="default" slot="reference">{{ $t('main.cancel') }}</el-button>
                      </el-popconfirm>

                      <el-popconfirm
                        :confirm-button-text="$t('main.ok')"
                        :cancel-button-text="$t('main.no')"
                        icon="el-icon-info"
                        icon-color="red"
                        style="margin-left: 10px;"
                        :title="$t('main.are_you_sure_to_do_want_this?')"
                        v-on:confirm="removeBoard"
                      >
                        <el-button type="danger" icon="el-icon-delete" slot="reference">{{
                            $t('main.remove')
                          }}
                        </el-button>
                      </el-popconfirm>

                    </div>
                  </el-card>
                </el-col>

              </el-row>
            </el-tab-pane>
            <el-tab-pane :label="$t('dashboard.editor.tabList')" key="1">
              <el-row :gutter="20">
                <el-col :span="12" :xs="12">
                  <el-card>
                    <div slot="header" class="clearfix">
                      <span>{{ $t('dashboard.editor.tabSettings') }}</span>
                    </div>

                    <div style="padding: 0 0 20px 0">
                      <el-button type="default" @click.prevent.stop="createTab"><i
                        class="el-icon-plus"/> {{ $t('dashboard.editor.addTab') }}
                      </el-button>
                    </div>

                    <div v-if="board.activeTab && board.tabs[board.activeTab]">

                      <el-form label-position="top"
                               :model="board.tabs[board.activeTab]"
                               style="width: 100%"
                               :rules="tabRules"
                               ref="tabForm">
                        <el-form-item :label="$t('dashboard.editor.name')" prop="name">
                          <el-input size="small" v-model="board.tabs[board.activeTab].name"></el-input>
                        </el-form-item>
                        <el-form-item :label="$t('dashboard.editor.icon')" prop="icon">
                          <el-input size="small" v-model="board.tabs[board.activeTab].icon"></el-input>
                        </el-form-item>
                        <el-form-item :label="$t('dashboard.editor.columnWidth')" prop="columnWidth">
                          <el-input-number size="small" v-model="board.tabs[board.activeTab].columnWidth" :min="50"
                                           :max="1024"></el-input-number>
                        </el-form-item>
                        <el-form-item :label="$t('dashboard.editor.gap')" prop="gap">
                          <el-switch v-model="board.tabs[board.activeTab].gap"></el-switch>
                        </el-form-item>
                        <el-form-item :label="$t('dashboard.editor.background')" prop="background">
                          <el-color-picker v-model="board.tabs[board.activeTab].background"></el-color-picker>
                        </el-form-item>
                        <el-form-item :label="$t('dashboard.editor.enabled')" prop="enabled">
                          <el-switch v-model="board.tabs[board.activeTab].enabled"></el-switch>
                        </el-form-item>
                        <el-form-item :label="$t('dashboard.editor.dragEnabled')" prop="dragEnabled">
                          <el-switch v-model="board.tabs[board.activeTab].dragEnabled"></el-switch>
                        </el-form-item>
                      </el-form>
                      <div style="text-align: right">
                        <el-button type="primary" @click.prevent.stop="updateTab">{{ $t('main.update') }}</el-button>
                        <el-popconfirm
                          :confirm-button-text="$t('main.ok')"
                          :cancel-button-text="$t('main.no')"
                          icon="el-icon-info"
                          icon-color="red"
                          style="margin-left: 10px;"
                          :title="$t('main.are_you_sure_to_do_want_this?')"
                          v-on:confirm="cancel()"
                        >
                          <el-button type="default" slot="reference">{{ $t('main.cancel') }}</el-button>
                        </el-popconfirm>
                        <el-popconfirm
                          :confirm-button-text="$t('main.ok')"
                          :cancel-button-text="$t('main.no')"
                          icon="el-icon-info"
                          icon-color="red"
                          style="margin-left: 10px;"
                          :title="$t('main.are_you_sure_to_do_want_this?')"
                          v-on:confirm="removeTab"
                        >
                          <el-button type="danger" icon="el-icon-delete" slot="reference">{{
                              $t('main.remove')
                            }}
                          </el-button>
                        </el-popconfirm>
                      </div>
                    </div>

                  </el-card>
                </el-col>
                <el-col :span="6" :xs="12">
                  <el-card>
                    <div slot="header" class="clearfix">
                      <span>{{ $t('dashboard.editor.tabList') }}</span>
                    </div>

                    <el-menu
                      v-if="board.tabs"
                      ref="tabMenu"
                      :default-active="board.activeTab"
                      v-model="board.activeTab"
                      class="el-menu-vertical-demo"
                    >
                      <el-menu-item :index="index + ''" v-for="(tab, index) in board.tabs"
                                    @click="menuTabClick(index, tab)">
                        <span>{{ tab.name }}</span>
                      </el-menu-item>
                    </el-menu>

                  </el-card>
                </el-col>
              </el-row>
            </el-tab-pane>
            <el-tab-pane :label="$t('dashboard.editor.cardList')" key="2">
              <el-row :gutter="20" v-if="board.activeTab && board.tabs[board.activeTab]">
                <el-col :span="12" :xs="12">
                  <el-card>
                    <div slot="header" class="clearfix">
                      <span>{{ $t('dashboard.editor.cardSettings') }}</span>
                    </div>

                    <div style="padding-bottom: 20px" v-if="board.tabs[board.activeTab].id">
                      <el-button type="default" @click.prevent.stop="addCard"><i class="el-icon-plus"/>
                        {{ $t('dashboard.editor.addNewCard') }}
                      </el-button>

                      <el-button type="default" @click.prevent.stop="showCardImport = true">{{
                          $t('main.import')
                        }}
                      </el-button>

                      <export-tool
                        :title="$t('main.import')"
                        :visible="showCardImport"
                        :value="importCardValue"
                        @on-close="showCardImport=false"
                        @on-import="importCard"
                        :import-dialog="true"/>

                    </div>

                    <div
                      v-if="board.tabs[board.activeTab].cards[board.activeCard] && board.tabs[board.activeTab].cards[board.activeCard].id">

                      <el-form label-position="top"
                               :model="board.tabs[board.activeTab].cards[board.activeCard]"
                               style="width: 100%"
                               :rules="cardRules"
                               ref="cardForm">
                        <el-form-item :label="$t('dashboard.editor.title')" prop="title">
                          <el-input size="small"
                                    v-model="board.tabs[board.activeTab].cards[board.activeCard].title"></el-input>
                        </el-form-item>

                        <el-form-item :label="$t('dashboard.editor.height')" prop="height">
                          <el-input-number size="small"
                                           v-model="board.tabs[board.activeTab].cards[board.activeCard].height"></el-input-number>
                        </el-form-item>

                        <el-form-item :label="$t('dashboard.editor.width')" prop="width">
                          <el-input-number size="small"
                                           v-model="board.tabs[board.activeTab].cards[board.activeCard].width"></el-input-number>
                        </el-form-item>

                        <el-form-item :label="$t('dashboard.editor.background')" prop="background">
                          <el-color-picker
                            v-model="board.tabs[board.activeTab].cards[board.activeCard].background"></el-color-picker>
                        </el-form-item>
                        <el-form-item :label="$t('dashboard.editor.enabled')" prop="enabled">
                          <el-switch v-model="board.tabs[board.activeTab].cards[board.activeCard].enabled"></el-switch>
                        </el-form-item>
                      </el-form>

                    </div>

                    <el-row
                      v-if="board.tabs[board.activeTab].cards[board.activeCard] && board.tabs[board.activeTab].cards[board.activeCard].id">
                      <el-col style="text-align: right">

                        <export-tool
                          :title="$t('main.export')"
                          :visible="showExportCard"
                          :value="exportCardValue"
                          @on-close="showExportCard=false"/>
                        <el-button type="primary" icon="el-icon-document" @click.prevent.stop='exportCard'>{{
                            $t('main.export')
                          }}
                        </el-button>

                        <el-button type="primary" @click.prevent.stop="copyCard">{{ $t('main.copy') }}</el-button>

                        <el-button type="primary" @click.prevent.stop="updateCard">{{ $t('main.update') }}</el-button>

                        <el-popconfirm
                          :confirm-button-text="$t('main.ok')"
                          :cancel-button-text="$t('main.no')"
                          icon="el-icon-info"
                          icon-color="red"
                          style="margin-left: 10px;"
                          :title="$t('main.are_you_sure_to_do_want_this?')"
                          v-on:confirm="cancel()"
                        >
                          <el-button type="default" slot="reference">{{ $t('main.cancel') }}</el-button>
                        </el-popconfirm>

                        <el-popconfirm
                          :confirm-button-text="$t('main.ok')"
                          :cancel-button-text="$t('main.no')"
                          icon="el-icon-info"
                          icon-color="red"
                          style="margin-left: 10px;"
                          :title="$t('main.are_you_sure_to_do_want_this?')"
                          v-on:confirm="removeCard"
                        >
                          <el-button type="danger" icon="el-icon-delete" slot="reference">{{
                              $t('main.remove')
                            }}
                          </el-button>
                        </el-popconfirm>
                      </el-col>
                    </el-row>

                  </el-card>
                </el-col>
                <el-col :span="8" :xs="12">

                  <el-card>
                    <div slot="header" class="clearfix">
                      <span>{{ $t('dashboard.editor.cardList') }}</span>
                    </div>

                    <el-menu
                      v-if="board.activeTab && board.tabs[board.activeTab].cards"
                      ref="tabMenu"
                      :default-active="board.activeCard + ''"
                      v-model="board.activeCard"
                      class="el-menu-vertical-demo"
                    >
                      <el-menu-item :index="index + ''" v-for="(card, index) in board.tabs[board.activeTab].cards"
                                    @click="menuCardsClick(card)">
                        <div style="clear: both">
                          <span>{{ card.title }}</span>
                          <span>
                            <el-button style="float: right" size="mini" @click.prevent.stop="sortCardUp(card, index)"><i
                              class="el-icon-upload2"/></el-button>
                          <el-button style="float: right" size="mini" @click.prevent.stop="sortCardDown(card, index)"><i
                            class="el-icon-download"/></el-button>
                          </span>
                        </div>

                      </el-menu-item>
                    </el-menu>

                  </el-card>
                </el-col>
              </el-row>
            </el-tab-pane>
            <el-tab-pane :label="$t('dashboard.editor.cardItems')" key="3">
              <el-row :gutter="20"
                      v-if="board.activeTab && board.tabs[board.activeTab] && board.tabs[board.activeTab].cards[board.activeCard] && board.tabs[board.activeTab].cards[board.activeCard].id">
                <el-col :span="15" :xs="15">

                  <el-card>
                    <div slot="header" class="clearfix">
                      <span>{{ $t('dashboard.editor.itemDetail') }}</span>
                    </div>

                    <div style="padding: 0 0 20px 0" v-if="board.activeCard !== undefined">
                      <el-button type="default" @click.prevent.stop="addCardItem"><i
                        class="el-icon-plus"/>{{ $t('dashboard.editor.addCardItem') }}
                      </el-button>
                    </div>

                    <div
                      v-if="board.tabs[board.activeTab].cards[board.activeCard] && board.tabs[board.activeTab].cards[board.activeCard].items && board.tabs[board.activeTab].cards[board.activeCard].selectedItem > -1"
                      style="padding: 0 0 20px;"
                    >

                      <el-card shadow="never" class="item-card-editor">

                        <el-form label-position="top"
                                 :model="board.tabs[board.activeTab].cards[board.activeCard].items[board.tabs[board.activeTab].cards[board.activeCard].selectedItem]"
                                 style="width: 100%"
                                 ref="cardItemForm">

                          <el-row :gutter="20">
                            <el-col
                              :span="8"
                              :xs="8"
                            >
                              <el-form-item :label="$t('dashboard.editor.type')" prop="type">
                                <el-select
                                  v-model="board.tabs[board.activeTab].cards[board.activeCard].items[board.tabs[board.activeTab].cards[board.activeCard].selectedItem].type"
                                  :placeholder="$t('dashboard.editor.pleaseSelectType')"
                                  style="width: 100%"
                                >
                                  <el-option
                                    v-for="item in itemTypes"
                                    :key="item.value"
                                    :label="item.label"
                                    :value="item.value">
                                  </el-option>
                                </el-select>
                              </el-form-item>

                            </el-col>
                            <el-col
                              :span="8"
                              :xs="8"
                            >
                              <el-form-item :label="$t('dashboard.editor.title')" prop="title">
                                <el-input size="small"
                                          v-model="board.tabs[board.activeTab].cards[board.activeCard].items[board.tabs[board.activeTab].cards[board.activeCard].selectedItem].title"></el-input>
                              </el-form-item>
                            </el-col>
                          </el-row>

                          <component
                            :is="getCardEditorName(board.tabs[board.activeTab].cards[board.activeCard].items[board.tabs[board.activeTab].cards[board.activeCard].selectedItem].type)"
                            :board="board"
                            :item="board.tabs[board.activeTab].cards[board.activeCard].items[board.tabs[board.activeTab].cards[board.activeCard].selectedItem]"
                            :index="board.tabs[board.activeTab].cards[board.activeCard].selectedItem"
                          />

                        </el-form>

                        <div style="text-align: right;">

                          <el-button type="primary" @click.prevent.stop="updateCard">{{ $t('main.update') }}</el-button>
                          <el-button
                            @click.prevent.stop="copyCardItem(board.tabs[board.activeTab].cards[board.activeCard].selectedItem)">
                            {{ $t('main.copy') }}
                          </el-button>
                          <el-popconfirm
                            :confirm-button-text="$t('main.ok')"
                            :cancel-button-text="$t('main.no')"
                            icon="el-icon-info"
                            icon-color="red"
                            style="margin-left: 10px;"
                            :title="$t('main.are_you_sure_to_do_want_this?')"
                            v-on:confirm="cancel()"
                          >
                            <el-button type="default" slot="reference">{{ $t('main.cancel') }}</el-button>
                          </el-popconfirm>
                          <el-popconfirm
                            :confirm-button-text="$t('main.ok')"
                            :cancel-button-text="$t('main.no')"
                            icon="el-icon-info"
                            icon-color="red"
                            style="margin-left: 10px;"
                            :title="$t('main.are_you_sure_to_do_want_this?')"
                            v-on:confirm="removeCardItem(board.tabs[board.activeTab].cards[board.activeCard].selectedItem)"
                          >
                            <el-button type="danger" icon="el-icon-delete" slot="reference">{{
                                $t('main.remove')
                              }}
                            </el-button>
                          </el-popconfirm>
                        </div>

                      </el-card>

                    </div>

                  </el-card>
                </el-col>

                <el-col :span="8" :xs="12">
                  <el-card>
                    <div slot="header" class="clearfix">
                      <span>{{ $t('dashboard.editor.itemList') }}</span>
                    </div>

                    <el-menu
                      v-if="board.tabs[board.activeTab].cards[board.activeCard] && board.tabs[board.activeTab].cards[board.activeCard].id"
                      ref="tabMenu"
                      :default-active="board.tabs[board.activeTab].cards[board.activeCard].selectedItem + ''"
                      v-model="board.tabs[board.activeTab].cards[board.activeCard].selectedItem"
                      class="el-menu-vertical-demo"
                    >
                      <el-menu-item :index="index + ''"
                                    v-for="(item, index) in board.tabs[board.activeTab].cards[board.activeCard].items"
                                    @click="menuCardItemClick(index)">
                    <span>
                      {{ item.title }}
                      <el-tag size="mini">{{ item.type }}</el-tag>
                      <el-tag v-if="item.hidden" size="mini" type="info" effect="plain">
                        {{ $t('dashboard.editor.hidden') }}
                      </el-tag>
                      <el-tag v-if="!item.enabled" size="mini" type="info" effect="plain">
                        {{ $t('dashboard.editor.disabled') }}
                      </el-tag>
                      <el-tag v-if="item.frozen" size="mini" type="info" effect="plain">
                        {{ $t('dashboard.editor.frozen') }}
                      </el-tag>
                    </span>
                      </el-menu-item>
                    </el-menu>

                  </el-card>


                  <!-- TODO: fix -->
                  <table style="margin-top: 20px" class="filter-list">
                    <thead style="background: #d7d7d7">
                    <tr>
                      <td>name</td>
                      <td>description</td>
                    </tr>
                    </thead>
                    <tbody>
                    <tr v-for="(filter, index) in getFilterList()">
                      <td><strong>{{ filter.name }}</strong></td>
                      <td>{{ filter.description }}</td>
                    </tr>
                    </tbody>
                  </table>

                </el-col>
              </el-row>
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
import {Form} from 'element-ui';
import {Card, Core, Tab} from '@/views/dashboard/core';
import ImagePreview from '@/views/images/preview.vue';
import CardWrapper from '@/components/card-wrapper/index.vue';
import EntitySearch from '@/views/entities/components/entity_search.vue';
import {
  CardEditorName,
  CardItemList,
  IButtonEditor,
  IChartEditor,
  IImageEditor,
  ILogsEditor,
  IProgressEditor,
  IStateEditor,
  ITextEditor
} from '@/views/dashboard/card_items';
import {EventStateChange} from '@/api/stream_types';
import ExportTool from '@/components/export-tool/index.vue';
import {ApiArea} from '@/api/stub';
import AreaSearch from '@/views/areas/components/areas_search.vue';
import SplitPane from 'vue-splitpane';
import {filterInfo, filterList} from '@/views/dashboard/filters';

@Component({
  name: 'DashboardEditor',
  components: {
    AreaSearch,
    ExportTool,
    EntitySearch,
    CardWrapper,
    Editor,
    EditorTabMuu,
    ImagePreview,
    IButtonEditor,
    ITextEditor,
    IImageEditor,
    IStateEditor,
    ILogsEditor,
    IProgressEditor,
    IChartEditor,
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

  private itemTypes = CardItemList;

  private tabRules = {
    name: [
      {required: true, trigger: 'blur'},
      {min: 4, max: 255, trigger: 'blur'}
    ]
  };

  private cardRules = {
    title: [
      {required: true, trigger: 'blur'},
      {min: 4, max: 255, trigger: 'blur'}
    ]
  };

  private boardRules = {
    name: [
      {required: true, trigger: 'blur'},
      {min: 4, max: 255, trigger: 'blur'}
    ],
    description: [
      {required: false, trigger: 'blur'},
      {min: 0, max: 255, trigger: 'blur'}
    ]
  };

  private mounted() {

  }

  created() {
    const uuid = new UUID();
    this.currentID = uuid.getDashFreeUUID();

    this.fetchDashboard();

    this.bus.$on('selected_card', (m: any) => {
      this.onSelectedCard(m);
    });

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

  private getCardEditorName(name: string): string {
    return CardEditorName(name);
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

  private showExport = false;
  private exportValue = {};

  private _exportDashbord() {
    this.exportValue = this.board.serialize();
    this.showExport = true;
  }

  private cancel() {
    this.$router.go(-1);
  }

  private async updateBoard() {
    const {data} = await this.board.update();
    if (data) {
      this.$notify({
        title: 'Success',
        message: 'dashboard updated successfully',
        type: 'success',
        duration: 2000
      });
    }
  }

  private changedArea(values: ApiArea, event?: any) {
    if (values) {
      this.$set(this.board.current, 'area', values);
      this.$set(this.board.current, 'areaId', values.id);
    } else {
      this.$set(this.board.current, 'area', undefined);
      this.$set(this.board.current, 'areaId', undefined);
    }
  }

  private async removeBoard() {
    await this.board.removeBoard();
    this.$notify({
      title: 'Success',
      message: 'dashboard deleted successfully',
      type: 'success',
      duration: 2000
    });

    this.$router.go(-1);
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

  private async updateTab() {
    if (!this.board.activeTab || !this.board.tabs[this.board.activeTab]) {
      return;
    }

    (this.$refs.tabForm as Form).validate(async valid => {
      if (!valid) {
        return;
      }

      const {data} = await this.board.updateTab();
      if (data) {
        this.$notify({
          title: 'Success',
          message: 'tab updated successfully',
          type: 'success',
          duration: 2000
        });
      }
    });
  }

  private async removeTab() {
    if (!this.board.tabs[this.board.activeTab]) {
      return;
    }
    await this.board.removeTab();

    await this.fetchDashboard();
  }

  private menuTabClick(index: string, tab: Tab) {
    this.board.activeTab = index + '';
    this.board.updateCurrentTab();
  }

  // ---------------------------------
  // cards
  // ---------------------------------
  private onSelectedCard(id: number) {
    this.board.onSelectedCard(id);
  }

  private async addCard() {
    return this.board.createCard();
  }

  private async updateCard() {
    (this.$refs.cardForm as Form).validate(async valid => {
      if (!valid) {
        return;
      }

      const {data} = await this.board.updateCard();

      if (data) {
        this.$notify({
          title: 'Success',
          message: 'card updated successfully',
          type: 'success',
          duration: 2000
        });
      }
    });
  }

  private async removeCard() {
    await this.board.removeCard();

    this.$notify({
      title: 'Success',
      message: 'card removed successfully',
      type: 'success',
      duration: 2000
    });
  }

  private menuCardsClick(card: Card) {
    this.bus.$emit('selected_card', card.id);
  }

  // export card
  private showExportCard = false;
  private exportCardValue = {};

  private exportCard() {
    if (this.board.activeTab == undefined || this.board.activeCard == undefined) {
      return;
    }
    this.exportCardValue = this.board.tabs[this.board.activeTab].cards[this.board.activeCard].serialize();
    this.showExportCard = true;
  }

  // import card
  private showCardImport = false;
  private importCardValue = '';

  private async importCard(value: string, event?: any) {
    const card = JSON.parse(value);
    const data = await this.board.importCard(card);
    if (data) {
      this.$notify({
        title: 'Success',
        message: 'card imported successfully',
        type: 'success',
        duration: 2000
      });
    }
  }

  private async copyCard() {
    if (this.board.activeTab == undefined || this.board.activeCard == undefined) {
      return;
    }

    const card = this.board.tabs[this.board.activeTab].cards[this.board.activeCard].serialize();

    const data = await this.board.importCard(card);
    if (data) {
      this.$notify({
        title: 'Success',
        message: 'card copied successfully',
        type: 'success',
        duration: 2000
      });
    }
  }

  private sortCardUp(card: Card, index: number) {
    if (!this.board.tabs || !this.board.activeTab) {
      return;
    }

    if (!this.board.tabs[this.board.activeTab].cards[index - 1]) {
      return;
    }

    let weight1 = this.board.tabs[this.board.activeTab].cards[index - 1].weight;
    const weight2 = card.weight;
    if (weight1 === weight2) {
      weight1 += 1;
    }
    this.board.tabs[this.board.activeTab].cards[index - 1].weight = weight2;
    this.board.tabs[this.board.activeTab].cards[index].weight = weight1;
    this.board.tabs[this.board.activeTab].cards[index].update();
    this.board.tabs[this.board.activeTab].cards[index - 1].update();

    this.board.tabs[this.board.activeTab].sortCards();
    this.board.updateCurrentTab();
  }

  private sortCardDown(card: Card, index: number) {
    if (!this.board.tabs || !this.board.activeTab) {
      return;
    }

    if (!this.board.tabs[this.board.activeTab].cards[index + 1]) {
      return;
    }

    const wd = this.board.tabs[this.board.activeTab].cards[index + 1].weight;
    this.board.tabs[this.board.activeTab].cards[index + 1].weight = card.weight;
    this.board.tabs[this.board.activeTab].cards[index].weight = wd;
    this.board.tabs[this.board.activeTab].cards[index].update();
    this.board.tabs[this.board.activeTab].cards[index + 1].update();

    this.board.tabs[this.board.activeTab].sortCards();
    this.board.updateCurrentTab();
  }

  // ---------------------------------
  // card items
  // ---------------------------------
  private addCardItem() {
    this.board.createCardItem();
  }

  private removeCardItem(index: number) {
    this.board.removeCardItem(index);
  }

  private itemPosition(): string {
    const defaultValue = 'L:0,T:0,W:0,H:0';

    if (!this.board.activeTab || this.board.activeCard == undefined) {
      return defaultValue;
    }

    if (!this.board.tabs || !this.board.tabs[this.board.activeTab] || !this.board.tabs[this.board.activeTab].cards) {
      return defaultValue;
    }

    const card = this.board.tabs[this.board.activeTab].cards[this.board.activeCard];
    if (!card) {
      return defaultValue;
    }

    if (card.selectedItem == -1 || card.selectedItem == undefined) {
      return defaultValue;
    }

    const currentItem = card.items[card.selectedItem];
    if (!currentItem) {
      return defaultValue;
    }
    const info = currentItem.positionInfo;
    return `L:${info.left},T:${info.top},W:${info.width},H:${info.height}`;
  }

  private copyCardItem(index: number) {
    if (!this.board.activeTab || this.board.activeCard == undefined) {
      return;
    }

    const card = this.board.tabs[this.board.activeTab].cards[this.board.activeCard];
    if (card.selectedItem == -1 || card.selectedItem == undefined) {
      return;
    }

    card.copyItem(index);
  }

  private menuCardItemClick(index: number) {
    if (!this.board.activeTab || this.board.activeCard == undefined) {
      return;
    }

    this.board.tabs[this.board.activeTab].cards[this.board.activeCard].selectedItem = index;
  }

  private getFilterList(): filterInfo[] {
    return filterList();
  }
}
</script>

<style lang="scss">

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
  overflow: scroll;
}

.bottom-container {
  width: 100%;
  height: 100%;
  padding: 0 20px;
  overflow: scroll;
}

.filter-list {

}
</style>
