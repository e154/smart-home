<script setup lang="ts">

import {computed, onMounted, PropType} from "vue";
import {FrameProp} from "@/views/Dashboard/components";
import {GetFullImageUrl} from "@/utils/serverId";

// ---------------------------------
// common
// ---------------------------------

const props = defineProps({
  frame: {
    type: Object as PropType<Nullable<FrameProp>>,
    default: () => null
  },
})

// ---------------------------------
// component methods
// ---------------------------------

onMounted(() => {

})

const imageUrl = computed(() => props.frame?.image ? `url(${GetFullImageUrl(props.frame?.image)})` : '')

// top
// ---------------------------------
const topLeftCornerStyle = computed(() => {
  if (!props.frame?.items || !props.frame?.items['top-left-corner']) {
    return {
      display: 'none'
    }
  }

  return {
    width: `${props.frame.items['top-left-corner'].width}px`,
    height: `${props.frame.items['top-left-corner'].height}px`,
    backgroundPosition: `-${props.frame.items['top-left-corner'].x}px -${props.frame.items['top-left-corner'].y}px`
  }
})

const topStyle = computed(() => {
  if (!props.frame?.items || !props.frame?.items['top']) {
    return {
      display: 'none'
    }
  }
  const result = {
    left: `${props.frame.items['top-left-corner'].width || 0}px`,
    right: `${props.frame.items['top-right-corner'].width || 0}px`,
    height: `${props.frame.items['top'].height}px`,
  }
  if (props.frame.items['top'].base64) {
    result['background-image'] = `url(${props.frame.items['top'].base64})`
  }

  return result
})

const topRightCornerStyle = computed(() => {
  if (!props.frame?.items || !props.frame?.items['top-right-corner']) {
    return {
      display: 'none'
    }
  }
  return {
    width: `${props.frame.items['top-right-corner'].width}px`,
    height: `${props.frame.items['top-right-corner'].height}px`,
    backgroundPosition: `-${props.frame.items['top-right-corner'].x}px -${props.frame.items['top-right-corner'].y}px`
  }
})

// middle
// ---------------------------------

const leftStyle = computed(() => {
  if (!props.frame?.items || !props.frame?.items['left']) {
    return {
      display: 'none'
    }
  }
  const result = {
    top: `${props.frame.items['top-left-corner'].height || 0}px`,
    bottom: `${props.frame.items['bottom-left-corner'].height || 0}px`,
    width: `${props.frame.items['left'].width}px`,
  }
  if (props.frame?.items['left'].base64) {
    result['background-image'] = `url(${props.frame.items['left'].base64})`
  }

  return result
})

const contentStyle = computed(() => {
  if (!props.frame?.items || !props.frame?.items['content']) {
    return {
      display: 'none'
    }
  }

  if (props.frame?.items['content'].base64) {
    return {
      backgroundImage: `url(${props.frame.items['content'].base64})`
    }
  }

  return {}
})

const rightStyle = computed(() => {
  if (!props.frame?.items || !props.frame?.items['right']) {
    return {
      display: 'none'
    }
  }
  const result = {
    top: `${props.frame.items['top-right-corner'].height || 0}px`,
    bottom: `${props.frame.items['bottom-right-corner'].height || 0}px`,
    width: `${props.frame.items['right'].width}px`,
  }
  if (props.frame.items['right'].base64) {
    result['background-image'] = `url(${props.frame.items['right'].base64})`
  }

  return result
})

// bottom
// ---------------------------------
const bottomLeftCornerStyle = computed(() => {
  if (!props.frame?.items || !props.frame?.items['bottom-left-corner']) {
    return {
      display: 'none'
    }
  }
  return {
    width: `${props.frame.items['bottom-left-corner'].width}px`,
    height: `${props.frame.items['bottom-left-corner'].height}px`,
    backgroundPosition: `-${props.frame.items['bottom-left-corner'].x}px -${props.frame.items['bottom-left-corner'].y}px`
  }
})

const bottomStyle = computed(() => {
  if (!props.frame?.items || !props.frame?.items['bottom']) {
    return {
      display: 'none'
    }
  }
  const result = {
    left: `${props.frame.items['bottom-left-corner'].width || 0}px`,
    right: `${props.frame.items['bottom-right-corner'].width || 0}px`,
    height: `${props.frame.items['bottom'].height}px`,
  }
  if (props.frame.items['bottom'].base64) {
    result['background-image'] = `url(${props.frame.items['bottom'].base64})`
  }

  return result
})

const bottomRightStyle = computed(() => {
  if (!props.frame?.items || !props.frame?.items['bottom-right-corner']) {
    return {
      display: 'none'
    }
  }
  return {
    width: `${props.frame.items['bottom-right-corner'].width}px`,
    height: `${props.frame.items['bottom-right-corner'].height}px`,
    backgroundPosition: `-${props.frame.items['bottom-right-corner'].x}px -${props.frame.items['bottom-right-corner'].y}px`
  }
})

</script>

<template>
  <div class="window" v-if="frame?.image">
    <div class="base" :style="contentStyle"></div>
    <div class="top-left-corner" :style="topLeftCornerStyle"></div>
    <div class="top" :style="topStyle"></div>
    <div class="top-right-corner" :style="topRightCornerStyle"></div>
    <div class="left" :style="leftStyle"></div>
    <div class="right" :style="rightStyle"></div>
    <div class="bottom-left-corner" :style="bottomLeftCornerStyle"></div>
    <div class="bottom" :style="bottomStyle"></div>
    <div class="bottom-right-corner" :style="bottomRightStyle"></div>
    <div class="content">
      <slot></slot>
    </div>
  </div>
  <div class="window" v-else>
    <slot></slot>
  </div>
</template>

<style scoped lang="less">
.window {
  width: 100%;
  height: 100%;
  position: relative;
}

.top-left-corner,
.top-right-corner,
.bottom-left-corner,
.bottom-right-corner,
.top,
.bottom,
.left,
.right,
.base {
  position: absolute;
  background-image: v-bind(imageUrl);
  background-repeat: no-repeat;
}

.top-left-corner {
  top: 0;
  left: 0;
}

.top-right-corner {
  top: 0;
  right: 0;
}

.bottom-left-corner {
  bottom: 0;
  left: 0;
}

.bottom-right-corner {
  bottom: 0;
  right: 0;
}

.top {
  top: 0;
  background-repeat: repeat-x;
}

.bottom {
  bottom: 0;
  background-repeat: repeat-x;
}

.left {
  left: 0;
  background-repeat: repeat-y;
}

.right {
  right: 0;
  background-repeat: repeat-y;
}

.base {
  background-repeat: repeat;
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
}

.content, base {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
}
</style>
