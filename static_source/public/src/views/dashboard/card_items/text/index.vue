<template>
  <transition name="fade">
    <div v-if="item.asButton"
         v-show="!item.hidden"
         @mouseover="mouseOver"
         @mouseleave="mouseLive()"
         class="device-menu"
         :class="[{'as-button': item.asButton && item.buttonActions.length > 0}]"
    >

      <div
        class="cursor-pointer"
        :style="item.style"
        v-html="value"
        :key="reloadKey"
        @click.prevent.stop="callAction(item.buttonActions[0])"></div>

      <div
        :class="[{'show': showMenu}]"
        class="device-menu-circle"
        v-if="item.asButton && item.buttonActions.length > 1"
      >
        <a href="#"
           v-for="action in item.buttonActions"
           @click.prevent.stop="callAction(action)">
          <img :src="item.getUrl(action.image)"/>
        </a>
      </div>
    </div>
    <div v-else
         :style="item.style"
         v-show="!item.hidden"
         v-html="value"
         :key="reloadKey"></div>
  </transition>
</template>

<script lang="ts">
import {Component, Prop, Vue, Watch} from 'vue-property-decorator'
import {ButtonAction, CardItem, requestCurrentState} from '@/views/dashboard/core';
import {Cache, Compare, GetTokens, RenderText, Resolve} from '@/views/dashboard/render'
import {Attribute, GetAttrValue} from '@/api/stream_types'
import api from "@/api/api";

@Component({
  name: 'IText',
  components: {}
})
export default class extends Vue {
  @Prop() private item?: CardItem;
  private value = '';
  private _cache!: Cache;
  private reloadKey = 0;
  private showMenu = false;

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

  @Watch('item', {deep: true})
  private onUpdateItem(item: CardItem) {
    this.update()
  }

  // todo fix
  resetCache() {
    this._cache.clear()
  }

  private timer: any;

  private mouseLive() {
    if (!this.showMenu) {
      return;
    }
    this.timer = setTimeout(() => {
      this.showMenu = false;
      this.timer = null;
    }, 2000);
  }

  private mouseOver() {
    this.showMenu = true;
    if (!this.timer) {
      return;
    }
    clearTimeout(this.timer);
  }

  private async callAction(action: ButtonAction) {
    if (!action) {
      return;
    }
    await api.v1.interactServiceEntityCallAction({
      id: action.entityId,
      name: action.action || ''
    });
    this.$notify({
      title: 'Success',
      message: 'Call Successfully',
      type: 'success',
      duration: 2000
    });
  }
}
</script>

<style>
.ql-align-center {
  text-align: center;
}
.cursor-pointer {
  cursor: pointer;
}
</style>
