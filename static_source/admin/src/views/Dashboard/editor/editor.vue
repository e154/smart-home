<script setup lang="ts">
import {computed, onMounted, onUnmounted, reactive, ref} from 'vue'
import {useI18n} from '@/hooks/web/useI18n'
import {ElButton, ElEmpty, ElMessage, ElTabPane, ElTabs} from 'element-plus'
import {useRoute} from 'vue-router'
import api from "@/api/api";
import {EventStateChange} from "@/api/stream_types";
import {UUID} from "uuid-generator-ts";
import stream from "@/api/stream";
import {Card, Core, Tab} from "@/views/Dashboard/core/core";
import {useBus} from "@/views/Dashboard/core/bus";
import ViewTab from "@/views/Dashboard/editor/ViewTab.vue";
import {useCache} from "@/hooks/web/useCache";
import {DraggableContainer} from "@/components/DraggableContainer";
import TabSettings from "@/views/Dashboard/editor/TabSettings.vue";
import TabEditor from "@/views/Dashboard/editor/TabEditor.vue";
import TabCardItem from "@/views/Dashboard/editor/TabCardItem.vue";
import TabCard from "@/views/Dashboard/editor/TabCard.vue";

const {emit} = useBus()
const route = useRoute();
const {t} = useI18n()
const {wsCache} = useCache()

// ---------------------------------
// common
// ---------------------------------

const loading = ref(true)
const dashboardId = computed(() => parseInt(route.params.id as string) as number);
const core = reactive<Core>(new Core());
const currentID = ref('')

const onStateChanged = (event: EventStateChange) => {
  emit('state_changed', event);
  core.onStateChanged(event);
}

onMounted(() => {
  const uuid = new UUID()
  currentID.value = uuid.getDashFreeUUID()

  stream.subscribe('state_changed', currentID.value, onStateChanged);

  fetchDashboard()
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

useBus({
  name: 'fetchDashboard',
  callback: () => {
    fetchDashboard()
  }
})
// ---------------------------------
// tabs
// ---------------------------------

const updateCurrentTab = (tab: any, ev: any) => {
  const {index} = tab;
  if (core.activeTabIdx === index) return;
  core.activeTabIdx = index;
  core.updateCurrentTab();
}

const activeTabIdx = computed({
  get(): string {
    return core.activeTabIdx + ''
  },
  set(value: string) {
    core.activeTabIdx = parseInt(value)
  }
})

const activeTab = computed<Tab>(() => core.getActiveTab as Tab)
const activeCard = computed<Card>(() => core.getActiveTab.cards[core.activeCard] as Card)

const getBackgroundColor = () => {
  return {backgroundColor: core.getActiveTab?.background}
}

const tagsView = computed(() => tagsView.value ? 37 : 0)

const createTab = async () => {
  await core.createTab();

  ElMessage({
    title: t('Success'),
    message: t('message.createdSuccessfully'),
    type: 'success',
    duration: 2000
  });
}

const addCard = () => {
  core.createCard();
}


</script>

<template>

  <div class="components-container dashboard-container" style="margin: 0" v-if="!loading" :style="getBackgroundColor()">

    <ElTabs
        v-model="activeTabIdx"
        @tab-click="updateCurrentTab"
        class="ml-20px"
        :lazy="true">
      <ElTabPane
          v-for="(tab, index) in core.tabs"
          :label="tab.name"
          :key="index"
          :disabled="!tab.enabled"
          :class="[{'gap': tab.gap}]">
        <ViewTab :tab="tab" :key="index" :core="core"/>
      </ElTabPane>
    </ElTabs>

    <DraggableContainer :name="'editor-main'">
      <template #header>
        <span>Main Options</span>
      </template>
      <template #default>
        <ElTabs v-model="core.mainTab">
          <!-- main -->
          <ElTabPane :label="$t('dashboard.mainTab')" name="main">
            <template #label>
              <Icon icon="wpf:maintenance"/>
            </template>
            <TabSettings v-if="core.current" :core="core"/>
          </ElTabPane>
          <!-- /main -->

          <!-- tabs -->
          <ElTabPane :label="$t('dashboard.tabsTab')" name="tabs">
            <template #label>
              <Icon icon="vaadin:tabs"/>
            </template>
            <TabEditor v-if="core.current && activeTab" :tab="activeTab" :core="core"/>
            <ElEmpty v-if="!core.tabs.length" :rows="5">
              <ElButton type="primary" @click="createTab()">
                {{ t('dashboard.addNewTab') }}
              </ElButton>
            </ElEmpty>
          </ElTabPane>
          <!-- /tabs -->

          <!-- cards -->
          <ElTabPane :label="$t('dashboard.cardsTab')" name="cards">
            <template #label>
              <Icon icon="material-symbols:cards-outline"/>
            </template>
            <TabCard v-if="core.current && activeTab" :tab="activeTab" :core="core"/>
            <ElEmpty v-if="!core.tabs.length" :rows="5">
              <ElButton type="primary" @click="createTab()">
                {{ t('dashboard.addNewTab') }}
              </ElButton>
            </ElEmpty>
          </ElTabPane>
          <!-- /cards -->

          <!-- cardItems -->
          <ElTabPane :label="$t('dashboard.cardItemsTab')" name="cardItems">
            <template #label>
              <Icon icon="icon-park-solid:add-item"/>
            </template>
            <TabCardItem v-if="core.current && activeTab && activeCard" :card="activeCard" :core="core"/>
            <ElEmpty v-if="!core.tabs.length" :rows="5">
              <ElButton type="primary" @click="createTab()">
                {{ t('dashboard.addNewTab') }}
              </ElButton>
            </ElEmpty>
            <ElEmpty v-if="core.tabs.length && !(core.activeCard >= 0)" :rows="5">
              <ElButton type="primary" @click="addCard()">
                {{ t('dashboard.addNewCard') }}
              </ElButton>
            </ElEmpty>
          </ElTabPane>
          <!-- /cardItems -->

        </ElTabs>
      </template>
    </DraggableContainer>

  </div>

</template>

<style lang="less">


/* Track */
::-webkit-scrollbar-track {
  background: #f1f1f1;
}

.dashboard-container {
  position: relative;
}

.components-container {
  height: calc(100vh - 87px);
//height: inherit; //height: -webkit-fill-available; //height: -moz-available; //height: fill-available; margin: 0; padding: 0;
}

.top-container {
  width: 100%;
  height: 100%;
  overflow-y: scroll;
}

.bottom-container {
  width: 100%;
  padding: 0 20px;
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

.el-tabs {
  height: inherit;
  height: -webkit-fill-available;
}

html.dark {
  .draggable-container.container-editor-main {

    .draggable-container-content,
    .el-divider__text {
      background-color: hsl(230, 7%, 17%);
    }

    .el-card {
      .el-divider__text {
        background-color: var(--el-bg-color-overlay);
      }
    }
  }
}

// custom style
.draggable-container.container-editor-main {
  .el-main {
    padding: 2px !important;
  }

  .el-card__header {
    padding: 18px 20px !important;
  }

  .el-card {
    --el-card-padding: 2px 5px;
  }

  .el-form-item--small {
    margin-bottom: 5px;
  }

  .el-divider--horizontal {
    margin: 11px 0;
  }

  .el-col.el-col-12 {
    padding-right: 6px;
    padding-left: 6px;
  }

  .el-menu-item {
    padding: 0 2px;
    line-height: 14px !important;
    height: 14px !important;
  }

  .el-menu--vertical:not(.el-menu--collapse):not(.el-menu--popup-container) .el-menu-item, .el-menu--vertical:not(.el-menu--collapse):not(.el-menu--popup-container) .el-menu-item-group__title, .el-menu--vertical:not(.el-menu--collapse):not(.el-menu--popup-container) .el-sub-menu__title {
    padding-left: 2px;
  }

  .el-col.el-col-24.is-guttered {
    padding: 0 !important;
  }

  .el-button {
    margin-bottom: 10px !important;
  }

  .el-collapse-item__content {
    padding-bottom: 10px !important;
  }
}
</style>
