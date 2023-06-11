<template>
  <div class="app-container logging">
    <div v-if="current && current.attributes && vis">
      <el-dialog
        :visible.sync="vis"
        width="50%"
        append-to-body
        destroy-on-close
        :title="'entity: ' + current.entityId"
      >
        <attributes-viewer v-model="current.attributes"/>
      </el-dialog>
    </div>

    <el-row>
      <el-col>
        <el-date-picker
          style="width: 100%; margin-bottom: 20px"
          v-model="dateFilter"
          type="daterange"
          align="right"
          unlink-panels
          range-separator="To"
          start-placeholder="Start date"
          end-placeholder="End date"
          :picker-options="pickerOptions"
          format="yyyy-MM-dd"
          @change="handleFilter"
        >
        </el-date-picker>
      </el-col>
    </el-row>

    <el-row>
      <el-col>
        <pagination
          v-show="total>0"
          :total="total"
          :pageSizes="pageSizes"
          :page.sync="listQuery.page"
          :limit.sync="listQuery.limit"
          @pagination="getList"
        />
      </el-col>
    </el-row>

    <el-row>
      <el-col>
        <el-table
          :data="list"
          style="width: 100%;"
          @sort-change="sortChange"
          highlight-current-row
          @current-change="handleCurrentChange"
        >
          <el-table-column
            :label="$t('entityStorage.table.createdAt')"
            prop="createdAt"
            sortable="custom"
            align="left"
            width="150px"
          >
            <template slot-scope="{row}">
              <span>{{ row.createdAt | parseTime }}</span>
            </template>
          </el-table-column>

          <el-table-column
            :label="$t('entityStorage.table.state')"
            prop="state"
            sortable="custom"
            align="left"
            width="200px"
          >
            <template slot-scope="{row}">
              <span>{{ row.state }}</span>
            </template>
          </el-table-column>

          <el-table-column
            :label="$t('entityStorage.table.attributes')"
            prop="attributes"
            align="left"
            width="auto"
          >
            <template slot-scope="{row}">
              <span>{{ Object.keys(row.attributes).length || $t('entityStorage.table.nothing') }}</span>
            </template>
          </el-table-column>

          <el-table-column
            :label="$t('entityStorage.table.entityId')"
            prop="state"
            sortable="custom"
            align="left"
            width="200px"
          >
            <template slot-scope="{row}">
              <span>{{ row.entityId }}</span>
            </template>
          </el-table-column>

        </el-table>
      </el-col>
    </el-row>

    <el-row style="margin-top: 20px">
      <el-col>
        <pagination
          v-show="total>0"
          :total="total"
          :pageSizes="pageSizes"
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
import { ApiEntity, ApiEntityStorage } from '@/api/stub'
import CardWrapper from '@/components/card-wrapper/index.vue'
import AttributesViewer from '@/views/entities/components/attributes_viewer.vue'

@Component({
  name: 'Storage',
  components: {
    AttributesViewer,
    CardWrapper,
    Pagination
  }
})
export default class extends Vue {
  @Prop({ required: true }) private entity!: ApiEntity;

  private list: ApiEntityStorage[] = [];
  private current?: ApiEntityStorage | undefined = {} as ApiEntityStorage;
  private showDialog = false;

  private total = 0;
  private listLoading = true;
  private listQuery: { page?: number, limit?: number, sort?: string, startDate?: string, endDate?: string } = {
    page: 1,
    limit: 100,
    sort: '-created_at'
  };

  private pageSizes = [50, 100, 150, 250];
  private dateFilter: Date[] = [];
  private pickerOptions: Object = {
    shortcuts: [{
      text: 'Last week',
      onClick(picker: any) {
        const end = new Date()
        const start = new Date()
        start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
        picker.$emit('pick', [start, end])
      }
    }, {
      text: 'Last month',
      onClick(picker: any) {
        const end = new Date()
        const start = new Date()
        start.setTime(start.getTime() - 3600 * 1000 * 24 * 30)
        picker.$emit('pick', [start, end])
      }
    }, {
      text: 'Last 3 months',
      onClick(picker: any) {
        const end = new Date()
        const start = new Date()
        start.setTime(start.getTime() - 3600 * 1000 * 24 * 90)
        picker.$emit('pick', [start, end])
      }
    }]
  };

  created() {
    this.getList()
  }

  private async getList() {
    this.listLoading = true
    const object: { page?: number, limit?: number, sort?: string, startDate?: string, endDate?: string } = {
      limit: this.listQuery.limit,
      page: this.listQuery.page,
      sort: this.listQuery.sort
    }
    if (this.listQuery.startDate) {
      object.startDate = this.listQuery.startDate
    }
    if (this.listQuery.endDate) {
      object.endDate = this.listQuery.endDate
    }

    if (!this.entity.id) {
      return
    }

    const { data } = await api.v1.entityStorageServiceGetEntityStorageList(this.entity.id, object)

    if (data) {
      this.list = data.items || []
      this.total = data?.meta?.total || 0
    }

    this.listLoading = false
  }

  private handleFilter() {
    if (this.dateFilter && this.dateFilter.length > 1) {
      this.listQuery.startDate = this.dateFilter[0].toISOString().substring(0, 10)
      this.listQuery.endDate = this.dateFilter[1].toISOString().substring(0, 10)
    } else {
      this.listQuery.startDate = undefined
      this.listQuery.endDate = undefined
    }

    this.listQuery.page = 1
    this.getList()
  }

  private sortChange(data: any) {
    const { prop, order } = data
    switch (prop) {
      case 'createdAt':
        if (order === 'ascending') {
          this.listQuery.sort = '+createdAt'
        } else {
          this.listQuery.sort = '-createdAt'
        }
        break
      case 'state':
        if (order === 'ascending') {
          this.listQuery.sort = '+state'
        } else {
          this.listQuery.sort = '-state'
        }
        break
      default:
        console.warn(`unknown field ${prop}`)
    }
    this.handleFilter()
  }

  get vis(): boolean {
    return this.showDialog
  }

  set vis(value: boolean) {
    this.showDialog = false
  }

  private handleCurrentChange(val?: ApiEntityStorage) {
    this.current = val
    this.showDialog = true
  }
}
</script>

<style lang="scss">

</style>
