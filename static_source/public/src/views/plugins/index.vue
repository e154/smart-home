<template>
    <div class="app-container">

      <el-table
        :key="tableKey"
        v-loading="listLoading"
        :data="list"
        style="width: 100%;"
        @sort-change="sortChange"
      >
        <el-table-column
          :label="$t('plugins.table.name')"
          prop="name"
          sortable="custom"
          align="left"
          width="auto"
          :class-name="getSortClass('name')"
        >
          <template slot-scope="{row}"
                    @click="goto(row)">
            <span class="cursor-pointer">{{ row.name }}</span>
          </template>
        </el-table-column>

        <el-table-column
          :label="$t('plugins.table.version')"
          width="140px"
          align="left"
        >
          <template slot-scope="{row}">
            <span>{{ row.version }}</span>
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
              :disabled="row.system"
            v-on:change="updateItem(row)">
            </el-switch>
          </template>
        </el-table-column>

        <el-table-column
          :label="$t('plugins.table.system')"
          class-name="status-col"
          width="150px"
        >
          <template slot-scope="{row}">

            <el-switch v-model="row.system"
                       disabled>
            </el-switch>
          </template>
        </el-table-column>

      </el-table>

      <pagination
        v-show="total>0"
        :total="total"
        :page.sync="listQuery.page"
        :limit.sync="listQuery.limit"
        @pagination="getList"
      />

    </div>
</template>

<script lang="ts">
import { Component, Vue, Watch } from 'vue-property-decorator'
import api from '@/api/api'
import { ApiPlugin } from '@/api/stub'
import Pagination from '@/components/Pagination/index.vue'
import router from '@/router'

@Component({
  name: 'Plugins',
  components: {
    Pagination
  }
})
export default class extends Vue {
  private tableKey = 0;
  private list: ApiPlugin[] = [];
  private total = 0;
  private listLoading = true;
  private listQuery = {
    page: 1,
    limit: 20,
    sort: '+name'
  };

  created() {
    this.getList()
  }

  private async getList() {
    this.listLoading = true
    const { data } = await api.v1.pluginServiceGetPluginList({
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
    if (prop === 'name') {
      this.sortByName(order)
    }
  }

  private sortByName(order: string) {
    if (order === 'ascending') {
      this.listQuery.sort = '+name'
    } else {
      this.listQuery.sort = '-name'
    }
    this.handleFilter()
  }

  private getSortClass(key: string) {
    const sort = this.listQuery.sort
    return sort === `+${key}` ? 'ascending' : 'descending'
  }

  private goto(plugin: ApiPlugin) {
    router.push({ path: `/plugins/edit/${plugin.name}` })
  }

  private updateItem(plugin: ApiPlugin) {
    if (plugin.enabled) {
      this.enable(plugin)
    } else {
      this.disable(plugin)
    }
  }

  private async enable(plugin: ApiPlugin) {
    const { data } = await api.v1.pluginServiceEnablePlugin(plugin.name)
    this.getList()
  }

  private async disable(plugin: ApiPlugin) {
    const { data } = await api.v1.pluginServiceDisablePlugin(plugin.name)
    this.getList()
  }
}
</script>

<style lang="scss" scoped>
.cursor-pointer {
  cursor: pointer;
}
</style>
