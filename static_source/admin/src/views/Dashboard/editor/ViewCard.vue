<script setup lang="ts">
import {computed, onMounted, onUnmounted, onUpdated, PropType, ref,} from "vue";
import {Card, CardItem, Core} from "@/views/Dashboard/core";
import {useBus} from "@/views/Dashboard/bus";
import debounce from 'lodash.debounce'
import Moveable from 'vue3-moveable'
import {deepFlat} from "@daybrush/utils";
import {VueSelecto} from "vue3-selecto";
import {CardItemName} from "@/views/Dashboard/card_items";
import {UUID} from "uuid-generator-ts";
import {KeystrokeCaptureViewer} from "@/views/Dashboard/components";

const {emit} = useBus()

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
const cardRef = ref(null)

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
  set(val: Card) {
  }
})

const hover = ref(false)

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
      setSelectedTargets([])
      return
    }
    const target = currentCard.value.items[itemIndex].target;
    // target.classList.add("selected");
    setSelectedTargets([target]);
  }
})

useBus({
  name: 'unselected_card_item',
  callback: () => {
    if (currentCard.value.active) {
      return
    }
    setSelectedTargets([]);
  }
})

const selectCardItem = (itemIndex: number) => {
  if (!currentCard.value.active) {
    props.core?.onSelectedCard(currentCard.value.id)
    emit('selected_card', currentCard.value.id)
  }

  currentCard.value.selectedItem = itemIndex;
  // if (itemIndex === -1 || !currentCard.value.items.length || !currentCard.value.items[itemIndex]) {
  //   targets.value = [];
  // } else {
  //   targets.value = [currentCard.value.items[itemIndex].target];
  // }
  // emit('unselected_card_item')
}

const onDrag = ({target, transform, beforeTranslate, left, top}: any) => {
  const classes = target.className.split(' ');
  target.style.transform = transform;
  for (const cl of classes) {
    if (cl.includes('item-index-')) {
      const index = parseInt(cl.replace("item-index-", ""))
      currentCard.value.items[index].transform = transform;
    }
  }
}

const setSelectedTargets = (target) => {
  selectoRef.value.setSelectedTargets(deepFlat(target));
  targets.value = target;
};

const onResize = ({target, width, height, clientX, clientY}: any) => {
  width = Math.round(width);
  height = Math.round(height);

  if (currentCard.value.selectedItem > -1) {
    currentCard.value.items[currentCard.value.selectedItem].width = width;
    currentCard.value.items[currentCard.value.selectedItem].height = height;
  }
  target.style.width = `${width}px`;
  target.style.height = `${height}px`;
}

const onRotate = ({target, transform, beforeRotate, clientX, clientY}: any) => {
  if (currentCard.value.selectedItem > -1) {
    currentCard.value.items[currentCard.value.selectedItem].transform = transform;
  }
  target.style.transform = transform;
}

const targets = ref([]);
const moveableRef = ref(null);

const onSnap = e => {
  // console.log(e.guidelines, e.elements);
};

// ---------------------------------
// selecto methods
// ---------------------------------

const onDragGroup = ({events}) => {
  events.forEach(ev => {
    const classes = ev.target.className.split(' ');
    ev.target.style.transform = ev.transform;
    for (const cl of classes) {
      if (cl.includes('item-index-')) {
        const index = parseInt(cl.replace("item-index-", ""))
        currentCard.value.items[index].transform = ev.transform;
      }
    }
  });
};

const onRenderGroup = e => {
  e.events.forEach(ev => {
    ev.target.style.transform = ev.transform;
  });
};

const onClickGroup = e => {
  if (!e.moveableTarget) {
    setSelectedTargets([]);
    selectCardItem(-1);
    return;
  }
  if (e.isTrusted) {
    selectoRef.value.clickTarget(e.inputEvent, e.moveableTarget);
  }
};

// ---------------------------------
// group methods
// ---------------------------------

//todo add group
// https://daybrush.com/moveable/storybook/index.html?path=/story/combination-with-other-components--components-selecto-with-multiple-group
// https://daybrush.com/moveable/storybook/index.html?path=/story/combination-with-other-components--components-selecto

const selectoRef = ref(null);

const onSelect = (e) => {
  e.added.forEach(el => {
    el.classList.add("selected");
  });
  e.removed.forEach(el => {
    el.classList.remove("selected");
  });
}

const onSelectEnd = (e) => {
  const {
    isDragStartEnd,
    isClick,
    added,
    removed,
    inputEvent,
    selected,
  } = e;
  const moveable = moveableRef.value;
  if (e.isDragStart) {
    e.inputEvent.preventDefault();
    moveable.waitToChangeTarget().then(() => {
      //   moveable.dragStart(e.inputEvent);
    });
  }
  targets.value = selected;
  if (selected && selected.length == 1) {
    const classes = selected[0].className.split(' ');
    for (const cl of classes) {
      if (cl.includes('item-index-')) {
        const index = parseInt(cl.replace("item-index-", ""))
        selectCardItem(index)
      }
    }
  }

  if (selected && selected.length == 0) {
    selectCardItem(-1)
  }
};

const onDragStart = (e) => {
  const moveable = moveableRef.value;
  const target = e.inputEvent.target;
  const flatted = targets.value.flat(3);
  if (moveable.isMoveableElement(target)
      || flatted.some(t => t === target || t.contains(target))
  ) {
    e.stop();
  }
};

</script>

<template>

  <div
      class="item-card elements selecto-area prevent-select"
      ref="cardRef"
      v-bind:class="'class-'+currentCard.currentID"
      :key="reloadKey"
      :style="{
        'transform': `scale(${zoom})`,
        'background-color': currentCard.background || 'inherit'}"
      @mouseover="hover = true"
      @touchstart="hover = true"
      @mouseleave="hover = false"
      @mouseout="hover = false"
  >

    <KeystrokeCaptureViewer :card="currentCard" :core="core" :hover="hover"/>

    <component
        v-for="(item, index) in currentCard.items"
        :key="index"
        class="movable"
        :style="item.position"
        v-bind:class="['item-index-'+index, 'item-id-'+item.id]"
        :is="getCardItemName(item)"
        :item="item"
        :core="core"
        :editor="true"
    />

    <Moveable
        ref="moveableRef"
        :target="targets"
        @drag="onDrag"
        @dragGroup="onDragGroup"
        @renderGroup="onRenderGroup"
        @clickGroup="onClickGroup"
        @resize="onResize"
        @rotate="onRotate"
        @onSnap="onSnap"
        v-bind="currentCard.settings()"
    />

    <VueSelecto
        ref="selectoRef"
        :rootContainer="cardRef"
        :selectableTargets="['.class-'+currentCard.currentID+' .movable']"
        :hitRate="0"
        :selectByClick="true"
        :selectFromInside="false"
        :toggleContinueSelect="['shift']"
        :ratio="0"
        @dragStart="onDragStart"
        @selectEnd="onSelectEnd"
        @select="onSelect"
    />
  </div>

</template>

<style lang="less">
.movable {
  position: absolute;
}

.item-card {
  position: relative;
  overflow: hidden;
  width: 100%;
  height: 100%;
}
</style>
