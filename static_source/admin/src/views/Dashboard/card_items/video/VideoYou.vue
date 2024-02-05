<script setup lang="ts">
import {computed, onBeforeUnmount, onMounted, PropType, ref} from "vue";
import {CardItem} from "@/views/Dashboard/core";
import {Cache, RenderVar} from "@/views/Dashboard/render";
import {debounce} from "lodash-es";
import {ItemPayloadSlider} from "@/views/Dashboard/card_items/slider/types";
import {ItemPayloadVideo} from "@/views/Dashboard/card_items/video/types";
import LiteYouTubeEmbed from 'vue-lite-youtube-embed'
import 'vue-lite-youtube-embed/style.css'

// ---------------------------------
// common
// ---------------------------------

const props = defineProps({
  item: {
    type: Object as PropType<Nullable<CardItem>>,
    default: () => null
  },
})

const iframe = ref()

onMounted(async () => {

  iframe.value.warmConnections()
  iframe.value.addIframe()

  setTimeout(() => {

  }, 2000);
})

onBeforeUnmount(() => {

})

const currentVideo = computed<ItemPayloadVideo>(() => props.item?.payload.video || {} as ItemPayloadSlider)

// ---------------------------------
// component methods
// ---------------------------------

const videId = ref()


const _cache = new Cache()
const getVideoId = debounce(() => {

  let token: string = props.item?.payload.video?.attribute || ''
  if (token) {
    const result = RenderVar(token, props.item?.lastEvent)
    videId.value = result
  }
})

getVideoId()

</script>

<template>
  <LiteYouTubeEmbed
      ref="iframe"
      :id="videId"
      :muted="true"
      title="youtube"/>
</template>

<style lang="less">

</style>
