<script setup lang="ts">
import {onMounted, PropType, ref,} from "vue";
import {Card, CardItem, Core} from "@/views/Dashboard/core";
import {useBus} from "@/views/Dashboard/bus";
import {CardItemName} from "@/views/Dashboard/card_items";
import {UUID} from "uuid-generator-ts";
import {KeystrokeCaptureViewer} from "@/views/Dashboard/components";

const {bus} = useBus()

const currentID = ref('')
onMounted(() => {
  const uuid = new UUID()
  currentID.value = uuid.getDashFreeUUID()
})

// ---------------------------------
// common
// ---------------------------------

const zoom = ref(1);

const props = defineProps({
  core: {
    type: Object as PropType<Core>,
  },
  card: {
    type: Object as PropType<Nullable<Card>>,
    default: () => null
  },
})

const hover = ref(false)

// ---------------------------------
// component methods
// ---------------------------------
const getCardItemName = (item: CardItem): string => {
  //todo: check if item disabled
  return CardItemName(item.type);
}

</script>

<template>

  <div
      class="item-card elements selecto-area"
      v-bind:class="'class-'+card.currentID"
      :style="{
        'transform': `scale(${zoom})`,
        'background-color': card.background || 'inherit'}"
      @mouseover="hover = true"
      @touchstart="hover = true"
      @mouseleave="hover = false"
      @mouseout="hover = false"
  >
    <KeystrokeCaptureViewer :card="card" :core="core" :hover="hover"/>
    <component
        v-for="(item, index) in card.items"
        :key="index"
        class="item-element"
        :style="item.position"
        :is="getCardItemName(item)"
        :item="item"
        :core="core"
    />
  </div>

</template>

<style lang="less">

.item-card {
  position: relative;
  overflow: hidden;
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
