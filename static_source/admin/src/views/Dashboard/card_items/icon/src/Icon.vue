<script setup lang="ts">
import {nextTick, ref, unref, watch} from "vue";
import {debounce} from "lodash-es";
import Iconify from "@purge-icons/generated";
import {propTypes} from "@/utils/propTypes";

// ---------------------------------
// common
// ---------------------------------
const props = defineProps({
  icon: propTypes.string.def(''),
  iconColor: propTypes.string.def(''),
})

const elRef = ref<ElRef>(null)

// ---------------------------------
// component methods
// ---------------------------------
const update = debounce(async () => {

  const el = unref(elRef)
  if (!el) return

  await nextTick()

  const svg = Iconify.renderSVG(props.icon, {})
  if (svg) {
    el.textContent = ''
    el.appendChild(svg)
  } else {
    const span = document.createElement('span')
    span.className = 'iconify'
    span.dataset.icon = props.icon
    el.textContent = ''
    el.appendChild(span)
  }
  return;
}, 100)

watch(
  () => [props.icon, props.iconColor],
  (val) => {
    if (!val) return;
    update()
  },
  {
    deep: true,
    immediate: true
  }
)

</script>

<template>
  <div ref="elRef" class="icon-item"></div>
</template>

<style lang="less" scoped>
.icon-item {
  :deep(svg) {
    height: 100%;
    width: 100%;
    color: v-bind(iconColor) !important;
  }
}
</style>
