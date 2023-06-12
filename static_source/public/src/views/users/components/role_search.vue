<template>
  <el-select
    style="width: 100%"
    v-model="currentValue"
    :multiple="multiple"
    filterable
    default-first-option
    remote
    clearable
    value-key="id"
    reserve-keyword
    placeholder="Please enter a keyword"
    :remote-method="remoteMethod"
    :loading="loading"
  >
    <el-option
      v-for="(item, index) in options"
      :key="item.id"
      :label="item.name"
      :value="item">
    </el-option>
  </el-select>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator'
import { ApiRole } from '@/api/stub'
import api from '@/api/api'

@Component({
  name: 'RoleSearch'
})
export default class extends Vue {
  @Prop({ required: true }) private value!: any;
  @Prop({ required: false, default: () => true }) private multiple?: boolean;

  private options?: ApiRole[] = [];
  private loading = true;

  get currentValue() {
    // array
    if (this.multiple) {
      this.options = this.value as ApiRole[]
      return this.value
    }
    // object
    if (this.value) {
      this.options = [this.value as ApiRole]
      return this.value as ApiRole
    } else {
      this.options = undefined
    }

    return undefined
  }

  set currentValue(value) {
    // array
    if (this.multiple) {
      const result: ApiRole[] = []
      for (const item in value) {
        result.push(value[+item])
      }
      this.$emit('update-value', result)
    }
    // object
    this.$emit('update-value', value)
  }

  private async remoteMethod(query: string) {
    if (query !== '') {
      this.loading = true
      const params = { query: query, limit: 25, offset: 0 }
      const { data } = await api.v1.roleServiceSearchRoleByName(params)
      this.options = data.items
      this.loading = false
    } else {
      this.options = []
    }
  }
}
</script>
