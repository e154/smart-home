<template>
  <div class="app-container">
    <el-row>
      <el-col>
        <el-button type="primary" @click.prevent.stop="add"><i class="el-icon-plus"/> {{ $t('automation.addNew') }}</el-button>
        <el-button type="primary" @click.prevent.stop="showImport = true">{{ $t('main.import') }}</el-button>

        <export-tool
          :title="$t('main.import')"
          :visible="showImport"
          :value="internal.importValue"
          @on-close="showImport=false"
          @on-import="onImport"
          :import-dialog="true"/>
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
            :label="$t('automation.table.id')"
            prop="id"
            sortable="custom"
            align="left"
            width="60px"
            :class-name="getSortClass('id')"
          >
            <template slot-scope="{row}">
              {{ row.id }}
            </template>
          </el-table-column>

          <el-table-column
            :label="$t('automation.table.name')"
            class-name="status-col"
            align="left"
            width="150px"
          >
            <template slot-scope="{row}">
              <div class="cursor-pointer"
                   @click="goto(row)">
                {{ row.name }}
              </div>
            </template>
          </el-table-column>

          <el-table-column
            :label="$t('automation.table.description')"
            width="auto"
            align="left"
          >
            <template slot-scope="{row}">
              <i v-if="row.description.length == 0" :class="'el-icon-minus'"/> {{ row.description }}
            </template>
          </el-table-column>

          <el-table-column
            :label="$t('plugins.table.enabled')"
            class-name="status-col"
            width="150px"
          >
            <template slot-scope="{row}">
              <el-switch
                v-model="row.enabled"
                v-on:change="onSwitch(row)"
                >

              </el-switch>
            </template>
          </el-table-column>

          <el-table-column
            :label="$t('automation.table.createdAt')"
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
            :label="$t('automation.table.updatedAt')"
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
import { Component, Vue } from 'vue-property-decorator'
import Pagination from '@/components/Pagination/index.vue'
import api from '@/api/api'
import { ApiArea, ApiTask } from '@/api/stub'
import router from '@/router'
import ExportTool from '@/components/export-tool/index.vue'

@Component({
  name: 'Index',
  components: {
    ExportTool,
    Pagination
  }
})
export default class extends Vue {
  private tableKey = 0;
  private list: ApiTask[] = [];
  private total = 0;
  private listLoading = true;
  private listQuery = {
    page: 1,
    limit: 20,
    sort: '-createdAt'
  };

  private internal = {
    importValue: ''
  };

  created() {
    this.getList()
  }

  private async getList() {
    this.listLoading = true
    const { data } = await api.v1.automationServiceGetTaskList({
      limit: this.listQuery.limit,
      page: this.listQuery.page,
      sort: this.listQuery.sort
    })

    this.list = data.items
    this.total = data.meta.total
    this.listLoading = false
  }

  private handleFilter() {
    this.listQuery.page = 1
    this.getList()
  }

  private sortChange(data: any) {
    const { prop, order } = data
    if (prop === 'id') {
      this.sortByID(order)
    } else if (prop === 'createdAt') {
      this.sortByCreatedAt(order)
    } else if (prop === 'updatedAt') {
      this.sortByUpdatedAt(order)
    }
  }

  private sortByCreatedAt(order: string) {
    if (order === 'ascending') {
      this.listQuery.sort = '+createdAt'
    } else {
      this.listQuery.sort = '-createdAt'
    }
    this.handleFilter()
  }

  private sortByUpdatedAt(order: string) {
    if (order === 'ascending') {
      this.listQuery.sort = '+updatedAt'
    } else {
      this.listQuery.sort = '-updatedAt'
    }
    this.handleFilter()
  }

  private sortByID(order: string) {
    if (order === 'ascending') {
      this.listQuery.sort = '+id'
    } else {
      this.listQuery.sort = '-id'
    }
    this.handleFilter()
  }

  private getSortClass(key: string) {
    const sort = this.listQuery.sort
    return sort === `+${key}` ? 'ascending' : 'descending'
  }

  private goto(entity: ApiTask) {
    router.push({ path: `/automation/edit/${entity.id}` })
  }

  private add() {
    router.push({ path: '/automation/new' })
  }

  private gotoArea(area: ApiArea) {
    router.push({ path: `/areas/edit/${area.id}` })
  }

  private async onSwitch(event: ApiTask) {
    if (event.enabled) {
      await api.v1.automationServiceEnableTask(event.id || 0)
    } else {
      await api.v1.automationServiceDisableTask(event.id || 0)
    }
    this.$notify({
      title: 'Success',
      message: 'Update Successfully',
      type: 'success',
      duration: 2000
    })
  }

  private showImport = false;

  private async onImport(value: string, event?: any) {
    const val: ApiTask = JSON.parse(value)
    const task = {
      name: val.name,
      description: val.description,
      enabled: val.enabled,
      condition: val.condition,
      triggers: val.triggers,
      conditions: val.conditions,
      actions: val.actions,
      area: val.area
    }
    const { data } = await api.v1.automationServiceAddTask(task)
    if (data) {
      this.$notify({
        title: 'Success',
        message: 'task imported successfully',
        type: 'success',
        duration: 2000
      })
      this.getList()
    }
  }
}
</script>

<style lang="scss" scoped>

.cursor-pointer {
  cursor: pointer;
}

.pagination-container {

}
</style>
