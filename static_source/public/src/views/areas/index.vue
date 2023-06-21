<template>
  <div class="app-container">

    <el-row>
      <el-col>
        <el-button type="primary" @click.prevent.stop="add">
          <i class="el-icon-plus"/> {{ $t('areas.addNew') }}
        </el-button>
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
            :label="$t('plugins.table.name')"
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
            :label="$t('scripts.table.description')"
            width="auto"
            align="left"
          >
            <template slot-scope="{row}">
          <span class="cursor-pointer">
            {{ row.description }}
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
import {Component, Vue} from 'vue-property-decorator'
import Pagination from '@/components/Pagination/index.vue'
import api from '@/api/api'
import {ApiArea} from '@/api/stub'
import router from '@/router'

@Component({
  name: 'AreasList',
  components: {
    Pagination
  }
})
export default class extends Vue {
  private tableKey = 0;
  private list: ApiArea[] = [];
  private total = 0;
  private listLoading = true;
  private listQuery = {
    page: 1,
    limit: 20,
    sort: '-name'
  };
  private defaultSort: Object = {prop: "name", order: "ascending"};
  private itemName = 'AreasTableSort'

  created() {
    this.restoreSort(); // Восстановление состояния сортировки при загрузке компонента
  }

  private async getList() {
    this.listLoading = true
    const {data} = await api.v1.areaServiceGetAreaList({
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
      if (column.property === 'name') {
        this.sort(column.property, column.order)
      }
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
    this.sort(prop, column.order)
  }

  private sort(column: string, order: string) {
    let pref: string = '-'
    if (order === 'ascending') {
      pref = '+'
    }
    switch (column) {
      case 'name':
        this.listQuery.sort = pref + 'name'
        break
    }
    this.handleFilter()
  }

  private getSortClass(key: string) {
    const sort = this.listQuery.sort
    return sort === `+${key}` ? 'ascending' : 'descending'
  }

  private goto(area: ApiArea) {
    router.push({path: `areas/edit/${area.id}`})
  }

  private add() {
    router.push({path: 'areas/new'})
  }
}
</script>

<style lang="scss" scoped>
.cursor-pointer {
  cursor: pointer;
}
</style>
