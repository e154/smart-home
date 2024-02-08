<script setup lang="ts">
import {useBus} from "@/views/Dashboard/bus";
import {computed, nextTick, onMounted, onUnmounted, PropType, ref} from "vue";
import {Card} from "@/views/Dashboard/core";
import GridStore from './GridStore';
import Muuri from "muuri";
import {UUID} from "uuid-generator-ts";
import {ItemDragHandle, ItemKey, ItemSize} from './constants';
import debounce from 'lodash.debounce'

const uuid = new UUID()
const muuri = ref<Muuri>({} as Muuri)

const props = defineProps({
  modelValue: {
    type: Array as PropType<Card>,
    default: () => []
  },
  /**
   * Optional class name to add to the grid. If not, one will be provided
   */
  className: {
    type: String,
    required: false,
    default: () => `class-${new UUID().getDashFreeUUID().replace(/-/g, '')}`
  },
  /**
   * To set up options of the grid
   * https://github.com/haltu/muuri#grid-options
   */
  options: {
    type: Object,
    required: false,
    default: () => ({})
  },
  /**
   * Array input for items to display (via v-model)
   */
  value: {
    type: Array,
    required: false
  },
  /**
   * Identifier property for each item
   */
  itemKey: {
    type: String,
    required: false,
    default: () => ItemKey.key
  },
  /**
   * Callback to fetch a dynamic width based on item
   */
  getItemWidth: {
    type: Function,
    required: false,
    default: () => ItemSize.width
  },
  /**
   * Callback to fetch a dynamic height based on item
   */
  getItemHeight: {
    type: Function,
    required: false,
    default: () => ItemSize.height
  },
  /**
   * Enable drag and drop feature on the grid
   */
  dragEnabled: {
    type: Boolean,
    required: false,
    default: false
  },
  /**
   * Selector for determining the handle to activate dragging
   */
  dragHandle: {
    type: String,
    required: false,
    default: ItemDragHandle.selector
  },
  /**
   * When dragEnabled is on, can control which other grid you can drag into
   */
  groupId: {
    type: [String, Number],
    required: false
  },
  /**
   * When dragEnabled is on, can control which other grid you can drag into
   */
  groupIds: {
    type: [Array],
    required: false
  }
})
const currentID = ref('')
const selector = computed(() => `.${props.className}`)
const muuriOptions = ref({})
const observer = ref(null)
const grid = ref(null);
const isInitiated = ref(false);
const genKey = () => {
  uuid.getDashFreeUUID().replace(/-/g, '')
}
const gridKey = genKey()

const emit = defineEmits(['updated'])

const update = () => {
  // console.log('update vuury')
  nextTick(() => {
    muuri.value
        .refreshItems()
        .layout(true, () => {})
        .layout(true, () => emit('updated'));
  });
}

const _getItemStyles = (item) => {
  return {
    width: props.getItemWidth(item),
    height: props.getItemHeight(item)
  };
}

const _sync = () => {

}

const _setup = () => {
  muuri.value = new Muuri(selector.value);
  if (props.groupId) {
    GridStore.addGrid(props.groupId, muuri.value);
  }
  if (props.groupIds) {
    GridStore.addGridToGroups(props.groupIds, muuri.value);
  }
  observer.value = new ResizeObserver(() => {
    _resizeOnLoad();
  });
  observer.value.observe(grid.value);
  _sync()
  nextTick(() => {
    GridStore.setItemsForGridId(gridKey, props.value);
  });
}

onMounted(() => {
  const uuid = new UUID()
  currentID.value = uuid.getDashFreeUUID()

  _setup()
})

onUnmounted(() => {
  // console.log('unmount')
})

const _resizeOnLoad = debounce(() => {
  nextTick(() => {
    update();
  });
}, 100)

useBus({
  name: 'updateVuuri',
  callback: () => {
    update();
  }
})

</script>

<template>
  <div
      ref="grid"
      class="muuri-grid"
      :class="className"
      :data-grid-key="gridKey">
    <div
        v-for="item in modelValue"
        :key="item[itemKey]"
        :style="_getItemStyles(item)"
        class="muuri-item"
        :data-item-key="item[itemKey]"
    >
      <div class="muuri-item-content">
        <slot name="item" :item="item"></slot>
      </div>
    </div>

  </div>
</template>


<style scoped>
.muuri-grid {
  position: relative;
  height: 100%;
  min-height: 300px;
  width: 100%;
}

.muuri-item {
  display: block;
  position: absolute;
  z-index: 1;
  width: 100px;
  height: 100px;
}

.muuri-item.muuri-item-dragging {
  z-index: 3;
}

.muuri-item.muuri-item-releasing {
  z-index: 2;
}

.muuri-item.muuri-item-hidden {
  z-index: 0;
}

.muuri-item-content {
  position: relative;
  width: 100%;
  height: 100%;
}

.muuri-item-placeholder {
  opacity: 0.5;
  margin: 0 !important;
}
</style>
