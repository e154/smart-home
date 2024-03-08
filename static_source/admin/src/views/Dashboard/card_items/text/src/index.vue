<script setup lang="ts">
import {onMounted, onUnmounted, PropType, ref, watch} from "vue";
import {
  ButtonAction,
  Cache,
  CardItem,
  Compare, eventBus,
  GetTokens,
  RenderText,
  requestCurrentState,
  Resolve
} from "@/views/Dashboard/core";
import api from "@/api/api";
import {ElMessage} from "element-plus";
import debounce from 'lodash.debounce'
import {useI18n} from "@/hooks/web/useI18n";
import {AttributeValue, GetAttributeValue} from "@/components/Attributes";

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

const el = ref<ElRef>(null)
onMounted(() => {

})

onUnmounted(() => {

})

// ---------------------------------
// component methods
// ---------------------------------

const reload = () => {
  reloadKey.value += 1
}

const update = debounce(async () => {
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
        val = GetAttributeValue(val as AttributeValue)
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

    value2 = await RenderText(prop.tokens, value2, props.item?.lastEvent)

    currentValue.value = value2 || value
    return
  }

  const tokens = GetTokens(value, _cache)
  if (tokens) {
    value = await RenderText(tokens, value, props.item?.lastEvent)
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
  if (action?.eventName) {
    eventBus.emit(action?.eventName, action?.eventArgs)
  }
  await api.v1.interactServiceEntityCallAction({
    id: action.entityId,
    name: action.action || '',
    tags: action.tags || [],
    areaId: action.areaId,
    attributes: {},
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

const getStyle = () => {
  return props.item?.style || {}
}

requestCurrentState(props.item.entityId!);

update()

</script>

<template>
  <div ref="el" :class="[{'hidden': item.hidden}]" class="w-[100%] h-[100%]">
    <div
      v-if="item.asButton"
      v-show="!item.hidden"
      @mouseover="mouseOver"
      @mouseleave="mouseLive()"
      class="device-menu w-[100%] h-[100%]"
      :class="[{'as-button': item.asButton && item.buttonActions.length > 0}]"
    >

      <div
        class="cursor-pointer w-[100%] h-[100%]"
        :style="item.style"
        v-html="currentValue"
        :key="reloadKey"
        @click.prevent.stop="callBaseAction()"></div>

    </div>
    <div v-else v-show="!item.hidden" class="w-[100%] h-[100%]">
      <div
        class="w-[100%] h-[100%]"
        :style="getStyle()"
        v-show="!item.hidden"
        v-html="currentValue"
        :key="reloadKey">

      </div>
    </div>
  </div>

</template>

<style lang="less" scoped>
.hidden {
  z-index: -99999;
}

.ql-align-center {
  text-align: center;
}

.cursor-pointer {
  cursor: pointer;
}

.device-menu.as-button {
  img.device {
    cursor: pointer;

    &:hover {
      transform: scale(1.1);
    }
  }
}

:deep(svg) {
  display: inline !important;
  vertical-align: middle;
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
