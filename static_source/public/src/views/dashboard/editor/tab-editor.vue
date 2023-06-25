<template>
  <el-row :gutter="20">
    <el-col :span="15" :xs="15">
      <el-card>
        <div slot="header" class="clearfix">
          <span>{{ $t('dashboard.editor.tabSettings') }}</span>
        </div>

        <div style="padding: 0 0 20px 0">
          <el-button type="default" @click.prevent.stop="createTab"><i
            class="el-icon-plus"/> {{ $t('dashboard.editor.addTab') }}
          </el-button>
        </div>

        <div v-if="tab">

          <el-form label-position="top"
                   :model="tab"
                   style="width: 100%"
                   :rules="tabRules"
                   ref="tabForm">
            <el-form-item :label="$t('dashboard.editor.name')" prop="name">
              <el-input size="small" v-model="tab.name"></el-input>
            </el-form-item>
            <el-form-item :label="$t('dashboard.editor.icon')" prop="icon">
              <el-input size="small" v-model="tab.icon"></el-input>
            </el-form-item>
            <el-form-item :label="$t('dashboard.editor.columnWidth')" prop="columnWidth">
              <el-input-number size="small" v-model="tab.columnWidth" :min="50"
                               :max="1024"></el-input-number>
            </el-form-item>
            <el-form-item :label="$t('dashboard.editor.gap')" prop="gap">
              <el-switch v-model="tab.gap"></el-switch>
            </el-form-item>
            <el-form-item :label="$t('dashboard.editor.background')" prop="background">
              <el-color-picker v-model="tab.background"></el-color-picker>
            </el-form-item>
            <el-form-item :label="$t('dashboard.editor.enabled')" prop="enabled">
              <el-switch v-model="tab.enabled"></el-switch>
            </el-form-item>
            <el-form-item :label="$t('dashboard.editor.dragEnabled')" prop="dragEnabled">
              <el-switch v-model="tab.dragEnabled"></el-switch>
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
    <el-col :span="8" :xs="12">
      <el-card>
        <div slot="header" class="clearfix">
          <span>{{ $t('dashboard.editor.tabList') }}</span>
        </div>

        <el-menu
          v-if="board.tabs"
          ref="tabMenu"
          :default-active="core.activeTab"
          v-model="core.activeTab"
          class="el-menu-vertical-demo"
        >
          <el-menu-item :index="index + ''" v-for="(t, index) in board.tabs"
                        @click="menuTabClick(index, t)">
            <span>{{ t.name }}</span>
          </el-menu-item>
        </el-menu>

      </el-card>
    </el-col>
  </el-row>
</template>

<script lang="ts">
import {Component, Prop, Vue} from 'vue-property-decorator';
import {
  CardItemName,
  Dummy,
  IButton,
  IChart,
  IImage,
  ILogs,
  IProgress,
  IState,
  IText
} from '@/views/dashboard/card_items';
import {Card, Core, Tab} from '@/views/dashboard/core';
import Moveable from 'vue-moveable';
import {ApiArea, ApiDashboard} from "@/api/stub";
import {UUID} from "uuid-generator-ts";
import AreaSearch from "@/views/areas/components/areas_search.vue";
import ExportTool from "@/components/export-tool/index.vue";
import {Form} from "element-ui";

class elementOption {
  public value = '';
  public label = '';
}

@Component({
  name: 'TabEditor',
  components: {
    ExportTool,
    AreaSearch,
    Moveable,
    Dummy,
    IText,
    IImage,
    IButton,
    IState,
    ILogs,
    IProgress,
    IChart
  }
})
export default class extends Vue {
  @Prop() private tab!: Tab;
  @Prop() private board!: ApiDashboard;
  @Prop() private core!: Core;
  @Prop() private bus!: Vue;

  // id for streaming subscribe
  private currentID = '';

  private mounted() {
  }

  private destroyed() {

  }

  created() {
    const uuid = new UUID();
    this.currentID = uuid.getDashFreeUUID();
  }

  // ---------------------------------
  // common
  // ---------------------------------
  private async updateTab() {
    if (!this.core.activeTab || !this.board.tabs[this.core.activeTab]) {
      return;
    }

    (this.$refs.tabForm as Form).validate(async valid => {
      if (!valid) {
        return;
      }

      const {data} = await this.core.updateTab();
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
    if (!this.board.tabs[this.core.activeTab]) {
      return;
    }
    await this.core.removeTab();

    //todo uncomment
    // await this.fetchDashboard();
    this.$emit('update-value')
  }

  private menuTabClick(index: string, tab: Tab) {
    this.core.activeTab = index + '';
    this.core.updateCurrentTab();
  }

  private async createTab() {
    await this.core.createTab();

    this.$notify({
      title: 'Success',
      message: 'tab created successfully',
      type: 'success',
      duration: 2000
    });
  }

  private tabRules = {
    name: [
      {required: true, trigger: 'blur'},
      {min: 4, max: 255, trigger: 'blur'}
    ]
  };

  private cancel() {
    this.$router.go(-1);
  }
}
</script>

<style lang="scss">



</style>
