<template>
  <div class="app-container logging">
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
        <el-checkbox-group
          v-model="levelFilter"
          size="mini"
          @change="handleFilter"
          style="margin-bottom: 20px"
        >
          <el-checkbox-button v-for="level in levels" :label="level" :key="level">{{ level }}</el-checkbox-button>
        </el-checkbox-group>
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
          :row-class-name="tableRowClassName"
          class="log-table"
        >
          <el-table-column
            :label="$t('log.table.createdAt')"
            prop="createdAt"
            sortable="custom"
            align="left"
            width="150px"
            :class-name="getSortClass('createdAt')"
          >
            <template slot-scope="{row}">
              <span>{{ row.createdAt | parseTime }}</span>
            </template>
          </el-table-column>

          <el-table-column
            :label="$t('log.table.level')"
            prop="level"
            sortable="custom"
            align="left"
            width="90px"
            :class-name="getSortClass('level')"
          >
            <template slot-scope="{row}">
              <span>{{ row.level }}</span>
            </template>
          </el-table-column>

          <el-table-column
            :label="$t('log.table.body')"
            prop="body"
            align="left"
            width="auto"
          >
            <template slot-scope="{row}">
              {{ row.body }}
            </template>
          </el-table-column>

          <el-table-column
            :label="$t('log.table.owner')"
            prop="owner"
            sortable="custom"
            align="left"
            width="150px"
            :class-name="getSortClass('owner')"
          >
            <template slot-scope="{row}">
              <span>{{ row.owner }}</span>
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
import { Component, Vue } from 'vue-property-decorator'
import Pagination from '@/components/Pagination/index.vue'
import api from '@/api/api'
import { ApiLog } from '@/api/stub'
import CardWrapper from '@/components/card-wrapper/index.vue'
import stream from '@/api/stream'
import { requestCurrentState } from '@/views/dashboard/core'
import { UUID } from 'uuid-generator-ts'

@Component({
  name: 'LogList',
  components: {
    CardWrapper,
    Pagination
  }
})
export default class extends Vue {
  private list: ApiLog[] = [];
  private total = 0;
  private listLoading = true;
  private listQuery: { page?: number, limit?: number, sort?: string, query?: string, startDate?: string, endDate?: string } = {
    page: 1,
    limit: 100,
    sort: '-created_at'
  };

  private pageSizes = [50, 100, 150, 250];
  private levels: string[] = ['Emergency', 'Alert', 'Critical', 'Error', 'Warning', 'Notice', 'Info', 'Debug'];
  private levelFilter: string[] = [];
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

  // id for streaming subscribe
  private currentID = '';

  created() {
    this.getList()

    const uuid = new UUID()
    this.currentID = uuid.getDashFreeUUID()

    setTimeout(() => {
      stream.subscribe('log', this.currentID, this.onLogs)
    }, 1000)
  }

  private destroyed() {
    stream.unsubscribe('log', this.currentID)
  }

  onLogs(log: ApiLog) {
    // this.list.push(log);
    // this.total += 1;
    this.getList()// todo optimize
  }

  private async getList() {
    this.listLoading = true
    const object: { page?: number, limit?: number, sort?: string, query?: string, startDate?: string, endDate?: string } = {
      limit: this.listQuery.limit,
      page: this.listQuery.page,
      sort: this.listQuery.sort
    }
    if (this.listQuery.query) {
      object.query = this.listQuery.query
    }
    if (this.listQuery.startDate) {
      object.startDate = this.listQuery.startDate
    }
    if (this.listQuery.endDate) {
      object.endDate = this.listQuery.endDate
    }
    const { data } = await api.v1.logServiceGetLogList(object)

    this.list = data.items
    this.total = data.meta.total
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

    if (this.levelFilter && this.levelFilter.length > 0) {
      this.listQuery.query = this.levelFilter.join(',')
    } else {
      this.listQuery.query = undefined
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
      case 'level':
        if (order === 'ascending') {
          this.listQuery.sort = '+level'
        } else {
          this.listQuery.sort = '-level'
        }
        break
      case 'owner':
        if (order === 'ascending') {
          this.listQuery.sort = '+owner'
        } else {
          this.listQuery.sort = '-owner'
        }
        break
      default:
        console.warn(`unknown field ${prop}`)
    }
    this.handleFilter()
  }

  private getSortClass(key: string) {
    const sort = this.listQuery.sort
    return sort === `+${key}` ? 'ascending' : 'descending'
  }

  private tableRowClassName(data: any): string {
    const { row, index } = data
    let style = ''
    switch (row.level) {
      case 'Emergency':
        style = 'log-emergency'
        break
      case 'Alert':
        style = 'log-alert'
        break
      case 'Critical':
        style = 'log-critical'
        break
      case 'Error':
        style = 'log-error'
        break
      case 'Warning':
        style = 'log-warning'
        break
      case 'Notice':
        style = 'log-notice'
        break
      case 'Info':
        style = 'log-info'
        break
      case 'Debug':
        style = 'log-debug'
        break
    }
    return style
  }
}
</script>

<style lang="scss">

.app-container.logging {

.pagination-container {
  padding: 5px 0;
}

}

.log-table {

td.el-table__cell {
  padding: 0;
  border-bottom: none !important;
}

.log-emergency {
  background-color: #ffc9c9;
}

.log-alert {
  background-color: #ffc9c9;
}

.log-critical {
  background-color: #ffc9c9;
}

.log-error {
  background-color: #ffc9c9;
}

.log-warning {
  background-color: #fff18e;
}

.log-notice {
  background-color: #c1ff89;
}

.log-info {
  background-color: inherit;
}

.log-debug {
  background-color: #82aeff;
}

}

</style>
