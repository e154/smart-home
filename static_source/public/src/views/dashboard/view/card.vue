<template>

  <div
    class="item-card"
    :style="{
        'transform': `scale(${this.zoom})`,
        'position': 'relative',
        'overflow': 'hidden',
        'width': '100%',
        'height': '100%',
        'background-color': card.background || 'white'
  }"
  >

    <component
      v-for="(item, index) in card.items"
      class="item-element"
      :is="getCardItemName(item.type)" :item="item"
      :style="item.position"
    />
  </div>

</template>

<script lang="ts">
import {Component, Prop, Vue} from 'vue-property-decorator';
import {
  CardItemName,
  Dummy,
  IButton,
  IChart,
  IImage,
  ILogs,
  IProgress,
  IState,
  IText
} from '@/views/dashboard/card_items';
import {Card} from '@/views/dashboard/core';
import Moveable from 'vue-moveable';

class elementOption {
  public value = '';
  public label = '';
}

@Component({
  name: 'DashboardCard',
  components: {
    Moveable,
    Dummy,
    IText,
    IImage,
    IButton,
    IState,
    ILogs,
    IProgress,
    IChart
  }
})
export default class extends Vue {
  @Prop() private card!: Card;
  @Prop() private bus!: Vue;

  private zoom = 1;

  private getCardItemName(name: string): string {
    return CardItemName(name);
  }

  private created() {
  }

  private destroyed() {
  }
}
</script>

<style lang="scss">

.item-card {
  position: relative;
  width: 100%;
  height: 100%;
}

.item-element {
  position: absolute;
  width: 100%;
  height: 100%;
  /*overflow: hidden;*/
  /*user-select: none;*/
}

</style>
