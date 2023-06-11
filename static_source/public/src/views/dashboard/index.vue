<template>
  <div class="app-container">
    <el-row>
      <el-col>
        <el-button type="primary" @click.prevent.stop="add">
          <i class="el-icon-plus"/> {{ $t('dashboard.addNew') }}
        </el-button>

        <el-button type="primary" @click.prevent.stop="showImport = true">{{ $t('main.import') }}</el-button>

        <export-tool
          :title="$t('main.import')"
          :visible="showImport"
          :value="importValue"
          @on-close="showImport=false"
          @on-import="onImport"
          :import-dialog="true"/>

      </el-col>
    </el-row>
    <el-row>
      <el-col :span="6" :xs="24" :lg="12">
        <el-form label-position="top"
                 ref="currentEntity"
                 :model="settings"
                 style="width: 100%">

          <el-form-item :label="$t('dashboard.mainDashboard')" prop="scripts">
            <dashboard_search
              v-model="settings.mainDashboard"
              @update-value="changedMainDashboard"
            />
          </el-form-item>

          <el-form-item :label="$t('dashboard.devDashboard')" prop="scripts">
            <dashboard_search
              v-model="settings.devDashboard"
              @update-value="changedDevDashboard"
            />
          </el-form-item>

        </el-form>
      </el-col>
    </el-row>
    <el-row>
      <el-col>
        <pagination
          v-show="total>listQuery.limit"
          :total="total"
          :page.sync="listQuery.page"
          :limit.sync="listQuery.limit"
          @pagination="getList"
        />
      </el-col>
    </el-row>
    <el-row>
      <el-col>
        <el-table
          :key="tableKey"
          v-loading="listLoading"
          :data="list"
          style="width: 100%;"
          @sort-change="sortChange"
        >
          <el-table-column
            :label="$t('dashboard.table.id')"
            prop="id"
            sortable="custom"
            align="center"
            width="70"
            :class-name="getSortClass('id')"
          >
            <template slot-scope="{row}">
              <span>{{ row.id }}</span>
            </template>
          </el-table-column>

          <el-table-column
            :label="$t('dashboard.table.name')"
            width="140px"
            align="left"
          >
            <template slot-scope="{row}">
          <span class="cursor-pointer"
                @click="goto(row)">
            {{ row.name }}
          </span>
            </template>
          </el-table-column>

          <el-table-column
            :label="$t('dashboard.table.description')"
            width="auto"
            align="left"
          >
            <template slot-scope="{row}">
              <span>{{ row.description }}</span>
            </template>
          </el-table-column>

          <!--          <el-table-column
                      :label="$t('dashboard.table.enabled')"
                      class-name="status-col"
                      width="150px"
                    >
                      <template slot-scope="{row}">

                        <el-switch
                          v-model="row.enabled"
                          :disabled="row.system"
                          v-on:change="activate(row)">
                        </el-switch>
                      </template>

                    </el-table-column>-->

          <el-table-column
            :label="$t('dashboard.table.operations')"
            width="90px"
            align="right"
          >
            <template slot-scope="{row, $index}">
              <el-button
                type="text" size="small"
                @click='edit(row, $index)'
              >
                {{ $t('dashboard.table.edit') }}
              </el-button>
            </template>
          </el-table-column>

          <el-table-column
            :label="$t('dashboard.table.createdAt')"
            width="160px"
            align="center"
            sortable="custom"
            prop="createdAt"
            :class-name="getSortClass('createdAt')"
          >
            <template slot-scope="{row}">
              <span>{{ row.createdAt | parseTime }}</span>
            </template>
          </el-table-column>

          <el-table-column
            :label="$t('dashboard.table.updatedAt')"
            width="160px"
            align="center"
            sortable="custom"
            prop="updatedAt"
            :class-name="getSortClass('updatedAt')"
          >
            <template slot-scope="{row}">
              <span>{{ row.updatedAt | parseTime }}</span>
            </template>
          </el-table-column>

        </el-table>
      </el-col>
    </el-row>
    <el-row>
      <el-col>
        <pagination
          v-show="total>0"
          :total="total"
          :page.sync="listQuery.page"
          :limit.sync="listQuery.limit"
          @pagination="getList"
        />
      </el-col>
    </el-row>
  </div>

</template>

<script lang="ts">
import {Component, Vue} from 'vue-property-decorator';
import api from '@/api/api';
import Pagination from '@/components/Pagination/index.vue';
import router from '@/router';
import {ApiDashboard, ApiDashboardShort, ApiVariable} from '@/api/stub';
import {Core} from '@/views/dashboard/core';
import ExportTool from '@/components/export-tool/index.vue';
import Scripts from '@/views/entities/components/scripts.vue';
import Dashboard_search from '@/views/dashboard/components/dashboard_search.vue';
import webpack from 'webpack';
import asString = webpack.Template.asString;

@Component({
  name: 'DashboardList',
  components: {
    Dashboard_search,
    Scripts,
    ExportTool,
    Pagination
  }
})
export default class extends Vue {
  private tableKey = 0;
  private list: ApiDashboardShort[] = [];
  private total = 0;
  private listLoading = true;
  private listQuery = {
    page: 1,
    limit: 20,
    sort: '-createdAt'
  };
  private settings: {
    mainDashboard?: ApiDashboard,
    mainVar?: ApiVariable,
    devDashboard?: ApiDashboard
    devVar?: ApiVariable,
  } = {};

  created() {
    this.getList();
    this.getSettings();
  }

  private async getSettings() {
    this.listLoading = true;
    api.v1.variableServiceGetVariableByName('mainDashboard').then((resp) => {
      this.$set(this.settings, 'mainVar', resp.data);
      if (!this.settings?.mainVar?.value) {
        return;
      }
      const id: number = parseInt(this.settings!.mainVar!.value);
      api.v1.dashboardServiceGetDashboardById(id).then((resp) => {
        this.$set(this.settings, 'mainDashboard', resp.data);
      });
    });

    api.v1.variableServiceGetVariableByName('devDashboard').then((resp) => {
      this.$set(this.settings, 'devVar', resp.data);
      if (!this.settings?.devVar?.value) {
        return;
      }
      const id: number = parseInt(this.settings!.devVar!.value);
      api.v1.dashboardServiceGetDashboardById(id).then((resp) => {
        this.$set(this.settings, 'devDashboard', resp.data);
      });
    });
  }

  private updateVariable(name: string, value: string) {
     api.v1.variableServiceUpdateVariable(name, {name: name, value: value})
  }

  private async getList() {
    this.listLoading = true;
    const {data} = await api.v1.dashboardServiceGetDashboardList({
      limit: this.listQuery.limit,
      page: this.listQuery.page,
      sort: this.listQuery.sort
    });

    this.list = data.items;
    this.total = data.meta.total;
    this.listLoading = false;
  }

  private handleFilter() {
    this.listQuery.page = 1;
    this.getList();
  }

  private sortChange(data: any) {
    const {prop, order} = data;
    if (prop === 'id') {
      this.sortByID(order);
    } else if (prop === 'createdAt') {
      this.sortByCreatedAt(order);
    } else if (prop === 'updatedAt') {
      this.sortByUpdatedAt(order);
    }
  }

  private sortByCreatedAt(order: string) {
    if (order === 'ascending') {
      this.listQuery.sort = '+createdAt';
    } else {
      this.listQuery.sort = '-createdAt';
    }
    this.handleFilter();
  }

  private sortByUpdatedAt(order: string) {
    if (order === 'ascending') {
      this.listQuery.sort = '+updatedAt';
    } else {
      this.listQuery.sort = '-updatedAt';
    }
    this.handleFilter();
  }

  private sortByID(order: string) {
    if (order === 'ascending') {
      this.listQuery.sort = '+id';
    } else {
      this.listQuery.sort = '-id';
    }
    this.handleFilter();
  }

  private getSortClass(key: string) {
    const sort = this.listQuery.sort;
    return sort === `+${key}` ? 'ascending' : 'descending';
  }

  private goto(board: ApiDashboardShort) {
    router.push({path: `/dashboards/view/${board.id}`});
  }

  private async add() {
    const counter: number = this.list.length;
    const dashboard = await Core.createNew('new' + counter);
    if (dashboard) {
      this.$notify({
        title: 'Success',
        message: 'dashboard added successfully',
        type: 'success',
        duration: 2000
      });
      router.push({path: `/dashboards/edit/${dashboard.id}`});
    }
  }

  private showImport = false;
  private importValue = '';

  private async onImport(value: string, event?: any) {
    const json = JSON.parse(value);
    const data = await Core._import(json);
    if (data) {
      this.getList();
      this.$notify({
        title: 'Success',
        message: 'dashboard imported successfully',
        type: 'success',
        duration: 2000
      });
    }
  }

  private edit(dashboard: ApiDashboardShort, index: number) {
    router.push({path: `/dashboards/edit/${dashboard.id}`});
  }

  private changedMainDashboard(values: ApiDashboard, event: any) {
    if (values) {
      this.$set(this.settings, 'mainDashboard', values);
      this.updateVariable("mainDashboard", values.id + '')
    } else {
      this.$set(this.settings, 'mainDashboard', undefined);
      this.updateVariable("mainDashboard", "")
    }
  }

  private changedDevDashboard(values: ApiDashboard, event: any) {
    if (values) {
      this.$set(this.settings, 'devDashboard', values);
      this.updateVariable("devDashboard", values.id + '')
    } else {
      this.$set(this.settings, 'devDashboard', undefined);
      this.updateVariable("devDashboard", "")
    }
  }
}
</script>

<style lang="scss" scoped>
.cursor-pointer {
  cursor: pointer;
}
</style>
