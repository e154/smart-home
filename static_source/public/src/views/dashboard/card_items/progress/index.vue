<template>
  <div>
    <el-progress
      v-if="item.payload.progress.type"
      :type="item.payload.progress.type"
      :percentage="value"
      :width="item.payload.progress.width"
      :stroke-width="item.payload.progress.strokeWidth"
      :text-inside="!item.payload.progress.textInside"></el-progress>
    <el-progress
      v-else
      :percentage="value"
      :width="item.payload.progress.width"
      :stroke-width="item.payload.progress.strokeWidth"
      :text-inside="!item.payload.progress.textInside"></el-progress>
  </div>

</template>

<script lang="ts">
import { Component, Prop, Vue, Watch } from 'vue-property-decorator'
import {CardItem, requestCurrentState} from '@/views/dashboard/core';
import { Cache, GetTokens, RenderText } from '@/views/dashboard/render'

@Component({
  name: 'Progress',
  components: {
  }
})
export default class extends Vue {
  @Prop() private item?: CardItem;
  private value = 0;
  private _cache!: Cache;

  private created() {
    this._cache = new Cache()
    // this.update()
    requestCurrentState(this.item?.entityId);
  }

  private mounted() {}

  private update(): void {
    let value: string = this.item?.payload.progress?.value || ''
    const tokens = GetTokens(value, this._cache)
    if (tokens) {
      value = RenderText(tokens, value, this.item?.lastEvent)
    }
    this.value = parseInt(value) || 0
  }

  @Watch('item', { deep: true })
  private onUpdateItem(item: CardItem) {
    this.update()
  }

  // todo fix
  resetCache() {
    this._cache.clear()
  }
}
</script>

<style scoped>

</style>
