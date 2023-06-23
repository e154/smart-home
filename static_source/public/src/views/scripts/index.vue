<template>
  <div class="app-container">
    <el-row>
      <el-col>
        <el-button type="primary" @click.prevent.stop="add">
          <i class="el-icon-plus"/> {{ $t('scripts.addNew') }}
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
            :label="$t('scripts.table.id')"
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
            :label="$t('scripts.table.name')"
            width="140px"
            align="left"
            prop="name"
            sortable="custom"
          >
            <template slot-scope="{row}">
          <span class="cursor-pointer"
                @click="goto(row)">
            {{ row.name }}
          </span>
            </template>
          </el-table-column>

          <el-table-column
            :label="$t('scripts.table.lang')"
            class-name="status-col"
            width="150px"
            prop="lang"
            sortable="custom"
          >
            <template slot-scope="{row}">
              <el-tag type="info">
                {{ row.lang }}
              </el-tag>
            </template>
          </el-table-column>

          <el-table-column
            :label="$t('scripts.table.description')"
            width="auto"
            align="left"
            prop="description"
            sortable="custom"
          >
            <template slot-scope="{row}">
              <span>{{ row.description }}</span>
            </template>
          </el-table-column>

          <el-table-column
            :label="$t('scripts.table.createdAt')"
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
            :label="$t('scripts.table.updatedAt')"
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
          @pagination="updatePagination"
        />
      </el-col>
    </el-row>
  </div>

</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import api from '@/api/api'
import Pagination from '@/components/Pagination/index.vue'
import router from '@/router'
import { ApiScript } from '@/api/stub'

@Component({
  name: 'ScriptList',
  components: {
    Pagination
  }
})
export default class extends Vue {
  private tableKey = 0;
  private list: ApiScript[] = [];
  private total = 0;
  private listLoading = true;
  private listQuery = {
    page: 1,
    limit: 20,
    sort: '-createdAt'
  };
  private defaultSort: Object = {prop: "createdAt", order: "ascending"};
  private itemName = 'ScriptTableSort'

  created() {
    this.restoreSort(); // Восстановление состояния сортировки при загрузке компонента
  }

  private async getList() {
    this.listLoading = true
    const { data } = await api.v1.scriptServiceGetScriptList({
      limit: this.listQuery.limit,
      page: this.listQuery.page,
      sort: this.listQuery.sort
    })

    this.list = data.items
    this.total = data.meta.total
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

  private goto(script: ApiScript) {
    router.push({ path: `/scripts/edit/${script.id}` })
  }

  private add() {
    router.push({ path: '/scripts/new' })
  }
}
</script>

<style lang="scss" scoped>
.cursor-pointer {
  cursor: pointer;
}
</style>
