<script setup lang="ts">
import {onMounted, onUnmounted, PropType, ref, watch} from "vue";
import {ButtonAction, CardItem, Compare, RenderVar, Resolve} from "@/views/Dashboard/core";
import {ElMessage} from "element-plus";
import {GetFullImageUrl, GetFullUrl} from "@/utils/serverId";
import {debounce} from "lodash-es";
import api from "@/api/api";
import {useI18n} from "@/hooks/web/useI18n";
import {AttributeValue, GetAttributeValue} from "@/components/Attributes";
import Image from './Image.vue'
import Icon from "@/views/Dashboard/card_items/icon/src/Icon.vue"

const {t} = useI18n()

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
onMounted(() => {
  if (props.item?.payload.image?.image) {
    image.value = GetFullImageUrl(props.item.payload.image.image)
  } else {
    image.value = null
  }
  background.value = props.item?.payload.image?.background || false
})

onUnmounted(() => {
  image.value = null
  background.value = false
})

// ---------------------------------
// component methods
// ---------------------------------
const image = ref<Nullable<string>>(null)
const background = ref<Nullable<boolean>>(false)

const update = debounce(async () => {
  if (props.item?.payload.image?.image) {
    image.value = GetFullImageUrl(props.item.payload.image.image)
  } else {
    image.value = null
  }
  background.value = props.item?.payload.image?.background || false
  if (props.item?.payload.image.attrField) {
    let token: string = props.item?.payload.image?.attrField || ''
    const url = await RenderVar(token, props.item?.lastEvent)
    image.value = GetFullUrl(url)
  }

  if (props.item?.payload.image?.items) {
    for (const prop of props.item?.payload.image?.items) {
      let val = Resolve(prop.key, props.item?.lastEvent);
      if (!val) {
        continue;
      }

      if (typeof val === 'object') {
        if (val && val.hasOwnProperty('type') && val.hasOwnProperty('name')) {
          val = GetAttributeValue(val as AttributeValue);
        }
      }

      if (val == undefined) {
        val = '[NO VALUE]';
      }

      const tr = Compare(val, prop.value, prop.comparison);
      if (tr) {
        if (prop.image) {
          image.value = GetFullImageUrl(prop.image);
        }
        if (prop.attrField) {
          let token: string = prop.attrField || ''
          image.value = GetFullUrl(await RenderVar(token, props.item?.lastEvent))
        }
        background.value = prop.background
      }
    }
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

// ---------------------------------
// button options
// ---------------------------------
const showMenu = ref(false)

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
    name: action.action || '',
    attributes: {},
    tags: action.tags || [],
    areaId: action.areaId,
  });
  ElMessage({
    title: t('Success'),
    message: t('message.callSuccessful'),
    type: 'success',
    duration: 2000
  });
}

const callBaseAction = async () => {
  if (props.item.buttonActions.length === 0 || props.item.buttonActions.length > 1) {
    return
  }
  callAction(props.item.buttonActions[0])
}

// ---------------------------------
// run
// ---------------------------------

update();

</script>

<template>
  <div ref="el" :class="[{'hidden': item.hidden}]">
    <div
      style="width: 100%; height: 100%; cursor: pointer"
      v-if="item.asButton"
      v-show="!item.hidden"
      @mouseover="mouseOver"
      @mouseleave="mouseLive"
      class="device-menu"
      :class="[{'as-button': item.asButton && item.buttonActions.length > 0}]"
    >
      <Image class="device"
             :image="image"
             :background="background"
             @click.prevent.stop="callBaseAction()"/>

      <div
        :class="[{'show': showMenu}]"
        class="device-menu-circle"
        v-if="item.asButton && item.buttonActions.length > 1"
      >
        <a
          href="#"
          class="device-menu-circle-item"
          v-for="(action, index) in item.buttonActions"
          @click.prevent.stop="callAction(action)"
          :key="index">

          <img
            v-if="action.image"
            :src="GetFullImageUrl(action.image)"/>
          <Icon
            v-else-if="action.icon"
            :icon="action.icon"
            :icon-color="action.iconColor"
          />

        </a>
      </div>

    </div>
    <div v-else v-show="!item.hidden"
         style="width: 100%; height: 100%">
      <Image class="device" :image="image" :background="background"/>
    </div>

  </div>

</template>

<style lang="less" scoped>


.hidden {
  z-index: -99999;
}

.device {
  width: 100%;
  height: auto;
}
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
        img, .el-icon {
          transform: scale(1.1);
        }
      }
    }
  }

  &.as-button {
    .device {
      cursor: pointer;

      &:hover {
        transform: scale(1.1);
      }
    }
  }
}
</style>
