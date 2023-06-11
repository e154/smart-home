<template>

  <div>
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
          @pagination="getList"
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

  created() {
    this.getList()
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
