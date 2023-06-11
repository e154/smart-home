<template>
  <div class="app-container">
    <el-row>
      <el-col>
        <el-button type="primary" @click.prevent.stop="add">
          <i class="el-icon-plus"/> {{ $t('zigbee2mqtt.addNew') }}
        </el-button>
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
            :label="$t('zigbee2mqtt.table.id')"
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
            :label="$t('zigbee2mqtt.table.name')"
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
            :label="$t('zigbee2mqtt.table.login')"
            width="auto"
            align="left"
          >
            <template slot-scope="{row}">
              {{ row.login }}
            </template>
          </el-table-column>

          <el-table-column
            :label="$t('zigbee2mqtt.table.permitJoin')"
            class-name="status-col"
            width="150px"
          >
            <template slot-scope="{row}">
              <el-switch
                v-model="row.permitJoin"
                @change="onSwitch"
              >

              </el-switch>
            </template>
          </el-table-column>

          <el-table-column
            :label="$t('zigbee2mqtt.table.createdAt')"
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
            :label="$t('zigbee2mqtt.table.updatedAt')"
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
import { ApiArea, ApiZigbee2MqttShort } from '@/api/stub'
import router from '@/router'

@Component({
  name: 'Index',
  components: {
    Pagination
  }
})
export default class extends Vue {
  private tableKey = 0;
  private list: ApiZigbee2MqttShort[] = [];
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
    const { data } = await api.v1.zigbee2MqttServiceGetBridgeList({
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

  private goto(entity: ApiZigbee2MqttShort) {
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
