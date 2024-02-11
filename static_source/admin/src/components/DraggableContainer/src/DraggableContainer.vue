<script setup lang="ts">

import {onBeforeUnmount, onMounted, ref, watch} from "vue";
import {propTypes} from "@/utils/propTypes";
import {useCache} from "@/hooks/web/useCache";
import {debounce} from "lodash-es";
import {useAppStore} from "@/store/modules/app";

const emit = defineEmits(['resize'])

const appStore = useAppStore()
const {wsCache} = useCache()

const props = defineProps({
  name: propTypes.string.def('main'),
  initialTop: propTypes.number.def(0),
  initialLeft: propTypes.number.def(0),
  initialWidth: propTypes.number.def(350),
  initialHeight: propTypes.number.def(350),
  maxWidth: propTypes.number.def(Infinity),
  maxHeight: propTypes.number.def(800),
  minWidth: propTypes.number.def(350),
  parentElement: {type: HTMLElement, default: null}
})

const menu = ref(null);
const top = ref(props.initialTop);
const left = ref(props.initialLeft);
const width = ref(props.initialWidth);
const height = ref(props.initialHeight);
const visible = ref(true);
const isDragging = ref(false);
let offsetX = ref(0);
let offsetY = ref(0);
let startWidth = ref(props.initialWidth);
let startHeight = ref(props.initialHeight);
const zIndex = ref(10);

let moveDirection: string;

onMounted(() => {
  restoreState();

  if (props.parentElement) {
    props.parentElement.appendChild(menu.value);
  } else {
    document.body.appendChild(menu.value);
  }
});

onBeforeUnmount(() => {
  if (props.parentElement) {
    props.parentElement.removeChild(menu.value);
  } else {
    document.body.removeChild(menu.value);
  }
});

watch([top, left, width, height], () => {
  saveState();
});

const startDragging = (dir: string, event: MouseEvent) => {
  moveDirection = dir
  isDragging.value = true;
  offsetX.value = event.clientX;
  offsetY.value = event.clientY;
  startWidth.value = width.value;
  startHeight.value = height.value;
  window.addEventListener('mousemove', draggingHandler);
  window.addEventListener('mouseup', stopDragging);
}

const draggingHandler = (e: MouseEvent) => {
  if (moveDirection === 'move') {
    dragMenu(e);
  } else if (moveDirection === 'right') {
    resizeRight(e);
  } else if (moveDirection === 'bottom') {
    resizeBottom(e);
  } else if (moveDirection === 'corner') {
    resizeCorner(e);
  }
}

const dragMenu = (event: MouseEvent) => {
  if (isDragging.value) {
    const deltaX = event.clientX - offsetX.value;
    const deltaY = event.clientY - offsetY.value;
    left.value += deltaX;
    top.value += deltaY;
    offsetX.value = event.clientX;
    offsetY.value = event.clientY;
    if (top.value < 0) top.value = 0
    if (left.value < 0) left.value = 0
  }
}

const resizeRight = (event: MouseEvent) => {
  if (isDragging.value) {
    const deltaX = event.clientX - offsetX.value;
    width.value = startWidth.value + deltaX;
    if (width.value < props.minWidth) width.value = props.minWidth
    onResize();
  }
}

const resizeBottom = (event: MouseEvent) => {
  if (isDragging.value) {
    const deltaY = event.clientY - offsetY.value;
    height.value = startHeight.value + deltaY;
    onResize();
  }
}

const resizeCorner = (event: MouseEvent) => {
  if (isDragging.value) {
    const deltaX = event.clientX - offsetX.value;
    const deltaY = event.clientY - offsetY.value;
    width.value = startWidth.value + deltaX;
    if (width.value < props.minWidth) width.value = props.minWidth
    height.value = startHeight.value + deltaY;
    onResize();
  }
}

const onResize = debounce(() => {
  emit('resize')
}, 100)

const stopDragging = () => {
  isDragging.value = false;
  window.removeEventListener('mousemove', draggingHandler);
  window.removeEventListener('mouseup', stopDragging);
}

const restoreState = () => {
  const position = wsCache.get(`${props.name}-position`);
  if (position) {
    top.value = position.top;
    left.value = position.left;
  }

  const size = wsCache.get(`${props.name}-size`);
  if (size) {
    width.value = size.width;
    height.value = size.height;
  }

  const _visible = wsCache.get(`${props.name}-visibility`);
  if (_visible != undefined) {
    visible.value = _visible
  }
}

const saveState = debounce(() => {
  wsCache.set(`${props.name}-position`, {top: top.value, left: left.value});
  wsCache.set(`${props.name}-size`, {width: width.value, height: height.value});
}, 100)

const toggleVisibility = () => {
  visible.value = !visible.value;
  wsCache.set(`${props.name}-visibility`, visible.value);
}

const bringToFront = () => {
  zIndex.value = appStore.getMaxZIndex(); // Устанавливаем z-index на 1 больше максимального
}

</script>

<template>
  <div
      class="draggable-container"
      :class="'container-' + name"
      :style="{ top: `${top}px`, left: `${left}px`, width: `${width}px`, height: `${visible?height:22}px`, zIndex: zIndex }"
      ref="menu"
      @mousedown="bringToFront"
  >
    <div class="draggable-container-header"
         @mousedown="startDragging('move', $event)"
         @dblclick="toggleVisibility"
    >
      <slot name="header"></slot>
    </div>
    <div class="resizer right" @mousedown="startDragging('right', $event)"></div>
    <div class="resizer bottom" @mousedown="startDragging('bottom', $event)"></div>
    <div class="resizer corner" @mousedown="startDragging('corner', $event)"></div>
    <div v-show="visible" class="draggable-container-content">
      <slot></slot>
    </div>
    <div v-show="visible" class="draggable-container-footer">
      <slot name="footer"></slot>
    </div>
  </div>
</template>

<style lang="less" scoped>
.draggable-container {
  position: absolute;
  width: 229px;
  z-index: 1000;
  background-color: var(--el-bg-color);
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.2);
  display: flex;
  flex-direction: column;
}

.draggable-container-content {
  position: relative;
//background-color: var(--el-bg-color); padding: 0 10px 10px 10px;
  padding: 0 10px 10px 10px;
  flex-grow: 1; /* Занимаем все оставшееся пространство */
  overflow: auto;
}

.draggable-container-header {
  color: var(--left-menu-text-active-color) !important;
  background-color: var(--left-menu-bg-color);
  font-size: 12px;
  padding: 5px;
  cursor: move; /* Устанавливаем курсор перемещения */
  user-select: none;
}

.draggable-container-footer {
  margin-top: auto; /* Footer всегда будет прижат к низу */
}

.resizer {
  position: absolute;
  user-select: none;
}

.right {
  width: 5px;
  height: 100%;
  top: 0;
  right: -2.5px;
  cursor: ew-resize;
}

.bottom {
  width: 100%;
  height: 5px;
  left: 0;
  bottom: -2.5px;
  cursor: ns-resize;
}

.corner {
  width: 10px;
  height: 10px;
  right: -5px;
  bottom: -5px;
  cursor: nwse-resize;
}
</style>
