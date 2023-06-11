<template>

  <div
    class="item-card"
    ref="canvas"
    :key="forceRefresh"
    :style="{
        'transform': `scale(${this.zoom})`,
        'position': 'relative',
        'overflow': 'hidden',
        'width': '100%',
        'height': `100%`,
        'background-color': card.background || 'white'}"
    v-on:click="select($event)"
    @mousedown.self="card.selectedItem = -1"
  >

    <div v-if="card.active">
      <Moveable
        v-for="(item, index) in card.items"
        :key="index"
        :container="$refs.canvas"
        class="moveable"
        v-bind="card.settings(index)"
        @drag="handleDrag"
        @resize="handleResize"
        @rotate="handleRotate"
        @warp="handleWarp"
        @renderEnd="hideDataLabel"
        :style="item.position"
      >

        <div class="item-element" @mousedown.capture="selectItem(index)"
             style="width: 100%;"
             v-bind:class="{'item-element-cursor-move': !item.frozen}">
          <component :is="getCardItemName(item.type)" :item="item" :disabled="true"/>
        </div>
        <!--      <span>[size x size]</span>-->

      </Moveable>
    </div>

    <div v-else style="width: 100%;">
      <component
        v-for="(item, index) in card.items"
        v-if="item.enabled"
        class="item-element"
        @mousedown.capture="selectItem(index)"
        :is="getCardItemName(item.type)" :item="item"
        :style="item.position"
      />
    </div>
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
  name: 'EditorCard',
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

  private forceRefresh = 0;
  private zoom = 1;

  private getCardItemName(name: string): string {
    return CardItemName(name);
  }

  private created() {
    this.bus.$on('selected_card', (cardId: number) => {
      this.card.active = this.card.id == cardId;
    });
  }

  private destroyed() {

  }

  private select(event?: any) {
    // console.log('select_card', this.itemIndex)
    this.bus.$emit('selected_card', this.card.id);
  }

  private selectItem(index: number) {
    this.card.selectedItem = index;
    // console.log('selected', index);
  }

  handleDrag({target, transform, beforeTranslate, left, top}: any) {
    // console.log('onDrag', transform, 'left, top', left, top);
    this.card.items[this.card.selectedItem].transform = transform;
    target.style.transform = transform;
  }

  handleResize({target, width, height, clientX, clientY}: any) {
    // console.log('resize');
    width = Math.round(width);
    height = Math.round(height);

    this.card.items[this.card.selectedItem].width = width;
    this.card.items[this.card.selectedItem].height = height;
    target.style.width = `${width}px`;
    target.style.height = `${height}px`;
    // this.setDataLabel(clientX, clientY, `${width} x ${height}`);
  }

  handleRotate({target, transform, beforeRotate, clientX, clientY}: any) {
    this.card.items[this.card.selectedItem].transform = transform;
    target.style.transform = transform;
    // this.setDataLabel(clientX, clientY, `${beforeRotate}°`);
  }

  handleWarp({target, transform}: any) {
    // console.log('onWarp', transform);
    this.card.items[this.card.selectedItem].transform = transform;
    target.style.transform = transform;
  }

  private dataLabel = null;

  private hideDataLabel() {
    this.dataLabel = null;
  }
}
</script>

<style lang="scss">

.item-card {
  position: relative;
  width: 100%;
  height: 100%;
}

.item-card-title {
  position: relative;

  width: 100%;
}

.container {
  position: relative;
  top: 50px;
  left: 100px;

  /* border: 3px solid red; */
  overflow: scroll;
  width: 70%;
  height: 500px;
}

.data-label {
  position: absolute;
  z-index: 99999999;
  top: 0px;
  left: 0px;
  transform: translate(-200%, -200%);
  background-color: #555;
  color: white;
  padding: 5px 10px;
  border-radius: 3px;
  font-size: 12px;
  white-space: nowrap;
}

.item-element {
  /* border: 1px solid #aaa; */
  position: absolute;
  width: 100%;
  height: 100%;
  /*overflow: hidden;*/

  user-select: none;
}

.item-element-cursor-move {
  cursor: move;
}

.moveable {
  position: absolute;
  width: 300px;
  height: 200px;
  margin: 0 auto;
}

.moveable-control-box .moveable-line,
.moveable-control-box .moveable-control,
.moveable-control-box .moveable-rotation {
  background: #7474742e !important;
}

.moveable-line.moveable-rotation-line .moveable-control {
  border-color: #fff !important;
}

</style>
