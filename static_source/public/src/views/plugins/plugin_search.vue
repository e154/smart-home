<template>
  <el-select
    style="width: 100%"
    v-model="currentValue"
    filterable
    default-first-option
    remote
    clearable
    value-key="name"
    reserve-keyword
    placeholder="Please enter a keyword"
    :remote-method="remoteMethod"
    :loading="loading"
  >
    <el-option
      v-for="(item, index) in options"
      :key="item.name"
      :label="item.name"
      :value="item">
    </el-option>
  </el-select>
</template>

<script lang="ts">
import { Component, Prop, Vue, Watch } from 'vue-property-decorator'
import { ApiPlugin } from '@/api/stub'
import api from '@/api/api'

@Component({
  name: 'PluginSearch'
})
export default class extends Vue {
  @Prop({ required: false }) private value?: ApiPlugin;

  private options?: ApiPlugin[] = [];
  private loading = true;

  get currentValue() {
    if (this.value) {
      this.options = [this.value]
      return this.value
    } else {
      this.options = []
    }

    return undefined
  }

  set currentValue(value) {
    this.$emit('update-value', value)
  }

  private async remoteMethod(query: string) {
    if (query !== '') {
      this.loading = true
      const params = { query: query, limit: 25, offset: 0 }
      const { data } = await api.v1.pluginServiceSearchPlugin(params)
      this.options = data.items
      this.loading = false
    } else {
      this.options = []
    }
  }
}
</script>
