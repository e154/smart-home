<template>

  <div>
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
            :label="$t('zigbee2mqtt.table.id')"
            prop="id"
            sortable="custom"
            align="left"
            width="160px"
            :class-name="getSortClass('id')"
          >
            <template slot-scope="{row}">
              {{ row.id }}
            </template>
          </el-table-column>

          <el-table-column
            :label="$t('zigbee2mqtt.table.model')"
            class-name="status-col"
            align="left"
            header-align="left"
            width="100px"
            prop="model"
            sortable="custom"
          >
            <template slot-scope="{row}">

                {{ row.model }}
            </template>
          </el-table-column>

          <el-table-column
            :label="$t('zigbee2mqtt.table.status')"
            width="70px"
            align="left"
            header-align="left"
            prop="status"
            sortable="custom"
          >
            <template slot-scope="{row}">
              {{ row.status }}
            </template>
          </el-table-column>

          <el-table-column
            :label="$t('zigbee2mqtt.table.description')"
            class-name="status-col"
            width="auto"
            align="left"
            header-align="left"
            prop="description"
            sortable="custom"
          >
            <template slot-scope="{row}">
              {{row.description}}
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
import { Component, Prop, Vue } from 'vue-property-decorator'
import Pagination from '@/components/Pagination/index.vue'
import api from '@/api/api'
import { ApiArea, ApiZigbee2MqttDevice } from '@/api/stub'
import router from '@/router'

@Component({
  name: 'Devices',
  components: {
    Pagination
  }
})
export default class extends Vue {
  @Prop({ required: true }) private id!: number;

  private tableKey = 0;
  private list: ApiZigbee2MqttDevice[] = [];
  private total = 0;
  private listLoading = true;
  private listQuery = {
    page: 1,
    limit: 20,
    sort: '-createdAt'
  };
  private defaultSort: Object = {prop: "createdAt", order: "ascending"};
  private itemName = 'Zigbee2mqttDevicesrdTableSort'

  created() {
    this.restoreSort(); // Восстановление состояния сортировки при загрузке компонента
  }

  private async getList() {
    this.listLoading = true
    const { data } = await api.v1.zigbee2MqttServiceDeviceList(this.id, {
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

  private goto(entity: ApiZigbee2MqttDevice) {
    router.push({ path: `/zigbee2mqtt/edit/${entity.id}` })
  }

  private add() {
    router.push({ path: '/zigbee2mqtt/new' })
  }

  private onSwitch(value: boolean) {

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
