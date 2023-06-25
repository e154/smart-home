<template>

  <el-row :gutter="20">
    <el-col :span="15" :xs="15">
      <el-card>
        <div slot="header" class="clearfix">
          <span>Main settings</span>
        </div>
        <el-form label-position="top"
                 :model="board"
                 style="width: 100%"
                 :rules="boardRules"
                 ref="tabForm">

          <el-form-item :label="$t('dashboard.editor.name')" prop="name">
            <el-input size="small" v-model="board.name"></el-input>
          </el-form-item>

          <el-form-item :label="$t('dashboard.editor.description')" prop="description">
            <el-input size="small" v-model="board.description"></el-input>
          </el-form-item>

          <el-form-item :label="$t('dashboard.editor.enabled')" prop="enabled">
            <el-switch v-model="board.enabled"></el-switch>
          </el-form-item>

          <el-form-item :label="$t('entities.table.area')" prop="area">
            <area-search
              :multiple=false
              v-model="board.area"
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

class elementOption {
  public value = '';
  public label = '';
}

@Component({
  name: 'TabDashboardSettings',
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
  private showExport = false;
  private exportValue = {};

  private _exportDashbord() {
    this.exportValue = this.core.serialize();
    this.showExport = true;
  }

  private async updateBoard() {
    const {data} = await this.core.update();
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
      this.$set(this.board, 'area', values);
      this.$set(this.board, 'areaId', values.id);
    } else {
      this.$set(this.board, 'area', undefined);
      this.$set(this.board, 'areaId', undefined);
    }
  }

  private async removeBoard() {
    await this.core.removeBoard();
    this.$notify({
      title: 'Success',
      message: 'dashboard deleted successfully',
      type: 'success',
      duration: 2000
    });

    this.$router.go(-1);
  }

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

  private cancel() {
    this.$router.go(-1);
  }
}
</script>

<style lang="scss">



</style>
