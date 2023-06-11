<template>
  <el-select
    style="width: 100%"
    v-model="currentValue"
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
      :label="item.id"
      :value="{id: item.id}">
    </el-option>
  </el-select>
</template>

<script lang="ts">
import { Component, Prop, Vue, Watch } from 'vue-property-decorator'
import { ApiEntityShort } from '@/api/stub'
import api from '@/api/api'

@Component({
  name: 'EntitySearch'
})
export default class extends Vue {
  @Prop() private value?: ApiEntityShort;

  private options?: ApiEntityShort[] = [];
  private loading = true;

  private update() {
    if (this.value) {
      this.options = [this.value]
    }
  }

  created() {
    this.update()
  }

  @Watch('value')
  private onValueChanged() {
    this.update()
  }

  get currentValue(): ApiEntityShort | undefined {
    if (this.value) {
      return this.value
    } else {
      return undefined
    }
  }

  set currentValue(value) {
    this.$emit('update-value', value)
  }

  private async remoteMethod(query: string) {
    if (query !== '') {
      this.loading = true
      const params = { query: query, limit: 25, offset: 0 }
      const { data } = await api.v1.entityServiceSearchEntity(params)
      this.options = data.items
      this.loading = false
    } else {
      this.options = []
    }
  }
}
</script>
