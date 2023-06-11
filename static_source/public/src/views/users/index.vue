<template>
  <div class="app-container">

    <el-row>
      <el-col>
        <el-button type="primary" @click.prevent.stop="add">
          <i class="el-icon-plus"/> {{ $t('users.addNew') }}
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
            :label="$t('users.table.id')"
            prop="id"
            sortable="custom"
            align="left"
            width="150px"
            :class-id="getSortClass('id')"
          >
            <template slot-scope="{row}">
              {{ row.id }}
            </template>
          </el-table-column>

          <el-table-column
            :label="$t('users.table.nickname')"
            prop="nickname"
            sortable="custom"
            align="left"
            width="150px"
            :class-nickname="getSortClass('nickname')"
          >
            <template slot-scope="{row}">
            <span class="cursor-pointer"
                  @click="goto(row)">
              {{ row.nickname }}
            </span>
            </template>
          </el-table-column>

          <el-table-column
            :label="$t('users.table.role')"
            prop="roleName"
            sortable="custom"
            align="left"
            width="150px"
            :class-roleName="getSortClass('role')"
          >
            <template slot-scope="{row}">
              <el-tag type="info">
                {{ row.roleName }}
              </el-tag>
            </template>
          </el-table-column>

          <el-table-column
            :label="$t('users.table.email')"
            prop="email"
            sortable="custom"
            align="left"
            width="150px"
            :class-email="getSortClass('email')"
          >
            <template slot-scope="{row}">
              {{ row.email }}
            </template>
          </el-table-column>

          <el-table-column
            :label="$t('users.table.status')"
            prop="status"
            sortable="custom"
            align="left"
            width="150px"
            :class-status="getSortClass('status')"
          >
            <template slot-scope="{row}">
              {{ row.status }}
            </template>
          </el-table-column>

          <el-table-column
            :label="$t('users.table.createdAt')"
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
            :label="$t('users.table.updatedAt')"
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
import { ApiUserShot } from '@/api/stub'
import router from '@/router'

@Component({
  name: 'UserList',
  components: {
    Pagination
  }
})
export default class extends Vue {
  private tableKey = 0;
  private list: ApiUserShot[] = [];
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
    const { data } = await api.v1.userServiceGetUserList({
      limit: this.listQuery.limit,
      page: this.listQuery.page,
      sort: this.listQuery.sort
    })

    this.list = data.items || []
    this.total = data?.meta?.total || 0
    this.listLoading = false
  }

  private handleFilter() {
    this.listQuery.page = 1
    this.getList()
  }

  private sortChange(data: any) {
    const { prop, order } = data
    switch (prop) {
      default:
        if (order === 'ascending') {
          this.listQuery.sort = `+${prop}`
        } else {
          this.listQuery.sort = `-${prop}`
        }
        this.handleFilter()
    }
  }

  private getSortClass(key: string) {
    const sort = this.listQuery.sort
    return sort === `+${key}` ? 'ascending' : 'descending'
  }

  private goto(area: ApiUserShot) {
    router.push({ path: `/etc/users/edit/${area.id}` })
  }

  private add() {
    router.push({ path: '/etc/users/new' })
  }
}
</script>

<style lang="scss" scoped>
.cursor-pointer {
  cursor: pointer;
}
</style>
