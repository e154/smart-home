<template>
  <el-row>
    <el-col>
      <el-table
        :key="reloadKey"
        :data="list"
        style="width: 100%;"
        @sort-change="sortChange"
        :row-class-name="tableRowClassName"
        class="log-table"
      >
        <el-table-column
          :label="$t('log.table.createdAt')"
          prop="name"
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
          prop="name"
          sortable="custom"
          align="left"
          width="100px"
          :class-name="getSortClass('level')"
        >
          <template slot-scope="{row}">
            <span>{{ row.level }}</span>
          </template>
        </el-table-column>

        <el-table-column
          :label="$t('log.table.body')"
          prop="name"
          align="left"
          width="auto"
        >
          <template slot-scope="{row}">
            {{ row.body }}
          </template>
        </el-table-column>

      </el-table>
    </el-col>
  </el-row>
</template>

<script lang="ts">
import { Component, Prop, Vue, Watch } from 'vue-property-decorator'
import {CardItem, requestCurrentState} from '@/views/dashboard/core';
import { ApiLog } from '@/api/stub'
import api from '@/api/api'
import { UUID } from 'uuid-generator-ts'
import stream from '@/api/stream'

interface list {
  page?: number
  limit?: number
  sort?: string
  query?: string
  startDate?: string
  endDate?: string
}

@Component({
  name: 'Logs',
  components: {}
})
export default class extends Vue {
  @Prop() private item?: CardItem;

  private reloadKey = 0;
  private list: ApiLog[] = [];
  private total = 0;
  private listLoading = true;
  private listQuery: list = {
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
    this.listQuery.limit = this.item?.payload.logs?.limit || 20
    this.getList()

    const uuid = new UUID()
    this.currentID = uuid.getDashFreeUUID()

    setTimeout(() => {
      stream.subscribe('log', this.currentID, this.onLogs)
    }, 1000)

    requestCurrentState(this.item?.entityId);
  }

  private destroyed() {
    stream.unsubscribe('log', this.currentID)
  }

  onLogs(log: ApiLog) {
    // this.list.push(log);
    // this.total += 1;
    this.getList()// todo optimize
  }

  reload() {
    this.reloadKey += 1
  }

  @Watch('item.uuid')
  private onUpdateItem(item: CardItem) {
    this.getList()
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
    if (prop === 'createdAt') {
      this.createdAt(order)
    }
  }

  private createdAt(order: string) {
    if (order === 'ascending') {
      this.listQuery.sort = '+createdAt'
    } else {
      this.listQuery.sort = '-createdAt'
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
