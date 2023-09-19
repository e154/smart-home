<script setup lang="ts">
import {computed, onMounted, onBeforeUnmount, PropType, ref, unref, watch} from "vue";
import {Card, CardItem, Core, Tab} from "@/views/Dashboard/core";
import {RenderText} from "@/views/Dashboard/render";

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

const constraints = {
  audio: false,
  video: {
    width: {min: 1024, ideal: 1280, max: 1920},
    height: {min: 576, ideal: 720, max: 1080},
    facingMode: 'environment',
  }
}

const src = ref(null)



// ---------------------------------
// component methods
// ---------------------------------

const getUrl = (): string => {
  // const uri = import.meta.env.VITE_API_BASEPATH as string + '/stream/onvif.cam1/mse';
  let uri = import.meta.env.VITE_API_BASEPATH as string + '/stream/onvif.cam1/mse';
  uri = uri.replace("http", "ws")
  uri = uri.replace("https", "wss")
  console.log(uri)
  return uri;
}

const mseQueue = []
let mseSourceBuffer
let mseStreamingStarted = false

function startPlay () {
  const mse = new MediaSource()
  videoEl.value.src = window.URL.createObjectURL(mse)
  // console.log('---0')
  // mse.addEventListener('sourceopen', function () {
  //   console.log('---01')
    const ws = new WebSocket(getUrl())
    ws.binaryType = 'arraybuffer'
    ws.onopen = function (event) {
      console.log('Connect to ws')
    }
    ws.onmessage = function (event) {
      // console.log('---1')
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
  // }, false)
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

// fix stalled video in safari
// videoEl.addEventListener('pause', () => {
//   if (videoEl.currentTime > videoEl.buffered.end(videoEl.buffered.length - 1)) {
//     videoEl.currentTime = videoEl.buffered.end(videoEl.buffered.length - 1) - 0.1
//     videoEl.play()
//   }
// })

const el = ref(null)
onMounted(async () => {
  // store dom element moveable
  props.item.setTarget(el.value)

  startPlay()
})

onBeforeUnmount(() => {

})

</script>

<template>
  <div ref="el">
    <video ref="videoEl" autoplay muted playsinline controls style="max-width: 100%; max-height: 100%;"></video>
  </div>
</template>

<style lang="less" >

</style>
