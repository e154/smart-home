<script setup lang="ts">
import {onBeforeUnmount, onMounted, PropType, ref} from "vue";
import {Cache, CardItem, RenderVar} from "@/views/Dashboard/core";
import {debounce} from "lodash-es";
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

// ---------------------------------
// component methods
// ---------------------------------

const videId = ref()


const _cache = new Cache()
const getVideoId = debounce(async () => {

  let token: string = props.item?.payload.video?.attribute || ''
  if (token) {
    const result = await RenderVar(token, props.item?.lastEvent)
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
