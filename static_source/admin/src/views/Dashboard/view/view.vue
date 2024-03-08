<script setup lang="ts">
import {computed, onMounted, onUnmounted, reactive, ref} from 'vue'
import {ElTabPane, ElTabs} from 'element-plus'
import api from "@/api/api";
import {UUID} from "uuid-generator-ts";
import stream from "@/api/stream";
import {Core, eventBus, stateService} from "@/views/Dashboard/core";
import ViewTab from "@/views/Dashboard/view/ViewTab.vue";
import {propTypes} from "@/utils/propTypes";
import {EventStateChange} from "@/api/types";
import {useAppStore} from "@/store/modules/app";
import {GetFullImageUrl} from "@/utils/serverId";

const appStore = useAppStore()

// ---------------------------------
// common
// ---------------------------------

const loading = ref(false)
const core = reactive<Core>(new Core());
const currentID = ref('')

const props = defineProps({
  id: propTypes.number.def(0),
})

const eventStateChanged = (eventName: string, event: EventStateChange) => {
  core.onStateChanged(event)
}

const eventBusHandler = (eventName: string, event: EventStateChange) => {
  core.eventBusHandler(eventName, event)
}

onMounted(() => {
  const uuid = new UUID()
  currentID.value = uuid.getDashFreeUUID()

  fetchDashboard()

  stream.subscribe('state_changed', currentID.value, stateService.onStateChanged);
  eventBus.subscribe('stateChanged', eventStateChanged)

  eventBus.subscribe(undefined, eventBusHandler)
})

onUnmounted(() => {

  stream.unsubscribe('state_changed', currentID.value);
  eventBus.unsubscribe('stateChanged', eventStateChanged)

  eventBus.unsubscribe(undefined, eventBusHandler)
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
    eventBus.emit('updateGrid', core.getActiveTab?.id)
  }
})

const getTabStyle = () => {
  const style = {}
  if (core.getActiveTab?.background) {
    style['background-color'] = core.getActiveTab?.background
  } else {
    if (core.getActiveTab?.backgroundAdaptive) {
      style['background-color'] = appStore.isDark ? '#333335' : '#FFF'
    }
  }

  if (core.getActiveTab?.backgroundImage) {
    style['background-image'] = `url(${GetFullImageUrl(core.getActiveTab.backgroundImage)})`
    style['background-repeat'] = 'repeat';
    style['background-position'] = 'center';
    // style['background-size'] = 'cover';
  }
  return style
}

</script>

<template>
  <ElTabs
      v-model="activeTabIdx"
      v-if="core.tabs.length > 1 && !loading"
      :style="getTabStyle()"
      class="pl-20px"
      :lazy="true">
    <ElTabPane
        v-for="(tab, index) in core.tabs"
        :label="tab.name"
        :key="index"
        :class="[{'gap': tab.gap}]"
        :disabled="!tab.enabled"
        :lazy="true">
      <ViewTab :tab="tab" :key="index" :core="core"/>
    </ElTabPane>
  </ElTabs>

  <div v-if="core.tabs.length && core.tabs.length === 1 && !loading" :class="[{'gap': core.tabs[0].gap}]"
       :style="getTabStyle()" class="pl-20px pt-20px !min-h-[100%] ">
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

</style>
