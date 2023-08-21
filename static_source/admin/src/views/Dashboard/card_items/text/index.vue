<script setup lang="ts">
import {onMounted, onUnmounted, PropType, ref, watch} from "vue";
import {ButtonAction, CardItem, requestCurrentState} from "@/views/Dashboard/core";
import {Cache, Compare, GetTokens, RenderText, Resolve} from "@/views/Dashboard/render";
import api from "@/api/api";
import {ElMessage} from "element-plus";
import {Attribute, GetAttrValue} from "@/api/stream_types";
import debounce from 'lodash.debounce'
import {useI18n} from "@/hooks/web/useI18n";

const {t} = useI18n()

// ---------------------------------
// common
// ---------------------------------
const reloadKey = ref(0)
const _cache = new Cache()
const currentValue = ref('')

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
})

onUnmounted(() => {

})

// ---------------------------------
// component methods
// ---------------------------------

const reload = () => {
  reloadKey.value += 1
}

const update = debounce(() => {
  // console.log('update value', item.value?.payload?.text);

  if (!props.item?.payload.text?.items) {
    currentValue.value = props.item?.payload.text?.default_text || ''
    return
  }

  let value = props.item?.payload.text?.default_text || ''
  let value2 = ''

  for (const prop of props.item?.payload.text?.items) {
    // select prop
    let val = Resolve(prop.key, props.item?.lastEvent)
    if (!val) {
      continue
    }

    if (typeof val === 'object') {
      if (val && val.hasOwnProperty('type') && val.hasOwnProperty('name')) {
        val = GetAttrValue(val as Attribute)
      }
    }

    if (val == undefined) {
      val = '[NO VALUE]'
    }

    const tr = Compare(val, prop.value, prop.comparison)
    if (!tr) {
      continue
    }

    if (!prop.tokens) {
      prop.tokens = []
    }

    // render text
    prop.tokens = GetTokens(prop.text, _cache)
    if (!prop.tokens.length) {
      currentValue.value = prop.text || ''
      return
    }

    if (prop.text) {
      value2 = prop.text
    }

    value2 = RenderText(prop.tokens, value2, props.item?.lastEvent)

    currentValue.value = value2 || value
    return
  }

  const tokens = GetTokens(value, _cache)
  if (tokens) {
    value = RenderText(tokens, value, props.item?.lastEvent)
  }
  currentValue.value = value
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

const resetCache = () => {
  _cache.clear()
}

const showMenu = ref(false)
const timer = ref<any>(null)

const mouseLive = () => {
  if (!showMenu.value) {
    return;
  }
  timer.value = setTimeout(() => {
    showMenu.value = false;
    timer.value = null;
  }, 2000);
}

const mouseOver = () => {
  showMenu.value = true;
  if (!timer.value) {
    return;
  }
  clearTimeout(timer.value);
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

const getStyle = () => {
  return props.item?.style || {}
}

requestCurrentState(props.item.entityId!);

update()

</script>

<template>
  <div
      v-if="item.asButton"
      ref="el"
      v-show="!item.hidden"
      @mouseover="mouseOver"
      @mouseleave="mouseLive()"
      class="device-menu"
      :class="[{'as-button': item.asButton && item.buttonActions.length > 0}]"
  >

    <div
        class="cursor-pointer"
        :style="item.style"
        v-html="currentValue"
        :key="reloadKey"
        @click.prevent.stop="callAction(item.buttonActions[0])"></div>

    <div
        :class="[{'show': showMenu}]"
        class="device-menu-circle"
        v-if="item.asButton && item.buttonActions.length > 1"
    >
      <a
          href="#"
          v-for="(action, index) in item.buttonActions"
          :key="index"
          @click.prevent.stop="callAction(action)">
        <img :src="item.getUrl(action.image)"/>
      </a>
    </div>
  </div>
  <div v-else
       ref="el"
       :style="getStyle()"
       v-show="!item.hidden"
       v-html="currentValue"
       :key="reloadKey"></div>
</template>

<style lang="less" >
.ql-align-center {
  text-align: center;
}
.cursor-pointer {
  cursor: pointer;
}
</style>
