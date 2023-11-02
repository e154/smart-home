<script setup lang="ts">
import {computed, onMounted, onUnmounted, PropType, reactive, ref, shallowReactive} from 'vue'
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
import 'splitpanes/dist/splitpanes.css'
import {useBus} from "@/views/Dashboard/bus";
import ViewTab from "@/views/Dashboard/view/ViewTab.vue";
import {propTypes} from "@/utils/propTypes";
import {ApiVariable} from "@/api/stub";

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
    return core.activeTab + ''
  },
  set(value: string) {
    core.activeTab = parseInt(value)
  }
})

const getBackgroundColor = () => {
  return {backgroundColor: core.tabs[core.activeTab]?.background}
}

fetchDashboard()

</script>

<template>
  <ElTabs v-model="activeTabIdx"  v-if="core.tabs.length > 1 && !loading" :style="getBackgroundColor()" class="pl-20px pt-20px !min-h-[calc(100%-var(--top-tool-height))]">
    <ElTabPane
        v-for="(tab, index) in core.tabs"
        :label="tab.name"
        :key="index"
        :class="[{'gap': tab.gap}]"
        :lazy="true">
      <ViewTab :tab="tab" :key="index" :core="core"/>
    </ElTabPane>
  </ElTabs>

  <div v-if="core.tabs.length && core.tabs.length === 1 && !loading" :class="[{'gap': core.tabs[0].gap}]" :style="getBackgroundColor()" class="pl-20px pt-20px !min-h-[calc(100%-var(--top-tool-height))] ">
    <ViewTab :tab="core.tabs[0]" :core="core"/>
  </div>
</template>

<style lang="less">


/* Track */
::-webkit-scrollbar-track {
  background: #f1f1f1;
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
  padding: 0 20px;
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
