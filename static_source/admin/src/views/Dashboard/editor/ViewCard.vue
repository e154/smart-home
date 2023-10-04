<script setup lang="ts">
import {
  computed,
  defineAsyncComponent,
  defineComponent,
  onBeforeUpdate,
  onMounted, onUnmounted,
  onUpdated,
  PropType,
  ref,
  watch
} from "vue";
import {Card, CardItem, Core, Tab} from "@/views/Dashboard/core";
import {Vuuri} from "@/views/Dashboard/Vuuri"
import {useBus} from "@/views/Dashboard/bus";
import debounce from 'lodash.debounce'
import Moveable, { OnRender } from 'vue3-moveable'
import { VueSelecto } from "vue3-selecto";
import {CardItemName} from "@/views/Dashboard/card_items";
import {UUID} from "uuid-generator-ts";

const {bus} = useBus()

const currentID = ref('')
onMounted(() => {
  const uuid = new UUID()
  currentID.value = uuid.getDashFreeUUID()

  currentCard.value.document = document
  currentCard.value.updateItemList()
})

onUnmounted(() => {

})

onUpdated(() => {

})

// ---------------------------------
// common
// ---------------------------------

const zoom = ref(1);
const canvas = ref(null)

const props = defineProps({
  core: {
    type: Object as PropType<Core>,
  },
  card: {
    type: Object as PropType<Nullable<Card>>,
    default: () => null
  },
})

const reloadKey = ref(0);
const reload = debounce(() => {
  reloadKey.value += 1
}, 100)


const currentCard = computed({
  get(): Card {
    return props.card as Card
  },
  set(val: Card) {}
})

// ---------------------------------
// component methods
// ---------------------------------
const getCardItemName = (item: CardItem): string => {
  //todo: check if item disabled
  return CardItemName(item.type);
}

useBus({
  name: 'selected_card_item',
  callback: (itemIndex: number) => {
    if (!currentCard.value.active) {
      return
    }
    if (itemIndex === -1 || !currentCard.value.items.length || !currentCard.value.items[itemIndex]) {
      targets.value = [];
      return
    }
    targets.value = [currentCard.value.items[itemIndex].target];
  }
})


useBus({
  name: 'unselected_card_item',
  callback: () => {
    if (currentCard.value.active) {
      return
    }

    targets.value = [];
  }
})

const selectCard = (event?: any) => {
  if (currentCard.value.active) return;
  props.core?.onSelectedCard(currentCard.value.id)
  bus.emit('selected_card', currentCard.value.id)
}

const selectCardItem = (itemIndex: number) => {
  if (!currentCard.value.active) {
    props.core?.onSelectedCard(currentCard.value.id)
    bus.emit('selected_card', currentCard.value.id)
  }

  currentCard.value.selectedItem = itemIndex;
  if (itemIndex === -1 || !currentCard.value.items.length || !currentCard.value.items[itemIndex]) {
    targets.value = [];
  } else {
    targets.value = [currentCard.value.items[itemIndex].target];
  }
  bus.emit('unselected_card_item')
}

const onRender = ({target, cssText}) => {
  target.style.cssText += cssText;
}

const onDrag = ({target, transform, beforeTranslate, left, top}: any) => {
  if (currentCard.value.selectedItem > -1) {
    currentCard.value.items[currentCard.value.selectedItem].transform = transform; //todo uncomment
  }
  target.style.transform = transform;
}

const onResize = ({target, width, height, clientX, clientY}: any) => {
  width = Math.round(width);
  height = Math.round(height);

  if (currentCard.value.selectedItem > -1) {
    currentCard.value.items[currentCard.value.selectedItem].width = width; //todo uncomment
    currentCard.value.items[currentCard.value.selectedItem].height = height; //todo uncomment
  }
  target.style.width = `${width}px`;
  target.style.height = `${height}px`;
}

const onRotate = ({target, transform, beforeRotate, clientX, clientY}: any) => {
  if (currentCard.value.selectedItem > -1) {
    currentCard.value.items[currentCard.value.selectedItem].transform = transform; //todo uncomment
  }
  target.style.transform = transform;
}

const handleWarp = ({target, transform}: any) => {
  // console.log('onWarp', transform);
  if (currentCard.value.selectedItem > -1) {
    currentCard.value.items[currentCard.value.selectedItem].transform = transform;
  }
  target.style.transform = transform;
}


const targets = ref([]);
const moveableRef = ref(null);
const selectoRef = ref(null);
const cubes = [];
for (let i = 0; i < 30; ++i) {
    cubes.push(i);
}

const onDragStart = (e) => {
  const moveable = moveableRef.value;
  const target = e.inputEvent.target;
  if (moveable.isMoveableElement(target)
        || targets.value.some(t => t === target || t.contains(target))
    ) {
        e.stop();
    }
};

const onSelectEnd = (e) => {
  const moveable = moveableRef.value;
  if (e.isDragStart) {
      e.inputEvent.preventDefault();
      moveable.waitToChangeTarget().then(() => {
          moveable.dragStart(e.inputEvent);
      });
  }
  targets.value = e.selected;
};

const onSnap = e => {
  // console.log(e.guidelines, e.elements);
};


</script>

<template>

<!--  <VueSelecto-->
<!--      ref="selectoRef"-->
<!--      :rootContainer="canvas"-->
<!--      :selectableTargets="['.class-'+currentCard.currentID+' .movable']"-->
<!--      :hitRate="0"-->
<!--      :selectByClick="true"-->
<!--      :selectFromInside="false"-->
<!--      :toggleContinueSelect="['shift']"-->
<!--      :ratio="0"-->
<!--      @dragStart="onDragStart"-->
<!--      @selectEnd="onSelectEnd"-->
<!--  />-->

  <div
      class="item-card elements selecto-area"
      ref="canvas"
      v-bind:class="'class-'+currentCard.currentID"
      :key="reloadKey"
      :style="{
        'transform': `scale(${zoom})`,
        'position': 'relative',
        'overflow': 'hidden',
        'width': '100%',
        'height': `100%`,
        'background-color': currentCard.background || 'inherit'}"
      @click="selectCard()"
      @mousedown.self="selectCardItem(-1)"
  >

    <component
        v-for="(item, index) in currentCard.items"
        :key="index"
        class="movable"
        :style="item.position"
        :is="getCardItemName(item)"
        :item="item"
        :core="core"
        :editor="true"
        @mousedown.capture="selectCardItem(index)"
    />

    <Moveable
        ref="moveableRef"
        :draggable="true"
        :target="targets"
        @drag="onDrag"
        @resize="onResize"
        @rotate="onRotate"
        @onSnap="onSnap"
        v-bind="currentCard.settings()"
    />
  </div>

</template>

<style lang="less">
.movable {
  position: absolute;
}

</style>
