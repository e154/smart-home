<script setup lang="ts">
import {computed, onMounted, onUnmounted, reactive, ref, shallowReactive} from 'vue'
import {useI18n} from '@/hooks/web/useI18n'
import {ElMessage, ElTabs, ElTabPane, ElSkeleton} from 'element-plus'
import {useForm} from '@/hooks/web/useForm'
import {useCache} from '@/hooks/web/useCache'
import {useRoute, useRouter} from 'vue-router'
import {useValidator} from '@/hooks/web/useValidator'
import api from "@/api/api";
import {EventStateChange} from "@/api/stream_types";
import {UUID} from "uuid-generator-ts";
import stream from "@/api/stream";
import {Card, Core, Tab} from "@/views/Dashboard/core";
import { Splitpanes, Pane } from 'splitpanes'
import 'splitpanes/dist/splitpanes.css'
import TabSettings from "@/views/Dashboard/editor/TabSettings.vue";
import {useBus} from "@/views/Dashboard/bus";
import TabEditor from "@/views/Dashboard/editor/TabEditor.vue";
import TabCard from "@/views/Dashboard/editor/TabCard.vue";
import ViewTab from "@/views/Dashboard/editor/ViewTab.vue";
import TabCardItem from "@/views/Dashboard/editor/TabCardItem.vue";

const {bus} = useBus()
const {register, elFormRef, methods} = useForm()
const {required} = useValidator()
const route = useRoute();
const {currentRoute, addRoute, push} = useRouter()
const {wsCache} = useCache()
const {t} = useI18n()

// ---------------------------------
// common
// ---------------------------------

const loading = ref(false)
const dashboardId = computed(() => route.params.id as number);
const core = reactive<Core>(new Core());
const currentID = ref('')

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
  const res = await api.v1.dashboardServiceGetDashboardById(dashboardId.value)
      .catch(() => {
      })
      .finally(() => {
        loading.value = false;
      })
  core.currentBoard(res.data);
}

fetchDashboard()

useBus({
  name: 'fetchDashboard',
  callback: () => {
    fetchDashboard()
  }
})
// ---------------------------------
// tabs
// ---------------------------------

const handleTabsEdit = (targetName: string, action: string) => {
  switch (action) {
    case 'add':
      createTab();
      break;
    case 'remove':
  }
}

const updateCurrentTab = (tab: any, ev: any) => {
  const {index} = tab;
  if (core.activeTab === index) return;
  core.activeTab = index;
  core.updateCurrentTab();
}

const createTab = async () => {
  await core.createTab();
  ElMessage({
    title: t('Success'),
    message: t('message.createdSuccessfully'),
    type: 'success',
    duration: 2000
  });
}

const activeTabIdx = computed({
  get(): string {
    return core.activeTab + ''
  },
  set(value: string) {
    core.activeTab = parseInt(value)
  }
})

const activeTab = computed<Tab>(() => core.tabs[core.activeTab] as Tab)
const activeCard = computed<Card>(() => core.tabs[core.activeTab].cards[core.activeCard] as Card)

const getBackgroundColor = () => {
  return {backgroundColor: core.tabs[core.activeTab]?.background}
}

</script>

<template>
  <div class="components-container dashboard-container" style="margin: 0" v-if="!loading" :style="getBackgroundColor()">

  <splitpanes class="default-theme" horizontal>
    <pane min-size="10" max-size="90" :size="50" class="top-container">
        <ElTabs
            v-model="activeTabIdx"
            @edit="handleTabsEdit"
            @tab-click="updateCurrentTab"
            class="ml-20px">
          <ElTabPane
              v-for="(tab, index) in core.tabs"
              :label="tab.name"
              :key="index"
              :class="[{'gap': tab.gap}]">
            <ViewTab :tab="tab" :key="index" :core="core"/>
          </ElTabPane>
        </ElTabs>
    </pane>
    <pane class="bottom-container">
      <ElTabs v-model="core.mainTab" >
        <!-- main -->
        <ElTabPane :label="$t('dashboard.mainTab')" name="main">
          <TabSettings v-if="core.current" :core="core"/>
        </ElTabPane>
        <!-- /main -->

        <!-- tabs -->
        <ElTabPane :label="$t('dashboard.tabsTab')" name="tabs">
          <TabEditor v-if="core.current" :tab="activeTab" :core="core"/>
        </ElTabPane>
        <!-- /tabs -->

        <!-- cards -->
        <ElTabPane :label="$t('dashboard.cardsTab')" name="cards">
          <TabCard v-if="core.current && activeTab" :tab="activeTab" :core="core"/>
        </ElTabPane>
        <!-- /cards -->

        <!-- cardItems -->
        <ElTabPane :label="$t('dashboard.cardItemsTab')" name="cardItems">
          <TabCardItem v-if="core.current && activeTab && activeCard" :card="activeCard" :core="core"/>
        </ElTabPane>
        <!-- /cardItems -->

      </ElTabs>
    </pane>
  </splitpanes>

  </div>

</template>

<style lang="less">

.splitpanes.default-theme .splitpanes__pane {
  background-color: inherit;
}


/* Track */
::-webkit-scrollbar-track {
  background: #f1f1f1;
}

.dashboard-container {
  position: relative;
}

.item-card-editor {
  /*background-color: #f1f1f1;*/
}

.components-container {
  height: calc(100vh - 50px);
  /*min-height: calc(100vh - 50px);*/
  margin: 0;
  padding: 0;
}

.top-container {
  width: 100%;
  height: 100%;
  overflow-y: scroll;
}

.bottom-container {
  width: 100%;
  height: 100%;
  padding: 0 20px;
  overflow-y: scroll;
}

.filter-list {

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
