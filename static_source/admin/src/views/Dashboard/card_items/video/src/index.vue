<script setup lang="ts">
import {onBeforeUnmount, onMounted, PropType, ref, watch} from "vue";
import {CardItem, requestCurrentState} from "@/views/Dashboard/core";
import VideoMse from "./VideoMse.vue";
import {debounce} from "lodash-es";
import {playerType} from "./types";
import VideoYou from "./VideoYou.vue";

// ---------------------------------
// common
// ---------------------------------

const props = defineProps({
  item: {
    type: Object as PropType<Nullable<CardItem>>,
    default: () => null
  },
})

const el = ref<ElRef>(null)
onMounted(async () => {

})

onBeforeUnmount(() => {

})

// ---------------------------------
// component methods
// ---------------------------------

const reloadKey = ref(0)
const reload = debounce(() => {
  reloadKey.value += 1
}, 500)

watch(
    () => props.item,
    (val?: CardItem) => {
      if (!val) return;
      reload()
    },
    {
      deep: true,
    }
)

requestCurrentState(props.item?.entityId);

</script>

<template>
  <div ref="el" :class="[{'hidden': item.hidden}]">
    <VideoMse :item="item" v-if="item && item.payload.video.playerType === playerType.onvifMse" :key="reloadKey"/>
    <VideoYou :item="item" v-if="item && item.payload.video.playerType === playerType.youtube" :key="reloadKey"/>
  </div>
</template>

<style lang="less">

</style>
