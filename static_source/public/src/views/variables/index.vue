<template>
  <div class="app-container">

    <el-row>
      <el-col>
        <el-button type="primary" @click.prevent.stop="add">
          <i class="el-icon-plus"/> {{ $t('variables.addNew') }}
        </el-button>
      </el-col>
    </el-row>

    <el-row>
      <el-col>
        <pagination
          v-show="total>listQuery.limit"
          :total="total"
          :page.sync="listQuery.page"
          :limit.sync="listQuery.limit"
          @pagination="updatePagination"
        />
      </el-col>
    </el-row>

    <el-row>
      <el-col>
        <el-table
          :key="tableKey"
          :default-sort="defaultSort"
          v-loading="listLoading"
          :data="list"
          style="width: 100%;"
          @sort-change="sortChange"
        >
          <el-table-column
            :label="$t('variables.table.name')"
            prop="name"
            sortable="custom"
            align="left"
            width="150px"
            :class-name="getSortClass('name')"
          >
            <template slot-scope="{row}">
            <span class="cursor-pointer"
                  @click="goto(row)">
              {{ row.name }}
            </span>
            </template>
          </el-table-column>

          <el-table-column
            :label="$t('variables.table.value')"
            width="auto"
            align="left"
          >
            <template slot-scope="{row}">
          <span class="cursor-pointer">
            {{ row.value }}
          </span>
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
          @pagination="updatePagination"
        />
      </el-col>
    </el-row>

  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import Pagination from '@/components/Pagination/index.vue'
import api from '@/api/api'
import {ApiVariable} from '@/api/stub';
import router from '@/router'

@Component({
  name: 'Variables',
  components: {
    Pagination
  }
})
export default class extends Vue {
  private tableKey = 0;
  private list: ApiVariable[] = [];
  private total = 0;
  private listLoading = true;
  private listQuery = {
    page: 1,
    limit: 20,
    sort: '-name'
  };
  private defaultSort: Object = {prop: "name", order: "ascending"};
  private itemName = 'VariablesTableSort'

  created() {
    this.restoreSort(); // Восстановление состояния сортировки при загрузке компонента
  }

  private async getList() {
    this.listLoading = true
    const { data } = await api.v1.variableServiceGetVariableList({
      limit: this.listQuery.limit,
      page: this.listQuery.page,
      sort: this.listQuery.sort
    })

    this.list = data.items || []
    this.total = data.meta?.total || 0
    this.listLoading = false
  }

  private handleFilter() {
    this.getList()
  }

  private updatePagination() {
    const sortData = localStorage.getItem(this.itemName);
    if (sortData) {
      const {column} = JSON.parse(sortData);
      localStorage.setItem(this.itemName, JSON.stringify({
        column: column,
        page: this.listQuery.page,
        limit: this.listQuery.limit
      }));
    }
    this.getList()
  }

  private restoreSort() {
    // Восстановление состояния сортировки при загрузке компонента
    const sortData = localStorage.getItem(this.itemName);
    if (sortData) {
      const {column, page, limit} = JSON.parse(sortData);
      this.defaultSort = {prop: column.property, order: column.order};
      this.listQuery.page = page
      this.listQuery.limit = limit
      this.sort(column.property, column.order)
    } else {
      this.getList()
    }
  }

  private sortChange(data: any) {
    // Обработчик изменения состояния сортировки
    const {column, prop, order} = data;
    this.defaultSort = {prop, order};
    // Сохраняем состояние сортировки в localStorage
    localStorage.setItem(this.itemName, JSON.stringify({
      column: column,
      page: this.listQuery.page,
      limit: this.listQuery.limit
    }));
    this.sort(prop, order)
  }

  private sort(column: string, order: string) {
    let pref: string = '-'
    if (order === 'ascending') {
      pref = '+'
    }
    this.listQuery.sort = pref + column
    this.handleFilter()
  }

  private getSortClass(key: string) {
    const sort = this.listQuery.sort
    return sort === `+${key}` ? 'ascending' : 'descending'
  }

  private goto(variable: ApiVariable) {
    router.push({ path: `/etc/variables/edit/${variable.name}` })
  }

  private add() {
    router.push({ path: '/etc/variables/new' })
  }
}
</script>

<style lang="scss" scoped>
.cursor-pointer {
  cursor: pointer;
}
</style>
