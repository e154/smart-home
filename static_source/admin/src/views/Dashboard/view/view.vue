<script setup lang="ts">
import {computed, onMounted, onUnmounted, PropType, reactive, ref, shallowReactive} from 'vue'
import {ElTabs, ElTabPane} from 'element-plus'
import api from "@/api/api";
import {EventStateChange} from "@/api/stream_types";
import {UUID} from "uuid-generator-ts";
import stream from "@/api/stream";
import {Core} from "@/views/Dashboard/core";
import 'splitpanes/dist/splitpanes.css'
import {useBus} from "@/views/Dashboard/bus";
import ViewTab from "@/views/Dashboard/view/ViewTab.vue";
import {propTypes} from "@/utils/propTypes";

const {bus} = useBus()

// ---------------------------------
// common
// ---------------------------------

const loading = ref(false)
const core = reactive<Core>(new Core());
const currentID = ref('')

const props = defineProps({
  id: propTypes.number.def(0),
})

const onStateChanged = (event: EventStateChange) => {
  bus.emit('state_changed', event);
  core.onStateChanged(event);
}

onMounted(() => {
  const uuid = new UUID()
  currentID.value = uuid.getDashFreeUUID()

  // setTimeout(() => {
  stream.subscribe('state_changed', currentID.value, onStateChanged);
  // }, 1000)
})

onUnmounted(() => {
  stream.unsubscribe('state_changed', currentID.value);
})

// ---------------------------------
// dashboard
// ---------------------------------

const fetchDashboard = async () => {
  loading.value = true;
  const res = await api.v1.dashboardServiceGetDashboardById(props.id)
      .catch(() => {
      })
      .finally(() => {
        loading.value = false;
      })
  core.currentBoard(res.data);
}

const activeTabIdx = computed({
  get(): string {
    return core.activeTabIdx + ''
  },
  set(value: string) {
    core.activeTabIdx = parseInt(value)
  }
})

const getBackgroundColor = () => {
  return {backgroundColor: core.getActiveTab?.background}
}

fetchDashboard()

</script>

<template>
  <ElTabs v-model="activeTabIdx"  v-if="core.tabs.length > 1 && !loading" :style="getBackgroundColor()" class="pl-20px !min-h-[100%]">
    <ElTabPane
        v-for="(tab, index) in core.tabs"
        :label="tab.name"
        :key="index"
        :class="[{'gap': tab.gap}]"
        :lazy="true">
      <ViewTab :tab="tab" :key="index" :core="core"/>
    </ElTabPane>
  </ElTabs>

  <div v-if="core.tabs.length && core.tabs.length === 1 && !loading" :class="[{'gap': core.tabs[0].gap}]" :style="getBackgroundColor()" class="pl-20px pt-20px !min-h-[100%] ">
    <ViewTab :tab="core.tabs[0]" :core="core"/>
  </div>
</template>

<style lang="less">


/* Track */
::-webkit-scrollbar-track {
  background: #f1f1f1;
}

p {
  display: block;
  margin-block-start: 1em;
  margin-block-end: 1em;
  margin-inline-start: 0;
  margin-inline-end: 0;
}

h1 {
  display: block;
  font-size: 2em;
  margin-block-start: 0.67em;
  margin-block-end: 0.67em;
  margin-inline-start: 0;
  margin-inline-end: 0;
  font-weight: 700;
}

h2 {
  display: block;
  font-size: 1.5em;
  margin-block-start: 0.67em;
  margin-block-end: 0.67em;
  margin-inline-start: 0;
  margin-inline-end: 0;
  font-weight: 700;
}

html {
  line-height: 1.15;
}

.splitpanes.default-theme .splitpanes__splitter {
  background-color: #bfbfbf6e;
}
</style>
