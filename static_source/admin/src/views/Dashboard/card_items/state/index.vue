<script setup lang="ts">
import {computed, onMounted, PropType, ref, unref, watch} from "vue";
import {ButtonAction, Card, CardItem, Core, requestCurrentState, Tab} from "@/views/Dashboard/core";
import {Cache, Compare, Resolve} from "@/views/Dashboard/render";
import {ApiImage} from "@/api/stub";
import api from "@/api/api";
import {ElMessage} from "element-plus";
import {Attribute, GetAttrValue} from "@/api/stream_types";
import {debounce} from "lodash-es";
import {useI18n} from "@/hooks/web/useI18n";

const {t} = useI18n()

// ---------------------------------
// common
// ---------------------------------
const reloadKey = ref(0)
const _cache = new Cache()
const showMenu = ref(false)
const currentImage = ref<Nullable<ApiImage>>(null)

const props = defineProps({
  item: {
    type: Object as PropType<Nullable<CardItem>>,
    default: () => null
  },
})

const el = ref(null)
onMounted(() => {
  // store dom element moveable
  props.item.setTarget(el.value)

  if (props.item?.payload.state.default_image) {
    currentImage.value = props.item?.payload.state.default_image;
  }

})

// ---------------------------------
// component methods
// ---------------------------------

const reload = () => {
  reloadKey.value += 1
}

const update = debounce(() => {
  let counter = 0;

  if (props.item?.payload.state?.items) {
    for (const prop of props.item?.payload.state?.items) {
      let val = Resolve(prop.key, props.item?.lastEvent);
      if (!val) {
        continue;
      }

      if (typeof val === 'object') {
        if (val && val.hasOwnProperty('type') && val.hasOwnProperty('name')) {
          val = GetAttrValue(val as Attribute);
        }
      }

      if (val == undefined) {
        val = '[NO VALUE]';
      }

      const tr = Compare(val, prop.value, prop.comparison);
      if (tr && prop.image) {
        counter++;
        currentImage.value = prop.image;
      }
    }
  }

  if (counter == 0) {
    currentImage.value = props.item?.payload?.state?.default_image;
  }
}, 100)

watch(
    () => props.item,
    (val?: CardItem) => {
      if (!val) return;
      update()
    },
    {
      deep: true,
      immediate: true
    }
)

let timer: any;
const mouseLive = () => {
  if (!showMenu.value) {
    return;
  }
  timer = setTimeout(() => {
    showMenu.value = false;
    timer = null;
  }, 2000);
}

const mouseOver = () => {
  showMenu.value = true;
  if (!timer) {
    return;
  }
  clearTimeout(timer);
}

const callAction = async (action: ButtonAction) => {
  if (!action) {
    return;
  }
  await api.v1.interactServiceEntityCallAction({
    id: action.entityId,
    name: action.action || ''
  });
  ElMessage({
    title: t('Success'),
    message: t('message.callSuccessful'),
    type: 'success',
    duration: 2000
  });
}



update()

requestCurrentState(props.item?.entityId);

</script>

<template>
  <div ref="el">
    <div
        v-if="item.asButton"
         v-show="!item.hidden"
         @mouseover="mouseOver"
         @mouseleave="mouseLive()"
         class="device-menu"
         :class="[{'as-button': item.asButton && item.buttonActions.length > 0}]"
    >

      <img
          class="device"
           style="width: 100%"
           :key="reloadKey"
           :src="item.getUrl(currentImage)"
           @click.prevent.stop="callAction(item.buttonActions[0])"/>

      <div
          :class="[{'show': showMenu}]"
          class="device-menu-circle"
          v-if="item.asButton && item.buttonActions.length > 1"
      >
        <a
            href="#"
           v-for="(action, index) in item.buttonActions"
           @click.prevent.stop="callAction(action)"
          :key="index">
          <img :src="item.getUrl(action.image)"/>
        </a>
      </div>
    </div>
    <div v-else v-show="!item.hidden">
      <img
          class="device"
           style="width: 100%"
           :key="reloadKey"
           :src="item.getUrl(currentImage)"/>
    </div>
  </div>
</template>

<style lang="less" scoped>

.device-menu {
  position: relative;
  transition: all 0.7s ease-in-out;
a {
  cursor: pointer;
  height: 40px;
  width: 40px;
img {
  height: 100%;
  width: 100%;
}
}
.device-menu-circle {
  position: absolute;
  height: 100%;
  width: 100%;
  left: 0;
  top: 0;
  opacity: 0;
  z-index: -10;
a {
  -moz-transition: all 0.1s ease-in;
  -webkit-transition: all 0.1s ease-in;
  display: inherit;
  height: inherit !important;
  position: absolute;
  transition: all 0.1s ease-in;
  width: inherit !important;
&:nth-of-type(1) {
   bottom: 0;
 }
&:nth-of-type(2) {
   right: 0;
   top: 0;
 }
&:nth-of-type(3) {
   right: 0;
 }
&:nth-of-type(4) {
   bottom: 0;
   right: 0;
 }
&:nth-of-type(5) {
   bottom: 0;
 }
&:nth-of-type(6) {
   bottom: 0;
   left: 0;
 }
&:nth-of-type(7) {
   left: 0;
 }
&:nth-of-type(8) {
   left: 0;
   top: 0;
 }
}
}
.device-menu-circle.show {
  opacity: 1;
a {
&:nth-child(+n+8) {
   display: none;
 }
&:nth-of-type(1) {
   bottom: 180%;
 }
&:nth-of-type(2) {
   right: -130%;
   top: -130%;
 }
&:nth-of-type(3) {
   right: -180%;
 }
&:nth-of-type(4) {
   bottom: -130%;
   right: -130%;
 }
&:nth-of-type(5) {
   bottom: -180%;
 }
&:nth-of-type(6) {
   bottom: -130%;
   left: -130%;
 }
&:nth-of-type(7) {
   left: -180%;
 }
&:nth-of-type(8) {
   left: -130%;
   top: -130%;
 }
&:hover {
img {
  transform: scale(1.1);
}
}
}
}
a.device-menu-button {
&:hover {
img {
  transform: scale(1.1);
}
}
}
}
.device-menu.as-button {
img.device {
  cursor: pointer;
&:hover {
   transform: scale(1.1);
 }
}
}
.unselectable {
  -khtml-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  -webkit-user-drag: none;
  -webkit-user-select: none;
  user-drag: none;
  user-select: none;
}
</style>
