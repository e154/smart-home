<template>
  <div class="app-container logging">

    <div v-if="current && current && vis">
      <el-dialog
        :visible.sync="vis"
        width="50%"
        append-to-body
        destroy-on-close
        :title="$t('message.attributes')"
      >
        <message-viewer v-model="current.message"/>
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
        <el-checkbox-group
          v-model="typeFilter"
          size="mini"
          @change="handleFilter"
          style="margin-bottom: 20px"
        >
          <el-checkbox-button v-for="type in messageTypes" :label="type" :key="type">{{ type }}</el-checkbox-button>
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
          :row-class-name="tableRowStatus"
          class="log-table"
          @current-change="handleCurrentChange"
        >

          <el-table-column
            :label="$t('message_delivery.table.id')"
            prop="id"
            sortable="custom"
            align="left"
            width="90px"
            :class-name="getSortClass('id')"
          >
            <template slot-scope="{row}">
              <span>{{ row.id }}</span>
            </template>
          </el-table-column>

          <el-table-column
            :label="$t('message.table.type')"
            prop="type"
            align="left"
            width="100px"
            :class-name="getSortClass('type')"
          >
            <template slot-scope="{row}">
              <span>{{ row.message.type }}</span>
            </template>
          </el-table-column>

          <el-table-column
            :label="$t('message_delivery.table.attributes')"
            prop="body"
            align="left"
            width="auto"
          >
            <template slot-scope="{row}">
              <span>{{ Object.keys(row.message.attributes).length || $t('message_delivery.table.nothing') }}</span>
            </template>
          </el-table-column>

          <el-table-column
            :label="$t('message_delivery.table.status')"
            prop="status"
            sortable="custom"
            align="left"
            width="150px"
            :class-name="getSortClass('status')"
          >
            <template slot-scope="{row}">
              <span>{{ row.status }}</span>
            </template>
          </el-table-column>

          <el-table-column
            :label="$t('message_delivery.table.createdAt')"
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
            :label="$t('message_delivery.table.updatedAt')"
            prop="updatedAt"
            sortable="custom"
            align="left"
            width="150px"
            :class-name="getSortClass('updatedAt')"
          >
            <template slot-scope="{row}">
              <span>{{ row.updatedAt | parseTime }}</span>
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
import {Component, Vue} from 'vue-property-decorator';
import Pagination from '@/components/Pagination/index.vue';
import api from '@/api/api';
import {ApiEntityStorage, ApiMessageDelivery} from '@/api/stub';
import CardWrapper from '@/components/card-wrapper/index.vue';
import stream from '@/api/stream';
import {UUID} from 'uuid-generator-ts';
import MessageViewer from '@/views/message_delivery/message_viewer.vue';

@Component({
  name: 'LogList',
  components: {
    MessageViewer,
    CardWrapper,
    Pagination,
  }
})
export default class extends Vue {

  private current?: ApiMessageDelivery | undefined = {} as ApiMessageDelivery;
  private showDialog = false;

  private list: ApiMessageDelivery[] = [];
  private total = 0;
  private listLoading = true;
  private listQuery: { page?: number, limit?: number, sort?: string, messageTypes?: string, startDate?: string, endDate?: string } = {
    page: 1,
    limit: 100,
    sort: '-created_at'
  };

  private pageSizes = [50, 100, 150, 250];
  private messageTypes: string[] = ['webpush', 'html5_notify', 'email', 'sms', 'telegram'];
  private typeFilter: string[] = [];
  private dateFilter: Date[] = [];
  private pickerOptions: Object = {
    shortcuts: [{
      text: 'Last week',
      onClick(picker: any) {
        const end = new Date();
        const start = new Date();
        start.setTime(start.getTime() - 3600 * 1000 * 24 * 7);
        picker.$emit('pick', [start, end]);
      }
    }, {
      text: 'Last month',
      onClick(picker: any) {
        const end = new Date();
        const start = new Date();
        start.setTime(start.getTime() - 3600 * 1000 * 24 * 30);
        picker.$emit('pick', [start, end]);
      }
    }, {
      text: 'Last 3 months',
      onClick(picker: any) {
        const end = new Date();
        const start = new Date();
        start.setTime(start.getTime() - 3600 * 1000 * 24 * 90);
        picker.$emit('pick', [start, end]);
      }
    }]
  };

  // id for streaming subscribe
  private currentID = '';

  created() {
    this.getList();

    const uuid = new UUID();
    this.currentID = uuid.getDashFreeUUID();

    setTimeout(() => {
      stream.subscribe('message_delivery', this.currentID, this.onLogs);
    }, 1000);
  }

  private destroyed() {
    stream.unsubscribe('message_delivery', this.currentID);
  }

  onLogs(item: ApiMessageDelivery) {
    // this.list.push(item);
    // this.total += 1;
    this.getList();// todo optimize
  }

  private async getList() {
    this.listLoading = true;
    const object: { page?: number, limit?: number, sort?: string, message_type?: string, startDate?: string, endDate?: string } = {
      limit: this.listQuery.limit,
      page: this.listQuery.page,
      sort: this.listQuery.sort
    };
    if (this.listQuery.messageTypes) {
      object.message_type = this.listQuery.messageTypes;
    }
    if (this.listQuery.startDate) {
      object.startDate = this.listQuery.startDate;
    }
    if (this.listQuery.endDate) {
      object.endDate = this.listQuery.endDate;
    }
    const {data} = await api.v1.messageDeliveryServiceGetList(object);

    this.list = data.items;
    this.total = data.meta.total;
    this.listLoading = false;
  }

  private handleFilter() {
    if (this.dateFilter && this.dateFilter.length > 1) {
      this.listQuery.startDate = this.dateFilter[0].toISOString().substring(0, 10);
      this.listQuery.endDate = this.dateFilter[1].toISOString().substring(0, 10);
    } else {
      this.listQuery.startDate = undefined;
      this.listQuery.endDate = undefined;
    }

    if (this.typeFilter && this.typeFilter.length > 0) {
      this.listQuery.messageTypes = this.typeFilter.join(',');
    } else {
      this.listQuery.messageTypes = undefined;
    }
    this.listQuery.page = 1;
    this.getList();
  }

  private sortChange(data: any) {
    const {prop, order} = data;
    switch (prop) {
      case 'id':
        if (order === 'ascending') {
          this.listQuery.sort = '+id';
        } else {
          this.listQuery.sort = '-id';
        }
        break;
      case 'createdAt':
        if (order === 'ascending') {
          this.listQuery.sort = '+createdAt';
        } else {
          this.listQuery.sort = '-createdAt';
        }
        break;
      case 'updatedAt':
        if (order === 'ascending') {
          this.listQuery.sort = '+updatedAt';
        } else {
          this.listQuery.sort = '-updatedAt';
        }
        break;
      case 'status':
        if (order === 'ascending') {
          this.listQuery.sort = '+status';
        } else {
          this.listQuery.sort = '-status';
        }
        break;
      default:
        console.warn(`unknown field ${prop}`);
    }
    this.handleFilter();
  }

  private getSortClass(key: string) {
    const sort = this.listQuery.sort;
    return sort === `+${key}` ? 'ascending' : 'descending';
  }

  private tableRowStatus(data: any): string {
    const {row, index} = data;
    let style = '';
    switch (row.status) {
      case 'in_progress':
        style = 'in_progress';
        break;
      case 'succeed':
        style = 'succeed';
        break;
      case 'error':
        style = 'error';
        break;
    }
    return style;
  }

  get vis(): boolean {
    return this.showDialog
  }

  set vis(value: boolean) {
    this.showDialog = false
  }

  private handleCurrentChange(val?: ApiMessageDelivery) {
    this.current = val
    this.showDialog = !this.showDialog
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

.error {
  background-color: #ffc9c9;
}

.succeed {
  background-color: inherit;
}

.in_progress {
  background-color: #82aeff;
}

}

</style>
