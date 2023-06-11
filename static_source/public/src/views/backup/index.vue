<template>
  <div class="app-container">

    <el-row>
      <el-col>
        <el-button type="primary" @click.prevent.stop="addNew">
          <i class="el-icon-plus"/> {{ $t('backup.addNew') }}
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
        >
          <el-table-column
            :label="$t('backup.table.name')"
            align="left"
          >
            <template slot-scope="{row}">
              {{ row }}
            </template>
          </el-table-column>

          <el-table-column
            :label="$t('entities.table.operations')"
            align="right"
          >
            <template slot-scope="{row, $index}">

              <el-popconfirm
                :confirm-button-text="$t('main.ok')"
                :cancel-button-text="$t('main.no')"
                icon="el-icon-info"
                icon-color="red"
                style="margin-left: 10px;"
                :title="$t('main.are_you_sure_to_do_want_this?')"
                v-on:confirm="restoreImage(row, $index)"
              >

                <el-button type="danger" slot="reference">{{ $t('backup.table.restore') }}</el-button>

              </el-popconfirm>

            </template>
          </el-table-column>


        </el-table>
      </el-col>
    </el-row>

  </div>
</template>

<script lang="ts">
import {Component, Vue} from 'vue-property-decorator';
import Pagination from '@/components/Pagination/index.vue';
import api from '@/api/api';
import stream from '@/api/stream';

@Component({
  name: 'AreasList',
  components: {
    Pagination
  }
})
export default class extends Vue {
  private tableKey = 0;
  private list: string[] = [];
  private listLoading = true;

  created() {
    this.getList();
  }

  private async getList() {
    this.listLoading = true;
    const {data} = await api.v1.backupServiceGetBackupList();
    this.list = data.items;
    this.listLoading = false;
  }

  private async restoreImage(name: string, index: number) {
    await api.v1.backupServiceRestoreBackup({name: name});
    setTimeout(async () => {
      await this.getList();
      this.$notify({
        title: 'Success',
        message: 'backup restored successfully',
        type: 'success',
        duration: 2000
      });
    }, 2000)
  }

  private async addNew() {
    await api.v1.backupServiceNewBackup({});
    await this.getList();
    this.$notify({
      title: 'Success',
      message: 'backup created successfully',
      type: 'success',
      duration: 2000
    });
  }

}
</script>

<style lang="scss" scoped>
.cursor-pointer {
  cursor: pointer;
}
</style>
