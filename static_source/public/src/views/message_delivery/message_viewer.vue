<template>
  <el-row>
    <el-col>
      <el-table
        :data="attributes">
        <el-table-column
          prop="name"
          :label="$t('entities.table.name')"
          width="180px">
          <template slot-scope="{row}">
            <span>{{ row.name }}</span>
          </template>
        </el-table-column>

        <el-table-column
          width="auto"
          :label="$t('entities.table.value')"
        >

          <template slot-scope="{row}">
            <span>{{ row.value }}</span>
          </template>

        </el-table-column>

      </el-table>
    </el-col>
  </el-row>
</template>

<script lang="ts">
import {Component, Prop, Vue} from 'vue-property-decorator';
import {ApiMessage} from '@/api/stub';

export interface MessageItem {
  name: string;
  value: string;
}

@Component({
  name: 'MessageViewer',
  components: {}
})
export default class extends Vue {
  @Prop({required: true}) private value?: ApiMessage;

  get attributes(): MessageItem[] {
    const attr: MessageItem[] = [];
    if (this.value) {
      for (const key in this.value.attributes) {
        attr.push({name: key, value: this.value.attributes[key]});
      }
    }
    return attr;
  }

  set attributes(value: MessageItem[]) {

  }
}
</script>
