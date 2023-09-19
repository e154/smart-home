<script setup lang="ts">
import {computed, onMounted, onBeforeUnmount, PropType, ref, unref, watch} from "vue";
import {Card, CardItem, Core, Tab} from "@/views/Dashboard/core";
import VideoMse from "@/views/Dashboard/card_items/stream_player/VideoMse.vue";

// ---------------------------------
// common
// ---------------------------------

const props = defineProps({
  item: {
    type: Object as PropType<Nullable<CardItem>>,
    default: () => null
  },
})

const el = ref(null)
onMounted(async () => {
  // store dom element moveable
  props.item.setTarget(el.value)

})

onBeforeUnmount(() => {

})

// ---------------------------------
// component methods
// ---------------------------------

const reloadKey = ref(0)
const reload = () => {
  reloadKey.value += 1
}

watch(
    () => props.item?.entityId,
    (val?: string) => {
      console.log('----')
      reload()
    },
    {
      deep: true,
    }
)

</script>

<template>
  <div ref="el">
    <VideoMse :item="item" v-if="item" :key="reloadKey"/>
  </div>
</template>

<style lang="less" >

</style>
