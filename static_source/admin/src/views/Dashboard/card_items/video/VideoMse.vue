<script setup lang="ts">
import {onMounted, onBeforeUnmount, PropType, ref, onUnmounted} from "vue";
import { CardItem } from "@/views/Dashboard/core";
import {Websocket, WebsocketBuilder} from "websocket-ts";
import {ApiSigninResponse} from "@/api/stub";
import {useCache} from "@/hooks/web/useCache";
const {wsCache} = useCache()

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

onUnmounted(() => {
  stopPlay()
})

// ---------------------------------
// component methods
// ---------------------------------

const getUrl = (): string => {
  if (!props.item?.entityId) {
    return ""
  }
  //todo: add channel select
  const accessToken = wsCache.get("accessToken")
  let uri = import.meta.env.VITE_API_BASEPATH as string || window.location.origin;
  uri = uri + '/stream/'+ props.item.entityId +'/channel/0/mse?access_token=' + accessToken;
  uri = uri.replace("https", "wss")
  uri = uri.replace("http", "ws")
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
  if (!videoEl.value) {
    return;
  }
  videoEl.value.src = window.URL.createObjectURL(mse)
  mse.addEventListener('sourceopen', function () {
    ws = new WebSocket(getUrl())
    ws.binaryType = 'arraybuffer'
    ws.onopen = function (event) {
      // console.log('Connect to ws')
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
    ws.onclose = function(e) {
      // console.log('Socket is closed. Reconnect will be attempted in 1 second.', e.reason);
      setTimeout(function() {
        startPlay();
      }, 1000);
    };
  }, false)
}

const stopPlay = () => {
  if (!ws) {
    return
  }
  ws.close()
  ws = null
  videoEl.value = null
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
