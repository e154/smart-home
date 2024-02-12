<script setup lang="ts">
import {nextTick, onMounted, onUnmounted, PropType, ref, unref, watch} from "vue";
import {CardItem} from "@/views/Dashboard/core/core";
import {RenderVar, Resolve} from "@/views/Dashboard/core/render";
import {debounce} from "lodash-es";
import Iconify from "@purge-icons/generated";
import {AttributeValue, GetAttributeValue} from "@/components/Attributes";
import {Compare} from "@/views/Dashboard/core/types";

// ---------------------------------
// common
// ---------------------------------
const props = defineProps({
  item: {
    type: Object as PropType<Nullable<CardItem>>,
    default: () => null
  },
})

const elRef = ref<ElRef>(null)
onMounted(() => {
  // store dom element moveable
  props.item.setTarget(elRef.value)
})

onUnmounted(() => {

})

// ---------------------------------
// component methods
// ---------------------------------
const icon = ref<Nullable<string>>(null)
const iconColor = ref<Nullable<string>>(null)

const update = debounce( async () => {
  let _icon = props.item?.payload.icon?.value || '';
  if (props.item?.payload.icon.attrField) {
    let token: string = props.item?.payload.icon?.attrField || ''
    _icon = RenderVar(token, props.item?.lastEvent)
  }

  iconColor.value = props.item?.payload?.icon?.iconColor || '#eee';

  if (props.item?.payload.icon?.items) {
    for (const prop of props.item?.payload.icon?.items) {
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
        if (prop.iconColor) iconColor.value = prop.iconColor;
        if (prop.icon) _icon = prop.icon;
        break
      }
    }
  }

  const el = unref(elRef)
  if (!el) return

  await nextTick()

  const svg = Iconify.renderSVG(_icon, {})
  if (svg) {
    el.textContent = ''
    el.appendChild(svg)
  } else {
    const span = document.createElement('span')
    span.className = 'iconify'
    span.dataset.icon = _icon
    el.textContent = ''
    el.appendChild(span)
  }
  return;
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

update();

</script>

<template>
  <div ref="elRef" :class="[{'hidden': item.hidden}]" class="icon-item"></div>
</template>

<style lang="less" scoped>
.icon-item {
  :deep(svg) {
    height: 100%;
    width: 100%;
    color: v-bind(iconColor)!important;
  }
}
</style>
