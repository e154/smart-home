<script setup lang="ts">
import { onMounted, onBeforeUnmount, PropType, ref} from "vue";
import { CardItem } from "@/views/Dashboard/core";
import {Websocket, WebsocketBuilder} from "websocket-ts";

// ---------------------------------
// common
// ---------------------------------

const props = defineProps({
  item: {
    type: Object as PropType<Nullable<CardItem>>,
    default: () => null
  },
})

const videoEl = ref()

onMounted(async () => {

  // fix stalled video in safari
  videoEl.value.addEventListener('pause', () => {
    if (!videoEl.value) {
      return
    }
    if (videoEl.value.currentTime > videoEl.value.buffered.end(videoEl.value.buffered.length - 1)) {
      videoEl.value.currentTime = videoEl.value.buffered.end(videoEl.value.buffered.length - 1) - 0.1
      videoEl.value.play()
    }
  })

  startPlay()
})

onBeforeUnmount(() => {
  stopPlay()
})

// ---------------------------------
// component methods
// ---------------------------------

const getUrl = (): string => {
  if (!props.item?.entityId) {
    return ""
  }
  let uri = (import.meta.env.VITE_API_BASEPATH as string || '/') + '/stream/'+ props.item.entityId +'/mse';
  uri = uri.replace("http", "ws")
  uri = uri.replace("https", "wss")
  return uri;
}

const mseQueue = []
let mseSourceBuffer
let mseStreamingStarted = false
let ws: Websocket;

const startPlay = () => {
  if (!props.item?.entityId) {
    return
  }
  const mse = new MediaSource()
  videoEl.value.src = window.URL.createObjectURL(mse)
  mse.addEventListener('sourceopen', function () {
    ws = new WebSocket(getUrl())
    ws.binaryType = 'arraybuffer'
    ws.onopen = function (event) {
      console.log('Connect to ws')
    }
    ws.onmessage = function (event) {
      const data = new Uint8Array(event.data)
      if (data[0] === 9) {
        let mimeCodec
        const decodedArr = data.slice(1)
        if (window.TextDecoder) {
          mimeCodec = new TextDecoder('utf-8').decode(decodedArr)
        } else {
          mimeCodec = Utf8ArrayToStr(decodedArr)
        }
        mseSourceBuffer = mse.addSourceBuffer('video/mp4; codecs="' + mimeCodec + '"')
        mseSourceBuffer.mode = 'segments'
        mseSourceBuffer.addEventListener('updateend', pushPacket)
      } else {
        readPacket(event.data)
      }
    }
  }, false)
}

const stopPlay = () => {
  if (!ws) {
    return
  }
  ws.close()
}

const pushPacket = () => {
  let packet

  if (!mseSourceBuffer.updating) {
    if (mseQueue.length > 0) {
      packet = mseQueue.shift()
      mseSourceBuffer.appendBuffer(packet)
    } else {
      mseStreamingStarted = false
    }
  }
  if (!videoEl.value) {
    return
  }
  if (videoEl.value.buffered.length > 0) {
    if (typeof document.hidden !== 'undefined' && document.hidden) {
      // no sound, browser paused video without sound in background
      videoEl.value.currentTime = videoEl.value.buffered.end((videoEl.value.buffered.length - 1)) - 0.5
    }
  }
}

const readPacket = (packet) => {
  if (!mseStreamingStarted) {
    mseSourceBuffer.appendBuffer(packet)
    mseStreamingStarted = true
    return
  }
  mseQueue.push(packet)
  if (!mseSourceBuffer.updating) {
    pushPacket()
  }
}

</script>

<template>
  <video ref="videoEl" autoplay playsinline controls style="max-width: 100%; max-height: 100%; height: 100%; width: 100%"></video>
</template>

<style lang="less" >

</style>
