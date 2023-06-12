<template>
  <transition name="fade">
    <div v-if="item.asButton"
         v-show="!item.hidden"
         @mouseover="mouseOver"
         @mouseleave="mouseLive()"
         class="device-menu"
         :class="[{'as-button': item.asButton && item.buttonActions.length > 0}]"
    >

      <img class="device"
           style="width: 100%"
           :key="reloadKey"
           :src="item.getUrl(internal.currentImage)"
           @click.prevent.stop="callAction(item.buttonActions[0])">

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
    <div v-else v-show="!item.hidden">
      <img class="device"
           style="width: 100%"
           :key="reloadKey"
           :src="item.getUrl(internal.currentImage)"
      >
    </div>
  </transition>
</template>

<script lang="ts">
import {Component, Prop, Vue, Watch} from 'vue-property-decorator';
import {ButtonAction, CardItem, requestCurrentState} from '@/views/dashboard/core';
import {Compare, Resolve} from '@/views/dashboard/render';
import {Attribute, GetAttrValue} from '@/api/stream_types';
import {ApiImage} from '@/api/stub';
import api from '@/api/api';

@Component({
  name: 'IState',
  components: {}
})
export default class extends Vue {
  @Prop() private item?: CardItem;
  private reloadKey = 0;
  private showMenu = false;

  private internal: { currentImage: ApiImage | undefined } = {
    currentImage: undefined
  };

  private created() {
    if (!this.item?.payload.state || !this.item?.payload.state.default_image) {
      return;
    }
    this.internal.currentImage = this.item?.payload.state.default_image;

    // setTimeout(() => {
    requestCurrentState(this.item?.entityId);
    // }, 1000);
    // this.update()
  }

  reload() {
    this.reloadKey += 1;
  }

  private mounted() {
  }

  private setImage(image: ApiImage | undefined) {
    this.internal.currentImage = image;
    // TODO fix reload
    // this.reload();
  }

  private update() {
    let counter = 0;

    if (this.item?.payload.state?.items) {
      for (const prop of this.item?.payload.state?.items) {
        let val = Resolve(prop.key, this.item?.lastEvent);
        if (!val) {
          continue;
        }

        if (typeof val === 'object') {
          if (val && val.hasOwnProperty('type') && val.hasOwnProperty('name')) {
            val = GetAttrValue(val as Attribute);
          }
        }

        if (val == undefined) {
          val = '[NO VALUE]';
        }

        const tr = Compare(val, prop.value, prop.comparison);
        if (tr && prop.image) {
          counter++;
          this.setImage(prop.image);
        }
      }
    }

    if (counter == 0) {
      this.setImage(this.item?.payload?.state?.default_image);
    }
  }

  @Watch('item', {deep: true})
  private onUpdateItem(item: CardItem) {
    this.update();
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

</style>
