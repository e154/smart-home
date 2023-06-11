<template>
  <transition name="fade">
    <div
      :style="item.style"
      v-show="!item.hidden"
      v-html="value"
      :key="reloadKey"></div>
  </transition>
</template>

<script lang="ts">
import { Component, Prop, Vue, Watch } from 'vue-property-decorator'
import {CardItem, requestCurrentState} from '@/views/dashboard/core';
import { Cache, Compare, GetTokens, RenderText, Resolve } from '@/views/dashboard/render'
import { Attribute, GetAttrValue } from '@/api/stream_types'

@Component({
  name: 'IText',
  components: {}
})
export default class extends Vue {
  @Prop() private item?: CardItem;
  private value = '';
  private _cache!: Cache;
  private reloadKey = 0;

  reload() {
    this.reloadKey += 1
  }

  private created() {
    this._cache = new Cache()
    this.update()
    requestCurrentState(this.item?.entityId!);
  }

  private mounted() {
  }

  private update(): void {
    // console.log('update value', this.item?.payload?.text);

    if (!this.item?.payload.text?.items) {
      this.value = this.item?.payload.text?.default_text || ''
      return
    }

    let value = this.item?.payload.text?.default_text || ''
    let value2 = ''

    for (const prop of this.item?.payload.text?.items) {
      // select prop
      let val = Resolve(prop.key, this.item?.lastEvent)
      if (!val) {
        continue
      }

      if (typeof val === 'object') {
        if (val && val.hasOwnProperty('type') && val.hasOwnProperty('name')) {
          val = GetAttrValue(val as Attribute)
        }
      }

      if (val == undefined) {
        val = '[NO VALUE]'
      }

      const tr = Compare(val, prop.value, prop.comparison)
      if (!tr) {
        continue
      }

      if (!prop.tokens) {
        prop.tokens = []
      }

      // render text
      prop.tokens = GetTokens(prop.text, this._cache)
      if (!prop.tokens.length) {
        this.value = prop.text || ''
        return
      }

      if (prop.text) {
        value2 = prop.text
      }

      value2 = RenderText(prop.tokens, value2, this.item?.lastEvent)

      this.value = value2 || value
      return
    }

    const tokens = GetTokens(value, this._cache)
    if (tokens) {
      value = RenderText(tokens, value, this.item?.lastEvent)
    }
    this.value = value
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

<style>
.ql-align-center {
  text-align: center;
}
</style>
